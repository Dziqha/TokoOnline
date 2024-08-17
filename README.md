
---

# 🛒 Toko Online API - Golang

Ini adalah proyek API untuk Toko Online yang dibangun menggunakan bahasa Go (Golang). API ini menyediakan berbagai endpoint untuk mengelola produk, kategori, pengguna, pesanan, dan fitur lain yang umum ditemukan dalam sebuah sistem e-commerce.

## 🚀 Fitur Utama

- **Manajemen Produk:** CRUD (Create, Read, Update, Delete) produk dengan detail seperti nama, deskripsi, harga, dan stok.
- **Kategori Produk:** Mengelola kategori produk untuk memudahkan pengelompokan barang.
- **Otentikasi Pengguna:** Sistem login dan registrasi pengguna menggunakan JWT (JSON Web Token).
- **Manajemen Pesanan:** Membuat, melihat, dan mengelola pesanan yang dilakukan oleh pengguna.
- **Keranjang Belanja:** Fitur untuk menambahkan dan mengelola produk dalam keranjang belanja sebelum melakukan checkout.


## 🛠️ Teknologi yang Digunakan

- 🖥️ **Golang:** Bahasa pemrograman utama yang digunakan untuk mengembangkan API ini.
- **Fiber Framework:** 🌐 Web framework minimalis dan cepat untuk Golang.
- 📦 **GORM:**  ORM (Object Relational Mapping) untuk Golang, memudahkan interaksi dengan database.
- **Redis:** ⚡ In-memory data structure store untuk cache dan manajemen sesi.
- 🔍 **Elasticsearch:**  Search engine distribusi untuk pencarian teks dan analitik data besar.
- **JWT:** 🔐 Untuk otentikasi dan otorisasi pengguna.
- 🗄️ **MySQL:** Sebagai database untuk menyimpan data.

## 📦 Struktur Proyek

```
├── configs/            # Database and Redis configuration
├── docs/               # Documentation API .yaml
├── pkg/                # Konfigurasi aplikasi
    ├── controllers/    # Logika untuk handling request
    ├── middleware/     # Middleware untuk otentikasi, logging, dll.
    ├── models/         # Struktur data dan interaksi database
    ├── response-code/ 
    ├── responses/ 
    ├── routes/         # Definisi rute API
    ├── utils/          # Fungsi utilitas yang digunakan di seluruh proyek
├── tests/              # Test endpoint menggunakan K6
├── temp/               # Temporary dari Air
├── .air.toml           # Konfigurasi Air
├── .env                # Environment
├── go.mod              # Go
├── go.sum              # Go sum
├── LICENSE             # License 
├── main.go             # Entry point aplikasi
└── README.md           # Documentation
```

## ⚙️ Cara Menjalankan Proyek  🚀

1. **Clone repository ini:**

   ```bash
   git clone https://github.com/Dziqha/TokoOnline.git

   cd TokoOnline
   ```

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

3. **Install Air (hot-reload tool):**

   ```bash
   go install github.com/cosmtrek/air@latest
   ```

   Pastikan $GOPATH/bin ada di dalam $PATH untuk menjalankan air.


4. **Buat file `.env` di root project:**

   Isi file `.env` dengan konfigurasi berikut:

   ```env
   DB_USER=your-username
   DB_PASSWORD=your-password
   DB_DB=your_db
   DB_HOST=your_host
   DB_PORT=your_port
   TOKEN_SECRET=your_secret
   REDIS_ADDR=your_address
   REDIS_PASSWORD=your-password
   REDIS_DB=your-db
   CACHE_KEY_INSERT_PRODUCT=your-key
   CACHE_KEY_PRODUCT_PREFIX=your-key
   CACHE_KEY_PRODUCT_ALL=your-key
   CACHE_KEY_ORDERS=your-key
   CACHE_KEY_ORDERS_ALL=your-key
   CACHE_KEY_ORDERS_PREFIX=your-key
   ```

5. **Jalankan aplikasi:**

   ```bash
   air
   ```

   API akan berjalan di `http://localhost:3000`. Gunakan Postman atau aplikasi serupa untuk mengakses endpoint API.


## 📚 Dokumentasi API

Dokumentasi lengkap API dapat diakses melalui [Swagger](http://localhost:3000/swagger/index.html) setelah server dijalankan.

## 🤝 Kontribusi

Kontribusi sangat terbuka! Silakan fork repo ini dan buat pull request, atau buka issue untuk diskusi lebih lanjut.

## 📝 Lisensi

Proyek ini dilisensikan di bawah lisensi MIT. Lihat file [LICENSE](LICENSE) untuk informasi lebih lanjut.

## 📌Catatan

Pastikan Redis dan Elasticsearch terinstal dan berjalan sebelum menjalankan aplikasi. Ikuti panduan instalasi resmi untuk kedua layanan sesuai dengan sistem operasi Anda. Untuk Windows, pastikan keduanya dapat diakses dari terminal atau command prompt. Jika menghadapi masalah, periksa dokumentasi resmi atau forum komunitas. 🚀

---
