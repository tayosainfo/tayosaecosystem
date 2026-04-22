import React, { useState } from 'react';
import { Shield, Lock, Mail, Eye, EyeOff, ArrowRight, Sparkles } from 'lucide-react';
import { useAuth } from '../../hooks/useAuth';
import { isPlatformApiError } from '../../lib/platformApi';

const BG_BLOBS = (
  <div className="absolute inset-0 overflow-hidden pointer-events-none">
    <div className="absolute -top-40 -right-40 w-80 h-80 bg-blue-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-pulse" />
    <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-indigo-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-pulse delay-1000" />
    <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-purple-400 rounded-full mix-blend-multiply filter blur-xl opacity-10 animate-pulse delay-500" />
  </div>
);

const LoginForm: React.FC = () => {
  const { login, isLoading } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    const identifier = email.trim();
    if (!identifier) { setError('Email is required'); return; }
    if (!password) { setError('Password is required'); return; }
    try {
      await login(identifier, password);
      window.location.href = '/home';
    } catch (err: unknown) {
      if (isPlatformApiError(err) && err.body.requireEmailVerification) {
        const em = String(err.body.email || identifier);
        window.location.href = `/verify?email=${encodeURIComponent(em)}`;
        return;
      }
      setError(err instanceof Error ? err.message : 'Invalid email or password');
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-indigo-900 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8 relative overflow-hidden">
      {BG_BLOBS}
      <div className="max-w-md w-full space-y-8 relative z-10">
        {/* Header */}
        <div className="text-center">
          <div className="flex justify-center mb-6">
            <div className="relative">
              <div className="absolute inset-0 bg-gradient-to-r from-blue-400 to-purple-500 rounded-full blur-lg opacity-75 animate-pulse" />
              <div className="relative bg-white rounded-full p-4 shadow-2xl">
                <Shield className="h-12 w-12 text-blue-600" />
              </div>
            </div>
          </div>
          <h2 className="text-4xl font-bold text-white mb-2">Welcome Back</h2>
          <p className="text-blue-100 text-lg">Sign in to your secure banking platform</p>
          <div className="flex items-center justify-center mt-4 space-x-2">
            <Sparkles className="h-4 w-4 text-yellow-400" />
            <span className="text-sm text-blue-200">Secure • Fast • Reliable</span>
            <Sparkles className="h-4 w-4 text-yellow-400" />
          </div>
        </div>

        {/* Form */}
        <form className="space-y-6" onSubmit={handleSubmit}>
          <div className="bg-white/10 backdrop-blur-md rounded-2xl p-8 shadow-2xl border border-white/20 space-y-5">
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-white mb-2">
                Email Address
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <Mail className="h-5 w-5 text-blue-300" />
                </div>
                <input
                  id="email"
                  type="email"
                  autoComplete="email"
                  required
                  className="block w-full pl-10 pr-3 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white placeholder-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent transition-all duration-200"
                  placeholder="you@example.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                />
              </div>
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
                  type={showPassword ? 'text' : 'password'}
                  autoComplete="current-password"
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
                  {showPassword
                    ? <EyeOff className="h-5 w-5 text-blue-300 hover:text-white transition-colors" />
                    : <Eye className="h-5 w-5 text-blue-300 hover:text-white transition-colors" />}
                </button>
              </div>
            </div>

            <div className="text-right">
              <button
                type="button"
                onClick={() => (window.location.href = '/forgot-password')}
                className="text-sm text-blue-200 hover:text-white transition-colors"
              >
                Forgot password?
              </button>
            </div>

            {error && (
              <div className="bg-red-500/20 border border-red-400/50 text-red-100 px-4 py-3 rounded-xl backdrop-blur-sm text-sm">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={isLoading}
              className="group relative w-full flex justify-center items-center py-3 px-4 rounded-xl text-blue-900 bg-gradient-to-r from-white to-blue-50 hover:from-blue-50 hover:to-white font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300 transform hover:scale-105 hover:shadow-2xl"
            >
              {isLoading
                ? <><div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-900 mr-2" />Signing In...</>
                : <><span>Sign In Securely</span><ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" /></>}
            </button>
          </div>
        </form>

        <div className="text-center space-y-2">
          <p className="text-blue-200 text-sm">
            Don't have an account?{' '}
            <button
              type="button"
              onClick={() => (window.location.href = '/register')}
              className="text-white font-semibold hover:underline"
            >
              Create one
            </button>
          </p>
          <p className="text-xs text-blue-400">🔒 Bank-grade Security &nbsp;•&nbsp; 🌍 Available 24/7 &nbsp;•&nbsp; 📱 Multi-platform</p>
        </div>
      </div>
    </div>
  );
};

export default LoginForm;
