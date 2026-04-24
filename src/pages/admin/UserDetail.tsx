import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeft, User as UserIcon, Mail, Phone, Calendar, Shield, Activity } from 'lucide-react';
import { makeAdminRequest } from '../../utils/api';

interface UserData {
  user: {
    user_id: string;
    full_name: string;
    auth_email: string;
    phone_e164: string;
    contact_email: string;
    role: string;
    status: string;
    created_at: string;
    last_login: string | null;
    role_assigned_at: string | null;
    role_assigned_by: string | null;
    date_of_birth: string | null;
    nationality: string;
  };
  kyc: {
    status: string;
    submitted_at: string | null;
    reviewed_at: string | null;
    reviewed_by: string | null;
  };
  onboarding: {
    phase: number;
    last_updated_at: string;
  };
  sacco: {
    status: string;
  };
  kibiina: {
    action: string;
  };
}

interface ActivityLog {
  action: string;
  details: string;
  timestamp: string;
  ip_address?: string;
  device_info?: string;
}

const UserDetail: React.FC = () => {
  const { userId } = useParams<{ userId: string }>();
  const navigate = useNavigate();
  const [userData, setUserData] = useState<UserData | null>(null);
  const [activity, setActivity] = useState<ActivityLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Modal states
  const [showStatusModal, setShowStatusModal] = useState(false);
  const [showRoleModal, setShowRoleModal] = useState(false);
  const [showPasswordResetModal, setShowPasswordResetModal] = useState(false);
  
  // Status change form
  const [newStatus, setNewStatus] = useState('');
  const [statusReason, setStatusReason] = useState('');
  
  // Role change form
  const [newRole, setNewRole] = useState('');
  const [showRoleConfirm, setShowRoleConfirm] = useState(false);

  useEffect(() => {
    if (userId) {
      fetchUserData();
      fetchUserActivity();
    }
  }, [userId]);

  const fetchUserData = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await makeAdminRequest(
        `/api/v1/admin/users/${userId}`,
        { method: 'GET' }
      );
      setUserData(response as UserData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load user data');
      console.error('Failed to fetch user data:', err);
    } finally {
      setLoading(false);
    }
  };

  const fetchUserActivity = async () => {
    try {
      const response = await makeAdminRequest(
        `/api/v1/admin/users/${userId}/activity`,
        { method: 'GET' }
      );
      setActivity((response as { activity: ActivityLog[] }).activity || []);
    } catch (err) {
      console.error('Failed to fetch user activity:', err);
    }
  };

  const handleStatusChange = () => {
    setNewStatus(userData?.user.status || 'active');
    setStatusReason('');
    setShowStatusModal(true);
  };

  const confirmStatusChange = async () => {
    if (!statusReason.trim()) {
      alert('Please provide a reason for the status change');
      return;
    }

    try {
      await makeAdminRequest(
        `/api/v1/admin/users/${userId}/status`,
        {
          method: 'PATCH',
          body: JSON.stringify({
            status: newStatus,
            reason: statusReason,
          }),
        }
      );
      alert('User status updated successfully');
      setShowStatusModal(false);
      fetchUserData();
      fetchUserActivity();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to update user status');
    }
  };

  const handleRoleChange = () => {
    setNewRole(userData?.user.role || 'user');
    setShowRoleConfirm(false);
    setShowRoleModal(true);
  };

  const confirmRoleChange = async () => {
    // Show confirmation dialog for admin role assignment
    if (newRole === 'admin' && !showRoleConfirm) {
      setShowRoleConfirm(true);
      return;
    }

    try {
      await makeAdminRequest(
        `/api/v1/admin/users/${userId}/role`,
        {
          method: 'PATCH',
          body: JSON.stringify({
            role: newRole,
          }),
        }
      );
      alert('User role updated successfully');
      setShowRoleModal(false);
      setShowRoleConfirm(false);
      fetchUserData();
      fetchUserActivity();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to update user role');
    }
  };

  const handlePasswordReset = () => {
    setShowPasswordResetModal(true);
  };

  const confirmPasswordReset = async () => {
    try {
      await makeAdminRequest(
        `/api/v1/admin/users/${userId}/reset-password`,
        { method: 'POST' }
      );
      alert('Password reset email sent successfully');
      setShowPasswordResetModal(false);
      fetchUserActivity();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to send password reset email');
    }
  };

  const getStatusBadge = (status: string) => {
    const colors = {
      active: 'bg-green-100 text-green-800',
      suspended: 'bg-yellow-100 text-yellow-800',
      deactivated: 'bg-red-100 text-red-800',
    };
    return colors[status as keyof typeof colors] || 'bg-gray-100 text-gray-800';
  };

  const getKycBadge = (status: string) => {
    const colors = {
      approved: 'bg-green-100 text-green-800',
      pending: 'bg-yellow-100 text-yellow-800',
      rejected: 'bg-red-100 text-red-800',
      not_started: 'bg-gray-100 text-gray-800',
    };
    return colors[status as keyof typeof colors] || 'bg-gray-100 text-gray-800';
  };

  const getRoleBadge = (role: string) => {
    return role === 'admin' ? 'bg-purple-100 text-purple-800' : 'bg-blue-100 text-blue-800';
  };

  const formatDate = (dateString: string | null) => {
    if (!dateString) return 'Never';
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-sm text-gray-600">Loading user details...</p>
        </div>
      </div>
    );
  }

  if (error || !userData) {
    return (
      <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto">
          <button
            onClick={() => navigate('/admin/users')}
            className="mb-4 inline-flex items-center text-sm text-gray-600 hover:text-gray-900"
          >
            <ArrowLeft className="h-4 w-4 mr-1" />
            Back to Users
          </button>
          <div className="bg-red-50 border border-red-200 rounded-lg p-4">
            <p className="text-sm text-red-800">{error || 'User not found'}</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto">
        {/* Back Button */}
        <button
          onClick={() => navigate('/admin/users')}
          className="mb-4 inline-flex items-center text-sm text-gray-600 hover:text-gray-900"
        >
          <ArrowLeft className="h-4 w-4 mr-1" />
          Back to Users
        </button>

        {/* User Header */}
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <div className="flex items-start justify-between">
            <div className="flex items-center">
              <div className="flex-shrink-0 h-16 w-16 bg-blue-100 rounded-full flex items-center justify-center">
                <span className="text-blue-600 font-bold text-2xl">
                  {userData.user.full_name.charAt(0).toUpperCase()}
                </span>
              </div>
              <div className="ml-4">
                <h1 className="text-2xl font-bold text-gray-900">{userData.user.full_name}</h1>
                <p className="text-sm text-gray-500">{userData.user.user_id}</p>
                <div className="mt-2 flex gap-2">
                  <span className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${getRoleBadge(userData.user.role)}`}>
                    {userData.user.role}
                  </span>
                  <span className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusBadge(userData.user.status)}`}>
                    {userData.user.status}
                  </span>
                  <span className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${getKycBadge(userData.kyc.status)}`}>
                    KYC: {userData.kyc.status.replace('_', ' ')}
                  </span>
                </div>
              </div>
            </div>
            <div className="flex gap-2">
              <button
                onClick={handleStatusChange}
                className="px-4 py-2 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700 text-sm font-medium"
              >
                Change Status
              </button>
              <button
                onClick={handleRoleChange}
                className="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 text-sm font-medium"
              >
                Change Role
              </button>
              <button
                onClick={handlePasswordReset}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 text-sm font-medium"
              >
                Reset Password
              </button>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* User Information */}
          <div className="lg:col-span-2 space-y-6">
            {/* Contact Information */}
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">Contact Information</h2>
              <div className="space-y-3">
                <div className="flex items-center">
                  <Mail className="h-5 w-5 text-gray-400 mr-3" />
                  <div>
                    <p className="text-sm font-medium text-gray-500">Email</p>
                    <p className="text-sm text-gray-900">{userData.user.auth_email}</p>
                    {userData.user.contact_email && userData.user.contact_email !== userData.user.auth_email && (
                      <p className="text-sm text-gray-600">Alt: {userData.user.contact_email}</p>
                    )}
                  </div>
                </div>
                <div className="flex items-center">
                  <Phone className="h-5 w-5 text-gray-400 mr-3" />
                  <div>
                    <p className="text-sm font-medium text-gray-500">Phone</p>
                    <p className="text-sm text-gray-900">{userData.user.phone_e164}</p>
                  </div>
                </div>
              </div>
            </div>

            {/* Account Details */}
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">Account Details</h2>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-sm font-medium text-gray-500">Joined</p>
                  <p className="text-sm text-gray-900">{formatDate(userData.user.created_at)}</p>
                </div>
                <div>
                  <p className="text-sm font-medium text-gray-500">Last Login</p>
                  <p className="text-sm text-gray-900">{formatDate(userData.user.last_login)}</p>
                </div>
                <div>
                  <p className="text-sm font-medium text-gray-500">Onboarding Phase</p>
                  <p className="text-sm text-gray-900">Phase {userData.onboarding.phase}</p>
                </div>
                <div>
                  <p className="text-sm font-medium text-gray-500">SACCO Status</p>
                  <p className="text-sm text-gray-900">{userData.sacco.status || 'Not enrolled'}</p>
                </div>
                {userData.user.date_of_birth && (
                  <div>
                    <p className="text-sm font-medium text-gray-500">Date of Birth</p>
                    <p className="text-sm text-gray-900">{userData.user.date_of_birth}</p>
                  </div>
                )}
                {userData.user.nationality && (
                  <div>
                    <p className="text-sm font-medium text-gray-500">Nationality</p>
                    <p className="text-sm text-gray-900">{userData.user.nationality}</p>
                  </div>
                )}
              </div>
            </div>

            {/* KYC Information */}
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">KYC Information</h2>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-sm font-medium text-gray-500">Status</p>
                  <p className="text-sm text-gray-900">{userData.kyc.status}</p>
                </div>
                <div>
                  <p className="text-sm font-medium text-gray-500">Submitted</p>
                  <p className="text-sm text-gray-900">{formatDate(userData.kyc.submitted_at)}</p>
                </div>
                <div>
                  <p className="text-sm font-medium text-gray-500">Reviewed</p>
                  <p className="text-sm text-gray-900">{formatDate(userData.kyc.reviewed_at)}</p>
                </div>
                <div>
                  <p className="text-sm font-medium text-gray-500">Reviewed By</p>
                  <p className="text-sm text-gray-900">{userData.kyc.reviewed_by || 'N/A'}</p>
                </div>
              </div>
            </div>
          </div>

          {/* Activity Log */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                <Activity className="h-5 w-5 mr-2" />
                Recent Activity
              </h2>
              <div className="space-y-4 max-h-96 overflow-y-auto">
                {activity.length === 0 ? (
                  <p className="text-sm text-gray-500">No recent activity</p>
                ) : (
                  activity.map((log, index) => (
                    <div key={index} className="border-l-2 border-blue-500 pl-4 pb-4">
                      <p className="text-sm font-medium text-gray-900">{log.action}</p>
                      <p className="text-xs text-gray-600">{log.details}</p>
                      <p className="text-xs text-gray-500 mt-1">{formatDate(log.timestamp)}</p>
                      {log.ip_address && (
                        <p className="text-xs text-gray-500">IP: {log.ip_address}</p>
                      )}
                    </div>
                  ))
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Password Reset Confirmation Modal */}
        {showPasswordResetModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Confirm Password Reset</h3>
              <p className="text-sm text-gray-600 mb-6">
                Are you sure you want to send a password reset email to {userData.user.auth_email}?
              </p>
              <div className="flex gap-3 justify-end">
                <button
                  onClick={() => setShowPasswordResetModal(false)}
                  className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
                >
                  Cancel
                </button>
                <button
                  onClick={confirmPasswordReset}
                  className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
                >
                  Send Reset Email
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Status Change Modal */}
        {showStatusModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Change User Status</h3>
              <div className="space-y-4">
                <div>
                  <label htmlFor="status" className="block text-sm font-medium text-gray-700 mb-1">
                    New Status
                  </label>
                  <select
                    id="status"
                    value={newStatus}
                    onChange={(e) => setNewStatus(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  >
                    <option value="active">Active</option>
                    <option value="suspended">Suspended</option>
                    <option value="deactivated">Deactivated</option>
                  </select>
                </div>
                <div>
                  <label htmlFor="reason" className="block text-sm font-medium text-gray-700 mb-1">
                    Reason <span className="text-red-500">*</span>
                  </label>
                  <textarea
                    id="reason"
                    value={statusReason}
                    onChange={(e) => setStatusReason(e.target.value)}
                    rows={3}
                    placeholder="Provide a reason for this status change..."
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>
              </div>
              <div className="flex gap-3 justify-end mt-6">
                <button
                  onClick={() => setShowStatusModal(false)}
                  className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
                >
                  Cancel
                </button>
                <button
                  onClick={confirmStatusChange}
                  className="px-4 py-2 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700"
                >
                  Update Status
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Role Change Modal */}
        {showRoleModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Change User Role</h3>
              {!showRoleConfirm ? (
                <>
                  <div className="space-y-4">
                    <div>
                      <label htmlFor="role" className="block text-sm font-medium text-gray-700 mb-1">
                        New Role
                      </label>
                      <select
                        id="role"
                        value={newRole}
                        onChange={(e) => setNewRole(e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      >
                        <option value="user">User</option>
                        <option value="admin">Admin</option>
                      </select>
                    </div>
                    {newRole === 'admin' && (
                      <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-3">
                        <p className="text-sm text-yellow-800">
                          ⚠️ Admin role grants full access to all system features and user data.
                        </p>
                      </div>
                    )}
                  </div>
                  <div className="flex gap-3 justify-end mt-6">
                    <button
                      onClick={() => setShowRoleModal(false)}
                      className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
                    >
                      Cancel
                    </button>
                    <button
                      onClick={confirmRoleChange}
                      className="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700"
                    >
                      {newRole === 'admin' ? 'Continue' : 'Update Role'}
                    </button>
                  </div>
                </>
              ) : (
                <>
                  <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-4">
                    <p className="text-sm font-semibold text-red-900 mb-2">Confirm Admin Role Assignment</p>
                    <p className="text-sm text-red-800">
                      You are about to grant admin privileges to {userData?.user.full_name}. This will give them:
                    </p>
                    <ul className="mt-2 text-sm text-red-800 list-disc list-inside space-y-1">
                      <li>Full access to all user data</li>
                      <li>Ability to manage other users</li>
                      <li>Access to admin-only features</li>
                    </ul>
                    <p className="mt-2 text-sm font-medium text-red-900">
                      Are you absolutely sure?
                    </p>
                  </div>
                  <div className="flex gap-3 justify-end">
                    <button
                      onClick={() => {
                        setShowRoleConfirm(false);
                        setShowRoleModal(false);
                      }}
                      className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
                    >
                      Cancel
                    </button>
                    <button
                      onClick={confirmRoleChange}
                      className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
                    >
                      Yes, Grant Admin Access
                    </button>
                  </div>
                </>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default UserDetail;
