package background

import (
	"gorm.io/gorm"
	"log/slog"
	"markup/internal/domain/models"
	"time"
)

type TaskManager struct {
	log *slog.Logger
	db  *gorm.DB
}

func NewTaskManager(log *slog.Logger, db *gorm.DB) *TaskManager {
	return &TaskManager{
		log: log,
		db:  db,
	}
}

func (tm *TaskManager) Run() {
	go tm.deleteOutdatedAssessments()
}

func (tm *TaskManager) deleteOutdatedAssessments() {
	for {
		err := tm.db.
			Where("created_at < ? AND hash IS NULL", time.Now().Add(-5*time.Minute)).
			Delete(models.Assessment{}).Error
		if err != nil {
			tm.log.Error("failed to delete outdated assessments", slog.Any("error", err))
		}
		time.Sleep(5 * time.Second)
	}
}
