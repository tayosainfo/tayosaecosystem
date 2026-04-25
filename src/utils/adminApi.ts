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
 * Direct database query for admin status via Supabase REST API
 */
export async function checkAdminStatusDirect(email: string): Promise<AdminStatusResponse> {
  try {
    const supabaseUrl = import.meta.env.VITE_SUPABASE_URL;
    const supabaseAnonKey = import.meta.env.VITE_SUPABASE_ANON_KEY;

    if (!supabaseUrl || !supabaseAnonKey) {
      console.error('Missing Supabase environment variables');
      return { isAdmin: false, role: 'user', email };
    }

    console.log('Querying Supabase REST API for admin status:', email);
    
    // Use REST API directly with proper headers
    const url = new URL(`${supabaseUrl}/rest/v1/users_identity`);
    url.searchParams.append('select', 'role,auth_email');
    url.searchParams.append('auth_email', `eq.${email}`);

    const response = await fetch(url.toString(), {
      method: 'GET',
      headers: {
        'apikey': supabaseAnonKey,
        'Authorization': `Bearer ${supabaseAnonKey}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      }
    });

    console.log('Supabase REST API response status:', response.status);

    if (!response.ok) {
      const errorText = await response.text();
      console.error('Supabase REST API error:', response.status, errorText);
      return { isAdmin: false, role: 'user', email };
    }

    const data = await response.json();
    console.log('Supabase REST API response data:', data);

    if (!Array.isArray(data) || data.length === 0) {
      console.log('User not found in database');
      return { isAdmin: false, role: 'user', email };
    }

    const user = data[0];
    console.log('User found with role:', user.role);
    
    return {
      isAdmin: user.role === 'admin',
      role: user.role || 'user',
      email: user.auth_email || email,
    };
  } catch (error) {
    console.error('Error in direct admin check:', error);
    return { isAdmin: false, role: 'user', email };
  }
}
