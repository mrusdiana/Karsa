# PRD: Karsa: Job Portal untuk Talenta Tech Pemula

## Context

Task 5 (Lanjutan) dalam program magang Auspify: membangun platform web yang menghubungkan pencari kerja dengan perusahaan.
Karsa berfokus pada lowongan **tech dan digital** untuk level **magang, fresh graduate, dan junior**.

**Stack yang disetujui Auspify:** Go, MySQL, Vue.js, Docker, CI/CD.
**Cara kerja:** pemilik project menulis kode sendiri dengan panduan Claude. Claude tidak menyentuh GitHub (commit, push, merge).

**Syarat Auspify untuk Task 5:**
- Alur pengerjaan:
  1. desain dashboard pengguna
  2. fitur posting lowongan
  3. alur lamaran kerja
  4. pencarian dan filter
  5. fitur manajemen admin
- Fitur utama: daftar lowongan, lamaran kerja, pencarian dan filter, dashboard admin.
- Keterampilan yang dinilai: backend lanjutan, manajemen pengguna, optimasi database, arsitektur aplikasi web.

**Keputusan (24 Nov 2026):**

| Topik | Keputusan |
|---|---|
| Nama | **Karsa**. Dalam bahasa Indonesia, "karsa" berarti kehendak atau niat |
| Fokus | Lowongan tech dan digital (software, data, UI/UX, produk, QA, DevOps, digital marketing) untuk level Magang, Entry, Junior, dan Mid |
| Peran | Pencari kerja, perusahaan (wajib diverifikasi admin), dan admin |
| Fitur tambahan | Upload CV (PDF), notifikasi email lewat worker, simpan lowongan, grafik di dashboard |
| Fitur rekrutmen | Pertanyaan screening + skor kecocokan; Cari Talenta + undangan melamar (kandidat opt-in) |
| Backend | Go: Gin (router dan validasi), GORM (ORM), golang-migrate (migrasi SQL), JWT (token Bearer), bcrypt |
| Worker dan scheduler | Dua goroutine di dalam proses backend. Antrean email disimpan di tabel `email_jobs` |
| Email | SMTP ke **Mailpit** di Docker. Email bisa dilihat di inbox web `localhost:8025`, tanpa layanan berbayar |
| Frontend | Vue 3, Vite, Vue Router, Pinia, Tailwind CSS 4, axios |
| Bahasa dan mata uang | Antarmuka **bahasa Indonesia**, IDR |
| Desain | Terang dan bersih, satu warna aksen, profesional |
| CI/CD | GitHub Actions: test dan build di setiap push, lalu push image ke GitHub Container Registry di branch `main`. **Tanpa deploy** |
| Definisi selesai | Berjalan lokal dengan Docker Compose, CI hijau, README, screenshot, dan test |

**Catatan data:** semua perusahaan dan lowongan di data demo adalah fiktif. Footer aplikasi dan README memuat keterangan:
*"Proyek demo untuk pembelajaran. Perusahaan dan lowongan bersifat fiktif."*

---

## 1. Problem statement

Mahasiswa dan fresh graduate di bidang tech kesulitan menemukan pekerjaan pertama mereka:

- Lowongan magang dan entry-level tersebar di banyak tempat (Instagram, grup Telegram, situs perusahaan) dan sering tercampur dengan lowongan senior.
- Setelah melamar, kandidat jarang tahu kelanjutan lamarannya ("di-ghosting").
- Perusahaan kecil kewalahan menyaring puluhan CV yang sebagian besar tidak cocok dengan skill yang dibutuhkan.
- Kandidat yang bagus tapi belum melamar tidak bisa ditemukan perusahaan.

Yang dirugikan adalah **pencari kerja**, yang membuang waktu dan tidak mendapat kepastian, dan **perusahaan**, yang menyaring secara manual dan kehilangan kandidat potensial.

## 2. Target user dan persona

| Peran | Siapa | Akses |
|---|---|---|
| **Tamu** | Pengunjung tanpa akun | Beranda, cari dan filter lowongan, detail lowongan, profil perusahaan |
| **Pencari kerja** | Mahasiswa, fresh graduate, talenta junior | Profil dan CV, melamar, lamaran saya, lowongan tersimpan, undangan |
| **Perusahaan** | Rekruter atau pemilik startup | Profil perusahaan, kelola lowongan, pelamar beserta skor, Cari Talenta, dashboard |
| **Admin** | Pengelola Karsa | Dashboard platform, verifikasi perusahaan, pengguna, moderasi lowongan, kategori |

Persona berikut adalah ilustrasi, bukan hasil riset pengguna.

**Persona 1: Nadia, 21 tahun, mahasiswa Informatika semester 7**
- Mencari magang frontend untuk syarat kelulusan.
- Masalah: lowongan magang sulit dibedakan dari lowongan senior, dan ia tidak pernah tahu apakah lamarannya dibaca.
- Butuh: filter "Magang", status lamaran yang jelas, dan bisa ditemukan perusahaan.

