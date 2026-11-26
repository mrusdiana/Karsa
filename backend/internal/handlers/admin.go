package handlers

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"karsa/internal/httpx"
	"karsa/internal/middleware"
	"karsa/internal/models"
)

type reasonInput struct {
	Reason string `json:"reason" binding:"required,notblank,max=500"`
}

type userStatusInput struct {
	IsActive *bool `json:"is_active" binding:"required"`
}

type adminCompany struct {
	companyView
	Owner gin.H `json:"owner"`
}

type nameCount struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type topCompany struct {
	Name         string `json:"name"`
	Applications int64  `json:"applications"`
}

func (handler *Handler) AdminDashboard(ctx *gin.Context) {
	users, err := countByStatus(handler.DB.Model(&models.User{}), "role", "seeker", "company", "admin")
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	companies, err := countByStatus(handler.DB.Model(&models.Company{}), "status", "pending", "verified", "rejected")
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	var openJobs, applications int64
	if err := handler.publicJobs().Count(&openJobs).Error; err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	if err := handler.DB.Model(&models.Application{}).Count(&applications).Error; err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	weekly, err := weeklyApplications(handler.DB.Model(&models.Application{}))
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	jobsByCategory := []nameCount{}
	err = handler.publicJobs().
		Select("categories.name, COUNT(*) AS count").
		Joins("JOIN categories ON categories.id = jobs.category_id").
		Group("categories.id").Order("count DESC, categories.name").Scan(&jobsByCategory).Error
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	topCompanies := []topCompany{}
	err = handler.DB.Model(&models.Application{}).
		Select("companies.name, COUNT(*) AS applications").
		Joins("JOIN jobs ON jobs.id = applications.job_id").
		Joins("JOIN companies ON companies.id = jobs.company_id").
		Group("companies.id").Order("applications DESC, companies.name").Limit(5).Scan(&topCompanies).Error
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	ctx.JSON(200, gin.H{
		"users":            users,
		"companies":        companies,
		"open_jobs":        openJobs,
		"applications":     applications,
		"weekly":           weekly,
		"jobs_by_category": jobsByCategory,
		"top_companies":    topCompanies,
	})
}

func (handler *Handler) AdminListCompanies(ctx *gin.Context) {
	query := handler.DB.Preload("User")
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	page, err := httpx.Paginate[models.Company](ctx, query.Order("created_at DESC, id DESC"), listPageSize)
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	views := make([]adminCompany, len(page.Data))
	for index, company := range page.Data {
		views[index] = handler.adminCompany(company)
	}
	ctx.JSON(200, httpx.WithData(page, views))
}

func (handler *Handler) AdminShowCompany(ctx *gin.Context) {
	if company, ok := handler.findCompany(ctx); ok {
		ctx.JSON(200, handler.adminCompany(company))
	}
}

func (handler *Handler) VerifyCompany(ctx *gin.Context) {
	company, ok := handler.findCompany(ctx)
	if !ok {
		return
	}
	changes := map[string]any{"status": "verified", "verified_at": time.Now(), "rejection_reason": nil}
	handler.decideCompany(ctx, company, changes, func(tx *gorm.DB) error {
		return handler.Mailer.CompanyVerified(tx, company.User.Email, company.Name)
	})
}

func (handler *Handler) RejectCompany(ctx *gin.Context) {
	company, ok := handler.findCompany(ctx)
	if !ok {
		return
	}
	var input reasonInput
	if !httpx.BindJSON(ctx, &input) {
		return
	}
	reason := strings.TrimSpace(input.Reason)
	changes := map[string]any{"status": "rejected", "verified_at": nil, "rejection_reason": reason}
	handler.decideCompany(ctx, company, changes, func(tx *gorm.DB) error {
		return handler.Mailer.CompanyRejected(tx, company.User.Email, company.Name, reason)
	})
}

