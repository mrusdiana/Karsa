# Kontrak Desain: Karsa (Job Portal)

Status: **DRAF v1, menunggu persetujuan.** Coding UI dimulai setelah kontrak ini disetujui.
Acuan fitur: `PRD.md`. Jika kontrak ini dan PRD berbeda, PRD yang berlaku, lalu kontrak ini diperbarui.

Aturan perubahan: setiap perubahan tampilan, alur, atau bentuk API dicatat di sini **sebelum** kodenya diubah.

Antarmuka berbahasa **Indonesia**. Kode, nama field API, dan nilai enum memakai bahasa Inggris.

---

## 1. Prinsip

1. **Terang, tenang, profesional.**
   - Latar putih hangat, satu warna aksen hijau teal, banyak ruang kosong.
   - Tanpa gradien mencolok, glow, atau ilustrasi berlebihan.
   - Kesan yang dituju: tempat yang serius untuk mencari kerja pertama.
2. **Informasi lowongan bisa dipindai dalam 3 detik.** Kartu lowongan selalu menampilkan urutan yang sama: logo, judul, perusahaan, lalu chip level, mode kerja, dan lokasi, lalu gaji, lalu skill.
3. **Status selalu jelas.** Status lamaran, verifikasi, dan undangan memakai teks + ikon + warna, bukan warna saja.
4. **Satu tindakan utama per layar.** Contoh: "Lamar sekarang", "Buka lowongan", "Kirim undangan". Tindakan lain memakai tombol sekunder atau ghost.
5. **Privasi terlihat.** Di Cari Talenta, kontak kandidat ditandai jelas sebagai disembunyikan. Di profil pencari kerja, toggle "Bisa ditemukan perusahaan" diberi penjelasan singkat tentang data apa saja yang terlihat.
6. **Skor kecocokan harus bisa dijelaskan.** Angka skor selalu bisa dibuka untuk melihat rinciannya (skill cocok/kurang, level, screening). Skor tidak pernah tampil sebagai angka tanpa penjelasan.

---

## 2. Design token

Semua warna, ukuran, dan jarak di komponen **wajib** memakai token. Token didefinisikan di `@theme` Tailwind 4.

### 2.1 Warna

| Token | Nilai | Dipakai untuk |
|---|---|---|
| `--bg` | `#FAFAF7` | Latar halaman (putih hangat) |
| `--surface` | `#FFFFFF` | Kartu, panel, modal, header |
| `--surface-2` | `#F3F4F1` | Input, chip nonaktif, baris hover tabel, latar sidebar dashboard |
| `--border` | `#E4E4DF` | Border kartu dan input, pemisah |
| `--border-strong` | `#CFCFC8` | Border saat hover |
| `--text` | `#14171F` | Teks utama |
| `--muted` | `#5F6672` | Teks sekunder, metadata (kontras 5.6:1 di atas `--surface`) |
| `--primary` | `#0F766E` | Teal 700: tombol utama, link, tab aktif, chip terpilih |
| `--primary-hover` | `#115E59` | Hover tombol utama |
| `--primary-soft` | `#E6F2F0` | Latar chip/badge aksen, highlight baris, kotak info |
| `--on-primary` | `#FFFFFF` | Teks di atas `--primary` (kontras 5.5:1) |
| `--info` | `#2563EB` | Status Ditinjau, sumber "Diundang" |
| `--success` | `#15803D` | Diterima, Terverifikasi, skor tinggi |
| `--warning` | `#B45309` | Interview, Menunggu verifikasi, deadline dekat, skor sedang |
| `--danger` | `#B91C1C` | Galat, Ditolak, Tidak memenuhi syarat, hapus |
| `--neutral` | `#6B7280` | Ditarik, Ditutup, Kedaluwarsa, Draft |
| `--focus` | `0 0 0 3px rgba(15, 118, 110, 0.35)` | Cincin fokus |
| `--shadow-card` | `0 1px 2px rgba(20, 23, 31, 0.04), 0 4px 12px -4px rgba(20, 23, 31, 0.08)` | Kartu saat hover, dropdown, modal |

- MVP hanya memakai **tema terang**.
- Badge status memakai latar warna status dengan opasitas 10% dan teks warna status penuh. Contoh: Diterima memakai `bg-success/10 text-success`.
- Semua pasangan teks dan latar wajib kontras ≥ 4.5:1.

### 2.2 Tipografi

| Peran | Font |
|---|---|
| Heading | **Manrope** (700, 800) |
| Teks antarmuka | **Inter** (400, 500, 600) |
| Angka (gaji, skor, statistik) | Inter dengan `tabular-nums` |

Cadangan: `system-ui, sans-serif`.

| Token | Ukuran / bobot | Dipakai untuk |
|---|---|---|
| `display` | 40px Manrope 800, tracking -0.02em (mobile 28px) | Judul hero beranda |
| `h1` | 28px Manrope 700 (mobile 24px) | Judul halaman, judul lowongan di detail |
| `h2` | 20px Manrope 700 | Judul bagian dan kartu dashboard |
| `h3` | 16px Inter 600 | Judul kartu lowongan |
| `body` | 15px Inter 400, line-height 1.6 | Teks umum, deskripsi lowongan |
| `small` | 13px Inter 500 | Metadata, chip, label |
| `stat` | 30px Manrope 800 | Angka di StatCard |

