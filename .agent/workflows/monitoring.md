---
description: Start/stop stack monitoring Grafana + Loki + Promtail
---

# Stack Monitoring

Kratify Backend menggunakan **Grafana + Loki + Promtail** untuk agregasi log dan monitoring via Docker Compose.

## Mulai Monitoring

// turbo

1. Jalankan semua service monitoring dalam mode detached:

```bash
docker-compose up -d
```

// turbo 2. Pastikan semua service berjalan:

```bash
docker-compose ps
```

3. Akses dashboard monitoring:
   - **Grafana**: http://localhost:3333 (username: `admin`, password: `admin`)
   - **Loki API**: http://localhost:3100

## Hentikan Monitoring

1. Hentikan dan hapus semua container monitoring:

```bash
docker-compose down
```

## Lihat Log

// turbo

1. Tampilkan log real-time dari semua service monitoring:

```bash
docker-compose logs -f
```

## Setup Lengkap (Pertama Kali)

1. Jalankan script setup yang membuat direktori dan menjalankan service:

```bash
bash setup-monitoring.sh
```
