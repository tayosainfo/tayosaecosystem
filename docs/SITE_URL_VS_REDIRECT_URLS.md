# 🎯 Understanding Site URL vs Redirect URLs

## Quick Answer

**Site URL:** Your main/default website (use Vercel)  
**Redirect URLs:** ALL allowed redirect destinations (include both Vercel AND localhost)

---

## 📊 Visual Explanation

```
┌─────────────────────────────────────────────────────────────┐
│                    SUPABASE CONFIGURATION                     │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Site URL (Default):                                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ https://tayosaecosystem.vercel.app                  │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  This is your "fallback" URL. Email links use this.         │
│                                                               │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Redirect URLs (All Allowed):                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ ✅ https://tayosaecosystem.vercel.app/auth/callback │   │
│  │ ✅ https://tayosaecosystem.vercel.app/auth/verify   │   │
│  │ ✅ https://tayosaecosystem.vercel.app/auth/reset    │   │
│  │ ✅ http://localhost:3000/auth/callback              │   │
│  │ ✅ http://localhost:3000/auth/verify-email          │   │
│  │ ✅ http://localhost:3000/auth/reset-password        │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  These are ALL the URLs Supabase will redirect to.          │
│  You can use ANY of these in your app!                      │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔄 How It Works

### Scenario 1: Developing Locally

```
You → http://localhost:3000 → Click "Register"
                                      ↓
                              Supabase sends email
                                      ↓
                              User clicks email link
                                      ↓
                    Supabase checks: Is localhost:3000 allowed?
                                      ↓
                              YES! (in Redirect URLs)
                                      ↓
                    Redirect to: http://localhost:3000/auth/verify-email ✅
```

### Scenario 2: Testing on Vercel

```
You → https://tayosaecosystem.vercel.app → Click "Register"
                                                    ↓
                                            Supabase sends email
                                                    ↓
                                            User clicks email link
                                                    ↓
                              Supabase checks: Is Vercel URL allowed?
                                                    ↓
                                            YES! (in Redirect URLs)
                                                    ↓
              Redirect to: https://tayosaecosystem.vercel.app/auth/verify-email ✅
```

### Scenario 3: Email Link (No Context)

```
User receives email → Clicks link → Supabase needs to redirect
                                              ↓
                                    Where should I send them?
                                              ↓
                                    Use "Site URL" as default
                                              ↓
                    Redirect to: https://tayosaecosystem.vercel.app ✅
```

---

## ❌ What Happens If You Use localhost as Site URL?

### Problem: Email Links Won't Work for Others

```
Site URL: http://localhost:3000  ← BAD IDEA!

Team Member → Gets email → Clicks link
                                ↓
                    Supabase redirects to: http://localhost:3000
                                ↓
                    ❌ ERROR! localhost doesn't exist on their computer
                                ↓
                    They can't verify their email!
```

### Solution: Use Vercel as Site URL

```
Site URL: https://tayosaecosystem.vercel.app  ← GOOD!

Team Member → Gets email → Clicks link
                                ↓
                    Supabase redirects to: https://tayosaecosystem.vercel.app
                                ↓
                    ✅ Works! They can verify their email!
```

---

## ✅ Recommended Configuration

### Site URL:
```
https://tayosaecosystem.vercel.app
```

**Why?**
- ✅ Email links work for everyone
- ✅ Team members can test
- ✅ It's your "production-like" environment
- ✅ It's publicly accessible

### Redirect URLs:
```
https://tayosaecosystem.vercel.app/auth/callback
https://tayosaecosystem.vercel.app/auth/verify-email
https://tayosaecosystem.vercel.app/auth/reset-password
http://localhost:3000/auth/callback
http://localhost:3000/auth/verify-email
http://localhost:3000/auth/reset-password
```

**Why include localhost?**
- ✅ You can still develop locally
- ✅ Local authentication works
- ✅ You can test without deploying
- ✅ Best of both worlds!

---

## 🤔 Common Questions

### Q: Can I develop locally with this setup?
**A:** YES! ✅

Even though Site URL is Vercel, you can still develop on localhost because localhost is in the Redirect URLs list.

### Q: Will email links work locally?
**A:** Partially.

- ✅ If you're testing locally and click the email link, it will go to Vercel (Site URL)
- ✅ But you can manually change the URL to localhost in your browser
- ✅ Or you can copy the token from the URL and use it locally

**Better approach:** Test email flows on Vercel, test other features locally.

### Q: Can I change Site URL later?
**A:** YES! ✅

You can change it anytime. When you deploy to production (`app.tayosa.com`), just update the Site URL to your production domain.

### Q: Do I need to remove localhost URLs later?
**A:** NO! ❌

You can keep localhost URLs in Redirect URLs forever. They don't hurt anything and let you keep developing locally.

### Q: What if I want email links to go to localhost?
**A:** You can, but it's not recommended.

If you really want this:
1. Set Site URL to `http://localhost:3000`
2. But remember: email links won't work for anyone else
3. Better: Keep Site URL as Vercel and test emails there

---

## 📋 Configuration Checklist

- [ ] Site URL set to: `https://tayosaecosystem.vercel.app`
- [ ] Added 3 Vercel redirect URLs
- [ ] Added 3 localhost redirect URLs
- [ ] Clicked "Save"
- [ ] Tested registration on Vercel
- [ ] Tested registration on localhost
- [ ] Email links work correctly

---

## 🎯 Summary

**Site URL = Default/Fallback**
- Use your live deployment (Vercel)
- Email links use this
- Should be publicly accessible

**Redirect URLs = All Allowed Destinations**
- Include ALL environments (Vercel + localhost)
- Your app can redirect to any of these
- Add more as you add environments

**Best Practice:**
- Site URL: Vercel (or production domain)
- Redirect URLs: Vercel + localhost + production + mobile
- This gives you maximum flexibility!

---

**Still confused?** Think of it this way:
- **Site URL** = Your home address (one main location)
- **Redirect URLs** = All places you're allowed to go (multiple locations)

You have one home, but you can go to many places! 🏠➡️🏢🏪🏫

---

**Last Updated:** April 22, 2026  
**Recommended Site URL:** `https://tayosaecosystem.vercel.app`  
**Recommended Redirect URLs:** Vercel + localhost (6 total)