- Deskripsi dan persyaratan lowongan dibatasi lebar maksimal 68ch, supaya nyaman dibaca.
- Gaji tidak terpotong ke baris baru (`white-space: nowrap`).

### 2.3 Jarak, bentuk, dan gerak

| Token | Nilai |
|---|---|
| Skala jarak | 4, 8, 12, 16, 24, 32, 48, 64px |
| Radius | 10px tombol/input; 14px kartu/panel; 9999px chip dan badge; 12px logo perusahaan |
| Logo perusahaan | 1:1. Ukuran 40px di kartu, 64px di detail, 96px di profil perusahaan. Kalau logo kosong, tampil inisial di atas `--primary-soft` |
| Transisi | 150–200ms ease-out untuk warna, border, dan bayangan. Kartu yang di-hover mendapat `--shadow-card` + `--border-strong`, tanpa bergeser |
| Gerak | Mati saat `prefers-reduced-motion` |
| Tinggi elemen interaktif | Minimal 44px |

### 2.4 Breakpoint

| Nama | Lebar | Daftar lowongan | Filter |
|---|---|---|---|
| Mobile | < 640px | 1 kolom | Tombol "Filter (n)" membuka drawer bawah |
| Tablet | 640–1023px | 1 kolom | Drawer |
| Desktop | ≥ 1024px | Sidebar filter 280px + 1 kolom kartu; konten maks. 1200px | Sidebar sticky |

Wajib diuji pada lebar 375, 768, 1024, dan 1440px, tanpa scroll horizontal.

---

## 3. Komponen

