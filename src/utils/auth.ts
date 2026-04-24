import { useEffect, useState } from 'react';
import { supabase } from '../lib/supabase';

export interface UserRole {
  isAdmin: boolean;
  role: string;
}

/**
 * Check if current user has admin role
 * Extracts role from JWT token claims
 */
export async function checkAdminStatus(): Promise<UserRole> {
  try {
    const { data: { user }, error } = await supabase.auth.getUser();
    
    if (error || !user) {
      return { isAdmin: false, role: 'user' };
    }

    // Extract role from app_metadata (set by Supabase custom claims hook)
    const userRole = user.app_metadata?.user_role || 'user';
    
    return {
      isAdmin: userRole === 'admin',
      role: userRole
    };
  } catch (error) {
    console.error('Failed to check admin status:', error);
    return { isAdmin: false, role: 'user' };
  }
}

/**
 * React hook for checking admin status
 * Returns admin status with loading state
 */
export function useAdminStatus() {
  const [adminStatus, setAdminStatus] = useState<UserRole>({
    isAdmin: false,
    role: 'user'
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Check admin status on mount
    checkAdminStatus().then(status => {
      setAdminStatus(status);
      setLoading(false);
    });

    // Listen for auth state changes and update admin status
    const { data: { subscription } } = supabase.auth.onAuthStateChange(() => {
      checkAdminStatus().then(setAdminStatus);
    });

    return () => subscription.unsubscribe();
  }, []);

  return { ...adminStatus, loading };
}
