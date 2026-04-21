import React, { useEffect, useState } from 'react';
import { ArrowRight } from 'lucide-react';
import { platformApi } from '../lib/platformApi';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

type GeoOption = { name: string };

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

const isSessionError = (message: string): boolean => {
  const msg = message.toLowerCase();
  return msg.includes('invalid or expired session') || msg.includes('missing bearer token');
};

const clearSessionAndReturnToLogin = () => {
  sessionStorage.removeItem('auth_token');
  sessionStorage.removeItem('auth_user');
  window.location.href = '/';
};

const normalizeGeo = (input: unknown): GeoOption[] => {
  const arr: unknown[] = Array.isArray(input)
    ? input
    : typeof input === 'object' && input !== null && Array.isArray((input as any).items)
      ? (input as any).items
      : [];
  return arr
    .map((v) => (typeof v === 'string' ? { name: v } : typeof v === 'object' && v !== null ? { name: String((v as any).name || '').trim() } : null))
    .filter((v): v is GeoOption => Boolean(v && v.name));
};

const SaccoSetup: React.FC = () => {
  const [districts, setDistricts] = useState<GeoOption[]>([]);
  const [counties, setCounties] = useState<GeoOption[]>([]);
  const [subCounties, setSubCounties] = useState<GeoOption[]>([]);
  const [parishes, setParishes] = useState<GeoOption[]>([]);
  const [villages, setVillages] = useState<GeoOption[]>([]);
  const [district, setDistrict] = useState('');
  const [county, setCounty] = useState('');
  const [subCounty, setSubCounty] = useState('');
  const [parish, setParish] = useState('');
  const [village, setVillage] = useState('');
  const [streetPlot, setStreetPlot] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [provider, setProvider] = useState('');
  const [momoNumber, setMomoNumber] = useState('');
  const [secondaryMomoNumber, setSecondaryMomoNumber] = useState('');
  const [frequency, setFrequency] = useState('weekly');
  const [goalAmount, setGoalAmount] = useState('');
  const [goalPurpose, setGoalPurpose] = useState('school_fees');
  const [sharesToPurchase, setSharesToPurchase] = useState('1');
  const [entranceMethod, setEntranceMethod] = useState('mobile_money');
  const [error, setError] = useState('');

  useEffect(() => {
    authedGet<any>('/api/v1/geo?level=district')
      .then((d) => setDistricts(normalizeGeo(d)))
      .catch((e) => {
        const message = e instanceof Error ? e.message : 'Failed to load districts';
        if (isSessionError(message)) {
          clearSessionAndReturnToLogin();
          return;
        }
        setError(message);
      });
  }, []);

  useEffect(() => {
    if (!district) {
      setCounties([]);
      setCounty('');
      return;
    }
    const q = new URLSearchParams({ level: 'county', parent: district }).toString();
    authedGet<any>(`/api/v1/geo?${q}`)
      .then((d) => setCounties(normalizeGeo(d)))
      .catch((e) => setError(e instanceof Error ? e.message : 'Failed to load counties'));
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
    const q = new URLSearchParams({ level: 'subcounty', parent: county }).toString();
    authedGet<any>(`/api/v1/geo?${q}`)
      .then((d) => setSubCounties(normalizeGeo(d)))
      .catch((e) => setError(e instanceof Error ? e.message : 'Failed to load sub-counties'));
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
    const q = new URLSearchParams({ level: 'parish', parent: subCounty }).toString();
    authedGet<any>(`/api/v1/geo?${q}`)
      .then((d) => setParishes(normalizeGeo(d)))
      .catch((e) => setError(e instanceof Error ? e.message : 'Failed to load parishes'));
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
    const q = new URLSearchParams({ level: 'village', parent: parish }).toString();
    authedGet<any>(`/api/v1/geo?${q}`)
      .then((d) => setVillages(normalizeGeo(d)))
      .catch((e) => setError(e instanceof Error ? e.message : 'Failed to load villages'));
    setVillage('');
  }, [parish]);

  const submit = async () => {
    if (!district || !county || !subCounty || !parish || !village) {
      setError('Please confirm location from district down to village.');
      return;
    }
    const token = sessionStorage.getItem('auth_token');
    if (!token) {
      setError('Session expired. Please sign in again.');
      return;
    }
    setIsSubmitting(true);
    try {
      await platformApi.submitSacco(token, {
        district,
        county,
        subCounty,
        parish,
        village,
        streetPlot: streetPlot || undefined,
        mobileMoneyProvider: provider,
        mobileMoneyNumber: momoNumber,
        secondaryMoMoNumber: secondaryMomoNumber || undefined,
        contributionFrequency: frequency || undefined,
        savingsGoalAmount: Number(goalAmount || 0),
        savingsGoalPurpose: goalPurpose || undefined,
        sharesToPurchase: Number(sharesToPurchase || 1),
        entranceFeePaymentMethod: entranceMethod || undefined,
      });
      window.location.href = '/home';
    } catch (e) {
      const message = e instanceof Error ? e.message : 'Failed to save SACCO setup';
      if (isSessionError(message)) {
        clearSessionAndReturnToLogin();
        return;
      }
      setError(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  const inputClass = 'w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95';
  const selectClass = 'w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95';

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-indigo-900 py-10 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto bg-white/10 backdrop-blur-md rounded-2xl border border-white/20 p-8">
        <div className="text-white mb-6">
          <p className="text-sm text-blue-100">Step 3 of 4 - SACCO setup</p>
          <h1 className="text-3xl font-bold">SACCO membership setup</h1>
          <p className="text-blue-100">Choose your parish and complete membership preferences.</p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <select className={selectClass} value={district} onChange={(e) => setDistrict(e.target.value)}>
            <option value="">Select district</option>
            {districts.map((d) => <option key={d.name} value={d.name}>{d.name}</option>)}
          </select>
          <select className={selectClass} value={county} onChange={(e) => setCounty(e.target.value)} disabled={!district}>
            <option value="">Select county</option>
            {counties.map((d) => <option key={d.name} value={d.name}>{d.name}</option>)}
          </select>
          <select className={selectClass} value={subCounty} onChange={(e) => setSubCounty(e.target.value)} disabled={!county}>
            <option value="">Select sub-county</option>
            {subCounties.map((d) => <option key={d.name} value={d.name}>{d.name}</option>)}
          </select>
          <select className={selectClass} value={parish} onChange={(e) => setParish(e.target.value)} disabled={!subCounty}>
            <option value="">Select parish</option>
            {parishes.map((d) => <option key={d.name} value={d.name}>{d.name}</option>)}
          </select>
          <select className={`${selectClass} md:col-span-2`} value={village} onChange={(e) => setVillage(e.target.value)} disabled={!parish}>
            <option value="">Select village</option>
            {villages.map((d) => <option key={d.name} value={d.name}>{d.name}</option>)}
          </select>
          <input className={`${inputClass} md:col-span-2`} placeholder="Street / plot number (optional)" value={streetPlot} onChange={(e) => setStreetPlot(e.target.value)} />
          <select className={selectClass} value={provider} onChange={(e) => setProvider(e.target.value)}>
            <option value="">Select mobile money provider</option>
            <option value="MTN_MOMO">MTN MoMo</option>
            <option value="AIRTEL_MONEY">Airtel Money</option>
          </select>
          <input className={inputClass} placeholder="Mobile money number" value={momoNumber} onChange={(e) => setMomoNumber(e.target.value)} />
          <input className={inputClass} placeholder="Secondary MoMo number (optional)" value={secondaryMomoNumber} onChange={(e) => setSecondaryMomoNumber(e.target.value)} />
          <select className={selectClass} value={frequency} onChange={(e) => setFrequency(e.target.value)}>
            <option value="weekly">Weekly</option>
            <option value="monthly">Monthly</option>
            <option value="per_cycle">Per-cycle</option>
          </select>
          <input className={inputClass} placeholder="Savings goal amount" value={goalAmount} onChange={(e) => setGoalAmount(e.target.value)} />
          <select className={selectClass} value={goalPurpose} onChange={(e) => setGoalPurpose(e.target.value)}>
            <option value="school_fees">School fees</option>
            <option value="business">Business</option>
            <option value="land">Land</option>
            <option value="health">Health</option>
            <option value="other">Other</option>
          </select>
          <input className={inputClass} placeholder="Shares to purchase" value={sharesToPurchase} onChange={(e) => setSharesToPurchase(e.target.value)} />
          <select className={`${selectClass} md:col-span-2`} value={entranceMethod} onChange={(e) => setEntranceMethod(e.target.value)}>
            <option value="mobile_money">Mobile money</option>
            <option value="bank_deposit">Bank deposit</option>
            <option value="cash_agent">Cash via SACCO agent</option>
          </select>
        </div>
        {error && <p className="text-red-200 mt-4 text-sm">{error}</p>}
        <div className="mt-6">
          <button disabled={isSubmitting} onClick={submit} className="w-full md:w-auto px-6 bg-white text-blue-900 rounded-xl h-12 font-semibold flex items-center justify-center gap-2 hover:bg-blue-50 transition disabled:opacity-60">
            Enrol in SACCO <ArrowRight className="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>
  );
};

export default SaccoSetup;