| Komponen | State wajib | Aturan |
|---|---|---|
| **Button** (`primary`, `secondary`, `ghost`, `danger`, `icon`) | default, hover, focus, disabled, loading | `primary`: `--primary` / `--on-primary`. `secondary`: putih + border `--border`. Saat loading, teks berganti ("Mengirim…", "Menyimpan…") dan tombol nonaktif. Tombol ikon wajib punya `aria-label`. |
| **Field** (text, email, password, number, select, textarea, file, toggle) | default, focus, error, disabled | Label di atas, petunjuk di bawah, galat menggantikan petunjuk (`--danger`, 13px). Galat muncul setelah field ditinggalkan atau setelah submit. |
| **SkillInput** | empty, filled, max | Ketik lalu Enter atau koma → jadi chip. Klik × untuk menghapus. Ada saran dari skill yang sudah ada (OQ-9). Batas 15 skill di profil, 10 di lowongan. |
| **JobCard** | default, hover, saved, applied, closed | Isi kartu: logo 40px, judul (h3), nama perusahaan + ikon centang kalau terverifikasi, chip level, mode kerja, dan lokasi, gaji (atau "Gaji dirahasiakan" dengan warna `--muted`), maksimal 4 chip skill + "+n", tanggal ("3 hari lalu"), tombol ikon bookmark. Badge "Sudah dilamar" (`--info`) atau "Ditutup" (`--neutral`) di kanan atas. Seluruh kartu adalah link. |
| **FilterPanel** | default, active | Isi: Kategori (select), Level (checkbox: Magang/Entry/Junior/Mid), Tipe (checkbox), Mode kerja (checkbox), Lokasi (select), Gaji minimal (select bertingkat: 3 jt, 5 jt, 8 jt, 12 jt+). Tombol "Hapus filter". Nilainya tersinkron dengan query URL. |
| **SearchBar** | empty, filled | Ikon kaca pembesar, placeholder "Cari posisi, skill, atau perusahaan". Di beranda tampil besar dengan tombol "Cari". Di header tampil ringkas. Submit dengan Enter. |
| **CategoryChip** | default, selected | Pill dengan nama kategori dan jumlah lowongan ("Data · 12"). |
| **StatusBadge** | lamaran (6), verifikasi (3), undangan (4), lowongan (3) | Pill dengan ikon + teks. Lihat §6. |
| **ApplicationTimeline** | — | Daftar vertikal: titik berwarna status, label, tanggal-waktu, dan "oleh Perusahaan" / "oleh kamu". Langkah terakhir ditebalkan. |
| **MatchScore** | high (≥ 75), medium (50–74), low (< 50) | Cincin progres 40px dengan angka di tengah. Warnanya `--success`, `--warning`, atau `--muted`, dengan label teks "Cocok tinggi/sedang/rendah". Kalau diklik, popover **ScoreBreakdown** terbuka. |
| **ScoreBreakdown** | — | Tiga baris berisi bar dan poin: Skill (x/60) dengan chip skill cocok (centang hijau) dan kurang (abu-abu), Level (x/20) dengan penjelasan, Screening (x/20) dengan jumlah jawaban lolos. |
| **RequirementFlag** | — | Badge "Tidak memenuhi syarat" (`--danger`, ikon `alert-triangle`). Tooltip-nya menyebutkan pertanyaan wajib yang gagal. |
| **ScreeningQuestionsEditor** (perusahaan) | 0–5 baris | Setiap baris berisi teks pertanyaan, tipe (select), toggle Wajib, dan jawaban yang diharapkan: radio Ya/Tidak untuk `yes_no`, operator (≥ / ≤) + angka untuk `number`, dan tersembunyi untuk `text`. Ada tombol tambah, hapus, naik, dan turun. |
| **ScreeningAnswerForm** (pelamar) | — | Field sesuai tipe pertanyaan: radio Ya/Tidak, input angka, atau textarea pendek. Pertanyaan wajib diberi tanda *. Jawaban yang diharapkan **tidak ditampilkan** ke pelamar. |
| **ApplyDialog** | ready, incomplete profile, submitting, success | Isi: ringkasan lowongan, CV yang dipakai (nama file + tautan "Ganti di profil"), surat pengantar (textarea + penghitung karakter), dan pertanyaan screening. Kalau profil belum lengkap, isinya diganti daftar yang kurang + tombol "Lengkapi profil". Setelah sukses: centang + "Lamaran terkirim" + tombol "Lihat Lamaran Saya". |
| **CandidateCard** (Cari Talenta) | default, invited | Isi: avatar inisial, nama, headline, lokasi, pendidikan · tahun lulus, pengalaman, chip skill (skill yang sedang difilter disorot), dan baris kecil `lock` "Kontak disembunyikan". Tombol "Undang". Kalau sudah diundang: "Diundang ke {lowongan}". |
| **InviteDialog** | — | Select lowongan yang masih dibuka (lowongan yang sudah pernah dipakai untuk mengundang kandidat ini dinonaktifkan), pesan opsional (500 karakter), dan sisa kuota hari ini ("17 undangan tersisa hari ini"). |
| **InvitationCard** (pencari kerja) | pending, accepted, declined, expired | Isi: logo + perusahaan, judul lowongan, pesan, tanggal, dan sisa hari sebelum kedaluwarsa. Tombol "Terima & lamar" (primary) dan "Tolak" (ghost). |
| **VerificationBanner** | pending, rejected | Kotak di atas dashboard perusahaan. Pending: `--warning`, "Profil perusahaan sedang ditinjau admin. Kamu bisa menyiapkan draft lowongan." Rejected: `--danger` + alasan + tombol "Perbaiki & ajukan ulang". |
| **DiscoverableToggle** | off, on, disabled | Switch + penjelasan: "Perusahaan terverifikasi bisa melihat nama, headline, skill, lokasi, dan pendidikanmu. Email, nomor HP, dan CV tetap tersembunyi." Toggle nonaktif selama profil belum lengkap, beserta alasannya. |
| **ProfileCompleteness** | — | Bar progres + daftar item yang belum diisi (dengan link). Tampil di halaman profil dan di dialog lamaran. |
| **CvUpload** | empty, uploading, uploaded, error | Area drop + tombol. Setelah upload: ikon PDF, nama file, ukuran, tanggal, tombol "Lihat" dan "Ganti". Galat: "Hanya PDF, maksimal 2 MB." |
| **StatCard** | loading, filled | Label (`--muted`), angka (`stat`), dan keterangan kecil (misalnya "+12 minggu ini"). |
| **ColumnChart** | loading, empty, filled | SVG tanpa library, dipakai ulang dari HeliozStore dengan warna `--primary`. Bar vertikal untuk lamaran per minggu. Setiap bar punya label nilai dan tooltip. |
| **BarList** | loading, empty, filled | Bar horizontal berisi label + nilai, untuk lowongan per kategori dan perusahaan teratas. |
| **DataTable** | loading, empty, filled | Header sticky, baris bisa diklik. Di mobile berubah jadi daftar kartu. |
| **Pagination** | — | Tombol Sebelumnya/Berikutnya + nomor halaman + "Menampilkan 1–12 dari 48". |
| **Modal / ConfirmDialog**, **Toast**, **EmptyState**, **Skeleton** | — | Pola yang sama dengan HeliozStore: Esc menutup modal, fokus dikembalikan ke elemen sebelumnya, toast tampil 3 detik di kanan atas (di mobile: bawah), skeleton seukuran konten. |

**Ikon:** lucide (`lucide-vue-next`), tanpa emoji.

---

## 4. Navigasi dan layout

### 4.1 Publik dan pencari kerja

- **Header** (putih, border bawah, sticky):
  - Kiri: logo Karsa.
  - Tengah, khusus desktop: "Cari Lowongan" dan "Perusahaan".
  - Kanan untuk tamu: "Masuk" (ghost), "Daftar" (primary), dan link kecil "Untuk Perusahaan".
  - Kanan untuk pencari kerja: ikon "Tersimpan", ikon "Undangan" dengan badge jumlah undangan pending, dan menu avatar (Profil & CV, Lamaran Saya, Undangan, Keluar).
