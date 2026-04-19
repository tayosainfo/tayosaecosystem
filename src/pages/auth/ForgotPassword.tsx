import React, { useState } from 'react';
import { Shield, Mail, ArrowRight, ArrowLeft } from 'lucide-react';
import { platformApi } from '../../lib/platformApi';

const ForgotPassword: React.FC = () => {
  const [email, setEmail] = useState('');
  const [code, setCode] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await platformApi.sendResetPasswordEmail({ email });
      setIsSubmitted(true);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to send reset code';
      setError(message);
    } finally {
      setIsLoading(false);
    }
  };

  const handleVerifyCode = async () => {
    setError('');
    setIsLoading(true);
    try {
      const data = await platformApi.exchangeResetPasswordToken({ email, code });
      window.location.href = `/reset?token=${encodeURIComponent(data.token)}`;
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Invalid or expired code';
      setError(message);
    } finally {
      setIsLoading(false);
    }
  };

  if (isSubmitted) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-indigo-900 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-md w-full text-center">
          <div className="bg-white/10 backdrop-blur-md rounded-2xl p-8 shadow-2xl border border-white/20">
            <div className="flex justify-center mb-6">
              <div className="bg-green-100 rounded-full p-4">
                <Mail className="h-12 w-12 text-green-600" />
              </div>
            </div>
            <h2 className="text-2xl font-bold text-white mb-4">Check Your Email</h2>
            <p className="text-blue-100 mb-4">
              We sent a 6-digit password reset code to {email}
            </p>
            <input
              value={code}
              onChange={(e) => setCode(e.target.value.replace(/[^0-9]/g, '').slice(0, 6))}
              className="w-full mb-4 px-4 py-3 rounded-xl bg-white/10 text-white border border-white/20"
              placeholder="Enter 6-digit code"
            />
            {error && <p className="text-red-200 text-sm mb-4">{error}</p>}
            <button
              onClick={handleVerifyCode}
              disabled={isLoading || code.length !== 6}
              className="w-full py-3 rounded-xl bg-white text-blue-900 font-semibold disabled:opacity-60 mb-3"
            >
              Verify code
            </button>
            <button
              onClick={() => setIsSubmitted(false)}
              className="text-blue-200 hover:text-white transition-colors"
            >
              Back
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-indigo-900 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <div className="flex justify-center mb-6">
            <div className="bg-white rounded-full p-4 shadow-2xl">
              <Shield className="h-12 w-12 text-blue-600" />
            </div>
          </div>
          <h2 className="text-3xl font-bold text-white mb-2">Forgot Password?</h2>
          <p className="text-blue-100">Enter your email to reset your password</p>
        </div>

        <form className="space-y-6" onSubmit={handleSubmit}>
          <div className="bg-white/10 backdrop-blur-md rounded-2xl p-8 shadow-2xl border border-white/20">
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
                  name="email"
                  type="email"
                  autoComplete="email"
                  required
                  className="block w-full pl-10 pr-3 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white placeholder-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent transition-all duration-200"
                  placeholder="Enter your email address"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                />
              </div>
            </div>

            <div className="mt-6">
              <button
                type="submit"
                disabled={isLoading}
                className="group relative w-full flex justify-center items-center py-3 px-4 border border-transparent text-sm font-medium rounded-xl text-blue-900 bg-gradient-to-r from-white to-blue-50 hover:from-blue-50 hover:to-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300 transform hover:scale-105 hover:shadow-2xl"
              >
                {isLoading ? (
                  <div className="flex items-center">
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-900 mr-2"></div>
                    Sending Reset Link...
                  </div>
                ) : (
                  <div className="flex items-center">
                    Send Reset Link
                    <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
                  </div>
                )}
              </button>
            </div>
            {error && <p className="text-red-200 text-sm mt-3">{error}</p>}

            <div className="mt-4 text-center">
              <button
                type="button"
                className="text-blue-200 hover:text-white transition-colors text-sm"
              >
                <ArrowLeft className="h-4 w-4 inline mr-1" />
                Back to login
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};

export default ForgotPassword;