# Supabase Email Templates Configuration

This document provides customized email templates for the Tayosa banking ecosystem.

## Configuration Instructions

1. Navigate to your Supabase dashboard: https://supabase.com/dashboard/project/[YOUR-PROJECT-REF]
2. Go to **Authentication > Email Templates**
3. Select the template you want to customize
4. Copy the template content from this document
5. Paste into the template editor
6. Click **Save** to apply changes

## Available Variables

All templates support these variables:

- `{{ .ConfirmationURL }}` - The confirmation/action URL
- `{{ .Token }}` - The verification/reset token
- `{{ .TokenHash }}` - Hashed version of the token
- `{{ .SiteURL }}` - Your site URL (configured in Auth settings)
- `{{ .Email }}` - The user's email address
- `{{ .RedirectTo }}` - Custom redirect URL (if provided)

## Template 1: Confirm Signup (Email Verification)

**Template Name:** Confirm signup

**Subject Line:**
```
Verify your Tayosa account
```

**HTML Template:**
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Verify Your Email</title>
  <style>
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
      line-height: 1.6;
      color: #333;
      max-width: 600px;
      margin: 0 auto;
      padding: 20px;
      background-color: #f4f4f4;
    }
    .container {
      background-color: #ffffff;
      border-radius: 8px;
      padding: 40px;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    .header {
      text-align: center;
      margin-bottom: 30px;
    }
    .logo {
      font-size: 32px;
      font-weight: bold;
      color: #2563eb;
      margin-bottom: 10px;
    }
    h1 {
      color: #1f2937;
      font-size: 24px;
      margin-bottom: 20px;
    }
    p {
      color: #4b5563;
      margin-bottom: 20px;
    }
    .button {
      display: inline-block;
      padding: 14px 32px;
      background-color: #2563eb;
      color: #ffffff !important;
      text-decoration: none;
      border-radius: 6px;
      font-weight: 600;
      margin: 20px 0;
    }
    .button:hover {
      background-color: #1d4ed8;
    }
    .footer {
      margin-top: 40px;
      padding-top: 20px;
      border-top: 1px solid #e5e7eb;
      font-size: 14px;
      color: #6b7280;
      text-align: center;
    }
    .security-note {
      background-color: #fef3c7;
      border-left: 4px solid #f59e0b;
      padding: 12px;
      margin: 20px 0;
      font-size: 14px;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <div class="logo">Tayosa</div>
      <p style="color: #6b7280; margin: 0;">Banking Made Simple</p>
    </div>
    
    <h1>Welcome to Tayosa!</h1>
    
    <p>Thank you for signing up for Tayosa. We're excited to have you on board!</p>
    
    <p>To complete your registration and start using your account, please verify your email address by clicking the button below:</p>
    
    <div style="text-align: center;">
      <a href="{{ .ConfirmationURL }}" class="button">Verify Email Address</a>
    </div>
    
    <p>Or copy and paste this link into your browser:</p>
    <p style="word-break: break-all; color: #2563eb; font-size: 14px;">{{ .ConfirmationURL }}</p>
    
    <div class="security-note">
      <strong>⏰ This link will expire in 24 hours</strong> for security reasons. If it expires, you can request a new verification email from the login page.
    </div>
    
    <p>If you didn't create an account with Tayosa, you can safely ignore this email.</p>
    
    <div class="footer">
      <p>Need help? Contact us at <a href="mailto:support@tayosa.com" style="color: #2563eb;">support@tayosa.com</a></p>
      <p style="margin-top: 10px;">© 2024 Tayosa. All rights reserved.</p>
    </div>
  </div>
</body>
</html>
```

## Template 2: Reset Password

**Template Name:** Reset password

**Subject Line:**
```
Reset your Tayosa password
```

**HTML Template:**
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Reset Your Password</title>
  <style>
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
      line-height: 1.6;
      color: #333;
      max-width: 600px;
      margin: 0 auto;
      padding: 20px;
      background-color: #f4f4f4;
    }
    .container {
      background-color: #ffffff;
      border-radius: 8px;
      padding: 40px;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    .header {
      text-align: center;
      margin-bottom: 30px;
    }
    .logo {
      font-size: 32px;
      font-weight: bold;
      color: #2563eb;
      margin-bottom: 10px;
    }
    h1 {
      color: #1f2937;
      font-size: 24px;
      margin-bottom: 20px;
    }
    p {
      color: #4b5563;
      margin-bottom: 20px;
    }
    .button {
      display: inline-block;
      padding: 14px 32px;
      background-color: #2563eb;
      color: #ffffff !important;
      text-decoration: none;
      border-radius: 6px;
      font-weight: 600;
      margin: 20px 0;
    }
    .button:hover {
      background-color: #1d4ed8;
    }
    .footer {
      margin-top: 40px;
      padding-top: 20px;
      border-top: 1px solid #e5e7eb;
      font-size: 14px;
      color: #6b7280;
      text-align: center;
    }
    .security-note {
      background-color: #fef3c7;
      border-left: 4px solid #f59e0b;
      padding: 12px;
      margin: 20px 0;
      font-size: 14px;
    }
    .warning {
      background-color: #fee2e2;
      border-left: 4px solid #ef4444;
      padding: 12px;
      margin: 20px 0;
      font-size: 14px;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <div class="logo">Tayosa</div>
      <p style="color: #6b7280; margin: 0;">Banking Made Simple</p>
    </div>
    
    <h1>Password Reset Request</h1>
    
    <p>We received a request to reset the password for your Tayosa account associated with <strong>{{ .Email }}</strong>.</p>
    
    <p>Click the button below to reset your password:</p>
    
    <div style="text-align: center;">
      <a href="{{ .ConfirmationURL }}" class="button">Reset Password</a>
    </div>
    
    <p>Or copy and paste this link into your browser:</p>
    <p style="word-break: break-all; color: #2563eb; font-size: 14px;">{{ .ConfirmationURL }}</p>
    
    <div class="security-note">
      <strong>⏰ This link will expire in 1 hour</strong> for security reasons. If it expires, you can request a new password reset from the login page.
    </div>
    
    <div class="warning">
      <strong>⚠️ Didn't request this?</strong><br>
      If you didn't request a password reset, please ignore this email. Your password will remain unchanged. For security, consider changing your password if you suspect unauthorized access.
    </div>
    
    <div class="footer">
      <p>Need help? Contact us at <a href="mailto:support@tayosa.com" style="color: #2563eb;">support@tayosa.com</a></p>
      <p style="margin-top: 10px;">© 2024 Tayosa. All rights reserved.</p>
    </div>
  </div>
</body>
</html>
```

## Template 3: Magic Link (Optional)

**Template Name:** Magic Link

**Subject Line:**
```
Sign in to Tayosa
```

**HTML Template:**
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Sign In to Tayosa</title>
  <style>
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
      line-height: 1.6;
      color: #333;
      max-width: 600px;
      margin: 0 auto;
      padding: 20px;
      background-color: #f4f4f4;
    }
    .container {
      background-color: #ffffff;
      border-radius: 8px;
      padding: 40px;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    .header {
      text-align: center;
      margin-bottom: 30px;
    }
    .logo {
      font-size: 32px;
      font-weight: bold;
      color: #2563eb;
      margin-bottom: 10px;
    }
    h1 {
      color: #1f2937;
      font-size: 24px;
      margin-bottom: 20px;
    }
    p {
      color: #4b5563;
      margin-bottom: 20px;
    }
    .button {
      display: inline-block;
      padding: 14px 32px;
      background-color: #2563eb;
      color: #ffffff !important;
      text-decoration: none;
      border-radius: 6px;
      font-weight: 600;
      margin: 20px 0;
    }
    .button:hover {
      background-color: #1d4ed8;
    }
    .footer {
      margin-top: 40px;
      padding-top: 20px;
      border-top: 1px solid #e5e7eb;
      font-size: 14px;
      color: #6b7280;
      text-align: center;
    }
    .security-note {
      background-color: #fef3c7;
      border-left: 4px solid #f59e0b;
      padding: 12px;
      margin: 20px 0;
      font-size: 14px;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <div class="logo">Tayosa</div>
      <p style="color: #6b7280; margin: 0;">Banking Made Simple</p>
    </div>
    
    <h1>Sign in to Tayosa</h1>
    
    <p>Click the button below to securely sign in to your Tayosa account:</p>
    
    <div style="text-align: center;">
      <a href="{{ .ConfirmationURL }}" class="button">Sign In</a>
    </div>
    
    <p>Or copy and paste this link into your browser:</p>
    <p style="word-break: break-all; color: #2563eb; font-size: 14px;">{{ .ConfirmationURL }}</p>
    
    <div class="security-note">
      <strong>⏰ This link will expire in 1 hour</strong> for security reasons. If it expires, you can request a new sign-in link from the login page.
    </div>
    
    <p>If you didn't request this sign-in link, you can safely ignore this email.</p>
    
    <div class="footer">
      <p>Need help? Contact us at <a href="mailto:support@tayosa.com" style="color: #2563eb;">support@tayosa.com</a></p>
      <p style="margin-top: 10px;">© 2024 Tayosa. All rights reserved.</p>
    </div>
  </div>
</body>
</html>
```

## Template 4: Email Change Confirmation

**Template Name:** Change Email Address

**Subject Line:**
```
Confirm your new email address
```

**HTML Template:**
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Confirm Email Change</title>
  <style>
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
      line-height: 1.6;
      color: #333;
      max-width: 600px;
      margin: 0 auto;
      padding: 20px;
      background-color: #f4f4f4;
    }
    .container {
      background-color: #ffffff;
      border-radius: 8px;
      padding: 40px;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    .header {
      text-align: center;
      margin-bottom: 30px;
    }
    .logo {
      font-size: 32px;
      font-weight: bold;
      color: #2563eb;
      margin-bottom: 10px;
    }
    h1 {
      color: #1f2937;
      font-size: 24px;
      margin-bottom: 20px;
    }
    p {
      color: #4b5563;
      margin-bottom: 20px;
    }
    .button {
      display: inline-block;
      padding: 14px 32px;
      background-color: #2563eb;
      color: #ffffff !important;
      text-decoration: none;
      border-radius: 6px;
      font-weight: 600;
      margin: 20px 0;
    }
    .button:hover {
      background-color: #1d4ed8;
    }
    .footer {
      margin-top: 40px;
      padding-top: 20px;
      border-top: 1px solid #e5e7eb;
      font-size: 14px;
      color: #6b7280;
      text-align: center;
    }
    .security-note {
      background-color: #fef3c7;
      border-left: 4px solid #f59e0b;
      padding: 12px;
      margin: 20px 0;
      font-size: 14px;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <div class="logo">Tayosa</div>
      <p style="color: #6b7280; margin: 0;">Banking Made Simple</p>
    </div>
    
    <h1>Confirm Your New Email Address</h1>
    
    <p>You recently requested to change the email address for your Tayosa account.</p>
    
    <p>To complete this change, please confirm your new email address by clicking the button below:</p>
    
    <div style="text-align: center;">
      <a href="{{ .ConfirmationURL }}" class="button">Confirm New Email</a>
    </div>
    
    <p>Or copy and paste this link into your browser:</p>
    <p style="word-break: break-all; color: #2563eb; font-size: 14px;">{{ .ConfirmationURL }}</p>
    
    <div class="security-note">
      <strong>⏰ This link will expire in 24 hours</strong> for security reasons.
    </div>
    
    <p>If you didn't request this email change, please contact our support team immediately at <a href="mailto:support@tayosa.com" style="color: #2563eb;">support@tayosa.com</a></p>
    
    <div class="footer">
      <p>Need help? Contact us at <a href="mailto:support@tayosa.com" style="color: #2563eb;">support@tayosa.com</a></p>
      <p style="margin-top: 10px;">© 2024 Tayosa. All rights reserved.</p>
    </div>
  </div>
</body>
</html>
```

## Testing Email Templates

### 1. Test in Supabase Dashboard

1. Navigate to **Authentication > Email Templates**
2. Select a template
3. Click **Send test email**
4. Enter your email address
5. Check your inbox (and spam folder)

### 2. Test Email Delivery

```bash
# Test registration (triggers verification email)
curl -X POST https://[YOUR-PROJECT-REF].supabase.co/auth/v1/signup \
  -H "apikey: [YOUR_ANON_KEY]" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpassword123"
  }'

# Test password reset (triggers reset email)
curl -X POST https://[YOUR-PROJECT-REF].supabase.co/auth/v1/recover \
  -H "apikey: [YOUR_ANON_KEY]" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com"
  }'
```

### 3. Verify Email Content

Check that emails contain:
- [ ] Correct branding (Tayosa logo and colors)
- [ ] Clear call-to-action button
- [ ] Working confirmation URL
- [ ] Expiration time mentioned
- [ ] Security warnings where appropriate
- [ ] Contact information
- [ ] Professional formatting

## SMTP Configuration (Production)

For production, configure a custom SMTP provider for better deliverability.

### Recommended Providers

1. **SendGrid** - Reliable, good free tier
2. **AWS SES** - Cost-effective, scalable
3. **Mailgun** - Developer-friendly
4. **Postmark** - High deliverability

### Configure SMTP in Supabase

1. Navigate to **Project Settings > Auth > SMTP Settings**
2. Enable custom SMTP
3. Enter provider details:

**Example (SendGrid):**
```
SMTP Host: smtp.sendgrid.net
SMTP Port: 587
SMTP User: apikey
SMTP Password: <your-sendgrid-api-key>
Sender Email: noreply@tayosa.com
Sender Name: Tayosa Banking
```

4. Click **Save**
5. Send test email to verify configuration

## Email Deliverability Best Practices

### 1. Domain Authentication

Configure SPF, DKIM, and DMARC records:

**SPF Record:**
```
v=spf1 include:sendgrid.net ~all
```

**DKIM Record:**
```
(Provided by your SMTP provider)
```

**DMARC Record:**
```
v=DMARC1; p=none; rua=mailto:dmarc@tayosa.com
```

### 2. Sender Reputation

- Use a dedicated sending domain (e.g., mail.tayosa.com)
- Warm up your sending IP gradually
- Monitor bounce and complaint rates
- Maintain clean email lists

### 3. Content Best Practices

- Use clear, concise subject lines
- Include plain text version
- Avoid spam trigger words
- Include unsubscribe link (for marketing emails)
- Test emails before sending

## Troubleshooting

### Emails Not Delivered

1. Check SMTP configuration
2. Verify sender email is verified
3. Check spam folder
4. Review Supabase auth logs
5. Test with different email providers

### Broken Links

1. Verify redirect URLs are whitelisted
2. Check Site URL configuration
3. Test confirmation URL format
4. Verify token is not expired

### Styling Issues

1. Test in multiple email clients
2. Use inline CSS (not external stylesheets)
3. Avoid complex layouts
4. Test on mobile devices

## Support

For email template issues:
- Supabase Docs: https://supabase.com/docs/guides/auth/auth-email-templates
- Support: support@tayosa.com

## Checklist

- [ ] All templates customized with Tayosa branding
- [ ] Test emails sent and received
- [ ] Links working correctly
- [ ] Styling displays properly in major email clients
- [ ] SMTP configured for production
- [ ] Domain authentication configured
- [ ] Sender email verified
- [ ] Email deliverability tested
