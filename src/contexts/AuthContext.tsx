import React, { createContext, useState, useEffect, useCallback } from 'react';
import { User } from '../types';
import { platformApi, SessionLoginPayload } from '../lib/platformApi';
import { supabase } from '../lib/supabase';

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  login: (identifier: string, password: string) => Promise<void>;
  applySession: (response: SessionLoginPayload) => void;
  logout: () => Promise<void>;
  isLoading: boolean;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const storedUser = sessionStorage.getItem('auth_user');
    if (storedUser) {
      try {
        setUser(JSON.parse(storedUser) as User);
      } catch (error) {
        console.error('Failed to parse cached auth user', error);
      }
    }
    setIsLoading(false);
  }, []);

  // Auto-refresh Supabase token before it expires (every 50 minutes)
  useEffect(() => {
    const token = sessionStorage.getItem('auth_token');
    if (!token) return;

    const refreshInterval = setInterval(async () => {
      console.log('Auto-refreshing Supabase token...');
      try {
        const { data, error } = await supabase.auth.refreshSession();
        if (error) throw error;
        
        if (data.session?.access_token) {
          sessionStorage.setItem('auth_token', data.session.access_token);
          console.log('Token refreshed successfully');
        }
      } catch (error) {
        console.error('Failed to refresh token:', error);
        // If refresh fails, log out
        await logout();
      }
    }, 50 * 60 * 1000); // 50 minutes

    return () => clearInterval(refreshInterval);
  }, [user]);

  // Handle 401 errors globally
  useEffect(() => {
    const handleUnauthorized = (event: Event) => {
      const customEvent = event as CustomEvent;
      if (customEvent.detail?.status === 401) {
        console.warn('401 Unauthorized detected, logging out...');
        logout();
      }
    };

    window.addEventListener('auth:unauthorized', handleUnauthorized);
    return () => window.removeEventListener('auth:unauthorized', handleUnauthorized);
  }, []);

  const applySession = useCallback((response: SessionLoginPayload) => {
    const [firstName, ...rest] = response.user.fullName.split(' ');
    const nextUser: User = {
      id: response.user.id,
      email: response.user.contactEmail,
      firstName: firstName || 'Member',
      lastName: rest.join(' '),
      phone: response.user.phoneE164,
      role: 'customer',
      isActive: true,
      createdAt: new Date(),
    };
    sessionStorage.setItem('auth_token', response.session.accessToken);
    sessionStorage.setItem('auth_user', JSON.stringify(nextUser));
    setUser(nextUser);
  }, []);

  const login = async (identifier: string, password: string) => {
    setIsLoading(true);
    try {
      const response = await platformApi.login({ identifier, password });
      applySession(response);
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async () => {
    setIsLoading(true);
    try {
      // Sign out from Supabase
      await supabase.auth.signOut();
      
      // Clear local storage
      sessionStorage.removeItem('auth_token');
      sessionStorage.removeItem('auth_user');
      setUser(null);
      
      // Redirect to login
      window.location.href = '/login';
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <AuthContext.Provider value={{
      user,
      isAuthenticated: !!user,
      login,
      applySession,
      logout,
      isLoading
    }}>
      {children}
    </AuthContext.Provider>
  );
};