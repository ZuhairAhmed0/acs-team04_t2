# Fuel Management System — وين الوقود؟

> Full-stack Fuel Management System encompassing project requirements, frontend & backend architecture, database design, API endpoints, and JSON contracts.
> - **Backend Stack:** Go 1.27 · `net/http` (`ServeMux`) · GORM · PostgreSQL · JWT (HS256, HttpOnly cookies) · bcrypt
> - **Frontend Stack:** React + Vite · MUI (RTL Arabic) · AOS · Axios (`withCredentials`) · Cairo Font

---

## Full Features Overview

### User Features (`role = 'user'`)
- **Registration & Welcome Bonus:** Digit-only Sudanese phone validation (max 10 digits) and `+10` welcome loyalty points at signup.
- **Wallet & Balance Breakdown:** Tracks `total_balance`, `allocated_balance`, and `available_balance` (`total = allocated + available`) in integer SDG units.
- **Loyalty Rewards Program:** Earn `+10` welcome points on registration and `+5` points on every wallet top-up or wallet-only fuel purchase (`10 points = 1 free gas gallon`, `1 point = 100 SDG`), with progress bar toward the next free gallon. Redeem loyalty points at checkout in multiples of `10` points (`10, 20, 30...`, worth `1,000 SDG` per 10 points) either as a discount alongside `available`/`allocated` balance or for a `0 SDG` points-only purchase (`payment_source: "loyalty"` or `amount: 0` with `loyalty_points_used`). Note: station allocations, station-to-station transfers, and purchases that redeem loyalty points earn `0` points.
- **Stations & Branch-Tap Fuel Purchase:** Browse active stations across Sudan, filter branches by city (`الخرطوم`, `أم درمان`, `بحري`, `بورتسودان`, `ود مدني`, `الأبيض`, `كسلا`), inspect real-time fuel availability (`full`, `low`, `empty`) and branch congestion percentage, and tap any non-empty branch to purchase fuel directly from `available` balance (credits station), `allocated` balance (deducts allocated balance without double-crediting), or `loyalty` points (`0 SDG` cash charge). Includes automatic duplicate-purchase alerts and a confirmation dialog if the same branch and amount were purchased within the last 10 minutes.
- **Transfers (`التحويلات`):** Dedicated workflows for wallet top-up, station balance allocation (shared across all branches of a station), and station-to-station balance transfers, with instant receipt dialogs.
- **Transaction History, Detail View & CSV Export:** Filterable by type, status, and date range, ordered newest first, with clickable transaction rows opening a full detail dialog (fetched from `GET /api/v1/transfers/{id}` for users) and authenticated CSV download (`transactions.csv`).

### Station Features (`role = 'station'`)
- **Registration & Approval Lifecycle:** Submit station and bank account details (creates a main branch and shared wallet in `pending` state awaiting admin approval).
- **Shared Station Dashboard:** View shared station balance, branch count, per-branch congestion bar chart, and recent transactions.
- **Branch Management:** Add branches by Sudanese city & address, toggle fuel availability (`full`, `low`, `empty`), and adjust branch congestion (`0`–`100%`, step `5`).
- **Cash Withdrawal Requests:** Submit idempotent withdrawal requests; open `pending` and `approved` withdrawals reserve against `available_balance` so over-requests are prevented (`402`). Track status and reference numbers.
- **Transaction History & CSV Export:** View and export station-scoped transactions as CSV.

### Admin Features (`role = 'admin'`)
- **System Dashboard:** Live metrics for total users, total/approved/pending stations, and pending cash-withdrawal requests.
- **User Management:** Activate/suspend user accounts and promote/revoke admin roles.
- **Station Moderation:** Review pending station applications (details, bank info, and branches with fuel/congestion status), approve, reject (with required reason), revoke approval, or delete stations.
- **Cash-Withdrawal Settlement:** Approve or reject pending withdrawals (with required reason), then settle approved withdrawals with an optional real bank `external_reference` that deducts the station wallet atomically and logs a `withdrawal` transaction.

### Security, Integrity & Reliability
- **HttpOnly Cookie Auth + Automatic Refresh:** `POST /auth/login` and `POST /auth/refresh` issue JWT access and SHA-256-hashed refresh tokens inside `HttpOnly` cookies (instead of exposing tokens in `localStorage` or JSON bodies), with `Authorization: Bearer` fallback on the API and automatic single-flight token refresh on `401` in Axios.
- **End-to-End Idempotency & Duplicate Purchase Protection:** All monetary mutations (`POST /transfers/*` and `POST /station/cash-withdrawals`) require an `Idempotency-Key` header (or `idempotency_key` body field, max 120 chars) paired with a SHA-256 `request_hash` and composite DB unique indexes (`idx_tx_actor_idempotency`, `idx_cw_station_idempotency`). Safe retries return the persisted record (`idempotent_replay: true`) and trigger a warning notice in the UI; parameter mismatches return `409 Conflict`. Fuel purchases also flag recent identical purchases (`duplicate_purchase: true` + `previous_reference` within a 10-minute window) and prompt the user for confirmation before submitting again.
- **Database Integrity & Concurrency Safety:** PostgreSQL `CHECK` constraint (`chk_wallets_non_negative`) forbids negative balances; GORM transactions lock rows with `SELECT ... FOR UPDATE` in deterministic deadlock-free order (`station` before `user`; lower `station_id` first on station-to-station transfers).
- **Hardening & Observability:** Per-IP token-bucket rate limiting (`RATE_LIMIT_RPS`, `RATE_LIMIT_BURST`), security headers (`X-Content-Type-Options`, `X-Frame-Options`, `X-XSS-Protection`, `Referrer-Policy`), strict CORS with credentials + `Idempotency-Key`, strict `.env` validation with zero hardcoded fallbacks, and daily structured request logging.

---

## Quick Start: Setup & Testing on Windows (after unzipping `Project.zip`)

Follow these steps in order in **Windows PowerShell** (or **Command Prompt / Windows Terminal**) after extracting the archive to run and test the full system (PostgreSQL, Go backend API, and React web app).

