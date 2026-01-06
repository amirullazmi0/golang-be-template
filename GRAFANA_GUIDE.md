# 📊 Panduan Lengkap Menggunakan Grafana untuk Monitoring

## 🚀 Quick Start

### 1. Start Monitoring Stack

```bash
# Start semua services (Grafana, Loki, Promtail)
docker-compose up -d

# Cek status services
docker-compose ps

# Output yang benar:
# kratify-grafana    running    0.0.0.0:3000->3000/tcp
# kratify-loki       running    0.0.0.0:3100->3100/tcp
# kratify-promtail   running
```

### 2. Start Aplikasi Backend

```bash
# Pastikan .env sudah dikonfigurasi dengan LOG_TO_FILE=true
go run main.go

# Atau dengan Air (hot reload)
air
```

### 3. Akses Grafana

1. Buka browser: **http://localhost:3000**
2. Login dengan:
     - **Username**: `admin`
     - **Password**: `admin`
3. Grafana akan minta ganti password (bisa skip atau ganti sesuai keinginan)

---

## 🔍 Cara Melihat Logs di Grafana

### Metode 1: Explore (Paling Mudah untuk Debugging)

1. **Klik icon kompas (🧭)** di sidebar kiri
2. Atau buka langsung: http://localhost:3000/explore
3. Pastikan **data source "Loki"** sudah terpilih di dropdown atas
4. Di **Label filters**, klik dan pilih:
     - `job` = `kratify-backend`
5. Klik **"Run query"** (tombol biru) atau tekan `Shift + Enter`

### Metode 2: Menggunakan LogQL Queries

Di halaman Explore, masukkan query di kolom query builder:

#### Query Dasar

**Lihat semua logs:**

```logql
{job="kratify-backend"}
```

**Lihat logs terakhir 5 menit:**

```logql
{job="kratify-backend"} [5m]
```

**Lihat logs dengan level tertentu:**

```logql
{job="kratify-backend", level="error"}
{job="kratify-backend", level="info"}
{job="kratify-backend", level="warn"}
```

#### Filter berdasarkan Content

**Search text di message:**

```logql
{job="kratify-backend"} |= "HTTP Request"
{job="kratify-backend"} |= "error"
{job="kratify-backend"} != "health"  # exclude health checks
```

**Multiple filters:**

```logql
{job="kratify-backend"} |= "HTTP Request" |= "POST"
```

#### Parse JSON dan Filter Field

**Filter by endpoint:**

```logql
{job="kratify-backend"} | json | path="/api/auth/login"
{job="kratify-backend"} | json | path=~"/api/users/.*"  # regex
```

**Filter by HTTP method:**

```logql
{job="kratify-backend"} | json | method="POST"
{job="kratify-backend"} | json | method=~"POST|PUT|DELETE"
```

**Filter by status code:**

```logql
{job="kratify-backend"} | json | status >= 400
{job="kratify-backend"} | json | status >= 500  # server errors only
{job="kratify-backend"} | json | status < 400   # success only
```

**Filter by latency (slow requests):**

```logql
{job="kratify-backend"} | json | latency_ms > 500
{job="kratify-backend"} | json | latency_ms > 1000  # lebih dari 1 detik
```

**Filter by IP:**

```logql
{job="kratify-backend"} | json | ip="127.0.0.1"
```

**Filter by request ID:**

```logql
{job="kratify-backend"} | json | request_id="abc123"
```

#### Kombinasi Complex Queries

**Semua error requests:**

```logql
{job="kratify-backend"} | json | status >= 400 | line_format "{{.timestamp}} [{{.level}}] {{.method}} {{.path}} - {{.status}} ({{.latency_ms}}ms)"
```

**Slow POST requests:**

```logql
{job="kratify-backend"} | json | method="POST" | latency_ms > 500
```

**Failed login attempts:**

```logql
{job="kratify-backend"} | json | path="/api/auth/login" | status >= 400
```

---

## 📈 Membuat Dashboard

### Step 1: Create New Dashboard

1. Klik icon **"+"** di sidebar kiri
2. Pilih **"Dashboard"**
3. Klik **"Add visualization"**
4. Pilih **data source "Loki"**

### Step 2: Add Panels

