/**
 * Admin API utilities
 * These functions check admin status from the app's stored user data
 */

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export interface AdminStatusResponse {
  isAdmin: boolean;
  role: string;
  email: string;
}

/**
 * Check if current user is admin by calling backend endpoint
 * Falls back to checking database directly if backend endpoint fails
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

    // Try calling backend endpoint first
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/admin/check-status`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${authToken}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      });

      if (response.ok) {
        const data = await response.json();
        return {
          isAdmin: data.role === 'admin',
          role: data.role || 'user',
          email: data.email || email,
        };
      }
    } catch (error) {
      console.warn('Backend endpoint not available, trying direct database query:', error);
    }

    // Fallback: Query database directly via Supabase
    return await checkAdminStatusDirect(email);
  } catch (error) {
    console.error('Error checking admin status:', error);
    return { isAdmin: false, role: 'user', email: '' };
  }
}

/**
 * Direct database query for admin status
 * This queries Supabase directly as a fallback
 */
export async function checkAdminStatusDirect(email: string): Promise<AdminStatusResponse> {
  try {
    // Import Supabase client
    const { supabase } = await import('../lib/supabase');
    
    // Query database for user role
    const { data, error } = await supabase
      .from('users_identity')
      .select('role, auth_email')
      .eq('auth_email', email)
      .single();

    if (error) {
      console.error('Failed to query database:', error);
      return { isAdmin: false, role: 'user', email };
    }

    if (!data) {
      return { isAdmin: false, role: 'user', email };
    }

    return {
      isAdmin: data.role === 'admin',
      role: data.role || 'user',
      email: data.auth_email || email,
    };
  } catch (error) {
    console.error('Error in direct admin check:', error);
    return { isAdmin: false, role: 'user', email };
  }
}
