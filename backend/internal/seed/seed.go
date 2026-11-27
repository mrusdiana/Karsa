package seed

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"karsa/internal/config"
	"karsa/internal/matching"
	"karsa/internal/models"
	"karsa/internal/util"
)

const (
	demoPassword = "password"
	million      = 1_000_000
)

type seeder struct {
	tx           *gorm.DB
	uploadDir    string
	passwordHash string
	now          time.Time
}

func Run(db *gorm.DB, cfg config.Config) error {
	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		log.Println("seed dilewati: tabel users sudah berisi data")
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(demoPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		return seeder{tx: tx, uploadDir: cfg.UploadDir, passwordHash: string(hash), now: time.Now()}.run()
	})
	if err == nil {
		log.Println("seed selesai: admin@karsa.test, perusahaan1@karsa.test, nadia@karsa.test (password: " + demoPassword + ")")
	}
	return err
}

func (seeder seeder) run() error {
	if _, err := seeder.createUser("Admin Karsa", "admin@karsa.test", "admin"); err != nil {
		return err
	}
	categoryIDs, err := seeder.createCategories()
	if err != nil {
		return err
	}
	createdCompanies, err := seeder.createCompanies()
	if err != nil {
		return err
	}
	createdJobs, err := seeder.createJobs(createdCompanies, categoryIDs)
	if err != nil {
		return err
	}
	profiles, err := seeder.createSeekers()
	if err != nil {
		return err
	}
	if err := seeder.createApplications(createdCompanies, createdJobs, profiles); err != nil {
		return err
	}
	return seeder.createExtras(createdJobs, profiles)
}

func (seeder seeder) createUser(name, email, role string) (models.User, error) {
	user := models.User{Name: name, Email: email, PasswordHash: seeder.passwordHash, Role: role, IsActive: true}
	return user, seeder.tx.Create(&user).Error
}

func (seeder seeder) createCategories() (map[string]uint, error) {
	ids := map[string]uint{}
	for _, name := range categoryNames {
		category := models.Category{Name: name, Slug: util.Slugify(name)}
		if err := seeder.tx.Create(&category).Error; err != nil {
			return nil, err
		}
		ids[name] = category.ID
	}
	return ids, nil
}

func (seeder seeder) createCompanies() ([]models.Company, error) {
	created := make([]models.Company, len(companies))
	for index, data := range companies {
		user, err := seeder.createUser("Rekruter "+data.name, fmt.Sprintf("perusahaan%d@karsa.test", index+1), "company")
		if err != nil {
			return nil, err
		}
		company := models.Company{
			UserID:      user.ID,
			Name:        data.name,
			Slug:        util.UniqueSlug(data.name),
			Industry:    &data.industry,
			Size:        &data.size,
			Location:    &data.location,
			Website:     &data.website,
			Description: &data.description,
			Status:      data.status,
		}
		if data.status == "verified" {
			company.VerifiedAt = &seeder.now
		}
		if err := seeder.tx.Create(&company).Error; err != nil {
			return nil, err
		}
		created[index] = company
	}
	return created, nil
}

func (seeder seeder) createJobs(createdCompanies []models.Company, categoryIDs map[string]uint) ([]models.Job, error) {
	created := make([]models.Job, len(jobs))
	for index, data := range jobs {
		salaryMin, salaryMax := data.salaryMin*million, data.salaryMax*million
		job := models.Job{
			CompanyID:    createdCompanies[data.company].ID,
			CategoryID:   categoryIDs[data.category],
			Title:        data.title,
			Slug:         util.UniqueSlug(data.title),
			Type:         data.jobType,
			Level:        data.level,
			WorkMode:     data.mode,
			SalaryMin:    &salaryMin,
			SalaryMax:    &salaryMax,
			ShowSalary:   index%5 != 4,
			Skills:       data.skills,
			Description:  "Kamu akan bergabung dengan tim " + createdCompanies[data.company].Name + " sebagai " + data.title + ". Kamu mengerjakan project nyata bersama mentor, mengikuti review mingguan, dan ikut merencanakan fitur berikutnya.",
			Requirements: "Memahami dasar skill yang tercantum.\nMau belajar dan terbuka pada masukan.\nBisa berkomunikasi dengan jelas dalam tim.",
			Deadline:     time.Date(seeder.now.Year(), seeder.now.Month(), seeder.now.Day()+data.deadlineDays, 0, 0, 0, 0, time.Local),
			Status:       data.status,
			Views:        uint(15 + index*7%90),
			CreatedAt:    seeder.now.AddDate(0, 0, -(len(jobs)-index)*3/2),
		}
		if data.location != "" {
			job.Location = &data.location
		}
		job.Questions = screeningFor(index, data)
		if err := seeder.tx.Create(&job).Error; err != nil {
			return nil, err
		}
		created[index] = job
	}
	return created, nil
}

func screeningFor(index int, data jobData) []models.ScreeningQuestion {
	if index%3 == 2 {
		return nil
	}
	yes, oneYear := "yes", ">=1"
	questions := []models.ScreeningQuestion{{Question: "Bisa mulai dalam 1 bulan?", Type: "yes_no", Required: true, ExpectedAnswer: &yes}}
	if data.level == "junior" || data.level == "mid" {
		questions = append(questions, models.ScreeningQuestion{Question: "Berapa tahun pengalaman kerjamu?", Type: "number", Required: true, ExpectedAnswer: &oneYear, Position: 1})
	}
	return questions
}