**Persona 2: Bima, 29 tahun, Engineering Lead di startup 30 orang**
- Merekrut 2 junior backend tanpa tim HR.
- Masalah: 80 pelamar, sebagian besar tidak punya skill Go, dan menyaring CV satu per satu memakan waktu.
- Butuh: pertanyaan screening, pelamar terurut berdasarkan kecocokan, dan bisa mengundang kandidat duluan.

**Persona 3: Sari, admin Karsa**
- Butuh: memastikan hanya perusahaan asli yang memasang lowongan, memantau aktivitas platform, dan menurunkan lowongan yang melanggar.

## 3. Goals dan non-goals

**Goals**
1. Tamu dan pencari kerja dapat menemukan lowongan yang relevan lewat kata kunci, filter (kategori, tipe, level, mode kerja, lokasi, gaji), dan urutan.
2. Pencari kerja dapat melamar dengan CV dari profil dan menjawab pertanyaan screening, lalu memantau status lamaran beserta timeline-nya.
3. Pencari kerja mendapat email setiap kali status lamarannya berubah atau ada undangan baru.
4. Perusahaan yang sudah terverifikasi dapat memasang lowongan, memproses pelamar lewat tahapan status, dan melihat pelamar terurut berdasarkan skor kecocokan.
5. Perusahaan dapat menemukan kandidat opt-in lewat Cari Talenta dan mengirim undangan melamar tanpa melihat kontak kandidat.
6. Admin dapat memverifikasi perusahaan, mengelola pengguna, kategori, dan lowongan, serta melihat statistik platform.
7. Pencarian lowongan tetap cepat dengan indeks database yang tepat.
8. Aplikasi berjalan dengan Docker Compose, dan setiap push diperiksa otomatis oleh CI.

**Non-goals**
- Chat langsung antara perusahaan dan kandidat.
- Penjadwalan interview, integrasi kalender, dan video call.
- Parsing isi CV atau penilaian berbasis AI.
- Pembayaran atau paket berbayar untuk perusahaan.
- Login dengan Google atau LinkedIn.
- Deploy ke server publik dan aplikasi mobile native.

## 4. User stories

**Pencari kerja**
- **Akun:** Sebagai pencari kerja, saya ingin mendaftar, masuk, dan keluar agar lamaran saya tersimpan.
- **Profil:** Sebagai pencari kerja, saya ingin mengisi headline, pendidikan, tahun lulus, pengalaman, skill, dan mengunggah CV sekali saja.
- **Cari:** Sebagai pencari kerja, saya ingin mencari lowongan dengan kata kunci dan memfilter berdasarkan level (misalnya Magang) dan mode kerja (misalnya Remote).
- **Detail:** Sebagai pencari kerja, saya ingin melihat deskripsi, persyaratan, skill, gaji (bila ditampilkan), deadline, dan profil perusahaan.
- **Simpan:** Sebagai pencari kerja, saya ingin menyimpan lowongan untuk dilamar nanti.
- **Lamar:** Sebagai pencari kerja, saya ingin melamar dengan CV dari profil, surat pengantar opsional, dan jawaban screening.
- **Status:** Sebagai pencari kerja, saya ingin melihat status dan riwayat setiap lamaran, serta mendapat email saat statusnya berubah.
- **Tarik lamaran:** Sebagai pencari kerja, saya ingin menarik lamaran yang belum diproses.
- **Ditemukan:** Sebagai pencari kerja, saya ingin memilih apakah profil saya bisa ditemukan perusahaan.
- **Undangan:** Sebagai pencari kerja, saya ingin menerima atau menolak undangan melamar dari perusahaan.

**Perusahaan**
- **Akun dan profil:** Sebagai perusahaan, saya ingin mendaftar dan mengisi profil perusahaan agar bisa diverifikasi.
- **Lowongan:** Sebagai perusahaan, saya ingin membuat, mengubah, membuka, dan menutup lowongan, termasuk skill yang dibutuhkan, deadline, dan pertanyaan screening.
- **Pelamar:** Sebagai perusahaan, saya ingin melihat pelamar terurut berdasarkan skor kecocokan, memfilternya, membuka CV, dan memindahkan status lamaran.
- **Catatan:** Sebagai perusahaan, saya ingin menulis catatan internal per pelamar.
- **Cari Talenta:** Sebagai perusahaan, saya ingin mencari kandidat berdasarkan skill dan lokasi, lalu mengundang mereka melamar.
- **Dashboard:** Sebagai perusahaan, saya ingin melihat jumlah pelamar per status dan tren lamaran per minggu.

**Admin**
- **Verifikasi:** Sebagai admin, saya ingin menyetujui atau menolak perusahaan baru beserta alasannya.
- **Pengguna:** Sebagai admin, saya ingin mencari pengguna dan menonaktifkan akun yang bermasalah.
- **Moderasi:** Sebagai admin, saya ingin menutup atau menghapus lowongan yang melanggar.
- **Kategori:** Sebagai admin, saya ingin mengelola kategori lowongan.
- **Dashboard:** Sebagai admin, saya ingin melihat statistik dan grafik aktivitas platform.

## 5. Daftar fitur

