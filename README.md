# 📝 Laporan Tugas 2 (Database, Relasi & layered architecture)

Di task kali ini, aku belajar menghubungkan kodingan Golang ke Database (PostgreSQL) dan membuat relasi antar tabel serta belajar cara membuat layered architecture dengan benar.

Fitur paling baru di sini adalah **Input Category**. Jadi sekarang kita bisa menambah jenis kategori barang langsung dari aplikasi, nggak perlu insert manual lewat database lagi! 🎉

## ✅ Yang Berhasil Dikerjakan

Berikut adalah checklist fitur yang sudah berhasil aku selesaikan:

- [x] **Koneksi Database:** Berhasil menyambungkan Golang dengan PostgreSQL.
- [x] **Input Kategori Baru, Update & Delete:** Bisa tambah, update & delete kategori (misal: "Makanan", "Minuman") langsung via POST.
- [x] **Menampilkan Kategori:** Berhasil menampilkan semua category dari table database
- [x] **Input Produk Relasi:** Saat tambah produk, kita tinggal masukkan ID dari kategori yang sudah kita buat tadi.
- [x] **JOIN Query:** Menampilkan data produk lengkap dengan nama kategorinya (hasil gabungan dua tabel).
- [x] **Validasi:** Kalau input ID kategori asal-asalan, program bakal nolak (aman dari error).

---

## 🛠️ Apa yang Dipelajari?

### 1. Fitur Tambah Kategori
Sekarang database `category_product` bisa diisi lewat API.
**Contoh Input Kategori (POST):**
```json
{
    "name": "Buah-Buahan"
}