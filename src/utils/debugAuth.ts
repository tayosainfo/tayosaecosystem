import { supabase } from '../lib/supabase';

/**
 * Debug utility to inspect JWT token claims
 * Use this in browser console to see what's in the token
 */
export async function debugTokenClaims() {
  try {
    const { data: { session }, error } = await supabase.auth.getSession();
    
    if (error || !session) {
      console.log('❌ No active session');
      return;
    }

    // Decode JWT token (just the payload, not verifying signature)
    const token = session.access_token;
    const payload = JSON.parse(atob(token.split('.')[1]));
    
    console.log('🔍 JWT Token Claims:', payload);
    console.log('📧 Email:', payload.email);
    console.log('👤 User ID:', payload.sub);
    console.log('🎭 Role from app_metadata:', payload.app_metadata?.user_role || 'NOT SET');
    console.log('⏰ Token expires:', new Date(payload.exp * 1000).toLocaleString());
    
    return payload;
  } catch (error) {
    console.error('Failed to decode token:', error);
  }
}

/**
 * Force refresh the JWT token to get updated claims
 */
export async function forceTokenRefresh() {
  try {
    console.log('🔄 Refreshing token...');
    const { data, error } = await supabase.auth.refreshSession();
    
    if (error) {
      console.error('❌ Token refresh failed:', error);
      return false;
    }
    
    if (data.session) {
      console.log('✅ Token refreshed successfully');
      
      // Decode and show new claims
      const token = data.session.access_token;
      const payload = JSON.parse(atob(token.split('.')[1]));
      console.log('🎭 New role:', payload.app_metadata?.user_role || 'NOT SET');
      
      // Reload the page to apply new token
      console.log('🔄 Reloading page to apply new token...');
      window.location.reload();
      
      return true;
    }
    
    return false;
  } catch (error) {
    console.error('Failed to refresh token:', error);
    return false;
  }
}

// Make functions available in browser console
if (typeof window !== 'undefined') {
  (window as any).debugTokenClaims = debugTokenClaims;
  (window as any).forceTokenRefresh = forceTokenRefresh;
}
