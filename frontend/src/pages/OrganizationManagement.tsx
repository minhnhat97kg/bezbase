import React, { useState, useEffect } from 'react';
import { useOrganization } from '../context/OrganizationContext';
import { organizationService, OrganizationUser, OrganizationInvitation, InviteUserRequest } from '../services/organizationService';
import TabLayout from '../components/common/TabLayout';
import Table from '../components/common/Table';

const OrganizationManagement: React.FC = () => {
  const {
    currentOrganization,
    currentOrgUser,
    userOrganizations,
    isLoadingOrganizations,
    canManageOrganization,
    canManageMembers,
    canInviteMembers,
    refreshCurrentOrganization,
  } = useOrganization();

  const [activeTab, setActiveTab] = useState('overview');
  const [members, setMembers] = useState<OrganizationUser[]>([]);
  const [invitations, setInvitations] = useState<OrganizationInvitation[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [showInviteModal, setShowInviteModal] = useState(false);

  useEffect(() => {
    if (currentOrganization) {
      loadOrganizationData();
    }
  }, [currentOrganization]);

  const loadOrganizationData = async () => {
    if (!currentOrganization) return;

    setIsLoading(true);
    try {
      await refreshCurrentOrganization();
      if (currentOrganization.users) {
        setMembers(currentOrganization.users);
      }
      if (currentOrganization.invitations) {
        setInvitations(currentOrganization.invitations);
      }
    } catch (error) {
      console.error('Failed to load organization data:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleInviteUser = async (inviteData: InviteUserRequest) => {
    if (!currentOrganization) return;

    try {
      await organizationService.inviteUser(currentOrganization.id, inviteData);
      setShowInviteModal(false);
      await loadOrganizationData();
    } catch (error) {
      console.error('Failed to invite user:', error);
    }
  };

  const handleRemoveMember = async (userId: number) => {
    if (!currentOrganization || !window.confirm('Are you sure you want to remove this member?')) return;

    try {
      await organizationService.removeUser(currentOrganization.id, userId);
      await loadOrganizationData();
    } catch (error) {
      console.error('Failed to remove member:', error);
    }
  };

  const handleUpdateMemberRole = async (userId: number, newRole: string) => {
    if (!currentOrganization) return;

    try {
      await organizationService.updateUserRole(currentOrganization.id, userId, newRole);
      await loadOrganizationData();
    } catch (error) {
      console.error('Failed to update member role:', error);
    }
  };

  if (isLoadingOrganizations) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto"></div>
          <h3 className="mt-4 text-lg font-medium text-gray-900">Loading Organizations...</h3>
          <p className="mt-2 text-sm text-gray-500">Please wait while we load your organizations.</p>
        </div>
      </div>
    );
  }

  if (!currentOrganization) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-gray-100">
            <svg className="h-6 w-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
            </svg>
          </div>
          <h3 className="mt-4 text-lg font-medium text-gray-900">No Organization Selected</h3>
          <p className="mt-2 text-sm text-gray-500">
            {userOrganizations.length === 0 
              ? "You don't belong to any organizations yet. Create your first organization to get started."
              : "Please select an organization from the header to manage."
            }
          </p>
          {userOrganizations.length === 0 && (
            <div className="mt-6">
              <button
                onClick={() => window.location.href = '/organizations/create'}
                className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                <svg className="-ml-1 mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                </svg>
                Create Organization
              </button>
            </div>
          )}
        </div>
      </div>
    );
  }

  const tabs = [
    { id: 'overview', label: 'Overview', icon: '🏢' },
    { id: 'members', label: 'Members', icon: '👥' },
    { id: 'settings', label: 'Settings', icon: '⚙️', disabled: !canManageOrganization() },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-white shadow">
        <div className="px-4 py-5 sm:p-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">{currentOrganization.name}</h1>
              <div className="mt-2 flex items-center space-x-4">
                <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${organizationService.getPlanColor(currentOrganization.plan_type)}`}>
                  {organizationService.formatPlanType(currentOrganization.plan_type)}
                </span>
                <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${organizationService.getRoleColor(currentOrgUser?.role || '')}`}>
                  {currentOrgUser ? organizationService.formatRole(currentOrgUser.role) : 'Unknown'}
                </span>
                {currentOrganization.domain && (
                  <span className="text-sm text-gray-500">
                    {currentOrganization.domain}
                  </span>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>

      <TabLayout
        tabs={tabs}
        activeTab={activeTab}
        onTabChange={setActiveTab}
      >
        {/* Overview Tab */}
        {activeTab === 'overview' && (
          <div className="space-y-6">
            <OrganizationOverview organization={currentOrganization} />
          </div>
        )}

        {/* Members Tab */}
        {activeTab === 'members' && (
          <div className="space-y-6">
            <MembersManagement
              members={members}
              invitations={invitations}
              canInvite={canInviteMembers()}
              canManage={canManageMembers()}
              onInvite={() => setShowInviteModal(true)}
              onRemoveMember={handleRemoveMember}
              onUpdateRole={handleUpdateMemberRole}
            />
          </div>
        )}

        {/* Settings Tab */}
        {activeTab === 'settings' && canManageOrganization() && (
          <div className="space-y-6">
            <OrganizationSettings
              organization={currentOrganization}
              onUpdate={loadOrganizationData}
            />
          </div>
        )}
      </TabLayout>

      {/* Invite User Modal */}
      {showInviteModal && (
        <InviteUserModal
          onInvite={handleInviteUser}
          onClose={() => setShowInviteModal(false)}
        />
      )}
    </div>
  );
};