- **Mobile:** ikon hamburger membuka drawer berisi link yang sama.
- **Footer:** logo, tagline "Langkah pertama kariermu di dunia tech.", link (Cari Lowongan, Untuk Perusahaan, Tentang), dan keterangan proyek demo dari PRD.

### 4.2 Dashboard perusahaan dan admin

- **Desktop:** sidebar kiri 240px (`--surface-2`) berisi logo dan menu, lalu konten dengan header halaman (judul + tombol utama).
- **Mobile:** sidebar menjadi drawer.
- **Menu perusahaan:** Dashboard, Lowongan, Cari Talenta, Undangan Terkirim, Profil Perusahaan, Keluar.
- **Menu admin:** Dashboard, Verifikasi Perusahaan (badge jumlah pending), Pengguna, Lowongan, Kategori, Keluar.
- Setiap tombol keluar membuka ConfirmDialog.

### 4.3 Route

| Path | Halaman | Akses |
|---|---|---|
| `/` | Beranda | Publik |
| `/lowongan` | Cari Lowongan | Publik |
| `/lowongan/:slug` | Detail Lowongan | Publik |
| `/perusahaan/:slug` | Profil Perusahaan publik | Publik |
| `/masuk`, `/daftar` | Masuk, Daftar (pilih peran) | Tamu |
| `/saya/profil` | Profil & CV | Pencari kerja |
| `/saya/lamaran`, `/saya/lamaran/:id` | Lamaran Saya, Detail + timeline | Pencari kerja |
| `/saya/tersimpan` | Lowongan Tersimpan | Pencari kerja |
| `/saya/undangan` | Undangan | Pencari kerja |
| `/perusahaan-saya` | Dashboard perusahaan | Perusahaan |
| `/perusahaan-saya/profil` | Profil Perusahaan | Perusahaan |
| `/perusahaan-saya/lowongan`, `/baru`, `/:id/edit` | Lowongan Saya, Form | Perusahaan |
| `/perusahaan-saya/lowongan/:id/pelamar` | Pelamar per lowongan | Perusahaan |
| `/perusahaan-saya/pelamar/:id` | Detail pelamar | Perusahaan |
| `/perusahaan-saya/talenta` | Cari Talenta | Perusahaan (terverifikasi) |
| `/perusahaan-saya/undangan` | Undangan Terkirim | Perusahaan |
| `/admin`, `/admin/perusahaan`, `/admin/pengguna`, `/admin/lowongan`, `/admin/kategori` | Panel admin | Admin |

Penjaga route (`beforeEach`), mengikuti pola HeliozStore:
- route `meta.role` yang tidak cocok → beranda sesuai peran + toast "Kamu tidak punya akses ke halaman itu."
- belum masuk → `/masuk?redirect=...`

---

## 5. Spesifikasi layar

### 5.1 Publik
- **Beranda:**
  1. Hero `--bg` dengan display "Temukan magang dan kerja pertamamu di dunia tech." beserta subjudul, SearchBar besar, dan chip pencarian populer (Frontend, Data Analyst, UI/UX, Magang Remote).
  2. Baris CategoryChip.
  3. "Lowongan terbaru": 6 JobCard dalam grid 2 kolom di desktop.
  4. "Magang & Fresh Graduate": 6 JobCard.
  5. "Perusahaan yang sedang merekrut": logo, nama, dan jumlah lowongan.
  6. CTA untuk perusahaan: "Pasang lowongan gratis".
- **Cari Lowongan:** judul "n lowongan ditemukan" + select urutan, FilterPanel, daftar JobCard, lalu Pagination. Chip filter aktif tampil di atas daftar dan bisa dihapus satu per satu.
- **Detail Lowongan:**
  - Kolom kiri (2/3): header (logo 64px, judul h1, perusahaan, chip), lalu bagian Deskripsi, Persyaratan, dan Skill.
  - Kolom kanan (1/3, sticky): kartu ringkasan berisi gaji, tipe, level, mode, lokasi, deadline ("Ditutup dalam 5 hari", berwarna `--warning` kalau ≤ 3 hari), jumlah pelamar, MatchScore (khusus pencari kerja yang sudah masuk), tombol **Lamar sekarang**, dan tombol Simpan.
  - Di mobile, kartu ringkasan pindah ke atas dan tombol Lamar menjadi bar bawah yang sticky.
  - Di bawahnya ada kartu "Tentang perusahaan".
- **Profil Perusahaan:** banner tipis `--primary-soft`, logo 96px, nama, industri · ukuran · lokasi, website, deskripsi, lalu daftar lowongan yang masih dibuka.
- **Masuk / Daftar:**
  - Kartu di tengah halaman.
  - Halaman Daftar dimulai dengan pilihan peran berupa dua kartu besar: "Saya mencari kerja" dan "Saya merekrut". Setelah memilih, tampil field sesuai peran. Peran perusahaan punya field tambahan "Nama perusahaan".

