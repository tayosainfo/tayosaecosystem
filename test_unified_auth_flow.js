import { createClient } from '@supabase/supabase-js';

// Supabase configuration
const supabaseUrl = process.env.VITE_SUPABASE_URL || 'https://ablvrbnbsdqshrorhmjf.supabase.co';
const supabaseAnonKey = process.env.VITE_SUPABASE_ANON_KEY || '';

if (!supabaseAnonKey) {
  console.error('Error: VITE_SUPABASE_ANON_KEY environment variable is required');
  process.exit(1);
}

// Initialize Supabase client
const supabase = createClient(supabaseUrl, supabaseAnonKey);

// Backend API base URL
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
  console.log('=== Testing Unified Auth Flow with Supabase ===\n');
  
  // Test 1: Verify Supabase connection
  console.log('1. Verifying Supabase connection...');
  const { data: healthCheck, error: healthError } = await supabase.auth.getSession();
  if (healthError && healthError.message !== 'Auth session missing!') {
    throw new Error(`Supabase connection failed: ${healthError.message}`);
  }
  console.log('✓ Supabase connection verified\n');

  // Test 2: Register new user via backend (backend handles Supabase signup)
  console.log('2. Registering new user via backend...');
  const register = await call('/api/v1/auth/register', 'POST', {
    fullName: 'Plan Test User',
    phone: '+256700111222',
    email: 'plantest@example.com',
    password: 'secret123',
    nationality: 'UG',
  });
  
  // Verify Supabase-compatible response format
  if (!register.user || !register.user.id) {
    throw new Error('Registration response missing user.id');
  }
  if (!register.session || !register.session.accessToken) {
    throw new Error('Registration response missing session.accessToken');
  }
  console.log('✓ Registration successful, user ID:', register.user.id);
  console.log('✓ Session token received (Supabase JWT)\n');

  // Test 3: Login with phone via backend
  console.log('3. Testing phone login...');
  const login = await call('/api/v1/auth/login', 'POST', {
    identifier: '+256700111222',
    password: 'secret123',
  });
  const token = login.session.accessToken;
  
  // Verify token is a valid JWT format (Supabase tokens are JWTs)
  const tokenParts = token.split('.');
  if (tokenParts.length !== 3) {
    throw new Error('Access token is not a valid JWT format');
  }
  console.log('✓ Phone login successful');
  console.log('✓ JWT token format verified\n');

  // Test 4: Login with email via backend
  console.log('4. Testing email login...');
  await call('/api/v1/auth/login', 'POST', {
    identifier: 'plantest@example.com',
    password: 'secret123',
  });
  console.log('✓ Email login successful\n');

  // Test 5: Authenticated request with Supabase token
  console.log('5. Testing authenticated request with Supabase token...');
  await call('/api/v1/onboarding/phase', 'POST', {
    userId: login.user.id,
    phase: 3,
    geo: { parish: 'Abongepach', village: 'VillageAbongepach' },
    membership: { membershipType: 'individual' },
  }, token);
  console.log('✓ Onboarding request successful with Supabase token\n');

  // Test 6: Public endpoint (no auth required)
  console.log('6. Testing public endpoint...');
  const geo = await call('/api/v1/geo?level=district', 'GET');
  if (!geo.count || geo.count === 0) {
    throw new Error('Geo districts query returned no results');
  }
  console.log('✓ Geo districts query successful, count:', geo.count);
  
  console.log('\n=== All tests passed! ===');
}

run().catch((err) => {
  console.error('\n❌ Test failed:', err.message);
  process.exit(1);
});

