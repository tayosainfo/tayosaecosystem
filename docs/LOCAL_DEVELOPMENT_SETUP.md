# 🛠️ Local Development Setup - Testing the Migration

**Platform:** Windows (your current system)  
**Purpose:** Test the InsForge → Supabase migration locally  
**Time Required:** 30-60 minutes setup

---

## 📋 Required Tools Overview

To test this project locally, you have **two options**:

### **Option A: Using Supabase Database (Recommended)**
1. **Node.js & npm** - For frontend
2. **Go** - For backend services
3. **PostgreSQL Client (psql)** - For database operations
4. **Git** - For version control
5. **Code Editor** - VS Code recommended

### **Option B: Using Local Database**
1. **Node.js & npm** - For frontend
2. **Go** - For backend services
3. **Docker & Docker Compose** - For local PostgreSQL database
4. **Git** - For version control
5. **Code Editor** - VS Code recommended

**Recommendation:** Use **Option A** (Supabase) since you've already configured it and the migration is designed for Supabase.

---

## 🔧 Installation Guide

### 1. Node.js & npm (Frontend)

**What it's for:** Running the web frontend locally

**Install:**
1. Go to: https://nodejs.org/
2. Download **LTS version** (recommended)
3. Run installer with default settings
4. Verify installation:

```bash
node --version    # Should show v18.x.x or v20.x.x
npm --version     # Should show 9.x.x or 10.x.x
```

**Alternative (if you have Chocolatey):**
```bash
choco install nodejs
```

---

### 2. Go (Backend Services)

**What it's for:** Running the 9 backend microservices

**Install:**
1. Go to: https://golang.org/dl/
2. Download **Windows installer** (.msi file)
3. Run installer with default settings
4. Verify installation:

```bash
go version    # Should show go1.21.x or go1.22.x
```

**Alternative (if you have Chocolatey):**
```bash
choco install golang
```

---

### 3. Flutter (Mobile App)

**What it's for:** Running the mobile app (optional for basic testing)

**Install:**
1. Go to: https://docs.flutter.dev/get-started/install/windows
2. Download **Flutter SDK**
3. Extract to `C:\flutter` (or your preferred location)
4. Add to PATH: `C:\flutter\bin`
5. Verify installation:

```bash
flutter --version    # Should show Flutter 3.x.x
flutter doctor       # Shows what's missing
```

**Note:** You can skip Flutter if you only want to test web/backend

---

### 4. Docker & Docker Compose (Optional - For Local Database)

**What it's for:** Running a local PostgreSQL database instead of using Supabase

**When you need it:**
- If you want to test with a completely local setup
- If you want to avoid using Supabase database during development
- If you're working offline

**Install:**
1. Go to: https://www.docker.com/products/docker-desktop/
2. Download **Docker Desktop for Windows**
3. Run installer and follow setup wizard
4. Restart your computer when prompted
5. Verify installation:

```bash
docker --version         # Should show Docker version 20.x.x+
docker-compose --version # Should show docker-compose version 2.x.x+
```

**Alternative (if you have Chocolatey):**
```bash
choco install docker-desktop
```

**Note:** Docker Desktop includes Docker Compose automatically.

---

### 5. PostgreSQL Client (psql)

**What it's for:** 
- **Option A (Supabase):** Connecting to Supabase database and running migrations
- **Option B (Local):** Connecting to local Docker database

**Install:**

**Option A: Install PostgreSQL (Full)**
1. Go to: https://www.postgresql.org/download/windows/
2. Download installer
3. Install with default settings
4. This includes `psql` command

**Option B: Install just psql (Lightweight)**
1. Download from: https://www.enterprisedb.com/download-postgresql-binaries
2. Extract and add `bin` folder to PATH

**Verify:**
```bash
psql --version    # Should show psql (PostgreSQL) 15.x or 16.x
```

**Alternative (if you have Chocolatey):**
```bash
choco install postgresql
```

---

### 5. Git (Version Control)

**What it's for:** Cloning the repository and managing code

**Install:**
1. Go to: https://git-scm.com/download/win
2. Download and run installer
3. Use default settings (Git Bash + Git GUI)
4. Verify:

```bash
git --version    # Should show git version 2.x.x
```

---