func (seeder seeder) createSeekers() ([]models.SeekerProfile, error) {
	profiles := make([]models.SeekerProfile, len(seekers))
	for index, data := range seekers {
		user, err := seeder.createUser(data.name, data.email, "seeker")
		if err != nil {
			return nil, err
		}
		cvPath, err := writeCV(seeder.uploadDir, "cv", data.name, data.headline)
		if err != nil {
			return nil, err
		}
		graduationYear := seeder.now.Year() + data.graduationOffset
		phone := fmt.Sprintf("0812000000%02d", index)
		profile := models.SeekerProfile{
			UserID:          user.ID,
			Headline:        &data.headline,
			Phone:           &phone,
			Location:        &data.location,
			Education:       &data.education,
			Major:           &data.major,
			GraduationYear:  &graduationYear,
			ExperienceYears: data.experience,
			Skills:          data.skills,
			CVPath:          &cvPath,
			IsDiscoverable:  data.isDiscoverable,
		}
		if err := seeder.tx.Create(&profile).Error; err != nil {
			return nil, err
		}
		profiles[index] = profile
	}
	return profiles, nil
}

func (seeder seeder) createApplications(createdCompanies []models.Company, createdJobs []models.Job, profiles []models.SeekerProfile) error {
	for _, data := range applications {
		job, profile := createdJobs[data.job], profiles[data.seeker]
		cvPath, err := writeCV(seeder.uploadDir, "applications", seekers[data.seeker].name, seekers[data.seeker].headline)
		if err != nil {
			return err
		}
		answers, passed := sampleAnswers(job, profile)
		score, breakdown := matching.Score(job, profile, passed, seeder.now.Year())
		appliedAt := seeder.now.AddDate(0, 0, -data.daysAgo)
		application := models.Application{
			JobID:             job.ID,
			UserID:            profile.UserID,
			CVPath:            cvPath,
			Status:            "applied",
			Source:            "applied",
			MatchScore:        score,
			ScoreBreakdown:    breakdown,
			MeetsRequirements: !slices.Contains(passed, false),
			Answers:           answers,
			Events:            []models.ApplicationEvent{{ToStatus: "applied", ChangedBy: profile.UserID, CreatedAt: appliedAt}},
			CreatedAt:         appliedAt,
			UpdatedAt:         appliedAt,
		}
		for step, status := range data.statuses {
			previous := application.Status
			changedBy := createdCompanies[jobs[data.job].company].UserID
			if status == "withdrawn" {
				changedBy = profile.UserID
			}
			changedAt := appliedAt.Add(time.Duration(step+1) * 20 * time.Hour)
			application.Events = append(application.Events, models.ApplicationEvent{FromStatus: &previous, ToStatus: status, ChangedBy: changedBy, CreatedAt: changedAt})
			application.Status, application.UpdatedAt = status, changedAt
		}
		if err := seeder.tx.Omit("Job", "User").Create(&application).Error; err != nil {
			return err
		}
	}
	return nil
}

func sampleAnswers(job models.Job, profile models.SeekerProfile) ([]models.ApplicationAnswer, []bool) {
	answers := []models.ApplicationAnswer{}
	passed := []bool{}
	for _, question := range job.Questions {
		answer := "yes"
		if question.Type == "number" {
			answer = strconv.Itoa(profile.ExperienceYears)
		}
		isPassed, _ := matching.Evaluate(question, answer)
		answers = append(answers, models.ApplicationAnswer{Question: question.Question, Required: question.Required, Answer: answer, Passed: isPassed})
		passed = append(passed, isPassed)
	}
	return answers, passed
}

func (seeder seeder) createExtras(createdJobs []models.Job, profiles []models.SeekerProfile) error {
	message := "Profilmu cocok dengan tim kami. Yuk, melamar!"
	for _, invitation := range []models.Invitation{
		{JobID: createdJobs[16].ID, CompanyID: createdJobs[16].CompanyID, UserID: profiles[7].UserID, Status: "pending"},
		{JobID: createdJobs[11].ID, CompanyID: createdJobs[11].CompanyID, UserID: profiles[0].UserID, Status: "pending"},
		{JobID: createdJobs[10].ID, CompanyID: createdJobs[10].CompanyID, UserID: profiles[4].UserID, Status: "declined"},
		{JobID: createdJobs[23].ID, CompanyID: createdJobs[23].CompanyID, UserID: profiles[6].UserID, Status: "expired"},
	} {
		invitation.Message = &message
		if err := seeder.tx.Omit(clause.Associations).Create(&invitation).Error; err != nil {
			return err
		}
	}
	return seeder.tx.Create(&[]models.SavedJob{
		{UserID: profiles[0].UserID, JobID: createdJobs[2].ID},
		{UserID: profiles[0].UserID, JobID: createdJobs[23].ID},
		{UserID: profiles[1].UserID, JobID: createdJobs[14].ID},
	}).Error
}