// Organization Overview Component
const OrganizationOverview: React.FC<{ organization: any }> = ({ organization }) => {
  const stats = [
    { name: 'Total Members', value: organization.users?.length || 0, icon: '👥' },
    { name: 'Pending Invitations', value: organization.invitations?.filter((inv: any) => !inv.accepted_at).length || 0, icon: '📤' },
    { name: 'Organization Roles', value: organization.roles?.filter((role: any) => role.org_id).length || 0, icon: '🎭' },
    { name: 'Plan Type', value: organizationService.formatPlanType(organization.plan_type), icon: '📊' },
  ];

  return (
    <div className="bg-white shadow rounded-lg">
      <div className="px-4 py-5 sm:p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">Organization Overview</h3>
        <dl className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
          {stats.map((stat) => (
            <div key={stat.name} className="bg-gray-50 px-4 py-5 rounded-lg">
              <dt className="text-sm font-medium text-gray-500 truncate flex items-center">
                <span className="mr-2">{stat.icon}</span>
                {stat.name}
              </dt>
              <dd className="mt-1 text-3xl font-semibold text-gray-900">{stat.value}</dd>
            </div>
          ))}
        </dl>
      </div>
    </div>
  );
};

// Members Management Component
interface MembersManagementProps {
  members: OrganizationUser[];
  invitations: OrganizationInvitation[];
  canInvite: boolean;
  canManage: boolean;
  onInvite: () => void;
  onRemoveMember: (userId: number) => void;
  onUpdateRole: (userId: number, role: string) => void;
}