| Fase | Fitur |
|---|---|
| **MVP** | Autentikasi 3 peran, profil pencari kerja + upload CV, profil perusahaan + verifikasi admin, CRUD lowongan + pertanyaan screening, cari dan filter lowongan, detail lowongan, simpan lowongan, melamar + skor kecocokan, pipeline status lamaran + timeline, Cari Talenta + undangan, notifikasi email lewat worker, scheduler penutup lowongan dan undangan kedaluwarsa, dashboard perusahaan dan admin dengan grafik, kelola pengguna, moderasi lowongan, kategori, Docker Compose, CI/CD |
| **v2** | Notifikasi di dalam aplikasi (ikon lonceng), lowongan serupa, ekspor pelamar ke CSV, reset password lewat email |
| **Nanti** | Chat kandidat dan perusahaan, jadwal interview, login sosial, deploy |

## 6. Functional requirements MVP

Aturan umum:
- Semua API memakai JSON di bawah `/api`, kecuali upload file yang memakai multipart.
- Endpoint selain publik mewajibkan token Bearer (JWT).
- Peran dicek di middleware:
  - peran yang salah → 403
  - data milik pengguna atau perusahaan lain → 404
  - akun nonaktif → 403 dengan pesan "Akun dinonaktifkan".
- Semua UI responsif.

### FR-1 Autentikasi dan peran
1. Daftar dengan nama, email unik, password (minimal 8 karakter), dan **peran**: `seeker` atau `company`. Perusahaan juga mengisi nama perusahaan saat mendaftar.
2. Akun admin tidak bisa dibuat lewat pendaftaran; admin dibuat oleh seeder.
3. Saat masuk, server mengembalikan token JWT (berlaku 7 hari) dan data pengguna. Password di-hash dengan bcrypt. Percobaan masuk dibatasi 5 kali per menit per IP.
4. Pesan gagal masuk tidak mengungkap apakah email terdaftar.
5. Setelah masuk, pengguna diarahkan sesuai peran: pencari kerja ke beranda, perusahaan ke dashboard perusahaan, admin ke dashboard admin.

### FR-2 Profil pencari kerja dan CV
1. Field profil: headline, nomor HP, lokasi (kota), pendidikan (SMA/SMK, D3, S1, S2), jurusan, tahun lulus, pengalaman (tahun, 0–30), skill (1–15 tag), bio, dan tautan portofolio (opsional).
2. **Upload CV:** hanya PDF, maksimal 2 MB, dicek dari isi file (magic bytes), tidak cukup dari ekstensinya. Mengganti CV menghapus file lama. CV yang sudah terlampir di lamaran tidak ikut terhapus, karena sudah disalin sebagai snapshot (FR-8.2).
3. **Profil lengkap** = headline, lokasi, pendidikan, tahun lulus, minimal 1 skill, dan CV. Profil yang belum lengkap tidak bisa dipakai melamar.
4. Toggle **"Bisa ditemukan perusahaan"** (`is_discoverable`, default mati) hanya bisa diaktifkan kalau profil sudah lengkap.
5. Skill dinormalisasi: spasi di awal dan akhir dibuang, perbandingan tidak peka huruf besar/kecil, dan disimpan dengan ejaan pertama yang dimasukkan.

### FR-3 Profil perusahaan dan verifikasi
1. Field: nama, logo (JPG/PNG/WebP, maksimal 1 MB), industri, ukuran (1–10, 11–50, 51–200, 201–500, >500), lokasi, website, dan deskripsi.
2. Status awal `pending`. Perusahaan yang `pending` atau `rejected` bisa mengisi profil dan membuat draft lowongan, tapi **tidak bisa membuka lowongan dan tidak bisa memakai Cari Talenta**.
3. Admin bisa **menyetujui** perusahaan (status jadi `verified`) atau **menolak** dengan alasan wajib (status jadi `rejected`). Kedua keputusan mengirim email ke perusahaan.
4. Perusahaan yang ditolak bisa memperbaiki profil lalu mengajukan ulang, dan statusnya kembali `pending`.
5. Halaman profil perusahaan publik (`/perusahaan/:slug`) hanya untuk perusahaan `verified`. Halaman ini menampilkan profil dan lowongan yang masih dibuka.

### FR-4 Kelola lowongan (perusahaan)
1. Field lowongan:
   - judul dan kategori
   - tipe: `full_time` | `part_time` | `internship` | `contract`
   - level: `internship` | `entry` | `junior` | `mid`
   - mode kerja: `onsite` | `hybrid` | `remote`
   - lokasi: kota, opsional kalau mode kerjanya remote
   - gaji min dan max dalam rupiah per bulan (opsional, max ≥ min), dan `show_salary`
   - skill: 1–10 tag
   - deskripsi dan persyaratan
   - deadline: hari ini sampai 90 hari ke depan
2. **Pertanyaan screening:** 0–5 per lowongan. Setiap pertanyaan berisi:
   - teks pertanyaan
   - tipe: `yes_no` | `number` | `text`
   - tanda wajib (`required`)
   - jawaban yang diharapkan:
     - `yes_no`: `yes` atau `no`
     - `number`: `>=N` atau `<=N`
     - `text`: tidak punya jawaban yang diharapkan; kalau `required`, jawabannya cukup tidak boleh kosong
