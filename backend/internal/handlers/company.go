package handlers

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"karsa/internal/httpx"
	"karsa/internal/middleware"
	"karsa/internal/models"
)

type companyInput struct {
	Name        string `json:"name" binding:"required,notblank,max=150"`
	Industry    string `json:"industry" binding:"max=100"`
	Size        string `json:"size" binding:"omitempty,oneof=1-10 11-50 51-200 201-500 >500"`
	Location    string `json:"location" binding:"max=100"`
	Website     string `json:"website" binding:"omitempty,http_url,max=255"`
	Description string `json:"description" binding:"max=5000"`
}

type companyView struct {
	models.Company
	LogoURL *string `json:"logo_url"`
}

type topJob struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Applicants int64  `json:"applicants"`
}

func (handler *Handler) ShowCompanyProfile(ctx *gin.Context) {
	handler.respondCompany(ctx, middleware.CurrentUser(ctx).Company.ID)
}

func (handler *Handler) UpdateCompanyProfile(ctx *gin.Context) {
	var input companyInput
	if !httpx.BindJSON(ctx, &input) {
		return
	}
	company := middleware.CurrentUser(ctx).Company
	changes := map[string]any{
		"name":        strings.TrimSpace(input.Name),
		"industry":    optional(input.Industry),
		"size":        optional(input.Size),
		"location":    optional(input.Location),
		"website":     optional(input.Website),
		"description": optional(input.Description),
	}
	if company.Status == "rejected" {
		changes["status"] = "pending"
		changes["rejection_reason"] = nil
	}
	if err := handler.DB.Model(&models.Company{}).Where("id = ?", company.ID).Updates(changes).Error; err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	handler.respondCompany(ctx, company.ID)
}

func (handler *Handler) UploadLogo(ctx *gin.Context) {
	company := middleware.CurrentUser(ctx).Company
	path, ok := handler.saveUpload(ctx, logoRule)
	if !ok {
		return
	}
	if err := handler.DB.Model(&models.Company{}).Where("id = ?", company.ID).Update("logo_path", path).Error; err != nil {
		handler.removeFile(&path)
		httpx.ServerError(ctx, err)
		return
	}
	handler.removeFile(company.LogoPath)
	handler.respondCompany(ctx, company.ID)
}

func (handler *Handler) CompanyDashboard(ctx *gin.Context) {
	company := middleware.CurrentUser(ctx).Company
	companyJobs := handler.DB.Model(&models.Job{}).Select("id").Where("company_id = ?", company.ID)

	var openJobs, total, recent, pendingInvitations int64
	for _, counter := range []struct {
		target *int64
		model  any
		where  string
		values []any
	}{
		{&openJobs, &models.Job{}, "company_id = ? AND status = ?", []any{company.ID, "open"}},
		{&total, &models.Application{}, "job_id IN (?)", []any{companyJobs}},
		{&recent, &models.Application{}, "job_id IN (?) AND created_at >= ?", []any{companyJobs, time.Now().AddDate(0, 0, -7)}},
		{&pendingInvitations, &models.Invitation{}, "company_id = ? AND status = ?", []any{company.ID, "pending"}},
	} {
		if err := handler.DB.Model(counter.model).Where(counter.where, counter.values...).Count(counter.target).Error; err != nil {
			httpx.ServerError(ctx, err)
			return
		}
	}
	byStatus, err := countByStatus(handler.DB.Model(&models.Application{}).Where("job_id IN (?)", companyJobs), "status", applicationStatuses...)
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	weekly, err := weeklyApplications(handler.DB.Model(&models.Application{}).Where("job_id IN (?)", companyJobs))
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	topJobs := []topJob{}
	err = handler.DB.Model(&models.Job{}).
		Select("jobs.id, jobs.title, COUNT(*) AS applicants").
		Joins("JOIN applications ON applications.job_id = jobs.id").
		Where("jobs.company_id = ?", company.ID).
		Group("jobs.id").Order("applicants DESC, jobs.id DESC").Limit(5).Scan(&topJobs).Error
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	ctx.JSON(200, gin.H{
		"open_jobs":           openJobs,
		"total_applicants":    total,
		"new_applicants_7d":   recent,
		"pending_invitations": pendingInvitations,
		"by_status":           byStatus,
		"weekly":              weekly,
		"top_jobs":            topJobs,
	})
}

func (handler *Handler) respondCompany(ctx *gin.Context, id uint) {
	var company models.Company
	if err := handler.DB.First(&company, id).Error; err != nil {
		httpx.FindError(ctx, err, "Perusahaan tidak ditemukan.")
		return
	}
	ctx.JSON(200, companyView{Company: company, LogoURL: handler.logoURL(company.LogoPath)})
}
