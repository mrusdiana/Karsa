package handlers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"karsa/internal/httpx"
	"karsa/internal/matching"
	"karsa/internal/middleware"
	"karsa/internal/models"
	"karsa/internal/util"
)

const maxDeadlineDays = 90

type questionInput struct {
	Question       string `json:"question" binding:"required,notblank,max=255"`
	Type           string `json:"type" binding:"required,oneof=yes_no number text"`
	Required       bool   `json:"required"`
	ExpectedAnswer string `json:"expected_answer" binding:"max=20"`
}

type jobInput struct {
	CategoryID         uint            `json:"category_id" binding:"required"`
	Title              string          `json:"title" binding:"required,notblank,max=120"`
	Type               string          `json:"type" binding:"required,oneof=full_time part_time internship contract"`
	Level              string          `json:"level" binding:"required,oneof=internship entry junior mid"`
	WorkMode           string          `json:"work_mode" binding:"required,oneof=onsite hybrid remote"`
	Location           string          `json:"location" binding:"max=100"`
	SalaryMin          *uint           `json:"salary_min" binding:"omitempty,max=4000000000"`
	SalaryMax          *uint           `json:"salary_max" binding:"omitempty,max=4000000000"`
	ShowSalary         bool            `json:"show_salary"`
	Skills             []string        `json:"skills" binding:"required,min=1,max=10,dive,required,notblank,max=40"`
	Description        string          `json:"description" binding:"required,notblank,max=10000"`
	Requirements       string          `json:"requirements" binding:"required,notblank,max=10000"`
	Deadline           string          `json:"deadline" binding:"required,datetime=2006-01-02"`
	Status             string          `json:"status" binding:"required,oneof=draft open closed"`
	ScreeningQuestions []questionInput `json:"screening_questions" binding:"max=5,dive"`
}

type companyJobDetail struct {
	managedJob
	CategoryID         uint                       `json:"category_id"`
	Description        string                     `json:"description"`
	Requirements       string                     `json:"requirements"`
	ScreeningQuestions []models.ScreeningQuestion `json:"screening_questions"`
}

func (handler *Handler) ListCompanyJobs(ctx *gin.Context) {
	query := withJobRelations(handler.DB).Where("company_id = ?", middleware.CurrentUser(ctx).Company.ID)
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	handler.respondManagedJobs(ctx, query.Order("created_at DESC, id DESC"))
}

func (handler *Handler) ShowCompanyJob(ctx *gin.Context) {
	if job, ok := handler.findOwnedJob(ctx); ok {
		handler.respondCompanyJob(ctx, 200, job.ID)
	}
}

func (handler *Handler) CreateJob(ctx *gin.Context) {
	var input jobInput
	if !httpx.BindJSON(ctx, &input) {
		return
	}
	company := middleware.CurrentUser(ctx).Company
	job := models.Job{CompanyID: company.ID, Slug: util.UniqueSlug(input.Title), Status: "draft"}
	if handler.fillJob(ctx, &job, input, *company) {
		handler.saveJob(ctx, job, 201)
	}
}

func (handler *Handler) UpdateJob(ctx *gin.Context) {
	job, ok := handler.findOwnedJob(ctx)
	if !ok {
		return
	}
	var input jobInput
	if !httpx.BindJSON(ctx, &input) {
		return
	}
	if handler.fillJob(ctx, &job, input, *middleware.CurrentUser(ctx).Company) {
		handler.saveJob(ctx, job, 200)
	}
}

