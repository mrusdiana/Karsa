package handlers

import (
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"

	"karsa/internal/config"
	"karsa/internal/mailer"
)

const (
	publicPageSize = 12
	listPageSize   = 20
)

type Handler struct {
	DB     *gorm.DB
	Config config.Config
	Mailer mailer.Mailer
}

type weekCount struct {
	WeekStart string `json:"week_start"`
	Count     int64  `json:"count"`
}

type statusCount struct {
	Status string
	Count  int64
}

func today() string {
	return time.Now().Format(time.DateOnly)
}

func optional(text string) *string {
	if text = strings.TrimSpace(text); text == "" {
		return nil
	}
	return &text
}

func splitList(text string) []string {
	items := []string{}
	for _, item := range strings.Split(text, ",") {
		if item = strings.TrimSpace(item); item != "" {
			items = append(items, item)
		}
	}
	return items
}

func likePattern(text string) string {
	escaped := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(strings.TrimSpace(text))
	return "%" + escaped + "%"
}

func normalizeSkills(skills []string) []string {
	result := []string{}
	for _, skill := range skills {
		skill = strings.TrimSpace(skill)
		isDuplicate := slices.ContainsFunc(result, func(existing string) bool { return strings.EqualFold(existing, skill) })
		if skill != "" && !isDuplicate {
			result = append(result, skill)
		}
	}
	return result
}

func hasSkill(query *gorm.DB, skill string) *gorm.DB {
	return query.Where("JSON_CONTAINS(LOWER(CAST(seeker_profiles.skills AS CHAR)), JSON_QUOTE(LOWER(?)))", skill)
}

func countByStatus(query *gorm.DB, column string, statuses ...string) (map[string]int64, error) {
	counts := map[string]int64{}
	for _, status := range statuses {
		counts[status] = 0
	}
	var rows []statusCount
	err := query.Select(column + " AS status, COUNT(*) AS count").Group(column).Scan(&rows).Error
	for _, row := range rows {
		counts[row.Status] = row.Count
	}
	return counts, err
}

func weeklyApplications(query *gorm.DB) ([]weekCount, error) {
	now := time.Now()
	daysSinceMonday := (int(now.Weekday()) + 6) % 7
	firstMonday := time.Date(now.Year(), now.Month(), now.Day()-daysSinceMonday-49, 0, 0, 0, 0, time.Local)

	var rows []weekCount
	err := query.
		Select("DATE_FORMAT(DATE_SUB(applications.created_at, INTERVAL WEEKDAY(applications.created_at) DAY), '%Y-%m-%d') AS week_start, COUNT(*) AS count").
		Where("applications.created_at >= ?", firstMonday).
		Group("week_start").Scan(&rows).Error

	weeks := make([]weekCount, 8)
	for index := range weeks {
		weeks[index].WeekStart = firstMonday.AddDate(0, 0, index*7).Format(time.DateOnly)
		for _, row := range rows {
			if row.WeekStart == weeks[index].WeekStart {
				weeks[index].Count = row.Count
			}
		}
	}
	return weeks, err
}
