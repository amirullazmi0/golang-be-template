---
description: Commit perubahan kode ke Git dengan conventional commit message
---

// turbo-all

# Git Commit

Workflow untuk commit perubahan kode ke repository Git menggunakan format **Conventional Commits**.

## Format Commit Message

```
<tipe>(<scope>): <deskripsi singkat>
```

### Tipe yang Tersedia

| Tipe       | Keterangan                                            |
| ---------- | ----------------------------------------------------- |
| `feat`     | Fitur baru                                            |
| `fix`      | Perbaikan bug                                         |
| `docs`     | Perubahan dokumentasi                                 |
| `style`    | Formatting, titik koma, dll (bukan perubahan kode)    |
| `refactor` | Refactoring kode tanpa tambah fitur atau perbaiki bug |
| `test`     | Menambah atau memperbaiki tes                         |
| `chore`    | Perubahan build, tools, dependensi, dll               |
| `perf`     | Peningkatan performa                                  |
| `ci`       | Perubahan CI/CD                                       |

### Contoh Scope (Opsional)

`auth`, `user`, `product`, `address`, `cart`, `order`, `db`, `middleware`, `config`, `docs`

## Langkah-langkah

1. Cek status perubahan:

```bash
git status --short
```

2. Cek diff perubahan:

```bash
git diff --stat
```

3. Tentukan mode commit:
   - **Semua sekaligus**: lanjut ke langkah 4a
   - **Per file**: lanjut ke langkah 4b

### 4a. Commit Semua Sekaligus

```bash
git add .
git commit -m "<tipe>(<scope>): <deskripsi>"
```

### 4b. Commit Per File

Untuk setiap file yang berubah, lakukan:

```bash
git add <path/ke/file>
git commit -m "<tipe>(<scope>): <deskripsi untuk file ini>"
```

Contoh commit per file:

```bash
git add internal/handler/user.handler.go
git commit -m "feat(user): tambah endpoint update role"

git add internal/repository/user.repository.go
git commit -m "feat(user): tambah query update role di repository"
```

> **Tips:** Bilang "commit per file" supaya setiap file di-commit terpisah dengan message masing-masing.

### Contoh Commit Message

```bash
# Fitur baru
git commit -m "feat(product): tambah endpoint CRUD produk"

# Perbaikan bug
git commit -m "fix(auth): perbaiki validasi token expired"

# Dokumentasi
git commit -m "docs(swagger): update anotasi endpoint user"

# Refactoring
git commit -m "refactor(repository): pindah query builder ke common helper"

# Migrasi database
git commit -m "chore(db): tambah migrasi tabel orders"
```

5. (Opsional) Push ke remote:

```bash
git push origin <nama-branch>
```
