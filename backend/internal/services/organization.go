package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"bezbase/internal/dto"
	"bezbase/internal/models"
	"bezbase/internal/pkg/contextx"
	"bezbase/internal/repository"

	"gorm.io/gorm"
)

type OrganizationService struct {
	orgRepo        repository.OrganizationRepository
	orgUserRepo    repository.OrganizationUserRepository
	invitationRepo repository.OrganizationInvitationRepository
	userRepo       repository.UserRepository
	rbacService    *RBACService
	emailService   *EmailService
	db             *gorm.DB
}

type CreateOrganizationRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Slug     string `json:"slug" validate:"required,min=2,max=100,alphanum"`
	Domain   string `json:"domain" validate:"omitempty,fqdn"`
	PlanType string `json:"plan_type" validate:"omitempty,oneof=free basic premium enterprise"`
}

type InviteUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin member viewer"`
}

type UpdateOrganizationRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=1,max=255"`
	Domain   *string `json:"domain" validate:"omitempty,fqdn"`
	PlanType *string `json:"plan_type" validate:"omitempty,oneof=free basic premium enterprise"`
	Settings *string `json:"settings"`
}

func NewOrganizationService(
	orgRepo repository.OrganizationRepository,
	orgUserRepo repository.OrganizationUserRepository,
	invitationRepo repository.OrganizationInvitationRepository,
	userRepo repository.UserRepository,
	rbacService *RBACService,
	emailService *EmailService,
	db *gorm.DB,
) *OrganizationService {
	return &OrganizationService{
		orgRepo:        orgRepo,
		orgUserRepo:    orgUserRepo,
		invitationRepo: invitationRepo,
		userRepo:       userRepo,
		rbacService:    rbacService,
		emailService:   emailService,
		db:             db,
	}
}

// CreateOrganization creates a new organization with the creator as owner
func (s *OrganizationService) CreateOrganization(ctx contextx.Contextx, userID uint, req CreateOrganizationRequest) (*models.Organization, error) {
	// Check if slug is already taken
	existingOrg, err := s.orgRepo.GetBySlug(ctx, req.Slug)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingOrg != nil {
		return nil, errors.New("organization slug already exists")
	}

	// Set default plan type
	if req.PlanType == "" {
		req.PlanType = "free"
	}

	org := &models.Organization{
		Name:     req.Name,
		Slug:     req.Slug,
		Domain:   req.Domain,
		PlanType: req.PlanType,
		IsActive: true,
	}

	// Create organization in transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Create organization
		if err := s.orgRepo.Create(ctx, org); err != nil {
			return err
		}

		// Add creator as owner
		orgUser := &models.OrganizationUser{
			OrgID:     org.ID,
			UserID:    userID,
			Role:      "owner",
			IsPrimary: true,
			JoinedAt:  time.Now(),
		}

		if err := s.orgUserRepo.Create(ctx, orgUser); err != nil {
			return err
		}

		// Update user's current organization if they don't have one
		user, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			return err
		}

		if user.CurrentOrgID == nil {
			user.CurrentOrgID = &org.ID
			if err := s.userRepo.Update(ctx, user); err != nil {
				return err
			}
		}

		// Create organization-specific admin role
		orgAdminRole := &models.Role{
			Name:           fmt.Sprintf("org_admin_%d", org.ID),
			DisplayName:    "Organization Administrator",
			Description:    fmt.Sprintf("Administrator for %s", org.Name),
			IsSystem:       false,
			IsActive:       true,
			OrgID:          &org.ID,
			HierarchyLevel: 1,
		}

		// Find global admin role to inherit from
		globalAdmin, err := s.rbacService.GetRoleByName(ctx, "admin")
		if err == nil {
			orgAdminRole.ParentRoleID = &globalAdmin.ID
		}

		roleReq := dto.CreateRoleRequest{
			Name:        orgAdminRole.Name,
			DisplayName: orgAdminRole.DisplayName,
			Description: orgAdminRole.Description,
		}
		createdRole, err := s.rbacService.CreateRole(ctx, roleReq)
		if err != nil {
			return err
		}

		// Assign the organization admin role to the owner
		if err := s.rbacService.AssignRoleToUser(userID, fmt.Sprint(createdRole.ID)); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return org, nil
}

