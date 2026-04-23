# 🎯 Simple Supabase Setup Guide (For Non-Technical Users)

**Don't worry!** This guide will walk you through everything step-by-step with simple instructions. Just follow along! ⏱️ Takes about 20-30 minutes.

---

## 📋 Before You Start

1. Open this link in your browser: https://supabase.com/dashboard/project/[YOUR-PROJECT-REF]
2. Make sure you're logged into your Supabase account
3. You should see your "Tayosa" project dashboard

**Ready? Let's go!** 🚀

---

## ⚡ Setup Steps (Follow in Order)


### 📝 Step 1: Check JWT Settings (2 minutes)

**What you're doing:** Checking how long users can stay logged in.

**Good news!** The default JWT settings in Supabase are already perfect for most apps. Let's just verify them:

1. **Click on the gear icon ⚙️** at the bottom of the left sidebar (this opens "Project Settings")
2. In the left menu, look for **"JWT Keys"** (it should be under the "CONFIGURATION" section)
3. **Click on "JWT Keys"**
4. You'll see information about your JWT tokens

**What you should see:**
- JWT Secret (this is automatically set - don't change it!)
- JWT expiry time (default is usually 3600 seconds = 1 hour)
- Refresh token expiry (default is usually 2592000 seconds = 30 days)

**You don't need to change anything here!** ✅ The defaults are already correct.

**Note:** If you see different numbers, that's okay! The defaults work great for most apps.

✅ **Done!** JWT settings are good to go.

---

### 🌐 Step 2: Set Up Website Links (5 minutes)

**What you're doing:** Telling Supabase which websites are allowed to use your login system.

1. **Click on "Authentication"** in the left sidebar (it has a 🔑 key icon)
2. Look at the menu items that appear under "Authentication"
3. **Click on "URL Configuration"** (it should be in the list)
   - If you don't see it immediately, look for a "CONFIGURATION" section
4. You'll see a box labeled **"Site URL"**
   - Click in the box and type: `https://tayosaecosystem.vercel.app`
   - This is your main website address (your Vercel deployment)
   
   **Why Vercel and not localhost?**
   - Site URL is your "default" redirect URL
   - Email links will use this URL
   - If you use localhost, email links won't work for other people testing your app
   - You can still develop locally by adding localhost to "Redirect URLs" below!

5. Scroll down to **"Redirect URLs"** or **"Additional Redirect URLs"**
   - You'll see a text box or an "Add URL" button
   - Add these URLs (one at a time if there's an "Add" button, or all at once if it's a text box):

**For Development (Current Setup):**
```
https://tayosaecosystem.vercel.app/auth/callback
https://tayosaecosystem.vercel.app/auth/verify-email
https://tayosaecosystem.vercel.app/auth/reset-password
http://localhost:3000/auth/callback
http://localhost:3000/auth/verify-email
http://localhost:3000/auth/reset-password
```

**Important:** By adding BOTH Vercel and localhost URLs to "Redirect URLs", you can:
- ✅ Develop and test locally on `http://localhost:3000`
- ✅ Test on your live deployment at `https://tayosaecosystem.vercel.app`
- ✅ Share the Vercel link with team members for testing
- ✅ Email links work for everyone (they go to Vercel, not localhost)

**For Production (Add Later When Ready):**
```
https://app.tayosa.com/auth/callback
https://app.tayosa.com/auth/verify-email
https://app.tayosa.com/auth/reset-password
```

**For Mobile App (Add Later When Testing):**
```
tayosa://auth/callback
tayosa://auth/verify-email
tayosa://auth/reset-password
```

**Note:** For now, just add the "Development" URLs. You can add the production and mobile URLs later when you're ready to deploy to production or test the mobile app.

6. **Click the "Save" button** at the bottom

✅ **Done!** Your website and app can now talk to Supabase.

---

### ✉️ Step 3: Turn On Email Login (3 minutes)

**What you're doing:** Allowing users to sign up and log in with their email.

1. **Click on "Authentication"** in the left sidebar (🔑 key icon)
2. Look for the menu items under Authentication - you should see:
   - Users
   - OAuth Apps
   - Policies
   - **Sign-in / Providers** (or just **Providers**)
3. **Click on "Sign-in / Providers"** (or "Providers")
4. Look for **"Email"** in the list of providers
5. **Click on "Email"** to expand it
6. You'll see a toggle switch - turn it **ON** (it should turn blue or green)
7. Make sure these settings are turned **ON**:
   - ✅ **"Confirm email"** or **"Enable email confirmations"** - Turn this ON
   - ✅ **"Secure email change"** - Turn this ON (if you see it)
8. Find **"Minimum password length"** or **"Password requirements"**
   - Make sure it's set to at least `8` characters
9. **Click "Save"** if there's a save button

✅ **Done!** Users can now sign up with email and password.

---

### 💌 Step 4: Customize Welcome Email (10 minutes)

**What you're doing:** Making the email that users receive look nice and professional.

#### Part A: Email Verification (Welcome Email)

1. **Click on "Authentication"** in the left sidebar (🔑 key icon)
2. Look for **"Email Templates"** in the menu (it might be under a "Configuration" section)
3. **Click on "Email Templates"**
4. You'll see a list of email types. **Click on "Confirm signup"** or **"Confirm email"**
5. You'll see two boxes:
   
   **Subject Line Box:**
   - Delete what's there
   - Type: `Verify your Tayosa account`
   
   **Body Box (the big one):**
   - Delete everything in the box
   - Copy this entire template and paste it:

```html
<!DOCTYPE html>
<html>
<head>
  <style>
    body { font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; }
    .header { text-align: center; margin-bottom: 30px; }
    .logo { font-size: 32px; font-weight: bold; color: #2563eb; }
    h1 { color: #1f2937; font-size: 24px; }
    .button { display: inline-block; padding: 14px 32px; background-color: #2563eb; color: white; text-decoration: none; border-radius: 6px; margin: 20px 0; }
    .note { background-color: #fef3c7; padding: 12px; margin: 20px 0; }
  </style>
</head>
<body>
  <div class="header">
    <div class="logo">Tayosa</div>
    <p>Banking Made Simple</p>
  </div>
  
  <h1>Welcome to Tayosa!</h1>
  
  <p>Thank you for signing up! Please verify your email address by clicking the button below:</p>
  
  <div style="text-align: center;">
    <a href="{{ .ConfirmationURL }}" class="button">Verify Email Address</a>
  </div>
  
  <div class="note">
    <strong>⏰ This link expires in 24 hours</strong>
  </div>
  
  <p>If you didn't create an account, you can ignore this email.</p>
  
  <p style="color: #6b7280; font-size: 14px;">Need help? Contact us at support@tayosa.com</p>
</body>
</html>
```

5. **Click the "Save" button**

#### Part B: Password Reset Email

1. **Stay on the "Email Templates" page**
2. **Click on "Reset password"** in the list
3. **Subject Line Box:**
   - Delete what's there
   - Type: `Reset your Tayosa password`
   
4. **Body Box:**
   - Delete everything
   - Copy this entire template and paste it:

```html
<!DOCTYPE html>
<html>
<head>
  <style>
    body { font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; }
    .header { text-align: center; margin-bottom: 30px; }
    .logo { font-size: 32px; font-weight: bold; color: #2563eb; }
    h1 { color: #1f2937; font-size: 24px; }
    .button { display: inline-block; padding: 14px 32px; background-color: #2563eb; color: white; text-decoration: none; border-radius: 6px; margin: 20px 0; }
    .note { background-color: #fef3c7; padding: 12px; margin: 20px 0; }
    .warning { background-color: #fee2e2; padding: 12px; margin: 20px 0; }
  </style>
</head>
<body>
  <div class="header">
    <div class="logo">Tayosa</div>
    <p>Banking Made Simple</p>
  </div>
  
  <h1>Password Reset Request</h1>
  
  <p>We received a request to reset your password. Click the button below to reset it:</p>
  
  <div style="text-align: center;">
    <a href="{{ .ConfirmationURL }}" class="button">Reset Password</a>
  </div>
  
  <div class="note">
    <strong>⏰ This link expires in 1 hour</strong>
  </div>
  
  <div class="warning">
    <strong>⚠️ Didn't request this?</strong><br>
    If you didn't request a password reset, you can ignore this email.
  </div>
  
  <p style="color: #6b7280; font-size: 14px;">Need help? Contact us at support@tayosa.com</p>
</body>
</html>
```

5. **Click the "Save" button**

✅ **Done!** Your emails now look professional and branded.

---

### 📁 Step 5: Create a Storage Folder (3 minutes)

**What you're doing:** Creating a place to store user documents (like ID photos).

1. **Click on "Storage"** in the left sidebar (it has a folder icon 📁)
2. **Click the green "New bucket" button**
3. A popup will appear. Fill in:
   - **Name:** Type `collateral_docs`
   - **Public bucket:** Make sure this is **OFF** (unchecked)
   - **File size limit:** Type `50` (this means 50 MB max per file)
4. **Click "Create bucket"**

✅ **Done!** Users can now upload documents securely.

---

### ✅ Step 6: Test That Everything Works (5 minutes)

**What you're doing:** Making sure emails are being sent correctly.

1. **Click on "Authentication"** in the left sidebar
2. **Click on "Email Templates"**
3. **Click on "Confirm signup"**
4. Look for a button that says **"Send test email"** (usually at the top right)
5. **Click it**
6. A box will appear asking for an email address
   - Type your own email address
7. **Click "Send"**
8. **Check your email inbox** (and spam folder!)
   - You should receive a welcome email
   - Make sure it looks good and the button works

**If you got the email:** ✅ Perfect! Everything is working!

**If you didn't get the email:** 
- Check your spam folder
- Wait 2-3 minutes and check again
- Make sure you typed your email correctly

---

## 🎉 Congratulations! You're Done!

You've successfully set up Supabase! Here's what you accomplished:

✅ Set up login times (1 hour sessions)  
✅ Added your website links  
✅ Turned on email login  
✅ Customized welcome and password reset emails  
✅ Created a storage folder for documents  
✅ Tested that emails work  

---

## 🔒 IMPORTANT: One More Step - Security Policies!

**Because you have automatic RLS enabled**, you need to set up security policies so users can access their data.

👉 **Follow this guide next:** `docs/SETUP_RLS_POLICIES.md`

**This takes 10-15 minutes** and is REQUIRED for your app to work properly!

**What it does:**
- Lets users see their own data
- Prevents users from seeing other people's data
- Keeps your database secure

**Don't skip this!** Without it, users won't be able to log in or see their information.

---

## 🤔 What Happens Next?

Now that Supabase is configured, your developer team can:
1. Run the database migration (updates the database)
2. Deploy the updated code (puts the new version live)
3. Test everything works end-to-end

---

## 📞 Need Help?

**If something doesn't look right or you get stuck:**

1. Take a screenshot of what you see
2. Email: support@tayosa.com
3. Or ask your developer team

**Common Questions:**

**Q: What if I make a mistake?**  
A: Don't worry! You can always go back and change any setting. Nothing is permanent.

**Q: Can I test this without affecting real users?**  
A: Yes! The test email feature lets you test without affecting anyone.

**Q: How do I know if I did it right?**  
A: If you received the test email and it looks good, you did it perfectly!

---

## 📋 Quick Checklist

Before you tell your team you're done, make sure:

- [ ] I clicked "Save" after each step
- [ ] I received the test email
- [ ] The test email looks professional (has Tayosa branding)
- [ ] The storage bucket `collateral_docs` was created
- [ ] All the website links were added

**All checked?** You're ready to go! 🚀

---

**⏱️ Total Time:** About 20-30 minutes  
**Difficulty:** Easy - Just follow the steps!  
**Status:** Ready for your development team to deploy! ✅
