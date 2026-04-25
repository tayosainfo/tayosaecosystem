/**
 * Admin API utilities
 * Checks admin status by querying Supabase directly
 */

export interface AdminStatusResponse {
  isAdmin: boolean;
  role: string;
  email: string;
}

/**
 * Check if current user is admin by querying Supabase directly
 * (No backend endpoint - direct Supabase query)
 */
export async function checkAdminStatusViaBackend(): Promise<AdminStatusResponse> {
  try {
    const userStr = sessionStorage.getItem('auth_user');
    
    if (!userStr) {
      console.log('No auth_user in sessionStorage');
      return { isAdmin: false, role: 'user', email: '' };
    }

    const user = JSON.parse(userStr);
    const email = user.email;

    if (!email) {
      console.log('No email found in auth_user');
      return { isAdmin: false, role: 'user', email: '' };
    }

    console.log('Checking admin status for email:', email);
    
    // Query Supabase directly
    const result = await checkAdminStatusDirect(email);
    console.log('Admin status result:', result);
    return result;
  } catch (error) {
    console.error('Error checking admin status:', error);
    return { isAdmin: false, role: 'user', email: '' };
  }
}

/**
 * Direct database query for admin status via Supabase
 */
export async function checkAdminStatusDirect(email: string): Promise<AdminStatusResponse> {
  try {
    // Import Supabase client
    const { supabase } = await import('../lib/supabase');
    
    console.log('Querying Supabase for admin status:', email);
    
    // Query database for user role
    const { data, error } = await supabase
      .from('users_identity')
      .select('role, auth_email')
      .eq('auth_email', email)
      .single();

    if (error) {
      console.error('Supabase query error:', error);
      return { isAdmin: false, role: 'user', email };
    }

    if (!data) {
      console.log('User not found in database');
      return { isAdmin: false, role: 'user', email };
    }

    console.log('User found with role:', data.role);
    
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
