# Golang Backend

Production-ready Gin backend with clean architecture, Prisma migrations, and custom query builder.

## 📁 Project Structure

```
golang-backend/
├── config/                    # Configuration
│   └── config.go             # Viper config loader with logger config
├── internal/
│   ├── handler/              # HTTP handlers (controllers)
│   │   ├── user_handler.go
│   │   ├── address_handler.go
│   │   └── routes.go
│   ├── usecase/              # Business logic layer
│   │   ├── user_usecase.go
│   │   └── address_usecase.go
│   ├── repository/           # Data access layer
│   │   ├── user_repository.go
│   │   └── address_repository.go
│   ├── model/                # Domain models
│   │   ├── user.go
│   │   └── address.go
│   ├── dto/                  # Data Transfer Objects (request/response)
│   │   ├── user_dto.go
│   │   └── address.dto.go
│   └── middleware/           # Custom middlewares
│       ├── auth.go          # JWT authentication
│       ├── logger.go        # Structured HTTP request logger
│       ├── recovery.go      # Panic recovery
│       └── role.go          # Role-based access control
├── pkg/
│   ├── database/             # Database connection & query builder
│   │   ├── db.go
│   │   ├── query_builder.go
│   │   └── README.md
│   ├── logger/               # Zap structured logger (Grafana Loki ready)
│   │   └── logger.go
│   ├── email/                # Email service (SMTP)
│   │   └── email.go
│   ├── validator/            # Request validation
│   │   └── validator.go
│   └── response/             # Standard API response
│       └── response.go
├── prisma/                   # Prisma schema & migrations
│   ├── schema.prisma         # Database schema
│   └── migrations/           # Generated migration files
│       ├── migration_lock.toml
│       └── 20XXXXXX_*.sql
├── monitoring/               # Grafana monitoring stack configs
│   ├── loki-config.yaml     # Loki configuration
│   ├── promtail-config.yaml # Promtail log collector config
│   └── grafana-datasources.yaml
├── logs/                     # Application logs (auto-generated)
│   ├── .gitkeep
│   └── app.log              # JSON formatted logs for Loki
├── docs/                     # Swagger documentation (auto-generated)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── bin/                      # Compiled binaries
├── tmp/                      # Air hot reload temp files
├── docker-compose.yml        # Grafana + Loki + Promtail stack
├── setup-monitoring.sh       # Quick setup script for monitoring
├── .env                      # Environment variables
├── .env.example              # Environment variables template
├── .gitignore
├── go.mod
├── go.sum
├── main.go                   # Application entry point
├── README.md                 # This file
└── GRAFANA_GUIDE.md          # Complete Grafana usage guide
```

## 🚀 Features

-    ✅ **Clean Architecture**: Handler → Usecase → Repository pattern
-    ✅ **Prisma**: Schema management & migrations
-    ✅ **Custom Query Builder**: Fluent SQL builder with audit trail support
-    ✅ **UUID**: UUID-based identifiers throughout
-    ✅ **Audit Trail**: Auto-tracking created_by, updated_by, deleted_by
-    ✅ **Soft Delete**: Built-in soft delete support
-    ✅ **JWT Authentication**: Access + refresh tokens with role-based authorization
-    ✅ **Email Verification**: SMTP email verification with HTML templates
-    ✅ **User Roles**: SUPERADMIN, ADMIN, USER with middleware protection
-    ✅ **Structured Logging**: Zap logger with JSON format & request ID tracking
-    ✅ **Grafana Loki Integration**: Centralized logging & monitoring with real-time log aggregation
-    ✅ **Log Management**: File rotation, compression, and retention (28 days default)
-    ✅ **Monitoring Stack**: Pre-configured Grafana + Loki + Promtail via Docker Compose
-    ✅ **Validation**: go-playground/validator with custom error formatting
-    ✅ **Middleware**: Recovery, CORS, RequestID, Logger, JWT Auth, Role-based Auth
-    ✅ **Swagger**: Auto-generated API documentation
-    ✅ **Graceful Shutdown**: Proper server shutdown handling
-    ✅ **Hot Reload**: Air for development (already configured)
-    ✅ **Standard Response**: Consistent JSON response format
-    ✅ **Production Ready**: Environment-based configuration with .env support

