package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"karsa/internal/httpx"
	"karsa/internal/middleware"
	"karsa/internal/models"
)

type savedJobView struct {
	jobSummary
	IsOpen bool `json:"is_open"`
}

func (handler *Handler) ListSavedJobs(ctx *gin.Context) {
	user := middleware.CurrentUser(ctx)
	query := withJobRelations(handler.DB).
		Joins("JOIN saved_jobs ON saved_jobs.job_id = jobs.id").
		Where("saved_jobs.user_id = ?", user.ID).
		Order("saved_jobs.created_at DESC, jobs.id DESC")
	page, err := httpx.Paginate[models.Job](ctx, query, publicPageSize)
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	summaries, err := handler.publicSummaries(page.Data, user)
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	var openIDs []uint
	if err := handler.publicJobs().Where("jobs.id IN ?", idsOf(page.Data)).Pluck("jobs.id", &openIDs).Error; err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	isOpen := map[uint]bool{}
	for _, id := range openIDs {
		isOpen[id] = true
	}
	views := make([]savedJobView, len(summaries))
	for index, summary := range summaries {
		views[index] = savedJobView{jobSummary: summary, IsOpen: isOpen[summary.ID]}
	}
	ctx.JSON(200, httpx.WithData(page, views))
}

func (handler *Handler) SaveJob(ctx *gin.Context) {
	id, ok := httpx.ParamID(ctx)
	if !ok {
		return
	}
	var job models.Job
	if err := handler.publicJobs().Where("jobs.id = ?", id).First(&job).Error; err != nil {
		httpx.FindError(ctx, err, "Lowongan tidak ditemukan.")
		return
	}
	saved := models.SavedJob{UserID: middleware.CurrentUser(ctx).ID, JobID: id}
	if err := handler.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&saved).Error; err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	ctx.Status(204)
}

func (handler *Handler) UnsaveJob(ctx *gin.Context) {
	id, ok := httpx.ParamID(ctx)
	if !ok {
		return
	}
	err := handler.DB.Where("user_id = ? AND job_id = ?", middleware.CurrentUser(ctx).ID, id).Delete(&models.SavedJob{}).Error
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	ctx.Status(204)
}
