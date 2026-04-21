import React, { useEffect, useMemo, useRef, useState } from 'react';
import { MapPin, ArrowRight, CheckCircle, AlertCircle, ChevronDown } from 'lucide-react';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

const getErrorMessage = (err: unknown, fallback: string): string => {
  if (err instanceof Error && err.message) {
    return err.message;
  }
  return fallback;
};

const isSessionError = (message: string): boolean => {
  const msg = message.toLowerCase();
  return msg.includes('invalid or expired session') || msg.includes('missing bearer token');
};

const clearSessionAndReturnToLogin = () => {
  sessionStorage.removeItem('auth_token');
  sessionStorage.removeItem('auth_user');
  window.location.href = '/';
};

async function authedGet<T>(path: string): Promise<T> {
  const token = sessionStorage.getItem('auth_token');
  const res = await fetch(`${API_BASE_URL}${path}`, {
    headers: token ? { Authorization: `Bearer ${token}` } : undefined,
  });
  const payload = (await res.json()) as any;
  if (!res.ok) {
    throw new Error(String(payload?.error || 'Request failed'));
  }
  return payload as T;
}

async function authedPost<T>(path: string, body: unknown): Promise<T> {
  const token = sessionStorage.getItem('auth_token');
  const res = await fetch(`${API_BASE_URL}${path}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: JSON.stringify(body),
  });
  const payload = (await res.json()) as any;
  if (!res.ok) {
    throw new Error(String(payload?.error || 'Request failed'));
  }
  return payload as T;
}

type GeoSelectOption = { id?: string; name: string };

function normalizeGeoItems(input: unknown): GeoSelectOption[] {
  const values: unknown[] = Array.isArray(input)
    ? input
    : typeof input === 'object' && input !== null && Array.isArray((input as any).items)
      ? (input as any).items
      : [];

  return values
    .map((v) => {
      if (typeof v === 'string') {
        return { name: v };
      }
      if (typeof v === 'object' && v !== null) {
        const anyV = v as any;
        const name = String(anyV.name ?? anyV.value ?? anyV.label ?? '').trim();
        const id = anyV.id != null ? String(anyV.id) : undefined;
        if (name) {
          return { name, id };
        }
      }
      return null;
    })
    .filter((x): x is GeoSelectOption => Boolean(x && x.name));
}

const GeoSelect: React.FC<{
  label: string;
  value: string;
  placeholder: string;
  disabled?: boolean;
  options: GeoSelectOption[];
  onChange: (next: string) => void;
}> = ({ label, value, placeholder, disabled, options, onChange }) => {
  const [open, setOpen] = useState(false);
  const wrapRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    const onDoc = (e: MouseEvent) => {
      if (!wrapRef.current) {
        return;
      }
      if (!wrapRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener('mousedown', onDoc);
    return () => document.removeEventListener('mousedown', onDoc);
  }, []);

  useEffect(() => {
    // Close the menu if it becomes disabled
    if (disabled) {
      setOpen(false);
    }
  }, [disabled]);

  const selectedLabel = value || '';

  return (
    <div ref={wrapRef} className="relative">
      <label className="block text-sm font-medium text-white mb-2">{label}</label>
      <button
        type="button"
        disabled={disabled}
        onClick={() => setOpen((v) => !v)}
        className="w-full flex items-center justify-between px-4 py-3 rounded-xl border border-white/30 bg-white/95 text-slate-900 focus:outline-none focus:ring-2 focus:ring-blue-400 disabled:opacity-60"
      >
        <span className={selectedLabel ? 'text-slate-900' : 'text-slate-500'}>{selectedLabel || placeholder}</span>
        <ChevronDown className="h-4 w-4 text-slate-700" />
      </button>

      {open && !disabled && (
        <div className="absolute z-50 mt-2 w-full rounded-xl border border-slate-200 bg-white shadow-xl overflow-hidden">
          <div className="max-h-72 overflow-auto">
            {options.length === 0 ? (
              <div className="px-4 py-3 text-sm text-slate-600">No options</div>
            ) : (
              options.map((opt, idx) => {
                const key = `${opt.id || opt.name}:${idx}`;
                const active = opt.name === value;
                return (
                  <button
                    key={key}
                    type="button"
                    onClick={() => {
                      onChange(opt.name);
                      setOpen(false);
                    }}
                    className={[
                      'w-full text-left px-4 py-3 text-sm',
                      active ? 'bg-blue-50 text-blue-900' : 'text-slate-800 hover:bg-slate-50',
                    ].join(' ')}
                  >
                    {opt.name}
                  </button>
                );
              })
            )}
          </div>
        </div>
      )}
    </div>
  );
};

const Onboarding: React.FC = () => {
  const [districts, setDistricts] = useState<GeoSelectOption[]>([]);
  const [counties, setCounties] = useState<GeoSelectOption[]>([]);
  const [subCounties, setSubCounties] = useState<GeoSelectOption[]>([]);
  const [parishes, setParishes] = useState<GeoSelectOption[]>([]);
  const [villages, setVillages] = useState<GeoSelectOption[]>([]);

  const [district, setDistrict] = useState('');
  const [county, setCounty] = useState('');
  const [subCounty, setSubCounty] = useState('');
  const [parish, setParish] = useState('');
  const [village, setVillage] = useState('');

  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const userId = useMemo(() => {
    try {
      const raw = sessionStorage.getItem('auth_user');
      if (!raw) {
        return '';
      }
      const u = JSON.parse(raw) as { id?: string };
      return u.id || '';
    } catch {
      return '';
    }
  }, []);

  useEffect(() => {
    let mounted = true;
    authedGet<any>('/api/v1/geo?level=district')
      .then((data) => {
        if (mounted) {
          setDistricts(normalizeGeoItems(data));
        }
      })
      .catch((e) => {
        const message = getErrorMessage(e, 'Failed to load districts');
        if (isSessionError(message)) {
          setError('Your session expired. Please sign in again.');
          clearSessionAndReturnToLogin();
          return;
        }
        setError(message);
      });
    return () => {
      mounted = false;
    };
  }, []);

  const loadChildren = async (level: string, parentName: string) => {
    const q = new URLSearchParams({ level, parent: parentName }).toString();
    return await authedGet<any>(`/api/v1/geo?${q}`);
  };

  useEffect(() => {
    if (!district) {
      setCounties([]);
      setCounty('');
      return;
    }
    setError('');
    loadChildren('county', district)
      .then((data) => {
        setCounties(normalizeGeoItems(data));
      })
      .catch((e) => setError(getErrorMessage(e, 'Failed to load counties')));
    setCounty('');
    setSubCounty('');
    setParish('');
    setVillage('');
    setSubCounties([]);
    setParishes([]);
    setVillages([]);
  }, [district]);

  useEffect(() => {
    if (!county) {
      setSubCounties([]);
      setSubCounty('');
      return;
    }
    setError('');
    loadChildren('subcounty', county)
      .then((data) => {
        setSubCounties(normalizeGeoItems(data));
      })
      .catch((e) => setError(getErrorMessage(e, 'Failed to load sub-counties')));
    setSubCounty('');
    setParish('');
    setVillage('');
    setParishes([]);
    setVillages([]);
  }, [county]);

  useEffect(() => {
    if (!subCounty) {
      setParishes([]);
      setParish('');
      return;
    }
    setError('');
    loadChildren('parish', subCounty)
      .then((data) => {
        setParishes(normalizeGeoItems(data));
      })
      .catch((e) => setError(getErrorMessage(e, 'Failed to load parishes')));
    setParish('');
    setVillage('');
    setVillages([]);
  }, [subCounty]);

  useEffect(() => {
    if (!parish) {
      setVillages([]);
      setVillage('');
      return;
    }
    setError('');
    loadChildren('village', parish)
      .then((data) => {
        setVillages(normalizeGeoItems(data));
      })
      .catch((e) => setError(getErrorMessage(e, 'Failed to load villages')));
    setVillage('');
  }, [parish]);

  const submit = async () => {
    setError('');
    setSuccess('');
    const token = sessionStorage.getItem('auth_token');
    if (!token) {
      setError('Your session expired. Please sign in again.');
      window.location.href = '/';
      return;
    }
    if (!userId) {
      setError('Missing user id. Please sign in again.');
      return;
    }
    if (!district || !county || !subCounty || !parish || !village) {
      setError('Please select your location (district → county → sub-county → parish → village).');
      return;
    }
    setIsLoading(true);
    try {
      await authedPost('/api/v1/onboarding/phase', {
        userId,
        // Backend expects phase 2-4; phase 2 is the location step.
        phase: 2,
        geo: { district, county, sub_county: subCounty, parish, village },
      });
      sessionStorage.setItem('onboarding_phase', '2');
      setSuccess('Onboarding started. Your location was saved.');
      window.location.href = '/onboarding/next';
    } catch (e) {
      const message = getErrorMessage(e, 'Failed to save onboarding details');
      if (isSessionError(message)) {
        setError('Your session expired. Please sign in again.');
        clearSessionAndReturnToLogin();
        return;
      }
      setError(message);
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

      <div className="max-w-2xl w-full space-y-8 relative z-10">
        <div className="text-center">
          <div className="flex justify-center mb-6">
            <div className="relative">
              <div className="absolute inset-0 bg-gradient-to-r from-blue-400 to-purple-500 rounded-full blur-lg opacity-75 animate-pulse"></div>
              <div className="relative bg-white rounded-full p-4 shadow-2xl">
                <MapPin className="h-12 w-12 text-blue-600" />
              </div>
            </div>
          </div>
          <h2 className="text-4xl font-bold text-white mb-2">Onboarding</h2>
          <p className="text-blue-100 text-lg">Start by confirming your location.</p>
        </div>

        <div className="bg-white/10 backdrop-blur-md rounded-2xl p-8 shadow-2xl border border-white/20">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
            <GeoSelect
              label="District"
              value={district}
              placeholder="Select district"
              options={districts}
              onChange={setDistrict}
            />

            <GeoSelect
              label="County"
              value={county}
              placeholder="Select county"
              options={counties}
              disabled={!district}
              onChange={setCounty}
            />

            <GeoSelect
              label="Sub-county"
              value={subCounty}
              placeholder="Select sub-county"
              options={subCounties}
              disabled={!county}
              onChange={setSubCounty}
            />

            <GeoSelect
              label="Parish"
              value={parish}
              placeholder="Select parish"
              options={parishes}
              disabled={!subCounty}
              onChange={setParish}
            />

            <div className="md:col-span-2">
              <GeoSelect
                label="Village"
                value={village}
                placeholder="Select village"
                options={villages}
                disabled={!parish}
                onChange={setVillage}
              />
            </div>
          </div>

          {error && (
            <div className="mt-5 flex items-center bg-red-500/20 text-red-100 px-4 py-3 rounded-xl backdrop-blur-sm border border-red-400/50">
              <AlertCircle className="h-5 w-5 mr-2 shrink-0" />
              <span className="text-sm">{error}</span>
            </div>
          )}
          {success && (
            <div className="mt-5 flex items-center bg-green-500/20 text-green-100 px-4 py-3 rounded-xl backdrop-blur-sm border border-green-400/50">
              <CheckCircle className="h-5 w-5 mr-2 shrink-0" />
              <span className="text-sm">{success}</span>
            </div>
          )}

          <div className="mt-6">
            <button
              type="button"
              disabled={isLoading}
              onClick={submit}
              className="group w-full flex justify-center items-center py-3 px-4 rounded-xl text-blue-900 bg-gradient-to-r from-white to-blue-50 font-medium hover:from-blue-50 hover:to-white transition-all duration-300 transform hover:scale-105 disabled:opacity-60 disabled:cursor-not-allowed"
            >
              {isLoading ? (
                <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-blue-900"></div>
              ) : (
                <>
                  Save & continue
                  <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
                </>
              )}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Onboarding;