// GetOrganization retrieves an organization by ID
func (s *OrganizationService) GetOrganization(ctx contextx.Contextx, orgID uint) (*models.Organization, error) {
	return s.orgRepo.GetByIDWithRelations(ctx, orgID)
}

// GetUserOrganizations retrieves all organizations for a user
func (s *OrganizationService) GetUserOrganizations(ctx contextx.Contextx, userID uint) ([]models.OrganizationUser, error) {
	return s.orgUserRepo.GetByUserID(ctx, userID)
}

// UpdateOrganization updates organization details
func (s *OrganizationService) UpdateOrganization(ctx contextx.Contextx, orgID uint, userID uint, req UpdateOrganizationRequest) (*models.Organization, error) {
	// Check if user has permission to update organization
	if !s.canManageOrganization(orgID, userID) {
		return nil, errors.New("insufficient permissions to update organization")
	}

	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		org.Name = *req.Name
	}
	if req.Domain != nil {
		org.Domain = *req.Domain
	}
	if req.PlanType != nil {
		org.PlanType = *req.PlanType
	}
	if req.Settings != nil {
		org.Settings = *req.Settings
	}

	if err := s.orgRepo.Update(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

// InviteUser invites a user to join the organization
func (s *OrganizationService) InviteUser(ctx contextx.Contextx, orgID uint, inviterID uint, req InviteUserRequest) (*models.OrganizationInvitation, error) {
	// Check if inviter has permission
	if !s.canInviteMembers(orgID, inviterID) {
		return nil, errors.New("insufficient permissions to invite users")
	}

	// Check if user is already a member
	existingMember, err := s.orgUserRepo.GetByOrgIDAndEmail(ctx, orgID, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingMember != nil {
		return nil, errors.New("user is already a member of this organization")
	}

	// Check if there's already a pending invitation
	existingInvitation, err := s.invitationRepo.GetByOrgIDAndEmail(ctx, orgID, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingInvitation != nil && !existingInvitation.IsExpired() && !existingInvitation.IsAccepted() {
		return nil, errors.New("invitation already sent and pending")
	}

	// Generate invitation token
	token, err := s.generateInvitationToken()
	if err != nil {
		return nil, err
	}

	invitation := &models.OrganizationInvitation{
		OrgID:     orgID,
		Email:     req.Email,
		Role:      req.Role,
		Token:     token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		InvitedBy: inviterID,
	}

	if err := s.invitationRepo.Create(ctx, invitation); err != nil {
		return nil, err
	}

	// Send invitation email
	org, _ := s.orgRepo.GetByID(ctx, orgID)
	inviter, _ := s.userRepo.GetByID(ctx, inviterID)

	if s.emailService != nil && org != nil && inviter != nil {
		// TODO: Implement SendOrganizationInvitation method in EmailService
		// For now, just log that we would send an invitation email
		log.Printf("Would send organization invitation email to %s for organization %s", req.Email, org.Name)
	}

	return invitation, nil
}

// AcceptInvitation accepts an organization invitation
func (s *OrganizationService) AcceptInvitation(ctx contextx.Contextx, token string, userID uint) error {
	invitation, err := s.invitationRepo.GetByToken(ctx, token)
	if err != nil {
		return errors.New("invalid invitation token")
	}

	if invitation.IsExpired() {
		return errors.New("invitation has expired")
	}

	if invitation.IsAccepted() {
		return errors.New("invitation has already been accepted")
	}

	// Check if user is already a member
	existingMember, err := s.orgUserRepo.GetByOrgIDAndUserID(ctx, invitation.OrgID, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingMember != nil {
		return errors.New("user is already a member of this organization")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Create organization membership
		orgUser := &models.OrganizationUser{
			OrgID:    invitation.OrgID,
			UserID:   userID,
			Role:     invitation.Role,
			JoinedAt: time.Now(),
		}

		if err := s.orgUserRepo.Create(ctx, orgUser); err != nil {
			return err
		}

		// Mark invitation as accepted
		now := time.Now()
		invitation.AcceptedAt = &now
		if err := s.invitationRepo.Update(ctx, invitation); err != nil {
			return err
		}

		// Update user's current organization if they don't have one
		user, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			return err
		}

		if user.CurrentOrgID == nil {
			user.CurrentOrgID = &invitation.OrgID
			if err := s.userRepo.Update(ctx, user); err != nil {
				return err
			}
		}

		return nil
	})
}

// SwitchOrganization switches user's current organization context
func (s *OrganizationService) SwitchOrganization(userID uint, orgID uint) error {
	// Check if user is a member of the organization
	ctx := contextx.NewContextx(context.Background())
	_, err := s.orgUserRepo.GetByOrgIDAndUserID(ctx, orgID, userID)
	if err != nil {
		return errors.New("user is not a member of this organization")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user.CurrentOrgID = &orgID
	return s.userRepo.Update(ctx, user)
}

// RemoveUser removes a user from an organization
func (s *OrganizationService) RemoveUser(orgID uint, targetUserID uint, actorUserID uint) error {
	ctx := contextx.Background()
	// Check if actor has permission
	if !s.canManageMembers(orgID, actorUserID) {
		return errors.New("insufficient permissions to remove users")
	}

	// Don't allow removing yourself if you're the only owner
	if targetUserID == actorUserID {
		owners, err := s.orgUserRepo.GetByOrgIDAndRole(ctx, orgID, "owner")
		if err != nil {
			return err
		}
		if len(owners) == 1 {
			return errors.New("cannot remove yourself as the only owner")
		}
	}

	return s.orgUserRepo.DeleteByOrgIDAndUserID(ctx, orgID, targetUserID)
}

// UpdateUserRole updates a user's role in an organization
func (s *OrganizationService) UpdateUserRole(orgID uint, targetUserID uint, newRole string, actorUserID uint) error {
	// Check if actor has permission
	if !s.canManageMembers(orgID, actorUserID) {
		return errors.New("insufficient permissions to update user roles")
	}
	ctx := contextx.Background()
	orgUser, err := s.orgUserRepo.GetByOrgIDAndUserID(ctx, orgID, targetUserID)
	if err != nil {
		return err
	}

	orgUser.Role = newRole
	return s.orgUserRepo.Update(ctx, orgUser)
}

// Helper methods

func (s *OrganizationService) canManageOrganization(orgID uint, userID uint) bool {
	orgUser, err := s.orgUserRepo.GetByOrgIDAndUserID(contextx.Background(), orgID, userID)
	if err != nil {
		return false
	}
	return orgUser.Role == "owner"
}

func (s *OrganizationService) canManageMembers(orgID uint, userID uint) bool {
	orgUser, err := s.orgUserRepo.GetByOrgIDAndUserID(contextx.Background(), orgID, userID)
	if err != nil {
		return false
	}
	return orgUser.CanManageMembers()
}

func (s *OrganizationService) canInviteMembers(orgID uint, userID uint) bool {
	orgUser, err := s.orgUserRepo.GetByOrgIDAndUserID(contextx.Background(), orgID, userID)
	if err != nil {
		return false
	}
	return orgUser.CanInviteMembers()
}

func (s *OrganizationService) generateInvitationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreateSlugFromName creates a URL-friendly slug from a name
func (s *OrganizationService) CreateSlugFromName(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove non-alphanumeric characters except hyphens
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	// Remove multiple consecutive hyphens
	slug = result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Ensure minimum length
	if len(slug) < 2 {
		slug = "org-" + slug
	}

	return slug
}