## 📦 Dependencies

### Core Framework & Libraries

-    **Gin** (`github.com/gin-gonic/gin`): HTTP web framework
-    **Viper** (`github.com/spf13/viper`): Configuration management from .env
-    **database/sql + lib/pq**: PostgreSQL driver

### Authentication & Security

-    **JWT** (`github.com/golang-jwt/jwt/v5`): Token-based authentication
-    **Bcrypt** (`golang.org/x/crypto`): Password hashing

### Database & ORM

-    **Prisma**: Schema management & migrations (Node.js)
-    **Custom Query Builder**: Fluent SQL builder with audit trail support

### Logging & Monitoring

-    **Zap** (`go.uber.org/zap`): High-performance structured logging with JSON format
-    **Grafana Loki**: Log aggregation system (Docker)
-    **Promtail**: Log collector & shipper to Loki (Docker)
-    **Grafana**: Metrics visualization & log exploration (Docker)

### Validation & Response

-    **Validator** (`github.com/go-playground/validator/v10`): Request validation
-    **Standard Response**: Consistent JSON response wrapper

### Middleware

-    **CORS** (`github.com/gin-contrib/cors`): Cross-origin resource sharing
-    **RequestID** (`github.com/gin-contrib/requestid`): Request tracking
-    **Custom Middleware**: Recovery, Logger, JWT Auth, Role-based Auth

### Email & Communication

-    **SMTP** (`net/smtp`): Email sending with HTML templates

### Documentation

-    **Swagger** (`github.com/swaggo/swag`): Auto-generated API documentation
-    **Gin Swagger** (`github.com/swaggo/gin-swagger`): Swagger UI integration

### Development Tools

-    **Air**: Hot reload for development
-    **godotenv** (`github.com/joho/godotenv`): Load .env files

## ⚙️ Setup

1. **Install dependencies:**

```bash
go mod tidy
npm install -D prisma
```

2. **Setup database:**

     - Install PostgreSQL
     - Update `.env` with your database credentials (database akan auto-create saat migration)

3. **Copy environment file:**

```bash
cp .env.example .env
```

4. **Edit `.env`:**

```env
# Server Configuration
APP_NAME=golang Backend
APP_ENV=development
APP_PORT=8080
APP_DEBUG=true

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=golang_db
DB_SSL_MODE=disable
DB_TIMEZONE=Asia/Jakarta

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRED_HOUR=24

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# SMTP Configuration (Gmail example)
SMTP_EMAIL=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_FROM_NAME="golang Backend"
SMTP_FROM_EMAIL=your-email@gmail.com

# Logger Configuration
LOG_TO_FILE=true
LOG_FILE_PATH=logs/app.log
LOG_MAX_SIZE=100
LOG_MAX_BACKUPS=3
LOG_MAX_AGE=28
LOG_COMPRESS=true
```

5. **Run Prisma migrations:**

```bash
npx prisma migrate dev --name init
npx prisma generate
```

6. **Run the application:**

```bash
# Development mode (dengan Air hot reload)
air

# Production mode
go run main.go
```

## 🔧 Custom Query Builder

Project ini menggunakan custom SQL query builder dengan fitur:

### QueryBuilder (SELECT)

```go
// Basic query
result := database.NewQueryBuilder("users").
    Select("id", "email", "name").
    Where("is_active = ?", true).
    OrderBy("created_at DESC").
    Limit(10).
    BuildResult()

rows, err := db.Query(result.Query, result.Args...)

// With joins
result := database.NewQueryBuilder("users u").
    Select("u.id", "u.name", "a.full_address").
    Join("LEFT JOIN addresses a ON a.user_id = u.id").
    Where("u.deleted_at IS NULL").
    BuildResult()
```

### InsertBuilder (INSERT dengan audit trail)

