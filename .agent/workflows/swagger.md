---
description: Generate atau update dokumentasi Swagger/OpenAPI
---

# Dokumentasi Swagger

Kratify Backend menggunakan **swag** untuk auto-generate dokumentasi Swagger dari anotasi Go.

## Prasyarat

- Install swag CLI: `go install github.com/swaggo/swag/cmd/swag@latest`

## Langkah-langkah

// turbo

1. Regenerate dokumentasi Swagger dari anotasi kode:

```bash
swag init
```

Ini akan mengupdate file di direktori `docs/`:

- `docs/docs.go`
- `docs/swagger.json`
- `docs/swagger.yaml`

// turbo 2. Pastikan hasil generate bisa di-compile:

```bash
go build ./...
```

3. Jalankan dev server dan akses Swagger UI:

```
http://localhost:8080/swagger/index.html
```

## Menambahkan Anotasi Swagger

Tambahkan anotasi di atas fungsi handler. Contoh:

```go
// @Summary      Ambil profil user
// @Description  Ambil profil dari user yang sudah login
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /api/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) { ... }
```