#### Panel 1: Request Rate (Requests per Second)

**Query:**

```logql
sum(rate({job="kratify-backend"} | json [1m]))
```

**Settings:**

-    Visualization: **Time series** atau **Stat**
-    Title: "Request Rate"
-    Unit: "ops" (operations per second)

#### Panel 2: Error Rate

**Query:**

```logql
sum(rate({job="kratify-backend", level="error"} [5m]))
```

**Settings:**

-    Visualization: **Time series**
-    Title: "Error Rate"
-    Color: Red
-    Alert: Set threshold jika error rate > 10

#### Panel 3: Average Latency

**Query:**

```logql
avg(avg_over_time({job="kratify-backend"} | json | unwrap latency_ms [5m]))
```

**Settings:**

-    Visualization: **Gauge** atau **Stat**
-    Title: "Average Latency"
-    Unit: "ms" (milliseconds)
-    Thresholds:
     -    Green: < 100ms
     -    Yellow: 100-500ms
     -    Red: > 500ms

#### Panel 4: Status Code Distribution

**Query:**

```logql
sum by (status) (count_over_time({job="kratify-backend"} | json [5m]))
```

**Settings:**

-    Visualization: **Pie chart** atau **Bar chart**
-    Title: "HTTP Status Codes"
-    Legend: Show

#### Panel 5: Top 10 Endpoints

**Query:**

```logql
topk(10, sum by (path) (count_over_time({job="kratify-backend"} | json [5m])))
```

**Settings:**

-    Visualization: **Bar chart**
-    Title: "Top 10 Endpoints"

#### Panel 6: Request by Method

**Query:**

```logql
sum by (method) (count_over_time({job="kratify-backend"} | json [5m]))
```

**Settings:**

-    Visualization: **Pie chart**
-    Title: "Requests by HTTP Method"

#### Panel 7: Error Logs Table

**Query:**

```logql
{job="kratify-backend", level="error"}
```

**Settings:**

-    Visualization: **Logs** (table view)
-    Title: "Recent Errors"
-    Show: Last 50 entries

#### Panel 8: P95 Latency

**Query:**

```logql
quantile_over_time(0.95, {job="kratify-backend"} | json | unwrap latency_ms [5m])
```

**Settings:**

-    Visualization: **Time series**
-    Title: "P95 Latency"
-    Unit: "ms"

### Step 3: Arrange & Save Dashboard

1. Drag panels untuk arrange layout
2. Resize panels sesuai kebutuhan
3. Klik **"Save dashboard"** (icon disket di atas)
4. Beri nama: "Kratify Backend Monitoring"
5. Klik **"Save"**

---

## ⚡ Features Grafana yang Berguna

### 1. Time Range Picker

Di pojok kanan atas, bisa pilih time range:

-    Last 5 minutes
-    Last 15 minutes
-    Last 1 hour
-    Last 24 hours
-    Last 7 days
-    Custom range

### 2. Auto Refresh

-    Klik dropdown refresh di pojok kanan atas
-    Pilih interval: 5s, 10s, 30s, 1m, 5m
-    Dashboard akan auto-refresh sesuai interval

### 3. Live Tail (Real-time Logs)

Di halaman Explore:

1. Toggle switch **"Live"** di pojok kanan atas
2. Logs akan streaming real-time
3. Sangat berguna untuk debugging

### 4. Log Details

-    Klik pada log entry untuk expand
-    Lihat full JSON dengan semua fields
-    Copy log, share permalink, atau filter by field

### 5. Export Data

-    Klik **"Inspector"** di panel
-    Tab **"Data"**
-    Download as CSV atau JSON

### 6. Annotations

-    Mark important events di timeline
-    Tambah notes untuk deployment, incident, dll

---

## 🎯 Use Cases Praktis

### 1. Debugging Error 500

```logql
{job="kratify-backend"} | json | status >= 500
```

Klik log entry → lihat detail error → trace dengan `request_id`

### 2. Monitor Performance Endpoint Specific

```logql
{job="kratify-backend"} | json | path="/api/users/profile" | latency_ms > 100
```

### 3. Track User Activity (by IP)

```logql
{job="kratify-backend"} | json | ip="192.168.1.100"
```

