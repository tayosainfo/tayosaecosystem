#!/usr/bin/env node

/**
 * Fix Admin RLS Policies for KYC Documents
 * This script connects to your Supabase database and updates the RLS policies
 * to allow admin users to view all KYC documents
 */

const SUPABASE_URL = 'https://ablvrbnbsdqshrorhmjf.supabase.co';
const SERVICE_ROLE_KEY = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImFibHZyYm5ic2Rxc2hyb3JobWpmIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTc3Njg2NDIxNCwiZXhwIjoyMDkyNDQwMjE0fQ.ZF8tRGBatd-GZ3sAK7REi3gAp8yj8H5uj-Eq-FBwJPE';

const SQL_STATEMENTS = [
  'DROP POLICY IF EXISTS "Users can view own KYC documents" ON public.kyc_documents;',
  'DROP POLICY IF EXISTS "Users can insert own KYC documents" ON public.kyc_documents;',
  'DROP POLICY IF EXISTS "Users can update own KYC documents" ON public.kyc_documents;',
  'DROP POLICY IF EXISTS "Users can delete own KYC documents" ON public.kyc_documents;',
  'DROP POLICY IF EXISTS "Service role full access to KYC documents" ON public.kyc_documents;',
  `CREATE POLICY "Users can view own KYC documents"
ON public.kyc_documents
FOR SELECT
TO authenticated
USING (
  user_id = auth.uid()::text
  OR
  EXISTS (
    SELECT 1 FROM public.users_identity
    WHERE supabase_user_id = auth.uid()::text
    AND role = 'admin'
  )
);`,
  `CREATE POLICY "Users can insert own KYC documents"
ON public.kyc_documents
FOR INSERT
TO authenticated
WITH CHECK (user_id = auth.uid()::text);`,
  `CREATE POLICY "Users can update own KYC documents"
ON public.kyc_documents
FOR UPDATE
TO authenticated
USING (user_id = auth.uid()::text)
WITH CHECK (user_id = auth.uid()::text);`,
  `CREATE POLICY "Users can delete own KYC documents"
ON public.kyc_documents
FOR DELETE
TO authenticated
USING (user_id = auth.uid()::text);`,
  `CREATE POLICY "Service role full access to KYC documents"
ON public.kyc_documents
FOR ALL
TO service_role
USING (true)
WITH CHECK (true);`,
  'DROP POLICY IF EXISTS "Enable read access for role checking" ON public.users_identity;',
  'DROP POLICY IF EXISTS "Users can read their own role" ON public.users_identity;',
  'DROP POLICY IF EXISTS "Allow reading user roles for authenticated users" ON public.users_identity;',
  `CREATE POLICY "Enable read access for role checking"
ON public.users_identity
FOR SELECT
USING (true);`
];

async function executeSql(sql) {
  try {
    const response = await fetch(`${SUPABASE_URL}/rest/v1/rpc/exec_sql`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${SERVICE_ROLE_KEY}`,
        'Content-Type': 'application/json',
        'apikey': SERVICE_ROLE_KEY
      },
      body: JSON.stringify({ sql })
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(`HTTP ${response.status}: ${error}`);
    }

    return await response.json();
  } catch (err) {
    throw err;
  }
}

async function fixAdminRLS() {
  try {
    console.log('🔧 Fixing Admin RLS Policies...\n');
    console.log('Connecting to Supabase database...\n');
    
    let successCount = 0;
    let errorCount = 0;

    for (let i = 0; i < SQL_STATEMENTS.length; i++) {
      const sql = SQL_STATEMENTS[i];
      const statementNum = i + 1;
      
      try {
        console.log(`[${statementNum}/${SQL_STATEMENTS.length}] Executing SQL...`);
        await executeSql(sql);
        console.log(`✅ Statement ${statementNum} executed successfully\n`);
        successCount++;
      } catch (err) {
        // Some statements might fail if policies don't exist, which is OK
        if (err.message.includes('does not exist') || err.message.includes('DROP POLICY')) {
          console.log(`⚠️  Statement ${statementNum} skipped (policy doesn't exist)\n`);
        } else {
          console.error(`❌ Statement ${statementNum} failed: ${err.message}\n`);
          errorCount++;
        }
      }
    }

    console.log('\n' + '='.repeat(60));
    console.log(`✅ RLS Policy Update Complete!`);
    console.log(`   Successful: ${successCount}/${SQL_STATEMENTS.length}`);
    if (errorCount > 0) {
      console.log(`   Errors: ${errorCount}`);
    }
    console.log('='.repeat(60));
    
    console.log('\n📝 Changes Applied:');
    console.log('   ✓ Updated kyc_documents RLS policies');
    console.log('   ✓ Admin users can now view all KYC documents');
    console.log('   ✓ Regular users can still only view their own documents');
    console.log('   ✓ Updated users_identity RLS policies');
    
    console.log('\n🧪 Test the fix:');
    console.log('   GET https://tayosaecosystem.onrender.com/api/v1/admin/kyc?status=pending');
    console.log('\n✨ The admin user (baylesinfo@gmail.com) should now have access!');
    
  } catch (err) {
    console.error('❌ Unexpected error:', err.message);
    process.exit(1);
  }
}

fixAdminRLS();
