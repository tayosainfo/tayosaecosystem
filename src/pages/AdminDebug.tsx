import { useState, useEffect } from 'react';
import { supabase } from '../lib/supabase';
import { useAuth } from '../hooks/useAuth';

interface DiagnosticResult {
  step: string;
  status: 'success' | 'error' | 'warning' | 'info';
  message: string;
  data?: any;
}

export default function AdminDebug() {
  const [results, setResults] = useState<DiagnosticResult[]>([]);
  const [loading, setLoading] = useState(false);
  const { user: appUser, isAuthenticated } = useAuth();

  const addResult = (result: DiagnosticResult) => {
    setResults(prev => [...prev, result]);
  };

  const runDiagnostics = async () => {
    setResults([]);
    setLoading(true);

    try {
      // Step 1: Check if user is logged in via app auth
      addResult({ step: '1', status: 'info', message: 'Checking application authentication status...' });
      
      if (!isAuthenticated || !appUser) {
        addResult({ 
          step: '1', 
          status: 'error', 
          message: 'Not logged in via application auth',
          data: { isAuthenticated, user: appUser }
        });
        setLoading(false);
        return;
      }

      addResult({ 
        step: '1', 
        status: 'success', 
        message: `Logged in as: ${appUser.email}`,
        data: { userId: appUser.id, email: appUser.email }
      });

      // Step 2: Check database role
      addResult({ step: '2', status: 'info', message: 'Checking database role...' });
      
      // Query users_identity table using app user email
      const { data: userData, error: dbError } = await supabase
        .from('users_identity')
        .select('user_id, full_name, auth_email, role, supabase_user_id, status')
        .eq('auth_email', appUser.email)
        .single();

      if (dbError) {
        addResult({ 
          step: '2', 
          status: 'error', 
          message: 'Failed to query database',
          data: dbError 
        });
      } else if (!userData) {
        addResult({ 
          step: '2', 
          status: 'error', 
          message: 'User not found in users_identity table' 
        });
      } else {
        addResult({ 
          step: '2', 
          status: userData.role === 'admin' ? 'success' : 'warning', 
          message: `Database role: ${userData.role}`,
          data: userData
        });

        // Step 3: Check supabase_user_id linkage
        addResult({ step: '3', status: 'info', message: 'Checking Supabase user ID linkage...' });
        
        if (!userData.supabase_user_id) {
          addResult({ 
            step: '3', 
            status: 'warning', 
            message: '⚠️ supabase_user_id is NULL (may not be needed for custom auth)',
            data: { 
              note: 'Your app uses custom authentication, not Supabase Auth directly'
            }
          });
        } else {
          addResult({ 
            step: '3', 
            status: 'success', 
            message: '✅ supabase_user_id is set' 
          });
        }
      }

      // Step 4: Summary
      addResult({ step: '4', status: 'info', message: 'Diagnosis complete' });

    } catch (error) {
      addResult({ 
        step: 'ERROR', 
        status: 'error', 
        message: 'Unexpected error during diagnostics',
        data: error 
      });
    } finally {
      setLoading(false);
    }
  };

  const checkAdminStatus = async () => {
    try {
      addResult({ step: 'CHECK', status: 'info', message: 'Checking admin status...' });
      
      if (!appUser) {
        addResult({ step: 'CHECK', status: 'error', message: 'No user logged in' });
        return;
      }

      const { data: userData, error } = await supabase
        .from('users_identity')
        .select('role')
        .eq('auth_email', appUser.email)
        .single();

      if (error) {
        addResult({ step: 'CHECK', status: 'error', message: 'Failed to check role', data: error });
        return;
      }

      if (userData?.role === 'admin') {
        addResult({ 
          step: 'CHECK', 
          status: 'success', 
          message: '✅ User IS an admin! Redirecting to admin dashboard...' 
        });
        setTimeout(() => window.location.href = '/admin', 2000);
      } else {
        addResult({ 
          step: 'CHECK', 
          status: 'error', 
          message: `❌ User role is '${userData?.role}', not 'admin'` 
        });
      }
    } catch (error) {
      addResult({ step: 'CHECK', status: 'error', message: 'Error checking admin status', data: error });
    }
  };

  useEffect(() => {
    runDiagnostics();
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-purple-900 to-pink-900 p-8">
      <div className="max-w-4xl mx-auto">
        <div className="bg-white/10 backdrop-blur-md rounded-xl border border-white/20 p-6 mb-6">
          <h1 className="text-3xl font-bold text-white mb-2">Admin Access Diagnostics</h1>
          <p className="text-blue-200">Troubleshooting admin role and JWT token issues</p>
        </div>

        <div className="bg-white/10 backdrop-blur-md rounded-xl border border-white/20 p-6 mb-6">
          <div className="flex gap-4 mb-6">
            <button
              onClick={runDiagnostics}
              disabled={loading}
              className="px-6 py-3 rounded-lg bg-blue-600 text-white font-semibold hover:bg-blue-700 disabled:opacity-50"
            >
              {loading ? 'Running...' : 'Run Diagnostics'}
            </button>
            <button
              onClick={checkAdminStatus}
              disabled={loading}
              className="px-6 py-3 rounded-lg bg-green-600 text-white font-semibold hover:bg-green-700 disabled:opacity-50"
            >
              Check Admin Status
            </button>
          </div>

          <div className="space-y-3">
            {results.map((result, index) => (
              <div
                key={index}
                className={`p-4 rounded-lg border ${
                  result.status === 'success' ? 'bg-green-500/20 border-green-400/50' :
                  result.status === 'error' ? 'bg-red-500/20 border-red-400/50' :
                  result.status === 'warning' ? 'bg-yellow-500/20 border-yellow-400/50' :
                  'bg-blue-500/20 border-blue-400/50'
                }`}
              >
                <div className="flex items-start gap-3">
                  <span className="text-white font-bold">Step {result.step}:</span>
                  <div className="flex-1">
                    <p className="text-white font-semibold">{result.message}</p>
                    {result.data && (
                      <pre className="mt-2 text-xs text-blue-100 bg-black/20 p-2 rounded overflow-x-auto">
                        {JSON.stringify(result.data, null, 2)}
                      </pre>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-white/10 backdrop-blur-md rounded-xl border border-white/20 p-6">
          <h2 className="text-xl font-bold text-white mb-4">Common Issues & Fixes</h2>
          <div className="space-y-4 text-white">
            <div>
              <h3 className="font-semibold text-yellow-300">Issue: User role is 'user' but should be 'admin'</h3>
              <p className="text-sm text-blue-200 mt-1">
                The user record in the database doesn't have the admin role assigned.
                Check your database to ensure the role field is set to 'admin' for this user.
              </p>
            </div>
            <div>
              <h3 className="font-semibold text-yellow-300">Issue: User not found in users_identity table</h3>
              <p className="text-sm text-blue-200 mt-1">
                The user's email doesn't match between the application auth and the users_identity table.
                Verify the email is exactly the same in both places.
              </p>
            </div>
            <div>
              <h3 className="font-semibold text-yellow-300">How to fix admin access</h3>
              <p className="text-sm text-blue-200 mt-1">
                1. Run diagnostics to see the current role<br/>
                2. If role is 'user', update it to 'admin' in the database<br/>
                3. Click "Check Admin Status" to verify<br/>
                4. Navigate to /admin to access the dashboard
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
