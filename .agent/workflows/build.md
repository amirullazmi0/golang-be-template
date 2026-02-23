---
description: Build binary production untuk Kratify Backend
---

# Build Binary Production

## Langkah-langkah

// turbo

1. Build binary Go:

```bash
go build -o kratify-backend.exe .
```

// turbo 2. Pastikan binary sudah terbuat:

```bash
ls -la kratify-backend.exe
```

3. (Opsional) Jalankan binary production:

```bash
./kratify-backend.exe
```

> **Catatan:** Pastikan file `.env` sudah dikonfigurasi dengan benar sebelum menjalankan binary.
