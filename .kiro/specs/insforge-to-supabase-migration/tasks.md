# Implementation Plan: InsForge to Supabase Migration

## Overview

This implementation plan provides a comprehensive, incremental approach to migrating the Tayosa banking ecosystem from InsForge to Supabase. The migration involves frontend refactoring, backend updates, mobile SDK integration, database migrations, test updates, and documentation. Each task is designed to be discrete and testable, with checkpoints to ensure system stability throughout the migration.

## Tasks

- [x] 1. Prepare database migration
  - [x] 1.1 Create database migration file for column renaming
    - Create `db/migrations/011_rename_insforge_to_supabase.sql`
    - Add ALTER TABLE statements to rename `insforge_user_id` to `supabase_user_id`
    - Add ALTER TABLE statements to rename `insforge_login_email` to `supabase_login_email`
    - Drop old index `idx_users_insforge_login_email`
    - Create new index `idx_users_supabase_login_email`
    - _Requirements: 6.1, 6.3, 6.4_
  
  - [ ]* 1.2 Write data integrity verification queries
    - Create SQL queries to check for null `supabase_user_id` values
    - Create SQL queries to check for duplicate emails
    - Document expected results before and after migration
    - _Requirements: 6.4_

- [x] 2. Frontend refactoring - File and import updates
  - [x] 2.1 Rename Supabase client file
    - Rename `src/lib/insforge.ts` to `src/lib/supabase.ts`
    - Update export statement to maintain `supabase` export name
    - _Requirements: 1.1, 1.2_
  
  - [x] 2.2 Update all imports referencing the renamed file
    - Search for all files importing from `src/lib/insforge`
    - Update import paths to `src/lib/supabase`
    - Verify no broken imports remain
    - _Requirements: 1.1_
  
  - [x] 2.3 Update comments and documentation in frontend code
    - Search for comments mentioning "InsForge" in frontend files
    - Replace with "Supabase" where appropriate
    - Update JSDoc comments if present
    - _Requirements: 1.1, 14.1, 14.2_

- [x] 3. Backend refactoring - User service updates
  - [x] 3.1 Update struct field names in user service
    - Rename `InsforgeUserID` to `SupabaseUserID` in Go structs
    - Rename `InsforgeEmail` to `SupabaseLoginEmail` in Go structs
    - Update all references to these fields throughout the codebase
    - _Requirements: 2.1, 2.5_
  
  - [x] 3.2 Update comments in supabase_client.go
    - Replace all "InsForge" references with "Supabase" in comments
    - Update function documentation to reflect Supabase integration
    - Verify error messages reference Supabase correctly
    - _Requirements: 2.8, 13.2_
  
  - [x] 3.3 Update database query field names
    - Update SQL queries to use `supabase_user_id` instead of `insforge_user_id`
    - Update SQL queries to use `supabase_login_email` instead of `insforge_login_email`
    - Verify all database operations use new column names
    - _Requirements: 2.5, 6.3_

- [x] 4. Backend refactoring - Other services token validation
  - [x] 4.1 Implement Supabase token validation in api-gateway-service
    - Add token validation function that calls Supabase `/auth/v1/user` endpoint
    - Use `SUPABASE_URL` and `SUPABASE_ANON_KEY` environment variables
    - Add middleware to validate tokens before routing requests
    - Forward `Authorization` header to downstream services
    - _Requirements: 3.1, 3.2, 12.1, 12.3_
  
  - [x] 4.2 Implement Supabase token validation in affiliate-service
    - Add token validation function using Supabase API
    - Update authentication middleware to use Supabase validation
    - Handle Supabase error responses appropriately
    - _Requirements: 3.1, 3.2, 3.6_
  
  - [x] 4.3 Implement Supabase token validation in audit-log-service
    - Add token validation function using Supabase API
    - Update authentication middleware to use Supabase validation
    - Handle Supabase error responses appropriately
    - _Requirements: 3.1, 3.2, 3.6_
  
  - [x] 4.4 Implement Supabase token validation in fee-service
    - Add token validation function using Supabase API
    - Update authentication middleware to use Supabase validation
    - Handle Supabase error responses appropriately
    - _Requirements: 3.1, 3.2, 3.6_
  
  - [x] 4.5 Implement Supabase token validation in kibiina-service
    - Add token validation function using Supabase API
    - Update authentication middleware to use Supabase validation
    - Handle Supabase error responses appropriately
    - _Requirements: 3.1, 3.2, 3.6_
  
  - [x] 4.6 Implement Supabase token validation in loan-credit-service
    - Add token validation function using Supabase API
    - Update authentication middleware to use Supabase validation
    - Handle Supabase error responses appropriately
    - _Requirements: 3.1, 3.2, 3.6_
  
  - [x] 4.7 Implement Supabase token validation in notification-service
    - Add token validation function using Supabase API
    - Update authentication middleware to use Supabase validation
    - Handle Supabase error responses appropriately
    - _Requirements: 3.1, 3.2, 3.6_
  
  - [x] 4.8 Implement Supabase token validation in object-storage-service
    - Add token validation function using Supabase API
    - Update authentication middleware to use Supabase validation
    - Handle Supabase error responses appropriately
    - _Requirements: 3.1, 3.2, 3.6_

