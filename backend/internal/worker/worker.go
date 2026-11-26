package worker

import (
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"karsa/internal/models"
)

const (
	interval    = 5 * time.Second
	maxAttempts = 3
)

var retryDelays = []time.Duration{time.Minute, 5 * time.Minute}

type Sender interface {
	Send(job models.EmailJob) error
}

func Run(ctx context.Context, db *gorm.DB, sender Sender) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		for ctx.Err() == nil {
			processed, err := ProcessNext(db, sender)
			if err != nil {
				log.Println("worker email:", err)
			}
			if err != nil || !processed {
				break
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func ProcessNext(db *gorm.DB, sender Sender) (bool, error) {
	processed := false
	err := db.Transaction(func(tx *gorm.DB) error {
		var job models.EmailJob
		err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("status = ? AND available_at <= ?", "pending", time.Now()).
			Order("available_at").Take(&job).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
		processed = true
		return tx.Model(&job).Updates(attempt(job, sender)).Error
	})
	return processed, err
}

func attempt(job models.EmailJob, sender Sender) map[string]any {
	err := sender.Send(job)
	if err == nil {
		return map[string]any{"status": "sent", "sent_at": time.Now(), "attempts": job.Attempts + 1, "last_error": nil}
	}
	log.Printf("email %d gagal: %v", job.ID, err)

	lastError := []rune(err.Error())
	attempts := job.Attempts + 1
	changes := map[string]any{"attempts": attempts, "last_error": string(lastError[:min(len(lastError), 500)])}
	if attempts >= maxAttempts {
		changes["status"] = "failed"
	} else {
		changes["available_at"] = time.Now().Add(retryDelays[attempts-1])
	}
	return changes
}
