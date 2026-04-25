import { useEffect, useState } from 'react';
import { useAuth } from '../hooks/useAuth';
import { checkAdminStatusViaBackend } from './adminApi';

export interface UserRole {
  isAdmin: boolean;
  role: string;
}

/**
 * Check if current user has admin role
 * Calls backend endpoint to avoid Supabase RLS/API issues
 */
export async function checkAdminStatus(): Promise<UserRole> {
  try {
    const result = await checkAdminStatusViaBackend();
    return {
      isAdmin: result.isAdmin,
      role: result.role
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
    const checkStatus = async () => {
      setLoading(true);
      try {
        if (user) {
          console.log('useAdminStatus: Checking admin status for user:', user.email);
          const status = await checkAdminStatus();
          console.log('useAdminStatus: Admin status result:', status);
          setAdminStatus(status);
        } else {
          console.log('useAdminStatus: No user found');
          setAdminStatus({ isAdmin: false, role: 'user' });
        }
      } catch (error) {
        console.error('useAdminStatus: Error checking admin status:', error);
        setAdminStatus({ isAdmin: false, role: 'user' });
      } finally {
        setLoading(false);
      }
    };

    checkStatus();
  }, [user]);

  return { ...adminStatus, loading };
}
