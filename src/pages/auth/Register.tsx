import React, { useEffect, useState } from 'react';
import { Shield, Lock, Mail, Eye, EyeOff, User as UserIcon, Phone, ArrowRight, Sparkles } from 'lucide-react';
import { platformApi, RegisterPendingResponse } from '../../lib/platformApi';
import { useAuth } from '../../hooks/useAuth';

const BG_BLOBS = (
  <div className="absolute inset-0 overflow-hidden pointer-events-none">
    <div className="absolute -top-40 -right-40 w-80 h-80 bg-blue-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-pulse" />
    <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-indigo-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-pulse delay-1000" />
    <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-purple-400 rounded-full mix-blend-multiply filter blur-xl opacity-10 animate-pulse delay-500" />
  </div>
);

function hasSession(v: unknown): v is { session: { accessToken: string }; user: { id: string } } {
  if (typeof v !== 'object' || v === null) return false;
  const o = v as any;
  return Boolean(o.session?.accessToken && o.user?.id);
}

const Register: React.FC = () => {
  const { applySession } = useAuth();
  const [form, setForm] = useState({
    firstName: '', lastName: '', email: '', phone: '',
    password: '', confirmPassword: '', referralCode: '',
  });
  const [termsAccepted, setTermsAccepted] = useState(false);
  const [privacyAccepted, setPrivacyAccepted] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const ref = new URLSearchParams(window.location.search).get('ref');
    if (ref) setForm(p => ({ ...p, referralCode: ref }));
  }, []);

  const set = (field: string) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm(p => ({ ...p, [field]: e.target.value }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!form.email.trim()) { setError('Email address is required'); return; }
    if (!form.phone.trim()) { setError('Phone number is required'); return; }
    if (form.phone.length < 10) { setError('Please enter a valid Uganda phone number'); return; }
    if (form.password.length < 6) { setError('Password must be at least 6 characters'); return; }
    if (form.password !== form.confirmPassword) { setError('Passwords do not match'); return; }
    if (!termsAccepted || !privacyAccepted) { setError('Please accept Terms and Privacy Policy'); return; }

    setIsLoading(true);
    try {
      const fullName = `${form.firstName} ${form.lastName}`.trim();
      const res = await platformApi.register({
        fullName,
        phone: form.phone,
        email: form.email,
        password: form.password,
        nationality: 'UG',
        referralCode: form.referralCode || undefined,
        termsAccepted: true,
        privacyAccepted: true,
        termsVersion: 'v1',
        privacyVersion: 'v1',
      });

      const pending = res as RegisterPendingResponse;
      if (pending.requireEmailVerification) {
        try {
          sessionStorage.setItem('tayosa_pending_profile', JSON.stringify({ fullName, phone: form.phone, nationality: 'UG' }));
        } catch { /* ignore */ }
        const em = pending.email || form.email.trim();
        window.location.href = `/verify?email=${encodeURIComponent(em)}`;
        return;
      }
      if (hasSession(res)) {
        applySession(res as any);
        window.location.href = '/onboarding';
        return;
      }
      window.location.href = `/verify?email=${encodeURIComponent(form.email.trim())}`;
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Registration failed. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const inputClass = "block w-full pl-10 pr-3 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white placeholder-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent transition-all duration-200";

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
          <h2 className="text-4xl font-bold text-white mb-2">Join TAYOSA</h2>
          <p className="text-blue-100 text-lg">Create your secure banking account</p>
          <div className="flex items-center justify-center mt-4 space-x-2">
            <Sparkles className="h-4 w-4 text-yellow-400" />
            <span className="text-sm text-blue-200">Secure • Fast • Reliable</span>
            <Sparkles className="h-4 w-4 text-yellow-400" />
          </div>
        </div>

        <form className="space-y-6" onSubmit={handleSubmit}>
          <div className="bg-white/10 backdrop-blur-md rounded-2xl p-8 shadow-2xl border border-white/20 space-y-5">
            {/* Name row */}
            <div className="grid grid-cols-2 gap-4">
              {(['firstName', 'lastName'] as const).map((field, i) => (
                <div key={field}>
                  <label className="block text-sm font-medium text-white mb-2">
                    {i === 0 ? 'First Name' : 'Last Name'}
                  </label>
                  <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                      <UserIcon className="h-5 w-5 text-blue-300" />
                    </div>
                    <input
                      type="text" required
                      className={inputClass}
                      placeholder={i === 0 ? 'John' : 'Doe'}
                      value={form[field]}
                      onChange={set(field)}
                    />
                  </div>
                </div>
              ))}
            </div>

            {/* Email */}
            <div>
              <label className="block text-sm font-medium text-white mb-2">Email Address *</label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <Mail className="h-5 w-5 text-blue-300" />
                </div>
                <input
                  type="email" required autoComplete="email"
                  className={inputClass}
                  placeholder="you@example.com"
                  value={form.email}
                  onChange={set('email')}
                />
              </div>
              <p className="text-xs text-blue-200 mt-1">A verification code will be sent to this address</p>
            </div>

            {/* Phone */}
            <div>
              <label className="block text-sm font-medium text-white mb-2">Phone Number *</label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <Phone className="h-5 w-5 text-blue-300" />
                </div>
                <input
                  type="tel" required
                  className={inputClass}
                  placeholder="+256 700 123 456"
                  value={form.phone}
                  onChange={set('phone')}
                />
              </div>
              <p className="text-xs text-blue-200 mt-1">Uganda format: 0700123456 or +256700123456</p>
            </div>

            {/* Referral */}
            <div>
              <label className="block text-sm font-medium text-white mb-2">Referral Code (Optional)</label>
              <input
                type="text"
                className="block w-full px-3 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white placeholder-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-400"
                placeholder="TAY-XXXX1234"
                value={form.referralCode}
                onChange={set('referralCode')}
              />
            </div>

            {/* Password */}
            <div>
              <label className="block text-sm font-medium text-white mb-2">Password</label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <Lock className="h-5 w-5 text-blue-300" />
                </div>
                <input
                  type={showPassword ? 'text' : 'password'}
                  required autoComplete="new-password"
                  className="block w-full pl-10 pr-12 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white placeholder-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent transition-all duration-200"
                  placeholder="At least 6 characters"
                  value={form.password}
                  onChange={set('password')}
                />
                <button type="button" className="absolute inset-y-0 right-0 pr-3 flex items-center" onClick={() => setShowPassword(!showPassword)}>
                  {showPassword ? <EyeOff className="h-5 w-5 text-blue-300 hover:text-white" /> : <Eye className="h-5 w-5 text-blue-300 hover:text-white" />}
                </button>
              </div>
            </div>

            {/* Confirm Password */}
            <div>
              <label className="block text-sm font-medium text-white mb-2">Confirm Password</label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <Lock className="h-5 w-5 text-blue-300" />
                </div>
                <input
                  type={showConfirm ? 'text' : 'password'}
                  required autoComplete="new-password"
                  className="block w-full pl-10 pr-12 py-3 border border-white/30 rounded-xl bg-white/10 backdrop-blur-sm text-white placeholder-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent transition-all duration-200"
                  placeholder="Repeat your password"
                  value={form.confirmPassword}
                  onChange={set('confirmPassword')}
                />
                <button type="button" className="absolute inset-y-0 right-0 pr-3 flex items-center" onClick={() => setShowConfirm(!showConfirm)}>
                  {showConfirm ? <EyeOff className="h-5 w-5 text-blue-300 hover:text-white" /> : <Eye className="h-5 w-5 text-blue-300 hover:text-white" />}
                </button>
              </div>
            </div>

            {/* Consents */}
            <div className="space-y-2 text-white text-sm">
              <label className="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" checked={termsAccepted} onChange={e => setTermsAccepted(e.target.checked)} className="rounded" />
                I accept the Terms and Conditions
              </label>
              <label className="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" checked={privacyAccepted} onChange={e => setPrivacyAccepted(e.target.checked)} className="rounded" />
                I accept the Privacy Policy
              </label>
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
                ? <><div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-900 mr-2" />Creating Account...</>
                : <><span>Create Account</span><ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" /></>}
            </button>
          </div>
        </form>

        <div className="text-center space-y-2">
          <p className="text-blue-200 text-sm">
            Already have an account?{' '}
            <button type="button" onClick={() => (window.location.href = '/')} className="text-white font-semibold hover:underline">
              Sign in
            </button>
          </p>
          <p className="text-xs text-blue-400">🔒 Bank-grade Security &nbsp;•&nbsp; 🌍 Available 24/7 &nbsp;•&nbsp; 📱 Multi-platform</p>
        </div>
      </div>
    </div>
  );
};

export default Register;