3. Status lowongan: `draft` → `open` → `closed`. Hanya perusahaan `verified` yang bisa membuka lowongan. Lowongan yang ditutup bisa dibuka lagi, selama deadline-nya masih di masa depan.
4. Lowongan yang sudah punya pelamar tidak bisa dihapus (409). Perusahaan diarahkan untuk menutupnya.
5. Kalau skill atau pertanyaan screening diubah setelah ada pelamar, perubahan hanya berlaku untuk lamaran baru. Skor dan jawaban lamaran lama tidak dihitung ulang.
6. Slug lowongan berbentuk `judul-kode` (misalnya `frontend-intern-7q3k`) dan unik.

### FR-5 Cari dan filter lowongan (publik)
1. Hanya lowongan `open` dengan deadline yang belum lewat dan perusahaan yang `verified` serta aktif.
2. Kata kunci mencari di judul, deskripsi, skill, dan nama perusahaan. Pencarian judul dan deskripsi memakai indeks FULLTEXT MySQL.
3. Filter (bisa digabung): kategori, tipe, level (bisa pilih lebih dari satu), mode kerja (bisa pilih lebih dari satu), lokasi, dan gaji minimal.
4. Urutan: terbaru (default), gaji tertinggi, dan deadline terdekat.
5. Paginasi 12 lowongan per halaman. Filter disimpan di query string URL, jadi tautannya bisa dibagikan dan tombol back browser tetap jalan.
6. Hasil kosong menampilkan empty state dengan tombol "Hapus filter".

### FR-6 Beranda dan detail lowongan
1. Beranda menampilkan:
   - hero dengan kolom pencarian
   - chip kategori beserta jumlah lowongannya
   - 6 lowongan terbaru
   - bagian "Magang & Fresh Graduate" (lowongan level `internship` dan `entry`)
   - perusahaan unggulan, yaitu perusahaan terverifikasi dengan lowongan terbuka terbanyak
2. Detail lowongan menampilkan:
   - judul, perusahaan (logo dan tautan ke profilnya), lokasi, mode kerja, tipe, level
   - gaji, kalau `show_salary` aktif; kalau tidak, ditulis "Gaji dirahasiakan"
   - skill, deskripsi, persyaratan, deadline, jumlah pelamar, dan tanggal dipasang
3. Setiap kali detail dibuka, `views` bertambah 1. Kunjungan berulang dari pengguna yang sama tidak dibatasi, karena ini hanya statistik kasar.
4. Untuk pencari kerja yang sudah masuk:
   - tombol **Lamar** (atau "Sudah dilamar" beserta statusnya)
   - tombol **Simpan** atau **Tersimpan**
   - perkiraan skor kecocokan
5. Tamu yang menekan Lamar atau Simpan diarahkan ke halaman masuk, lalu kembali ke lowongan itu setelah masuk.

### FR-7 Simpan lowongan
1. Pencari kerja bisa menyimpan dan menghapus lowongan. Menyimpan lowongan yang sudah tersimpan tidak membuat data ganda (idempoten).
2. Halaman "Lowongan Tersimpan" menandai lowongan yang sudah ditutup atau lewat deadline, dan tombol Lamar-nya dinonaktifkan.

### FR-8 Melamar dan skor kecocokan
1. Syarat melamar:
   - peran pencari kerja dan profil lengkap (FR-2.3)
   - lowongan `open` dan deadline belum lewat
   - belum pernah melamar lowongan itu (UNIQUE `job_id` + `user_id`)
   - lamaran yang pernah **ditarik** juga tidak bisa diajukan ulang
2. Isi lamaran:
   - surat pengantar opsional, maksimal 2.000 karakter
   - jawaban untuk setiap pertanyaan screening; pertanyaan `required` wajib dijawab
   - **snapshot CV**: file CV disalin, sehingga lamaran tetap memakai CV saat melamar
3. **Skor kecocokan (0–100)** dihitung saat melamar, lalu disimpan. Komponennya:
   - **Skill (60 poin):** 60 × (jumlah skill lowongan yang dimiliki pelamar ÷ jumlah skill lowongan). Perbandingan tidak peka huruf besar/kecil.
   - **Level (20 poin):**

     | Level lowongan | Pelamar yang cocok |
     |---|---|
     | `internship` | Masih kuliah (tahun lulus ≥ tahun ini) atau lulus ≤ 1 tahun |
     | `entry` | Pengalaman 0–1 tahun |
     | `junior` | Pengalaman 1–3 tahun |
     | `mid` | Pengalaman 3–6 tahun |

     Cocok dengan rentangnya = 20 poin. Selisih satu tingkat = 10 poin. Selain itu 0.
   - **Screening (20 poin):** 20 × (jumlah jawaban lolos ÷ jumlah pertanyaan). Lowongan tanpa pertanyaan screening mendapat 20.