### 5.2 Pencari kerja
- **Profil & CV:** ProfileCompleteness, form profil dua kolom di desktop, kartu CvUpload, dan kartu DiscoverableToggle.
- **Lamaran Saya:**
  - Tab status: Semua, Aktif (Dilamar/Ditinjau/Interview), Diterima, Ditolak/Ditarik.
  - Setiap baris berisi logo, judul, perusahaan, StatusBadge, dan tanggal.
  - Detail lamaran berisi ringkasan lowongan, ApplicationTimeline, jawaban yang dikirim, CV yang dikirim, dan tombol "Tarik lamaran" (hanya saat status Dilamar, lewat ConfirmDialog).
- **Lowongan Tersimpan:** daftar JobCard dengan tombol hapus dari tersimpan.
- **Undangan:** tab Menunggu dan Riwayat, berisi InvitationCard. "Terima & lamar" membuka ApplyDialog dengan sumber `invited`.

### 5.3 Perusahaan
- **Dashboard:**
  - VerificationBanner (kalau belum terverifikasi)
  - 4 StatCard: Lowongan dibuka, Total pelamar, Pelamar baru 7 hari, Undangan menunggu
  - ColumnChart lamaran per minggu
  - pelamar per status (BarList)
  - tabel 5 lowongan dengan pelamar terbanyak
- **Lowongan Saya:** tabel berisi judul, status, pelamar, views, deadline, dan aksi (Lihat pelamar, Edit, Buka/Tutup, Hapus). Tombol utama "Buat lowongan".
- **Form Lowongan:**
  - Bagian 1 Info dasar, Bagian 2 Detail & gaji, Bagian 3 Skill, Bagian 4 Deskripsi & persyaratan, Bagian 5 Pertanyaan screening.
  - Tombol "Simpan draft" dan "Buka lowongan". Tombol "Buka lowongan" nonaktif selama perusahaan belum terverifikasi, dengan tooltip penjelasan.
- **Pelamar per lowongan:**
  - Header berisi judul lowongan + jumlah pelamar per status, ditampilkan sebagai tab yang juga berfungsi sebagai filter.
  - Baris filter: skor minimal (slider 0–100), skill, checkbox "Hanya yang memenuhi syarat", dan sumber.
  - Tabel berisi pelamar, MatchScore, RequirementFlag, sumber, status, dan tanggal.
- **Detail pelamar:**
  - Kolom kiri: profil lengkap (dengan email dan nomor HP), tombol "Unduh CV", surat pengantar, dan jawaban screening (ikon lolos/gagal per jawaban).
  - Kolom kanan: MatchScore + ScoreBreakdown terbuka, aksi status (hanya perpindahan yang sah sesuai PRD FR-9.2; "Tolak" lewat ConfirmDialog), catatan internal (textarea + simpan), dan ApplicationTimeline.
- **Cari Talenta:**
  - Panel filter (skill, lokasi, pendidikan, tahun lulus, pengalaman maks.) dan grid CandidateCard 2 kolom.
  - Kalau perusahaan belum terverifikasi, halaman berisi EmptyState "Fitur ini tersedia setelah perusahaanmu terverifikasi."
- **Undangan Terkirim:** tabel berisi kandidat, lowongan, status undangan, dan tanggal.
- **Profil Perusahaan:** form + upload logo + status verifikasi.

### 5.4 Admin
- **Dashboard:**
  - 4 StatCard: Pengguna, Perusahaan terverifikasi, Lowongan dibuka, Total lamaran
  - kartu "n perusahaan menunggu verifikasi" + link
  - ColumnChart lamaran per minggu
  - BarList lowongan per kategori
  - BarList 5 perusahaan teratas
- **Verifikasi Perusahaan:** tab Menunggu, Terverifikasi, Ditolak. Detail berupa drawer kanan berisi profil lengkap, tombol "Setujui", dan tombol "Tolak" (wajib alasan).
- **Pengguna:** pencarian + filter peran, tabel berisi nama, email, peran, status, tanggal daftar, dan aksi Nonaktifkan/Aktifkan (lewat ConfirmDialog).
- **Lowongan:** filter status dan perusahaan, aksi Tutup (wajib alasan) dan Hapus.
- **Kategori:** tabel + form sebaris (nama, slug otomatis), aksi Edit dan Hapus.

---

## 6. Format dan teks

| Hal | Aturan | Contoh |
|---|---|---|
| Gaji | Rentang dalam juta, per bulan | `Rp 5–8 jt/bulan`, `Rp 4 jt/bulan`, "Gaji dirahasiakan" |
| Angka besar | `Intl.NumberFormat('id-ID')` | `1.250` |
| Tanggal | Tanggal, bulan singkat (Indonesia), tahun | `24 Nov 2026` |
| Tanggal relatif | Untuk tanggal dipasang dan aktivitas | "Hari ini", "3 hari lalu" |
| Deadline | Sisa hari | "Ditutup dalam 5 hari", "Ditutup hari ini", "Sudah ditutup" |
| Bahasa | Indonesia, sapaan "kamu", kalimat pendek | "Kamu sudah melamar lowongan ini." |

