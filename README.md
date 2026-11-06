# wa-bot-go 🤖

WhatsApp bot modern yang dibangun menggunakan bahasa pemrograman Go dengan library [whatsmeow](https://github.com/tulir/whatsmeow).

## ✨ Fitur

- **Auto-Reply**: Balasan otomatis untuk pesan masuk dengan sistem cooldown harian per chat.
- **Admin System**: Perintah khusus yang hanya bisa dijalankan oleh admin/owner yang terdaftar di database.
- **Sticker Maker**: Konversi Foto, GIF, dan Video menjadi stiker WhatsApp secara instan.
    - **Rounded Edges**: Stiker otomatis memiliki sudut membulat yang estetik.
    - **Square Crop**: Otomatis melakukan crop tengah agar stiker proporsional.
- **Command Router**: Sistem perintah yang rapi dengan dukungan prefix `.` dan `/`.

## 🛠️ Persyaratan

- **Go** (versi terbaru disarankan)
- **FFmpeg**: Diperlukan untuk pemrosesan video dan pembuatan stiker rounded.
- **SQLite**: Digunakan untuk penyimpanan data session dan pengaturan.

## 🚀 Cara Menjalankan

1. Clone repositori ini.
2. Pastikan FFmpeg sudah terinstal di sistem Anda.
3. Jalankan perintah:
   ```bash
   go run cmd/bot/main.go
   ```
4. Scan QR code yang muncul di terminal menggunakan aplikasi WhatsApp Anda.

## 📜 Perintah yang Tersedia

### Publik (Prefix `.`)
- `.about` - Informasi tentang bot.
- `.help` - Menampilkan daftar perintah yang tersedia.
- `.sticker` atau `.s` - Mengubah foto/video/gif (langsung atau reply) menjadi stiker.

### Admin (Prefix `/`)
- `/autoreply on/off` - Mengaktifkan atau menonaktifkan fitur auto-reply.
- `/autoreply status` - Melihat status aktif auto-reply saat ini.

---
Co-Authored-By: Claude Code <noreply@anthropic.com>