4. `meets_requirements` bernilai `false` kalau ada pertanyaan `required` yang jawabannya tidak lolos. Lamaran seperti ini diberi label **"Tidak memenuhi syarat"**, tapi tidak otomatis ditolak.
5. Rincian skor (skill yang cocok dan yang kurang, poin level, dan jawaban yang lolos) disimpan bersama lamaran dan ditampilkan ke perusahaan.
6. Fungsi penghitung skor adalah fungsi murni, dan wajib punya unit test untuk semua cabang di atas.
7. Setelah melamar, sistem mencatat event `applied` dan mengirim email konfirmasi ke pelamar.

### FR-9 Pipeline status lamaran
1. Status yang ada: `applied` (Dilamar), `reviewed` (Ditinjau), `interview` (Interview), `accepted` (Diterima), `rejected` (Ditolak), dan `withdrawn` (Ditarik).
2. Perpindahan status yang diizinkan:
   - Perusahaan: `applied` → `reviewed` | `rejected`; `reviewed` → `interview` | `rejected`; `interview` → `accepted` | `rejected`.
   - Pencari kerja: `applied` → `withdrawn`.
   - `accepted`, `rejected`, dan `withdrawn` bersifat final.
   - Perpindahan lain ditolak dengan 422.
3. Setiap perpindahan mencatat `application_events` (status asal, status tujuan, siapa yang mengubah, dan waktu), yang ditampilkan sebagai timeline.
4. Setiap perpindahan yang dilakukan perusahaan **mengantrikan email** ke pelamar (FR-13).
5. Perusahaan bisa menulis catatan internal (maksimal 1.000 karakter) yang tidak pernah terlihat oleh pelamar.
6. Kalau dua rekruter mengubah status bersamaan, hanya satu perubahan yang berhasil. Ini dijaga dengan pengecekan status asal di klausa `WHERE` saat update. Yang kalah mendapat 409 "Status sudah berubah, muat ulang halaman."

### FR-10 Pelamar per lowongan (perusahaan)
1. Daftar pelamar menampilkan nama, headline, skor, label syarat, sumber (Melamar / Diundang), status, dan tanggal melamar.
2. Bisa difilter berdasarkan status, skor minimal, skill, `meets_requirements`, dan sumber. Bisa diurutkan berdasarkan skor tertinggi (default) atau terbaru. Paginasi 20 per halaman.
3. Detail pelamar menampilkan profil lengkap (termasuk email dan nomor HP), rincian skor, jawaban screening, surat pengantar, timeline, catatan, dan tombol unduh CV.
4. **Unduh CV** lewat endpoint terlindungi. File hanya bisa diakses pemilik CV, perusahaan pemilik lowongan, dan admin. File CV tidak disajikan dari folder publik.

### FR-11 Cari Talenta dan undangan
1. Hanya untuk perusahaan `verified`. Yang bisa dicari hanya pencari kerja yang aktif, profilnya lengkap, dan `is_discoverable = true`.
2. Filter: skill (bisa lebih dari satu; semua skill yang dipilih harus dimiliki), lokasi, pendidikan, rentang tahun lulus, dan pengalaman maksimal. Paginasi 20 per halaman.
3. Data yang dikirim ke perusahaan hanya nama, headline, skill, lokasi, pendidikan, jurusan, tahun lulus, dan pengalaman. **Email, nomor HP, dan CV tidak pernah dikirim** dari endpoint ini.
4. **Undangan melamar:**
   - Perusahaan memilih salah satu lowongan `open` miliknya dan menulis pesan opsional (maksimal 500 karakter).
   - Batasan:
     - satu undangan per kandidat per lowongan (UNIQUE)
     - kandidat yang sudah melamar lowongan itu tidak bisa diundang
     - maksimal 20 undangan per perusahaan per hari (429)
   - Undangan baru mengantrikan email ke kandidat.
5. Kandidat bisa **menerima** atau **menolak** undangan:
   - **Menerima** membuka form lamaran untuk lowongan itu (surat pengantar dan jawaban screening). Saat dikirim, lamaran dibuat dengan `source = invited`, status undangan menjadi `accepted`, dan **kontak kandidat terlihat oleh perusahaan lewat lamaran itu**.
   - Lamaran dari undangan tetap butuh profil lengkap, dan skornya dihitung dengan aturan yang sama (FR-8).
   - **Menolak** mengubah status undangan menjadi `declined`, dan statusnya terlihat di halaman "Undangan Terkirim".
6. Undangan `pending` berubah menjadi `expired` saat lowongannya ditutup atau lewat deadline, atau 14 hari setelah dikirim (FR-14).

### FR-12 Dashboard
1. **Dashboard perusahaan:**
   - jumlah lowongan terbuka, total pelamar, pelamar baru 7 hari terakhir, dan undangan pending
   - jumlah pelamar per status
   - grafik lamaran per minggu (8 minggu terakhir)
   - 5 lowongan dengan pelamar terbanyak
   - banner status verifikasi, kalau perusahaan belum `verified`
2. **Dashboard admin:**
   - total pengguna (per peran), perusahaan (per status), lowongan terbuka, dan lamaran
   - grafik lamaran per minggu (8 minggu)
   - grafik lowongan per kategori (bar horizontal)
   - 5 perusahaan teratas berdasarkan jumlah lamaran
   - jumlah perusahaan yang menunggu verifikasi, dengan tautan langsung ke halaman verifikasi