**Label enum:**

| Field | Nilai → Label |
|---|---|
| `type` | `full_time` Penuh waktu · `part_time` Paruh waktu · `internship` Magang · `contract` Kontrak |
| `level` | `internship` Magang · `entry` Entry level · `junior` Junior · `mid` Mid level |
| `work_mode` | `onsite` Di kantor · `hybrid` Hybrid · `remote` Remote |

**Status lamaran:**

| Nilai (API) | Label (UI) | Warna | Ikon |
|---|---|---|---|
| `applied` | Dilamar | `--primary` | send |
| `reviewed` | Ditinjau | `--info` | eye |
| `interview` | Interview | `--warning` | calendar-clock |
| `accepted` | Diterima | `--success` | check-circle |
| `rejected` | Ditolak | `--danger` | x-circle |
| `withdrawn` | Ditarik | `--neutral` | undo-2 |

**Status lain:**

| Jenis | Nilai → Label (warna) |
|---|---|
| Verifikasi | `pending` Menunggu verifikasi (`--warning`) · `verified` Terverifikasi (`--success`) · `rejected` Ditolak (`--danger`) |
| Undangan | `pending` Menunggu (`--warning`) · `accepted` Diterima (`--success`) · `declined` Ditolak (`--neutral`) · `expired` Kedaluwarsa (`--neutral`) |
| Lowongan | `draft` Draft (`--neutral`) · `open` Dibuka (`--success`) · `closed` Ditutup (`--neutral`) |

---

## 7. Aksesibilitas (wajib lolos)

- [ ] Kontras teks ≥ 4.5:1.
- [ ] Semua fungsi bisa dipakai dengan keyboard, dan fokus selalu terlihat (`--focus`).
- [ ] Setiap input punya label, dan galat terhubung lewat `aria-describedby`.
- [ ] Logo perusahaan punya `alt` berisi nama perusahaan.
- [ ] Tombol ikon (bookmark, menu, tutup) punya `aria-label`. Bookmark memakai `aria-pressed`.
- [ ] Hasil pencarian dan perubahan status diumumkan dengan `aria-live="polite"` ("12 lowongan ditemukan", "Status diubah menjadi Interview").
- [ ] MatchScore punya label teks, bukan hanya cincin warna.
- [ ] Drawer dan modal mengunci fokus di dalamnya, bisa ditutup dengan Esc, dan mengembalikan fokus saat ditutup.
- [ ] Target sentuh ≥ 44×44px; tidak ada scroll horizontal di lebar 375px.
- [ ] `prefers-reduced-motion` mematikan transisi.

---

## 8. Arsitektur

```
Browser ──► frontend (Vue 3 SPA, nginx)  :5173
   │
   └──REST /api──► backend (Go, satu binary)  :8080
                     ├─ HTTP server (Gin)
                     ├─ goroutine WORKER    → email_jobs → SMTP
                     └─ goroutine SCHEDULER → tutup lowongan lewat deadline, undangan kedaluwarsa
                     │
                     ├──► MySQL 8.4 (volume dbdata)
                     ├──► volume uploads (cv/, logos/)
                     └──► Mailpit SMTP :1025 (inbox web :8025)
```

**Struktur folder yang diusulkan:**

```
Job Portal/
├── backend/
│   ├── cmd/server/main.go       # rakit config, DB, router, worker, scheduler, graceful shutdown
│   ├── internal/
│   │   ├── config/              # baca env
│   │   ├── database/            # koneksi GORM + jalankan migrasi
│   │   ├── models/              # struct GORM
│   │   ├── middleware/          # auth JWT, cek peran, cek perusahaan terverifikasi
│   │   ├── handlers/            # handler Gin per fitur (auth, jobs, applications, ...)
│   │   ├── services/            # logika bisnis (apply, matching, invitations)
│   │   ├── matching/            # fungsi skor murni + unit test
│   │   ├── mailer/              # template + SMTP
│   │   ├── worker/              # goroutine email
│   │   └── scheduler/           # goroutine terjadwal
│   ├── migrations/              # 0001_create_users.up.sql / .down.sql, ...
│   ├── seed/                    # data demo + contoh CV PDF
│   └── Dockerfile
├── frontend/                    # Vue 3 (pola HeliozStore)
├── docker-compose.yml
└── .github/workflows/ci.yml
```

---

## 9. Kontrak data (frontend ↔ backend)

**Aturan umum:**
- Semua endpoint berada di bawah `/api` dan memakai JSON. Autentikasi: `Authorization: Bearer <jwt>`.
- Gaji ditulis sebagai bilangan bulat rupiah. Tanggal-waktu memakai ISO 8601. Gambar dikirim sebagai URL lengkap (`logo_url`).
- **Daftar berhalaman:** `{ data: [...], page, per_page, total, last_page }`
- **Galat validasi (422):** `{ message, errors: { "field": "pesan" } }`. Bentuknya sama dengan Laravel, supaya `errorMessage()` dan `firstError()` dari HeliozStore bisa dipakai ulang.

