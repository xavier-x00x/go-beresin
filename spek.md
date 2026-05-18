# BERESIN — Product Requirements Document (PRD)
## Marketplace Jasa Sektor Real | Versi 1.0

---

> **Dokumen ini adalah panduan teknis dan produk lengkap untuk membangun aplikasi Beresin dari nol hingga MVP, kemudian ke versi full-scale.**

---

## DAFTAR ISI

1. [Ringkasan Eksekutif](#1-ringkasan-eksekutif)
2. [Konsep Bisnis & Model Pendapatan](#2-konsep-bisnis--model-pendapatan)
3. [Target Pengguna & Persona](#3-target-pengguna--persona)
4. [Arsitektur Sistem](#4-arsitektur-sistem)
5. [Tech Stack Rekomendasi](#5-tech-stack-rekomendasi)
6. [Fitur & Modul Lengkap](#6-fitur--modul-lengkap)
7. [Database Schema](#7-database-schema)
8. [API Endpoints](#8-api-endpoints)
9. [Alur Pengguna (User Flow)](#9-alur-pengguna-user-flow)
10. [Desain & UI/UX Spec](#10-desain--uiux-spec)
11. [Sistem Keamanan & Escrow](#11-sistem-keamanan--escrow)
12. [Notifikasi & Komunikasi](#12-notifikasi--komunikasi)
13. [Admin Panel](#13-admin-panel)
14. [Roadmap Pengembangan](#14-roadmap-pengembangan)
15. [Estimasi Biaya & Tim](#15-estimasi-biaya--tim)
16. [Risiko & Mitigasi](#16-risiko--mitigasi)

---

## 1. RINGKASAN EKSEKUTIF

### Apa itu Beresin?

Beresin adalah platform marketplace **dua arah** untuk jasa sektor real/lapangan di Indonesia. Platform ini menghubungkan **Talent/Freelancer** (penyedia jasa) dengan **User** (pengguna jasa) melalui dua mekanisme utama:

- **Mode A — Direct Booking:** User mencari dan langsung memesan jasa Talent yang telah terdaftar.
- **Mode B — Tender/Bidding:** User memposting kebutuhan pekerjaan, lalu Talent mengajukan penawaran.

### Differensiasi dari Kompetitor

| Platform | Fokus | Kekurangan |
|---|---|---|
| Fastwork / Sribulancer | Jasa digital (desain, konten) | Tidak cover jasa fisik/lapangan |
| GoService (Gojek) | Jasa rumah tangga terbatas | Tidak ada negosiasi/kontrak |
| Upwork / Fiverr | Jasa remote internasional | Bukan untuk talent lokal lapangan |
| **Beresin** | **Jasa fisik + digital lokal** | **Solusi lengkap dari negosiasi → kontrak → escrow** |

### Nilai Utama Platform

- **Trust:** Verifikasi Talent, rating, review, dan kontrak digital
- **Fleksibilitas:** Booking langsung ATAU posting tender
- **Keamanan:** Sistem Escrow — uang ditahan platform hingga pekerjaan selesai
- **Negosiasi:** Chat real-time dengan fitur penawaran terstruktur

---

## 2. KONSEP BISNIS & MODEL PENDAPATAN

### Kategori Jasa yang Dilayani

**Rumah & Properti:**
- Tukang bangunan (cat, plester, renovasi)
- Tukang listrik
- Tukang ledeng / plumber
- Servis AC / elektronik
- Cleaning service
- Tukang kebun
- Pest control

**Event & Entertainment:**
- Grup Band / Musisi
- DJ
- MC (Master of Ceremony)
- Dancer / Penari
- Fotografer & Videografer
- Dekorator event
- Catering

**Kesehatan & Konseling:**
- Psikolog
- Psikiater
- Terapis
- Personal trainer
- Nutrisionis

**Pendidikan:**
- Guru privat (semua mata pelajaran)
- Instruktur musik
- Pelatih olahraga

**Kecantikan & Perawatan:**
- Make-up artist
- Stylist rambut
- Nail artist

**Teknisi:**
- Servis komputer / laptop
- Teknisi HP
- Instalasi CCTV / smart home

### Model Pendapatan

#### 1. Service Fee (Komisi Transaksi)
- **User:** 5% dari nilai transaksi
- **Talent:** 10–15% dari nilai transaksi
- Dihitung saat release escrow

#### 2. Subscription Talent (Premium)
| Tier | Harga/bulan | Benefit |
|---|---|---|
| Gratis | Rp 0 | Maks 5 penawaran/bulan, listing biasa |
| Pro | Rp 99.000 | Penawaran unlimited, badge Pro, prioritas pencarian |
| Business | Rp 299.000 | Semua Pro + analitik lanjutan, multi-akun tim |

#### 3. Fitur Berbayar
- **Boost Listing:** Rp 25.000–75.000 untuk tampil di posisi teratas pencarian (3/7/14 hari)
- **Featured Talent:** Rp 150.000/minggu untuk tampil di homepage
- **Kontrak Premium:** Template kontrak resmi bernotaris (Rp 50.000/kontrak)

#### 4. Asuransi Pekerjaan (Opsional)
- User dapat membeli proteksi pekerjaan: Rp 10.000–50.000 per transaksi
- Menanggung kerugian hingga nilai tertentu jika terjadi sengketa

---

## 3. TARGET PENGGUNA & PERSONA

### Persona User (Pencari Jasa)

**Persona 1 — Budi, 35 tahun, karyawan swasta Jakarta**
- Butuh tukang renovasi, tidak tahu harus cari di mana
- Khawatir ditipu atau hasil kerja buruk
- Ingin tahu harga sebelum deal
- Pain: Susah cari tukang terpercaya, tidak ada jaminan kualitas

**Persona 2 — Sinta, 28 tahun, event organizer freelance**
- Sering butuh band, DJ, fotografer untuk acara klien
- Butuh proses booking yang cepat dan ada kontrak
- Pain: Koordinasi talent acara repot, tidak ada rekam jejak

**Persona 3 — Pak Hendra, 52 tahun, pengusaha properti**
- Butuh tim cleaning service rutin, teknisi berkala
- Butuh invoice dan dokumentasi untuk pembukuan
- Pain: Tidak ada sistem yang profesional untuk jasa lapangan

### Persona Talent (Penyedia Jasa)

**Persona 1 — Agus, 30 tahun, tukang listrik berpengalaman**
- Ingin dapat klien lebih banyak tanpa bergantung mulut ke mulut
- Ingin pembayaran aman (tidak kabur setelah dikerjakan)
- Pain: Sering tidak dibayar lunas setelah pekerjaan selesai

**Persona 2 — Rizky, 25 tahun, fotografer wedding**
- Punya portofolio bagus tapi susah promosi
- Ingin sistem booking yang rapi dan otomatis
- Pain: Deal lewat WhatsApp tidak profesional, tidak ada kontrak

**Persona 3 — Grup Band "The Acoustic", 4 personel**
- Ingin satu platform untuk terima booking manggung
- Butuh bisa setting paket harga berbeda (café, wedding, corporate)
- Pain: Sering klien batalkan tanpa kompensasi

---

## 4. ARSITEKTUR SISTEM

### Overview Arsitektur

```
[Flutter App (iOS + Android)]
           ↕
     [API Gateway / Load Balancer]
           ↕
      [Go Backend Services]
   ┌─────────────────────────────────────┐
   │  Auth    │  User    │  Job/Tender   │
   │  Service │  Service │  Service      │
   ├──────────┼──────────┼───────────────┤
   │  Chat    │  Payment │  Notification │
   │  Service │  Service │  Service      │
   ├──────────┼──────────┼───────────────┤
   │  Review  │  Search  │  Admin        │
   │  Service │  Service │  Service      │
   └─────────────────────────────────────┘
           ↕
   [Database Layer]
   PostgreSQL (primary) + Redis (cache/pubsub) + Elasticsearch (search)
           ↕
   [Storage]
   Cloudflare R2 (foto, video, dokumen)
           ↕
   [Third Party Services]
   Midtrans (payment) │ Firebase (push notif) │ Twilio (WhatsApp/SMS)
   Google Maps API │ SendGrid (email) │ Privy ID (e-sign)
```

### Microservices Detail

**1. Auth Service**
- JWT-based authentication
- OAuth2 (Google, Apple)
- WhatsApp OTP via Twilio
- Refresh token rotation
- Role-based access control (RBAC): user / talent / admin
- API contract versioning untuk stabilitas frontend Flutter

**2. User Service**
- Profile management (user & talent)
- KYC / verifikasi identitas
- Portfolio management
- Availability calendar

**3. Job/Tender Service**
- Direct service listing (Mode A)
- Tender posting & bidding (Mode B)
- Matching engine (talent ↔ tender)
- Booking management

**4. Chat Service**
- WebSocket real-time messaging
- File sharing (gambar, dokumen, audio)
- Quotation sharing dalam chat
- Chat history archiving
- Redis pub/sub untuk scaling horizontal

**5. Payment Service**
- Escrow management
- Integrasi Midtrans (VA, e-wallet, QRIS)
- DP / pelunasan flow
- Disbursement ke talent (via Flip.id atau Midtrans payout)
- Invoice generation (PDF otomatis)
- Idempotency key untuk semua callback pembayaran

**6. Notification Service**
- Push notification (Firebase FCM)
- Email (SendGrid)
- WhatsApp notification (via Twilio / WABA)
- In-app notification center
- Semua pengiriman dilakukan async via queue

**7. Review Service**
- Rating system (1–5 bintang)
- Multi-dimensi: ketepatan waktu, kualitas, komunikasi, keramahan
- Foto hasil kerja
- Moderasi review (AI-assisted)

**8. Search Service**
- Elasticsearch untuk pencarian full-text
- Geo-search (cari talent terdekat)
- Faceted filtering (kategori, harga, rating, ketersediaan)
- Ranking algorithm (rating + responsivitas + kelengkapan profil)

**9. Admin Service**
- Dashboard monitoring
- Verifikasi talent
- Dispute resolution
- Escrow control
- Analytics & reporting

---

## 5. TECH STACK REKOMENDASI

### Frontend — Mobile

**Framework:** Flutter 3.x
- Satu codebase untuk iOS dan Android
- Routing, state, dan networking dipisah jelas agar mudah dirawat
- Fokus pada performance, offline-friendly UI, dan integrasi API yang stabil

**Library Utama:**
```
go_router                  → navigasi
dio                        → data fetching & caching
riverpod                   → state management
google_maps_flutter        → peta & geolokasi
web_socket_channel         → real-time chat
image_picker / file_picker  → upload foto & file
flutter_secure_storage      → token storage aman
table_calendar             → kalender booking/ketersediaan
flutter_form_builder       → form management
```

### Frontend — Ops Dashboard / Internal Web

**Framework:** Next.js 14 (App Router) untuk ops dashboard terpisah, atau Flutter Web jika phase berikutnya membutuhkan satu codebase tambahan
- Dipisah dari aplikasi mobile
- Mengonsumsi OpenAPI yang sama dari backend Go
- Dipakai untuk admin, operasional, dan analitik internal

**Library Utama:**
```
tailwindcss                → styling
shadcn/ui                  → komponen UI
tanstack-query             → data fetching
zod                        → validasi schema
mapbox / leaflet           → peta
framer-motion              → animasi
```

### Backend

**Runtime:** Go 1.23

**Rekomendasi Utama:** Go service terpisah dengan clean architecture / layered architecture

```
HTTP Layer:      chi / gin / fiber
Database:        PostgreSQL + pgx + sqlc
Message Queue:   Redis + Asynq (untuk async jobs)
Real-time:       WebSocket + Redis pub/sub
Validation:      go-playground/validator
Auth:            JWT + OAuth2
File Upload:     multipart + S3 SDK + image processing
Email:           SendGrid SDK
PDF:             chromedp / wkhtmltopdf
```

### Database

```
PostgreSQL 15+     → primary database (relational data)
Redis 7+           → session, cache, queue, real-time pub/sub
Elasticsearch 8+   → full-text search, geo-search
```

### Infrastructure

```
Cloud:            AWS / Google Cloud / DigitalOcean (untuk MVP: DigitalOcean lebih murah)
Container:        Docker + Docker Compose (dev) → Kubernetes (production scale)
CI/CD:            GitHub Actions → deploy ke server
CDN:              Cloudflare (web) + CloudFront (asset)
Storage:          AWS S3 atau Cloudflare R2
Monitoring:       Sentry (error) + Grafana + Prometheus (metrics)
Logging:          ELK Stack (Elasticsearch + Logstash + Kibana)
SSL:              Let's Encrypt / Cloudflare
```

### Third-Party Services

| Kebutuhan | Service | Harga Estimasi |
|---|---|---|
| Payment gateway | Midtrans | 0.7–2.9% per transaksi |
| Payout talent | Flip.id atau Midtrans Payout | Rp 3.500–6.500/transfer |
| WhatsApp notif | Twilio / Wablas | Rp 500–1.500/pesan |
| Push notif | Firebase FCM | Gratis (hingga limit) |
| Email | SendGrid | Gratis 100 email/hari |
| Maps | Google Maps Platform | $200 kredit/bulan gratis |
| SMS OTP | Twilio / Nexmo | ~Rp 800/SMS |
| KYC/Verifikasi | Verihubs / Privy | Rp 3.000–10.000/verifikasi |
| E-sign kontrak | Privy ID / DocuSign | Rp 5.000–15.000/dokumen |

---

## 6. FITUR & MODUL LENGKAP

### 6.1 MODUL AUTENTIKASI & ONBOARDING

#### Halaman Onboarding (4 slides)
- Slide 1: Ilustrasi pencarian jasa + tagline "Semua Jasa, Tinggal Beres."
- Slide 2: Ilustrasi negosiasi + "Negosiasi langsung, harga transparan."
- Slide 3: Ilustrasi kontrak + "Aman dengan kontrak digital & escrow."
- Slide 4: Ilustrasi bintang + "Rating & review nyata dari pengguna."
- Tombol: "Daftar Sekarang" dan "Sudah punya akun? Masuk"

#### Registrasi
**Pilihan Role saat daftar:**
- User (Pencari Jasa)
- Talent/Freelancer (Penyedia Jasa)

**Metode:**
- Email + password (dengan verifikasi email)
- Google OAuth
- Nomor WhatsApp + OTP

**Data wajib saat registrasi:**
- Nama lengkap
- Nomor HP (wajib untuk semua transaksi)
- Lokasi kota/kabupaten
- (Talent) Kategori jasa utama

#### Login
- Email / Google / WhatsApp
- "Ingat saya" (persistent session 30 hari)
- Lupa password: reset via email atau WhatsApp OTP

#### Verifikasi Talent (KYC)
- Upload KTP / SIM / Paspor
- Foto selfie dengan dokumen (liveness check)
- Verifikasi nomor rekening bank (untuk disbursement)
- Status: Pending → Verified → Badge "Terverifikasi"
- Integrasi: Verihubs atau Privy ID

---

### 6.2 MODUL HOME & DISCOVERY

#### Halaman Home
**Hero Section:**
- Search bar besar dengan placeholder dinamis (berganti setiap 3 detik)
- Contoh: "Cari tukang...", "Cari fotografer...", "Cari band wedding..."
- Tombol kategori cepat di bawah search bar

**Kategori Utama (Icon Grid 4 kolom):**
```
Tukang     | Teknisi   | Fotografer | Band/Musik
DJ         | MC        | Dancer     | Psikolog
Cleaning   | Guru      | MUA        | Lainnya...
```

**Section Cards (horizontal scroll):**
1. "Trending Minggu Ini" — berdasarkan volume booking
2. "Terpopuler di [Kota User]" — berdasarkan lokasi
3. "Rating Terbaik" — min rating 4.8, min 10 ulasan
4. "Baru Bergabung" — talent baru dengan harga kompetitif
5. "Promo & Paket Spesial" — talent yang sedang promo

**Talent Card Komponen:**
```
┌─────────────────────────────────┐
│ [Foto Profil]  Nama Talent      │
│                Profesi           │
│ ⭐ 4.9  (127 ulasan)            │
│ 📍 Jakarta Selatan              │
│ ⚡ Respons < 1 jam              │
│ [Badge Verified] [Badge Pro]    │
│ Mulai dari Rp 250.000           │
│ [Lihat Profil] [Booking]        │
└─────────────────────────────────┘
```

**Bottom Navigation Bar:**
- Home (house icon)
- Explore (compass icon)
- Tender (megaphone icon)
- Chat (chat bubble icon)
- Order (clipboard icon)
- Profile (person icon)

---

### 6.3 MODUL EXPLORE / PENCARIAN

#### Fitur Pencarian
- Full-text search (nama talent, deskripsi jasa, keahlian)
- Autocomplete dengan suggestions
- Recent searches (disimpan lokal)
- Voice search (opsional, fase 2)

#### Filter Panel (bottom sheet di mobile)
```
Kategori Jasa:     [pilihan multi-select]
Lokasi:            [kota/radius dari posisi saya]
Harga:             Rp [min] — Rp [max] (slider)
Rating Minimum:    [1★ 2★ 3★ 4★ 4.5★+]
Ketersediaan:      [Hari ini] [Minggu ini] [Pilih tanggal]
Tipe Layanan:      [On-site] [Remote] [Keduanya]
Hanya Verified:    [toggle]
```

#### Sorting Options
- Relevansi (default — algoritma Beresin)
- Rating tertinggi
- Harga termurah
- Harga tertinggi
- Jarak terdekat
- Respons tercepat
- Terbanyak direview

#### Tampilan Hasil
- List view (default mobile)
- Grid view (2 kolom mobile, 3–4 kolom web)
- Map view (tampilkan lokasi talent di peta)

---

### 6.4 MODUL PROFIL TALENT

#### Informasi Profil
- Foto profil (wajib, min 300x300px)
- Banner/cover portfolio (opsional, 1200x400px)
- Nama lengkap / nama usaha
- Kategori jasa (maks 3 kategori)
- Tagline singkat (maks 80 karakter)
- Deskripsi lengkap (maks 2000 karakter)
- Lokasi operasional (kota + radius layanan dalam km)
- Tahun mulai berkarir

#### Verifikasi & Badge
- ✅ Terverifikasi KTP
- ⭐ Pro Member
- 🏆 Top Talent (otomatis jika memenuhi kriteria)
- 🚀 Rising Star (talent baru dengan rating cepat naik)
- 🎯 Spesialist (jika hanya satu kategori, depth lebih tinggi)

#### Paket Harga (Service Packages)
Talent dapat membuat hingga **5 paket** per listing:

```
Paket BASIC:
- Nama paket (contoh: "Sesi 2 jam")
- Deskripsi singkat
- Harga (tetap ATAU mulai dari)
- Durasi layanan
- Apa yang termasuk (checklist)
- Maksimal revisi (jika ada)
- Estimasi penyelesaian
- [Booking Paket Ini]

Paket STANDARD & PREMIUM: (sama strukturnya)
```

#### Portfolio
- Upload foto (maks 20 foto per listing)
- Upload video (maks 5 video, maks 100MB per video)
- Embed YouTube / Instagram link
- Deskripsi per item portfolio
- Tag proyek (wedding, corporate, café, dll)

#### Sertifikat & Pengalaman
- Upload sertifikat / ijazah (maks 10 file)
- Verifikasi sertifikat oleh admin (manual atau otomatis)
- Pengalaman kerja (format: Nama perusahaan, posisi, tahun)

#### Kalender Ketersediaan
- Tandai hari libur / tidak tersedia
- Set jam kerja per hari
- Blokir tanggal otomatis saat ada booking terkonfirmasi
- Sinkronisasi dengan Google Calendar (fase 2)

---

### 6.5 MODUL POSTING JASA (TALENT)

**Form Multi-Step (5 langkah):**

**Step 1 — Informasi Dasar**
- Judul jasa (maks 100 karakter) — contoh: "Grup Band Akustik untuk Wedding & Café"
- Kategori utama + sub-kategori
- Deskripsi lengkap (dengan rich text editor: bold, list, dll)
- Tag/keyword (maks 10 tag)

**Step 2 — Media**
- Upload foto cover (wajib, min 1 foto)
- Upload galeri foto (maks 15 foto)
- Upload video demo (opsional)
- YouTube / Instagram embed link

**Step 3 — Paket Harga**
- Buat minimal 1, maks 5 paket
- Per paket: nama, harga, deskripsi, durasi, inklusi

**Step 4 — Area Layanan**
- Pilih kota/kabupaten yang dilayani
- Set radius maksimal (dalam km dari lokasi talent)
- Tentukan biaya perjalanan di luar radius (Rp/km atau flat fee)

**Step 5 — Setting Lanjutan**
- Respons waktu target (dalam jam)
- Pertanyaan pra-booking (maks 5 pertanyaan custom untuk user isi sebelum booking)
- Persyaratan khusus (misalnya: "Dibutuhkan lokasi dengan colokan listrik 3 fase")
- Kebijakan pembatalan (Fleksibel / Moderat / Ketat)

**Preview & Publish:**
- Preview tampilan listing sebelum publish
- Status: Draft / Published / Paused / Deleted

---

### 6.6 MODUL TENDER (POST PROJECT)

#### Buat Tender (User)

**Step 1 — Informasi Pekerjaan**
- Judul pekerjaan (contoh: "Butuh Fotografer untuk Pernikahan 15 Agustus")
- Kategori jasa
- Deskripsi detail kebutuhan (dengan tips panduan pengisian)
- Upload referensi: foto, video (maks 10 file, maks 50MB total)
- Lokasi pekerjaan (dengan peta interaktif)
- Tanggal pekerjaan (single date atau range)
- Estimasi budget: Rp [min] — Rp [max] ATAU "Saya terbuka dengan penawaran"
- Opsi: Urgent (perlu respons dalam 24 jam)

**Step 2 — Preferensi Talent**
- Pengalaman minimum (tidak ada / 1 tahun / 3 tahun / 5 tahun+)
- Rating minimum (tidak ada / min 4.0 / min 4.5 / min 4.8)
- Tipe talent: Individual / Tim / Keduanya
- Hanya Talent Terverifikasi: [toggle]
- Jumlah penawaran yang diinginkan: Maks [5/10/20/unlimited]

**Step 3 — Review & Publish**
- Preview posting tender
- Estimasi jumlah talent yang akan melihat
- Pilihan visibilitas: Publik / Hanya Talent Undangan (fase 2)
- Durasi tender aktif: 3 / 7 / 14 / 30 hari
- Tombol "Publikasikan Tender"

#### List Tender (untuk Talent)

**Filter:**
- Kategori
- Lokasi (kota + radius)
- Budget (Rp range)
- Urgensi
- Jumlah penawaran yang sudah masuk (sedikit persaingan vs banyak)

**Tender Card:**
```
┌──────────────────────────────────────────┐
│ 🔥 URGENT                                │
│ Fotografer Wedding — 15 Agustus 2025     │
│ 📍 Jakarta Selatan  💰 Rp 2–5 juta      │
│ 📅 Deadline tender: 3 hari lagi          │
│ 👤 12 penawaran masuk                    │
│ "Butuh fotografer berpengalaman untuk    │
│  acara outdoor di villa..."              │
│ [Lihat Detail] [Ajukan Penawaran]        │
└──────────────────────────────────────────┘
```

**Badge Status Tender:**
- 🔥 Urgent
- 💼 Budget Besar (> Rp 5 juta)
- 🆕 Baru (< 6 jam)
- ✅ Sedikit Pesaing (< 3 penawaran)

#### Detail Tender

**Informasi Proyek:**
- Judul, deskripsi lengkap
- Galeri referensi foto/video
- Peta lokasi interaktif
- Budget range
- Tanggal pekerjaan
- Preferensi talent yang diminta
- Durasi tender tersisa

**List Penawaran yang Masuk (visible ke User, hidden ke Talent lain):**
```
┌────────────────────────────────────┐
│ [Foto] Rizky Photography            │
│ ⭐ 4.9 • 87 project selesai        │
│ Harga: Rp 3.500.000                │
│ Durasi: 8 jam (termasuk editing)   │
│ "Saya memiliki pengalaman 5 tahun   │
│  wedding photography..."            │
│ [Chat] [Pilih Talent Ini]          │
└────────────────────────────────────┘
```

**Tombol Aksi untuk User:**
- Pilih Talent (→ langsung ke alur negosiasi/kontrak)
- Chat Talent (→ obrolan khusus tender ini)
- Bandingkan Penawaran (→ lihat tabel perbandingan)
- Tutup Tender (hentikan penerimaan penawaran)

#### Ajukan Penawaran (Talent)

**Form Penawaran:**
- Cover letter / pesan perkenalan (maks 500 karakter, dengan tips)
- Harga penawaran (input nominal Rupiah)
- Estimasi durasi pengerjaan
- DP minimum yang diminta (persentase atau nominal)
- Kebijakan revisi yang ditawarkan
- Portofolio relevan yang dilampirkan (pilih dari portfolio yang sudah ada)

**Sidebar Preview:**
- Profil singkat talent sendiri (rating, badge, foto)
- Perbandingan harga dengan rata-rata tender sejenis

---

### 6.7 MODUL NEGOSIASI & CHAT

#### Chat Interface

**Fitur Chat Dasar:**
- Bubble chat real-time (WebSocket)
- Status pesan: Terkirim / Dibaca (centang)
- Timestamp per pesan
- Reply ke pesan tertentu

**Attachment yang Bisa Dikirim:**
- Gambar (JPEG, PNG, maks 10MB)
- Video (MP4, maks 50MB)
- Dokumen (PDF, DOC, maks 20MB)
- Voice note (rekam langsung, maks 2 menit)
- Lokasi (share pin peta)

**Fitur Khusus Marketplace:**
- **Kirim Quotation:** Form terstruktur yang bisa dikirim dalam chat
  - Nama pekerjaan, detail, harga, jadwal, syarat
  - User bisa Accept / Counter / Reject
- **Counter Offer:** User atau Talent bisa ajukan penawaran balik
- **Sticky Summary:** Di atas chat, tampil ringkasan deal terkini (harga, jadwal, status)

**Status Negosiasi:**
- Diskusi awal
- Penawaran diajukan
- Counter offer
- Penawaran diterima ✅
- Kontrak dibuat
- Kontrak ditandatangani

**Tombol Aksi Cepat (Quick Action Bar):**
- Ajukan Penawaran
- Terima Penawaran
- Buat Kontrak
- Tolak / Akhiri Chat

**Perbandingan Penawaran (untuk Tender):**

Tabel visual yang membandingkan maks 5 talent sekaligus:

| Kriteria | Talent A | Talent B | Talent C |
|---|---|---|---|
| Harga | Rp 3,5 jt | Rp 4 jt | Rp 2,8 jt |
| Rating | 4.9 ⭐ | 4.7 ⭐ | 4.8 ⭐ |
| Durasi | 8 jam | 10 jam | 8 jam |
| Respons | 30 mnt | 2 jam | 15 mnt |
| Pengalaman | 5 tahun | 3 tahun | 4 tahun |
| Highlight | Best Rating | — | Best Price / Fastest |

---

### 6.8 MODUL KONTRAK DIGITAL

#### Pembuatan Kontrak

**Inisiator:** Talent atau User (keduanya bisa)
**Trigger:** Setelah penawaran diterima di chat

**Form Kontrak (auto-filled dari hasil negosiasi):**

**Bagian 1 — Identitas Para Pihak**
- Nama user (Pihak Pertama / Pemberi Kerja)
- Nama talent (Pihak Kedua / Penerima Kerja)
- Nomor identitas masing-masing

**Bagian 2 — Detail Pekerjaan**
- Nama/judul pekerjaan
- Deskripsi detail scope pekerjaan
- Lokasi pelaksanaan
- Tanggal mulai
- Tanggal selesai (deadline)
- Deliverable yang diharapkan

**Bagian 3 — Pembayaran**
- Harga total yang disepakati: Rp ___
- Skema pembayaran:
  - DP: ___% (Rp ___) dibayar sebelum pekerjaan dimulai
  - Pelunasan: ___% (Rp ___) dibayar setelah pekerjaan selesai
  - Atau: Cicilan 3 termin (opsional untuk proyek besar)
- Metode pembayaran: via escrow Beresin

**Bagian 4 — Ketentuan**
- Jumlah revisi yang disepakati: ___ kali
- Syarat pembatalan oleh user: potongan ___% dari DP
- Syarat pembatalan oleh talent: kembalikan DP 100% + denda ___% (opsional)
- Force majeure clause (standar)
- Hukum yang berlaku: Hukum Indonesia

**Bagian 5 — Tanda Tangan Digital**
- User menandatangani (dengan PIN 6 digit atau biometrik)
- Talent menandatangani
- Timestamp + hash dokumen (untuk validitas hukum)
- Integrasi: Privy ID (sudah tersertifikasi Kominfo) untuk tanda tangan digital yang sah secara hukum

**Status Kontrak:**
- Draft (sedang dibuat)
- Menunggu Tanda Tangan User
- Menunggu Tanda Tangan Talent
- Aktif (kedua pihak sudah tanda tangan)
- Selesai
- Dispute
- Dibatalkan

**Output Kontrak:**
- File PDF yang bisa didownload
- Tersimpan di akun keduanya secara permanen
- Nomor kontrak unik (BRS-YYYYMMDD-XXXX)

---

### 6.9 MODUL PEMBAYARAN & ESCROW

#### Alur Pembayaran

```
User membayar DP
       ↓
Dana masuk ke Escrow Beresin (BUKAN langsung ke Talent)
       ↓
Notifikasi ke Talent: "DP diterima, silakan mulai pekerjaan"
       ↓
Pekerjaan berjalan...
       ↓
Talent upload bukti penyelesaian
       ↓
User konfirmasi pekerjaan selesai (atau otomatis setelah 3x24 jam)
       ↓
User bayar pelunasan
       ↓
Dana pelunasan masuk ke Escrow
       ↓
Beresin release dana ke Talent (dikurangi service fee)
       ↓
Talent menerima pembayaran dalam 1x24 jam kerja
```

#### Metode Pembayaran yang Diterima
- **Transfer Bank:** BCA, Mandiri, BNI, BRI, CIMB (Virtual Account)
- **E-Wallet:** GoPay, OVO, DANA, ShopeePay, LinkAja
- **QRIS** (semua dompet digital yang mendukung QRIS)
- **Kartu Kredit/Debit** (Visa, Mastercard — khusus web)
- **Paylater:** Kredivo, Akulaku (fase 2)

#### Skema Escrow
- Dana user → rekening escrow terpisah (bukan rekening operasional Beresin)
- Interest dari dana escrow (jika ada) → ke CSR fund atau charity
- Disbursement ke talent: maksimal 1x24 jam kerja setelah konfirmasi
- Platform mengambil komisi saat disbursement

#### Proteksi Pembeli
- Jika talent tidak mulai kerja dalam 48 jam: user bisa cancel, DP dikembalikan 100%
- Jika talent batalkan sepihak: DP dikembalikan 100% + kompensasi 10% (dari dana talent)
- Jika ada dispute: dana di-hold sampai resolusi

#### Disbursement ke Talent
- Nomor rekening terverifikasi saat KYC
- Transfer via Flip.id (biaya minimal) atau Midtrans Payout
- Laporan disbursement tersedia di dashboard talent

#### Invoice & Laporan Keuangan
- Invoice PDF otomatis per transaksi
- Rekapitulasi bulanan (untuk talent, untuk akuntansi)
- Laporan pajak tahunan (untuk talent dengan pendapatan > PTKP)

---

### 6.10 MODUL TRACKING PEKERJAAN

#### Timeline Progress (untuk User)

```
✅ Kontrak Ditandatangani   [Tgl & Jam]
✅ DP Dibayar               [Tgl & Jam]
✅ Pekerjaan Dimulai        [Tgl & Jam]
🔄 Progress 25%             [Tgl & Jam] — dengan foto bukti
🔄 Progress 50%             [Tgl & Jam] — dengan foto bukti
⬜ Progress 75%             (belum)
⬜ Selesai                  (belum)
⬜ Review Diberikan         (belum)
```

**Fitur Update Progress (Talent):**
- Tombol "Update Progress" di dashboard order aktif
- Upload foto/video dokumentasi (wajib jika progress > 50%)
- Catatan update (opsional, maks 500 karakter)
- Persentase penyelesaian (slider 0–100%)

**Notifikasi Otomatis:**
- User dapat notifikasi setiap ada update progress
- Reminder ke talent jika tidak ada update selama 24 jam
- Alert ke admin jika tidak ada aktivitas 48 jam sebelum deadline

---

### 6.11 MODUL REVIEW & RATING

#### Siapa yang Bisa Review?
- **User** me-review Talent (setelah pekerjaan selesai)
- **Talent** me-review User (setelah pekerjaan selesai) — membantu talent lain menilai klien

#### Komponen Rating (User → Talent)
Skala 1–5 bintang untuk setiap dimensi:
1. Kualitas kerja secara keseluruhan
2. Ketepatan waktu (ontime vs deadline)
3. Komunikasi & responsivitas
4. Keramahan & profesionalisme
5. Value for money (sepadan dengan harga)

**Rating agregat** = rata-rata tertimbang dari semua dimensi
**Tampilan publik** = satu angka (contoh: 4.8) + badge distribusi rating (grafik batang)

#### Form Review
- Rating bintang per dimensi
- Teks review (min 20 karakter, maks 1000 karakter)
- Upload foto hasil kerja (maks 5 foto, opsional)
- Tandai: "Apakah kamu akan memakai talent ini lagi?" (Yes/No)
- Tag highlight: #tepat-waktu #ramah #kualitas-bagus #komunikatif #recommended

#### Kebijakan Review
- Review bisa ditulis dalam 7 hari setelah pekerjaan selesai
- Review tidak bisa dihapus oleh user (hanya oleh admin jika melanggar)
- Talent bisa balas review (respons publik, maks 300 karakter)
- Review yang mengandung SARA/kata kasar akan otomatis di-flag untuk moderasi
- Fake review detection (analisis pola, IP, akun baru)

---

### 6.12 MODUL DISPUTE RESOLUTION

#### Kapan Dispute Bisa Diajukan?
- Pekerjaan tidak sesuai kontrak
- Talent tidak muncul / menghilang
- User menolak membayar pelunasan tanpa alasan jelas
- Kualitas kerja di bawah standar yang dijanjikan

#### Alur Dispute

**Langkah 1 — Pengajuan Dispute (oleh User atau Talent)**
- Pilih alasan dispute (dropdown kategorisasi)
- Deskripsi masalah (maks 2000 karakter)
- Upload bukti: foto, video, screenshot chat, dokumen
- Sistem otomatis freeze escrow saat dispute diajukan

**Langkah 2 — Respons Pihak Lawan (48 jam)**
- Pihak yang disputing diberi notifikasi
- Harus submit respons + bukti balasan dalam 48 jam
- Jika tidak respons dalam 48 jam → keputusan otomatis ke pihak pengirim dispute

**Langkah 3 — Mediasi Beresin**
- Tim Dispute Beresin review semua bukti
- Bisa request informasi tambahan dari kedua pihak
- Timeline resolusi: 3–7 hari kerja
- Keputusan: Dana dirilis ke User / Talent / Split

**Langkah 4 — Keputusan & Eksekusi**
- Keputusan final dikirim via email + in-app
- Eksekusi escrow sesuai keputusan
- Kedua pihak bisa banding satu kali (dalam 24 jam)

**Track Record Dispute:**
- User dengan banyak dispute yang dikalahkan akan mendapat warning
- Talent dengan banyak dispute akan di-suspend → review → reaktivasi atau banned

---

### 6.13 DASHBOARD USER

#### Halaman Utama Dashboard User

**Summary Cards:**
- Booking Aktif: [X] order sedang berjalan
- Menunggu Review: [X] order belum di-review
- Total Jasa Dipesan: [X] sepanjang waktu
- Total Pengeluaran: Rp [X] (bulan ini)

**Tabs:**
1. **Order Aktif** — daftar pekerjaan yang sedang berlangsung
2. **Tender Saya** — tender yang diposting + statusnya
3. **Riwayat** — semua order yang sudah selesai/dibatalkan
4. **Favorit** — talent yang di-save/di-follow
5. **Dompet** — kredit Beresin (jika ada refund/cashback)

**Order Card (di list):**
```
Nama Talent + Profesi
Nama Pekerjaan
Status: [Berjalan] / [Menunggu Konfirmasi] / [Selesai]
Tanggal Pekerjaan
Progress Bar (jika ada)
Harga Final
[Lihat Detail] [Chat]
```

---

### 6.14 DASHBOARD TALENT

#### Halaman Utama Dashboard Talent

**Summary Cards (periode: hari ini / 7 hari / 30 hari / semua waktu):**
- Pendapatan Bersih: Rp [X]
- Order Masuk: [X]
- Order Selesai: [X]
- Rating Rata-rata: [X.X] ⭐
- Conversion Rate: [X]% (penawaran → deal)
- Repeat Customer: [X]%

**Chart Pendapatan:**
- Grafik garis: pendapatan 30 hari terakhir
- Breakdown per kategori (jika talent punya multiple kategori)

**Tabs:**
1. **Order Baru** — notifikasi order langsung + tender yang relevan
2. **Aktif** — pekerjaan sedang berjalan
3. **Menunggu** — booking yang belum dikonfirmasi
4. **Selesai** — arsip pekerjaan selesai
5. **Listing Saya** — kelola posting jasa

**Kalender Booking:**
- Tampilan bulan / minggu
- Hari yang ada booking ditandai
- Tap tanggal → lihat detail order di hari itu

**Notifikasi Penting:**
- "Ada tender baru yang sesuai keahlianmu"
- "Review baru dari [User]"
- "Pembayaran Rp X diterima"
- "Kontrak telah ditandatangani"

---

### 6.15 ADMIN PANEL

#### Dashboard Utama
- Total transaksi hari ini / minggu / bulan (dengan grafik)
- Gross Merchandise Value (GMV) real-time
- Pendapatan platform (service fee)
- Jumlah user aktif
- Jumlah talent aktif
- Order dalam dispute
- Escrow amount (total dana yang di-hold)

#### Modul Verifikasi Talent
- Antrian verifikasi KYC baru
- Detail dokumen yang diupload (foto KTP, selfie)
- Tombol: Approve / Reject (dengan alasan)
- Riwayat verifikasi

#### Modul Dispute Management
- Daftar dispute aktif (sorted by: umur, nilai, urgency)
- Detail kasus: chat, kontrak, bukti kedua pihak
- Tool timeline pekerjaan
- Keputusan: Release ke User / Release ke Talent / Split [X]% : [Y]%
- Catatan admin (internal)

#### Modul User & Talent Management
- Cari user/talent by name/email/ID
- Lihat profil lengkap + riwayat transaksi
- Aksi: Warning / Suspend (temp) / Banned (permanent) / Unban
- Reset password manual
- Impersonate (untuk debugging)

#### Modul Content Moderation
- Queue review yang di-flag (konten tidak pantas)
- Queue foto portfolio yang terindikasi melanggar
- Approve / Remove + log keputusan

#### Modul Keuangan
- Laporan escrow: dana masuk, dana keluar, saldo
- Disbursement history
- Refund log
- Komisi yang diterima platform
- Export ke Excel/CSV

#### Modul Analytics
- Funnel konversi: Discovery → View → Chat → Deal → Selesai
- Kategori terpopuler
- Kota dengan transaksi terbanyak
- Rata-rata nilai transaksi
- Churn rate (talent/user yang tidak aktif lagi)
- NPS (Net Promoter Score) tracker

---

## 7. DATABASE SCHEMA

### Tabel Utama (PostgreSQL)

```sql
-- USERS
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE,
  phone VARCHAR(20) UNIQUE NOT NULL,
  password_hash VARCHAR(255),
  full_name VARCHAR(100) NOT NULL,
  role ENUM('user', 'talent', 'admin') NOT NULL,
  avatar_url TEXT,
  city VARCHAR(100),
  is_verified BOOLEAN DEFAULT FALSE,
  is_active BOOLEAN DEFAULT TRUE,
  google_id VARCHAR(255),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- TALENT PROFILES
CREATE TABLE talent_profiles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) UNIQUE NOT NULL,
  bio TEXT,
  tagline VARCHAR(100),
  years_experience INT,
  service_radius_km INT DEFAULT 20,
  latitude DECIMAL(10,8),
  longitude DECIMAL(11,8),
  is_kyc_verified BOOLEAN DEFAULT FALSE,
  kyc_document_url TEXT,
  kyc_selfie_url TEXT,
  kyc_reviewed_at TIMESTAMPTZ,
  kyc_reviewed_by UUID REFERENCES users(id),
  subscription_tier ENUM('free', 'pro', 'business') DEFAULT 'free',
  subscription_expires_at TIMESTAMPTZ,
  average_rating DECIMAL(3,2) DEFAULT 0,
  total_reviews INT DEFAULT 0,
  total_completed_jobs INT DEFAULT 0,
  response_time_hours INT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- SERVICE CATEGORIES
CREATE TABLE categories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(100) UNIQUE NOT NULL,
  parent_id UUID REFERENCES categories(id),
  icon_url TEXT,
  sort_order INT DEFAULT 0,
  is_active BOOLEAN DEFAULT TRUE
);

-- SERVICE LISTINGS
CREATE TABLE service_listings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  talent_id UUID REFERENCES talent_profiles(id) NOT NULL,
  category_id UUID REFERENCES categories(id) NOT NULL,
  title VARCHAR(200) NOT NULL,
  description TEXT NOT NULL,
  tags TEXT[], -- array of tags
  cover_image_url TEXT,
  gallery_urls TEXT[],
  video_urls TEXT[],
  status ENUM('draft', 'published', 'paused', 'deleted') DEFAULT 'draft',
  view_count INT DEFAULT 0,
  booking_count INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- SERVICE PACKAGES (per listing)
CREATE TABLE service_packages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  listing_id UUID REFERENCES service_listings(id) NOT NULL,
  name VARCHAR(100) NOT NULL, -- "Basic", "Standard", "Premium"
  description TEXT,
  price_amount BIGINT NOT NULL, -- dalam rupiah
  price_type ENUM('fixed', 'starting_from') DEFAULT 'fixed',
  duration_hours INT, -- durasi layanan dalam jam
  inclusions TEXT[], -- apa yang termasuk
  max_revisions INT DEFAULT 0,
  sort_order INT DEFAULT 0
);

-- TENDERS
CREATE TABLE tenders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) NOT NULL,
  category_id UUID REFERENCES categories(id) NOT NULL,
  title VARCHAR(200) NOT NULL,
  description TEXT NOT NULL,
  reference_media_urls TEXT[],
  location_address TEXT,
  location_lat DECIMAL(10,8),
  location_lng DECIMAL(11,8),
  work_date_start DATE,
  work_date_end DATE,
  budget_min BIGINT,
  budget_max BIGINT,
  budget_is_negotiable BOOLEAN DEFAULT TRUE,
  is_urgent BOOLEAN DEFAULT FALSE,
  min_experience_years INT DEFAULT 0,
  min_rating DECIMAL(3,2) DEFAULT 0,
  require_verified BOOLEAN DEFAULT FALSE,
  prefer_team BOOLEAN DEFAULT FALSE,
  max_bids INT,
  expires_at TIMESTAMPTZ,
  status ENUM('active', 'closed', 'completed', 'cancelled') DEFAULT 'active',
  selected_bid_id UUID, -- FK ke bids (setelah user pilih)
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- BIDS (penawaran pada tender)
CREATE TABLE bids (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tender_id UUID REFERENCES tenders(id) NOT NULL,
  talent_id UUID REFERENCES talent_profiles(id) NOT NULL,
  cover_letter TEXT NOT NULL,
  bid_amount BIGINT NOT NULL,
  estimated_duration_hours INT,
  dp_percentage INT DEFAULT 50,
  max_revisions INT DEFAULT 0,
  portfolio_ids UUID[], -- referensi ke portfolio items
  status ENUM('pending', 'shortlisted', 'accepted', 'rejected', 'withdrawn') DEFAULT 'pending',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(tender_id, talent_id)
);

-- BOOKINGS / ORDERS
CREATE TABLE orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_number VARCHAR(20) UNIQUE NOT NULL, -- BRS-20240801-0001
  user_id UUID REFERENCES users(id) NOT NULL,
  talent_id UUID REFERENCES talent_profiles(id) NOT NULL,
  listing_id UUID REFERENCES service_listings(id),
  package_id UUID REFERENCES service_packages(id),
  tender_id UUID REFERENCES tenders(id),
  bid_id UUID REFERENCES bids(id),
  title VARCHAR(200) NOT NULL,
  description TEXT,
  work_date_start DATE,
  work_date_end DATE,
  location_address TEXT,
  final_amount BIGINT NOT NULL,
  dp_amount BIGINT NOT NULL,
  remaining_amount BIGINT NOT NULL,
  platform_fee BIGINT NOT NULL,
  talent_receive_amount BIGINT NOT NULL,
  status ENUM(
    'pending_confirmation',  -- menunggu talent konfirmasi
    'confirmed',             -- talent konfirmasi
    'dp_pending',            -- menunggu DP
    'active',                -- DP dibayar, pekerjaan dimulai
    'completed_pending',     -- talent tandai selesai
    'completed',             -- user konfirmasi selesai
    'cancelled',
    'dispute'
  ) DEFAULT 'pending_confirmation',
  progress_percentage INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- CONTRACTS
CREATE TABLE contracts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  contract_number VARCHAR(20) UNIQUE NOT NULL,
  order_id UUID REFERENCES orders(id) UNIQUE NOT NULL,
  user_id UUID REFERENCES users(id) NOT NULL,
  talent_id UUID REFERENCES talent_profiles(id) NOT NULL,
  work_title VARCHAR(200) NOT NULL,
  work_description TEXT NOT NULL,
  work_start_date DATE NOT NULL,
  work_end_date DATE NOT NULL,
  location TEXT,
  total_amount BIGINT NOT NULL,
  dp_percentage INT NOT NULL,
  dp_amount BIGINT NOT NULL,
  max_revisions INT DEFAULT 0,
  cancellation_terms TEXT,
  user_signed_at TIMESTAMPTZ,
  talent_signed_at TIMESTAMPTZ,
  user_signature_hash VARCHAR(255),
  talent_signature_hash VARCHAR(255),
  document_hash VARCHAR(255),
  pdf_url TEXT,
  status ENUM('draft', 'awaiting_user', 'awaiting_talent', 'active', 'completed', 'cancelled', 'dispute') DEFAULT 'draft',
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- PAYMENTS
CREATE TABLE payments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID REFERENCES orders(id) NOT NULL,
  payment_type ENUM('dp', 'remaining', 'full') NOT NULL,
  amount BIGINT NOT NULL,
  method ENUM('bank_transfer', 'ewallet', 'qris', 'credit_card', 'va') NOT NULL,
  provider VARCHAR(50), -- "midtrans", "gopay", "ovo", dll
  external_transaction_id VARCHAR(255),
  status ENUM('pending', 'success', 'failed', 'expired', 'refunded') DEFAULT 'pending',
  paid_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- DISBURSEMENTS (pembayaran ke talent)
CREATE TABLE disbursements (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID REFERENCES orders(id) NOT NULL,
  talent_id UUID REFERENCES talent_profiles(id) NOT NULL,
  amount BIGINT NOT NULL,
  bank_code VARCHAR(10),
  bank_account VARCHAR(30),
  account_name VARCHAR(100),
  external_reference VARCHAR(255),
  status ENUM('pending', 'processing', 'success', 'failed') DEFAULT 'pending',
  disbursed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- CHATS
CREATE TABLE chat_rooms (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  type ENUM('order', 'bid_negotiation', 'inquiry') NOT NULL,
  order_id UUID REFERENCES orders(id),
  tender_id UUID REFERENCES tenders(id),
  bid_id UUID REFERENCES bids(id),
  user_id UUID REFERENCES users(id) NOT NULL,
  talent_id UUID REFERENCES talent_profiles(id) NOT NULL,
  last_message_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE chat_messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  room_id UUID REFERENCES chat_rooms(id) NOT NULL,
  sender_id UUID REFERENCES users(id) NOT NULL,
  message_type ENUM('text', 'image', 'video', 'audio', 'file', 'quotation', 'offer', 'system') DEFAULT 'text',
  content TEXT,
  media_url TEXT,
  metadata JSONB, -- untuk quotation/offer: {amount, terms, etc}
  is_read BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- REVIEWS
CREATE TABLE reviews (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID REFERENCES orders(id) NOT NULL,
  reviewer_id UUID REFERENCES users(id) NOT NULL,
  reviewee_id UUID REFERENCES users(id) NOT NULL,
  review_type ENUM('user_to_talent', 'talent_to_user') NOT NULL,
  rating_overall DECIMAL(3,2) NOT NULL,
  rating_quality DECIMAL(3,2),
  rating_timeliness DECIMAL(3,2),
  rating_communication DECIMAL(3,2),
  rating_friendliness DECIMAL(3,2),
  rating_value DECIMAL(3,2),
  comment TEXT,
  photo_urls TEXT[],
  would_recommend BOOLEAN,
  talent_response TEXT,
  talent_responded_at TIMESTAMPTZ,
  is_flagged BOOLEAN DEFAULT FALSE,
  is_published BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- DISPUTES
CREATE TABLE disputes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID REFERENCES orders(id) NOT NULL,
  filed_by UUID REFERENCES users(id) NOT NULL,
  reason_category VARCHAR(100) NOT NULL,
  description TEXT NOT NULL,
  evidence_urls TEXT[],
  respondent_response TEXT,
  respondent_evidence_urls TEXT[],
  admin_notes TEXT,
  decision ENUM('user_wins', 'talent_wins', 'split', 'pending') DEFAULT 'pending',
  split_user_percentage INT,
  resolved_at TIMESTAMPTZ,
  resolved_by UUID REFERENCES users(id),
  status ENUM('open', 'under_review', 'resolved', 'appealed') DEFAULT 'open',
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- NOTIFICATIONS
CREATE TABLE notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) NOT NULL,
  type VARCHAR(100) NOT NULL, -- "new_order", "payment_received", dll
  title VARCHAR(200) NOT NULL,
  body TEXT NOT NULL,
  action_url TEXT,
  metadata JSONB,
  is_read BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 8. API ENDPOINTS

### Konvensi API

```
Base URL:     https://api.beresin.id/v1
Auth Header:  Authorization: Bearer {jwt_token}
Format:       JSON
Pagination:   ?page=1&limit=20
Response:     { success: bool, data: {}, meta: {pagination} }
Error:        { success: false, error: { code, message } }
```

### Auth Endpoints

```
POST   /auth/register               Registrasi user baru
POST   /auth/login                  Login
POST   /auth/login/google           OAuth Google
POST   /auth/login/whatsapp/send-otp  Kirim OTP WhatsApp
POST   /auth/login/whatsapp/verify  Verifikasi OTP
POST   /auth/refresh-token          Refresh JWT
POST   /auth/logout                 Logout
POST   /auth/forgot-password        Request reset password
POST   /auth/reset-password         Reset password dengan token
```

### User Endpoints

```
GET    /users/me                    Profil sendiri
PUT    /users/me                    Update profil
PUT    /users/me/avatar             Upload avatar
GET    /users/:id                   Profil publik user (limited)
```

### Talent Endpoints

```
GET    /talents                     List semua talent (dengan filter & pagination)
GET    /talents/:id                 Profil talent lengkap
PUT    /talents/me                  Update profil talent sendiri
POST   /talents/me/kyc              Submit dokumen KYC
GET    /talents/me/stats            Statistik talent sendiri
GET    /talents/me/calendar         Kalender ketersediaan
PUT    /talents/me/calendar         Update kalender
```

### Listing Endpoints

```
GET    /listings                    List jasa (dengan filter: category, location, price, rating)
GET    /listings/:id                Detail listing
POST   /listings                    Buat listing baru (talent only)
PUT    /listings/:id                Update listing
DELETE /listings/:id                Hapus listing
PUT    /listings/:id/status         Pause/unpause listing
GET    /listings/:id/packages       Semua paket dari listing
POST   /listings/:id/packages       Tambah paket
PUT    /packages/:id                Update paket
DELETE /packages/:id                Hapus paket
```

### Tender Endpoints

```
GET    /tenders                     List tender (filter: category, location, budget, status)
GET    /tenders/:id                 Detail tender
POST   /tenders                     Buat tender (user only)
PUT    /tenders/:id                 Update tender
DELETE /tenders/:id                 Hapus tender
PUT    /tenders/:id/close           Tutup tender

POST   /tenders/:id/bids            Ajukan penawaran (talent only)
GET    /tenders/:id/bids            List penawaran (user only, view all bids)
PUT    /bids/:id                    Update penawaran
DELETE /bids/:id                    Tarik penawaran
PUT    /bids/:id/accept             User pilih pemenang bid
```

### Order Endpoints

```
GET    /orders                      List order (user & talent masing2 lihat yang relevan)
GET    /orders/:id                  Detail order
POST   /orders                      Buat order (dari listing langsung)
PUT    /orders/:id/confirm          Talent konfirmasi order
PUT    /orders/:id/cancel           Batalkan order
PUT    /orders/:id/progress         Update progress (talent only) — termasuk upload bukti
PUT    /orders/:id/complete         Tandai selesai (talent)
PUT    /orders/:id/confirm-complete User konfirmasi selesai
```

### Contract Endpoints

```
POST   /contracts                   Buat kontrak (dari order)
GET    /contracts/:id               Detail kontrak
PUT    /contracts/:id               Edit kontrak (sebelum ditandatangani)
POST   /contracts/:id/sign          Tandatangani kontrak
GET    /contracts/:id/pdf           Download PDF kontrak
```

### Payment Endpoints

```
POST   /payments/initiate           Inisiasi pembayaran (get payment URL)
POST   /payments/callback           Webhook dari Midtrans (internal)
GET    /payments/order/:order_id    Status pembayaran per order
GET    /payments/history            Riwayat pembayaran
GET    /disbursements/history       Riwayat pencairan (talent only)
```

### Chat Endpoints

```
GET    /chats                       List chat rooms
GET    /chats/:room_id              Detail room + messages
POST   /chats/rooms                 Buat chat room baru
POST   /chats/:room_id/messages     Kirim pesan (teks/file)
PUT    /chats/:room_id/read         Tandai pesan sudah dibaca
```

### Review Endpoints

```
POST   /reviews                     Buat review (setelah order selesai)
GET    /reviews/talent/:talent_id   Semua review untuk talent tertentu
PUT    /reviews/:id/respond         Talent balas review
```

### Dispute Endpoints

```
POST   /disputes                    Ajukan dispute
GET    /disputes/:id                Detail dispute
POST   /disputes/:id/respond        Respons dispute (pihak lawan)
```

### Admin Endpoints

```
GET    /admin/dashboard             Statistik utama
GET    /admin/kyc/queue             Antrian verifikasi KYC
PUT    /admin/kyc/:id/approve       Approve KYC
PUT    /admin/kyc/:id/reject        Reject KYC (dengan alasan)
GET    /admin/disputes              Semua dispute aktif
PUT    /admin/disputes/:id/resolve  Resolve dispute
GET    /admin/users                 Semua user
PUT    /admin/users/:id/suspend     Suspend user
PUT    /admin/users/:id/ban         Ban user
GET    /admin/escrow/summary        Ringkasan escrow
GET    /admin/analytics/revenue     Laporan pendapatan
```

---

## 9. ALUR PENGGUNA (USER FLOW)

### Flow A — User Booking Langsung

```
Buka App → Home
   ↓
Cari jasa (search / browse kategori)
   ↓
Lihat listing talent → Pilih paket
   ↓
[Jika butuh kustomisasi] Chat dengan talent
   ↓
Klik "Booking Sekarang"
   ↓
Isi form booking: tanggal, lokasi, catatan
   ↓
Review summary + harga + fee
   ↓
Klik "Lanjut ke Pembayaran DP"
   ↓
Pilih metode pembayaran → Bayar
   ↓
[Sistem kirim notif ke Talent]
   ↓
Talent konfirmasi (maks 24 jam)
   ↓
Kontrak otomatis dibuat → Kedua pihak tanda tangan
   ↓
Pekerjaan berjalan (tracking progress)
   ↓
User terima hasil → Konfirmasi selesai
   ↓
User bayar pelunasan
   ↓
Escrow release ke Talent
   ↓
User tulis review → Selesai ✅
```

### Flow B — User Post Tender

```
Buka App → Tab Tender → Buat Tender
   ↓
Isi form 3 langkah → Publish
   ↓
[Sistem matching → notif ke talent relevan]
   ↓
Talent-talent ajukan penawaran
   ↓
User terima notif "Ada X penawaran masuk"
   ↓
User buka tender → Lihat daftar penawaran
   ↓
User chat dengan kandidat (opsional negosiasi)
   ↓
User bandingkan penawaran (compare view)
   ↓
User pilih talent → "Terima Penawaran"
   ↓
[Masuk ke alur kontrak + pembayaran — sama seperti Flow A]
```

### Flow C — Talent Menerima Order

```
Notifikasi masuk: "Ada order baru dari [User]!"
   ↓
Buka dashboard → Lihat order detail
   ↓
Review: nama klien, pekerjaan, lokasi, tanggal, harga
   ↓
Konfirmasi / Tolak (dengan alasan jika ditolak)
   ↓
[Jika dikonfirmasi] → Chat dengan user jika perlu
   ↓
[Kontrak dibuat] → Tanda tangan
   ↓
[Notif: DP diterima] → Mulai kerjakan
   ↓
Update progress (dengan foto bukti)
   ↓
Tandai selesai → Upload hasil akhir
   ↓
Tunggu konfirmasi user (maks 3x24 jam, otomatis jika tidak ada respons)
   ↓
Dana di-release → Terima pembayaran 💰
```

### Flow D — Talent Bidding Tender

```
Buka Tab Tender → Browse listing tender
   ↓
Filter berdasarkan kategori, lokasi, budget
   ↓
Buka detail tender yang menarik
   ↓
Klik "Ajukan Penawaran"
   ↓
Tulis cover letter + input harga + estimasi
   ↓
Submit penawaran
   ↓
Tunggu respon dari user
   ↓
[Jika user mau chat] → Negosiasi
   ↓
[Jika terpilih] → Notif "Penawaran Anda Diterima!"
   ↓
[Lanjut ke alur kontrak]
```

---

## 10. DESAIN & UI/UX SPEC

### Brand Identity

**Nama:** Beresin
**Tagline:** "Semua Jasa, Tinggal Beres."
**Nada:** Profesional tapi ramah, percaya diri, lokal Indonesia

### Color Palette

```
Primary — Navy Blue:
  Main:       #1A2B5E
  Light:      #2D4A8C
  Dark:       #0F1A3A

Secondary — Warm Orange:
  Main:       #F97316
  Light:      #FB923C
  Dark:       #EA580C

Neutral:
  Background: #F8F9FC
  Surface:    #FFFFFF
  Border:     #E2E8F0
  Text Primary:    #0F172A
  Text Secondary:  #475569
  Text Muted:      #94A3B8

Semantic:
  Success:    #16A34A
  Warning:    #D97706
  Error:      #DC2626
  Info:       #0EA5E9
```

### Typography

```
Font Utama:  Plus Jakarta Sans (Google Fonts)
  - Heading H1: 28px, Bold (700), letter-spacing: -0.5px
  - Heading H2: 22px, Bold (700)
  - Heading H3: 18px, SemiBold (600)
  - Heading H4: 16px, SemiBold (600)
  - Body Large: 16px, Regular (400), line-height: 1.6
  - Body:       14px, Regular (400), line-height: 1.5
  - Caption:    12px, Regular (400)
  - Label:      12px, SemiBold (600), letter-spacing: 0.5px
```

### Spacing System (8px grid)

```
xs:  4px
sm:  8px
md:  16px
lg:  24px
xl:  32px
2xl: 48px
3xl: 64px
```

### Component Specs

**Border Radius:**
- Button: 12px
- Card: 16px
- Input: 10px
- Badge: 100px (pill)
- Modal: 20px (top corners)

**Shadow:**
```
Card Shadow: 0 1px 3px rgba(0,0,0,0.08), 0 4px 12px rgba(0,0,0,0.05)
Elevated: 0 4px 20px rgba(0,0,0,0.12)
Dropdown: 0 8px 30px rgba(0,0,0,0.15)
```

**Button Styles:**
- Primary: bg #1A2B5E, text white, hover darken 10%
- Secondary: bg transparent, border 1.5px #1A2B5E, text #1A2B5E
- Accent: bg #F97316, text white (untuk CTA utama)
- Danger: bg #DC2626, text white
- Ghost: bg transparent, text #475569

**Input Fields:**
- Height: 48px (mobile) / 44px (web)
- Border: 1.5px solid #E2E8F0
- Focus: border-color #1A2B5E, box-shadow 0 0 0 3px rgba(26,43,94,0.1)
- Error: border-color #DC2626

### Panduan Desain Mobile

- **Minimum tap target:** 44x44px
- **Safe area:** padding bottom 16px + tinggi tab bar
- **Scroll:** selalu vertikal, hindari horizontal scroll dalam konten utama
- **Skeleton loading:** gunakan skeleton screen, bukan spinner untuk list konten
- **Empty states:** selalu ada ilustrasi + teks + CTA di setiap empty state
- **Error states:** pesan error yang jelas + tombol retry
- **Pull to refresh:** di semua list

### Animasi & Transisi

- Durasi standar: 200–300ms
- Easing: ease-out untuk masuk, ease-in untuk keluar
- Page transition: slide dari kanan (push), slide ke kanan (pop)
- Modal: slide dari bawah (bottom sheet) atau fade (modal center)
- Skeleton to content: fade in
- Button tap: scale 0.97 + darken

---

## 11. SISTEM KEAMANAN & ESCROW

### Keamanan Autentikasi

- Password: bcrypt dengan salt rounds 12
- JWT: access token 15 menit, refresh token 30 hari
- OTP: 6 digit, expired 5 menit, maks 3 kali percobaan
- Rate limiting: maks 5 request login gagal → block 15 menit
- 2FA opsional (via authenticator app) — fase 2

### Keamanan Data

- All data in transit: TLS 1.3
- PII encryption at rest (KTP number, rekening bank)
- Database: enkripsi full-disk
- File storage: private bucket, akses via signed URL (expired 1 jam)
- No raw card number disimpan (handled by Midtrans PCI DSS)

### Keamanan Escrow

- Dana escrow disimpan di rekening terpisah (Virtual Account dedicated)
- Setiap transaksi memiliki escrow record tersendiri
- Tidak ada akses manual ke escrow selain melalui sistem terkontrol
- Audit log setiap perubahan status escrow
- Release escrow hanya terjadi setelah:
  - Konfirmasi user (manual), ATAU
  - Auto-release 3x24 jam setelah talent tandai selesai
  - Admin decision (untuk dispute)

### Fraud Prevention

- Device fingerprinting (untuk deteksi multiple akun)
- IP monitoring (flag IP luar negeri yang tidak biasa)
- Velocity check (transaksi abnormal dalam waktu singkat)
- ML-based review fraud detection (fake review pattern)
- Manual review untuk transaksi > Rp 10 juta

---

## 12. NOTIFIKASI & KOMUNIKASI

### Push Notification (FCM)

| Event | Penerima | Contoh Pesan |
|---|---|---|
| Order baru | Talent | "Ada order baru dari Budi untuk Fotografer Wedding" |
| Order dikonfirmasi | User | "Rizky Photography mengkonfirmasi ordermu!" |
| DP diterima | Talent | "DP Rp 1.500.000 diterima. Silakan mulai pekerjaan." |
| Progress update | User | "Rizky mengupdate progress ke 50% — lihat foto terbaru" |
| Pekerjaan selesai | User | "Pekerjaan selesai! Harap konfirmasi dalam 3 hari." |
| Pembayaran diterima | Talent | "Rp 3.200.000 sudah masuk ke rekeningmu!" |
| Penawaran baru | User | "Ada 3 penawaran baru untuk tender Fotografer-mu" |
| Pesan baru | Keduanya | "[Nama]: [preview pesan]" |
| Dispute diajukan | Keduanya | "Ada pengajuan dispute pada order #BRS-001" |

### Email Notification

- Registrasi berhasil (welcome email + panduan memulai)
- Reset password
- Konfirmasi order (dengan detail lengkap + attachment kontrak PDF)
- Invoice / kuitansi pembayaran
- Disbursement sukses (ke talent)
- Weekly summary (ke talent: pendapatan, statistik)
- Reminder review (3 hari setelah order selesai jika belum review)

### WhatsApp Notification (via Twilio/WABA)

- OTP login & verifikasi
- Konfirmasi booking singkat
- Reminder pembayaran
- Notif transaksi penting

### In-App Notification Center

- Bell icon di header, menampilkan jumlah unread
- List semua notifikasi dengan timestamp
- Filter: Semua / Order / Pembayaran / Sistem
- Tap notif → navigasi ke halaman relevan

---

## 13. ADMIN PANEL

### Teknologi Admin Panel

**Rekomendasi:** React + Ant Design Pro (atau Refine.dev)
- Terpisah dari aplikasi user (subdomain: admin.beresin.id)
- Akses hanya via VPN / IP whitelist

### Role Admin

- **Super Admin:** Akses penuh semua fitur
- **Moderator:** Verifikasi KYC, moderasi konten, resolve dispute
- **Finance:** Pantau escrow, approve disbursement, laporan keuangan
- **Customer Support:** Lihat detail transaksi, chat dengan user/talent (read-only order/chat)

### Fitur Admin (Detail)

**Dashboard:**
- GMV hari ini / 7 hari / 30 hari / YTD (grafik realtime)
- Jumlah registrasi baru (user & talent)
- Order baru vs selesai vs dispute
- Top 10 kategori, top 10 kota, top 10 talent
- Escrow balance real-time

**KYC Queue:**
- Antrian dengan prioritas (diurutkan by tanggal submit)
- Preview foto KTP & selfie side by side
- Input catatan jika reject
- Target SLA: review dalam 24 jam kerja

**Order Monitor:**
- Semua order dengan filter status
- Flag otomatis untuk: order yang stagnan > 48 jam, approaching deadline
- Override status (emergency)

**Dispute Resolution:**
- Timeline pekerjaan visual
- Chat history viewer
- File evidence viewer
- Form keputusan dengan logika kalkulasi escrow
- Catatan internal (tidak terlihat oleh user/talent)

**User Management:**
- Search: nama, email, HP, nomor KTP
- View full profile + transaction history
- Action log (semua tindakan admin pada user ini)

**Analytics:**
- Funnel conversion rates
- Cohort analysis (retensi user & talent)
- Category growth
- Geographic heatmap (kota dengan transaksi terbanyak)
- Revenue attribution

---

## 14. ROADMAP PENGEMBANGAN

### 14.1 Prinsip Delivery
- Sprint berjalan 2 minggu.
- Dev B selalu 1 sprint lebih awal untuk API dan kontrak data.
- Flutter mobile dan Go backend adalah dua aplikasi terpisah.
- Semua perubahan perilaku harus melalui OpenAPI, contoh payload, dan acceptance criteria.
- Tidak ada shared runtime atau shared state langsung antar aplikasi.

### 14.2 Fase 0 — Persiapan (Sprint 0, Minggu 1–2)
- Finalisasi PRD, scope MVP, dan milestone launch.
- Setup repo terpisah untuk Flutter app dan Go backend.
- Setup environment dev, staging, production.
- Setup OpenAPI baseline, secret management, CI/CD, logging, monitoring.
- Siapkan Figma untuk design system dan screen utama.
- Buat test data, seed, dan dummy account untuk QA awal.

### 14.3 Fase 1 — MVP Core (Sprint 1–6, Bulan 2–4)
Target: platform bisa dipakai transaksi nyata pada jalur booking, chat, payment, dan review.

**Sprint 1 — Onboarding & Autentikasi**
- Backend Go: register, login, refresh token, OTP, lupa password, audit log, rate limit.
- Flutter: onboarding, register/login, OTP, forgot password, token storage.
- Output: user bisa daftar dan login sesuai role.

**Sprint 2 — Home & Discovery**
- Backend Go: listing, filter, sort, pagination, search dasar.
- Flutter: home, kategori, search, filter sheet, empty/loading state.
- Output: user bisa cari dan membuka listing.

**Sprint 3 — Profil Talent & Posting Jasa**
- Backend Go: profil publik, portfolio, kalender, KYC upload, CRUD listing.
- Flutter: profil talent, gallery, review, form posting jasa multi-step.
- Output: talent bisa publish listing draft.

**Sprint 4 — Chat & Negosiasi**
- Backend Go: WebSocket auth, room, message, read receipt, typing, quotation.
- Flutter: chat list, chat room, reply mode, attachment, voice note, quotation card.
- Output: percakapan realtime stabil.

**Sprint 5 — Order, Escrow, Payment**
- Backend Go: order create/confirm/cancel, payment initiate, webhook, escrow ledger, disbursement.
- Flutter: booking flow, summary, Midtrans WebView, order timeline, payment success.
- Output: sandbox transaction end-to-end.

**Sprint 6 — Review, Notification, Dashboard Dasar**
- Backend Go: review create/respond, notifikasi push/email, rating aggregation.
- Flutter: review form, review list, badge notifikasi, dashboard ringkas user/talent.
- Output: review tersimpan dan notifikasi tampil.

### 14.4 Fase 2 — Fitur Lengkap (Sprint 7–10, Bulan 5–7)
Target: tender, kontrak digital, dispute, subscription, dan admin capability.

**Sprint 7 — Dispute, Subscription, Hardening**
- Backend Go: dispute create/resolve, freeze escrow, subscription billing, reminder jobs.
- Flutter: dispute view, evidence upload, subscription status, role guard.
- Output: skenario sengketa dan subscription siap operasi.

**Sprint 8 — Tender, Bid, Kontrak Digital**
- Backend Go: tender CRUD, bid accept/reject, matching engine, kontrak auto-fill, sign flow, PDF output.
- Flutter: tender list/detail, create tender form, bid flow, kontrak screen, download kontrak.
- Output: alur tender-to-order-to-contract lengkap.

**Sprint 9 — Admin API & Analitik**
- Backend Go: admin endpoints, metric aggregation, dispute queue, audit exports, analytics summary.
- Flutter atau web ops terpisah: admin/ops dashboard bila dibutuhkan.
- Output: tim internal bisa memantau operasi utama.

**Sprint 10 — Scale & Release Readiness**
- Backend Go: profiling query, index tuning, queue monitoring, backup drill, incident playbook.
- Flutter: performance pass, offline state, caching, crash fix, release preparation.
- Output: sistem siap launch dengan rollback plan.

### 14.5 Fase 3 — Growth & Scale (Sprint 11+)
- Boost listing dan featured placement.
- Referral program.
- Insurance integration.
- Talent team management.
- Voucher, promo, loyalty points.
- B2B / korporat booking flow.
- Google Calendar sync.
- Video call awal bila dibutuhkan.
- Ekspansi kota berbasis data demand.

### 14.6 Definition of Done Roadmap
- Fitur punya kontrak API yang jelas.
- Frontend Flutter dan backend Go lulus integration test.
- Staging bisa didemo tanpa manual workaround.
- Logging, metric, dan error tracking aktif.
- Dokumentasi singkat untuk fitur yang baru ditutup.

---

## 15. ESTIMASI BIAYA & TIM

### Komposisi Tim Minimal (untuk MVP)

| Role | Jumlah | Estimasi Gaji/bulan |
|---|---|---|
| Product Manager | 1 | Rp 10–15 juta |
| UI/UX Designer | 1 | Rp 7–12 juta |
| Backend Developer (Go senior) | 1–2 | Rp 15–25 juta/orang |
| Mobile Developer (Flutter) | 1–2 | Rp 12–20 juta/orang |
| Ops/Web Dashboard Developer | 1 | Rp 10–15 juta |
| QA Engineer | 1 | Rp 7–10 juta |
| DevOps (part-time / freelance) | 1 | Rp 5–10 juta |

**Total Tim:** 6–8 orang
**Estimasi Burn Rate:** Rp 75–125 juta/bulan (tim + operational)

### Biaya Infrastruktur (MVP)

| Item | Estimasi/bulan |
|---|---|
| Cloud Server (DigitalOcean/AWS) | $200–400 (~Rp 3–6 juta) |
| Database (Managed PostgreSQL) | $50–100 |
| Redis (Managed) | $30–60 |
| S3/R2 Storage | $20–50 |
| CDN (Cloudflare) | $20 |
| Monitoring (Sentry, Grafana Cloud) | $50 |
| Email (SendGrid) | $15–25 |
| Push Notification (FCM) | Free – $10 |
| **Total Infra** | **~Rp 7–12 juta/bulan** |

### Biaya Third-Party Services

| Service | Model | Estimasi |
|---|---|---|
| Midtrans | 0.7–2.9% per transaksi | Bergantung GMV |
| Flip.id (payout) | Rp 3.500–5.000/transfer | Bergantung volume |
| Privy (e-sign) | Rp 5.000–10.000/kontrak | Bergantung volume |
| Verihubs (KYC) | Rp 3.000–5.000/verifikasi | Bergantung volume |
| Twilio (OTP) | Rp 800–1.200/OTP | Bergantung volume |

### Estimasi Total Investasi MVP (4 bulan)

```
Tim Pengembang (4 bulan):   Rp  400–520 juta
Infrastruktur (4 bulan):    Rp   30–50 juta
Desain & Aset:              Rp   20–30 juta
Legal & Perizinan:          Rp   10–20 juta
Marketing MVP Launch:       Rp   50–100 juta
Contingency (15%):          Rp   75–105 juta
──────────────────────────────────────────────
TOTAL INVESTASI MVP:        Rp 585–825 juta
```

*Catatan: Bisa ditekan dengan outsourcing beberapa role atau menggunakan bootcamp developer.*

### Break-Even Analysis (Proyeksi Kasar)

Asumsi:
- Rata-rata nilai transaksi: Rp 500.000
- Service fee: 15% (5% user + 10% talent) = Rp 75.000/transaksi
- Break-even pada ~1.000 transaksi/bulan
- Target bulan ke-6 setelah launch: 500 transaksi/bulan

---

## 16. RISIKO & MITIGASI

### Risiko Teknis

| Risiko | Probabilitas | Dampak | Mitigasi |
|---|---|---|---|
| Downtime payment gateway | Medium | Tinggi | Fallback ke payment provider lain, circuit breaker |
| Data breach | Low | Sangat Tinggi | Enkripsi, penetration testing rutin, bug bounty |
| Skala tidak cukup saat viral | Medium | Tinggi | Auto-scaling, load testing berkala |
| Bug pada escrow logic | Low | Sangat Tinggi | Unit test coverage 90%+, code review ketat |

### Risiko Bisnis

| Risiko | Probabilitas | Dampak | Mitigasi |
|---|---|---|---|
| Rendahnya adopsi talent | Tinggi | Tinggi | Program onboarding aktif, komisi 0% 3 bulan pertama |
| Trust issue (penipuan) | Medium | Tinggi | KYC ketat, escrow, review transparan |
| Kompetitor masuk | Medium | Medium | Fokus pada komunitas lokal, fitur yang dalam |
| Regulasi fintech | Low | Medium | Konsultasi hukum, kemitraan dengan fintech berizin |
| Talent tidak profesional | Tinggi | Medium | SLA response, pelatihan, sistem reputasi ketat |

### Risiko Operasional

| Risiko | Probabilitas | Dampak | Mitigasi |
|---|---|---|---|
| Dispute volume tinggi | Medium | Tinggi | Tim CS yang cukup, panduan kontrak yang jelas |
| Seasonal demand | Medium | Medium | Promosi off-peak, diversifikasi kategori |
| Talent churn | Medium | Tinggi | Program loyalitas, analitik retensi |

---

## LAMPIRAN

### Checklist Pre-Launch

**Teknis:**
- [ ] Semua endpoint API memiliki autentikasi
- [ ] Rate limiting diaktifkan
- [ ] SSL certificate terpasang
- [ ] Backup database otomatis (daily)
- [ ] Error monitoring (Sentry) aktif
- [ ] Load testing dilakukan (target: 100 concurrent users)
- [ ] OWASP Top 10 vulnerabilities dicek
- [ ] Penetration testing (minimal manual)

**Legal & Compliance:**
- [ ] Syarat & Ketentuan disiapkan (oleh lawyer)
- [ ] Kebijakan Privasi (UUPD)
- [ ] Perjanjian Escrow dengan bank
- [ ] Perizinan usaha (PT + izin fintech jika perlu)
- [ ] Template kontrak divalidasi lawyer

**Bisnis:**
- [ ] Sistem onboarding talent (awal: 50–100 talent di 1–2 kota)
- [ ] Strategi launch (media sosial, komunitas, press release)
- [ ] Customer support channel aktif (WA, email, in-app)
- [ ] SLA response time ditetapkan dan disosialisasikan
- [ ] Panduan penggunaan (video tutorial + FAQ)

### Glosarium

| Term | Definisi |
|---|---|
| Talent | Freelancer / penyedia jasa yang terdaftar di Beresin |
| User | Pengguna yang mencari / memesan jasa |
| Listing | Posting jasa yang dibuat oleh Talent |
| Tender | Kebutuhan pekerjaan yang dipost oleh User untuk direspons Talent |
| Bid | Penawaran yang diajukan Talent pada sebuah Tender |
| Escrow | Dana yang ditahan platform sebagai jaminan selama pekerjaan berlangsung |
| DP (Down Payment) | Uang muka yang dibayar sebelum pekerjaan dimulai |
| Dispute | Sengketa antara User dan Talent yang dimediasi oleh Beresin |
| KYC | Know Your Customer — proses verifikasi identitas Talent |
| GMV | Gross Merchandise Value — total nilai transaksi yang terjadi di platform |
| Disbursement | Pencairan dana dari escrow ke rekening Talent |

---

*Dokumen ini adalah PRD hidup (living document) yang akan terus diperbarui seiring perkembangan produk.*
*Versi: 1.0 | Tanggal: Mei 2025 | Status: Draft untuk Review*