3. Angka dashboard dihitung dengan query agregat (`COUNT`, `GROUP BY`), tidak dengan memuat seluruh baris data ke memori.

### FR-13 Notifikasi email (worker)
1. Email tidak dikirim langsung di dalam request. Request hanya menulis baris `email_jobs` (status `pending`) di transaksi database yang sama dengan perubahan datanya.
2. **Worker** berjalan sebagai goroutine yang mengambil job `pending` setiap 5 detik:
   - Job diambil dengan `SELECT ... FOR UPDATE SKIP LOCKED`, supaya tidak pernah diproses dua kali.
   - Kalau berhasil, status job menjadi `sent`.
   - Kalau gagal, `attempts` bertambah dan `available_at` diundur 1, 5, lalu 15 menit. Setelah 3 kali gagal, statusnya menjadi `failed` beserta `last_error`.
3. Email yang dikirim:

   | Pemicu | Penerima |
   |---|---|
   | Lamaran terkirim (konfirmasi) | Pelamar |
   | Status lamaran berubah (Ditinjau / Interview / Diterima / Ditolak) | Pelamar |
   | Undangan melamar baru | Kandidat |
   | Perusahaan disetujui atau ditolak | Perusahaan |

4. Template email berupa HTML sederhana berbahasa Indonesia, ditambah versi teks, dan selalu berisi tautan ke halaman terkait di Karsa.
5. Di Docker, email dikirim ke Mailpit dan bisa dilihat di `localhost:8025`. Konfigurasi SMTP diambil dari environment variable.

### FR-14 Scheduler
1. **Scheduler** berjalan sebagai goroutine dengan `time.Ticker` yang jalan setiap 1 jam, dan juga sekali saat server baru menyala:
   - lowongan `open` yang deadline-nya sudah lewat diubah menjadi `closed`
   - undangan `pending` diubah menjadi `expired` kalau lowongannya `closed` atau sudah lebih dari 14 hari
2. Setiap operasi berupa satu `UPDATE` massal yang aman dijalankan berulang (idempoten).
3. Worker dan scheduler berhenti dengan rapi saat server dimatikan (SIGTERM), memakai `context` dan graceful shutdown. Email yang sedang dikirim tidak terpotong di tengah jalan.

### FR-15 Admin: pengguna, lowongan, kategori
1. **Pengguna:** daftar dengan pencarian (nama/email), filter peran, dan paginasi. Admin bisa menonaktifkan atau mengaktifkan akun.
   - Akun nonaktif tidak bisa masuk, dan token yang sudah ada ditolak.
   - Lowongan milik perusahaan yang nonaktif hilang dari pencarian publik.
   - Admin tidak bisa menonaktifkan dirinya sendiri.
2. **Lowongan:** daftar semua lowongan dengan filter status dan perusahaan. Admin bisa menutup lowongan (alasan wajib, disimpan di `moderation_note`) atau menghapusnya. Lowongan yang sudah punya lamaran hanya bisa ditutup.
3. **Kategori:** CRUD nama dan slug. Kategori yang masih dipakai lowongan tidak bisa dihapus (409).

### FR-16 Docker
1. `docker compose up --build` menjalankan `mysql`, `mailpit`, `backend`, dan `frontend` tanpa perlu memasang Go atau MySQL di komputer.
2. Healthcheck MySQL memakai `mysqladmin ping -h 127.0.0.1`. Backend baru jalan setelah MySQL `healthy`.
3. Saat start, backend menjalankan migrasi. Seed hanya dijalankan kalau `RUN_SEED=true` dan tabel `users` masih kosong.
4. Data MySQL dan file upload (CV dan logo) tetap ada walaupun container dimatikan, karena disimpan di volume.
5. Image backend dibuat multi-stage: tahap build dengan `golang`, lalu hasilnya dijalankan di `alpine` sebagai user non-root.

### FR-17 CI/CD
1. Workflow `.github/workflows/ci.yml` berjalan di setiap push dan pull request.
2. Job **backend:** `go vet ./...` lalu `go test ./...` dengan service MySQL 8.4 di GitHub Actions.
3. Job **frontend:** `npm ci`, `npm run lint`, lalu `npm run build`.
4. Job **docker:** hanya berjalan di push ke `main`, setelah dua job di atas lolos. Job ini membangun image backend dan frontend, lalu push ke `ghcr.io/<owner>/karsa-backend` dan `karsa-frontend` dengan tag `latest` dan SHA commit. Login memakai `GITHUB_TOKEN`, tanpa secret tambahan.
5. Badge status CI ditampilkan di README.

## 7. Sketsa data model

