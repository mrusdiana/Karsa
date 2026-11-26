package handlers

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"karsa/internal/httpx"
	"karsa/internal/models"
	"karsa/internal/util"
)

type categoryInput struct {
	Name string `json:"name" binding:"required,notblank,max=60"`
}

type categoryView struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	JobCount int64  `json:"job_count"`
}

func (handler *Handler) ListCategories(ctx *gin.Context) {
	categories, err := handler.categoriesWithCounts()
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	ctx.JSON(200, categories)
}

func (handler *Handler) CreateCategory(ctx *gin.Context) {
	handler.saveCategory(ctx, models.Category{}, 201)
}

func (handler *Handler) UpdateCategory(ctx *gin.Context) {
	if category, ok := handler.findCategory(ctx); ok {
		handler.saveCategory(ctx, category, 200)
	}
}

func (handler *Handler) DeleteCategory(ctx *gin.Context) {
	category, ok := handler.findCategory(ctx)
	if !ok {
		return
	}
	err := handler.DB.Delete(&category).Error
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		httpx.Error(ctx, 409, "Kategori masih dipakai lowongan.")
		return
	}
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	ctx.Status(204)
}

func (handler *Handler) categoriesWithCounts() ([]categoryView, error) {
	categories := []categoryView{}
	counts := handler.publicJobs().Select("jobs.category_id, COUNT(*) AS job_count").Group("jobs.category_id")
	err := handler.DB.Table("categories").
		Select("categories.id, categories.name, categories.slug, COALESCE(counts.job_count, 0) AS job_count").
		Joins("LEFT JOIN (?) AS counts ON counts.category_id = categories.id", counts).
		Order("categories.name").Scan(&categories).Error
	return categories, err
}

func (handler *Handler) saveCategory(ctx *gin.Context, category models.Category, status int) {
	var input categoryInput
	if !httpx.BindJSON(ctx, &input) {
		return
	}
	category.Name = strings.TrimSpace(input.Name)
	category.Slug = util.Slugify(category.Name)
	err := handler.DB.Save(&category).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		httpx.FieldErrors(ctx, map[string]string{"name": "Kategori sudah ada."})
		return
	}
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	ctx.JSON(status, category)
}

func (handler *Handler) findCategory(ctx *gin.Context) (models.Category, bool) {
	var category models.Category
	id, ok := httpx.ParamID(ctx)
	if !ok {
		return category, false
	}
	if err := handler.DB.First(&category, id).Error; err != nil {
		httpx.FindError(ctx, err, "Kategori tidak ditemukan.")
		return category, false
	}
	return category, true
}