func (handler *Handler) OpenJob(ctx *gin.Context) {
	job, ok := handler.findOwnedJob(ctx)
	if !ok {
		return
	}
	if job.ModerationNote != nil {
		httpx.Error(ctx, 403, "Lowongan ditutup admin: "+*job.ModerationNote)
		return
	}
	if job.Deadline.Format(time.DateOnly) < today() {
		httpx.FieldErrors(ctx, map[string]string{"deadline": "Deadline sudah lewat. Ubah deadline sebelum membuka lowongan."})
		return
	}
	result := handler.DB.Model(&job).Where("moderation_note IS NULL").Update("status", "open")
	if result.Error != nil {
		httpx.ServerError(ctx, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		httpx.Error(ctx, 409, "Status sudah berubah, muat ulang halaman.")
		return
	}
	handler.respondCompanyJob(ctx, 200, job.ID)
}

func (handler *Handler) CloseJob(ctx *gin.Context) {
	job, ok := handler.findOwnedJob(ctx)
	if !ok {
		return
	}
	if job.Status == "draft" {
		httpx.Error(ctx, 422, "Draft belum dibuka, jadi tidak bisa ditutup.")
		return
	}
	err := handler.DB.Transaction(func(tx *gorm.DB) error {
		return closeJob(tx, job.ID, nil)
	})
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	handler.respondCompanyJob(ctx, 200, job.ID)
}

func (handler *Handler) DeleteJob(ctx *gin.Context) {
	if job, ok := handler.findOwnedJob(ctx); ok {
		handler.deleteJob(ctx, job)
	}
}

func (handler *Handler) deleteJob(ctx *gin.Context, job models.Job) {
	err := handler.DB.Delete(&job).Error
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		httpx.Error(ctx, 409, "Lowongan sudah punya pelamar. Tutup lowongan ini saja.")
		return
	}
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	ctx.Status(204)
}

func closeJob(tx *gorm.DB, jobID uint, moderationNote *string) error {
	changes := map[string]any{"status": "closed"}
	if moderationNote != nil {
		changes["moderation_note"] = *moderationNote
	}
	if err := tx.Model(&models.Job{}).Where("id = ?", jobID).Updates(changes).Error; err != nil {
		return err
	}
	return expireInvitations(tx, jobID)
}

func expireInvitations(tx *gorm.DB, jobID uint) error {
	return tx.Model(&models.Invitation{}).Where("job_id = ? AND status = ?", jobID, "pending").Update("status", "expired").Error
}

func (handler *Handler) fillJob(ctx *gin.Context, job *models.Job, input jobInput, company models.Company) bool {
	isOpening := input.Status == "open" && job.Status != "open"
	if isOpening && company.Status != "verified" {
		httpx.Error(ctx, 403, "Perusahaan belum terverifikasi")
		return false
	}
	if isOpening && job.ModerationNote != nil {
		httpx.Error(ctx, 403, "Lowongan ditutup admin: "+*job.ModerationNote)
		return false
	}

	fields := map[string]string{}
	if input.Status == "draft" && job.Status != "draft" {
		fields["status"] = "Lowongan yang sudah pernah dibuka tidak bisa jadi draft."
	}
	if input.Status == "closed" && job.Status == "draft" {
		fields["status"] = "Draft belum dibuka, jadi tidak bisa ditutup."
	}
	deadline, err := time.ParseInLocation(time.DateOnly, input.Deadline, time.Local)
	latest := time.Now().AddDate(0, 0, maxDeadlineDays).Format(time.DateOnly)
	isOutOfRange := input.Deadline < today() || input.Deadline > latest
	if err != nil || (input.Status != "closed" && isOutOfRange) {
		fields["deadline"] = "Deadline harus antara hari ini dan 90 hari ke depan."
	}
	if input.WorkMode != "remote" && optional(input.Location) == nil {
		fields["location"] = "Wajib diisi untuk kerja di kantor atau hybrid."
	}
	if input.SalaryMin != nil && input.SalaryMax != nil && *input.SalaryMax < *input.SalaryMin {
		fields["salary_max"] = "Gaji maksimal harus lebih besar atau sama dengan gaji minimal."
	}
	questions := make([]models.ScreeningQuestion, len(input.ScreeningQuestions))
	for index, question := range input.ScreeningQuestions {
		questions[index] = models.ScreeningQuestion{
			Question:       strings.TrimSpace(question.Question),
			Type:           question.Type,
			Required:       question.Required,
			ExpectedAnswer: optional(question.ExpectedAnswer),
			Position:       index,
		}
		if !matching.ValidExpected(question.Type, questions[index].ExpectedAnswer) {
			fields[fmt.Sprintf("screening_questions[%d].expected_answer", index)] = "Jawaban yang diharapkan tidak valid."
		}
	}
	var categoryCount int64
	if err := handler.DB.Model(&models.Category{}).Where("id = ?", input.CategoryID).Count(&categoryCount).Error; err != nil {
		httpx.ServerError(ctx, err)
		return false
	}
	if categoryCount == 0 {
		fields["category_id"] = "Kategori tidak ditemukan."
	}
	if len(fields) > 0 {
		httpx.FieldErrors(ctx, fields)
		return false
	}

	job.CategoryID = input.CategoryID
	job.Title = strings.TrimSpace(input.Title)
	job.Type, job.Level, job.WorkMode = input.Type, input.Level, input.WorkMode
	job.Location = optional(input.Location)
	job.SalaryMin, job.SalaryMax, job.ShowSalary = input.SalaryMin, input.SalaryMax, input.ShowSalary
	job.Skills = normalizeSkills(input.Skills)
	job.Description, job.Requirements = strings.TrimSpace(input.Description), strings.TrimSpace(input.Requirements)
	job.Deadline = deadline
	job.Status = input.Status
	job.Questions = questions
	return true
}

