---
description: Setup awal project dari nol (install dependencies, konfigurasi env, jalankan migrasi)
---

# Setup Project

Setup lengkap untuk project Kratify Backend dari fresh clone.

## Prasyarat

- Go 1.25.1+ sudah terinstall
- PostgreSQL sudah berjalan
- Node.js & npm sudah terinstall (untuk Prisma CLI)
- Docker & Docker Compose (opsional, untuk monitoring)

## Langkah-langkah

// turbo

1. Install dependensi Go:

```bash
go mod download
```

// turbo 2. Install dependensi Node.js (Prisma CLI):

```bash
npm install
```

3. Buat file `.env` dari contoh:

```bash
cp .env.example .env
```

4. **Edit `.env`** — sesuaikan nilai berikut dengan environment lokal kamu:
   - `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
   - `JWT_SECRET` (ubah untuk production)
   - `SMTP_*` (untuk layanan email, opsional)
   - `IMAGEKIT_*` (untuk upload file, opsional)

// turbo 5. Jalankan migrasi Prisma untuk setup schema database:

```bash
npm run migrate:dev
```

// turbo 6. Install Air untuk hot-reload saat development:

```bash
go install github.com/air-verse/air@latest
```

// turbo 7. Install swag untuk generate dokumentasi Swagger:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

// turbo 8. Pastikan semua source code bisa di-compile:

```bash
go build ./...
```

9. Jalankan development server:

```bash
air
```

Server akan berjalan di `http://localhost:8080`.  
Swagger UI di `http://localhost:8080/swagger/index.html`.
