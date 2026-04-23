const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export type SessionLoginPayload = {
  session: { accessToken: string; userId: string; refreshToken?: string; csrfToken?: string };
  user: { id: string; fullName: string; phoneE164: string; contactEmail?: string };
};

export type RegisterPendingResponse = {
  requireEmailVerification: boolean;
  pendingLocalProfile: boolean;
  email: string;
  message?: string;
};

export type MePayload = {
  user: SessionLoginPayload['user'] & { contactEmailVerified?: boolean };
  kyc?: { status?: string; submittedAt?: string; reviewedAt?: string };
  sacco?: { status?: string; parish?: string };
  kibiina?: { action?: string };
  shares?: { balanceUnits?: number };
  referralCode?: string;
  featureAccess?: { canTransact?: boolean; canJoinKibiina?: boolean };
};


/** Thrown on non-2xx API responses; includes parsed JSON for flags like requireEmailVerification. */
export class PlatformApiError extends Error {
  readonly status: number;
  readonly body: Record<string, unknown>;

  constructor(message: string, status: number, body: Record<string, unknown>) {
    super(message);
    this.name = 'PlatformApiError';
    this.status = status;
    this.body = body;
  }
}

export function isPlatformApiError(err: unknown): err is PlatformApiError {
  return err instanceof PlatformApiError;
}

function isRecord(v: unknown): v is Record<string, unknown> {
  return typeof v === 'object' && v !== null;
}

/** Match Supabase web SDK: anon auth routes often expect client_type=web for cookies / token behavior. */
function withAuthClientType(path: string): string {
  if (!path.startsWith('/api/v1/auth') || path.includes('client_type=')) {
    return path;
  }
  return path + (path.includes('?') ? '&' : '?') + 'client_type=web';
}

function withoutClientType(path: string): string {
  if (!path.includes('client_type=')) {
    return path;
  }
  const [base, query = ''] = path.split('?');
  const params = new URLSearchParams(query);
  params.delete('client_type');
  const next = params.toString();
  return next ? `${base}?${next}` : base;
}

