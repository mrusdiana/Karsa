package handlers

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"karsa/internal/httpx"
	"karsa/internal/matching"
	"karsa/internal/middleware"
	"karsa/internal/models"
)

const maxAnswerLength = 500

var errInvitationUsed = errors.New("undangan sudah tidak berlaku")

type answerInput struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	Answer     string `json:"answer"`
}

type applyInput struct {
	CoverLetter  string        `json:"cover_letter" binding:"max=2000"`
	Answers      []answerInput `json:"answers" binding:"max=20,dive"`
	InvitationID *uint         `json:"invitation_id"`
}

type applicationSummary struct {
	ID        uint       `json:"id"`
	Job       jobSummary `json:"job"`
	Status    string     `json:"status"`
	Source    string     `json:"source"`
	HasUpdate bool       `json:"has_update"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type sentAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type eventView struct {
	FromStatus *string   `json:"from_status"`
	ToStatus   string    `json:"to_status"`
	By         string    `json:"by"`
	CreatedAt  time.Time `json:"created_at"`
}

type applicationDetail struct {
	applicationSummary
	CoverLetter *string      `json:"cover_letter"`
	Answers     []sentAnswer `json:"answers"`
	CVURL       string       `json:"cv_url"`
	Events      []eventView  `json:"events"`
}

func (handler *Handler) Apply(ctx *gin.Context) {
	id, ok := httpx.ParamID(ctx)
	if !ok {
		return
	}
	var input applyInput
	if !httpx.BindJSON(ctx, &input) {
		return
	}
	user := middleware.CurrentUser(ctx)

	var job models.Job
	err := withJobRelations(handler.DB).Preload("Questions", orderedQuestions).Where("status <> ?", "draft").First(&job, id).Error
	if err != nil {
		httpx.FindError(ctx, err, "Lowongan tidak ditemukan.")
		return
	}
	var applied int64
	if err := handler.DB.Model(&models.Application{}).Where("job_id = ? AND user_id = ?", job.ID, user.ID).Count(&applied).Error; err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	if applied > 0 {
		httpx.Error(ctx, 409, "Kamu sudah melamar lowongan ini")
		return
	}
	var available int64
	if err := handler.publicJobs().Where("jobs.id = ?", job.ID).Count(&available).Error; err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	if available == 0 {
		httpx.Error(ctx, 422, "Lowongan sudah ditutup")
		return
	}

	profile, ok := handler.currentProfile(ctx)
	if !ok {
		return
	}
	if missing := profile.Missing(); len(missing) > 0 {
		ctx.AbortWithStatusJSON(422, gin.H{
			"message": "Lengkapi profil dan CV sebelum melamar.",
			"errors":  gin.H{"profile": strings.Join(missing, ", ")},
			"missing": missing,
		})
		return
	}
	invitation, ok := handler.activeInvitation(ctx, input.InvitationID, job.ID, user.ID)
	if !ok {
		return
	}
	answers, fields := checkAnswers(job.Questions, input.Answers)
	if len(fields) > 0 {
		httpx.FieldErrors(ctx, fields)
		return
	}

	application := newApplication(job, profile, answers)
	application.UserID = user.ID
	application.CoverLetter = optional(input.CoverLetter)
	if invitation != nil {
		application.Source = "invited"
	}
	application.CVPath, err = handler.copyFile(*profile.CVPath, "applications")
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}

	err = handler.DB.Transaction(func(tx *gorm.DB) error {
		var lockedJob models.Job
		if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).Where("status = ?", "open").First(&lockedJob, job.ID).Error; err != nil {
			return err
		}
		if err := tx.Omit("Job", "User").Create(&application).Error; err != nil {
			return err
		}
		if invitation != nil {
			result := tx.Model(invitation).Where("status = ?", "pending").Update("status", "accepted")
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errInvitationUsed
			}
		}
		return handler.Mailer.ApplicationReceived(tx, user.Email, user.Name, job.Title, job.Company.Name, application.ID)
	})
	if err != nil {
		handler.removeFile(&application.CVPath)
	}
	switch {
	case err == nil:
		application.Job = &job
		ctx.JSON(201, gin.H{"application": handler.applicationSummary(application)})
	case errors.Is(err, gorm.ErrDuplicatedKey):
		httpx.Error(ctx, 409, "Kamu sudah melamar lowongan ini")
	case errors.Is(err, gorm.ErrRecordNotFound):
		httpx.Error(ctx, 422, "Lowongan sudah ditutup")
	case errors.Is(err, errInvitationUsed):
		httpx.Error(ctx, 422, "Undangan sudah tidak berlaku")
	default:
		httpx.ServerError(ctx, err)
	}
}

func (handler *Handler) activeInvitation(ctx *gin.Context, requestedID *uint, jobID, userID uint) (*models.Invitation, bool) {
	var invitations []models.Invitation
	err := handler.DB.
		Where("job_id = ? AND user_id = ? AND status = ?", jobID, userID, "pending").
		Where("created_at >= ?", time.Now().Add(-models.InvitationLifetime)).
		Limit(1).Find(&invitations).Error
	if err != nil {
		httpx.ServerError(ctx, err)
		return nil, false
	}
	if requestedID != nil && (len(invitations) == 0 || invitations[0].ID != *requestedID) {
		httpx.Error(ctx, 422, "Undangan sudah tidak berlaku")
		return nil, false
	}
	if len(invitations) == 0 {
		return nil, true
	}
	return &invitations[0], true
}

func checkAnswers(questions []models.ScreeningQuestion, inputs []answerInput) ([]models.ApplicationAnswer, map[string]string) {
	given := map[uint]string{}
	for _, input := range inputs {
		given[input.QuestionID] = strings.TrimSpace(input.Answer)
	}
	fields := map[string]string{}
	answers := make([]models.ApplicationAnswer, len(questions))
	for index, question := range questions {
		answer := given[question.ID]
		passed, valid := matching.Evaluate(question, answer)
		key := fmt.Sprintf("answers.%d", question.ID)
		switch {
		case question.Required && answer == "":
			fields[key] = "Wajib diisi."
		case utf8.RuneCountInString(answer) > maxAnswerLength:
			fields[key] = "Maksimal 500."
		case !valid:
			fields[key] = "Jawaban tidak valid."
		}
		answers[index] = models.ApplicationAnswer{Question: question.Question, Required: question.Required, Answer: answer, Passed: passed}
	}
	return answers, fields
}

func newApplication(job models.Job, profile models.SeekerProfile, answers []models.ApplicationAnswer) models.Application {
	passed := make([]bool, len(answers))
	meetsRequirements := true
	for index, answer := range answers {
		passed[index] = answer.Passed
		if answer.Required && !answer.Passed {
			meetsRequirements = false
		}
	}
	score, breakdown := matching.Score(job, profile, passed, time.Now().Year())
	return models.Application{
		JobID:             job.ID,
		UserID:            profile.UserID,
		Status:            "applied",
		Source:            "applied",
		MatchScore:        score,
		ScoreBreakdown:    breakdown,
		MeetsRequirements: meetsRequirements,
		Answers:           answers,
		Events:            []models.ApplicationEvent{{ToStatus: "applied", ChangedBy: profile.UserID}},
	}
}

func (handler *Handler) ListMyApplications(ctx *gin.Context) {
	query := handler.DB.Preload("Job.Company").Preload("Job.Category").Where("user_id = ?", middleware.CurrentUser(ctx).ID)
	if statuses := splitList(ctx.Query("status")); len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
	}
	page, err := httpx.Paginate[models.Application](ctx, query.Order("created_at DESC, id DESC"), listPageSize)
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	summaries := make([]applicationSummary, len(page.Data))
	for index, application := range page.Data {
		summaries[index] = handler.applicationSummary(application)
	}
	ctx.JSON(200, httpx.WithData(page, summaries))
}

func (handler *Handler) ShowMyApplication(ctx *gin.Context) {
	application, ok := handler.findMyApplication(ctx)
	if !ok {
		return
	}
	if application.HasUpdate {
		if err := handler.DB.Model(&models.Application{}).Where("id = ?", application.ID).Update("has_update", false).Error; err != nil {
			httpx.ServerError(ctx, err)
			return
		}
		application.HasUpdate = false
	}
	ctx.JSON(200, handler.applicationDetail(application))
}

func (handler *Handler) MyNotifications(ctx *gin.Context) {
	userID := middleware.CurrentUser(ctx).ID
	var applications, invitations int64
	if err := handler.DB.Model(&models.Application{}).Where("user_id = ? AND has_update = ?", userID, true).Count(&applications).Error; err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	err := handler.DB.Model(&models.Invitation{}).
		Where("user_id = ? AND status = ? AND created_at >= ?", userID, "pending", time.Now().Add(-models.InvitationLifetime)).
		Count(&invitations).Error
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	ctx.JSON(200, gin.H{"applications": applications, "invitations": invitations})
}

func (handler *Handler) DownloadMyApplicationCV(ctx *gin.Context) {
	if application, ok := handler.findMyApplication(ctx); ok {
		handler.sendCV(ctx, application.CVPath, cvName(middleware.CurrentUser(ctx).Name))
	}
}

func (handler *Handler) WithdrawApplication(ctx *gin.Context) {
	application, ok := handler.findMyApplication(ctx)
	if !ok {
		return
	}
	if application.Status != "applied" {
		httpx.Error(ctx, 422, "Lamaran yang sudah diproses tidak bisa ditarik.")
		return
	}
	err := handler.DB.Transaction(func(tx *gorm.DB) error {
		return moveStatus(tx, application, "withdrawn", application.UserID)
	})
	if err != nil {
		respondMoveError(ctx, err)
		return
	}
	handler.ShowMyApplication(ctx)
}

func (handler *Handler) findMyApplication(ctx *gin.Context) (models.Application, bool) {
	var application models.Application
	id, ok := httpx.ParamID(ctx)
	if !ok {
		return application, false
	}
	err := handler.DB.
		Preload("Job.Company").Preload("Job.Category").Preload("Answers").Preload("Events", orderedEvents).
		Where("user_id = ?", middleware.CurrentUser(ctx).ID).First(&application, id).Error
	if err != nil {
		httpx.FindError(ctx, err, "Lamaran tidak ditemukan.")
	}
	return application, err == nil
}

func (handler *Handler) applicationSummary(application models.Application) applicationSummary {
	return applicationSummary{
		ID:        application.ID,
		Job:       handler.publicSummary(*application.Job),
		Status:    application.Status,
		Source:    application.Source,
		HasUpdate: application.HasUpdate,
		CreatedAt: application.CreatedAt,
		UpdatedAt: application.UpdatedAt,
	}
}

func (handler *Handler) applicationDetail(application models.Application) applicationDetail {
	answers := make([]sentAnswer, len(application.Answers))
	for index, answer := range application.Answers {
		answers[index] = sentAnswer{Question: answer.Question, Answer: answer.Answer}
	}
	return applicationDetail{
		applicationSummary: handler.applicationSummary(application),
		CoverLetter:        application.CoverLetter,
		Answers:            answers,
		CVURL:              fmt.Sprintf("%s/api/me/applications/%d/cv", handler.Config.APIURL, application.ID),
		Events:             eventViews(application, "you", "company"),
	}
}

func eventViews(application models.Application, seekerLabel, companyLabel string) []eventView {
	views := make([]eventView, len(application.Events))
	for index, event := range application.Events {
		by := companyLabel
		if event.ChangedBy == application.UserID {
			by = seekerLabel
		}
		views[index] = eventView{FromStatus: event.FromStatus, ToStatus: event.ToStatus, By: by, CreatedAt: event.CreatedAt}
	}
	return views
}

func orderedEvents(query *gorm.DB) *gorm.DB {
	return query.Order("id")
}
