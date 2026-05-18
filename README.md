# Go Beresin 🚀

Project Go ini telah diinisialisasi dan siap untuk dikembangkan! Project ini dilengkapi dengan server HTTP bawaan, JSON endpoints, logger middleware, serta **graceful shutdown** untuk manajemen proses yang aman.

## 🛠️ Persyaratan
* **Go**: versi `1.22.2` atau lebih baru.

## 📁 Struktur Project
```text
go-beresin/
├── main.go      # Entry point utama aplikasi dengan HTTP Server
├── go.mod       # File konfigurasi modul Go
└── README.md    # Dokumentasi project (file ini)
```

## 🚀 Cara Menjalankan Project

### 1. Menjalankan Server secara Langsung (Development)
Untuk menjalankan server secara langsung tanpa melakukan build, jalankan perintah berikut di terminal:
```bash
go run main.go
```
Secara default, server akan berjalan di port `8080`. Anda dapat menggantinya menggunakan environment variable `PORT`:
```bash
PORT=9000 go run main.go
```

### 2. Melakukan Build untuk Produksi
Jika ingin mem-compile kode Go menjadi binary executable:
```bash
go build -o go-beresin
```
Kemudian jalankan binary yang dihasilkan:
```bash
./go-beresin
```

## 🌐 Endpoints yang Tersedia

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| **GET** | `/` | Endpoint selamat datang (JSON) |
| **GET** | `/health` | Memeriksa status kesehatan server (Health Check) |

### Contoh Pemanggilan dengan Curl:
```bash
# Get Welcome Message
curl http://localhost:8080/

# Get Health Check
curl http://localhost:8080/health
```

---
Dibuat secara otomatis dengan penuh dedikasi oleh **Antigravity**. Selamat mengoding! 💻🔥