### 6. VS Code (Recommended Editor)

**What it's for:** Editing code and managing the project

**Install:**
1. Go to: https://code.visualstudio.com/
2. Download Windows installer
3. Install with default settings

**Recommended Extensions:**
- Go (for backend development)
- Flutter (for mobile development)
- TypeScript and JavaScript (for frontend)
- PostgreSQL (for database work)

---

## 🚀 Project Setup Steps

You have **two options** for setting up the database:

### **Option A: Using Supabase Database (Recommended)**
This uses your configured Supabase project - **recommended** since the migration is designed for Supabase.

### **Option B: Using Local Docker Database**
This runs a local PostgreSQL database using Docker - useful for completely offline development.

---

## 🎯 **Option A: Supabase Database Setup**

### Step 1: Clone and Setup

```bash
# Clone your repository (replace with your actual repo URL)
git clone https://github.com/your-username/tayosaecosystem.git
cd tayosaecosystem

# Install frontend dependencies
npm install

# Install Go dependencies for backend services
cd services/user-service
go mod tidy
cd ../..

# If testing mobile app
cd app/mobile_app
flutter pub get
cd ../..
```

### Step 2: Environment Configuration (Supabase)

**Create `.env` file in project root:**

```bash
# Frontend Environment Variables
VITE_SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
VITE_SUPABASE_ANON_KEY=[YOUR_ANON_KEY]
VITE_API_BASE_URL=http://localhost:8080

# Backend Environment Variables
SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
SUPABASE_ANON_KEY=[YOUR_ANON_KEY]
SUPABASE_SERVICE_ROLE_KEY=[YOUR_SERVICE_ROLE_KEY]
DATABASE_URL=postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres

# Optional: Allow memory fallback for testing without DB
USER_SERVICE_ALLOW_MEMORY_FALLBACK=1
```

**Replace `[YOUR-PASSWORD]` with your actual Supabase database password.**

### Step 3: Database Migration (Supabase)

**Run the migration to set up database schema:**

```bash
# Connect to Supabase and run migration
psql "postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres" -f db/migrations/011_rename_insforge_to_supabase.sql
```

**Or use Supabase SQL Editor:**
1. Open: https://supabase.com/dashboard/project/[YOUR-PROJECT-REF]
2. Go to SQL Editor
3. Copy contents of `db/migrations/011_rename_insforge_to_supabase.sql`
4. Paste and run

### Step 4: Set Up RLS Policies (Supabase)

**CRITICAL:** Follow the RLS setup guide:

```bash
# Open the guide
code docs/SETUP_RLS_POLICIES.md

# Or follow the steps in Supabase SQL Editor
# This is REQUIRED because automatic RLS is enabled
```

---

## 🐳 **Option B: Local Docker Database Setup**

### Step 1: Clone and Setup (Same as Option A)

```bash
# Clone your repository
git clone https://github.com/your-username/tayosaecosystem.git
cd tayosaecosystem

# Install dependencies
npm install
cd services/user-service && go mod tidy && cd ../..
```

### Step 2: Start Local Database with Docker

```bash
# Start PostgreSQL database using Docker Compose
docker-compose up -d postgres

# Verify it's running
docker ps
# Should show tayosa-postgres container running on port 5432
```

### Step 3: Environment Configuration (Local)

**Create `.env` file in project root:**

```bash
# Frontend Environment Variables
VITE_SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
VITE_SUPABASE_ANON_KEY=[YOUR_ANON_KEY]
VITE_API_BASE_URL=http://localhost:8080

# Backend Environment Variables (Local Database)
SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
SUPABASE_ANON_KEY=[YOUR_ANON_KEY]
SUPABASE_SERVICE_ROLE_KEY=[YOUR_SERVICE_ROLE_KEY]
DATABASE_URL=postgres://tayosa:password@localhost:5432/tayosa_banking?sslmode=disable

# Note: Using local Docker database but still need Supabase for auth
```

### Step 4: Database Migration (Local)

```bash
# Apply all migrations to local database
psql "postgres://tayosa:password@localhost:5432/tayosa_banking?sslmode=disable" -f db/migrations/001_initial_schema.sql
psql "postgres://tayosa:password@localhost:5432/tayosa_banking?sslmode=disable" -f db/migrations/002_add_users_table.sql
# ... apply all migration files in order
psql "postgres://tayosa:password@localhost:5432/tayosa_banking?sslmode=disable" -f db/migrations/011_rename_insforge_to_supabase.sql
```

