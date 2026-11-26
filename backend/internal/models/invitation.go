package models

import "time"

const InvitationLifetime = 14 * 24 * time.Hour

type Invitation struct {
	ID        uint `gorm:"primaryKey"`
	CompanyID uint
	JobID     uint
	UserID    uint
	Message   *string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time

	Job  *Job
	User *User
}

func (invitation Invitation) ExpiresAt() time.Time {
	return invitation.CreatedAt.Add(InvitationLifetime)
}
