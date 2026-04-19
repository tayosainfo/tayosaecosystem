const base = process.env.API_BASE_URL || 'http://localhost:8080';

async function call(path, method, body, token) {
  const res = await fetch(`${base}${path}`, {
    method,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: body ? JSON.stringify(body) : undefined,
  });
  const data = await res.json();
  if (!res.ok) {
    throw new Error(`${path} failed: ${JSON.stringify(data)}`);
  }
  return data;
}

async function run() {
  const register = await call('/api/v1/auth/register', 'POST', {
    fullName: 'Plan Test User',
    phone: '+256700111222',
    email: 'plantest@example.com',
    password: 'secret123',
    nationality: 'UG',
  });
  console.log('register ok', register.user.id);

  const login = await call('/api/v1/auth/login', 'POST', {
    identifier: '+256700111222',
    password: 'secret123',
  });
  const token = login.session.accessToken;
  console.log('phone login ok', token);

  await call('/api/v1/auth/login', 'POST', {
    identifier: 'plantest@example.com',
    password: 'secret123',
  });
  console.log('email login ok');

  await call('/api/v1/onboarding/phase', 'POST', {
    userId: login.user.id,
    phase: 3,
    geo: { parish: 'Abongepach', village: 'VillageAbongepach' },
    membership: { membershipType: 'individual' },
  }, token);
  console.log('onboarding ok');

  const geo = await call('/api/v1/geo?level=district', 'GET');
  console.log('geo districts count', geo.count);
}

run().catch((err) => {
  console.error(err.message);
  process.exit(1);
});

