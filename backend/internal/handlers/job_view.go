package handlers

import (
	"time"

	"gorm.io/gorm"

	"karsa/internal/models"
)

type companySummary struct {
	Name     string  `json:"name"`
	Slug     string  `json:"slug"`
	LogoURL  *string `json:"logo_url"`
	Verified bool    `json:"verified"`
}

type categorySummary struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type jobSummary struct {
	ID                uint            `json:"id"`
	Slug              string          `json:"slug"`
	Title             string          `json:"title"`
	Company           companySummary  `json:"company"`
	Category          categorySummary `json:"category"`
	Type              string          `json:"type"`
	Level             string          `json:"level"`
	WorkMode          string          `json:"work_mode"`
	Location          *string         `json:"location"`
	SalaryMin         *uint           `json:"salary_min"`
	SalaryMax         *uint           `json:"salary_max"`
	ShowSalary        bool            `json:"show_salary"`
	Skills            []string        `json:"skills"`
	Deadline          string          `json:"deadline"`
	CreatedAt         time.Time       `json:"created_at"`
	IsSaved           *bool           `json:"is_saved,omitempty"`
	ApplicationStatus *string         `json:"application_status,omitempty"`
}

type managedJob struct {
	jobSummary
	Status         string  `json:"status"`
	ModerationNote *string `json:"moderation_note"`
	ApplicantCount int64   `json:"applicant_count"`
	Views          uint    `json:"views"`
}

func withJobRelations(query *gorm.DB) *gorm.DB {
	return query.Preload("Company").Preload("Category")
}

func (handler *Handler) companySummary(company models.Company) companySummary {
	return companySummary{
		Name:     company.Name,
		Slug:     company.Slug,
		LogoURL:  handler.logoURL(company.LogoPath),
		Verified: company.Status == "verified",
	}
}

func (handler *Handler) fullSummary(job models.Job) jobSummary {
	return jobSummary{
		ID:         job.ID,
		Slug:       job.Slug,
		Title:      job.Title,
		Company:    handler.companySummary(*job.Company),
		Category:   categorySummary{Name: job.Category.Name, Slug: job.Category.Slug},
		Type:       job.Type,
		Level:      job.Level,
		WorkMode:   job.WorkMode,
		Location:   job.Location,
		SalaryMin:  job.SalaryMin,
		SalaryMax:  job.SalaryMax,
		ShowSalary: job.ShowSalary,
		Skills:     job.Skills,
		Deadline:   job.Deadline.Format(time.DateOnly),
		CreatedAt:  job.CreatedAt,
	}
}

func (handler *Handler) publicSummary(job models.Job) jobSummary {
	summary := handler.fullSummary(job)
	if !job.ShowSalary {
		summary.SalaryMin, summary.SalaryMax = nil, nil
	}
	return summary
}

func (handler *Handler) publicSummaries(jobs []models.Job, viewer *models.User) ([]jobSummary, error) {
	summaries := make([]jobSummary, len(jobs))
	for index, job := range jobs {
		summaries[index] = handler.publicSummary(job)
	}
	if viewer == nil || viewer.Role != "seeker" || len(jobs) == 0 {
		return summaries, nil
	}

	jobIDs := idsOf(jobs)
	var savedIDs []uint
	if err := handler.DB.Model(&models.SavedJob{}).Where("user_id = ? AND job_id IN ?", viewer.ID, jobIDs).Pluck("job_id", &savedIDs).Error; err != nil {
		return nil, err
	}
	var applications []models.Application
	if err := handler.DB.Select("job_id", "status").Where("user_id = ? AND job_id IN ?", viewer.ID, jobIDs).Find(&applications).Error; err != nil {
		return nil, err
	}
	saved := map[uint]bool{}
	for _, id := range savedIDs {
		saved[id] = true
	}
	statuses := map[uint]string{}
	for _, application := range applications {
		statuses[application.JobID] = application.Status
	}
	for index := range summaries {
		isSaved := saved[summaries[index].ID]
		summaries[index].IsSaved = &isSaved
		if status, ok := statuses[summaries[index].ID]; ok {
			summaries[index].ApplicationStatus = &status
		}
	}
	return summaries, nil
}

func (handler *Handler) managedJobs(jobs []models.Job) ([]managedJob, error) {
	counts, err := handler.applicantCounts(idsOf(jobs))
	views := make([]managedJob, len(jobs))
	for index, job := range jobs {
		views[index] = managedJob{
			jobSummary:     handler.fullSummary(job),
			Status:         job.Status,
			ModerationNote: job.ModerationNote,
			ApplicantCount: counts[job.ID],
			Views:          job.Views,
		}
	}
	return views, err
}

func (handler *Handler) applicantCounts(jobIDs []uint) (map[uint]int64, error) {
	var rows []struct {
		JobID uint
		Count int64
	}
	err := handler.DB.Model(&models.Application{}).
		Select("job_id, COUNT(*) AS count").
		Where("job_id IN ?", jobIDs).Group("job_id").Scan(&rows).Error
	counts := map[uint]int64{}
	for _, row := range rows {
		counts[row.JobID] = row.Count
	}
	return counts, err
}

func idsOf(jobs []models.Job) []uint {
	ids := make([]uint, len(jobs))
	for index, job := range jobs {
		ids[index] = job.ID
	}
	return ids
}
