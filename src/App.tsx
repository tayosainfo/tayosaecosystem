import React from 'react';
import { AuthProvider } from './contexts/AuthContext';
import { useAuth } from './hooks/useAuth';
import Login from './pages/auth/Login';
import Register from './pages/auth/Register';
import { VerifyEmail } from './pages/auth/VerifyEmail';
import ForgotPassword from './pages/auth/ForgotPassword';
import ResetPassword from './pages/auth/ResetPassword';
import Onboarding from './pages/onboarding/Onboarding';
import OnboardingNext from './pages/onboarding/OnboardingNext';
import Home from './pages/Home';
import KycStep from './pages/KycStep';
import SaccoSetup from './pages/SaccoSetup';
import Affiliate from './pages/Affiliate';
import KibiinaSetup from './pages/KibiinaSetup';
import Admin from './pages/Admin';
import Users from './pages/admin/Users';

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

  if (window.location.pathname.startsWith('/onboarding')) {
    if (window.location.pathname.startsWith('/onboarding/next')) {
      return <OnboardingNext />;
    }
    return <Onboarding />;
  }
  if (window.location.pathname.startsWith('/kyc')) {
    return <KycStep />;
  }
  if (window.location.pathname.startsWith('/sacco')) {
    return <SaccoSetup />;
  }
  if (window.location.pathname.startsWith('/affiliate')) {
    return <Affiliate />;
  }
  if (window.location.pathname.startsWith('/admin/users')) {
    return <Users />;
  }
  if (window.location.pathname.startsWith('/admin')) {
    return <Admin />;
  }
  if (window.location.pathname.startsWith('/kibiina')) {
    return <KibiinaSetup />;
  }
  if (window.location.pathname.startsWith('/home')) {
    return <Home />;
  }

  window.location.replace('/home');
  return null;
};

function App() {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  );
}

export default App;