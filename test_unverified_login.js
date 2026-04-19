import fs from 'fs';
import path from 'path';
import { createClient } from '@insforge/sdk';

const envContent = fs.readFileSync(path.resolve('./.env'), 'utf-8');
const anonKeyMatch = envContent.match(/VITE_INSFORGE_ANON_KEY=(.*)/);
const anonKey = anonKeyMatch ? anonKeyMatch[1].trim() : '';

const client = createClient('https://74qj9u5z.us-east.insforge.app', anonKey);

async function testUnverifiedLogin() {
    console.log(`--- TEST: Login with known unverified email: test_fail_7@example.com ---`);
    const res = await client.auth.signInWithPassword({
        email: 'test_fail_7@example.com',
        password: 'password123'
    });
    console.log('SignIn raw response:');
    console.dir(res, { depth: null });
}

testUnverifiedLogin();