func (handler *Handler) decideCompany(ctx *gin.Context, company models.Company, changes map[string]any, notify func(*gorm.DB) error) {
	isChanged := false
	err := handler.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.Company{}).Where("id = ? AND status <> ?", company.ID, changes["status"]).Updates(changes)
		if result.Error != nil || result.RowsAffected == 0 {
			return result.Error
		}
		isChanged = true
		return notify(tx)
	})
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	if !isChanged {
		httpx.Error(ctx, 409, "Status sudah berubah, muat ulang halaman.")
		return
	}
	handler.AdminShowCompany(ctx)
}

func (handler *Handler) AdminListUsers(ctx *gin.Context) {
	query := handler.DB.Preload("Company")
	if keyword := ctx.Query("q"); keyword != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", likePattern(keyword), likePattern(keyword))
	}
	if role := ctx.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}
	page, err := httpx.Paginate[models.User](ctx, query.Order("created_at DESC, id DESC"), listPageSize)
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	ctx.JSON(200, page)
}

func (handler *Handler) AdminUpdateUser(ctx *gin.Context) {
	id, ok := httpx.ParamID(ctx)
	if !ok {
		return
	}
	var input userStatusInput
	if !httpx.BindJSON(ctx, &input) {
		return
	}
	if id == middleware.CurrentUser(ctx).ID && !*input.IsActive {
		httpx.Error(ctx, 422, "Kamu tidak bisa menonaktifkan akunmu sendiri.")
		return
	}
	var user models.User
	if err := handler.DB.Preload("Company").First(&user, id).Error; err != nil {
		httpx.FindError(ctx, err, "Pengguna tidak ditemukan.")
		return
	}
	if err := handler.DB.Model(&models.User{}).Where("id = ?", id).Update("is_active", *input.IsActive).Error; err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	user.IsActive = *input.IsActive
	ctx.JSON(200, user)
}

func (handler *Handler) AdminListJobs(ctx *gin.Context) {
	query := withJobRelations(handler.DB).Joins("JOIN companies ON companies.id = jobs.company_id")
	if status := ctx.Query("status"); status != "" {
		query = query.Where("jobs.status = ?", status)
	}
	if company := ctx.Query("company"); company != "" {
		query = query.Where("companies.name LIKE ?", likePattern(company))
	}
	handler.respondManagedJobs(ctx, query.Order("jobs.created_at DESC, jobs.id DESC"))
}

func (handler *Handler) AdminCloseJob(ctx *gin.Context) {
	job, ok := handler.findJob(ctx)
	if !ok {
		return
	}
	var input reasonInput
	if !httpx.BindJSON(ctx, &input) {
		return
	}
	if job.Status == "draft" {
		httpx.Error(ctx, 422, "Draft belum dibuka, jadi tidak bisa ditutup.")
		return
	}
	reason := strings.TrimSpace(input.Reason)
	err := handler.DB.Transaction(func(tx *gorm.DB) error {
		return closeJob(tx, job.ID, &reason)
	})
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	handler.respondCompanyJob(ctx, 200, job.ID)
}

func (handler *Handler) AdminDeleteJob(ctx *gin.Context) {
	if job, ok := handler.findJob(ctx); ok {
		handler.deleteJob(ctx, job)
	}
}

func (handler *Handler) findCompany(ctx *gin.Context) (models.Company, bool) {
	var company models.Company
	id, ok := httpx.ParamID(ctx)
	if !ok {
		return company, false
	}
	err := handler.DB.Preload("User").First(&company, id).Error
	if err != nil {
		httpx.FindError(ctx, err, "Perusahaan tidak ditemukan.")
	}
	return company, err == nil
}

func (handler *Handler) findJob(ctx *gin.Context) (models.Job, bool) {
	var job models.Job
	id, ok := httpx.ParamID(ctx)
	if !ok {
		return job, false
	}
	err := handler.DB.First(&job, id).Error
	if err != nil {
		httpx.FindError(ctx, err, "Lowongan tidak ditemukan.")
	}
	return job, err == nil
}

func (handler *Handler) adminCompany(company models.Company) adminCompany {
	return adminCompany{
		companyView: companyView{Company: company, LogoURL: handler.logoURL(company.LogoPath)},
		Owner:       gin.H{"name": company.User.Name, "email": company.User.Email, "is_active": company.User.IsActive},
	}
}
