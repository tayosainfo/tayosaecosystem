import { supabase } from '../lib/supabase';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

/**
 * Make an authenticated admin API request
 * Automatically includes JWT token and handles token refresh
 */
export async function makeAdminRequest(
  endpoint: string,
  options: RequestInit = {}
): Promise<any> {
  // Get current session
  const { data: { session }, error } = await supabase.auth.getSession();
  
  if (error || !session) {
    throw new Error('Not authenticated');
  }

  // Make request with JWT token
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      'Authorization': `Bearer ${session.access_token}`,
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  // Handle token expiration
  if (response.status === 401) {
    // Try to refresh token
    const { data: { session: newSession } } = await supabase.auth.refreshSession();
    if (newSession) {
      // Retry with new token
      return makeAdminRequest(endpoint, options);
    }
    throw new Error('Authentication failed');
  }

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
  const { data: { session }, error } = await supabase.auth.getSession();
  
  if (error || !session) {
    throw new Error('Not authenticated');
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      'Authorization': `Bearer ${session.access_token}`,
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  if (response.status === 401) {
    const { data: { session: newSession } } = await supabase.auth.refreshSession();
    if (newSession) {
      return makeAuthenticatedRequest(endpoint, options);
    }
    throw new Error('Authentication failed');
  }

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || `Request failed: ${response.statusText}`);
  }

  return response.json();
}
