---
description: Kelola migrasi database Prisma (buat, terapkan, reset, push)
---

# Workflow Migrasi Database

Kratify Backend menggunakan **Prisma** sebagai ORM untuk migrasi PostgreSQL.  
File schema: `prisma/schema.prisma`

## Buat Migrasi Baru (Development)

1. Edit schema Prisma di `prisma/schema.prisma` sesuai perubahan model.

// turbo 2. Jalankan migrasi dengan nama deskriptif:

```bash
npm run migrate:dev -- --name <nama_migrasi>
```

Ganti `<nama_migrasi>` dengan nama deskriptif seperti `tambah_kolom_status_order`.

## Terapkan Migrasi (Production/Deploy)

// turbo

1. Deploy migrasi yang tertunda ke database:

```bash
npm run migrate:deploy
```

## Reset Database (⚠️ Destruktif)

> **PERINGATAN:** Ini akan menghapus SEMUA data di database.

1. Reset database dan terapkan ulang semua migrasi:

```bash
npm run migrate:reset
```

## Push Schema Tanpa Migrasi

// turbo

1. Push perubahan schema secara langsung (cocok untuk prototyping, melewati histori migrasi):

```bash
npm run db:push
```

## Pull Schema dari Database

// turbo

1. Introspeksi database yang sudah ada dan update schema Prisma:

```bash
npm run db:pull
```

## Buka Prisma Studio (GUI)

1. Jalankan Prisma Studio untuk melihat/edit data secara visual:

```bash
npm run studio
```

Prisma Studio akan terbuka di browser.