Bentuk ringkas yang dipakai berulang:
- `JobSummary` = `{ id, slug, title, company: { name, slug, logo_url, verified }, category: { name, slug }, type, level, work_mode, location, salary_min, salary_max, show_salary, skills, deadline, created_at, is_saved?, application_status? }`
- `Me` = `{ id, name, email, role, is_active, company?: { id, name, slug, status } }`

### 9.1 Publik

| Method dan path | Respons |
|---|---|
| `POST /auth/register` `{ name, email, password, role, company_name? }` | 201 `{ user: Me, token }` |
| `POST /auth/login` `{ email, password }` | `{ user: Me, token }` · 401 kredensial salah · 403 akun nonaktif · 429 |
| `GET /home` | `{ categories: [ { name, slug, job_count } ], latest: [JobSummary], entry_level: [JobSummary], companies: [ { name, slug, logo_url, open_jobs } ] }` |
| `GET /categories` | `[ { id, name, slug, job_count } ]` |
| `GET /jobs?q&category&type&level&mode&location&salary_min&sort&page` | Daftar berhalaman berisi `JobSummary` (`level`, `type`, dan `mode` boleh dipisah koma; `sort`: `newest` / `salary` / `deadline`) |
| `GET /jobs/:slug` | `JobSummary` + `{ description, requirements, applicant_count, views, screening_questions: [ { id, question, type, required } ], company: { ..., industry, size, description, website }, match_preview? }`. Lowongan yang tidak publik → 404 |
| `GET /companies/:slug` | `{ name, slug, logo_url, industry, size, location, website, description, jobs: [JobSummary] }` |

`match_preview` hanya dikirim untuk pencari kerja yang sudah masuk: `{ score, breakdown }`.
`expected_answer` **tidak pernah** dikirim ke endpoint publik maupun pencari kerja.

### 9.2 Pencari kerja (`role = seeker`)

| Method dan path | Respons |
|---|---|
| `GET /auth/me` | `{ user: Me }` |
| `GET /me/profile`, `PUT /me/profile` | `{ headline, phone, location, education, major, graduation_year, experience_years, skills, bio, portfolio_url, cv: { name, size, uploaded_at, url } \| null, is_discoverable, completeness: { percent, missing: [ "cv", ... ] } }` |
| `POST /me/cv` (multipart `cv`) | Profil terbaru |
| `GET /me/cv` | File PDF milik sendiri |
| `POST /jobs/:id/apply` `{ cover_letter?, answers: [ { question_id, answer } ], invitation_id? }` | 201 `{ application }` · 409 sudah melamar · 422 profil belum lengkap (`errors.profile`) / jawaban tidak valid |
| `GET /me/applications?status&page` | Daftar berhalaman berisi `{ id, job: JobSummary, status, source, created_at, updated_at }` |
| `GET /me/applications/:id` | `{ ..., cover_letter, answers: [ { question, answer } ], cv_url, events: [ { from_status, to_status, by: "company" \| "you", created_at } ] }` |
| `GET /me/applications/:id/cv` | File PDF snapshot yang dikirim bersama lamaran (`cv_url` menunjuk ke sini) |
| `POST /me/applications/:id/withdraw` | Lamaran terbaru · 422 kalau status bukan `applied` |
| `GET /me/saved-jobs?page` | Daftar berhalaman berisi `JobSummary` + `is_open` |
| `PUT /me/saved-jobs/:jobId`, `DELETE /me/saved-jobs/:jobId` | 204 (idempoten) |
| `GET /me/invitations?status` | `[ { id, job: JobSummary, company: { name, slug, logo_url }, message, status, created_at, expires_at } ]` |
| `POST /me/invitations/:id/decline` | Undangan terbaru |

Menerima undangan dilakukan lewat `POST /jobs/:id/apply` dengan `invitation_id`. Server memeriksa bahwa undangan itu milik pengguna, masih `pending`, dan untuk lowongan yang sama.

### 9.3 Perusahaan (`role = company`)

Endpoint bertanda **(V)** mewajibkan perusahaan `verified`. Kalau belum, server membalas 403 `{ message: "Perusahaan belum terverifikasi" }`.

