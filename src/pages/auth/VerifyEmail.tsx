import React, { useEffect, useState } from 'react';
import { Shield, Key, ArrowRight, AlertCircle, User as UserIcon, Phone } from 'lucide-react';
import { platformApi, SessionLoginPayload } from '../../lib/platformApi';
import { useAuth } from '../../hooks/useAuth';

const getErrorMessage = (err: unknown, fallback: string): string => {
  if (err instanceof Error && err.message) {
    return err.message;
  }
  return fallback;
};

function isSessionLoginPayload(x: unknown): x is SessionLoginPayload {
  if (typeof x !== 'object' || x === null) {
    return false;
  }
  const o = x as SessionLoginPayload;
  return (
    typeof o.session?.accessToken === 'string' &&
    typeof o.user?.fullName === 'string' &&
    typeof o.user?.phoneE164 === 'string'
  );
}

type PendingProfile = { fullName?: string; phone?: string; nationality?: string };

export const VerifyEmail: React.FC = () => {
  const { applySession } = useAuth();
  const [otp, setOtp] = useState('');
  const [email, setEmail] = useState('');
  const [fullName, setFullName] = useState('');
  const [phone, setPhone] = useState('');
  const [nationality] = useState('UG');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const q = params.get('email');
    if (q) {
      setEmail(q);
    }
    try {
      const raw = sessionStorage.getItem('tayosa_pending_profile');
      if (raw) {
        const p = JSON.parse(raw) as PendingProfile;
        if (p.fullName) {
          setFullName(p.fullName);
        }
        if (p.phone) {
          setPhone(p.phone);
        }
      }
    } catch {
      /* ignore */
    }
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (otp.length !== 6) {
      setError('Please enter a valid 6-digit code');
      return;
    }

    if (!email.trim() || !email.includes('@')) {
      setError('Please provide your registered email');
      return;
    }

    setIsLoading(true);

    try {
      const extras: Record<string, string> = {};
      if (fullName.trim() && phone.trim()) {
        extras.fullName = fullName.trim();
        extras.phone = phone.trim();
        extras.nationality = nationality.trim() || 'UG';
      }
      const data = await platformApi.verifyEmail({
        email,
        otp,
        ...extras,
      });
      if (isSessionLoginPayload(data)) {
        try {
          sessionStorage.removeItem('tayosa_pending_profile');
        } catch {
          /* ignore */
        }
        applySession(data);
        window.location.href = '/';
        return;
      }
      window.location.href = '/login';
    } catch (err: unknown) {
      setError(getErrorMessage(err, 'Verification failed. Please double check the code.'));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-indigo-900 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8 relative overflow-hidden">
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-blue-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-pulse"></div>
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-indigo-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-pulse delay-1000"></div>
      </div>

      <div className="max-w-md w-full space-y-8 relative z-10">
        <div className="text-center">
          <div className="flex justify-center mb-6">
            <div className="relative">
              <div className="absolute inset-0 bg-gradient-to-r from-blue-400 to-purple-500 rounded-full blur-lg opacity-75 animate-pulse"></div>
              <div className="relative bg-white rounded-full p-4 shadow-2xl">
                <Shield className="h-12 w-12 text-blue-600" />
              </div>
            </div>
          </div>
          <h2 className="text-4xl font-bold text-white mb-2">Verify Registration</h2>
          <p className="text-blue-100 text-lg">We&apos;ve sent a 6-digit verification code to your email.</p>
          <p className="text-blue-200 text-sm mt-2">
            If you just signed up with InsForge, enter the same full name and phone you used on the sign-up form once,
            so we can finish your Tayosa profile after the code is accepted.
          </p>
        </div>

        <form className="space-y-6" onSubmit={handleSubmit}>
          <div className="bg-white/10 backdrop-blur-md rounded-2xl p-8 shadow-2xl border border-white/20">
            <div className="space-y-5">
              <div>
                <label className="block text-sm font-medium text-white mb-2">Confirm Email Address</label>
                <input
                  type="email"
                  required
                  className="block w-full px-4 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white focus:outline-none focus:ring-2 focus:ring-blue-400 mb-4 transition-all duration-200"
                  placeholder="Enter your registered email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-white mb-2">Full name (for first-time setup)</label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <UserIcon className="h-5 w-5 text-blue-300" />
                  </div>
                  <input
                    type="text"
                    className="block w-full pl-10 pr-4 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white focus:outline-none focus:ring-2 focus:ring-blue-400"
                    placeholder="Same as sign-up"
                    value={fullName}
                    onChange={(e) => setFullName(e.target.value)}
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-white mb-2">Phone (for first-time setup)</label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Phone className="h-5 w-5 text-blue-300" />
                  </div>
                  <input
                    type="tel"
                    className="block w-full pl-10 pr-4 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white focus:outline-none focus:ring-2 focus:ring-blue-400"
                    placeholder="+256 …"
                    value={phone}
                    onChange={(e) => setPhone(e.target.value)}
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-white mb-2">6-Digit Verification Code</label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Key className="h-5 w-5 text-blue-300" />
                  </div>
                  <input
                    type="text"
                    required
                    maxLength={6}
                    className="block w-full pl-10 pr-4 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white text-center font-mono text-xl tracking-widest placeholder-blue-300 focus:outline-none focus:ring-2 focus:ring-blue-400 transition-all duration-200"
                    placeholder="------"
                    value={otp}
                    onChange={(e) => setOtp(e.target.value.replace(/[^0-9]/g, ''))}
                  />
                </div>
              </div>
            </div>

            {error && (
              <div className="mt-4 flex items-center bg-red-500/20 text-red-100 px-4 py-3 rounded-xl backdrop-blur-sm border border-red-400/50">
                <AlertCircle className="h-5 w-5 mr-2 shrink-0" />
                <span className="text-sm">{error}</span>
              </div>
            )}

            <div className="mt-6 flex flex-col space-y-4">
              <button
                type="submit"
                disabled={isLoading}
                className="group w-full flex justify-center items-center py-3 px-4 rounded-xl text-blue-900 bg-gradient-to-r from-white to-blue-50 font-medium hover:from-blue-50 hover:to-white transition-all duration-300 transform hover:scale-105"
              >
                {isLoading ? (
                  <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-blue-900"></div>
                ) : (
                  <>
                    Verify & continue
                    <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
                  </>
                )}
              </button>

              <button
                type="button"
                onClick={async () => {
                  try {
                    await platformApi.resendVerificationEmail({ email });
                    setError('Verification code resent. Check your inbox.');
                  } catch (err: unknown) {
                    setError(getErrorMessage(err, 'Failed to resend verification code.'));
                  }
                }}
                className="text-white hover:text-blue-200 text-sm font-medium transition-colors"
              >
                Resend verification code
              </button>

              <button
                type="button"
                onClick={() => (window.location.href = '/login')}
                className="text-white hover:text-blue-200 text-sm font-medium transition-colors"
              >
                Go back to Login
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};
