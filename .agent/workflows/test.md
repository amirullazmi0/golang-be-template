---
description: Jalankan tes Go (unit test dan integration test)
---

# Menjalankan Tes

## Jalankan Semua Tes

// turbo

1. Jalankan semua tes di seluruh project dengan output verbose:

```bash
go test ./... -v
```

## Jalankan Tes untuk Package Tertentu

// turbo

1. Jalankan tes hanya untuk package tertentu (ganti `<package>` dengan path-nya):

```bash
go test ./internal/usecase/... -v
```

Contoh:

```bash
go test ./internal/handler/... -v
go test ./internal/repository/... -v
go test ./pkg/database/... -v
```

## Jalankan Tes dengan Coverage

// turbo

1. Jalankan tes dengan laporan coverage:

```bash
go test ./... -v -cover
```

// turbo 2. Generate laporan coverage dalam format HTML:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Jalankan Fungsi Tes Tertentu

// turbo

1. Jalankan fungsi tes tertentu berdasarkan nama:

```bash
go test ./... -v -run NamaFungsiTes
```

## Deteksi Race Condition

// turbo

1. Jalankan tes dengan race detector:

```bash
go test ./... -race -v
```