```go
// Insert dengan auto audit trail
userID := "uuid-dari-jwt-context"

result := database.NewInsertBuilder("users").
    Columns("id", "email", "password", "name", "is_active").
    Values(newID, email, hashedPassword, name, true).
    SetCreatedBy(userID). // Auto add created_by + created_at
    BuildResult()

_, err := db.Exec(result.Query, result.Args...)
```

### UpdateBuilder (UPDATE dengan audit trail)

```go
// Update dengan auto audit trail
result := database.NewUpdateBuilder("users").
    Set("name", "John Doe").
    Set("email", "john@example.com").
    Where("id = ?", userID).
    SetUpdatedBy(currentUserID). // Auto add updated_by + updated_at
    BuildResult()

_, err := db.Exec(result.Query, result.Args...)
```

### DeleteBuilder (Soft/Hard Delete)

```go
// Soft delete (default) - UPDATE deleted_at
result := database.NewDeleteBuilder("users").
    Where("id = ?", userID).
    SetDeletedBy(currentUserID). // Auto add deleted_by + deleted_at
    BuildResult()

// Hard delete - actual DELETE
result := database.NewDeleteBuilder("users").
    Where("id = ?", userID).
    HardDelete(). // DELETE FROM users WHERE...
    BuildResult()

_, err := db.Exec(result.Query, result.Args...)
```

### Raw Query Helpers

```go
// Raw query with multiple rows
rows, err := database.RawQuery(db,
    "SELECT * FROM users WHERE email LIKE ?",
    "%@example.com")

// Raw exec (INSERT/UPDATE/DELETE)
result, err := database.RawExec(db,
    "UPDATE users SET is_active = ? WHERE id = ?",
    false, userID)

// Raw query single row
row := database.RawQueryRow(db,
    "SELECT id, email FROM users WHERE id = ?",
    userID)
```

### Audit Trail Fields

Semua table memiliki audit trail otomatis:

-    `created_by` (UUID) - User yang membuat record
-    `created_at` (timestamp) - Waktu dibuat
-    `updated_by` (UUID) - User yang update terakhir
-    `updated_at` (timestamp) - Waktu update terakhir
-    `deleted_by` (UUID) - User yang soft delete
-    `deleted_at` (timestamp) - Waktu soft delete

## 📖 API Documentation

#### User Model

-    UUID-based ID
-    Email (unique) + password (bcrypt hashed)
-    Name
-    Role (SUPERADMIN, ADMIN, USER) - default: USER
-    Refresh token + token expiry
-    Verification token + verification expiry (for email verification)
-    Is active flag (default: false - requires email verification)
-    Full audit trail (created_by, updated_by, deleted_by)
-    One-to-many relationship dengan Address
-    Email (unique) + password (bcrypt hashed)
-    Name
-    Refresh token + token expiry
-    Is active flag
-    Full audit trail (created_by, updated_by, deleted_by)
-    One-to-many relationship dengan Address

#### Address Model

-    UUID-based ID
-    User relationship (many-to-one)
-    Complete address fields: label, recipient_name, phone
-    Location: province, city, district, sub_district, postal_code
-    Full address (text)
-    Is primary flag (untuk set alamat utama)
-    Is active flag

### Available Endpoints

#### Auth (Public)

-    `POST /api/auth/register` - Register user baru (kirim email verification)
-    `GET /api/auth/verify-email?token=xxx` - Verify email address
-    `POST /api/auth/login` - Login user (requires verified email)
-    `POST /api/auth/refresh` - Refresh access token
-    `POST /api/auth/logout` - Logout user (clear refresh token)

#### Users (Protected - butuh Bearer Token)

-    `GET /api/users/profile` - Get current user profile
-    `PUT /api/users/profile` - Update user profile
-    `PUT /api/users/change-password` - Change password
-    `GET /api/users` - Get all users (Admin only)
-    `DELETE /api/users/:id` - Delete user (SuperAdmin only)

#### Addresses (Protected - butuh Bearer Token)

-    `GET /api/addresses` - Get all addresses for authenticated user

## 🔐 Authentication

### Email Verification Flow

