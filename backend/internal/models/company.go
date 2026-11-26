package models

import "time"

type Company struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	UserID          uint       `json:"-"`
	Name            string     `json:"name"`
	Slug            string     `json:"slug"`
	LogoPath        *string    `json:"-"`
	Industry        *string    `json:"industry"`
	Size            *string    `json:"size"`
	Location        *string    `json:"location"`
	Website         *string    `json:"website"`
	Description     *string    `json:"description"`
	Status          string     `json:"status"`
	RejectionReason *string    `json:"rejection_reason"`
	VerifiedAt      *time.Time `json:"verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	User *User `json:"-"`
}