| Method dan path | Respons |
|---|---|
| `GET /company/profile`, `PUT /company/profile` | `{ id, name, slug, logo_url, industry, size, location, website, description, status, rejection_reason }`. Kalau statusnya `rejected`, PUT juga mengubah status kembali ke `pending` |
| `POST /company/logo` (multipart `logo`) | Profil terbaru |
| `GET /company/dashboard` | `{ open_jobs, total_applicants, new_applicants_7d, pending_invitations, by_status: { applied, reviewed, ... }, weekly: [ { week_start, count } ], top_jobs: [ { id, title, applicants } ] }` |
| `GET /company/jobs?status&page` | Daftar berhalaman berisi `{ ...JobSummary, status, applicant_count, views }` |
| `GET /company/jobs/:id` | Lowongan lengkap + `screening_questions` (termasuk `expected_answer`) |
| `POST /company/jobs`, `PUT /company/jobs/:id` | Body: field lowongan + `status` (`draft` / `open`) + `screening_questions: [ { question, type, required, expected_answer } ]`. `status: open` **(V)** |
| `POST /company/jobs/:id/close`, `POST /company/jobs/:id/open` **(V)** | Lowongan terbaru |
| `DELETE /company/jobs/:id` | 204 · 409 kalau sudah ada pelamar |
| `GET /company/jobs/:id/applications?status&min_score&skill&meets&source&sort&page` | `{ counts: { applied, reviewed, ... }, ...berhalaman: [ { id, candidate: { name, headline, location }, match_score, meets_requirements, source, status, created_at } ] }` |
| `GET /company/applications/:id` | `{ ..., candidate: { name, email, phone, headline, location, education, major, graduation_year, experience_years, skills, bio, portfolio_url }, cover_letter, answers: [ { question, answer, passed, required } ], score_breakdown: { skill: { points, matched, missing }, level: { points, note }, screening: { points, passed, total } }, company_note, events } }` |
| `PATCH /company/applications/:id` `{ status?, company_note? }` | Lamaran terbaru · 422 perpindahan tidak sah · 409 status sudah berubah |
| `GET /company/applications/:id/cv` | File PDF (snapshot) |
| `GET /company/talents?skills&location&education&grad_from&grad_to&max_experience&page` **(V)** | Daftar berhalaman berisi `{ id, name, headline, location, education, major, graduation_year, experience_years, skills, invited_job_ids }`. **Tanpa email, nomor HP, dan CV** |
| `POST /company/invitations` `{ user_id, job_id, message? }` **(V)** | 201 `{ invitation }` · 409 sudah diundang atau sudah melamar · 429 kuota harian habis |
| `GET /company/invitations?status&page` | Daftar berhalaman berisi `{ id, candidate: { name, headline }, job: { id, title }, status, created_at }` + `{ remaining_today }` |

### 9.4 Admin (`role = admin`)

| Method dan path | Keterangan |
|---|---|
| `GET /admin/dashboard` | `{ users: { seeker, company, admin }, companies: { pending, verified, rejected }, open_jobs, applications, weekly: [ { week_start, count } ], jobs_by_category: [ { name, count } ], top_companies: [ { name, applications } ] }` |
| `GET /admin/companies?status&page`, `GET /admin/companies/:id` | Daftar dan detail perusahaan |
| `POST /admin/companies/:id/verify`, `POST /admin/companies/:id/reject` `{ reason }` | Perusahaan terbaru + email diantrikan |
| `GET /admin/users?q&role&page`, `PATCH /admin/users/:id` `{ is_active }` | Daftar dan aktif/nonaktif pengguna · 422 kalau admin menonaktifkan diri sendiri |
| `GET /admin/jobs?status&company&page`, `POST /admin/jobs/:id/close` `{ reason }`, `DELETE /admin/jobs/:id` | Moderasi · 409 kalau menghapus lowongan yang punya lamaran |
| `GET /admin/categories`, `POST /admin/categories`, `PUT /admin/categories/:id`, `DELETE /admin/categories/:id` | CRUD · 409 kalau kategori masih dipakai |

### 9.5 Galat

| Kode | Arti | Perilaku UI |
|---|---|---|
| 401 | Belum masuk / token kedaluwarsa | Ke `/masuk?redirect=...` |
| 403 | Peran salah / akun nonaktif / perusahaan belum terverifikasi | Tampilkan `message` dari server. Kalau akun nonaktif: keluar, lalu ke halaman masuk |
| 404 | Tidak ada atau milik orang lain | "Data tidak ditemukan." |
| 409 | Konflik (duplikat, status sudah berubah, data masih dipakai) | `message` dari server |
| 413 / 422 file | File terlalu besar atau bukan PDF | Galat di bawah komponen upload |
| 422 | Validasi | Galat ditampilkan di bawah field terkait |
| 429 | Terlalu banyak percobaan / kuota undangan habis | `message` dari server |
| 5xx / jaringan | Server atau jaringan gagal | "Terjadi kesalahan. Coba lagi." |

---

## 10. Keputusan desain (terkunci)

- **D-1, Logo:**
  - Ikon: kotak membulat (radius 10px) `--primary` berisi huruf **K** putih geometris, dengan kaki kanan bawah yang memanjang menjadi panah naik ke kanan, melambangkan langkah karier.
  - Wordmark: "Karsa", Manrope 800, warna `--text`.
  - Ikon yang sama dipakai sebagai favicon SVG.
- **D-2, Hero beranda:** latar polos `--bg` dengan judul besar dan SearchBar. Tanpa carousel dan tanpa foto stok.
- **D-3, Grafik:** SVG buatan sendiri (ColumnChart dari HeliozStore + BarList), tanpa library grafik.
- **D-4, Skor:** cincin + angka + label, dan rinciannya selalu bisa dibuka (Prinsip 6).
