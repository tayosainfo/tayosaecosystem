import React, { useEffect, useMemo, useState } from 'react';
import { CheckCircle2, Clock3, Lock, ArrowRight, Users, PiggyBank, Landmark, Wallet } from 'lucide-react';
import { useAuth } from '../hooks/useAuth';
import { platformApi } from '../lib/platformApi';
import { makeAdminRequest } from '../utils/api';

type KycStatus = 'not_started' | 'pending' | 'approved';
type SaccoStatus = 'not_started' | 'enrolled';

const Home: React.FC = () => {
  const { user } = useAuth();
  const [kyc, setKyc] = useState<KycStatus>('not_started');
  const [sacco, setSacco] = useState<SaccoStatus>('not_started');
  const [shares, setShares] = useState(0);
  const [refCode, setRefCode] = useState('AMARA-7742');

  useEffect(() => {
    const token = sessionStorage.getItem('auth_token');
    if (!token) return;
    platformApi.getMe(token).then((me) => {
      const k = String(me.kyc?.status || 'not_started') as KycStatus;
      const s = String(me.sacco?.status || 'not_started') as SaccoStatus;
      setKyc(k);
      setSacco(s);
      setShares(Number(me.shares?.balanceUnits || 0));
      if (me.referralCode) {
        setRefCode(me.referralCode);
      }
    }).catch(() => {
      // ignore transient errors on first render
    });
  }, []);

  const canTransact = useMemo(() => kyc === 'approved' && sacco === 'enrolled', [kyc, sacco]);
  const canJoinKibiina = useMemo(() => sacco === 'enrolled', [sacco]);

  const checklist = [
    { label: 'Account created', status: 'done', hint: 'Phone and email verified' },
    {
      label: 'Identity check',
      status: kyc === 'approved' ? 'done' : kyc === 'pending' ? 'pending' : 'todo',
      hint: kyc === 'approved' ? 'KYC approved' : kyc === 'pending' ? 'Submitted, under review' : 'Complete KYC to unlock savings',
    },
    {
      label: 'SACCO membership',
      status: sacco === 'enrolled' ? 'done' : 'todo',
      hint: sacco === 'enrolled' ? 'Parish SACCO active' : 'Select district and enrol',
    },
  ] as const;

  const statusPill =
    kyc === 'approved' ? 'All checks complete' : kyc === 'pending' ? 'KYC under review - usually within 24 hours' : 'Finish setup to unlock all features';

  const markKycApprovedForPreview = async () => {
    const userRaw = sessionStorage.getItem('auth_user');
    if (!userRaw) return;
    
    try {
      const parsed = JSON.parse(userRaw) as { id: string };
      await makeAdminRequest(
        `/api/v1/admin/kyc?userId=${encodeURIComponent(parsed.id)}`,
        {
          method: 'PATCH',
          body: JSON.stringify({ status: 'approved', reviewedBy: 'web-admin-preview' }),
        }
      );
      window.location.reload();
    } catch (error) {
      console.error('Failed to approve KYC:', error);
      alert('Failed to approve KYC. You may not have admin permissions.');
    }
  };

  const markSaccoEnrolledForPreview = () => {
    sessionStorage.setItem('sacco_status', 'enrolled');
    window.location.reload();
  };

  const shareReferral = async () => {
    const code = refCode;
    const link = `${window.location.origin}/register?ref=${encodeURIComponent(code)}`;
    try {
      await navigator.clipboard.writeText(`${code} ${link}`);
      window.alert('Referral code and link copied.');
    } catch {
      window.alert(`Copy this code manually: ${code}`);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-indigo-900 py-10 px-4 sm:px-6 lg:px-8">
      <div className="max-w-6xl mx-auto space-y-6">
        <div className="bg-white/10 backdrop-blur-md rounded-2xl border border-white/20 p-6 md:p-8 text-white">
          <p className="text-sm text-blue-100">{statusPill}</p>
          <div className="mt-2 flex flex-col md:flex-row md:items-end md:justify-between gap-3">
            <div>
              <h1 className="text-4xl font-bold">UGX 0</h1>
              <p className="text-blue-100">Welcome, {user?.firstName || 'Member'}</p>
            </div>
            <div className="flex gap-2">
              {kyc !== 'approved' && (
                <button onClick={markKycApprovedForPreview} className="inline-flex items-center justify-center gap-2 px-4 py-2 rounded-xl bg-white text-blue-900 font-semibold hover:bg-blue-50 transition">Approve KYC (preview)</button>
              )}
              {sacco !== 'enrolled' && (
                <button
                  onClick={markSaccoEnrolledForPreview}
                  className="inline-flex items-center justify-center gap-2 px-4 py-2 rounded-xl bg-white text-blue-900 font-semibold hover:bg-blue-50 transition"
                >
                  Enrol SACCO (preview)
                </button>
              )}
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="bg-white/10 backdrop-blur-md rounded-xl border border-white/20 p-4 text-white">
            <p className="text-xs text-blue-100">SACCO</p>
            <p className="text-xl font-bold mt-1">{sacco === 'enrolled' ? 'Parish SACCO active' : 'Not enrolled yet'}</p>
          </div>
          <div className="bg-white/10 backdrop-blur-md rounded-xl border border-white/20 p-4 text-white">
            <p className="text-xs text-blue-100">Shares</p>
            <p className="text-xl font-bold mt-1">{shares} units</p>
          </div>
          <div className="bg-white/10 backdrop-blur-md rounded-xl border border-white/20 p-4 text-white">
            <p className="text-xs text-blue-100">Kibiina</p>
            <p className="text-xl font-bold mt-1">{canJoinKibiina ? 'Nakawa A' : 'Locked'}</p>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <button disabled={!canTransact} className={`rounded-xl p-5 border border-white/20 text-left backdrop-blur-md ${canTransact ? 'bg-white text-slate-900' : 'bg-white/10 text-blue-100'}`}>
            <PiggyBank className="h-5 w-5 mb-2" />
            <p className="font-semibold">Save</p>
            {!canTransact && <p className="text-xs opacity-90">Locked until KYC and SACCO setup</p>}
          </button>
          <button disabled={!canTransact} className={`rounded-xl p-5 border border-white/20 text-left backdrop-blur-md ${canTransact ? 'bg-white text-slate-900' : 'bg-white/10 text-blue-100'}`}>
            <Landmark className="h-5 w-5 mb-2" />
            <p className="font-semibold">Borrow</p>
            {!canTransact && <p className="text-xs opacity-90">Locked until KYC and SACCO setup</p>}
          </button>
          <button disabled={!canTransact} className={`rounded-xl p-5 border border-white/20 text-left backdrop-blur-md ${canTransact ? 'bg-white text-slate-900' : 'bg-white/10 text-blue-100'}`}>
            <Wallet className="h-5 w-5 mb-2" />
            <p className="font-semibold">Withdraw</p>
            {!canTransact && <p className="text-xs opacity-90">Locked until KYC and SACCO setup</p>}
          </button>
          <button onClick={() => canJoinKibiina && (window.location.href = '/kibiina')} disabled={!canJoinKibiina} className={`rounded-xl p-5 border border-white/20 text-left backdrop-blur-md ${canJoinKibiina ? 'bg-white text-slate-900' : 'bg-white/10 text-blue-100'}`}>
            <Users className="h-5 w-5 mb-2" />
            <p className="font-semibold">Kibiina</p>
            {!canJoinKibiina && <p className="text-xs opacity-90">Join SACCO first</p>}
          </button>
        </div>

        <div className="bg-white/10 backdrop-blur-md rounded-2xl border border-white/20 p-6">
          <h2 className="text-sm font-semibold text-blue-100 mb-4">COMPLETE YOUR PROFILE</h2>
          <div className="space-y-3">
            {checklist.map((item) => (
              <div key={item.label} className="flex items-start justify-between bg-white rounded-xl p-4">
                <div className="flex items-start gap-2 text-slate-900">
                  {item.status === 'done' ? <CheckCircle2 className="h-5 w-5 text-green-600 mt-0.5" /> : item.status === 'pending' ? <Clock3 className="h-5 w-5 text-amber-500 mt-0.5" /> : <Lock className="h-5 w-5 text-slate-400 mt-0.5" />}
                  <div>
                    <p className="font-medium text-sm">{item.label}</p>
                    <p className="text-xs text-slate-500">{item.hint}</p>
                  </div>
                </div>
                {item.label === 'Identity check' && item.status !== 'done' && (
                  <button onClick={() => (window.location.href = '/kyc')} className="text-xs text-blue-700 font-semibold">
                    Open
                  </button>
                )}
                {item.label === 'SACCO membership' && kyc === 'approved' && item.status !== 'done' && (
                  <button onClick={() => (window.location.href = '/sacco')} className="text-xs text-blue-700 font-semibold">
                    Open
                  </button>
                )}
              </div>
            ))}
          </div>
        </div>

        <div className="bg-white/10 backdrop-blur-md rounded-2xl border border-white/20 p-6">
          <h2 className="text-sm font-semibold text-blue-100 mb-4">SACCO SECTION</h2>
          <div className="bg-white rounded-xl p-4 flex flex-col md:flex-row md:items-center md:justify-between gap-3">
            <div>
              <p className="text-sm font-semibold text-slate-900">Grow your SACCO with referrals</p>
              <p className="text-xs text-slate-500">You can register a new member directly or share your referral code/link.</p>
            </div>
            <div className="flex gap-2">
              <button
                onClick={() => (window.location.href = `/register?ref=${encodeURIComponent(refCode)}`)}
                className="inline-flex items-center justify-center gap-2 px-4 py-2 rounded-xl bg-blue-700 text-white text-sm font-semibold hover:bg-blue-800"
              >
                Register new member
                <ArrowRight className="h-4 w-4" />
              </button>
              <button
                onClick={shareReferral}
                className="inline-flex items-center justify-center px-4 py-2 rounded-xl border border-slate-300 text-slate-700 text-sm font-semibold hover:bg-slate-50"
              >
                Share link/code
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Home;
