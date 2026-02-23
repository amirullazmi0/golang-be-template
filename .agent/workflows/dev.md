---
description: Jalankan development server dengan Air hot-reload
---

# Development Server

Menjalankan Kratify Backend development server menggunakan Air untuk hot-reload.

## Prasyarat

- Go 1.25.1+ sudah terinstall
- Air sudah terinstall (`go install github.com/air-verse/air@latest`)
- PostgreSQL berjalan dengan konfigurasi `.env` yang benar
- Migrasi Prisma sudah diterapkan (`npm run migrate:dev`)

## Langkah-langkah

// turbo

1. Jalankan development server dengan Air:

```bash
air
```

Server akan berjalan di `http://localhost:8080`

### Tanpa Hot-Reload

Jika Air belum tersedia, jalankan secara langsung:

```bash
go run main.go
```

### Verifikasi

- Health check: `GET http://localhost:8080/health`
- Swagger UI: `http://localhost:8080/swagger/index.html`
