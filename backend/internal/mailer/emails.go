package mailer

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"karsa/internal/models"
)

var statusLines = map[string]string{
	"reviewed":  "Lamaranmu sedang ditinjau oleh perusahaan.",
	"interview": "Selamat, kamu lanjut ke tahap interview. Perusahaan akan menghubungimu untuk jadwalnya.",
	"accepted":  "Selamat, lamaranmu diterima!",
	"rejected":  "Perusahaan memutuskan untuk tidak melanjutkan lamaranmu kali ini. Tetap semangat!",
}

var statusLabels = map[string]string{
	"reviewed":  "Ditinjau",
	"interview": "Interview",
	"accepted":  "Diterima",
	"rejected":  "Ditolak",
}

func (mailer Mailer) ApplicationReceived(tx *gorm.DB, email, name, jobTitle, companyName string, applicationID uint) error {
	return queue(tx, email, "application_received", "Lamaran terkirim: "+jobTitle, models.EmailPayload{
		Heading: "Lamaranmu sudah terkirim",
		Lines:   []string{"Halo " + name + ",", fmt.Sprintf("Lamaranmu untuk %s di %s sudah kami teruskan ke perusahaan.", jobTitle, companyName)},
		Button:  "Lihat lamaran",
		URL:     mailer.applicationURL(applicationID),
	})
}

func (mailer Mailer) StatusChanged(tx *gorm.DB, email, name, jobTitle, companyName, status string, applicationID uint) error {
	return queue(tx, email, "application_status", fmt.Sprintf("Lamaran %s: %s", statusLabels[status], jobTitle), models.EmailPayload{
		Heading: "Status lamaranmu: " + statusLabels[status],
		Lines:   []string{"Halo " + name + ",", fmt.Sprintf("Ada kabar untuk lamaranmu sebagai %s di %s.", jobTitle, companyName), statusLines[status]},
		Button:  "Lihat lamaran",
		URL:     mailer.applicationURL(applicationID),
	})
}

func (mailer Mailer) InvitationReceived(tx *gorm.DB, email, name, jobTitle, companyName string, message *string) error {
	lines := []string{"Halo " + name + ",", fmt.Sprintf("%s mengundangmu melamar sebagai %s.", companyName, jobTitle)}
	if message != nil {
		lines = append(lines, "Pesan: "+*message)
	}
	return queue(tx, email, "invitation", "Undangan melamar dari "+companyName, models.EmailPayload{
		Heading: "Kamu diundang melamar",
		Lines:   lines,
		Button:  "Lihat undangan",
		URL:     mailer.Config.FrontendURL + "/saya/undangan",
	})
}

func (mailer Mailer) CompanyVerified(tx *gorm.DB, email, companyName string) error {
	return queue(tx, email, "company_verified", companyName+" sudah terverifikasi", models.EmailPayload{
		Heading: "Perusahaanmu sudah terverifikasi",
		Lines:   []string{companyName + " sekarang bisa membuka lowongan dan memakai Cari Talenta."},
		Button:  "Buka dashboard",
		URL:     mailer.Config.FrontendURL + "/perusahaan-saya",
	})
}

func (mailer Mailer) CompanyRejected(tx *gorm.DB, email, companyName, reason string) error {
	return queue(tx, email, "company_rejected", "Verifikasi "+companyName+" ditolak", models.EmailPayload{
		Heading: "Verifikasi belum disetujui",
		Lines:   []string{"Alasan: " + reason, "Perbaiki profil perusahaan untuk mengajukan ulang."},
		Button:  "Perbaiki profil",
		URL:     mailer.Config.FrontendURL + "/perusahaan-saya/profil",
	})
}

func (mailer Mailer) applicationURL(id uint) string {
	return fmt.Sprintf("%s/saya/lamaran/%d", mailer.Config.FrontendURL, id)
}

func queue(tx *gorm.DB, email, template, subject string, payload models.EmailPayload) error {
	return tx.Create(&models.EmailJob{
		ToEmail:     email,
		Subject:     subject,
		Template:    template,
		Payload:     payload,
		Status:      "pending",
		AvailableAt: time.Now(),
	}).Error
}
