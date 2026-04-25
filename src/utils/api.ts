const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

/**
 * Make an authenticated admin API request
 * Automatically includes JWT token from sessionStorage
 */
export async function makeAdminRequest(
  endpoint: string,
  options: RequestInit = {}
): Promise<any> {
  // Get JWT token from sessionStorage (custom auth, not Supabase Auth)
  const token = sessionStorage.getItem('auth_token');
  
  if (!token) {
    throw new Error('Not authenticated');
  }

  // Make request with JWT token
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  // Handle authorization failure
  if (response.status === 403) {
    throw new Error('Insufficient permissions');
  }

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || `Request failed: ${response.statusText}`);
  }

  return response.json();
}

/**
 * Make a regular authenticated API request (non-admin)
 * Similar to makeAdminRequest but for regular user endpoints
 */
export async function makeAuthenticatedRequest(
  endpoint: string,
  options: RequestInit = {}
): Promise<any> {
  // Get JWT token from sessionStorage (custom auth, not Supabase Auth)
  const token = sessionStorage.getItem('auth_token');
  
  if (!token) {
    throw new Error('Not authenticated');
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || `Request failed: ${response.statusText}`);
  }

  return response.json();
}