### Prerequisites (Windows)
- **Go** (`1.24+` — `go.mod` targets Go `1.27`, download from https://go.dev/dl/)
- **Node.js (v18+ LTS) & npm** (download from https://nodejs.org/)
- **PostgreSQL 17** via **Docker Desktop** / **Podman Desktop**, **or** the **PostgreSQL Windows Installer** (https://www.postgresql.org/download/windows/)

---

### Step 1: Unzip and Install Dependencies

`Frontend/node_modules` is excluded from `Project.zip`.

1. Right-click `Project.zip` → **Extract All...** (or run in PowerShell):
   ```powershell
   Expand-Archive -Path Project.zip -DestinationPath .
   cd Project
   ```
2. Install Frontend packages and download Backend Go modules:
   ```powershell
   # Install Frontend dependencies (restores Frontend\node_modules)
   cd Frontend
   npm install
   cd ..

   # Download Go modules for Backend (first-time setup on a new PC)
   cd Backend
   go mod download
   cd ..
   ```

---

### Step 2: Verify Environment Files

Both `.env` (project root) and `Frontend/.env` are included in `Project.zip`. Verify or recreate them if needed:

1. **Backend (`.env` in the `Project\` root):**
   Check that `.env` exists in the project root with your database, JWT, logging, rate-limit, and CORS settings:
   ```env
   PORT=8080
   APP_ENV=development
   JWT_SECRET_KEY=fuel_management_super_secret_jwt_key_2026_secure
   JWT_EXPIRATION_HOURS=24
   REFRESH_TOKEN_EXPIRATION_DAYS=30
   POSTGRES_HOST=localhost
   POSTGRES_PORT=5588
   POSTGRES_USERNAME=Server
   POSTGRES_PASSWORD=server
   POSTGRES_DATABASE=fuel_db
   POSTGRES_TIMEZONE=Africa/Khartoum
   LOG_PATH=logs
   RATE_LIMIT_RPS=10
   RATE_LIMIT_BURST=20
   CORS_ALLOWED_ORIGINS=*
   ```
   > Running the backend or seeder from `Backend\` automatically reads `Backend\.env` or falls back to `..\.env` in the project root, and resolves `LOG_PATH=logs` to `Project\logs\`.
   > If you use a local Windows PostgreSQL installation running on the default port `5432` instead of Docker/Podman, change `POSTGRES_PORT=5432` (and make sure the user `Server` / password `server` / DB `fuel_db` exist).

2. **Frontend (`Frontend\.env`):**
   If `Frontend\.env` is missing, copy it from `Frontend\.env.example`:
   ```powershell
   Copy-Item Frontend\.env.example Frontend\.env
   ```
   *(In Command Prompt `cmd`, use `copy Frontend\.env.example Frontend\.env`.)*

   Default content:
   ```env
   VITE_API_BASE_URL=http://localhost:8080
   ```
   > Use `http://localhost:8080` and open the app at `http://localhost:5173` (do not mix `localhost` and `127.0.0.1`, so `HttpOnly` auth cookies work properly in your browser).

---

### Step 3: Start PostgreSQL

Choose **one** option below:

#### Option A: Using Docker Desktop (or Podman) on Windows (Recommended)
Run this single-line command in **PowerShell** or **Command Prompt** (replace `docker` with `podman` if using Podman):
```powershell
# First time: create and start PostgreSQL 17 on port 5588 matching .env
docker run -d --name karary-postgres -e POSTGRES_USER=Server -e POSTGRES_PASSWORD=server -e POSTGRES_DB=fuel_db -p 5588:5432 postgres:17

# On subsequent runs, if the container already exists:
docker start karary-postgres
```
Wait until PostgreSQL is ready before seeding:
```powershell
docker exec karary-postgres pg_isready -U Server -d fuel_db
```

#### Option B: Using Local PostgreSQL Installed on Windows (No Docker)
1. Open **pgAdmin 4** or **SQL Shell (psql)** as the `postgres` superuser.
2. Run:
   ```sql
   CREATE USER "Server" WITH PASSWORD 'server';
   CREATE DATABASE fuel_db OWNER "Server";
   ```
3. In `.env` at the project root, set `POSTGRES_PORT=5432` (or whatever port your Windows PostgreSQL service listens on — check Windows **Services** for `postgresql-x64-17`).

---

### Step 4: Initialize Tables & Seed the Database

Run the seeder from `Backend\` (so Go resolves `Backend\go.mod`) to auto-migrate all 7 PostgreSQL tables, enforce wallet constraints, and populate Sudanese Arabic test data (2 admins, 10 users with loyalty points, 8 Sudanese stations across all statuses, 24 branches across 7 cities in Arabic, reconciled wallets, 28 transactions, and 6 cash withdrawals):

```powershell
cd Backend
go run assets/seed.go
cd ..
```

> By default `seed.go` clears existing data first (`-clean=true`). Pass `-clean=false` to preserve existing rows.

---

### Step 5: Start the Backend Server (`http://localhost:8080`)

> **Windows Note:** `Backend\bin\server` is a Linux `x86_64` binary and cannot run directly on Windows. Use **Option A** below on Windows.

#### Option A: Run from Source on Windows (Recommended)
In **Terminal 1** (PowerShell / Windows Terminal):
```powershell
cd Backend
go run main.go
```
Or build and run a native Windows `.exe`:
```powershell
cd Backend
go build -o bin\server.exe main.go
.\bin\server.exe
```

#### Option B: Run the Container Image (`fuel-management-1.0.1.tar`)
- **On Windows (Docker Desktop):** `--network host` cannot reach `localhost:5588` directly inside the Linux VM. Temporarily set `POSTGRES_HOST=host.docker.internal` in `.env` and run in PowerShell:
  ```powershell
  docker load -i Backend\assets\fuel-management-1.0.1.tar
  docker run -d --name fuel-management --env-file .env -v "${PWD}\logs:/app/logs" -p 8080:8080 localhost/fuel-management:1.0.1
  ```
- **On Linux (Podman):**
  ```bash
  podman load -i Backend/assets/fuel-management-1.0.1.tar
  podman run -d --name fuel-management --env-file .env -v "$(pwd)/logs:/app/logs" --network host localhost/fuel-management:1.0.1
  ```

---

### Step 6: Start the Frontend (`http://localhost:5173`)

Open a **second terminal (Terminal 2)** in the `Project\` directory (keep Terminal 1 running the backend):

```powershell
cd Frontend
npm run dev
```

Open **http://localhost:5173** in your browser (Edge, Chrome, or Firefox).

> **Testing from a Mobile Phone on the Same Wi-Fi (Windows):**
> 1. Run `ipconfig` in PowerShell and note your PC's **IPv4 Address** (e.g. `192.168.1.25`).
> 2. Set `VITE_API_BASE_URL=http://<YOUR_PC_LAN_IP>:8080` in `Frontend\.env`.
> 3. Start Vite with `npm run dev -- --host` and open `http://<YOUR_PC_LAN_IP>:5173` on your phone.
> 4. If Windows Defender Firewall prompts for Go or Node.js, allow access on **Private networks**.

---

### Step 7: Verify & Test the Full System

#### 1. Interactive Browser Walkthrough
Open **http://localhost:5173** and sign in with any seeded account (default password for all accounts: `Password123!`):

| Role | Email | Password | What to Test |
|------|-------|----------|--------------|
| **User (مستخدم)** | `ahmed.khalil@example.com` | `Password123!` | Wallet balances, loyalty points (`10 pts = 1 gallon`), top-up (`+5 pts`), station allocation (`0 pts`), tapping a non-empty station branch on **المحطات** to purchase fuel (`+5 pts`), station-to-station transfers (`0 pts`), history & CSV export |
| **Station (محطة)** | `al-nile@petroleum.sd` | `Password123!` | Shared wallet, congestion chart, add branches, toggle fuel status (`full`/`low`/`empty`), update congestion %, submit cash-withdrawal requests, CSV export |
| **Admin (مدير)** | `admin@fuelms.com` | `Password123!` | System metrics, user management (activate/suspend, promote/revoke admin), approve/reject/revoke/delete stations, approve/reject/complete withdrawals with external reference |

#### 2. Automated Checks on Windows
```powershell
# Backend static analysis and build check (from Project\Backend)
cd Backend
go vet ./...
go build ./...
cd ..

# Frontend production bundle check (from Project\Frontend)
cd Frontend
npm run build
cd ..
```

#### 3. Quick API Smoke Test in Windows PowerShell (`curl.exe`)
> **PowerShell Tip:** Always use `curl.exe` (not `curl`, which is an alias for `Invoke-WebRequest` in PowerShell) and pass `-c` / `-b` because authentication tokens are stored in `HttpOnly` cookies.

```powershell
# Login as a user and save HttpOnly auth cookies to cookies.txt
curl.exe -s -c cookies.txt -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d "{\"email\":\"ahmed.khalil@example.com\",\"password\":\"Password123!\"}"

# Verify authenticated profile using saved HttpOnly cookie
curl.exe -s -b cookies.txt http://localhost:8080/api/v1/auth/me

# Test idempotent top-up with Idempotency-Key header
curl.exe -s -b cookies.txt -X POST http://localhost:8080/api/v1/transfers/topup -H "Content-Type: application/json" -H "Idempotency-Key: win-test-topup-001" -d "{\"amount\":50000}"

# Login as admin
curl.exe -s -c admin-cookies.txt -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d "{\"email\":\"admin@fuelms.com\",\"password\":\"Password123!\"}"
```

---

## 1. Project Structure

```
Project/
├── .env                                 # Environment variables
├── logs/                                # Daily log files (auto-created)
├── README.md                            # Unified project documentation
├── Backend/                             # Go REST API Server
│   ├── main.go                          # Entry point → server.Run()
│   ├── go.mod / go.sum                  # Go module (module: wa, Go 1.27)
│   ├── assets/                          # Seed script, Dockerfile & image archive
│   │   ├── dockerfile
│   │   ├── seed.go
│   │   └── fuel-management-1.0.1.tar
│   ├── bin/                             # Compiled server binary
│   └── app/
│       ├── config/
│       │   └── config.go                # Config struct, Load() from .env / ../.env
│       ├── database/
│       │   ├── db.go                    # InitDB(), AutoMigrate all models
│       │   └── models.go                # GORM models (7 tables)
│       ├── helpers/
│       │   └── http_helper.go           # Generic Encode[T], Decode[T], DecodeValid[T], RespondWithError
│       ├── security/
│       │   ├── jwt.go                   # JWT generation/validation, refresh tokens
│       │   ├── passwords.go             # bcrypt hash & verify
│       │   ├── logger.go                # File + stdout logging middleware
│       │   └── ratelimiting.go          # Per-IP token bucket rate limiter
│       ├── server/
│       │   ├── server.go                # NewServer(), Run(), graceful shutdown
│       │   ├── routers.go               # addRoutes() — all route registrations
│       │   └── middleware.go            # AuthMiddleware, RequireRole
│       └── packages/
│           ├── auth/                    # Registration, login, logout, refresh, me
│           │   ├── handler.go
│           │   ├── service.go
│           │   ├── repo.go
│           │   └── dto.go
│           ├── user/                    # Dashboard, wallet, stations list, transactions
│           │   ├── handler.go
│           │   ├── service.go
│           │   ├── repo.go
│           │   └── dto.go
│           ├── station/                 # Dashboard, branches, fuel status, withdrawals
│           │   ├── handler.go
│           │   ├── service.go
│           │   ├── repo.go
│           │   └── dto.go
│           ├── transfers/               # Top-up, allocate, fuel purchase, station-to-station
│           │   ├── handler.go
│           │   ├── service.go
│           │   ├── repo.go
│           │   └── dto.go
│           └── admin/                   # Users, stations, withdrawal management
│               ├── handler.go
│               ├── service.go
│               ├── repo.go
│               └── dto.go
└── Frontend/                            # React + Vite RTL Web Application
    ├── .env / .env.example              # VITE_API_BASE_URL configuration
    ├── index.html
    ├── package.json
    ├── vite.config.js
    └── src/
        ├── app/                         # Routes and providers (المسارات والـ providers)
        ├── components/                  # brand · common · layout
        ├── config/                      # Constants (roles, banks, cities, statuses) & navigation
        ├── context/                     # Auth · Toast
        ├── features/                    # auth · user · station · admin
        ├── hooks/                       # useRequest · useSubmit · useFormState · useConfirmedAction
        ├── services/                    # http.js + service per endpoint group
        ├── theme/                       # Colors, MUI theme, RTL cache
        └── utils/                       # Formatting, errors, CSV download, session
```

### Backend 3-Layer Architecture

Each backend package follows a **3-layer architecture**:

| Layer | File | Responsibility |
|-------|------|----------------|
| Handler | `handler.go` | Decode HTTP request → call service → encode response |
| Service | `service.go` | Business logic, validation, orchestration |
| Repository | `repo.go` | Database access via GORM queries and transactions |
| DTO | `dto.go` | Request and response structs with JSON tags |

### Frontend Integration Notes (ملاحظات الواجهة الأمامية)

- **Authentication (`POST /auth/login`):** Uses `HttpOnly` cookies (`access_token`, `refresh_token`) with Axios `withCredentials: true` (no tokens in `localStorage`), automatically refreshes expired access tokens on `401` via `POST /auth/refresh`, and also supports `Authorization: Bearer <token>` headers for CLI/API clients.
- **Admin Dashboard (`GET /admin/dashboard`):** Reads `users_count`, `stations_count` / `registered_stations_count`, `approved_stations_count`, `pending_stations_count`, and `pending_withdrawals_count`.
- **Moderation & Role Payloads:** Station or withdrawal rejection sends `{ "rejection_reason": "..." }`. `PATCH /admin/users/:id/status` sends `{ "status": "active|suspended" }`, and `/role` sends `{ "role": "user|admin" }`.
- **Branch Cities:** Stored and filtered in Arabic (`الخرطوم`, `أم درمان`, `بحري`, `بورتسودان`, `ود مدني`, `الأبيض`, `كسلا`), matching `branches.city` directly in the database with no hardcoded backend city mappings.
- **Idempotency Keys:** Monetary mutations (`POST /transfers/*` and `POST /station/cash-withdrawals`) attach a stable per-attempt UUID key in both the `Idempotency-Key` header and `idempotency_key` JSON body field so retries never duplicate transactions.
- **External Transaction Reference:** Optional bank reference (`external_reference`) submitted when completing an approved cash withdrawal is persisted on the withdrawal record and shown in both admin and station withdrawal lists.

---

## 2. Account Types and Registration

There are three account roles:

- **User**: owns a wallet and loyalty points, views stations, tops up the wallet, allocates balance to stations, purchases fuel at station branches, transfers balance between stations, and reviews transactions.
- **Station**: one account per station. The station account manages its branches, views the shared station balance, submits cash-withdrawal requests, and reviews transactions.
- **Admin**: manages users and stations, approves or rejects station registrations, and reviews and completes cash-withdrawal requests.

Registration begins by choosing **User** or **Station**. Admin is not a public registration option.

### User Registration

| Field | Requirement |
|-------|-------------|
| `name` | Required |
| `password` | Required |
| `email` | Required and unique |
| `phone` | Required; digits only, maximum 10 digits |

### Station Registration

| Field | Requirement |
|-------|-------------|
| `station_name` / `name` | Required |
| `location` | Required; used for the main branch |
| `bank_account_number` | Required |
| `bank_account_name` | Required |
| `bank_name` | Required; select one of: Bank of Khartoum, Omdurman National Bank, Faisal Islamic Bank |
| `email` | Required and unique |
| `password` | Required |

A newly registered station remains pending ("قيد المراجعة") until an admin approves it. Registration creates its main branch automatically using the submitted location.

---

## 3. User Pages and Functions

The user account has four pages.

### 3.1 Dashboard

Shows:
- Total balance
- Allocated balance
- Available balance
- Loyalty points balance, equivalent fuel gallons (`10 points = 1 gallon`, `1 point = 100 SDG`), and progress toward the next gallon (regular users earn `10` welcome points on registration and `+5` points per top-up or wallet-only fuel purchase; redeemable at fuel checkout in multiples of `10` points)
- Quick actions for wallet top-up, station allocation, and station-to-station transfer
- Latest three transactions
- A link/button to the full transaction history

`total_balance = allocated_balance + available_balance`

### 3.2 Stations

- Lists stations and their branches.
- Supports filtering by city.
- Shows only branches matching the selected city. If a station has branches in multiple cities, only the matching branch is shown for that city.
- Each branch shows its fuel status (`full`, `low`, `empty`) and live congestion percentage (`congestion_percentage`). Tapping any available branch opens a **Purchase Fuel** (`شراء وقود`) dialog for that branch (`+5` loyalty points for wallet-only purchases, `0` when redeeming loyalty points), where the user chooses the payment source (`available` wallet balance, `allocated` balance, or `loyalty` points only) and amount (`0` SDG allowed when paying strictly with loyalty points), and can redeem loyalty points in multiples of `10` (`1 point = 100 SDG`, `10 points = 1,000 SDG`); branches with `fuel_status = "empty"` are disabled and rejected on the server. If the user attempts an identical purchase (same branch and amount within 10 minutes), the form shows an inline warning banner, prompts a confirmation dialog (`نعم، كرر الشراء`), and surfaces a duplicate notice in both the toast notification and receipt dialog (`duplicate_purchase` / `idempotent_replay`).

### 3.3 Transfers

Contains three operations:

1. **Wallet top-up** (`+5` loyalty points for regular users)
2. **Allocate balance to a station**
3. **Transfer balance from one station to another**

For station allocation, the user selects a station from a select control; no manual station-name entry. The allocation applies to the station's shared balance and is available across all its branches. **A branch is not selected for allocation.**

For fuel purchase (`fuel_purchase`) initiated from tapping a branch in Stations, paying from `available` credits the station's wallet; paying from `allocated` deducts from the already-allocated balance without double-crediting the station.

For station-to-station transfer, select the source station and destination station. **Branches are not selected for the transfer.** The user's `allocated_balance` must be greater than or equal to the transfer amount (`allocated_balance >= amount`).

After an operation, show a receipt/notification dialog containing:
- Reference number
- Amount
- Source account name (when applicable)
- Destination account/station name (when applicable)
- Branch name when the transaction occurs at a specific branch (such as `fuel_purchase`)
- Redeemed loyalty points and loyalty discount (when loyalty points are applied)
- Earned loyalty points and updated total loyalty points (when applicable)

### 3.4 All Transactions

Shows the user's transaction history ordered from **newest to oldest** (`created_at DESC, id DESC`), with relevant transaction details and filters. Clicking any transaction record opens a detail dialog (`GET /api/v1/transfers/{id}` for regular users) showing reference, parties, branch, amount, loyalty discount/points used, status, and timestamp.

The user can download their transaction history as a `.csv` file for analysis.

---

## 4. Admin Pages and Functions

The admin account has four pages.

### 4.1 Admin Dashboard

Shows general system information, including:
- Number of users
- Number of registered stations
- Number of approved stations
- Number of pending station registrations
- Number of pending cash-withdrawal requests

### 4.2 User Management

- List users.
- Activate or suspend a user.
- Promote a user to admin.
- Revoke admin role and return the account to user role.

### 4.3 Station Management

Two sections:

- **New registration requests:** View submitted station details, approve a station, or reject a station.
- **Existing stations:** List stations, view station, bank details, and branch fuel/congestion statuses, revoke station approval, or delete a station.

### 4.4 Cash-Withdrawal Management

- View cash-withdrawal requests submitted by stations.
- Review each requested amount.
- Approve or reject pending requests.
- After carrying out the transfer in the external financial app, mark the **approved** request as completed and optionally record the real external transfer reference (`external_reference`, persisted on the withdrawal record).

---

## 5. Station Pages and Functions

The station account has three pages: dashboard, branch management, and all transactions.

### 5.1 Station Dashboard

Shows:
- Shared station balance
- Number of branches
- A chart showing congestion percentage (`0`–`100%`) for each branch
- Cash-withdrawal request form (amount only)
- Latest five transactions (newest to oldest)

Each branch's congestion percentage is stored in the database (`branches.congestion_percentage`) and updated directly by the station via Branch Management (`PATCH /api/v1/station/branches/{id}/congestion`).

### 5.2 Branch Management

- The main branch is created automatically during station registration.
- The station can add additional branches.
- Adding a branch requires only:
  - City
  - Address
- Station name and bank details are inherited from the parent station and are not entered again.
- Each branch has a fuel status:
  - `full` — ممتلئة
  - `low` — ممتلئة قليلاً
  - `empty` — فارغة
- Each branch also has a congestion percentage slider (`0`–`100%`, step `5`).

All branches share the same station balance. Transactions can identify the branch where the transaction took place, but allocations and station-to-station transfers do not require a branch selection.

### 5.3 All Transactions

The station has a dedicated page showing all transactions associated with that station, ordered from **newest to oldest** (`created_at DESC, id DESC`).

The station can download its transaction history as a `.csv` file for analysis.

---

## 6. Database Design

The following tables cover the specified functionality.

### `users`

| Column | Suggested Type | Notes |
|--------|----------------|-------|
| `id` | `BIGINT PK` | |
| `name` | `VARCHAR(150)` | |
| `email` | `VARCHAR(255) UNIQUE` | Unique across user and station accounts |
| `password_hash` | `VARCHAR(255)` | Store a password hash, never a plain password |
| `phone` | `VARCHAR(30)` | Digits only, max 10 digits on registration |
| `role` | `ENUM('user','admin')` | |
| `status` | `ENUM('active','suspended')` | |
| `loyalty_points` | `BIGINT DEFAULT 0` | Loyalty points for `role = 'user'` (`+10` on registration, `+5` on top-up or wallet-only fuel purchase; `10 points = 1 gallon`, `1 point = 100 SDG`; constrained by `chk_users_loyalty_non_negative`) |
| `created_at` | `TIMESTAMP` | |

### `stations`

| Column | Suggested Type | Notes |
|--------|----------------|-------|
| `id` | `BIGINT PK` | |
| `name` | `VARCHAR(180)` | Station name |
| `location` | `VARCHAR(255)` | Main-branch location at registration |
| `bank_account_number` | `VARCHAR(100)` | |
| `bank_account_name` | `VARCHAR(180)` | |
| `bank_name` | `ENUM` | Bank of Khartoum, Omdurman National Bank, Faisal Islamic Bank |
| `email` | `VARCHAR(255) UNIQUE` | Unique across user and station accounts |
| `password_hash` | `VARCHAR(255)` | |
| `status` | `ENUM('pending','active','rejected','revoked')` | |
| `reviewed_by` | `BIGINT NULL FK -> users.id` | Admin who reviewed registration |
| `reviewed_at` | `TIMESTAMP NULL` | |
| `rejection_reason` | `TEXT NULL` | |
| `created_at` | `TIMESTAMP` | |

### `branches`

| Column | Suggested Type | Notes |
|--------|----------------|-------|
| `id` | `BIGINT PK` | |
| `station_id` | `BIGINT FK -> stations.id` | Parent station |
| `city` | `VARCHAR(120)` | |
| `address` | `VARCHAR(255)` | |
| `is_main` | `BOOLEAN` | One main branch per station |
| `fuel_status` | `ENUM('full','low','empty')` | |
| `congestion_percentage` | `INTEGER` | `0`–`100` congestion percentage |
| `created_at` | `TIMESTAMP` | |

### `wallets`

| Column | Suggested Type | Notes |
|--------|----------------|-------|
| `id` | `BIGINT PK` | |
| `owner_type` | `ENUM('user','station')` | |
| `owner_id` | `BIGINT` | User or station ID |
| `total_balance` | `BIGINT` | |
| `allocated_balance` | `BIGINT` | |
| `available_balance` | `BIGINT` | |
| `created_at` | `TIMESTAMP` | |

Each user/station has one wallet. Each station has one shared wallet used by all its branches. Amounts are stored as integer SDG units, not floating-point values.

### `transactions`

| Column | Suggested Type | Notes |
|--------|----------------|-------|
| `id` | `BIGINT PK` | |
| `reference` | `VARCHAR(80) UNIQUE` | |
| `actor_key` | `VARCHAR(60) NULL` | Composite unique key `(actor_key, idempotency_key)` (`idx_tx_actor_idempotency`) |
| `idempotency_key` | `VARCHAR(120) NULL` | Client-supplied key preventing duplicate transfer execution |
| `request_hash` | `VARCHAR(64) NULL` | SHA-256 hex digest of normalized request parameters |
| `type` | `ENUM` | `topup`, `allocation`, `fuel_purchase`, `station_transfer`, `withdrawal` |
| `amount` | `BIGINT` | |
| `from_type` | `ENUM('user','station') NULL` | |
| `from_id` | `BIGINT NULL` | Source account ID |
| `to_type` | `ENUM('user','station') NULL` | |
| `to_id` | `BIGINT NULL` | Destination account ID |
| `branch_id` | `BIGINT NULL FK -> branches.id` | Optional; only when a transaction occurs at a particular branch |
| `loyalty_points_used` | `BIGINT NOT NULL DEFAULT 0` | Loyalty points redeemed on a `fuel_purchase` (`0` or a multiple of `10`) |
| `loyalty_discount` | `BIGINT NOT NULL DEFAULT 0` | Discount value in SDG (`loyalty_points_used * 100`) applied to a `fuel_purchase` |
| `status` | `VARCHAR(30)` | |
| `created_at` | `TIMESTAMP` | |

A branch is **not** required for allocation or station-to-station transfer. `branch_id` and `branch_name` record the specific branch for branch-level transactions such as `fuel_purchase`.

### `cash_withdrawals`

| Column | Suggested Type | Notes |
|--------|----------------|-------|
| `id` | `BIGINT PK` | |
| `station_id` | `BIGINT FK -> stations.id` | Requesting station (composite unique index `idx_cw_station_idempotency` with `idempotency_key`) |
| `idempotency_key` | `VARCHAR(120) NULL` | Client-supplied key preventing duplicate withdrawal requests |
| `request_hash` | `VARCHAR(64) NULL` | SHA-256 hex digest of normalized withdrawal parameters |
| `amount` | `BIGINT` | |
| `status` | `ENUM('pending','approved','rejected','completed')` | |
| `external_reference` | `VARCHAR(120) NULL` | Optional transfer reference recorded at completion |
| `reviewed_by` | `BIGINT NULL FK -> users.id` | Admin |
| `reviewed_at` | `TIMESTAMP NULL` | |
| `rejection_reason` | `TEXT NULL` | |
| `created_at` | `TIMESTAMP` | |

### `refresh_tokens`

| Column | Suggested Type | Notes |
|--------|----------------|-------|
| `id` | `TEXT PK` | Token identifier |
| `user_id` | `TEXT` | `user:<id>` or `station:<id>` |
| `token_hash` | `VARCHAR(255) UNIQUE` | SHA-256 hash of refresh token |
| `expires_at` | `TIMESTAMP` | Expiration timestamp |
| `revoked` | `BOOLEAN` | Revocation flag (default: false) |
| `created_at` | `TIMESTAMP` | Creation timestamp |

---

## 7. API Endpoints

Base path: `/api/v1`

### Authentication

| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/auth/register/user` | Register user |
| POST | `/auth/register/station` | Register station |
| POST | `/auth/login` | Login user, station, or admin |
| POST | `/auth/refresh` | Refresh JWT access token |
| GET | `/auth/me` | Current account details |
| POST | `/auth/logout` | Logout and revoke refresh token |

### User

| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/user/dashboard` | User dashboard and latest 3 transactions |
| GET | `/user/wallet` | Wallet balances |
| GET | `/user/stations` | Stations and branches |
| GET | `/user/stations?city=أم درمان` | Filter branches by city |
| GET | `/user/transactions` | User transaction history |
| GET | `/user/transactions/export` | Download user transactions as CSV |

### Transfers

| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/transfers/topup` | Wallet top-up |
| POST | `/transfers/allocate` | Allocate to a station (no branch selection) |
| POST | `/transfers/purchase` | Purchase fuel at a specific station branch |
| POST | `/transfers/station-to-station` | Transfer between stations (no branch selection) |
| GET | `/transfers` | Transfer history |
| GET | `/transfers/:id` | Transfer details |

### Station

| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/station/dashboard` | Station dashboard and latest 5 transactions |
| GET | `/station/wallet` | Shared station wallet |
| GET | `/station/transactions` | All station transactions |
| GET | `/station/transactions/export` | Download station transactions as CSV |
| GET | `/station/branches` | List station branches |
| POST | `/station/branches` | Add branch |
| PATCH | `/station/branches/:id/fuel-status` | Update branch fuel status |
| PATCH | `/station/branches/:id/congestion` | Update branch congestion percentage |
| POST | `/station/cash-withdrawals` | Submit cash-withdrawal request |
| GET | `/station/cash-withdrawals` | View station withdrawal requests |

### Admin

| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/admin/dashboard` | System summary |
| GET | `/admin/users` | List users |
| PATCH | `/admin/users/:id/status` | Activate/suspend user |
| PATCH | `/admin/users/:id/role` | Promote/revoke admin role |
| GET | `/admin/stations/pending` | Pending station registrations |
| GET | `/admin/stations` | All stations |
| PATCH | `/admin/stations/:id/approve` | Approve station |
| PATCH | `/admin/stations/:id/reject` | Reject station |
| PATCH | `/admin/stations/:id/revoke` | Revoke station approval |
| DELETE | `/admin/stations/:id` | Delete station |
| GET | `/admin/cash-withdrawals` | List withdrawal requests |
| PATCH | `/admin/cash-withdrawals/:id/approve` | Approve request |
| PATCH | `/admin/cash-withdrawals/:id/reject` | Reject request |
| PATCH | `/admin/cash-withdrawals/:id/complete` | Mark transfer as completed |

---

## 8. Request and Response Examples

All endpoints follow a unified response structure:
- **Success response**: `{"success": true, "data": ...}` or `{"success": true, "message": "..."}`
- **Error response**: `{"success": false, "message": "error description"}`

Authenticated endpoints require the JWT bearer token in the `Authorization` header:
```http
Authorization: Bearer <jwt_access_token>
Content-Type: application/json
```

### 8.1 Authentication

#### Register user
`POST /api/v1/auth/register/user`

Request:
```json
{
  "name": "Monzer",
  "email": "user@example.com",
  "password": "example-password",
  "phone": "0990000000"
}
```

Response (`201 Created`):
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "name": "Monzer",
    "email": "user@example.com",
    "phone": "0990000000",
    "role": "user",
    "status": "active",
    "loyalty_points": 10,
    "created_at": "2026-11-14T05:00:00Z"
  }
}
```

#### Register station
`POST /api/v1/auth/register/station`

Request:
```json
{
  "name": "محطة النيل",
  "location": "الخرطوم",
  "bank_account_number": "123456789",
  "bank_account_name": "محمد أحمد",
  "bank_name": "bank_of_khartoum",
  "email": "station@example.com",
  "password": "example-password"
}
```

Response (`201 Created`):
```json
{
  "success": true,
  "message": "Station registration submitted",
  "data": {
    "id": 1,
    "name": "محطة النيل",
    "status": "pending"
  }
}
```

#### Login
`POST /api/v1/auth/login`

Issues `access_token` and `refresh_token` as `HttpOnly` cookies (`Set-Cookie`) and omits token strings from the JSON payload.

Request:
```json
{
  "email": "user@example.com",
  "password": "example-password"
}
```

Response (`200 OK`):
```json
{
  "success": true,
  "expiresIn": 86400,
  "user": {
    "id": 1,
    "role": "user",
    "name": "Monzer",
    "email": "user@example.com",
    "status": "active",
    "phone": "0990000000",
    "loyalty_points": 10
  }
}
```

#### Refresh token
`POST /api/v1/auth/refresh`

Reads the refresh token from the `refresh_token` `HttpOnly` cookie (or optional `refreshToken` JSON body fallback), rotates/issues a new `access_token` `HttpOnly` cookie, and returns the updated account summary.

Request (optional JSON fallback when cookies are not used):
```json
{
  "refreshToken": "4a7f28bc9e31d041c2e68f8a192b67..."
}
```

Response (`200 OK`):
```json
{
  "success": true,
  "expiresIn": 86400,
  "user": {
    "id": 1,
    "role": "user",
    "name": "Monzer",
    "email": "user@example.com",
    "status": "active",
    "loyalty_points": 10
  }
}
```

#### Current account profile
`GET /api/v1/auth/me`

Response (`200 OK`):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "role": "user",
    "name": "Monzer",
    "email": "user@example.com",
    "status": "active",
    "phone": "0990000000",
    "loyalty_points": 10
  }
}
```

#### Logout
`POST /api/v1/auth/logout`

Request:
```json
{
  "refreshToken": "4a7f28bc9e31d041c2e68f8a192b67..."
}
```

Response (`200 OK`):
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

### 8.2 User Operations

#### User dashboard
`GET /api/v1/user/dashboard`

Response (`200 OK`):
```json
{
  "success": true,
  "data": {
    "wallet": {
      "total_balance": 8000000,
      "allocated_balance": 3000000,
      "available_balance": 5000000,
      "loyalty_points": 20,
      "free_gas_gallons": 2,
      "loyalty_point_value": 100
    },
    "recent_transactions": [
      {
        "id": 1,
        "reference": "TRF-1731557000-0001",
        "type": "allocation",
        "amount": 500000,
        "from_name": "Monzer",
        "to_name": "محطة النيل",
        "branch_name": null,
        "status": "completed",
        "created_at": "2026-11-13T18:00:00Z"
      }
    ]
  }
}
```

#### User wallet
`GET /api/v1/user/wallet`

Response (`200 OK`):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "owner_type": "user",
    "owner_id": 1,
    "total_balance": 8000000,
    "allocated_balance": 3000000,
    "available_balance": 5000000,
    "loyalty_points": 20,
    "free_gas_gallons": 2,
    "loyalty_point_value": 100,
    "created_at": "2026-11-14T05:00:00Z"
  }
}
```

#### List and filter stations
`GET /api/v1/user/stations`
`GET /api/v1/user/stations?city=أم درمان`

Response (`200 OK`):
```json
{
  "success": true,
  "data": [
    {
      "station_id": 1,
      "station_name": "شركة النيل للبترول",
      "branch_id": 2,
      "branch_city": "أم درمان",
      "branch_address": "شارع الموردة، مقابل النادي الثقافي",
      "fuel_status": "low",
      "congestion_percentage": 80
    }
  ]
}
```

#### User transaction history
`GET /api/v1/user/transactions?type=allocation&status=completed&from_date=2026-11-01&to_date=2026-11-30`

Response (`200 OK`):
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "reference": "TRF-1731557000-0001",
      "type": "allocation",
      "amount": 500000,
      "from_name": "Monzer",
      "to_name": "محطة النيل",
      "branch_name": null,
      "status": "completed",
      "created_at": "2026-11-13T18:00:00Z"
    }
  ]
}
```

#### Export user transactions CSV
`GET /api/v1/user/transactions/export?from_date=2026-11-01&to_date=2026-11-30`

Response (`200 OK`, `Content-Type: text/csv; charset=utf-8`):
```csv
reference,type,amount,from_name,to_name,branch_name,status,created_at
TRF-1731557000-0001,allocation,500000,Monzer,محطة النيل,,completed,2026-11-13T18:00:00Z
```

### 8.3 Transfers and Allocations

> **Idempotency Requirement:** Every `POST /api/v1/transfers/*` endpoint (`topup`, `allocate`, `purchase`, `station-to-station`) requires a non-empty `Idempotency-Key` header (preferred) or `idempotency_key` JSON field (max 120 characters). Replaying the same key with identical parameters returns the original recorded result (`200 OK`) without re-executing the transfer; reusing the same key with different parameters returns `409 Conflict` (`"idempotency key already used with different request parameters"`).

#### Wallet top-up
`POST /api/v1/transfers/topup`

Request:
```json
{
  "amount": 1000000
}
```

Response (`200 OK`):
```json
{
  "success": true,
  "data": {
    "reference": "TOP-1731557000-0002",
    "amount": 1000000,
    "from_name": "",
    "to_name": "Monzer",
    "branch_name": null,
    "status": "completed",
    "points_earned": 5,
    "loyalty_points": 15
  }
}
```

#### Allocate balance to station
`POST /api/v1/transfers/allocate`

Allocations apply to the station's shared balance and do not require branch selection.

Request:
```json
{
  "station_id": 1,
  "amount": 500000
}
```

Response (`200 OK`):
```json
{
  "success": true,
  "data": {
    "reference": "TRF-1731557000-0003",
    "amount": 500000,
    "from_name": "Monzer",
    "to_name": "محطة النيل",
    "branch_name": null,
    "status": "completed"
  }
}
```

#### Purchase fuel at a branch
`POST /api/v1/transfers/purchase`

Records a branch-level `fuel_purchase` transaction from the user's `available` wallet balance (default, crediting the station's balance), `allocated` balance (previously allocated to stations), or `loyalty` points only (`0` SDG wallet charge). Optional `loyalty_points_used` must be a multiple of `10` (`1 point = 100 SDG`), and users may send `amount: 0` with `loyalty_points_used` (or `payment_source: "loyalty"`) to buy gas purely with loyalty points. Wallet-only purchases award `+5` loyalty points for regular users, while purchases redeeming loyalty points award `0` points. Rejects inactive stations and branches with `fuel_status = "empty"`. When the same user repeats an identical `(branch_id, amount)` purchase within 10 minutes, `duplicate_purchase: true` and `previous_reference` are returned; replaying the same `Idempotency-Key` returns `idempotent_replay: true` without re-charging.

Request:
```json
{
  "branch_id": 1,
  "amount": 50000,
  "payment_source": "available",
  "loyalty_points_used": 10
}
```

Response (`200 OK`):
```json
{
  "success": true,
  "data": {
    "reference": "PUR-1731557000-0003",
    "amount": 49000,
    "total_amount": 50000,
    "from_name": "Monzer",
    "to_name": "محطة النيل",
    "branch_name": "الخرطوم - شارع المطار",
    "status": "completed",
    "loyalty_points_used": 10,
    "loyalty_discount": 1000,
    "loyalty_points": 10
  }
}
```
*(For wallet-only purchases without `loyalty_points_used`, `data` includes `"points_earned": 5`. When redeeming points purely (`payment_source: "loyalty"` or `amount: 0` with `loyalty_points_used: 10`), `"amount"` is `0` and `"total_amount"` equals `"loyalty_discount"`. When repeating the same branch and amount within 10 minutes, `data` also includes `"duplicate_purchase": true` and `"previous_reference"`. On an idempotent retry with the same key, `data` includes `"idempotent_replay": true` and `"points_earned"` is `0`.)*

#### Transfer between stations
`POST /api/v1/transfers/station-to-station`

Transfers move balance between two stations' shared balances. No branch selection is involved.

Request:
```json
{
  "from_station_id": 1,
  "to_station_id": 2,
  "amount": 200000
}
```

Response (`200 OK`):
```json
{
  "success": true,
  "data": {
    "reference": "STF-1731557000-0004",
    "amount": 200000,
    "from_name": "محطة النيل",
    "to_name": "محطة الخرطوم",
    "branch_name": null,
    "status": "completed"
  }
}
```

#### Transfer history
`GET /api/v1/transfers`

Response (`200 OK`):
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "reference": "TRF-1731557000-0003",
      "type": "allocation",
      "amount": 500000,
      "from_name": "Monzer",
      "to_name": "محطة النيل",
      "branch_name": null,
      "status": "completed"
    }
  ]
}
```

#### Transfer details
`GET /api/v1/transfers/{id}`

Response (`200 OK`):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "reference": "TRF-1731557000-0003",
    "type": "allocation",
    "amount": 500000,
    "from_type": "user",
    "from_id": 1,
    "from_name": "Monzer",
    "to_type": "station",
    "to_id": 1,
    "to_name": "محطة النيل",
    "branch_id": null,
    "branch_name": null,
    "loyalty_points_used": 0,
    "loyalty_discount": 0,
    "status": "completed",
    "created_at": "2026-11-13T18:00:00Z"
  }
}
```

### 8.4 Station Operations

#### Station dashboard
`GET /api/v1/station/dashboard`

Response (`200 OK`):
```json
{
  "success": true,
  "data": {
    "balance": 1500000,
    "branches_count": 2,
    "branches": [
      {
        "id": 1,
        "city": "الخرطوم",
        "address": "شارع المطار",
        "fuel_status": "full",
        "congestion_percentage": 75
      },
      {
        "id": 2,
        "city": "أم درمان",
        "address": "شارع الوادي",
        "fuel_status": "low",
        "congestion_percentage": 40
      }
    ],
    "recent_transactions": []
  }
}
```

#### Station wallet
`GET /api/v1/station/wallet`

Response (`200 OK`):
```json
{
  "success": true,
  "data": {
    "id": 2,
    "owner_type": "station",
    "owner_id": 1,
    "total_balance": 1500000,
    "allocated_balance": 0,
    "available_balance": 1500000,
    "created_at": "2026-11-14T05:00:00Z"
  }
}
```

#### List station branches
`GET /api/v1/station/branches`

Response (`200 OK`):
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "station_id": 1,
      "city": "الخرطوم",
      "address": "شارع المطار",
      "is_main": true,
      "fuel_status": "full",
      "created_at": "2026-11-14T05:00:00Z"
    }
  ]
}
```

#### Add a branch
`POST /api/v1/station/branches`

Request:
```json
{
  "city": "أم درمان",
  "address": "شارع الوادي"
}
```

Response (`201 Created`):
```json
{
  "success": true,
  "data": {
    "id": 2,
    "station_id": 1,
    "city": "أم درمان",
    "address": "شارع الوادي",
    "is_main": false,
    "fuel_status": "full",
    "created_at": "2026-11-14T05:15:00Z"
  }
}
```

#### Update branch fuel status
`PATCH /api/v1/station/branches/{id}/fuel-status`

Request:
```json
{
  "fuel_status": "low"
}
```

Response (`200 OK`):
```json
{
  "success": true,
  "message": "Fuel status updated successfully"
}
```

#### Update branch congestion percentage
`PATCH /api/v1/station/branches/{id}/congestion`

Request:
```json
{
  "congestion_percentage": 65
}
```

Response (`200 OK`):
```json
{
  "success": true,
  "message": "Congestion percentage updated successfully"
}
```

#### Submit cash-withdrawal request
`POST /api/v1/station/cash-withdrawals`

Requires an `Idempotency-Key` header (preferred) or `idempotency_key` body field (max 120 characters); replays with identical parameters return the existing pending/settled withdrawal, and mismatched parameter replays return `409 Conflict`.

Request:
```json
{
  "amount": 300000
}
```

Response (`200 OK`):
```json
{
  "success": true,
  "data": {
    "id": 10,
    "amount": 300000,
    "status": "pending"
  }
}
```

#### Station withdrawal requests list
`GET /api/v1/station/cash-withdrawals`

Response (`200 OK`):
```json
{
  "success": true,
  "data": [
    {
      "id": 10,
      "station_id": 1,
      "amount": 300000,
      "status": "pending",
      "created_at": "2026-11-14T05:20:00Z"
    }
  ]
}
```

#### Export station transactions CSV
`GET /api/v1/station/transactions/export?from_date=2026-11-01&to_date=2026-11-30`

Response (`200 OK`, `Content-Type: text/csv; charset=utf-8`):
```csv
reference,type,amount,from_name,to_name,branch_name,status,created_at
TRF-1731557000-0003,allocation,500000,Monzer,محطة النيل,,completed,2026-11-13T18:00:00Z
```

### 8.5 Admin Operations

#### Admin dashboard
`GET /api/v1/admin/dashboard`

Response (`200 OK`):
```json
{
  "success": true,
  "data": {
    "users_count": 45,
    "registered_stations_count": 12,
    "approved_stations_count": 10,
    "pending_stations_count": 2,
    "pending_withdrawals_count": 3
  }
}
```

#### List users
`GET /api/v1/admin/users`

Response (`200 OK`):
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Monzer",
      "email": "user@example.com",
      "phone": "0990000000",
      "role": "user",
      "status": "active",
      "created_at": "2026-11-14T05:00:00Z"
    }
  ]
}
```

#### Update user status
`PATCH /api/v1/admin/users/{id}/status`

Request:
```json
{
  "status": "suspended"
}
```

Response (`200 OK`):
```json
{
  "success": true,
  "message": "User status updated successfully"
}
```

#### Update user role
`PATCH /api/v1/admin/users/{id}/role`

Request:
```json
{
  "role": "admin"
}
```

Response (`200 OK`):
```json
{
  "success": true,
  "message": "User role updated successfully"
}
```

#### List pending station registrations
`GET /api/v1/admin/stations/pending`

Response (`200 OK`):
```json
{
  "success": true,
  "data": [
    {
      "id": 2,
      "name": "محطة النيل",
      "location": "الخرطوم",
      "bank_account_number": "123456789",
      "bank_account_name": "محمد أحمد",
      "bank_name": "bank_of_khartoum",
      "email": "station@example.com",
      "status": "pending",
      "created_at": "2026-11-14T05:00:00Z"
    }
  ]
}
```

#### List all stations
`GET /api/v1/admin/stations`

Response (`200 OK`):
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "محطة النيل",
      "location": "الخرطوم",
      "bank_account_number": "123456789",
      "bank_account_name": "محمد أحمد",
      "bank_name": "bank_of_khartoum",
      "email": "station@example.com",
      "status": "active",
      "branches": [
        {
          "id": 1,
          "city": "الخرطوم",
          "address": "شارع المطار",
          "is_main": true,
          "fuel_status": "full"
        }
      ]
    }
  ]
}
```

#### Approve station registration
`PATCH /api/v1/admin/stations/{id}/approve`

Response (`200 OK`):
```json
{
  "success": true,
  "message": "Station approved successfully"
}
```

#### Reject station registration
`PATCH /api/v1/admin/stations/{id}/reject`

Request:
```json
{
  "rejection_reason": "Incomplete registration documents"
}
```

Response (`200 OK`):
```json
{
  "success": true,
  "message": "Station rejected"
}
```

#### Revoke station approval
`PATCH /api/v1/admin/stations/{id}/revoke`

Response (`200 OK`):
```json
{
  "success": true,
  "message": "Station approval revoked"
}
```

#### Delete station
`DELETE /api/v1/admin/stations/{id}`

Response (`200 OK`):
```json
{
  "success": true,
  "message": "Station deleted successfully"
}
```

#### List cash-withdrawal requests
`GET /api/v1/admin/cash-withdrawals`

Response (`200 OK`):
```json
{
  "success": true,
  "data": [
    {
      "id": 10,
      "station_id": 1,
      "station_name": "محطة النيل",
      "amount": 300000,
      "status": "pending",
      "created_at": "2026-11-14T05:20:00Z"
    }
  ]
}
```

#### Approve withdrawal request
`PATCH /api/v1/admin/cash-withdrawals/{id}/approve`

Response (`200 OK`):
```json
{
  "success": true,
  "message": "Cash withdrawal approved"
}
```

#### Reject withdrawal request
`PATCH /api/v1/admin/cash-withdrawals/{id}/reject`

Request:
```json
{
  "rejection_reason": "Insufficient station reserve"
}
```

Response (`200 OK`):
```json
{
  "success": true,
  "message": "Cash withdrawal rejected"
}
```

#### Complete withdrawal
`PATCH /api/v1/admin/cash-withdrawals/{id}/complete`

Request (optional body):
```json
{
  "external_reference": "BOK-20261114-9981"
}
```

Response (`200 OK`):
```json
{
  "success": true,
  "data": {
    "id": 10,
    "status": "completed"
  }
}
```

Only withdrawals in `approved` status can be completed; `external_reference` is optional and persisted on the withdrawal record when provided.

### 8.6 Common Error Responses

All error responses adhere to the standard JSON structure:

```json
{
  "success": false,
  "message": "Detailed error description"
}
```

| HTTP Status | Condition | Example `message` |
|---|---|---|
| `400 Bad Request` | Missing required fields, missing/invalid idempotency key, invalid input format | `"All fields are required"`, `"amount must be greater than zero"`, `"idempotency_key is required"` |
| `401 Unauthorized` | Invalid/expired JWT token, invalid credentials | `"Invalid email or password"`, `"Unauthorized"`, `"Invalid or expired token"` |
| `402 Payment Required` | Insufficient balance for allocation or transfer | `"Insufficient balance"` |
| `403 Forbidden` | Insufficient role, account suspended, station not approved | `"Forbidden"`, `"Your account is suspended"`, `"Station registration is pending approval"` |
| `404 Not Found` | Entity (user, station, branch, withdrawal) not found | `"User not found"`, `"Station not found"`, `"Transfer not found"` |
| `409 Conflict` | Unique email conflict, idempotency parameter mismatch, duplicate withdrawal settlement, empty branch on purchase | `"Email already registered"`, `"idempotency key already used with different request parameters"`, `"Withdrawal already settled"`, `"selected branch is out of fuel"` |
| `422 Unprocessable Entity` | Station is not active | `"Station is not active"` |
| `429 Too Many Requests` | IP rate limit exceeded | `"Rate limit exceeded"` |
| `500 Internal Server Error` | Database connection or unhandled server error | `"Internal server error"`, `"Database is temporarily unavailable"` |

---

## 9. CSV Transaction Exports

Both user and station accounts can export their own transaction history.

| Account | Endpoint | Scope |
|---------|----------|-------|
| User | `GET /api/v1/user/transactions/export` | Transactions visible to the authenticated user |
| Station | `GET /api/v1/station/transactions/export` | Transactions associated with the authenticated station |

Suggested query filters supported by both endpoints:
- `from_date`
- `to_date`
- `type`
- `status`

Example:

`GET /api/v1/user/transactions/export?from_date=2026-11-01&to_date=2026-11-30`

Response:
- Content-Type: `text/csv; charset=utf-8`
- Content-Disposition: `attachment; filename="transactions.csv"`

Suggested CSV columns:

```csv
reference,type,amount,from_name,to_name,branch_name,status,created_at
TRF-001,allocation,500000,محمد أحمد,محطة النيل,,completed,2026-11-13T18:00:00Z
```

CSV export must be scoped to the authenticated account on the server; a user or station must not be able to export another account's transactions.

---

## 10. Important Implementation Rules

- User registration `phone` must be digits only with a maximum of 10 digits.
- Regular users (`role = 'user'`, excluding stations and admins) earn `10` welcome loyalty points at registration and `+5` loyalty points on each wallet top-up or wallet-only fuel purchase (`10 points = 1 gallon of fuel`, `1 point = 100 SDG`). Loyalty points can only be redeemed in multiples of `10` (`10, 20, 30...`), either as a discount or as a `0 SDG` points-only fuel purchase (`payment_source = 'loyalty'` or `amount = 0` with `loyalty_points_used`); purchases that redeem loyalty points award `0` points.
- Fuel purchases (`fuel_purchase`) are triggered by tapping a non-empty station branch on the Stations page (not from Dashboard quick actions or Transfers) and use `available` balance (credits station wallet), `allocated` balance (deducts from allocated balance without double-crediting the station), or `loyalty` points only (`0 SDG` cash charge, credits station only for any remaining wallet charge); branches with `fuel_status = 'empty'` are rejected (`409`), and repeating an identical purchase at the same branch within 10 minutes warns and asks confirmation before submitting (`duplicate_purchase` / `idempotent_replay`).
- Station-to-station transfers require the user's `allocated_balance >= amount`.
- Branch city filtering queries `branches.city` directly from the database (`LOWER(branches.city) = LOWER(?)`) without hardcoded backend lookup lists.
- Transaction and withdrawal lists are returned and rendered newest to oldest (`created_at DESC, id DESC`).
- Email addresses must be unique across user and station accounts.
- Admin is created/assigned administratively; it is not a public registration choice.
- A station is inactive/pending until approved by admin.
- One station account manages one station.
- One station has one main branch created at registration.
- All branches of a station share one wallet/balance.
- Allocations and station-to-station transfers select stations only, never branches.
- A transaction may record a branch when it actually took place at a particular branch.
- Station and user transaction exports are CSV downloads.
- Store passwords as secure bcrypt hashes, hash refresh tokens with SHA-256, and issue browser session tokens in `HttpOnly` cookies (never exposing tokens in JSON responses or `localStorage`).
- Use database transactions (`SELECT ... FOR UPDATE` in deadlock-free lock order: `station` before `user`, lower `station_id` first on station transfers), enforce the PostgreSQL `chk_wallets_non_negative` and `chk_users_loyalty_non_negative` constraints, reserve open (`pending` and `approved`) cash withdrawals against `available_balance`, and require idempotency keys (`Idempotency-Key` header or `idempotency_key` body field, unique per actor/station + SHA-256 `request_hash`) on all transfer and cash-withdrawal mutations to prevent negative balances and duplicate financial requests.
- Enforce role-based authorization, per-IP rate limiting, security hardening headers, and strict environment-based configuration on the backend.

---

## 11. Page Summary

| Account | Pages |
|---------|-------|
| **User** | Dashboard, Stations, Transfers, All Transactions |
| **Admin** | Dashboard, User Management, Station Management, Cash Withdrawals |
| **Station** | Dashboard, Branch Management, All Transactions |

---

## 12. Unspecified Behavior

The following details have intentionally not been invented:
- Automated external sensor/telemetry calculation for branch congestion (values are entered and updated by the station).
- The exact external process/payment provider for wallet top-ups.
- Any additional cash-withdrawal settlement or balance rules beyond the specified admin review and completion flow.
