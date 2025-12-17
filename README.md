# Kratify Backend

Production-ready Gin backend with clean architecture pattern.

## 📁 Project Structure

```
kratify-backend/
├── config/              # Configuration & database setup
│   ├── config.go       # Viper config loader
│   └── database.go     # GORM database connection
├── internal/
│   ├── handler/        # HTTP handlers (controllers)
│   ├── usecase/        # Business logic layer
│   ├── repository/     # Data access layer
│   ├── model/          # Domain models (GORM entities)
│   ├── dto/            # Data Transfer Objects (request/response)
│   └── middleware/     # Custom middlewares (auth, logger, etc)
├── pkg/
│   ├── logger/         # Zap structured logger
│   ├── validator/      # Request validation
│   └── response/       # Standard API response
├── docs/               # Swagger documentation (auto-generated)
├── migrations/         # Database migrations
├── .env                # Environment variables
├── .env.example        # Environment variables template
└── main.go             # Application entry point
```

## 🚀 Features

-    ✅ **Clean Architecture**: Handler → Usecase → Repository pattern
-    ✅ **GORM**: PostgreSQL with auto-migration
-    ✅ **JWT Authentication**: Secure token-based auth
-    ✅ **Structured Logging**: Zap logger with request ID
-    ✅ **Validation**: go-playground/validator with custom error formatting
-    ✅ **Middleware**: Recovery, CORS, RequestID, Logger, JWT Auth
-    ✅ **Swagger**: Auto-generated API documentation
-    ✅ **Graceful Shutdown**: Proper server shutdown handling
-    ✅ **Hot Reload**: Air for development (already configured)
-    ✅ **Standard Response**: Consistent JSON response format

## 📦 Dependencies

-    **Gin**: HTTP web framework
-    **GORM**: ORM for database operations
-    **Viper**: Configuration management
-    **Zap**: Structured logging
-    **JWT**: Token-based authentication
-    **Validator**: Request validation
-    **Swagger**: API documentation
-    **CORS**: Cross-origin resource sharing

## ⚙️ Setup

1. **Install dependencies:**

```bash
go mod tidy
```

2. **Setup database:**

     - Install PostgreSQL
     - Create database: `kratify_db`
     - Update `.env` with your database credentials

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

5. **Run the application:**

```bash
# Development mode (dengan Air hot reload)
air

# Production mode
go run main.go
```

## 📖 API Documentation

Swagger UI tersedia di: **http://localhost:8080/swagger/index.html**

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

**Header format:**

```
Authorization: Bearer <your-jwt-token>
```

**Cara pakai:**

1. Register user baru via `/api/auth/register`
2. Login via `/api/auth/login` → dapat token
3. Gunakan token di header untuk endpoint protected

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

1. Buat model di `internal/model/`
2. Buat DTO di `internal/dto/`
3. Buat repository di `internal/repository/`
4. Buat usecase di `internal/usecase/`
5. Buat handler di `internal/handler/`
6. Register routes di `main.go`
7. Run `swag init` untuk update docs

## 📊 Database Migration

Auto migration akan berjalan otomatis saat aplikasi start.

Untuk manual migration, update di `main.go`:

```go
db.AutoMigrate(&model.User{}, &model.YourNewModel{})
```

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
