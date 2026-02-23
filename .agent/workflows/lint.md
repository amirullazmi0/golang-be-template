---
description: Lint dan format kode Go menggunakan go vet dan gofmt
---

# Lint & Format

## Format Kode

// turbo

1. Format semua file Go di project:

```bash
gofmt -w .
```

## Vet (Analisis Statis)

// turbo

1. Jalankan `go vet` untuk mendeteksi kesalahan umum:

```bash
go vet ./...
```

## Rapihkan Dependensi

// turbo

1. Bersihkan `go.mod` dan `go.sum`:

```bash
go mod tidy
```

## Cek Semua Sekaligus

// turbo

1. Jalankan format, vet, dan build secara berurutan:

```bash
gofmt -w . && go vet ./... && go build ./...
```