**Or apply all at once:**
```bash
# Apply all migrations (if you have a script)
for file in db/migrations/*.sql; do
  psql "postgres://tayosa:password@localhost:5432/tayosa_banking?sslmode=disable" -f "$file"
done
```

### Step 5: No RLS Setup Needed (Local)

**Note:** Local PostgreSQL doesn't have RLS enabled by default, so you can skip the RLS policies setup for local testing.

---

## 🧪 Testing the Migration

### Test 1: Frontend Only

**Start the frontend:**
```bash
npm run dev
```

**Expected result:**
- Opens http://localhost:3000
- You can see the login/register pages
- Supabase client is initialized

**Test actions:**
- Try to register a new user
- Check if email verification works
- Try to log in

---

### Test 2: Backend Services

**Start user-service:**
```bash
cd services/user-service
go run .
```

**Expected result:**
- Service starts on port 8080
- Logs show "Supabase auth enabled"
- No compilation errors

**Test actions:**
```bash
# Test health endpoint
curl http://localhost:8080/health

# Test registration (replace with actual data)
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","fullName":"Test User"}'
```

---

### Test 3: Full Stack Integration

**Terminal 1 - Backend:**
```bash
cd services/user-service
go run .
```

**Terminal 2 - Frontend:**
```bash
npm run dev
```

**Test flow:**
1. Open http://localhost:3000
2. Register a new user
3. Check email for verification
4. Verify email and log in
5. Check that user data is stored correctly

---

### Test 4: Mobile App (Optional)

**If you installed Flutter:**

```bash
cd app/mobile_app

# Run on web (easiest for testing)
flutter run -d web-server --web-port 3001

# Or run on Android emulator (if you have Android Studio)
flutter run
```

---

## 🔍 Verification Checklist

### ✅ Installation Verification

Run these commands to verify everything is installed:

**Essential tools:**
```bash
node --version        # Node.js
npm --version         # npm
go version           # Go
psql --version       # PostgreSQL client
git --version        # Git
```

**If using Docker (Option B):**
```bash
docker --version         # Docker
docker-compose --version # Docker Compose
```

**Optional tools:**
```bash
flutter --version    # Flutter (optional)
code --version       # VS Code (optional)
```

### ✅ Project Setup Verification

**For both options:**
```bash
# Check if dependencies are installed
ls node_modules      # Should show many folders
cd services/user-service && go mod verify && cd ../..  # Should show no errors

# Check if environment variables are set
echo $VITE_SUPABASE_URL  # Should show Supabase URL (on Git Bash)
# Or on Windows CMD: echo %VITE_SUPABASE_URL%
```

**For Option B (Docker) only:**
```bash
# Check if Docker database is running
docker ps            # Should show tayosa-postgres container
docker-compose ps     # Should show postgres service as Up
```

### ✅ Migration Verification

**Option A (Supabase):**
```bash
# Check if migration ran successfully
psql "postgresql://postgres:[PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres" -c "SELECT column_name FROM information_schema.columns WHERE table_name = 'users' AND column_name IN ('supabase_user_id', 'supabase_login_email');"

# Should show both columns
```

**Option B (Local Docker):**
```bash
# Check if migration ran successfully
psql "postgres://tayosa:password@localhost:5432/tayosa_banking?sslmode=disable" -c "SELECT column_name FROM information_schema.columns WHERE table_name = 'users' AND column_name IN ('supabase_user_id', 'supabase_login_email');"

# Should show both columns
```

---

## 🚨 Common Issues & Solutions

### Issue: "docker: command not found"
**Solution:** 
- Install Docker Desktop for Windows
- Restart your terminal after installation
- Make sure Docker Desktop is running (check system tray)

### Issue: "Cannot connect to Docker daemon"
**Solution:**
- Start Docker Desktop application
- Wait for Docker to fully start (whale icon in system tray should be steady)
- Try the command again

