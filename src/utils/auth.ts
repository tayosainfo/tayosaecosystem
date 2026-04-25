import { useEffect, useState } from 'react';
import { supabase } from '../lib/supabase';
import { useAuth } from '../hooks/useAuth';

export interface UserRole {
  isAdmin: boolean;
  role: string;
}

/**
 * Check if current user has admin role
 * Extracts role from database based on app user email
 */
export async function checkAdminStatus(): Promise<UserRole> {
  try {
    // Get the current user from app context
    const authToken = sessionStorage.getItem('auth_user');
    if (!authToken) {
      return { isAdmin: false, role: 'user' };
    }

    const user = JSON.parse(authToken);
    const userEmail = user.email;

    if (!userEmail) {
      return { isAdmin: false, role: 'user' };
    }

    // Query database for user role
    const { data, error } = await supabase
      .from('users_identity')
      .select('role')
      .eq('auth_email', userEmail)
      .single();

    if (error || !data) {
      console.error('Failed to fetch user role:', error);
      return { isAdmin: false, role: 'user' };
    }

    const userRole = data.role || 'user';
    
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
  const { user } = useAuth();

  useEffect(() => {
    // Check admin status on mount and when user changes
    if (user) {
      checkAdminStatus().then(status => {
        setAdminStatus(status);
        setLoading(false);
      });
    } else {
      setLoading(false);
    }
  }, [user]);

  return { ...adminStatus, loading };
}