| Entitas | Field kunci |
|---|---|
| **users** | id, name, email (unik), password_hash, role (`seeker` / `company` / `admin`), is_active, timestamps |
| **seeker_profiles** | user_id (PK), headline, phone, location, education, major, graduation_year, experience_years, skills (JSON), bio, portfolio_url, cv_path, is_discoverable, timestamps |
| **companies** | id, user_id (unik), name, slug (unik), logo_path, industry, size, location, website, description, status (`pending` / `verified` / `rejected`), rejection_reason, verified_at, timestamps |
| **categories** | id, name, slug (unik) |
| **jobs** | id, company_id, category_id, title, slug (unik), type, level, work_mode, location, salary_min, salary_max, show_salary, skills (JSON), description, requirements, deadline, status (`draft` / `open` / `closed`), moderation_note, views, timestamps |
| **screening_questions** | id, job_id, question, type (`yes_no` / `number` / `text`), required, expected_answer, position |
| **applications** | id, job_id, user_id, cv_path (snapshot), cover_letter, status, source (`applied` / `invited`), match_score, score_breakdown (JSON), meets_requirements, company_note, timestamps; unik (job_id, user_id) |
| **application_answers** | id, application_id, question_id, answer, passed |
| **application_events** | id, application_id, from_status (nullable), to_status, changed_by, created_at |
| **saved_jobs** | user_id, job_id, created_at; PK (user_id, job_id) |
| **invitations** | id, company_id, job_id, user_id, message, status (`pending` / `accepted` / `declined` / `expired`), timestamps; unik (job_id, user_id) |
| **email_jobs** | id, to_email, subject, template, payload (JSON), status (`pending` / `sent` / `failed`), attempts, available_at, sent_at, last_error, created_at |

**Indeks (optimasi database):**

| Tabel | Indeks | Untuk |
|---|---|---|
| jobs | (status, deadline, created_at) | Pencarian publik + urutan terbaru |
| jobs | (category_id), (company_id), (level), (work_mode) | Filter |
| jobs | FULLTEXT (title, description) | Kata kunci |
| applications | (job_id, match_score), (user_id, created_at), (status) | Daftar pelamar terurut skor, Lamaran Saya, dashboard |
| invitations | (user_id, status), (company_id, created_at) | Halaman undangan, batas 20 per hari |
| email_jobs | (status, available_at) | Worker mengambil job |
| seeker_profiles | (is_discoverable, location) | Cari Talenta |

- Gaji disimpan sebagai `INT UNSIGNED` dalam rupiah per bulan.
- Kolom JSON `skills` difilter dengan `JSON_CONTAINS`. Untuk ukuran data demo cara ini cukup cepat; kalau datanya membesar, skill dipindah ke tabel `job_skills` (dicatat sebagai keterbatasan).

**Status lamaran:**

```
applied (Dilamar) ──► reviewed (Ditinjau) ──► interview (Interview) ──► accepted (Diterima)
   │                       │                         │
   ├──► withdrawn          └──► rejected ◄───────────┘
   │   (oleh pelamar)            (Ditolak)
   └──► rejected
```

## 8. Edge case dan failure state

| Kasus | Perilaku |
|---|---|
| Melamar dengan profil atau CV belum lengkap | 422 berisi daftar field yang kurang; UI mengarahkan ke halaman profil |
| Melamar dua kali / tombol Lamar ditekan dua kali | Hanya satu lamaran (UNIQUE); permintaan kedua → 409 "Kamu sudah melamar lowongan ini" |
| Melamar setelah deadline lewat atau lowongan ditutup | 422 "Lowongan sudah ditutup" |
| Pertanyaan `required` tidak dijawab | 422 dengan galat per pertanyaan |
| Jawaban angka bukan angka | 422 |
| Upload CV bukan PDF (walaupun ekstensinya `.pdf`) atau > 2 MB | Ditolak dengan pesan |
| Perusahaan `pending` mencoba membuka lowongan atau memakai Cari Talenta | 403 "Perusahaan belum terverifikasi" |
| Perusahaan membuka pelamar atau lowongan perusahaan lain | 404 |
| Pencari kerja membuka endpoint perusahaan atau admin | 403 |
| Perpindahan status tidak sah (misalnya `applied` → `accepted`) | 422 |
| Dua rekruter mengubah status bersamaan | Satu berhasil, yang lain 409 |
| Mengundang kandidat yang sudah melamar atau sudah diundang | 409 |
| Lebih dari 20 undangan per hari | 429 |
| Kandidat mematikan "Bisa ditemukan" | Langsung hilang dari Cari Talenta; undangan yang sudah ada tetap berlaku |
| Lowongan ditutup saat masih ada undangan pending | Undangan menjadi `expired` |
| Menerima undangan yang sudah kedaluwarsa | 422 "Undangan sudah tidak berlaku" |
| Mengubah skill atau pertanyaan setelah ada pelamar | Hanya berlaku untuk lamaran baru (FR-4.5) |
| Menghapus lowongan yang punya pelamar / kategori yang dipakai | 409, disarankan ditutup |
| SMTP (Mailpit) mati | Request tetap sukses; email dicoba ulang oleh worker sampai 3 kali, lalu `failed` |
| Server restart saat ada email pending | Email tetap terkirim setelah server menyala, karena antreannya tersimpan di database |
| Akun dinonaktifkan saat sedang masuk | Request berikutnya → 403, frontend keluar dan kembali ke halaman masuk |
| Token kedaluwarsa | 401, ke halaman masuk, lalu kembali ke halaman semula |
| Pencarian dengan karakter khusus (`%`, `'`) | Aman karena memakai query berparameter (tanpa SQL injection) |
| Gaji max lebih kecil dari min | 422 di field gaji max |
| Pencarian tanpa hasil | Empty state + "Hapus filter" |
| Jaringan atau server gagal | Pesan galat + tombol coba lagi |

