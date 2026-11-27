package handlers

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"karsa/internal/httpx"
	"karsa/internal/middleware"
	"karsa/internal/models"
)

const dailyInvitationLimit = 20

var errInvitationLimit = errors.New("kuota undangan habis")

type talentView struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url"`
	publicProfile
	InvitedJobIDs []uint `json:"invited_job_ids"`
}

type invitationInput struct {
	UserID  uint   `json:"user_id" binding:"required"`
	JobID   uint   `json:"job_id" binding:"required"`
	Message string `json:"message" binding:"max=500"`
}

type sentInvitation struct {
	ID        uint      `json:"id"`
	Candidate gin.H     `json:"candidate"`
	Job       gin.H     `json:"job"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type sentInvitationPage struct {
	httpx.Page[sentInvitation]
	RemainingToday int64 `json:"remaining_today"`
}

func (handler *Handler) discoverableSeekers() *gorm.DB {
	return handler.DB.Model(&models.User{}).Preload("Profile").
		Joins("JOIN seeker_profiles ON seeker_profiles.user_id = users.id").
		Where("users.role = ? AND users.is_active = ? AND seeker_profiles.is_discoverable = ?", "seeker", true, true)
}

func (handler *Handler) ListTalents(ctx *gin.Context) {
	query := handler.discoverableSeekers()
	for _, skill := range splitList(ctx.Query("skills")) {
		query = hasSkill(query, skill)
	}
	if location := ctx.Query("location"); location != "" {
		query = query.Where("seeker_profiles.location LIKE ?", likePattern(location))
	}
	if education := ctx.Query("education"); education != "" {
		query = query.Where("seeker_profiles.education = ?", education)
	}
	for parameter, condition := range map[string]string{
		"grad_from":      "seeker_profiles.graduation_year >= ?",
		"grad_to":        "seeker_profiles.graduation_year <= ?",
		"max_experience": "seeker_profiles.experience_years <= ?",
	} {
		if number, err := strconv.Atoi(ctx.Query(parameter)); err == nil {
			query = query.Where(condition, number)
		}
	}
	page, err := httpx.Paginate[models.User](ctx, query.Order("seeker_profiles.updated_at DESC, users.id DESC"), listPageSize)
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}

	userIDs := make([]uint, len(page.Data))
	for index, user := range page.Data {
		userIDs[index] = user.ID
	}
	var invitations []models.Invitation
	err = handler.DB.Where("company_id = ? AND user_id IN ?", middleware.CurrentUser(ctx).Company.ID, userIDs).Find(&invitations).Error
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	talents := make([]talentView, len(page.Data))
	for index, user := range page.Data {
		talents[index] = talentView{ID: user.ID, Name: user.Name, AvatarURL: user.AvatarURL, publicProfile: newPublicProfile(*user.Profile), InvitedJobIDs: []uint{}}
		for _, invitation := range invitations {
			if invitation.UserID == user.ID {
				talents[index].InvitedJobIDs = append(talents[index].InvitedJobIDs, invitation.JobID)
			}
		}
	}
	ctx.JSON(200, httpx.WithData(page, talents))
}

func (handler *Handler) CreateInvitation(ctx *gin.Context) {
	var input invitationInput
	if !httpx.BindJSON(ctx, &input) {
		return
	}
	company := middleware.CurrentUser(ctx).Company

	var job models.Job
	if err := handler.publicJobs().Where("jobs.id = ? AND jobs.company_id = ?", input.JobID, company.ID).First(&job).Error; err != nil {
		httpx.FindError(ctx, err, "Lowongan tidak ditemukan atau sudah ditutup.")
		return
	}
	var candidate models.User
	if err := handler.discoverableSeekers().Where("users.id = ?", input.UserID).First(&candidate).Error; err != nil {
		httpx.FindError(ctx, err, "Kandidat tidak ditemukan.")
		return
	}
	var applied int64
	if err := handler.DB.Model(&models.Application{}).Where("job_id = ? AND user_id = ?", job.ID, candidate.ID).Count(&applied).Error; err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	if applied > 0 {
		httpx.Error(ctx, 409, "Kandidat sudah melamar lowongan ini.")
		return
	}

	invitation := models.Invitation{CompanyID: company.ID, JobID: job.ID, UserID: candidate.ID, Message: optional(input.Message), Status: "pending"}
	err := handler.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&models.Company{}, company.ID).Error; err != nil {
			return err
		}
		remaining, err := remainingInvitations(tx, company.ID)
		if err != nil {
			return err
		}
		if remaining <= 0 {
			return errInvitationLimit
		}
		if err := tx.Create(&invitation).Error; err != nil {
			return err
		}
		return handler.Mailer.InvitationReceived(tx, candidate.Email, candidate.Name, job.Title, company.Name, invitation.Message)
	})
	switch {
	case err == nil:
		invitation.Job, invitation.User = &job, &candidate
		ctx.JSON(201, gin.H{"invitation": newSentInvitation(invitation)})
	case errors.Is(err, errInvitationLimit):
		httpx.Error(ctx, 429, "Kuota undangan hari ini sudah habis (20 per hari).")
	case errors.Is(err, gorm.ErrDuplicatedKey):
		httpx.Error(ctx, 409, "Kandidat sudah diundang ke lowongan ini.")
	default:
		httpx.ServerError(ctx, err)
	}
}

func (handler *Handler) ListSentInvitations(ctx *gin.Context) {
	company := middleware.CurrentUser(ctx).Company
	query := handler.DB.Preload("Job").Preload("User.Profile").Where("company_id = ?", company.ID)
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	page, err := httpx.Paginate[models.Invitation](ctx, query.Order("created_at DESC, id DESC"), listPageSize)
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	remaining, err := remainingInvitations(handler.DB, company.ID)
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	views := make([]sentInvitation, len(page.Data))
	for index, invitation := range page.Data {
		views[index] = newSentInvitation(invitation)
	}
	ctx.JSON(200, sentInvitationPage{Page: httpx.WithData(page, views), RemainingToday: remaining})
}

func remainingInvitations(db *gorm.DB, companyID uint) (int64, error) {
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	var sent int64
	err := db.Model(&models.Invitation{}).Where("company_id = ? AND created_at >= ?", companyID, midnight).Count(&sent).Error
	return max(0, dailyInvitationLimit-sent), err
}

func newSentInvitation(invitation models.Invitation) sentInvitation {
	candidate := newCandidateSummary(*invitation.User)
	return sentInvitation{
		ID:        invitation.ID,
		Candidate: gin.H{"name": candidate.Name, "headline": candidate.Headline, "avatar_url": candidate.AvatarURL},
		Job:       gin.H{"id": invitation.Job.ID, "title": invitation.Job.Title},
		Status:    invitation.Status,
		CreatedAt: invitation.CreatedAt,
	}
}