### 4. Monitor Login Activity

```logql
{job="kratify-backend"} | json | path=~"/api/auth/(login|register)"
```

### 5. Find All Failed Requests

```logql
{job="kratify-backend"} | json | status >= 400 | line_format "{{.timestamp}} {{.method}} {{.path}} - {{.status}}"
```

### 6. Trace Request by ID

```logql
{job="kratify-backend"} | json | request_id="your-request-id-here"
```

---

## 🔧 Tips & Tricks

### 1. Keyboard Shortcuts

-    `Ctrl/Cmd + K`: Command palette
-    `Shift + Enter`: Run query
-    `Ctrl/Cmd + S`: Save dashboard
-    `G + E`: Go to Explore
-    `G + D`: Go to Dashboards

### 2. Query Optimization

**❌ Slow:**

```logql
{job="kratify-backend"} | json | status >= 400
```

**✅ Faster:**

```logql
{job="kratify-backend", level="error"}
```

Gunakan labels (`level="error"`) instead of parsing JSON untuk filter basic.

### 3. Share Dashboard

1. Klik icon **"Share"** di dashboard
2. Pilih **"Link"** atau **"Snapshot"**
3. Copy URL untuk share ke team

### 4. Create Alerts

1. Edit panel di dashboard
2. Tab **"Alert"**
3. Buat rule, contoh: "Alert if error rate > 10/minute"
4. Set notification channel (email, Slack, dll)

### 5. Dark/Light Theme

1. Klik profile icon (pojok kiri bawah)
2. **Preferences**
3. Pilih **UI Theme**

---

## 🧪 Testing Monitoring Setup

### 1. Generate Test Logs

**Success requests:**

```bash
curl http://localhost:8080/api/health
```

**Error requests (404):**

```bash
curl http://localhost:8080/api/nonexistent
```

**Unauthorized (401):**

```bash
curl http://localhost:8080/api/users/profile
```

**Load test (banyak requests):**

```bash
for i in {1..100}; do curl http://localhost:8080/api/health; done
```

### 2. Verify di Grafana

1. Buka Explore
2. Query: `{job="kratify-backend"}`
3. Pastikan logs muncul
4. Check latency, status codes, dll

---

## 📚 Sample Dashboard JSON

Buat file `grafana-dashboard.json`:

```json
{
     "dashboard": {
          "title": "Kratify Backend Monitoring",
          "panels": [
               {
                    "title": "Request Rate",
                    "targets": [
                         {
                              "expr": "sum(rate({job=\"kratify-backend\"} | json [1m]))"
                         }
                    ]
               }
          ]
     }
}
```

Import:

1. **Dashboards** → **Import**
2. Upload JSON file
3. Klik **"Load"**

---

## 🐛 Troubleshooting

### Logs Tidak Muncul

**1. Cek Loki status:**

```bash
curl http://localhost:3100/ready
# Output: ready
```

**2. Cek Promtail logs:**

```bash
docker logs kratify-promtail
```

**3. Cek aplikasi sudah generate logs:**

```bash
ls -la logs/
cat logs/app.log
```

**4. Restart services:**

```bash
docker-compose restart
```

### Grafana Tidak Bisa Connect ke Loki

**1. Cek data source settings:**

-    Settings → Data sources → Loki
-    URL harus: `http://loki:3100`
-    Click **"Save & test"**

**2. Cek network:**

```bash
docker network ls
docker network inspect kratify-backend_monitoring
```

### Query Error

-    Pastikan syntax LogQL benar
-    Check typo di label names
-    Verify time range (jangan terlalu lama di development)

---

## 📖 Resources

-    **LogQL Documentation**: https://grafana.com/docs/loki/latest/logql/
-    **Grafana Docs**: https://grafana.com/docs/grafana/latest/
-    **Query Examples**: https://grafana.com/docs/loki/latest/logql/query_examples/

---

## 🎉 Next Steps

1. ✅ Setup monitoring stack
2. ✅ View logs di Explore
3. ✅ Create dashboard dengan panels
4. ✅ Setup alerts untuk error rates
5. ✅ Share dashboard dengan team
6. 🚀 Integrate dengan Slack/PagerDuty untuk notifications

Selamat monitoring! 🔥
