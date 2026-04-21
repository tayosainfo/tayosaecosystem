import React, { useState } from 'react';
import { Shield, Lock, Mail, Eye, EyeOff, User as UserIcon, Phone, ArrowRight, Sparkles } from 'lucide-react';
import { useAuth } from '../../hooks/useAuth';
import { isPlatformApiError, platformApi, RegisterPendingResponse } from '../../lib/platformApi';

const useInsForgeWeb = () => Boolean(String(import.meta.env.VITE_INSFORGE_BASE_URL || '').trim());

const hasSession = (v: unknown): v is { session: { accessToken: string }; user: { id: string } } => {
  if (typeof v !== 'object' || v === null) {
    return false;
  }
  const o = v as any;
  return Boolean(o.session && typeof o.session.accessToken === 'string' && o.user && typeof o.user.id === 'string');
};

const getErrorMessage = (err: unknown, fallback: string): string => {
  if (err instanceof Error && err.message) {
    return err.message;
  }
  return fallback;
};

const LoginForm: React.FC = () => {
  const [isSignUp, setIsSignUp] = useState(false);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [phone, setPhone] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [error, setError] = useState('');
  const { login, applySession, isLoading } = useAuth();
  const insForgeWeb = useInsForgeWeb();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (isSignUp) {
      if (!phone.trim()) {
        setError('Phone number is required for account creation');
        return;
      }
      if (phone.length < 10) {
        setError('Please enter a valid phone number');
        return;
      }
      if (password !== confirmPassword) {
        setError('Passwords do not match');
        return;
      }
      if (password.length < 6) {
        setError('Password must be at least 6 characters long');
        return;
      }
      if (insForgeWeb && !email.trim()) {
        setError('Email is required when using InsForge so you can receive verification codes');
        return;
      }

      try {
        if (!phone.trim()) {
          setError('Please enter a valid phone number or email');
          return;
        }
        const fullName = `${firstName} ${lastName}`.trim();
        const res = await platformApi.register({
          fullName,
          phone,
          email,
          password,
          nationality: 'UG',
          termsAccepted: true,
          privacyAccepted: true,
          termsVersion: 'v1',
          privacyVersion: 'v1',
        });
        const pending = res as RegisterPendingResponse;
        if (pending.pendingLocalProfile && pending.requireEmailVerification) {
          try {
            sessionStorage.setItem(
              'tayosa_pending_profile',
              JSON.stringify({ fullName, phone, nationality: 'UG' }),
            );
          } catch {
            /* ignore quota / private mode */
          }
          const em = pending.email || email.trim();
          window.location.href = `/verify?email=${encodeURIComponent(em)}`;
          return;
        }
        if (hasSession(res)) {
          applySession(res as any);
          window.location.href = '/onboarding';
          return;
        }
        await login(phone || email, password);
        window.location.href = '/onboarding';

      } catch (err: unknown) {
        setError(getErrorMessage(err, 'Error occurred during sign up.'));
      }
      return;
    }

    const loginIdentifier = (email || phone).trim();
    try {
      // Allow login with either phone or email
      if (!loginIdentifier) {
        setError('Please enter your phone number or email');
        return;
      }
      await login(loginIdentifier, password);
      window.location.href = '/home';
    } catch (err: unknown) {
      if (isPlatformApiError(err) && err.body.requireEmailVerification) {
        const em = String(err.body.email || (loginIdentifier.includes('@') ? loginIdentifier : email.trim()));
        if (em.includes('@')) {
          window.location.href = `/verify?email=${encodeURIComponent(em)}`;
        } else {
          window.location.href = '/verify';
        }
        return;
      }
      setError(getErrorMessage(err, 'Login failed (401). Check your identifier and password.'));
    }
  };


  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-indigo-900 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8 relative overflow-hidden">
      {/* Background decorative elements */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-blue-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-pulse"></div>
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-indigo-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-pulse delay-1000"></div>
        <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-purple-400 rounded-full mix-blend-multiply filter blur-xl opacity-10 animate-pulse delay-500"></div>
      </div>

      <div className="max-w-md w-full space-y-8 relative z-10">
        {/* Header */}
        <div className="text-center">
          <div className="flex justify-center mb-6">
            <div className="relative">
              <div className="absolute inset-0 bg-gradient-to-r from-blue-400 to-purple-500 rounded-full blur-lg opacity-75 animate-pulse"></div>
              <div className="relative bg-white rounded-full p-4 shadow-2xl">
                <Shield className="h-12 w-12 text-blue-600" />
              </div>
            </div>
          </div>
          <h2 className="text-4xl font-bold text-white mb-2">
            {isSignUp ? 'Join TAYOSA' : 'Welcome Back'}
          </h2>
          <p className="text-blue-100 text-lg">
            {isSignUp ? 'Create your secure banking account' : 'Sign in to your secure banking platform'}
          </p>
          <div className="flex items-center justify-center mt-4 space-x-2">
            <Sparkles className="h-4 w-4 text-yellow-400" />
            <span className="text-sm text-blue-200">Secure • Fast • Reliable</span>
            <Sparkles className="h-4 w-4 text-yellow-400" />
          </div>
        </div>

        {/* Auth Toggle */}
        <div className="flex bg-white/10 backdrop-blur-sm rounded-xl p-1 mb-8">
          <button
            type="button"
            onClick={() => setIsSignUp(false)}
            className={`flex-1 py-3 px-4 rounded-lg text-sm font-medium transition-all duration-300 ${!isSignUp
              ? 'bg-white text-blue-900 shadow-lg transform scale-105'
              : 'text-white hover:bg-white/10'
              }`}
          >
            Sign In
          </button>
          <button
            type="button"
            onClick={() => setIsSignUp(true)}
            className={`flex-1 py-3 px-4 rounded-lg text-sm font-medium transition-all duration-300 ${isSignUp
              ? 'bg-white text-blue-900 shadow-lg transform scale-105'
              : 'text-white hover:bg-white/10'
              }`}
          >
            Sign Up
          </button>
        </div>

        {/* Form */}
        <form className="space-y-6" onSubmit={handleSubmit}>
          <div className="bg-white/10 backdrop-blur-md rounded-2xl p-8 shadow-2xl border border-white/20">
            <div className="space-y-5">
              {isSignUp && (
                <>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label htmlFor="firstName" className="block text-sm font-medium text-white mb-2">
                        First Name
                      </label>
                      <div className="relative">
                        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                          <UserIcon className="h-5 w-5 text-blue-300" />
                        </div>
                        <input
                          id="firstName"
                          name="firstName"
                          type="text"
                          required={isSignUp}
                          className="block w-full pl-10 pr-3 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white placeholder-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent transition-all duration-200"
                          placeholder="John"
                          value={firstName}
                          onChange={(e) => setFirstName(e.target.value)}
                        />
                      </div>
                    </div>
                    <div>
                      <label htmlFor="lastName" className="block text-sm font-medium text-white mb-2">
                        Last Name
                      </label>
                      <div className="relative">
                        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                          <UserIcon className="h-5 w-5 text-blue-300" />
                        </div>
                        <input
                          id="lastName"
                          name="lastName"
                          type="text"
                          required={isSignUp}
                          className="block w-full pl-10 pr-3 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white placeholder-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent transition-all duration-200"
                          placeholder="Doe"
                          value={lastName}
                          onChange={(e) => setLastName(e.target.value)}
                        />
                      </div>
                    </div>
                  </div>

                  <div>
                    <label htmlFor="phone" className="block text-sm font-medium text-white mb-2">
                      Phone Number *
                    </label>
                    <div className="relative">
                      <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                        <Phone className="h-5 w-5 text-blue-300" />
                      </div>
                      <input
                        id="phone"
                        name="phone"
                        type="tel"
                        required
                        className="block w-full pl-10 pr-3 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white placeholder-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent transition-all duration-200"
                        placeholder="+256 700 123 456"
                        value={phone}
                        onChange={(e) => setPhone(e.target.value)}
                      />
                    </div>
                    <p className="text-xs text-blue-200 mt-1">Required for account creation and login</p>
                  </div>
                </>
              )}

              <div>
                <label htmlFor="email" className="block text-sm font-medium text-white mb-2">
                  {isSignUp
                    ? insForgeWeb
                      ? 'Email Address'
                      : 'Email Address (Optional)'
                    : 'Email or Phone Number'}
                </label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    {isSignUp ? <Mail className="h-5 w-5 text-blue-300" /> : <Phone className="h-5 w-5 text-blue-300" />}
                  </div>
                  <input
                    id="email"
                    name="email"
                    type={isSignUp ? "email" : "text"}
                    autoComplete="email"
                    required={!isSignUp || insForgeWeb}
                    className="block w-full pl-10 pr-3 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white placeholder-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent transition-all duration-200"
                    placeholder={
                      isSignUp
                        ? insForgeWeb
                          ? 'you@example.com'
                          : 'Enter your email (optional)'
                        : 'Enter phone number or email'
                    }
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                </div>
                {isSignUp && (
                  <p className="text-xs text-blue-200 mt-1">
                    {insForgeWeb
                      ? 'InsForge sends a verification code to this address'
                      : 'Email is optional but recommended for account recovery'}
                  </p>
                )}
              </div>

              <div>
                <label htmlFor="password" className="block text-sm font-medium text-white mb-2">
                  Password
                </label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Lock className="h-5 w-5 text-blue-300" />
                  </div>
                  <input
                    id="password"
                    name="password"
                    type={showPassword ? 'text' : 'password'}
                    autoComplete={isSignUp ? 'new-password' : 'current-password'}
                    required
                    className="block w-full pl-10 pr-12 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white placeholder-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent transition-all duration-200"
                    placeholder="Enter your password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                  />
                  <button
                    type="button"
                    className="absolute inset-y-0 right-0 pr-3 flex items-center"
                    onClick={() => setShowPassword(!showPassword)}
                  >
                    {showPassword ? (
                      <EyeOff className="h-5 w-5 text-blue-300 hover:text-white transition-colors" />
                    ) : (
                      <Eye className="h-5 w-5 text-blue-300 hover:text-white transition-colors" />
                    )}
                  </button>
                </div>
              </div>

              {isSignUp && (
                <div>
                  <label htmlFor="confirmPassword" className="block text-sm font-medium text-white mb-2">
                    Confirm Password
                  </label>
                  <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                      <Lock className="h-5 w-5 text-blue-300" />
                    </div>
                    <input
                      id="confirmPassword"
                      name="confirmPassword"
                      type={showConfirmPassword ? 'text' : 'password'}
                      autoComplete="new-password"
                      required={isSignUp}
                      className="block w-full pl-10 pr-12 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white placeholder-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent transition-all duration-200"
                      placeholder="Confirm your password"
                      value={confirmPassword}
                      onChange={(e) => setConfirmPassword(e.target.value)}
                    />
                    <button
                      type="button"
                      className="absolute inset-y-0 right-0 pr-3 flex items-center"
                      onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    >
                      {showConfirmPassword ? (
                        <EyeOff className="h-5 w-5 text-blue-300 hover:text-white transition-colors" />
                      ) : (
                        <Eye className="h-5 w-5 text-blue-300 hover:text-white transition-colors" />
                      )}
                    </button>
                  </div>
                </div>
              )}
            </div>

            {error && (
              <div className="mt-4 bg-red-500/20 border border-red-400/50 text-red-100 px-4 py-3 rounded-xl backdrop-blur-sm">
                {error}
              </div>
            )}

            <div className="mt-6">
              <button
                type="submit"
                disabled={isLoading}
                className="group relative w-full flex justify-center items-center py-3 px-4 border border-transparent text-sm font-medium rounded-xl text-blue-900 bg-gradient-to-r from-white to-blue-50 hover:from-blue-50 hover:to-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300 transform hover:scale-105 hover:shadow-2xl"
              >
                {isLoading ? (
                  <div className="flex items-center">
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-900 mr-2"></div>
                    {isSignUp ? 'Creating Account...' : 'Signing In...'}
                  </div>
                ) : (
                  <div className="flex items-center">
                    {isSignUp ? 'Create Account' : 'Sign In Securely'}
                    <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
                  </div>
                )}
              </button>
            </div>
          </div>
        </form>



        {/* Footer */}
        <div className="text-center">
          <p className="text-xs text-blue-300">
            By continuing, you agree to TAYOSA's Terms of Service and Privacy Policy
          </p>
          <p className="text-xs text-blue-400 mt-2">
            📱 Phone number required • 📧 Email optional • 🔒 Secure banking
          </p>
          <div className="mt-4 flex items-center justify-center space-x-4 text-xs text-blue-400">
            <span>🔒 Bank-grade Security</span>
            <span>•</span>
            <span>🌍 Available 24/7</span>
            <span>•</span>
            <span>📱 Multi-platform</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginForm;