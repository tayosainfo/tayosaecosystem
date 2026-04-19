import React, { useState } from 'react';
import { Shield, Lock, Eye, EyeOff, ArrowRight } from 'lucide-react';
import { platformApi } from '../../lib/platformApi';

const ResetPassword: React.FC = () => {
  const [formData, setFormData] = useState({
    password: '',
    confirmPassword: ''
  });
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (formData.password !== formData.confirmPassword) {
      setError('Passwords do not match');
      return;
    }
    if (formData.password.length < 6) {
      setError('Password must be at least 6 characters long');
      return;
    }

    setIsLoading(true);
    const params = new URLSearchParams(window.location.search);
    const token = params.get('token') || '';
    if (!token) {
      setError('Missing reset token. Please request a new reset code.');
      setIsLoading(false);
      return;
    }
    try {
      await platformApi.resetPassword({
        newPassword: formData.password,
        otp: token,
      });
      setMessage('Password reset successful. Redirecting to login...');
      setTimeout(() => {
        window.location.href = '/login';
      }, 1200);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Password reset failed';
      setError(msg);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-indigo-900 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <div className="flex justify-center mb-6">
            <div className="bg-white rounded-full p-4 shadow-2xl">
              <Shield className="h-12 w-12 text-blue-600" />
            </div>
          </div>
          <h2 className="text-3xl font-bold text-white mb-2">Reset Password</h2>
          <p className="text-blue-100">Enter your new password</p>
        </div>

        <form className="space-y-6" onSubmit={handleSubmit}>
          <div className="bg-white/10 backdrop-blur-md rounded-2xl p-8 shadow-2xl border border-white/20">
            <div className="space-y-5">
              <div>
                <label htmlFor="password" className="block text-sm font-medium text-white mb-2">
                  New Password
                </label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Lock className="h-5 w-5 text-blue-300" />
                  </div>
                  <input
                    id="password"
                    name="password"
                    type={showPassword ? 'text' : 'password'}
                    autoComplete="new-password"
                    required
                    className="block w-full pl-10 pr-12 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white placeholder-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent transition-all duration-200"
                    placeholder="Enter new password"
                    value={formData.password}
                    onChange={(e) => setFormData(prev => ({ ...prev, password: e.target.value }))}
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

              <div>
                <label htmlFor="confirmPassword" className="block text-sm font-medium text-white mb-2">
                  Confirm New Password
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
                    required
                    className="block w-full pl-10 pr-12 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white placeholder-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent transition-all duration-200"
                    placeholder="Confirm new password"
                    value={formData.confirmPassword}
                    onChange={(e) => setFormData(prev => ({ ...prev, confirmPassword: e.target.value }))}
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
            </div>

            {error && (
              <div className="mt-4 bg-red-500/20 border border-red-400/50 text-red-100 px-4 py-3 rounded-xl backdrop-blur-sm">
                {error}
              </div>
            )}
            {message && (
              <div className="mt-4 bg-green-500/20 border border-green-400/50 text-green-100 px-4 py-3 rounded-xl backdrop-blur-sm">
                {message}
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
                    Resetting Password...
                  </div>
                ) : (
                  <div className="flex items-center">
                    Reset Password
                    <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
                  </div>
                )}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};

export default ResetPassword;