import React from 'react';
import { AuthProvider } from './contexts/AuthContext';
import { useAuth } from './hooks/useAuth';
import Login from './pages/auth/Login';
import Register from './pages/auth/Register';
import { VerifyEmail } from './pages/auth/VerifyEmail';
import ForgotPassword from './pages/auth/ForgotPassword';
import ResetPassword from './pages/auth/ResetPassword';

const AppContent: React.FC = () => {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    if (window.location.pathname.startsWith('/forgot-password')) {
      return <ForgotPassword />;
    }
    if (window.location.pathname.startsWith('/reset')) {
      return <ResetPassword />;
    }
    if (window.location.pathname.startsWith('/verify')) {
      return <VerifyEmail />;
    }
    if (window.location.pathname.startsWith('/register')) {
      return <Register />;
    }
    return <Login />;
  }

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col items-center justify-center p-4">
      <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-8 text-center max-w-md w-full">
        <h1 className="text-2xl font-bold text-gray-900 mb-4">Authentication Successful</h1>
        <p className="text-gray-600">The application is being rebuilt to use InsForge. Sacco backend and frontend components have been removed.</p>
      </div>
    </div>
  );
};

function App() {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  );
}

export default App;