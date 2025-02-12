package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"markup/internal/domain/models"
)

func New(host string, port int, user string, pass string, dbname string) (*gorm.DB, error) {
	const op = "storage.postgres.New"

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		host,
		user,
		pass,
		dbname,
		port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.AutoMigrate(
		&models.User{}, &models.Role{}, &models.Permission{},
		&models.Batch{}, &models.Markup{}, &models.MarkupType{},
		&models.MarkupTypeField{}, &models.AssessmentType{},
		&models.Assessment{}, &models.AssessmentField{},
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}

func Close(db *gorm.DB) {
	instance, _ := db.DB()
	_ = instance.Close()
}