### Issue: "Port 5432 already in use"
**Solution:**
- Stop any existing PostgreSQL services: `net stop postgresql-x64-15`
- Or use a different port in docker-compose.yml
- Or stop the Docker container: `docker-compose down`

### Issue: "go: command not found"
**Solution:** 
- Restart your terminal after installing Go
- Check PATH includes Go bin directory
- On Windows: `echo %PATH%` should include Go

### Issue: "psql: command not found"
**Solution:**
- Install PostgreSQL or just psql client
- Add PostgreSQL bin directory to PATH
- Restart terminal

### Issue: "Cannot connect to database"
**Solution:**
- Check your DATABASE_URL has correct password
- Verify Supabase project is accessible
- Use `USER_SERVICE_ALLOW_MEMORY_FALLBACK=1` for testing without DB

### Issue: "RLS policy error"
**Solution:**
- You MUST set up RLS policies first
- Follow `docs/SETUP_RLS_POLICIES.md`
- This is required because automatic RLS is enabled

### Issue: "CORS error in browser"
**Solution:**
- Make sure backend is running on port 8080
- Check VITE_API_BASE_URL points to http://localhost:8080
- Restart both frontend and backend

---

## 📱 Mobile Testing (Advanced)

### If you want to test the mobile app:

**Additional requirements:**
- **Android Studio** (for Android emulator)
- **Xcode** (for iOS simulator - Mac only)
- **Chrome** (for web testing)

**Setup:**
```bash
# Install Android Studio
# Download from: https://developer.android.com/studio

# Set up Flutter for mobile development
flutter doctor
flutter config --enable-web

# Create and start Android emulator (in Android Studio)
# Or use web for easier testing
flutter run -d web-server --web-port 3001
```

---

## 🎯 Quick Test Script

**Create this test script to verify everything works:**

**File: `test-migration.bat` (Windows)**
```batch
@echo off
echo Testing InsForge to Supabase Migration...

echo.
echo 1. Checking tools...
node --version || echo ERROR: Node.js not installed
go version || echo ERROR: Go not installed
psql --version || echo ERROR: PostgreSQL client not installed

echo.
echo 2. Installing dependencies...
call npm install

echo.
echo 3. Starting backend...
start cmd /k "cd services/user-service && go run ."

echo.
echo 4. Starting frontend...
timeout /t 5
start cmd /k "npm run dev"

echo.
echo 5. Opening browser...
timeout /t 10
start http://localhost:3000

echo.
echo Migration test setup complete!
echo Check both terminal windows for any errors.
echo Test registration and login at http://localhost:3000
pause
```

**Run it:**
```bash
test-migration.bat
```

---

## 📊 Expected Results

### ✅ Successful Migration Test

**You should see:**
1. **Frontend:** Loads at http://localhost:3000
2. **Backend:** Runs without errors, shows "Supabase auth enabled"
3. **Registration:** Works and sends verification email
4. **Login:** Works after email verification
5. **Database:** User data is stored in Supabase
6. **RLS:** Users can only see their own data

### ❌ Common Test Failures

**If registration fails:**
- Check RLS policies are set up
- Verify environment variables
- Check Supabase dashboard logs

**If backend won't start:**
- Check Go installation
- Verify DATABASE_URL or use memory fallback
- Check for compilation errors

**If frontend won't load:**
- Check Node.js installation
- Verify npm dependencies installed
- Check for port conflicts

---

## 🎉 Success Criteria

**Your migration test is successful if:**

1. ✅ All tools install without errors
2. ✅ Project dependencies install correctly
3. ✅ Backend services start and connect to Supabase
4. ✅ Frontend loads and shows Supabase integration
5. ✅ User registration works end-to-end
6. ✅ Email verification works
7. ✅ Login works after verification
8. ✅ RLS policies protect user data

**If all these work, your migration is 100% successful!** 🎉

---

## 📞 Need Help?

**If you get stuck:**
1. Check the error messages carefully
2. Verify all tools are installed correctly
3. Check environment variables are set
4. Review the RLS policies setup
5. Check Supabase dashboard for errors
6. Email: support@tayosa.com

---

**Estimated Setup Time:** 30-60 minutes  
**Testing Time:** 15-30 minutes  
**Total Time:** 45-90 minutes  

**Ready to test your migration?** Follow the steps above! 🚀
