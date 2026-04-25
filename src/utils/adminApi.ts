/**
 * Admin API utilities
 * These functions call backend endpoints instead of querying Supabase directly
 * This avoids RLS and API permission issues
 */

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export interface AdminStatusResponse {
  isAdmin: boolean;
  role: string;
  email: string;
}

/**
 * Check if current user is admin by calling backend endpoint
 * This bypasses Supabase RLS and API restrictions
 */
export async function checkAdminStatusViaBackend(): Promise<AdminStatusResponse> {
  try {
    const authToken = sessionStorage.getItem('auth_token');
    const userStr = sessionStorage.getItem('auth_user');
    
    if (!authToken || !userStr) {
      return { isAdmin: false, role: 'user', email: '' };
    }

    const user = JSON.parse(userStr);
    const email = user.email;

    // Call backend endpoint to check admin status
    const response = await fetch(`${API_BASE_URL}/api/v1/admin/check-status`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authToken}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email }),
    });

    if (!response.ok) {
      console.error('Failed to check admin status:', response.status);
      return { isAdmin: false, role: 'user', email };
    }

    const data = await response.json();
    return {
      isAdmin: data.role === 'admin',
      role: data.role || 'user',
      email: data.email || email,
    };
  } catch (error) {
    console.error('Error checking admin status:', error);
    return { isAdmin: false, role: 'user', email: '' };
  }
}

/**
 * Direct database query for admin status (for debugging)
 * This is a fallback if the backend endpoint is not available
 */
export async function checkAdminStatusDirect(email: string): Promise<AdminStatusResponse> {
  try {
    // This would need to be implemented on the backend
    // For now, return a placeholder
    console.warn('Direct admin status check not implemented');
    return { isAdmin: false, role: 'user', email };
  } catch (error) {
    console.error('Error in direct admin check:', error);
    return { isAdmin: false, role: 'user', email };
  }
}
