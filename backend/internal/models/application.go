package models

import "time"

type SkillScore struct {
	Points  int      `json:"points"`
	Matched []string `json:"matched"`
	Missing []string `json:"missing"`
}

type LevelScore struct {
	Points int    `json:"points"`
	Note   string `json:"note"`
}

type ScreeningScore struct {
	Points int `json:"points"`
	Passed int `json:"passed"`
	Total  int `json:"total"`
}

type ScoreBreakdown struct {
	Skill     SkillScore     `json:"skill"`
	Level     LevelScore     `json:"level"`
	Screening ScreeningScore `json:"screening"`
}

type Application struct {
	ID                uint `gorm:"primaryKey"`
	JobID             uint
	UserID            uint
	CVPath            string `gorm:"column:cv_path"`
	CoverLetter       *string
	Status            string
	Source            string
	MatchScore        int
	ScoreBreakdown    ScoreBreakdown `gorm:"serializer:json"`
	MeetsRequirements bool
	CompanyNote       *string
	CreatedAt         time.Time
	UpdatedAt         time.Time

	Job     *Job
	User    *User
	Answers []ApplicationAnswer
	Events  []ApplicationEvent
}

type ApplicationAnswer struct {
	ID            uint   `gorm:"primaryKey" json:"-"`
	ApplicationID uint   `json:"-"`
	Question      string `json:"question"`
	Required      bool   `json:"required"`
	Answer        string `json:"answer"`
	Passed        bool   `json:"passed"`
}

type ApplicationEvent struct {
	ID            uint `gorm:"primaryKey"`
	ApplicationID uint
	FromStatus    *string
	ToStatus      string
	ChangedBy     uint
	CreatedAt     time.Time
}
