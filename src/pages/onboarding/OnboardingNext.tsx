import React from 'react';
import { CheckCircle, ArrowRight } from 'lucide-react';

const OnboardingNext: React.FC = () => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-indigo-900 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8 relative overflow-hidden">
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-blue-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-pulse"></div>
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-indigo-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-pulse delay-1000"></div>
      </div>

      <div className="max-w-xl w-full space-y-8 relative z-10">
        <div className="bg-white/10 backdrop-blur-md rounded-2xl p-8 shadow-2xl border border-white/20 text-center">
          <div className="flex justify-center mb-4">
            <div className="relative">
              <div className="absolute inset-0 bg-gradient-to-r from-green-400 to-blue-400 rounded-full blur-lg opacity-75 animate-pulse"></div>
              <div className="relative bg-white rounded-full p-4 shadow-2xl">
                <CheckCircle className="h-12 w-12 text-green-600" />
              </div>
            </div>
          </div>
          <h2 className="text-3xl font-bold text-white mb-2">Location saved</h2>
          <p className="text-blue-100">
            Next onboarding steps (KYC, membership, and Kibiina preferences) are ready to be added.
          </p>

          <div className="mt-6 flex flex-col gap-3">
            <button
              type="button"
              onClick={() => {
                sessionStorage.removeItem('onboarding_phase');
                window.location.href = '/onboarding';
              }}
              className="w-full py-3 px-4 rounded-xl text-blue-900 bg-gradient-to-r from-white to-blue-50 font-medium hover:from-blue-50 hover:to-white transition-all duration-300"
            >
              Back to location
            </button>
            <button
              type="button"
              onClick={() => (window.location.href = '/home')}
              className="group w-full flex justify-center items-center py-3 px-4 rounded-xl text-white bg-white/10 border border-white/20 font-medium hover:bg-white/15 transition-all duration-300"
            >
              Continue
              <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default OnboardingNext;

