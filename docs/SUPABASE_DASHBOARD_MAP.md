# 🗺️ Supabase Dashboard - Where to Find Everything

**Dashboard URL:** https://supabase.com/dashboard/project/[YOUR-PROJECT-REF]

This guide shows you exactly where to find each setting in your Supabase dashboard.

---

## 📍 Left Sidebar Navigation

When you open your Supabase dashboard, you'll see a left sidebar with these main sections:

```
┌─────────────────────────────┐
│ 🏠 Home                     │
│ 📊 Table Editor             │
│ 🔑 Authentication           │  ← You'll use this a lot!
│ 🗄️  Database                │
│ 📁 Storage                  │  ← You'll use this too!
│ 🔧 Edge Functions           │
│ </> SQL Editor              │  ← For RLS policies!
│ ⚙️  Project Settings        │  ← For JWT settings!
└─────────────────────────────┘
```

---

## 🔑 Authentication Section

**Click on "Authentication" in the left sidebar**

You'll see these sub-menu items:

```
Authentication
├── 👥 Users                    (View registered users)
├── 🔐 OAuth Apps               (Not needed for this setup)
├── 🛡️  Policies                (View RLS policies after setup)
├── 🌐 URL Configuration        ← Step 2: Add website URLs
├── 📧 Sign-in / Providers      ← Step 3: Enable email login
└── 💌 Email Templates          ← Step 4: Customize emails
```

### What You'll Do Here:

**Step 2: URL Configuration**
- Click: Authentication → URL Configuration
- Add your website URLs

**Step 3: Sign-in / Providers**
- Click: Authentication → Sign-in / Providers
- Enable email login

**Step 4: Email Templates**
- Click: Authentication → Email Templates
- Customize welcome and password reset emails

---

## ⚙️ Project Settings Section

**Click on the gear icon ⚙️ at the bottom of the left sidebar**

You'll see these sub-menu items:

```
Project Settings
├── ⚙️  General                 (Project name, region)
├── 🗄️  Database                (Connection strings)
├── 🔐 API                      (API keys - already configured)
├── 🔑 JWT Keys                 ← Step 1: Verify JWT settings
├── 🔗 Integrations             (Not needed)
└── 📊 Usage                    (Monitor usage)
```

### What You'll Do Here:

**Step 1: JWT Keys**
- Click: Project Settings → JWT Keys
- Just verify the settings (don't change anything!)

---

## 📁 Storage Section

**Click on "Storage" in the left sidebar**

You'll see:

```
Storage
├── 📁 Buckets                  ← Step 5: Create storage bucket
└── ⚙️  Settings                (Not needed)
```

### What You'll Do Here:

**Step 5: Create Storage Bucket**
- Click: Storage (in left sidebar)
- Click: "New bucket" button
- Name it: `collateral_docs`

---

## </> SQL Editor Section

**Click on "SQL Editor" in the left sidebar (has a `</>` icon)**

You'll see:

```
SQL Editor
├── 📝 New query                ← You'll use this!
├── 📚 Templates                (Pre-made queries)
└── 📜 Query history            (Your past queries)
```

### What You'll Do Here:

**RLS Policies Setup (CRITICAL!)**
- Click: SQL Editor (in left sidebar)
- You'll see a big text box
- Copy and paste the SQL code blocks from `SETUP_RLS_POLICIES.md`
- Click "Run" after each one

---

## 🗄️ Database Section (For Verification)

**Click on "Database" in the left sidebar**

You'll see:

```
Database
├── 📊 Tables                   (View your tables)
├── 🛡️  Policies                ← Verify RLS policies here!
├── 🔗 Replication              (Not needed)
├── 🔧 Functions                (Not needed)
├── 🔌 Extensions               (Not needed)
└── 🔄 Backups                  (Good to check later)
```

### What You'll Do Here:

**After RLS Setup - Verification:**
- Click: Database → Policies
- You should see ~13 policies listed
- This confirms your RLS setup worked!

---

## 📋 Step-by-Step Navigation Guide

### For Quick Supabase Setup:

**Step 1: JWT Settings**
```
1. Click: ⚙️ (gear icon at bottom of left sidebar)
2. Click: JWT Keys
3. Verify settings (don't change anything)
```

**Step 2: URL Configuration**
```
1. Click: 🔑 Authentication (in left sidebar)
2. Click: URL Configuration
3. Add your website URLs
```

**Step 3: Email Login**
```
1. Click: 🔑 Authentication (in left sidebar)
2. Click: Sign-in / Providers
3. Enable email login
```

**Step 4: Email Templates**
```
1. Click: 🔑 Authentication (in left sidebar)
2. Click: Email Templates
3. Customize welcome and password reset emails
```

**Step 5: Storage Bucket**
```
1. Click: 📁 Storage (in left sidebar)
2. Click: "New bucket" button
3. Name: collateral_docs
```

**Step 6: Test Email**
```
1. Click: 🔑 Authentication (in left sidebar)
2. Click: Email Templates
3. Click: Confirm signup
4. Click: "Send test email" button
```

---

### For RLS Policies Setup:

**All Steps:**
```
1. Click: </> SQL Editor (in left sidebar)
2. Copy SQL code block from guide
3. Paste into the text box
4. Click: "Run" button
5. Repeat for all 4 code blocks
```

**Verification:**
```
1. Click: 🗄️ Database (in left sidebar)
2. Click: Policies
3. You should see ~13 policies
```

---

## 🎯 Visual Reference

### Main Dashboard Layout:

```
┌──────────────────────────────────────────────────────────────┐
│  Supabase Logo    [Project: Tayosa]           [User Menu]    │
├──────────────────────────────────────────────────────────────┤
│                                                                │
│  ┌──────────────┐  ┌────────────────────────────────────┐   │
│  │              │  │                                      │   │
│  │  Left        │  │  Main Content Area                  │   │
│  │  Sidebar     │  │  (This is where settings appear)    │   │
│  │              │  │                                      │   │
│  │  🏠 Home     │  │  When you click something in the    │   │
│  │  📊 Tables   │  │  left sidebar, the content shows    │   │
│  │  🔑 Auth     │  │  up here.                           │   │
│  │  🗄️ Database │  │                                      │   │
│  │  📁 Storage  │  │                                      │   │
│  │  </> SQL     │  │                                      │   │
│  │  ⚙️ Settings │  │                                      │   │
│  │              │  │                                      │   │
│  └──────────────┘  └────────────────────────────────────┘   │
│                                                                │
└──────────────────────────────────────────────────────────────┘
```

---

## 🔍 Can't Find Something?

### If you can't find a menu item:

1. **Make sure you clicked the right main section first**
   - Example: To find "URL Configuration", first click "Authentication"

2. **Look for sub-sections**
   - Some items are grouped under "CONFIGURATION" or "SETTINGS" headers

3. **Try scrolling down**
   - Some menus have more items below the fold

4. **Use the search bar**
   - There's usually a search bar at the top of the dashboard

---

## 📞 Common Questions

**Q: I don't see "JWT Keys" under Authentication**  
A: JWT Keys is under "Project Settings" (gear icon ⚙️), not under "Authentication"!

**Q: I don't see "URL Configuration" under Authentication**  
A: Make sure you clicked "Authentication" in the left sidebar first. Then look for "URL Configuration" in the menu that appears.

**Q: Where is the SQL Editor?**  
A: Look for the `</>` icon in the left sidebar. It might be labeled "SQL Editor" or just show the icon.

**Q: I can't find "Policies" to verify my RLS setup**  
A: Click "Database" in the left sidebar, then look for "Policies" in the menu that appears.

---

## ✅ Quick Checklist

Use this to track your progress:

**Quick Supabase Setup:**
- [ ] Found and verified JWT Keys (⚙️ Project Settings → JWT Keys)
- [ ] Found URL Configuration (🔑 Authentication → URL Configuration)
- [ ] Found Sign-in / Providers (🔑 Authentication → Sign-in / Providers)
- [ ] Found Email Templates (🔑 Authentication → Email Templates)
- [ ] Found Storage section (📁 Storage)
- [ ] Sent test email successfully

**RLS Policies Setup:**
- [ ] Found SQL Editor (</> SQL Editor)
- [ ] Ran all 4 SQL code blocks
- [ ] Verified policies (🗄️ Database → Policies)
- [ ] See ~13 policies listed

---

## 🎉 You're Ready!

Now that you know where everything is, follow these guides:

1. **First:** `docs/QUICK_SUPABASE_SETUP.md`
2. **Second:** `docs/SETUP_RLS_POLICIES.md`

**You've got this!** 💪

---

**Last Updated:** April 22, 2026  
**Dashboard URL:** https://supabase.com/dashboard/project/[YOUR-PROJECT-REF]
