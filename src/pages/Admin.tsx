import React, { useEffect, useMemo, useState } from 'react';
import { platformApi } from '../lib/platformApi';

type Tab = 'kyc' | 'fees';

const Admin: React.FC = () => {
  const [tab, setTab] = useState<Tab>('kyc');
  const [status, setStatus] = useState('pending');
  const [items, setItems] = useState<any[]>([]);
  const [error, setError] = useState('');
  const [fees, setFees] = useState({ registrationFeeUGX: 0, saccoEntranceFeeUGX: 0, transactionFeePct: 0 });
  const token = useMemo(() => sessionStorage.getItem('auth_token') || '', []);
  const adminSecret = String(import.meta.env.VITE_ADMIN_API_KEY || '').trim();

  const loadKyc = async () => {
    if (!token || !adminSecret) return;
    setError('');
    try {
      const res = await platformApi.adminListKyc(token, adminSecret, status);
      setItems(res.items || []);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load KYC queue');
    }
  };

  const loadFees = async () => {
    if (!token || !adminSecret) return;
    setError('');
    try {
      const res = await platformApi.adminGetFees(token, adminSecret);
      const v = res.value || {};
      setFees({
        registrationFeeUGX: Number(v.registrationFeeUGX || 0),
        saccoEntranceFeeUGX: Number(v.saccoEntranceFeeUGX || 0),
        transactionFeePct: Number(v.transactionFeePct || 0),
      });
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load fees');
    }
  };

  useEffect(() => {
    if (tab === 'kyc') loadKyc();
    if (tab === 'fees') loadFees();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [tab, status]);

  const decide = async (userId: string, decision: 'approved' | 'rejected') => {
    if (!token || !adminSecret) return;
    setError('');
    try {
      await platformApi.adminDecideKyc(token, adminSecret, userId, decision);
      await loadKyc();
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to submit decision');
    }
  };

  const saveFees = async () => {
    if (!token || !adminSecret) return;
    setError('');
    try {
      await platformApi.adminSetFees(token, adminSecret, fees as any);
      await loadFees();
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to save fees');
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-indigo-900 py-10 px-4 sm:px-6 lg:px-8">
      <div className="max-w-6xl mx-auto space-y-6">
        <div className="bg-white/10 backdrop-blur-md rounded-2xl border border-white/20 p-6 text-white flex flex-col md:flex-row md:items-center md:justify-between gap-3">
          <div>
            <h1 className="text-3xl font-bold">Admin Panel</h1>
            <p className="text-blue-100 text-sm">KYC review and platform settings.</p>
          </div>
          <div className="flex gap-2">
            <button onClick={() => setTab('kyc')} className={`px-4 py-2 rounded-xl text-sm font-semibold ${tab === 'kyc' ? 'bg-white text-blue-900' : 'bg-white/10 text-white border border-white/20'}`}>
              KYC Queue
            </button>
            <button onClick={() => setTab('fees')} className={`px-4 py-2 rounded-xl text-sm font-semibold ${tab === 'fees' ? 'bg-white text-blue-900' : 'bg-white/10 text-white border border-white/20'}`}>
              Fees & Charges
            </button>
          </div>
        </div>

        {!adminSecret && (
          <div className="bg-red-500/20 border border-red-400/50 text-red-100 px-4 py-3 rounded-xl">
            Missing `VITE_ADMIN_API_KEY` in `.env`. Set it to the same value as `ADMIN_API_KEY` and restart the frontend.
          </div>
        )}

        {error && (
          <div className="bg-red-500/20 border border-red-400/50 text-red-100 px-4 py-3 rounded-xl">{error}</div>
        )}

        {tab === 'kyc' && (
          <div className="bg-white/10 backdrop-blur-md rounded-2xl border border-white/20 p-6">
            <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-3 mb-4">
              <h2 className="text-white font-semibold">KYC review queue</h2>
              <div className="flex gap-2">
                <select value={status} onChange={(e) => setStatus(e.target.value)} className="h-10 rounded-xl px-3 bg-white/95 text-slate-900">
                  <option value="pending">Pending</option>
                  <option value="approved">Approved</option>
                  <option value="rejected">Rejected</option>
                  <option value="all">All users</option>
                </select>
                <button onClick={loadKyc} className="h-10 px-4 rounded-xl bg-white text-blue-900 font-semibold">Refresh</button>
              </div>
            </div>

            <div className="space-y-3">
              {items.length === 0 ? (
                <div className="text-blue-100 text-sm">No records.</div>
              ) : (
                items.map((k) => (
                  <div key={String(k.userId)} className="bg-white rounded-xl p-4 flex flex-col md:flex-row md:items-center md:justify-between gap-3">
                    <div className="text-slate-900">
                      <div className="font-semibold">{String(k.fullName || 'Unnamed user')} ({String(k.userId)})</div>
                      <div className="text-xs text-slate-500">Phone: {String(k.phoneE164 || 'N/A')}</div>
                      <div className="text-xs text-slate-500">Email: {String(k.contactEmail || 'N/A')}</div>
                      <div className="text-xs text-slate-500">Status: {String(k.status)}</div>
                      <div className="text-xs text-slate-500">ID: {String(k.idType || '')} {String(k.idNumber || '')}</div>
                    </div>
                    <div className="flex gap-2">
                      <button onClick={() => decide(String(k.userId), 'approved')} className="px-3 h-10 rounded-xl bg-green-600 text-white text-sm font-semibold">Approve</button>
                      <button onClick={() => decide(String(k.userId), 'rejected')} className="px-3 h-10 rounded-xl bg-red-600 text-white text-sm font-semibold">Reject</button>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        )}

        {tab === 'fees' && (
          <div className="bg-white/10 backdrop-blur-md rounded-2xl border border-white/20 p-6">
            <h2 className="text-white font-semibold mb-4">Fees & charges</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-white mb-2">Registration fee (UGX)</label>
                <input className="w-full h-11 rounded-xl px-3 bg-white/95" value={fees.registrationFeeUGX} onChange={(e) => setFees((p) => ({ ...p, registrationFeeUGX: Number(e.target.value || 0) }))} />
              </div>
              <div>
                <label className="block text-sm font-medium text-white mb-2">SACCO entrance fee (UGX)</label>
                <input className="w-full h-11 rounded-xl px-3 bg-white/95" value={fees.saccoEntranceFeeUGX} onChange={(e) => setFees((p) => ({ ...p, saccoEntranceFeeUGX: Number(e.target.value || 0) }))} />
              </div>
              <div>
                <label className="block text-sm font-medium text-white mb-2">Transaction fee (%)</label>
                <input className="w-full h-11 rounded-xl px-3 bg-white/95" value={fees.transactionFeePct} onChange={(e) => setFees((p) => ({ ...p, transactionFeePct: Number(e.target.value || 0) }))} />
              </div>
            </div>
            <div className="mt-5">
              <button onClick={saveFees} className="h-11 px-5 rounded-xl bg-white text-blue-900 font-semibold">Save</button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Admin;