async function request<T>(path: string, init?: RequestInit, attachWebClientType = true): Promise<T> {
  const urlPath = attachWebClientType ? withAuthClientType(path) : path;
  const response = await fetch(`${API_BASE_URL}${urlPath}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers || {}),
    },
    ...init,
  });
  let payload: unknown = {};
  try {
    payload = await response.json();
  } catch {
    payload = {};
  }
  if (!response.ok) {
    const rec = isRecord(payload) ? payload : {};
    const err = String(rec.error ?? 'Request failed');
    const hint = rec.hint != null ? String(rec.hint) : '';
    const msg = hint ? `${err} ${hint}` : err;
    throw new PlatformApiError(msg, response.status, rec as Record<string, unknown>);
  }
  return payload as T;
}

export const platformApi = {
  register: (payload: {
    fullName: string;
    phone: string;
    email?: string;
    password: string;
    dateOfBirth?: string;
    nationality?: string;
    referralCode?: string;
    termsAccepted: boolean;
    privacyAccepted: boolean;
    termsVersion?: string;
    privacyVersion?: string;
  }) =>
    request<SessionLoginPayload | RegisterPendingResponse>('/api/v1/auth/register', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  login: (payload: { identifier: string; password: string }) =>
    request<SessionLoginPayload>('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify(payload),
    }).catch(async (err: unknown) => {
      // Some local auth deployments reject client_type=web.
      // Retry once without the query param so login remains compatible.
      if (!(err instanceof PlatformApiError) || err.status !== 401) {
        throw err;
      }
      const fallbackPath = withoutClientType(withAuthClientType('/api/v1/auth/login'));
      return await request<SessionLoginPayload>(fallbackPath, {
        method: 'POST',
        body: JSON.stringify(payload),
      }, false);
    }),

  resendVerificationEmail: (payload: { email: string }) =>
    request<{ success: boolean; message: string }>('/api/v1/auth/resend-verification-email', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  verifyEmail: (payload: {
    email: string;
    otp: string;
    fullName?: string;
    phone?: string;
    dateOfBirth?: string;
    nationality?: string;
  }) =>
    request<SessionLoginPayload>('/api/v1/auth/verify-email', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  sendResetPasswordEmail: (payload: { email: string }) =>
    request<{ success: boolean; message: string }>('/api/v1/auth/send-reset-password-email', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  exchangeResetPasswordToken: (payload: { email: string; code: string }) =>
    request<{ token: string; expiresAt: string }>('/api/v1/auth/exchange-reset-password-token', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  resetPassword: (payload: { newPassword: string; otp: string }) =>
    request<{ message: string }>('/api/v1/auth/reset-password', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  updateOnboardingPhase: (
    token: string,
    payload: {
      userId: string;
      phase: number;
      kyc?: Record<string, unknown>;
      membership?: Record<string, unknown>;
      kibiina?: Record<string, unknown>;
      referralCode?: string;
      geo: Record<string, string>;
    },
  ) =>
    request('/api/v1/onboarding/phase', {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify(payload),
    }),

  getMe: (token: string) =>
    request<MePayload>('/api/v1/users/me', {
      method: 'GET',
      headers: { Authorization: `Bearer ${token}` },
    }),

  submitKYC: (token: string, payload: Record<string, unknown>) =>
    request('/api/v1/onboarding/kyc', {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify(payload),
    }),

  submitSacco: (token: string, payload: Record<string, unknown>) =>
    request('/api/v1/onboarding/sacco', {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify(payload),
    }),

  submitKibiina: (token: string, payload: Record<string, unknown>) =>
    request('/api/v1/onboarding/kibiina', {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify(payload),
    }),

  adminListKyc: (token: string, adminSecret: string, status = 'pending') =>
    request<{ items: any[]; count: number }>(`/api/v1/admin/kyc?status=${encodeURIComponent(status)}`, {
      method: 'GET',
      headers: { Authorization: `Bearer ${token}`, 'X-Admin-Secret': adminSecret },
    }, false),

  adminDecideKyc: (token: string, adminSecret: string, userId: string, decision: 'approved' | 'rejected', reviewNote?: string) =>
    request(`/api/v1/admin/kyc?userId=${encodeURIComponent(userId)}`, {
      method: 'PATCH',
      headers: { Authorization: `Bearer ${token}`, 'X-Admin-Secret': adminSecret },
      body: JSON.stringify({ status: decision, reviewNote: reviewNote || '', reviewedBy: 'admin_panel' }),
    }, false),

  adminGetFees: (token: string, adminSecret: string) =>
    request<{ key: string; value: Record<string, unknown> }>('/api/v1/admin/settings?key=fees', {
      method: 'GET',
      headers: { Authorization: `Bearer ${token}`, 'X-Admin-Secret': adminSecret },
    }, false),

  adminSetFees: (token: string, adminSecret: string, value: Record<string, unknown>) =>
    request('/api/v1/admin/settings?key=fees', {
      method: 'PATCH',
      headers: { Authorization: `Bearer ${token}`, 'X-Admin-Secret': adminSecret },
      body: JSON.stringify(value),
    }, false),

  /**
   * Upload a file to the object-storage-service, which proxies the bytes
   * directly to Supabase storage. Returns the storage key for use in KYC
   * and other document references.
   *
   * NOTE: Do NOT pass Content-Type in the headers — the browser must set it
   * automatically so the multipart boundary is included.
   */
  uploadFile: async (token: string, file: File, category = 'kyc'): Promise<{ key: string }> => {
    const form = new FormData();
    form.append('file', file);
    form.append('category', category);
    const response = await fetch(`${API_BASE_URL}/api/v1/storage/upload`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: form,
    });
    let payload: unknown = {};
    try { payload = await response.json(); } catch { payload = {}; }
    if (!response.ok) {
      const rec = isRecord(payload) ? payload : {};
      throw new PlatformApiError(
        String(rec.error ?? 'File upload failed'),
        response.status,
        rec as Record<string, unknown>,
      );
    }
    return payload as { key: string };
  },
};
