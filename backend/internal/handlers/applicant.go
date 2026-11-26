package handlers

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"karsa/internal/httpx"
	"karsa/internal/middleware"
	"karsa/internal/models"
)

var (
	errStatusChanged    = errors.New("status sudah berubah")
	applicationStatuses = []string{"applied", "reviewed", "interview", "accepted", "rejected", "withdrawn"}
	companyMoves        = map[string][]string{
		"applied":   {"reviewed", "rejected"},
		"reviewed":  {"interview", "rejected"},
		"interview": {"accepted", "rejected"},
	}
)

type applicantInput struct {
	Status      *string `json:"status" binding:"omitempty,oneof=reviewed interview accepted rejected"`
	CompanyNote *string `json:"company_note" binding:"omitempty,max=1000"`
}

type candidateSummary struct {
	Name     string  `json:"name"`
	Headline *string `json:"headline"`
	Location *string `json:"location"`
}

type candidateDetail struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	contactProfile
}

type applicantRow struct {
	ID                uint             `json:"id"`
	Candidate         candidateSummary `json:"candidate"`
	MatchScore        int              `json:"match_score"`
	MeetsRequirements bool             `json:"meets_requirements"`
	Source            string           `json:"source"`
	Status            string           `json:"status"`
	CreatedAt         time.Time        `json:"created_at"`
}

type applicantPage struct {
	httpx.Page[applicantRow]
	Counts map[string]int64 `json:"counts"`
}

