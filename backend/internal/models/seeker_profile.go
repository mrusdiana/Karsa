package models

import "time"

type SeekerProfile struct {
	UserID          uint `gorm:"primaryKey;autoIncrement:false"`
	Headline        *string
	Phone           *string
	Location        *string
	Education       *string
	Major           *string
	GraduationYear  *int
	ExperienceYears int
	Skills          []string `gorm:"serializer:json"`
	Bio             *string
	PortfolioURL    *string `gorm:"column:portfolio_url"`
	CVPath          *string `gorm:"column:cv_path"`
	IsDiscoverable  bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (profile SeekerProfile) Missing() []string {
	missing := []string{}
	for _, field := range []struct {
		name   string
		filled bool
	}{
		{"headline", profile.Headline != nil},
		{"location", profile.Location != nil},
		{"education", profile.Education != nil},
		{"graduation_year", profile.GraduationYear != nil},
		{"skills", len(profile.Skills) > 0},
		{"cv", profile.CVPath != nil},
	} {
		if !field.filled {
			missing = append(missing, field.name)
		}
	}
	return missing
}
