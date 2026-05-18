# Doppler Secret Management Guide

Dokumen ini berisi panduan lengkap tentang cara mengatur, menyinkronkan, dan menggunakan **Doppler CLI** sebagai pengelola variabel lingkungan (secret management) terpusat pada project **Go Beresin**.

---

## 1. Instalasi Doppler CLI
Sebelum memulai, Anda perlu menginstal Doppler CLI di komputer lokal Anda:

### Linux (Ubuntu/Debian)
```bash
# Tambahkan repo gpg key
curl -sLf --retry 3 https://packages.doppler.com/public/cli/gpg.DEB-GPG-KEY-doppler | sudo gpg --dearmor -o /usr/share/keyrings/doppler-archive-keyring.gpg

# Tambahkan repo apt
echo "deb [signed-by=/usr/share/keyrings/doppler-archive-keyring.gpg] https://packages.doppler.com/public/cli/deb/debian any-version main" | sudo tee /etc/apt/sources.list.d/doppler-cli.list

# Install Doppler
sudo apt update && sudo apt install doppler
```

### macOS (Homebrew)
```bash
brew install dopplerhq/cli/doppler
```

### Verifikasi Instalasi
```bash
doppler --version
```

---

## 2. Autentikasi & Setup Project
Setelah Doppler CLI terinstal, ikuti langkah-langkah berikut:

1. **Login ke Akun Doppler**:
   Jalankan perintah berikut untuk mengautentikasi CLI dengan akun Doppler Anda:
   ```bash
   doppler login
   ```
   *Ikuti instruksi di terminal untuk membuka browser dan memasukkan auth code.*

2. **Inisialisasi Project di Workspace**:
   Project ini sudah dilengkapi dengan file `doppler.yaml` di root direktori yang mendefinisikan pemetaan project ke `beresin` dan config ke `dev`. Untuk menghubungkan workspace lokal Anda, cukup jalankan:
   ```bash
   doppler setup
   ```
   *Pilih opsi default (Project: `beresin`, Config: `dev`) saat diminta.*

---

## 3. Unggah (Sync) Secret Awal dari `.env`
Jika Anda baru pertama kali membuat project `beresin` di dashboard Doppler, Anda dapat mengunggah seluruh isi berkas `.env` lokal Anda saat ini ke Doppler dev environment dengan perintah:

```bash
doppler secrets upload .env
```

Ini secara otomatis akan mengimpor semua variabel seperti `DATABASE_URL`, `REDIS_ADDR`, `PORT`, `REDIS_PASSWORD`, `JWT_SECRET`, dan lain-lain ke dashboard cloud Doppler Anda. Setelah terunggah, Anda dapat menghapus berkas `.env` lokal Anda demi keamanan.

---

## 4. Cara Menjalankan Aplikasi
Gunakan perintah `doppler run --` di depan perintah startup Go Anda agar Doppler menyuntikkan semua variabel lingkungan secara dinamis pada saat runtime:

### Menjalankan API Server
```bash
doppler run -- go run cmd/api/main.go
```

### Menjalankan Database Migrasi
```bash
doppler run -- go run cmd/migrate/main.go up
```

### Menjalankan Database Seeder
```bash
doppler run -- go run cmd/seed/main.go
```

Dengan cara ini, Anda tidak memerlukan lagi file `.env` fisik pada komputer lokal Anda, menjadikannya sangat aman dan konsisten di seluruh tim pengembang!
