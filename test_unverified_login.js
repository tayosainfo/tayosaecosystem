import fs from 'fs';
import path from 'path';
import { createClient } from '@supabase/supabase-js';

const envContent = fs.readFileSync(path.resolve('./.env'), 'utf-8');
const urlMatch = envContent.match(/VITE_SUPABASE_URL=(.*)/);
const anonKeyMatch = envContent.match(/VITE_SUPABASE_ANON_KEY=(.*)/);
const supabaseUrl = urlMatch ? urlMatch[1].trim() : '';
const anonKey = anonKeyMatch ? anonKeyMatch[1].trim() : '';

if (!supabaseUrl || !anonKey) {
    throw new Error('Missing VITE_SUPABASE_URL or VITE_SUPABASE_ANON_KEY in .env file');
}

const client = createClient(supabaseUrl, anonKey);

async function testUnverifiedLogin() {
    console.log(`--- TEST: Login with known unverified email: test_fail_7@example.com ---`);
    const { data, error } = await client.auth.signInWithPassword({
        email: 'test_fail_7@example.com',
        password: 'password123'
    });
    
    console.log('SignIn response:');
    if (error) {
        console.log('Error:', error);
        console.log('Error message:', error.message);
        console.log('Error status:', error.status);
    }
    if (data) {
        console.log('Data:');
        console.dir(data, { depth: null });
    }
}

testUnverifiedLogin();
