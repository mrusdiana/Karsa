package handlers

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"karsa/internal/httpx"
	"karsa/internal/matching"
	"karsa/internal/middleware"
	"karsa/internal/models"
)

const homeJobCount = 6

var (
	searchWord = regexp.MustCompile(`[\p{L}\p{N}]+`)
	jobOrders  = map[string]string{
		"newest":   "jobs.created_at DESC, jobs.id DESC",
		"salary":   "IF(jobs.show_salary, COALESCE(jobs.salary_max, jobs.salary_min, 0), 0) DESC, jobs.id DESC",
		"deadline": "jobs.deadline, jobs.id DESC",
	}
)

type companyProfile struct {
	companySummary
	Industry    *string `json:"industry"`
	Size        *string `json:"size"`
	Location    *string `json:"location"`
	Website     *string `json:"website"`
	Description *string `json:"description"`
}

type publicQuestion struct {
	ID       uint   `json:"id"`
	Question string `json:"question"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type matchPreview struct {
	Score     int                   `json:"score"`
	Breakdown models.ScoreBreakdown `json:"breakdown"`
}

type jobDetail struct {
	jobSummary
	Company            companyProfile   `json:"company"`
	Description        string           `json:"description"`
	Requirements       string           `json:"requirements"`
	ApplicantCount     int64            `json:"applicant_count"`
	Views              uint             `json:"views"`
	ScreeningQuestions []publicQuestion `json:"screening_questions"`
	MatchPreview       *matchPreview    `json:"match_preview,omitempty"`
}

type hiringCompany struct {
	Name     string  `json:"name"`
	Slug     string  `json:"slug"`
	LogoPath *string `json:"-"`
	LogoURL  *string `gorm:"-" json:"logo_url"`
	OpenJobs int64   `json:"open_jobs"`
}

func (handler *Handler) publicJobs() *gorm.DB {
	return handler.DB.Model(&models.Job{}).
		Joins("JOIN companies ON companies.id = jobs.company_id").
		Joins("JOIN users ON users.id = companies.user_id").
		Where("jobs.status = ? AND jobs.deadline >= ?", "open", today()).
		Where("companies.status = ? AND users.is_active = ?", "verified", true)
}

func (handler *Handler) Home(ctx *gin.Context) {
	categories, err := handler.categoriesWithCounts()
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	latest, err := handler.homeJobs(ctx, handler.publicJobs())
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	entryLevel, err := handler.homeJobs(ctx, handler.publicJobs().Where("jobs.level IN ?", []string{"internship", "entry"}))
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	companies := []hiringCompany{}
	err = handler.publicJobs().
		Select("companies.name, companies.slug, companies.logo_path, COUNT(*) AS open_jobs").
		Group("companies.id").Order("open_jobs DESC, companies.name").Limit(homeJobCount).Scan(&companies).Error
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	for index := range companies {
		companies[index].LogoURL = handler.logoURL(companies[index].LogoPath)
	}
	ctx.JSON(200, gin.H{"categories": categories, "latest": latest, "entry_level": entryLevel, "companies": companies})
}

func (handler *Handler) homeJobs(ctx *gin.Context, query *gorm.DB) ([]jobSummary, error) {
	var jobs []models.Job
	if err := withJobRelations(query).Order(jobOrders["newest"]).Limit(homeJobCount).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return handler.publicSummaries(jobs, middleware.CurrentUser(ctx))
}

func (handler *Handler) ListJobs(ctx *gin.Context) {
	query := handler.publicJobs()
	if keyword := strings.TrimSpace(ctx.Query("q")); keyword != "" {
		query = query.Where(handler.matchesKeyword(keyword))
	}
	if category := ctx.Query("category"); category != "" {
		query = query.Where("jobs.category_id IN (?)", handler.DB.Model(&models.Category{}).Select("id").Where("slug = ?", category))
	}
	for parameter, column := range map[string]string{"type": "jobs.type", "level": "jobs.level", "mode": "jobs.work_mode"} {
		if values := splitList(ctx.Query(parameter)); len(values) > 0 {
			query = query.Where(column+" IN ?", values)
		}
	}
	if location := ctx.Query("location"); location != "" {
		query = query.Where("jobs.location LIKE ?", likePattern(location))
	}
	if salary, err := strconv.ParseUint(ctx.Query("salary_min"), 10, 32); err == nil && salary > 0 {
		query = query.Where("jobs.show_salary AND COALESCE(jobs.salary_max, jobs.salary_min) >= ?", salary)
	}
	order, ok := jobOrders[ctx.Query("sort")]
	if !ok {
		order = jobOrders["newest"]
	}

	page, err := httpx.Paginate[models.Job](ctx, withJobRelations(query).Order(order), publicPageSize)
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	summaries, err := handler.publicSummaries(page.Data, middleware.CurrentUser(ctx))
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	ctx.JSON(200, httpx.WithData(page, summaries))
}

func (handler *Handler) matchesKeyword(keyword string) *gorm.DB {
	pattern := likePattern(keyword)
	conditions := handler.DB.
		Where("jobs.title LIKE ?", pattern).
		Or("companies.name LIKE ?", pattern).
		Or("LOWER(CAST(jobs.skills AS CHAR)) LIKE LOWER(?)", `%"`+strings.TrimPrefix(pattern, "%"))
	if words := searchWord.FindAllString(keyword, -1); len(words) > 0 {
		conditions = conditions.Or("MATCH(jobs.title, jobs.description) AGAINST (? IN BOOLEAN MODE)", "+"+strings.Join(words, "* +")+"*")
	}
	return conditions
}