- [x] 5. Checkpoint - Backend services validation
  - Ensure all backend services compile without errors
  - Verify environment variables are correctly referenced
  - Ask the user if questions arise

- [x] 6. Mobile app - Supabase SDK integration
  - [x] 6.1 Add Supabase Flutter SDK dependency
    - Add `supabase_flutter: ^2.0.0` to `pubspec.yaml` dependencies
    - Remove any InsForge SDK dependencies if present
    - Run `flutter pub get` to install dependencies
    - _Requirements: 4.1, 4.2_
  
  - [x] 6.2 Create Supabase client configuration file
    - Create `app/mobile_app/lib/core/network/supabase_client.dart`
    - Implement `SupabaseConfig` class with initialization method
    - Configure Supabase URL and anon key from environment variables
    - Export `client` getter for accessing Supabase instance
    - _Requirements: 4.3_
  
  - [x] 6.3 Update main.dart to initialize Supabase
    - Import Supabase configuration
    - Call `SupabaseConfig.initialize()` before `runApp()`
    - Handle initialization errors appropriately
    - _Requirements: 4.3_
  
  - [x] 6.4 Update API client to integrate Supabase authentication
    - Modify `app/mobile_app/lib/core/network/api_client.dart`
    - Add Supabase client instance to ApiClient class
    - Add Dio interceptor to attach Supabase access token to requests
    - Update authentication methods to use backend endpoints
    - _Requirements: 4.4, 4.5_
  
  - [x] 6.5 Update login screen to use Supabase-aware API client
    - Modify `app/mobile_app/lib/features/auth/presentation/login_screen.dart`
    - Ensure login calls use updated API client
    - Maintain existing UI and user experience
    - Handle Supabase authentication errors
    - _Requirements: 4.6, 4.8_
  
  - [x] 6.6 Update register screen to use Supabase-aware API client
    - Modify `app/mobile_app/lib/features/auth/presentation/register_screen.dart`
    - Ensure registration calls use updated API client
    - Maintain existing UI and user experience
    - Handle Supabase authentication errors
    - _Requirements: 4.6, 4.7_
  
  - [x] 6.7 Update forgot password screen for Supabase integration
    - Modify `app/mobile_app/lib/features/auth/presentation/forgot_password_screen.dart`
    - Update password reset flow to use Supabase
    - Ensure email sending works through Supabase
    - _Requirements: 4.6_
  
  - [x] 6.8 Update reset password screen for Supabase integration
    - Modify `app/mobile_app/lib/features/auth/presentation/reset_password_screen.dart`
    - Update password reset completion to use Supabase
    - Handle Supabase password reset tokens
    - _Requirements: 4.6_

- [x] 7. Checkpoint - Mobile app compilation
  - Ensure Flutter app compiles without errors
  - Verify Supabase SDK is properly initialized
  - Test basic authentication flow manually
  - Ask the user if questions arise

- [x] 8. Test file updates
  - [x] 8.1 Update test_unified_auth_flow.js for Supabase
    - Replace InsForge SDK import with Supabase SDK import
    - Update client initialization to use Supabase configuration
    - Update test assertions to match Supabase response format
    - Verify test connects to Supabase test environment
    - _Requirements: 7.1, 7.2, 7.3_
  
  - [x] 8.2 Update test_unverified_login.js for Supabase
    - Replace InsForge SDK import with Supabase SDK import
    - Update client initialization to use Supabase configuration
    - Update test assertions to match Supabase response format
    - Verify test connects to Supabase test environment
    - _Requirements: 7.1, 7.2, 7.3_
  
  - [ ]* 8.3 Run integration tests to verify authentication flows
    - Execute `test_unified_auth_flow.js` against Supabase test environment
    - Execute `test_unverified_login.js` against Supabase test environment
    - Verify all test cases pass
    - Document any test failures for investigation
    - _Requirements: 7.4, 7.5_

