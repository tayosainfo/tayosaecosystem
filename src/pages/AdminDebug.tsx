import { useState, useEffect } from 'react';
import { supabase } from '../lib/supabase';

interface DiagnosticResult {
  step: string;
  status: 'success' | 'error' | 'warning' | 'info';
  message: string;
  data?: any;
}

export default function AdminDebug() {
  const [results, setResults] = useState<DiagnosticResult[]>([]);
  const [loading, setLoading] = useState(false);

  const addResult = (result: DiagnosticResult) => {
    setResults(prev => [...prev, result]);
  };

  const runDiagnostics = async () => {
    setResults([]);
    setLoading(true);

    try {
      // Step 1: Check if user is logged in
      addResult({ step: '1', status: 'info', message: 'Checking authentication status...' });
      const { data: { user }, error: userError } = await supabase.auth.getUser();
      
      if (userError || !user) {
        addResult({ 
          step: '1', 
          status: 'error', 
          message: 'Not logged in or session expired',
          data: userError 
        });
        setLoading(false);
        return;
      }

      addResult({ 
        step: '1', 
        status: 'success', 
        message: `Logged in as: ${user.email}`,
        data: { userId: user.id, email: user.email }
      });

      // Step 2: Check JWT token claims
      addResult({ step: '2', status: 'info', message: 'Checking JWT token claims...' });
      const { data: { session } } = await supabase.auth.getSession();
      
      if (!session) {
        addResult({ step: '2', status: 'error', message: 'No active session' });
        setLoading(false);
        return;
      }

      // Decode JWT
      const token = session.access_token;
      const payload = JSON.parse(atob(token.split('.')[1]));
      const userRole = payload.app_metadata?.user_role || payload.user_role || 'NOT SET';

      addResult({ 
        step: '2', 
        status: userRole === 'admin' ? 'success' : 'warning', 
        message: `JWT token role: ${userRole}`,
        data: { 
          app_metadata: payload.app_metadata,
          user_role: userRole,
          expires: new Date(payload.exp * 1000).toLocaleString()
        }
      });

      // Step 3: Check database role
      addResult({ step: '3', status: 'info', message: 'Checking database role...' });
      
      // Query users_identity table
      const { data: userData, error: dbError } = await supabase
        .from('users_identity')
        .select('user_id, full_name, auth_email, role, supabase_user_id, status')
        .eq('auth_email', user.email)
        .single();

      if (dbError) {
        addResult({ 
          step: '3', 
          status: 'error', 
          message: 'Failed to query database',
          data: dbError 
        });
      } else if (!userData) {
        addResult({ 
          step: '3', 
          status: 'error', 
          message: 'User not found in users_identity table' 
        });
      } else {
        addResult({ 
          step: '3', 
          status: userData.role === 'admin' ? 'success' : 'warning', 
          message: `Database role: ${userData.role}`,
          data: userData
        });

        // Step 4: Check supabase_user_id linkage
        addResult({ step: '4', status: 'info', message: 'Checking Supabase user ID linkage...' });
        
        if (!userData.supabase_user_id) {
          addResult({ 
            step: '4', 
            status: 'error', 
            message: '❌ supabase_user_id is NULL - This is the problem!',
            data: { 
              fix: 'Run: UPDATE users_identity SET supabase_user_id = \'' + user.id + '\' WHERE auth_email = \'' + user.email + '\';'
            }
          });
        } else if (userData.supabase_user_id !== user.id) {
          addResult({ 
            step: '4', 
            status: 'error', 
            message: '❌ supabase_user_id mismatch',
            data: { 
              expected: user.id,
              actual: userData.supabase_user_id,
              fix: 'Run: UPDATE users_identity SET supabase_user_id = \'' + user.id + '\' WHERE auth_email = \'' + user.email + '\';'
            }
          });
        } else {
          addResult({ 
            step: '4', 
            status: 'success', 
            message: '✅ supabase_user_id correctly linked' 
          });
        }
      }

      // Step 5: Summary
      addResult({ step: '5', status: 'info', message: 'Diagnosis complete' });

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

  const forceRefresh = async () => {
    try {
      addResult({ step: 'REFRESH', status: 'info', message: 'Forcing token refresh...' });
      const { data, error } = await supabase.auth.refreshSession();
      
      if (error) {
        addResult({ step: 'REFRESH', status: 'error', message: 'Token refresh failed', data: error });
        return;
      }

      if (data.session) {
        const token = data.session.access_token;
        const payload = JSON.parse(atob(token.split('.')[1]));
        const userRole = payload.app_metadata?.user_role || 'NOT SET';
        
        addResult({ 
          step: 'REFRESH', 
          status: 'success', 
          message: `Token refreshed. New role: ${userRole}`,
          data: { user_role: userRole }
        });

        if (userRole === 'admin') {
          addResult({ 
            step: 'REFRESH', 
            status: 'success', 
            message: '✅ Admin access should now work! Reloading page...' 
          });
          setTimeout(() => window.location.reload(), 2000);
        } else {
          addResult({ 
            step: 'REFRESH', 
            status: 'warning', 
            message: '⚠️ Token refreshed but role is still not admin. Check database linkage.' 
          });
        }
      }
    } catch (error) {
      addResult({ step: 'REFRESH', status: 'error', message: 'Refresh error', data: error });
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
              onClick={forceRefresh}
              disabled={loading}
              className="px-6 py-3 rounded-lg bg-green-600 text-white font-semibold hover:bg-green-700 disabled:opacity-50"
            >
              Force Token Refresh
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
              <h3 className="font-semibold text-yellow-300">Issue: supabase_user_id is NULL</h3>
              <p className="text-sm text-blue-200 mt-1">
                The users_identity record is not linked to the Supabase Auth user.
                Run the SQL fix shown in Step 4 above.
              </p>
            </div>
            <div>
              <h3 className="font-semibold text-yellow-300">Issue: JWT role is 'user' but database role is 'admin'</h3>
              <p className="text-sm text-blue-200 mt-1">
                Token was issued before role was assigned. Click "Force Token Refresh" button above.
              </p>
            </div>
            <div>
              <h3 className="font-semibold text-yellow-300">Issue: Custom claims hook not configured</h3>
              <p className="text-sm text-blue-200 mt-1">
                Go to Supabase Dashboard → Authentication → Hooks → Enable "Custom Access Token" hook
                → Set function to: public.custom_access_token_hook
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