func (handler *Handler) ShowJob(ctx *gin.Context) {
	var job models.Job
	err := withJobRelations(handler.publicJobs()).Preload("Questions", orderedQuestions).Where("jobs.slug = ?", ctx.Param("id")).First(&job).Error
	if err != nil {
		httpx.FindError(ctx, err, "Lowongan tidak ditemukan.")
		return
	}
	if err := handler.DB.Model(&models.Job{}).Where("id = ?", job.ID).UpdateColumn("views", gorm.Expr("views + 1")).Error; err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	viewer := middleware.CurrentUser(ctx)
	summaries, err := handler.publicSummaries([]models.Job{job}, viewer)
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	counts, err := handler.applicantCounts([]uint{job.ID})
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}

	detail := jobDetail{
		jobSummary:         summaries[0],
		Company:            handler.companyProfile(*job.Company),
		Description:        job.Description,
		Requirements:       job.Requirements,
		ApplicantCount:     counts[job.ID],
		Views:              job.Views + 1,
		ScreeningQuestions: make([]publicQuestion, len(job.Questions)),
	}
	for index, question := range job.Questions {
		detail.ScreeningQuestions[index] = publicQuestion{ID: question.ID, Question: question.Question, Type: question.Type, Required: question.Required}
	}
	if viewer != nil && viewer.Role == "seeker" {
		profile, err := handler.seekerProfile(viewer.ID)
		if err != nil {
			httpx.ServerError(ctx, err)
			return
		}
		score, breakdown := matching.Score(job, profile, nil, time.Now().Year())
		detail.MatchPreview = &matchPreview{Score: score, Breakdown: breakdown}
	}
	ctx.JSON(200, detail)
}

func (handler *Handler) ShowCompany(ctx *gin.Context) {
	var company models.Company
	err := handler.DB.
		Joins("JOIN users ON users.id = companies.user_id").
		Where("companies.slug = ? AND companies.status = ? AND users.is_active = ?", ctx.Param("slug"), "verified", true).
		First(&company).Error
	if err != nil {
		httpx.FindError(ctx, err, "Perusahaan tidak ditemukan.")
		return
	}
	var jobs []models.Job
	err = withJobRelations(handler.publicJobs()).Where("jobs.company_id = ?", company.ID).Order(jobOrders["newest"]).Find(&jobs).Error
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	summaries, err := handler.publicSummaries(jobs, middleware.CurrentUser(ctx))
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	ctx.JSON(200, struct {
		companyProfile
		Jobs []jobSummary `json:"jobs"`
	}{handler.companyProfile(company), summaries})
}

func (handler *Handler) companyProfile(company models.Company) companyProfile {
	return companyProfile{
		companySummary: handler.companySummary(company),
		Industry:       company.Industry,
		Size:           company.Size,
		Location:       company.Location,
		Website:        company.Website,
		Description:    company.Description,
	}
}
