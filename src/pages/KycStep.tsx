import React, { useMemo, useState } from 'react';
import { ArrowRight } from 'lucide-react';
import { platformApi } from '../lib/platformApi';

const KycStep: React.FC = () => {
  const [idType, setIdType] = useState('National ID');
  const [idNumber, setIdNumber] = useState('');
  const [dobDay, setDobDay] = useState('');
  const [dobMonth, setDobMonth] = useState('');
  const [dobYear, setDobYear] = useState('');
  const [gender, setGender] = useState('female');
  const [occupationStatus, setOccupationStatus] = useState('self_employed');
  const [nokName, setNokName] = useState('');
  const [nokRelationship, setNokRelationship] = useState('spouse');
  const [nokPhone, setNokPhone] = useState('');
  const [nokEmail, setNokEmail] = useState('');
  const [sourceOfFunds, setSourceOfFunds] = useState('');
  const [idFrontKey, setIdFrontKey] = useState('');
  const [idBackKey, setIdBackKey] = useState('');
  const [selfieKey, setSelfieKey] = useState('');
  const [uploadingFront, setUploadingFront] = useState(false);
  const [uploadingBack, setUploadingBack] = useState(false);
  const [uploadingSelfie, setUploadingSelfie] = useState(false);
  const [pep, setPep] = useState(false);
  const [disclosures, setDisclosures] = useState('none');
  const [disclosuresOther, setDisclosuresOther] = useState('');
  const [error, setError] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const isSessionError = (message: string): boolean => {
    const msg = message.toLowerCase();
    return msg.includes('invalid or expired session') || msg.includes('missing bearer token');
  };

  const clearSessionAndReturnToLogin = () => {
    sessionStorage.removeItem('auth_token');
    sessionStorage.removeItem('auth_user');
    window.location.href = '/';
  };

  const years = useMemo(() => {
    const y = new Date().getFullYear();
    return Array.from({ length: 100 }, (_, i) => String(y - i));
  }, []);
  const months = useMemo(() => Array.from({ length: 12 }, (_, i) => String(i + 1).padStart(2, '0')), []);
  const days = useMemo(() => Array.from({ length: 31 }, (_, i) => String(i + 1).padStart(2, '0')), []);

  const uploadFileAndGetKey = async (file: File): Promise<string> => {
    const token = sessionStorage.getItem('auth_token');
    if (!token) {
      throw new Error('Session expired. Please sign in again.');
    }
    const result = await platformApi.uploadFile(token, file, 'kyc');
    if (!result.key) {
      throw new Error('Upload completed but storage key was missing');
    }
    return result.key;
  };

  const onUpload = async (
    file: File | undefined,
    setBusy: (v: boolean) => void,
    setKey: (k: string) => void,
  ) => {
    if (!file) return;
    setError('');
    setBusy(true);
    try {
      const key = await uploadFileAndGetKey(file);
      if (!key) {
        throw new Error('Upload completed but storage key was missing');
      }
      setKey(key);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to upload file');
    } finally {
      setBusy(false);
    }
  };

  const submit = async () => {
    const dob = `${dobYear}-${dobMonth}-${dobDay}`;
    const token = sessionStorage.getItem('auth_token');
    if (!token) {
      setError('Session expired. Please sign in again.');
      return;
    }
    if (!dobDay || !dobMonth || !dobYear) {
      setError('Please select date of birth');
      return;
    }
    if (!idFrontKey || !idBackKey || !selfieKey) {
      setError('Please upload ID front, ID back, and selfie photo');
      return;
    }
    setSubmitting(true);
    try {
      await platformApi.submitKYC(token, {
        dateOfBirth: dob,
        gender,
        nationality: 'UG',
        occupationStatus,
        idType,
        idNumber,
        idDocumentFrontKey: idFrontKey,
        idDocumentBackKey: idBackKey,
        selfieKey,
        nokFullName: nokName,
        nokRelationship,
        nokPhone,
        nokEmail: nokEmail || undefined,
        sourceOfFunds,
        pepStatus: pep,
        saccoMembershipDisclosures: disclosures === 'other' ? disclosuresOther : disclosures,
      });
      window.location.href = '/home';
    } catch (e) {
      const message = e instanceof Error ? e.message : 'Failed to submit KYC';
      if (isSessionError(message)) {
        clearSessionAndReturnToLogin();
        return;
      }
      setError(message);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-indigo-900 py-10 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto bg-white/10 backdrop-blur-md rounded-2xl border border-white/20 p-8">
        <div className="text-white mb-6">
          <p className="text-sm text-blue-100">Step 2 of 4 - KYC</p>
          <h1 className="text-3xl font-bold">Verify your identity</h1>
          <p className="text-blue-100">This unlocks savings, borrowing, and withdrawals.</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-white mb-2">ID type</label>
            <select value={idType} onChange={(e) => setIdType(e.target.value)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95">
              <option>National ID</option>
              <option>Passport</option>
              <option>Driver License</option>
              <option>Voter Card</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-white mb-2">ID number</label>
            <input value={idNumber} onChange={(e) => setIdNumber(e.target.value)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95" placeholder="CF92001003..." />
          </div>
          <div className="md:col-span-2">
            <label className="block text-sm font-medium text-white mb-2">Date of birth</label>
            <div className="grid grid-cols-3 gap-2">
              <select value={dobDay} onChange={(e) => setDobDay(e.target.value)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95">
                <option value="">Day</option>
                {days.map((d) => <option key={d} value={d}>{d}</option>)}
              </select>
              <select value={dobMonth} onChange={(e) => setDobMonth(e.target.value)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95">
                <option value="">Month</option>
                {months.map((m) => <option key={m} value={m}>{m}</option>)}
              </select>
              <select value={dobYear} onChange={(e) => setDobYear(e.target.value)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95">
                <option value="">Year</option>
                {years.map((y) => <option key={y} value={y}>{y}</option>)}
              </select>
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-white mb-2">Gender</label>
            <select value={gender} onChange={(e) => setGender(e.target.value)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95">
              <option value="female">Female</option>
              <option value="male">Male</option>
              <option value="other">Other</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-white mb-2">Occupation / employment status</label>
            <select value={occupationStatus} onChange={(e) => setOccupationStatus(e.target.value)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95">
              <option value="employed">Employed</option>
              <option value="self_employed">Self-employed</option>
              <option value="student">Student</option>
              <option value="farmer">Farmer</option>
              <option value="unemployed">Unemployed</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-white mb-2">Source of funds</label>
            <select value={sourceOfFunds} onChange={(e) => setSourceOfFunds(e.target.value)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95">
              <option value="">Select source</option>
              <option value="salary">Salary</option>
              <option value="business">Business</option>
              <option value="farming">Farming</option>
              <option value="remittance">Remittance</option>
              <option value="other">Other</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-white mb-2">Next of kin full name</label>
            <input value={nokName} onChange={(e) => setNokName(e.target.value)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95" />
          </div>
          <div>
            <label className="block text-sm font-medium text-white mb-2">Next of kin relationship</label>
            <select value={nokRelationship} onChange={(e) => setNokRelationship(e.target.value)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95">
              <option value="spouse">Spouse</option>
              <option value="parent">Parent</option>
              <option value="sibling">Sibling</option>
              <option value="child">Child</option>
              <option value="relative">Relative</option>
              <option value="friend">Friend</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-white mb-2">Next of kin phone</label>
            <input value={nokPhone} onChange={(e) => setNokPhone(e.target.value)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95" />
          </div>
          <div>
            <label className="block text-sm font-medium text-white mb-2">Next of kin email (optional)</label>
            <input type="email" value={nokEmail} onChange={(e) => setNokEmail(e.target.value)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95" />
          </div>
          <div>
            <label className="block text-sm font-medium text-white mb-2">ID front photo</label>
            <input type="file" accept="image/*,.pdf" capture="environment" onChange={(e) => onUpload(e.target.files?.[0], setUploadingFront, setIdFrontKey)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95 file:mr-3 file:rounded-md file:border-0 file:bg-blue-100 file:px-3 file:py-2" />
            <p className="text-xs text-blue-100 mt-1">{uploadingFront ? 'Uploading...' : idFrontKey ? 'Uploaded' : 'Upload the front side of your ID'}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-white mb-2">ID back photo</label>
            <input type="file" accept="image/*,.pdf" capture="environment" onChange={(e) => onUpload(e.target.files?.[0], setUploadingBack, setIdBackKey)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95 file:mr-3 file:rounded-md file:border-0 file:bg-blue-100 file:px-3 file:py-2" />
            <p className="text-xs text-blue-100 mt-1">{uploadingBack ? 'Uploading...' : idBackKey ? 'Uploaded' : 'Upload the back side of your ID'}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-white mb-2">Selfie / live photo</label>
            <input type="file" accept="image/*" capture="user" onChange={(e) => onUpload(e.target.files?.[0], setUploadingSelfie, setSelfieKey)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95 file:mr-3 file:rounded-md file:border-0 file:bg-blue-100 file:px-3 file:py-2" />
            <p className="text-xs text-blue-100 mt-1">{uploadingSelfie ? 'Uploading...' : selfieKey ? 'Uploaded' : 'Take or upload your selfie'}</p>
          </div>
          <div className="md:col-span-2">
            <label className="block text-sm font-medium text-white mb-2">SACCO membership disclosures</label>
            <select value={disclosures} onChange={(e) => setDisclosures(e.target.value)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95">
              <option value="none">None</option>
              <option value="existing_membership">I belong to another SACCO</option>
              <option value="past_default">I have had past loan default(s)</option>
              <option value="other">Other</option>
            </select>
            {disclosures === 'other' && (
              <input value={disclosuresOther} onChange={(e) => setDisclosuresOther(e.target.value)} className="w-full border border-white/30 rounded-xl h-12 px-3 bg-white/95 mt-2" placeholder="Please describe" />
            )}
          </div>
          <div className="md:col-span-2 text-white text-sm">
            <label className="flex items-center gap-2">
              <input type="checkbox" checked={pep} onChange={(e) => setPep(e.target.checked)} />
              Politically exposed person (PEP)
            </label>
          </div>
        </div>
        {error && <p className="text-red-200 mt-4 text-sm">{error}</p>}
        <div className="mt-6">
          <button disabled={submitting || uploadingFront || uploadingBack || uploadingSelfie} onClick={submit} className="w-full md:w-auto px-6 bg-white text-blue-900 rounded-xl h-12 font-semibold flex items-center justify-center gap-2 hover:bg-blue-50 transition disabled:opacity-60">
            {submitting ? 'Submitting...' : 'Continue'} <ArrowRight className="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>
  );
};

export default KycStep;
