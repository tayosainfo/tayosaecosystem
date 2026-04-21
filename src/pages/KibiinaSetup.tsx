import React, { useState } from 'react';
import { ArrowRight } from 'lucide-react';
import { platformApi } from '../lib/platformApi';

const KibiinaSetup: React.FC = () => {
  const [action, setAction] = useState('join');
  const [inviteCode, setInviteCode] = useState('');
  const [groupName, setGroupName] = useState('');
  const [contributionAmount, setContributionAmount] = useState('');
  const [cycleFrequency, setCycleFrequency] = useState('weekly');
  const [maxGroupSize, setMaxGroupSize] = useState('15');
  const [payoutOrderPreference, setPayoutOrderPreference] = useState('agreed_order');
  const [notificationPreference, setNotificationPreference] = useState('sms');
  const [languagePreference, setLanguagePreference] = useState('english');
  const [error, setError] = useState('');

  const submit = async () => {
    const token = sessionStorage.getItem('auth_token');
    if (!token) return;
    try {
      await platformApi.submitKibiina(token, {
        action,
        inviteCode: inviteCode || undefined,
        groupName: groupName || undefined,
        contributionAmount: Number(contributionAmount || 0),
        cycleFrequency: cycleFrequency || undefined,
        maxGroupSize: Number(maxGroupSize || 0),
        payoutOrderPreference: payoutOrderPreference || undefined,
        notificationPreference: notificationPreference || undefined,
        languagePreference: languagePreference || undefined,
      });
      window.location.href = '/home';
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to submit Kibiina setup');
    }
  };

  const inputClass = 'w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95';
  const selectClass = 'w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95';
  const isCreate = action === 'create';

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-indigo-900 py-10 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto bg-white/10 backdrop-blur-md rounded-2xl border border-white/20 p-8">
        <div className="text-white mb-6">
          <p className="text-sm text-blue-100">Step 4 of 4 (Optional)</p>
          <h1 className="text-3xl font-bold">Kibiina group setup</h1>
          <p className="text-blue-100">Choose whether to create a new group or join an existing one.</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-3 mb-5">
          <button
            type="button"
            onClick={() => setAction('join')}
            className={`h-12 rounded-xl font-semibold ${!isCreate ? 'bg-white text-blue-900' : 'bg-white/10 text-white border border-white/20'}`}
          >
            Join existing group
          </button>
          <button
            type="button"
            onClick={() => setAction('create')}
            className={`h-12 rounded-xl font-semibold ${isCreate ? 'bg-white text-blue-900' : 'bg-white/10 text-white border border-white/20'}`}
          >
            Create new group
          </button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {!isCreate && (
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-white mb-2">Invite code</label>
              <input className={inputClass} placeholder="Enter group invite code" value={inviteCode} onChange={(e) => setInviteCode(e.target.value)} />
            </div>
          )}

          {isCreate && (
            <>
              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-white mb-2">Group name</label>
                <input className={inputClass} placeholder="Example: Nakawa A Group" value={groupName} onChange={(e) => setGroupName(e.target.value)} />
              </div>
              <div>
                <label className="block text-sm font-medium text-white mb-2">Contribution amount (UGX)</label>
                <input className={inputClass} placeholder="Example: 20000" value={contributionAmount} onChange={(e) => setContributionAmount(e.target.value)} />
              </div>
              <div>
                <label className="block text-sm font-medium text-white mb-2">Cycle frequency</label>
                <select className={selectClass} value={cycleFrequency} onChange={(e) => setCycleFrequency(e.target.value)}>
                  <option value="weekly">Weekly</option>
                  <option value="bi_weekly">Bi-weekly</option>
                  <option value="monthly">Monthly</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-white mb-2">Max group size</label>
                <select className={selectClass} value={maxGroupSize} onChange={(e) => setMaxGroupSize(e.target.value)}>
                  <option value="10">10</option>
                  <option value="15">15</option>
                  <option value="20">20</option>
                  <option value="25">25</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-white mb-2">Payout order</label>
                <select className={selectClass} value={payoutOrderPreference} onChange={(e) => setPayoutOrderPreference(e.target.value)}>
                  <option value="agreed_order">Agreed order</option>
                  <option value="random_draw">Random draw</option>
                  <option value="seniority">Seniority</option>
                </select>
              </div>
            </>
          )}

          <div>
            <label className="block text-sm font-medium text-white mb-2">Reminder preference</label>
            <select className={selectClass} value={notificationPreference} onChange={(e) => setNotificationPreference(e.target.value)}>
              <option value="sms">SMS</option>
              <option value="email">Email</option>
              <option value="push">Push notification</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-white mb-2">Language</label>
            <select className={selectClass} value={languagePreference} onChange={(e) => setLanguagePreference(e.target.value)}>
              <option value="english">English</option>
              <option value="luganda">Luganda</option>
            </select>
          </div>
        </div>
        {error && <p className="text-red-200 mt-4 text-sm">{error}</p>}
        <div className="mt-6">
          <button onClick={submit} className="w-full md:w-auto px-6 bg-white text-blue-900 rounded-xl h-12 font-semibold flex items-center justify-center gap-2 hover:bg-blue-50 transition">
            {isCreate ? 'Create group' : 'Join group'} <ArrowRight className="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>
  );
};

export default KibiinaSetup;

