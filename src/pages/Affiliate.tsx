import React from 'react';

const Affiliate: React.FC = () => {
  const code = 'AMARA-7742';
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-indigo-900 py-10 px-4 sm:px-6 lg:px-8">
      <div className="max-w-5xl mx-auto bg-white/10 backdrop-blur-md rounded-2xl border border-white/20 p-8">
        <div className="text-white">
          <p className="text-sm text-blue-100">Add a member - your rewards</p>
          <h1 className="text-4xl font-bold">UGX 45,000</h1>
          <p className="text-blue-100">Total rewards earned</p>
        </div>
        <div className="mt-6 grid grid-cols-1 lg:grid-cols-2 gap-4">
          <div className="bg-white rounded-xl p-4">
            <p className="text-xs text-slate-500">Your referral code</p>
            <p className="text-2xl font-bold tracking-wide text-slate-900">{code}</p>
            <div className="grid grid-cols-3 gap-2 mt-4">
              <button className="border rounded-xl py-3 text-sm">Share link</button>
              <button className="border rounded-xl py-3 text-sm">Send SMS</button>
              <button className="border rounded-xl py-3 text-sm">WhatsApp</button>
            </div>
          </div>

          <div className="bg-white rounded-xl p-4">
            <p className="font-semibold text-sm mb-2">Members you added</p>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between"><span>Grace Nakato</span><span className="text-emerald-700">+15,000</span></div>
              <div className="flex justify-between"><span>Joseph Oluwa</span><span className="text-emerald-700">+15,000</span></div>
              <div className="flex justify-between"><span>Aisha Mutesi</span><span className="text-amber-600">Pending</span></div>
            </div>
          </div>
        </div>
        <div className="mt-5">
          <button onClick={() => (window.location.href = '/home')} className="bg-white text-blue-900 rounded-xl px-6 h-11 text-sm font-semibold">
            Back to home
          </button>
        </div>
      </div>
    </div>
  );
};

export default Affiliate;