type applicantDetail struct {
	ID                uint                       `json:"id"`
	Job               gin.H                      `json:"job"`
	Candidate         candidateDetail            `json:"candidate"`
	MatchScore        int                        `json:"match_score"`
	MeetsRequirements bool                       `json:"meets_requirements"`
	Source            string                     `json:"source"`
	Status            string                     `json:"status"`
	CoverLetter       *string                    `json:"cover_letter"`
	Answers           []models.ApplicationAnswer `json:"answers"`
	ScoreBreakdown    models.ScoreBreakdown      `json:"score_breakdown"`
	CompanyNote       *string                    `json:"company_note"`
	CVURL             string                     `json:"cv_url"`
	Events            []eventView                `json:"events"`
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"updated_at"`
}

func (handler *Handler) ListApplicants(ctx *gin.Context) {
	job, ok := handler.findOwnedJob(ctx)
	if !ok {
		return
	}
	counts, err := countByStatus(handler.DB.Model(&models.Application{}).Where("job_id = ?", job.ID), "status", applicationStatuses...)
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}

	query := handler.DB.Preload("User.Profile").
		Joins("LEFT JOIN seeker_profiles ON seeker_profiles.user_id = applications.user_id").
		Where("applications.job_id = ?", job.ID)
	if statuses := splitList(ctx.Query("status")); len(statuses) > 0 {
		query = query.Where("applications.status IN ?", statuses)
	}
	if minimum, err := strconv.Atoi(ctx.Query("min_score")); err == nil {
		query = query.Where("applications.match_score >= ?", minimum)
	}
	for _, skill := range splitList(ctx.Query("skill")) {
		query = hasSkill(query, skill)
	}
	if meets, err := strconv.ParseBool(ctx.Query("meets")); err == nil {
		query = query.Where("applications.meets_requirements = ?", meets)
	}
	if source := ctx.Query("source"); source != "" {
		query = query.Where("applications.source = ?", source)
	}
	order := "applications.match_score DESC, applications.id DESC"
	if ctx.Query("sort") == "newest" {
		order = "applications.created_at DESC, applications.id DESC"
	}

	page, err := httpx.Paginate[models.Application](ctx, query.Order(order), listPageSize)
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	rows := make([]applicantRow, len(page.Data))
	for index, application := range page.Data {
		rows[index] = applicantRow{
			ID:                application.ID,
			Candidate:         newCandidateSummary(*application.User),
			MatchScore:        application.MatchScore,
			MeetsRequirements: application.MeetsRequirements,
			Source:            application.Source,
			Status:            application.Status,
			CreatedAt:         application.CreatedAt,
		}
	}
	ctx.JSON(200, applicantPage{Page: httpx.WithData(page, rows), Counts: counts})
}

func (handler *Handler) ShowApplicant(ctx *gin.Context) {
	application, ok := handler.findApplicant(ctx)
	if !ok {
		return
	}
	ctx.JSON(200, applicantDetail{
		ID:  application.ID,
		Job: gin.H{"id": application.Job.ID, "title": application.Job.Title, "slug": application.Job.Slug},
		Candidate: candidateDetail{
			Name:           application.User.Name,
			Email:          application.User.Email,
			contactProfile: newContactProfile(*application.User.Profile),
		},
		MatchScore:        application.MatchScore,
		MeetsRequirements: application.MeetsRequirements,
		Source:            application.Source,
		Status:            application.Status,
		CoverLetter:       application.CoverLetter,
		Answers:           application.Answers,
		ScoreBreakdown:    application.ScoreBreakdown,
		CompanyNote:       application.CompanyNote,
		CVURL:             fmt.Sprintf("%s/api/company/applications/%d/cv", handler.Config.APIURL, application.ID),
		Events:            eventViews(application, "candidate", "company"),
		CreatedAt:         application.CreatedAt,
		UpdatedAt:         application.UpdatedAt,
	})
}

func (handler *Handler) UpdateApplicant(ctx *gin.Context) {
	application, ok := handler.findApplicant(ctx)
	if !ok {
		return
	}
	var input applicantInput
	if !httpx.BindJSON(ctx, &input) {
		return
	}
	if input.Status != nil && !slices.Contains(companyMoves[application.Status], *input.Status) {
		httpx.FieldErrors(ctx, map[string]string{"status": "Perpindahan status tidak diizinkan."})
		return
	}
	user := middleware.CurrentUser(ctx)
	err := handler.DB.Transaction(func(tx *gorm.DB) error {
		if input.CompanyNote != nil {
			if err := tx.Model(&models.Application{}).Where("id = ?", application.ID).Update("company_note", optional(*input.CompanyNote)).Error; err != nil {
				return err
			}
		}
		if input.Status == nil {
			return nil
		}
		if err := moveStatus(tx, application, *input.Status, user.ID); err != nil {
			return err
		}
		return handler.Mailer.StatusChanged(tx, application.User.Email, application.User.Name, application.Job.Title, user.Company.Name, *input.Status, application.ID)
	})
	if err != nil {
		respondMoveError(ctx, err)
		return
	}
	handler.ShowApplicant(ctx)
}

func (handler *Handler) DownloadApplicantCV(ctx *gin.Context) {
	if application, ok := handler.findApplicant(ctx); ok {
		handler.sendCV(ctx, application.CVPath, cvName(application.User.Name))
	}
}

func (handler *Handler) findApplicant(ctx *gin.Context) (models.Application, bool) {
	var application models.Application
	id, ok := httpx.ParamID(ctx)
	if !ok {
		return application, false
	}
	err := handler.DB.
		Preload("Job").Preload("User.Profile").Preload("Answers").Preload("Events", orderedEvents).
		Joins("JOIN jobs ON jobs.id = applications.job_id").
		Where("jobs.company_id = ?", middleware.CurrentUser(ctx).Company.ID).
		First(&application, id).Error
	if err != nil {
		httpx.FindError(ctx, err, "Lamaran tidak ditemukan.")
	}
	return application, err == nil
}

func moveStatus(tx *gorm.DB, application models.Application, status string, changedBy uint) error {
	result := tx.Model(&models.Application{}).
		Where("id = ? AND status = ?", application.ID, application.Status).
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errStatusChanged
	}
	return tx.Create(&models.ApplicationEvent{
		ApplicationID: application.ID,
		FromStatus:    &application.Status,
		ToStatus:      status,
		ChangedBy:     changedBy,
	}).Error
}

func respondMoveError(ctx *gin.Context, err error) {
	if errors.Is(err, errStatusChanged) {
		httpx.Error(ctx, 409, "Status sudah berubah, muat ulang halaman.")
		return
	}
	httpx.ServerError(ctx, err)
}

func newCandidateSummary(user models.User) candidateSummary {
	summary := candidateSummary{Name: user.Name}
	if user.Profile != nil {
		summary.Headline = user.Profile.Headline
		summary.Location = user.Profile.Location
	}
	return summary
}