func (handler *Handler) saveJob(ctx *gin.Context, job models.Job, status int) {
	questions := job.Questions
	err := handler.DB.Transaction(func(tx *gorm.DB) error {
		query := tx.Omit(clause.Associations, "moderation_note", "views")
		if job.ID == 0 {
			query = tx.Omit(clause.Associations)
		}
		if err := query.Save(&job).Error; err != nil {
			return err
		}
		if err := tx.Where("job_id = ?", job.ID).Delete(&models.ScreeningQuestion{}).Error; err != nil {
			return err
		}
		for index := range questions {
			questions[index].JobID = job.ID
		}
		if len(questions) > 0 {
			if err := tx.Create(&questions).Error; err != nil {
				return err
			}
		}
		if job.Status == "closed" {
			return expireInvitations(tx, job.ID)
		}
		return nil
	})
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	handler.respondCompanyJob(ctx, status, job.ID)
}

func (handler *Handler) findOwnedJob(ctx *gin.Context) (models.Job, bool) {
	var job models.Job
	id, ok := httpx.ParamID(ctx)
	if !ok {
		return job, false
	}
	err := handler.DB.Where("id = ? AND company_id = ?", id, middleware.CurrentUser(ctx).Company.ID).First(&job).Error
	if err != nil {
		httpx.FindError(ctx, err, "Lowongan tidak ditemukan.")
		return job, false
	}
	return job, true
}

func (handler *Handler) respondCompanyJob(ctx *gin.Context, status int, jobID uint) {
	var job models.Job
	err := withJobRelations(handler.DB).Preload("Questions", orderedQuestions).First(&job, jobID).Error
	if err != nil {
		httpx.FindError(ctx, err, "Lowongan tidak ditemukan.")
		return
	}
	views, err := handler.managedJobs([]models.Job{job})
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	ctx.JSON(status, companyJobDetail{
		managedJob:         views[0],
		CategoryID:         job.CategoryID,
		Description:        job.Description,
		Requirements:       job.Requirements,
		ScreeningQuestions: job.Questions,
	})
}

func (handler *Handler) respondManagedJobs(ctx *gin.Context, query *gorm.DB) {
	page, err := httpx.Paginate[models.Job](ctx, query, listPageSize)
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	views, err := handler.managedJobs(page.Data)
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	ctx.JSON(200, httpx.WithData(page, views))
}

func orderedQuestions(query *gorm.DB) *gorm.DB {
	return query.Order("position")
}