1. **Register** → User dibuat dengan `is_active = false`
2. **Email Sent** → Verification email dikirim otomatis (background goroutine)
3. **Verify Email** → User klik link di email → `is_active = true`
4. **Login** → User bisa login setelah email verified

### JWT Authentication

Semua endpoint di bawah `/api/users` dan `/api/addresses` membutuhkan JWT token.

**Token Types:**

-    **Access Token**: Short-lived (sesuai config `EXPIRED_HOUR`)
-    **Refresh Token**: Long-lived (7 hari)

**Header format:**

```
Authorization: Bearer <your-access-token>
```

### Role-Based Authorization

**User Roles:**

-    `SUPERADMIN` - Full access to all endpoints
-    `ADMIN` - Can manage users and data
-    `USER` - Basic user access (default)

**Protected Endpoints:**

-    `GET /api/users` - Requires ADMIN role

## 📝 Example Request

### Register (with Email Verification)

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'

# Response: User created, check email for verification link
# Email berisi link: http://localhost:8080/api/auth/verify-email?token=xxx
```

### Verify Email

```bash
# Klik link di email atau:
curl -X GET "http://localhost:8080/api/auth/verify-email?token=your-verification-token"

# Response: Email verified successfully. You can now login.
```

### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# Response: access_token, refresh_token, user info
```

### Refresh Token

```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your-refresh-token"
  }'

# Response: new access_token and refresh_token
```

### Logout

```bash
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer <your-access-token>"

# Response: Logged out successfully
```

### Get Profile (Protected)

```bash
curl -X GET http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer <your-token>"
```

### Create Address (Protected)

````bash
curl -X POST http://localhost:8080/api/addresses \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "label": "Home",
    "recipient_name": "John Doe",
    "phone": "081234567890",
    "province": "DKI Jakarta",
    "city": "Jakarta Selatan",
    "district": "Kebayoran Baru",
    "sub_district": "Senayan",
    "postal_code": "12190",
    "full_address": "Jl. Sudirman No. 123",
    "is_primary": true
  }'
