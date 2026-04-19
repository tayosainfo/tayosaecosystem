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

/** Match InsForge web SDK: anon auth routes often expect client_type=web for cookies / token behavior. */
function withAuthClientType(path: string): string {
  if (!path.startsWith('/api/v1/auth') || path.includes('client_type=')) {
    return path;
  }
  return path + (path.includes('?') ? '&' : '?') + 'client_type=web';
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const urlPath = withAuthClientType(path);
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
  }) =>
    request<SessionLoginPayload | RegisterPendingResponse>('/api/v1/auth/register', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),

  login: (payload: { identifier: string; password: string }) =>
    request<SessionLoginPayload>('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify(payload),
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
};
