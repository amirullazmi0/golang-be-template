---
description: Tambahkan fitur/modul baru mengikuti pola clean architecture Kratify
---

# Menambahkan Fitur Baru

Ikuti panduan ini untuk menambahkan modul fitur baru (misal: "Order", "Cart", "Review") ke Kratify Backend.

## Gambaran Arsitektur

```
Handler (HTTP) → Usecase (Business Logic) → Repository (Data Access)
     ↑                   ↑                         ↑
   DTO              Model/DTO                   Model + DB
```

## Langkah-langkah

### 1. Definisikan Model

Buat file model baru di `internal/model/<fitur>.go`:

```go
package model

type Fitur struct {
    ID        string     `json:"id"`
    Name      string     `json:"name"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
```

### 2. Definisikan DTO

Buat DTO request/response di `internal/dto/<fitur>.dto.go`:

```go
package dto

type CreateFiturRequest struct {
    Name string `json:"name" validate:"required,min=3,max=255"`
}

type UpdateFiturRequest struct {
    Name string `json:"name" validate:"required,min=3,max=255"`
}
```

### 3. Buat Repository

Buat data access layer di `internal/repository/<fitur>.repository.go`:

```go
package repository

type FiturRepository struct {
    db *database.DB
}

func NewFiturRepository(db *database.DB) *FiturRepository {
    return &FiturRepository{db: db}
}
```

- Gunakan `database.QueryBuilder` dari `pkg/database` untuk membangun query
- Return `(model.Fitur, error)` atau `([]model.Fitur, int, error)` untuk endpoint list

### 4. Buat Usecase

Buat business logic layer di `internal/usecase/<fitur>.usecase.go`:

```go
package usecase

type FiturUsecase struct {
    repo *repository.FiturRepository
    jwt  *config.JWTConfig
}

func NewFiturUsecase(repo *repository.FiturRepository, jwt *config.JWTConfig) *FiturUsecase {
    return &FiturUsecase{repo: repo, jwt: jwt}
}
```

- Validasi business rules di sini
- Gunakan `validator.ValidateStruct()` untuk validasi DTO

### 5. Buat Handler

Buat HTTP handler di `internal/handler/<fitur>.handler.go`:

```go
package handler

type FiturHandler struct {
    usecase *usecase.FiturUsecase
}

func NewFiturHandler(uc *usecase.FiturUsecase) *FiturHandler {
    return &FiturHandler{usecase: uc}
}
```

- Tambahkan anotasi Swagger untuk semua endpoint
- Gunakan `response.Success()`, `response.Error()`, `response.SuccessWithMeta()` dari `pkg/response`
- Bind dan validasi DTO menggunakan `c.ShouldBindJSON()`

### 6. Daftarkan Route

Tambahkan route di `internal/handler/routes.go` dalam fungsi `SetupRoutes()`:

```go
fiturs := api.Group("/fiturs")
fiturs.Use(middleware.JWTAuth(&cfg.JWT))
{
    fiturs.GET("", fiturHandler.GetAll)
    fiturs.POST("", fiturHandler.Create)
    fiturs.GET("/:id", fiturHandler.GetByID)
    fiturs.PUT("/:id", fiturHandler.Update)
    fiturs.DELETE("/:id", fiturHandler.Delete)
}
```

### 7. Hubungkan di main.go

Inisialisasi dan inject dependency di `main.go`:

```go
fiturRepo := repository.NewFiturRepository(db.DB)
fiturUsecase := usecase.NewFiturUsecase(fiturRepo, &cfg.JWT)
fiturHandler := handler.NewFiturHandler(fiturUsecase)
```

Update `handler.SetupRoutes(...)` untuk menyertakan `fiturHandler`.

### 8. Update Schema Prisma (jika tabel baru)

Tambahkan model ke `prisma/schema.prisma` lalu jalankan migrasi:

```bash
npm run migrate:dev -- --name tambah_tabel_fitur
```

### 9. Regenerate Dokumentasi Swagger

```bash
swag init
```

### 10. Tes

```bash
go build ./...
air
```

Verifikasi endpoint melalui Swagger UI di `http://localhost:8080/swagger/index.html`.
