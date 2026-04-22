import React, { useEffect, useState } from 'react';
import { Shield, Key, Mail, ArrowRight, AlertCircle } from 'lucide-react';
import { platformApi, SessionLoginPayload } from '../../lib/platformApi';
import { useAuth } from '../../hooks/useAuth';

function isSessionPayload(x: unknown): x is SessionLoginPayload {
  if (typeof x !== 'object' || x === null) return false;
  const o = x as any;
  return typeof o.session?.accessToken === 'string' && typeof o.user?.fullName === 'string';
}

const BG_BLOBS = (
  <div className="absolute inset-0 overflow-hidden pointer-events-none">
    <div className="absolute -top-40 -right-40 w-80 h-80 bg-blue-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-pulse" />
    <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-indigo-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-pulse delay-1000" />
  </div>
);

export const VerifyEmail: React.FC = () => {
  const { applySession } = useAuth();
  const [otp, setOtp] = useState('');
  const [email, setEmail] = useState('');
  const [error, setError] = useState('');
  const [info, setInfo] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const q = new URLSearchParams(window.location.search).get('email');
    if (q) setEmail(q);
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(''); setInfo('');
    if (!email.trim() || !email.includes('@')) { setError('Please enter your registered email'); return; }
    if (otp.length !== 6) { setError('Please enter the 6-digit code'); return; }

    setIsLoading(true);
    try {
      // Pull pending profile data that register stored in sessionStorage
      let extras: Record<string, string> = {};
      try {
        const raw = sessionStorage.getItem('tayosa_pending_profile');
        if (raw) {
          const p = JSON.parse(raw) as { fullName?: string; phone?: string; nationality?: string };
          if (p.fullName) extras.fullName = p.fullName;
          if (p.phone) extras.phone = p.phone;
          if (p.nationality) extras.nationality = p.nationality;
        }
      } catch { /* ignore */ }

      const data = await platformApi.verifyEmail({ email, otp, ...extras });
      if (isSessionPayload(data)) {
        try { sessionStorage.removeItem('tayosa_pending_profile'); } catch { /* ignore */ }
        applySession(data);
        window.location.href = '/onboarding';
        return;
      }
      window.location.href = '/';
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Verification failed. Please check the code and try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const resend = async () => {
    if (!email.trim()) { setError('Please enter your email first'); return; }
    setError(''); setInfo('');
    try {
      await platformApi.resendVerificationEmail({ email });
      setInfo('Verification code resent — check your inbox.');
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to resend. Please try again.');
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
          <h2 className="text-4xl font-bold text-white mb-2">Verify Your Email</h2>
          <p className="text-blue-100 text-lg">
            {email ? <>We sent a 6-digit code to <span className="font-semibold text-white">{email}</span></> : 'Enter the code we sent to your email'}
          </p>
        </div>

        <form className="space-y-6" onSubmit={handleSubmit}>
          <div className="bg-white/10 backdrop-blur-md rounded-2xl p-8 shadow-2xl border border-white/20 space-y-5">
            {/* Email field */}
            <div>
              <label className="block text-sm font-medium text-white mb-2">Email Address</label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <Mail className="h-5 w-5 text-blue-300" />
                </div>
                <input
                  type="email" required
                  className="block w-full pl-10 pr-3 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white placeholder-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-400 transition-all duration-200"
                  placeholder="you@example.com"
                  value={email}
                  onChange={e => setEmail(e.target.value)}
                />
              </div>
            </div>

            {/* OTP field */}
            <div>
              <label className="block text-sm font-medium text-white mb-2">6-Digit Verification Code</label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <Key className="h-5 w-5 text-blue-300" />
                </div>
                <input
                  type="text" required maxLength={6} inputMode="numeric"
                  className="block w-full pl-10 pr-3 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white text-center font-mono text-2xl tracking-[0.5em] placeholder-blue-300 focus:outline-none focus:ring-2 focus:ring-blue-400 transition-all duration-200"
                  placeholder="------"
                  value={otp}
                  onChange={e => setOtp(e.target.value.replace(/\D/g, ''))}
                />
              </div>
            </div>

            {error && (
              <div className="flex items-center bg-red-500/20 border border-red-400/50 text-red-100 px-4 py-3 rounded-xl backdrop-blur-sm">
                <AlertCircle className="h-5 w-5 mr-2 shrink-0" />
                <span className="text-sm">{error}</span>
              </div>
            )}
            {info && (
              <div className="bg-green-500/20 border border-green-400/50 text-green-100 px-4 py-3 rounded-xl text-sm">
                {info}
              </div>
            )}

            <button
              type="submit"
              disabled={isLoading || otp.length !== 6}
              className="group w-full flex justify-center items-center py-3 px-4 rounded-xl text-blue-900 bg-gradient-to-r from-white to-blue-50 hover:from-blue-50 hover:to-white font-medium transition-all duration-300 transform hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading
                ? <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-blue-900" />
                : <><span>Verify &amp; Continue</span><ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" /></>}
            </button>

            <div className="flex flex-col items-center gap-3 pt-2">
              <button type="button" onClick={resend} className="text-blue-200 hover:text-white text-sm transition-colors">
                Resend verification code
              </button>
              <button type="button" onClick={() => (window.location.href = '/')} className="text-blue-300 hover:text-white text-sm transition-colors">
                ← Back to sign in
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};
