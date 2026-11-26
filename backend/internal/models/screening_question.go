package models

type ScreeningQuestion struct {
	ID             uint    `gorm:"primaryKey" json:"id"`
	JobID          uint    `json:"-"`
	Question       string  `json:"question"`
	Type           string  `json:"type"`
	Required       bool    `json:"required"`
	ExpectedAnswer *string `json:"expected_answer"`
	Position       int     `json:"-"`
}
