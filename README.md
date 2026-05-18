# Go Beresin Backend 🚀

Selamat datang di repositori backend **Go Beresin**! Aplikasi ini dirancang menggunakan bahasa pemrograman Go dengan performa tinggi menggunakan framework **Fiber v2**, didukung oleh **PostgreSQL + PostGIS** untuk pencarian spasial geospasial, **Redis** untuk rate-limiting & caching, serta integrasi **Doppler** untuk manajemen secret terpusat.

---

## 🛠️ Persyaratan Sistem (Prerequisites)
Sebelum menjalankan project secara lokal, pastikan Anda telah menginstal perkakas berikut:
* **Go**: Versi `1.25.0` atau yang lebih baru.
* **Docker & Docker Compose**: Untuk menjalankan service pendukung secara instan.
* **Doppler CLI** (Opsional tapi Direkomendasikan): Untuk secret management.

---

## 📁 Struktur Repositori Utama
```text
go-beresin/
├── cmd/
│   ├── api/          # Entry point utama web server Fiber
│   ├── migrate/      # CLI tool migrasi database SQL (Up/Down)
│   └── seed/         # Script seeder otomatis untuk data awal realistik
├── docs/
│   ├── doppler.md    # Panduan lengkap setup Doppler Secret Management
│   ├── go-beresin.postman_collection.json # Kontrak API Postman tim frontend/backend
│   └── swagger.*     # Auto-generated specs Swagger UI (JSON/YAML)
├── internal/
│   ├── repository/   # Repository layer & model hasil generate sqlc
│   ├── transport/    # Router, Handlers (endpoints), dan Middlewares
│   └── ...
├── migrations/       # DDL SQL skema database (termasuk PostGIS)
├── pkg/
│   └── database/     # DB Wrapper untuk inisialisasi pgxpool.Pool
├── docker-compose.yml# Konfigurasi container lokal dev (MinIO, ES, Kibana, dll)
├── sqlc.yaml         # Konfigurasi generator sqlc untuk pgx/v5
└── .env.example      # Template variabel lingkungan lokal
```

---

## 🚀 Panduan Setup & Eksekusi Cepat

### 1. Jalankan Environment Pendukung (Docker Compose)
Aplikasi ini menggunakan beberapa service pendukung untuk development lokal. PostgreSQL dan Redis diasumsikan sudah terinstal dan berjalan secara lokal/terpisah di komputer Anda. Jalankan service tambahan lainnya dengan Docker Compose:
```bash
docker-compose up -d
```
Service tambahan yang diaktifkan meliputi:
* **MinIO** (Port `9000` / `9001`): Penyimpanan berkas S3-Compatible lokal.
* **Elasticsearch** (Port `9200`): Service pencarian teks penuh.
* **Kibana** (Port `5601`): Panel visualisasi data Elasticsearch.

### 2. Konfigurasi Environment & Secret
Project ini mendukung **Doppler CLI** untuk manajemen secret terpusat tanpa perlu file `.env` statis secara lokal. 
* Silakan ikuti panduan setup lengkap di dokumen: **[Panduan Doppler](file:///home/noto/projects/webdev/go-beresin/docs/doppler.md)**.
* Sebagai alternatif tradisional, Anda dapat menyalin `.env.example` menjadi `.env` dan menyesuaikan nilainya secara manual.

### 3. Jalankan Migrasi Skema Database
Untuk memigrasikan skema 16 tabel lengkap beserta PostGIS geospasial dan index performa tinggi:
```bash
# Menggunakan Doppler:
doppler run -- go run cmd/migrate/main.go up

# Menggunakan file .env lokal:
go run cmd/migrate/main.go up
```

### 4. Jalankan Database Seeder (Data Awal Uji)
Untuk menyuntikkan data uji (kategori, sub-kategori, user dummy, talent geospasial nyata Jabodetabek/Bandung, service listing, dan package harga):
```bash
# Menggunakan Doppler:
doppler run -- go run cmd/seed/main.go

# Menggunakan file .env lokal:
go run cmd/seed/main.go
```

### 5. Jalankan API Server
Untuk mengaktifkan API server Fiber di port `8080`:
```bash
# Menggunakan Doppler:
doppler run -- go run cmd/api/main.go

# Menggunakan file .env lokal:
go run cmd/api/main.go
```

---

## 🌐 Kontrak API & Dokumentasi

### Swagger UI (Dokumentasi Endpoint Interaktif)
Dokumentasi endpoint interaktif berbasis spesifikasi OpenAPI 2.0 (Swagger) dapat langsung diakses saat server berjalan pada alamat:
👉 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

Jika Anda melakukan perubahan anotasi pada handler, Anda dapat me-regenerate file spesifikasi Swagger dengan perintah:
```bash
# Pastikan swag CLI terinstal (go install github.com/swaggo/swag/cmd/swag@latest)
$HOME/go/bin/swag init -g cmd/api/main.go
```

### Postman Collection (Kontrak Integrasi Flutter & Backend)
Kami telah menyediakan API Collection lengkap yang siap diimpor ke Postman atau Bruno untuk mempermudah integrasi dengan tim Mobile Frontend (Flutter):
👉 **[Go Beresin Postman Collection](file:///home/noto/projects/webdev/go-beresin/docs/go-beresin.postman_collection.json)**

Collection ini sudah dilengkapi dengan test scripts otomatis untuk menyisipkan Bearer JWT Token ke dalam variabel secara otomatis setelah proses mock login berhasil.

---
*Dikembangkan dengan penuh dedikasi oleh tim **Go Beresin**.* 💻🔥
