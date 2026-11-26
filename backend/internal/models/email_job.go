package models

import "time"

type EmailPayload struct {
	Heading string   `json:"heading"`
	Lines   []string `json:"lines"`
	Button  string   `json:"button"`
	URL     string   `json:"url"`
}

type EmailJob struct {
	ID          uint `gorm:"primaryKey"`
	ToEmail     string
	Subject     string
	Template    string
	Payload     EmailPayload `gorm:"serializer:json"`
	Status      string
	Attempts    int
	AvailableAt time.Time
	SentAt      *time.Time
	LastError   *string
	CreatedAt   time.Time
}