```Login via `/api/auth/login` → dapat access_token + refresh_token
3. Gunakan access_token di header untuk endpoint protected
4. Jika access_token expired, gunakan refresh_token untuk get new access_token

## 📝 Example Request

### Register

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'
````

### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Get Profile (Protected)

```bash
curl -X GET http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer <your-token>"
```

## 🔄 Update Swagger Docs

Setiap kali update handler dengan comment swagger:

```bash
swag init
```

## 🛠️ Development

### Hot Reload

Project sudah dikonfigurasi dengan Air. Cukup jalankan:

```bash
air
```

### Add New Module

1. Buat model di `prisma/schema.prisma`
2. Run migration: `npx prisma migrate dev --name add_your_model`
3. Buat Go model di `internal/model/`
4. Buat DTO di `internal/dto/`
5. Buat repository di `internal/repository/` (gunakan query builder)
6. Buat usecase di `internal/usecase/`
7. Buat handler di `internal/handler/`
8. Register routes di `main.go`
9. Run `swag init` untuk update docs

### Example: Add New Table

Template sudah include **User** dan **Address** models. Untuk add table baru:

Edit `prisma/schema.prisma`:

````prisma
model Product {
  id          String    @id @default(uuid()) @db.Uuid
  name        String    @db.VarChar(255)
  price       Decimal   @db.Decimal(10, 2)
  stock       Int       @default(0)
  isActive    Boolean   @default(true) @map("is_active")
  createdAt   DateTime  @default(now()) @map("created_at")
## 📊 Monitoring dengan Grafana & Loki

Project ini sudah dikonfigurasi dengan **Grafana Loki** untuk centralized logging dan monitoring.

### Setup Monitoring Stack

1. **Start Grafana, Loki & Promtail:**

```bash
docker-compose up -d
```

2. **Akses Grafana:**
   - URL: http://localhost:3000
   - Username: `admin`
   - Password: `admin`

3. **View Logs di Grafana:**
   - Buka **Explore** menu (ikon kompas)
   - Pilih **Loki** sebagai data source
   - Query contoh:
     ```logql
     {job="golang-backend"}
     ```
   - Filter by level:
     ```logql
     {job="golang-backend", level="error"}
     ```
   - Search by message:
     ```logql
     {job="golang-backend"} |= "HTTP Request"
     ```

### Log Format

Logger menggunakan **JSON format** untuk kompatibilitas dengan Loki:

```json
{
  "timestamp": "2026-01-06T10:30:45.123Z",
  "level": "info",
  "message": "HTTP Request",
  "method": "GET",
  "path": "/api/users/profile",
  "status": 200,
  "ip": "127.0.0.1",
  "latency_ms": 15,
  "request_id": "abc123"
}
```

### Environment Variables untuk Logger

Tambahkan di `.env`:

```env
LOG_TO_FILE=true
LOG_FILE_PATH=logs/app.log
LOG_MAX_SIZE=100       # MB
LOG_MAX_BACKUPS=3      # jumlah backup files
LOG_MAX_AGE=28         # days
LOG_COMPRESS=true      # compress rotated logs
```

### Monitoring Features

- ✅ **Structured JSON Logging**: Semua log dalam format JSON
- ✅ **Log Aggregation**: Promtail collect logs → Loki store → Grafana visualize
- ✅ **Multiple Output**: Console (colored) + File (JSON)
- ✅ **Log Levels**: Debug, Info, Warn, Error, Fatal
- ✅ **HTTP Request Tracking**: Method, path, status, latency, IP, user agent
- ✅ **Request ID**: Trace individual requests
- ✅ **Error Tracking**: Auto-log errors dengan level sesuai status code
- ✅ **Log Retention**: 31 days default

### Query Examples di Grafana

**1. Semua HTTP errors (4xx & 5xx):**
```logql
{job="golang-backend"} | json | status >= 400
```

**2. Slow requests (>1 second):**
```logql
{job="golang-backend"} | json | latency_ms > 1000
```

**3. Requests per endpoint:**
```logql
sum by (path) (count_over_time({job="golang-backend"}[5m]))
```

**4. Error rate:**
```logql
sum(rate({job="golang-backend", level="error"}[5m]))
```

### Stop Monitoring Stack

```bash
docker-compose down
```

Untuk hapus data juga:
```bash
docker-compose down -v
```

## 🎯 Next Steps

-    [x] Add role-based access control (RBAC) ✅
-    [x] Add email verification ✅
-    [x] Add refresh token system ✅
-    [x] Add Grafana Loki monitoring ✅
-    [ ] Add rate limiting
-    [ ] Add caching (Redis)
-    [ ] Add unit tests
-    [ ] Add Docker support for app
-    [ ] Add CI/CD pipeline
-    [ ] Add Prometheus metrics
-    [ ] Add distributed tracing
-    [ ] Add forgot password / reset password
-    [ ] Add resend verification email
Run migration:

```bash
npx prisma migrate dev --name add_product_table
````

## 📊 Database Migration

Project ini menggunakan **Prisma** untuk database migrations.

### Create Migration

Setelah update `prisma/schema.prisma`:

```bash
npx prisma migrate dev --name migration_name
```

### Apply Migration (Production)

```bash
npx prisma migrate deploy
```

### Reset Database (Development only!)

```bash
npx prisma migrate reset
```

### View Migration Status

```bash
npx prisma migrate status
```

### Manual Migration

1. Edit `prisma/schema.prisma`
2. Run `npx prisma migrate dev --name your_migration_name`
3. Run `npx prisma generate` (optional, untuk Prisma Client)
4. Prisma akan auto-generate SQL migration files

**Note:** Database akan auto-created saat first migration jika belum exist.

## 🎯 Next Steps

-    [ ] Add role-based access control (RBAC)
-    [ ] Add rate limiting
-    [ ] Add caching (Redis)
-    [ ] Add unit tests
-    [ ] Add Docker support
-    [ ] Add CI/CD pipeline
-    [ ] Add Prometheus metrics
-    [ ] Add distributed tracing

## 📄 License

MIT
