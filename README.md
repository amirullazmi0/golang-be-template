# Kratify Backend

Production-ready Gin backend with clean architecture, Prisma migrations, and custom query builder.

## 📁 Project Structure

```
root-backend/
├── config/              # Configuration
│   └── config.go       # Viper config loader
├── internal/
│   ├── handler/        # HTTP handlers (controllers)
│   ├── usecase/        # Business logic layer
│   ├── repository/     # Data access layer
│   ├── model/          # Domain models
│   ├── dto/            # Data Transfer Objects (request/response)
│   └── middleware/     # Custom middlewares (auth, logger, etc)
├── pkg/
│   ├── database/       # Database connection & query builder
│   ├── logger/         # Zap structured logger
│   ├── validator/      # Request validation
│   └── response/       # Standard API response
├── prisma/             # Prisma schema & migrations
│   ├── schema.prisma   # Database schema
│   └── migrations/     # Generated migrations
├── docs/               # Swagger documentation (auto-generated)
├── .env                # Environment variables
├── .env.example        # Environment variables template
└── main.go             # Application entry point
```

## 🚀 Features

-    ✅ **Clean Architecture**: Handler → Usecase → Repository pattern
-    ✅ **Prisma**: Schema management & migrations
-    ✅ **Custom Query Builder**: Fluent SQL builder with audit trail support
-    ✅ **UUID**: UUID-based identifiers throughout
-    ✅ **Audit Trail**: Auto-tracking created_by, updated_by, deleted_by
-    ✅ **Soft Delete**: Built-in soft delete support
-    ✅ **JWT Authentication**: Access + refresh tokens
-    ✅ **Structured Logging**: Zap logger with request ID
-    ✅ **Validation**: go-playground/validator with custom error formatting
-    ✅ **Middleware**: Recovery, CORS, RequestID, Logger, JWT Auth
-    ✅ **Swagger**: Auto-generated API documentation
-    ✅ **Graceful Shutdown**: Proper server shutdown handling
-    ✅ **Hot Reload**: Air for development (already configured)
-    ✅ **Standard Response**: Consistent JSON response format

## 📦 Dependencies

-    **Gin**: HTTP web framework
-    **Prisma**: Schema management & migrations
-    **database/sql + lib/pq**: PostgreSQL driver
-    **Viper**: Configuration management
-    **Zap**: Structured logging
-    **JWT**: Token-based authentication (golang-jwt/jwt/v5)
-    **Validator**: Request validation
-    **Swagger**: API documentation
-    **CORS**: Cross-origin resource sharing

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
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=kratify_db
JWT_SECRET=your-secret-key-change-this
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

Swagger UI tersedia di: **http://localhost:8080/swagger/index.html**

### Default Schema

Template ini sudah include 2 model utama:

#### User Model

-    UUID-based ID
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
-    Full audit trail

### Available Endpoints

#### Auth (Public)

-    `POST /api/auth/register` - Register user baru
-    `POST /api/auth/login` - Login user

#### Users (Protected - butuh Bearer Token)

-    `GET /api/users/profile` - Get current user profile
-    `PUT /api/users/profile` - Update user profile
-    `PUT /api/users/change-password` - Change password
-    `GET /api/users` - Get all users
-    `DELETE /api/users/:id` - Delete user

#### Health Check

-    `GET /health` - Check server status

## 🔐 Authentication

Semua endpoint di bawah `/api/users` membutuhkan JWT token.

**Token Types:**

-    **Access Token**: Short-lived (sesuai config `EXPIRED_HOUR`)
-    **Refresh Token**: Long-lived (7 hari)

**Header format:**

```
Authorization: Bearer <your-access-token>
```

**Cara pakai:**

1. Register user baru via `/api/auth/register` → dapat access_token + refresh_token
2. Login via `/api/auth/login` → dapat access_token + refresh_token
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
```

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

```prisma
model Product {
  id          String    @id @default(uuid()) @db.Uuid
  name        String    @db.VarChar(255)
  price       Decimal   @db.Decimal(10, 2)
  stock       Int       @default(0)
  isActive    Boolean   @default(true) @map("is_active")
  createdAt   DateTime  @default(now()) @map("created_at")
  updatedAt   DateTime  @updatedAt @map("updated_at")
  deletedAt   DateTime? @map("deleted_at")
  createdBy   String?   @map("created_by") @db.Uuid
  updatedBy   String?   @map("updated_by") @db.Uuid
  deletedBy   String?   @map("deleted_by") @db.Uuid

  @@map("products")
}
```

Run migration:

```bash
npx prisma migrate dev --name add_product_table
```

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
