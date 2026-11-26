package models

import "time"

type Job struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CompanyID      uint      `json:"-"`
	CategoryID     uint      `json:"category_id"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	Type           string    `json:"type"`
	Level          string    `json:"level"`
	WorkMode       string    `json:"work_mode"`
	Location       *string   `json:"location"`
	SalaryMin      *uint     `json:"salary_min"`
	SalaryMax      *uint     `json:"salary_max"`
	ShowSalary     bool      `json:"show_salary"`
	Skills         []string  `gorm:"serializer:json" json:"skills"`
	Description    string    `json:"description"`
	Requirements   string    `json:"requirements"`
	Deadline       time.Time `gorm:"type:date" json:"deadline"`
	Status         string    `json:"status"`
	ModerationNote *string   `json:"-"`
	Views          uint      `json:"views"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Company   *Company            `json:"company,omitempty"`
	Category  *Category           `json:"category,omitempty"`
	Questions []ScreeningQuestion `json:"-"`
}
