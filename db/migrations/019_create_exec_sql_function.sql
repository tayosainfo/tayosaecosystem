-- Migration: Create exec_sql function for executing arbitrary SQL
-- Description: Allows executing raw SQL statements via RPC
-- Date: 2026-04-25

-- Create a function that can execute arbitrary SQL
-- This is only accessible to the service role for security
CREATE OR REPLACE FUNCTION public.exec_sql(sql text)
RETURNS void
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
BEGIN
  EXECUTE sql;
END;
$$;

-- Grant execute permission to service role only
GRANT EXECUTE ON FUNCTION public.exec_sql(text) TO service_role;

-- Add comment for documentation
COMMENT ON FUNCTION public.exec_sql(text) IS 'Execute arbitrary SQL statements (service role only)';