## 9. Success metrics

Keberhasilan diukur dengan pengujian. Target kinerja berikut adalah usulan.

**Kelengkapan**
- 100% FR MVP lolos uji manual, dan logika intinya punya test otomatis.
- Tiga alur end-to-end berhasil:
  1. **Melamar:** pencari kerja daftar → lengkapi profil dan CV → cari dengan filter Magang → lamar dan jawab screening → perusahaan memindahkan status sampai Diterima → email setiap perubahan muncul di Mailpit.
  2. **Verifikasi dan posting:** perusahaan daftar → admin menyetujui → perusahaan membuka lowongan → lowongan muncul di pencarian.
  3. **Undangan:** perusahaan mencari talenta → mengundang → kandidat menerima → lamaran `invited` muncul dengan skor.
- `docker compose up --build` dari kondisi bersih langsung menjalankan aplikasi beserta data demo.
- Workflow CI hijau di GitHub, dan image ada di GitHub Container Registry.
- README, screenshot, dan daftar teknologi lengkap.

**Akurasi dan keamanan**
- 0 kasus lamaran ganda pada uji klik ganda atau permintaan bersamaan.
- 0 kasus email terkirim dua kali untuk satu job.
- 0 kasus kontak atau CV kandidat bocor lewat Cari Talenta (dicek dengan test respons endpoint).
- 0 kasus akses data lintas pengguna atau perusahaan (dicek dengan test 403/404).
- Skor kecocokan sesuai rumus FR-8.3 di semua kasus unit test.

**Kinerja (usulan)**
- Pencarian lowongan dengan 1.000 lowongan uji merespons < 300 ms di lokal. Rencana query-nya (`EXPLAIN`) memakai indeks.
- Tidak ada scroll horizontal pada lebar 375 px dan 1440 px.

## 10. Open questions

**Sudah dijawab (24 Nov 2026)**
- **OQ-1, Nama:** Karsa.
- **OQ-2, Fokus:** tech dan digital, termasuk magang dan fresh graduate.
- **OQ-3, Peran:** 3 peran, dengan verifikasi perusahaan oleh admin.
- **OQ-4, Deploy:** tidak; cukup Docker + CI dan push image.
- **OQ-5, Fitur rekrutmen:** screening + skor kecocokan, serta Cari Talenta + undangan.

**Bisa diputuskan selama implementasi**
- **OQ-6, Tenggat:** kapan Task 5 harus dikumpulkan?
- **OQ-7, Daftar kategori awal:** Software Engineering, Data, UI/UX Design, Product, QA, DevOps & Cloud, Digital Marketing, dan perlu tambahan lain?
- **OQ-8, Daftar kota:** input bebas, atau pilihan kota besar (Jakarta, Bandung, Surabaya, Yogyakarta, Semarang, Medan, Bali, Remote)?
- **OQ-9, Saran skill:** apakah input skill menampilkan saran dari skill yang sudah ada di database?
- **OQ-10, Struktur kode Go:** per fitur (`internal/jobs`, `internal/applications`) atau per lapisan (`handlers`, `services`, `repositories`)?

## Verifikasi PRD terhadap syarat Auspify

| Syarat Auspify | Tercakup di |
|---|---|
| Alur 1: desain dashboard pengguna | FR-12, desain di DESIGN.md |
| Alur 2: fitur posting lowongan | FR-3, FR-4 |
| Alur 3: alur lamaran kerja | FR-2, FR-8, FR-9, FR-10, FR-11 |
| Alur 4: pencarian dan filter | FR-5 (ditambah filter pelamar di FR-10 dan Cari Talenta di FR-11) |
| Alur 5: fitur manajemen admin | FR-3.3, FR-12.2, FR-15 |
| Fitur: daftar lowongan | FR-5, FR-6 |
| Fitur: lamaran kerja | FR-8, FR-9 |
| Fitur: pencarian & filter | FR-5 |
| Fitur: dashboard admin | FR-12.2, FR-15 |
| Skill: backend lanjutan | Worker dan scheduler goroutine (FR-13, FR-14), transaksi, `SKIP LOCKED`, graceful shutdown |
| Skill: manajemen pengguna | FR-1 (3 peran + JWT), FR-3 (verifikasi), FR-15.1 (nonaktifkan akun) |
| Skill: optimasi database | Indeks di bagian 7, FULLTEXT, query agregat (FR-12.3), skor disimpan agar bisa diurutkan lewat indeks |
| Skill: arsitektur aplikasi web | Arsitektur di DESIGN.md, Docker (FR-16), CI/CD (FR-17) |
| Stack: Go, MySQL, Vue.js, Docker, CI/CD | Context, FR-16, FR-17 |
