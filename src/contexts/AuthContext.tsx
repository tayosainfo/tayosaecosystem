import React, { createContext, useState, useEffect, useCallback } from 'react';
import { User } from '../types';
import { platformApi, SessionLoginPayload } from '../lib/platformApi';

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
      sessionStorage.removeItem('auth_token');
      sessionStorage.removeItem('auth_user');
      setUser(null);
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