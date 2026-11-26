package httpx

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Page[Item any] struct {
	Data     []Item `json:"data"`
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
	Total    int64  `json:"total"`
	LastPage int    `json:"last_page"`
}

func Paginate[Row any](ctx *gin.Context, query *gorm.DB, perPage int) (Page[Row], error) {
	result := Page[Row]{Data: []Row{}, PerPage: perPage}
	query = query.Session(&gorm.Session{})
	if err := query.Model(new(Row)).Count(&result.Total).Error; err != nil {
		return result, err
	}
	result.LastPage = max(1, (int(result.Total)+perPage-1)/perPage)
	requested, _ := strconv.Atoi(ctx.Query("page"))
	result.Page = min(max(1, requested), result.LastPage)
	err := query.Offset((result.Page - 1) * perPage).Limit(perPage).Find(&result.Data).Error
	return result, err
}

func WithData[Row, Item any](source Page[Row], data []Item) Page[Item] {
	return Page[Item]{Data: data, Page: source.Page, PerPage: source.PerPage, Total: source.Total, LastPage: source.LastPage}
}