- [x] 9. Environment configuration updates
  - [x] 9.1 Update .env.example with Supabase variables
    - Verify `VITE_SUPABASE_URL` is documented
    - Verify `VITE_SUPABASE_ANON_KEY` is documented
    - Verify `SUPABASE_SERVICE_ROLE_KEY` is documented
    - Add comments explaining each variable's purpose
    - Remove any InsForge-specific variables
    - _Requirements: 5.1, 5.2, 5.3, 5.4_
  
  - [x] 9.2 Create environment variable documentation
    - Document all required Supabase environment variables
    - Provide instructions for obtaining Supabase credentials
    - Document Supabase project URL format
    - Include examples for development and production
    - _Requirements: 5.5, 5.6_
  
  - [x] 9.3 Update mobile app environment configuration
    - Document `--dart-define` flags for Flutter builds
    - Provide example build commands with Supabase configuration
    - Document how to configure Supabase URL and anon key
    - _Requirements: 5.1, 5.2_

- [x] 10. Documentation updates
  - [x] 10.1 Update README.md with Supabase setup instructions
    - Add section on Supabase project setup
    - Document how to obtain Supabase URL and keys
    - Provide step-by-step migration instructions
    - Include troubleshooting section for common issues
    - _Requirements: 14.1, 14.5_
  
  - [x] 10.2 Create migration runbook
    - Document pre-migration checklist
    - Provide step-by-step deployment instructions
    - Document database migration execution steps
    - Include rollback procedures
    - Document post-migration verification steps
    - _Requirements: 14.3_
  
  - [x] 10.3 Update API documentation
    - Update authentication endpoint documentation
    - Document Supabase token format and validation
    - Update error response documentation
    - Document session management with Supabase
    - _Requirements: 14.4_
  
  - [x] 10.4 Document Supabase project configuration
    - Document email template configuration in Supabase
    - Document RLS policy setup if required
    - Document OAuth provider configuration
    - Document allowed redirect URLs
    - _Requirements: 5.7, 14.5_

- [x] 11. Deployment preparation
  - [x] 11.1 Verify Supabase project configuration
    - Confirm Supabase project URL is correct
    - Verify anon key and service role key are available
    - Test email delivery in Supabase dashboard
    - Verify database connection string works
    - _Requirements: 5.6, 8.2_
  
  - [x] 11.2 Prepare database migration execution plan
    - Document migration file execution order
    - Create backup plan for existing data
    - Document expected migration duration
    - Prepare rollback SQL statements
    - _Requirements: 6.1, 6.2, 6.4_
  
  - [x] 11.3 Update deployment scripts for Supabase
    - Update CI/CD pipeline with Supabase environment variables
    - Update Docker configurations if applicable
    - Update deployment documentation
    - Test deployment process in staging environment
    - _Requirements: 18.1, 18.2, 18.3_
  
  - [x] 11.4 Configure Supabase email templates
    - Customize email verification template in Supabase dashboard
    - Customize password reset template in Supabase dashboard
    - Test email delivery for verification and reset flows
    - Configure SMTP settings if using custom email provider
    - _Requirements: 8.1, 8.2, 9.1, 9.2_
  
  - [x] 11.5 Configure Supabase authentication settings
    - Enable email confirmations in Supabase dashboard
    - Configure allowed redirect URLs for email links
    - Set up OAuth providers if required
    - Configure session timeout settings
    - _Requirements: 8.3, 10.1, 10.4_

- [x] 12. Final checkpoint - Pre-deployment verification
  - Ensure all code changes are committed
  - Verify all tests pass in test environment
  - Confirm Supabase project is fully configured
  - Review deployment checklist
  - Ask the user if ready to proceed with deployment

## Notes

- Tasks marked with `*` are optional and can be skipped for faster implementation
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation and allow for user feedback
- Database migration should be executed carefully with backup procedures in place
- Mobile app changes require testing on both iOS and Android platforms
- All environment variables must be configured before deployment
- Supabase email templates should be customized before enabling email verification
- The migration maintains backward compatibility with existing API contracts
- No InsForge SDK dependencies remain after migration completion
