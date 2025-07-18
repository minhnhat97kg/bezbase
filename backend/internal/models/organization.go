package models

import (
	"time"

	"gorm.io/gorm"
)

type Organization struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null;size:255" validate:"required,min=1,max=255"`
	Slug      string         `json:"slug" gorm:"uniqueIndex;not null;size:100" validate:"required,min=2,max=100,alphanum"`
	Domain    string         `json:"domain" gorm:"size:255" validate:"omitempty,fqdn"`
	Settings  string         `json:"settings" gorm:"type:jsonb"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	PlanType  string         `json:"plan_type" gorm:"default:'free';size:50" validate:"oneof=free basic premium enterprise"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Users        []OrganizationUser    `json:"users,omitempty" gorm:"foreignKey:OrgID"`
	Invitations  []OrganizationInvitation `json:"invitations,omitempty" gorm:"foreignKey:OrgID"`
	Roles        []Role                `json:"roles,omitempty" gorm:"foreignKey:OrgID"`
}

type OrganizationUser struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	OrgID     uint           `json:"org_id" gorm:"not null;index"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Role      string         `json:"role" gorm:"not null;default:'member';size:50" validate:"required,oneof=owner admin member viewer"`
	IsPrimary bool           `json:"is_primary" gorm:"default:false"`
	JoinedAt  time.Time      `json:"joined_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Organization Organization `json:"organization,omitempty" gorm:"foreignKey:OrgID"`
	User         User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type OrganizationInvitation struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	OrgID      uint           `json:"org_id" gorm:"not null;index"`
	Email      string         `json:"email" gorm:"not null;size:255" validate:"required,email"`
	Role       string         `json:"role" gorm:"not null;default:'member';size:50" validate:"required,oneof=admin member viewer"`
	Token      string         `json:"token" gorm:"not null;uniqueIndex;size:255"`
	ExpiresAt  time.Time      `json:"expires_at"`
	InvitedBy  uint           `json:"invited_by" gorm:"not null"`
	AcceptedAt *time.Time     `json:"accepted_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Organization Organization `json:"organization,omitempty" gorm:"foreignKey:OrgID"`
	InvitedByUser User        `json:"invited_by_user,omitempty" gorm:"foreignKey:InvitedBy"`
}

// TableName returns the table name for Organization
func (Organization) TableName() string {
	return "organizations"
}

// TableName returns the table name for OrganizationUser
func (OrganizationUser) TableName() string {
	return "organization_users"
}

// TableName returns the table name for OrganizationInvitation
func (OrganizationInvitation) TableName() string {
	return "organization_invitations"
}

// BeforeCreate sets default values
func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	if o.PlanType == "" {
		o.PlanType = "free"
	}
	return nil
}

// BeforeCreate sets default values for OrganizationUser
func (ou *OrganizationUser) BeforeCreate(tx *gorm.DB) error {
	if ou.Role == "" {
		ou.Role = "member"
	}
	if ou.JoinedAt.IsZero() {
		ou.JoinedAt = time.Now()
	}
	return nil
}

// BeforeCreate sets default values for OrganizationInvitation
func (oi *OrganizationInvitation) BeforeCreate(tx *gorm.DB) error {
	if oi.Role == "" {
		oi.Role = "member"
	}
	if oi.ExpiresAt.IsZero() {
		// Default expiration is 7 days
		oi.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
	}
	return nil
}

// IsExpired checks if the invitation has expired
func (oi *OrganizationInvitation) IsExpired() bool {
	return time.Now().After(oi.ExpiresAt)
}

// IsAccepted checks if the invitation has been accepted
func (oi *OrganizationInvitation) IsAccepted() bool {
	return oi.AcceptedAt != nil
}

// CanManageMembers checks if the organization user can manage other members
func (ou *OrganizationUser) CanManageMembers() bool {
	return ou.Role == "owner" || ou.Role == "admin"
}

// CanInviteMembers checks if the organization user can invite new members
func (ou *OrganizationUser) CanInviteMembers() bool {
	return ou.Role == "owner" || ou.Role == "admin"
}

// CanManageSettings checks if the organization user can manage organization settings
func (ou *OrganizationUser) CanManageSettings() bool {
	return ou.Role == "owner"
}