const MembersManagement: React.FC<MembersManagementProps> = ({
  members,
  invitations,
  canInvite,
  canManage,
  onInvite,
  onRemoveMember,
  onUpdateRole,
}) => {
  const memberColumns = [
    {
      key: 'user',
      label: 'User',
      render: (member: OrganizationUser) => (
        <div className="flex items-center">
          <div className="flex-shrink-0 h-10 w-10">
            <div className="h-10 w-10 rounded-full bg-gray-300 flex items-center justify-center">
              <span className="text-sm font-medium text-gray-700">
                {member.user?.first_name?.charAt(0)}{member.user?.last_name?.charAt(0)}
              </span>
            </div>
          </div>
          <div className="ml-4">
            <div className="text-sm font-medium text-gray-900">
              {member.user?.first_name} {member.user?.last_name}
            </div>
            <div className="text-sm text-gray-500">{member.user?.email}</div>
          </div>
        </div>
      ),
    },
    {
      key: 'role',
      label: 'Role',
      render: (member: OrganizationUser) => (
        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${organizationService.getRoleColor(member.role)}`}>
          {organizationService.formatRole(member.role)}
        </span>
      ),
    },
    {
      key: 'joined_at',
      label: 'Joined',
      render: (member: OrganizationUser) => new Date(member.joined_at).toLocaleDateString(),
    },
    {
      key: 'actions',
      label: 'Actions',
      render: (member: OrganizationUser) => (
        canManage && member.role !== 'owner' ? (
          <div className="flex space-x-2">
            <select
              value={member.role}
              onChange={(e) => onUpdateRole(member.user_id, e.target.value)}
              className="text-xs border rounded px-2 py-1"
            >
              <option value="admin">Admin</option>
              <option value="member">Member</option>
              <option value="viewer">Viewer</option>
            </select>
            <button
              onClick={() => onRemoveMember(member.user_id)}
              className="text-red-600 hover:text-red-900 text-xs"
            >
              Remove
            </button>
          </div>
        ) : null
      ),
    },
  ];

  const invitationColumns = [
    { key: 'email', label: 'Email' },
    {
      key: 'role',
      label: 'Role',
      render: (invitation: OrganizationInvitation) => (
        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${organizationService.getRoleColor(invitation.role)}`}>
          {organizationService.formatRole(invitation.role)}
        </span>
      ),
    },
    {
      key: 'status',
      label: 'Status',
      render: (invitation: OrganizationInvitation) => {
        if (invitation.accepted_at) {
          return <span className="text-green-600 text-xs">Accepted</span>;
        }
        if (new Date(invitation.expires_at) < new Date()) {
          return <span className="text-red-600 text-xs">Expired</span>;
        }
        return <span className="text-yellow-600 text-xs">Pending</span>;
      },
    },
    {
      key: 'created_at',
      label: 'Invited',
      render: (invitation: OrganizationInvitation) => new Date(invitation.created_at).toLocaleDateString(),
    },
  ];

  return (
    <div className="space-y-6">
      {/* Members Section */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium text-gray-900">Members ({members.length})</h3>
            {canInvite && (
              <button
                onClick={onInvite}
                className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700"
              >
                Invite User
              </button>
            )}
          </div>
        </div>
        <Table
          data={members}
          columns={memberColumns}
          emptyMessage="No members found"
        />
      </div>

      {/* Invitations Section */}
      {invitations.length > 0 && (
        <div className="bg-white shadow rounded-lg">
          <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
            <h3 className="text-lg font-medium text-gray-900">Pending Invitations ({invitations.length})</h3>
          </div>
          <Table
            data={invitations}
            columns={invitationColumns}
            emptyMessage="No pending invitations"
          />
        </div>
      )}
    </div>
  );
};

// Organization Settings Component
const OrganizationSettings: React.FC<{ organization: any; onUpdate: () => void }> = ({ organization, onUpdate }) => {
  const [formData, setFormData] = useState({
    name: organization.name,
    domain: organization.domain || '',
    plan_type: organization.plan_type,
  });
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      await organizationService.updateOrganization(organization.id, formData);
      await onUpdate();
    } catch (error) {
      console.error('Failed to update organization:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="bg-white shadow rounded-lg">
      <div className="px-4 py-5 sm:p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">Organization Settings</h3>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700">Organization Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700">Domain</label>
            <input
              type="text"
              value={formData.domain}
              onChange={(e) => setFormData({ ...formData, domain: e.target.value })}
              className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
              placeholder="example.com"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700">Plan Type</label>
            <select
              value={formData.plan_type}
              onChange={(e) => setFormData({ ...formData, plan_type: e.target.value as any })}
              className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="free">Free</option>
              <option value="basic">Basic</option>
              <option value="premium">Premium</option>
              <option value="enterprise">Enterprise</option>
            </select>
          </div>

          <div className="flex justify-end">
            <button
              type="submit"
              disabled={isLoading}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50"
            >
              {isLoading ? 'Saving...' : 'Save Changes'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

// Invite User Modal Component
interface InviteUserModalProps {
  onInvite: (data: InviteUserRequest) => void;
  onClose: () => void;
}

const InviteUserModal: React.FC<InviteUserModalProps> = ({ onInvite, onClose }) => {
  const [formData, setFormData] = useState<InviteUserRequest>({
    email: '',
    role: 'member',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onInvite(formData);
  };

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={onClose}></div>

        <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <form onSubmit={handleSubmit}>
            <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Invite User</h3>
              
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Email Address</label>
                  <input
                    type="email"
                    value={formData.email}
                    onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                    required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700">Role</label>
                  <select
                    value={formData.role}
                    onChange={(e) => setFormData({ ...formData, role: e.target.value as any })}
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                  >
                    <option value="admin">Admin</option>
                    <option value="member">Member</option>
                    <option value="viewer">Viewer</option>
                  </select>
                </div>
              </div>
            </div>

            <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
              <button
                type="submit"
                className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm"
              >
                Send Invitation
              </button>
              <button
                type="button"
                onClick={onClose}
                className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default OrganizationManagement;