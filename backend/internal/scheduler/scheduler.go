package scheduler

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"

	"karsa/internal/models"
)

func Run(ctx context.Context, db *gorm.DB) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for {
		if err := Sweep(db); err != nil {
			log.Println("scheduler:", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func Sweep(db *gorm.DB) error {
	now := time.Now()
	err := db.Model(&models.Job{}).
		Where("status = ? AND deadline < ?", "open", now.Format(time.DateOnly)).
		Update("status", "closed").Error
	if err != nil {
		return err
	}
	closedJobs := db.Model(&models.Job{}).Select("id").Where("status = ?", "closed")
	return db.Model(&models.Invitation{}).
		Where("status = ?", "pending").
		Where("job_id IN (?) OR created_at < ?", closedJobs, now.Add(-models.InvitationLifetime)).
		Update("status", "expired").Error
}
