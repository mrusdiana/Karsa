package models

import "time"

type SavedJob struct {
	UserID    uint `gorm:"primaryKey;autoIncrement:false"`
	JobID     uint `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time
}
