package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"markup/internal/domain/enums/markupStatus"
	"markup/internal/domain/enums/roles"
	"markup/internal/domain/models"
	"markup/internal/lib/auth"
	"markup/internal/lib/responses"
	"markup/internal/lib/validation/query"
	"net/http"
	"time"
)

type Honeypot struct {
	log *slog.Logger
	db  *gorm.DB
}

func NewHoneypot(
	log *slog.Logger,
	db *gorm.DB,
) *Honeypot {
	return &Honeypot{
		log: log,
		db:  db,
	}
}

func (con *Honeypot) Index(c *gin.Context) {
	const op = "HoneypotController.Index"
	log := con.log.With(slog.String("op", op))

	var batches []models.Batch

	var page int
	var perPage int
	var err error

	if page, err = query.DefaultInt(c, log, "page", "1"); err != nil {
		return
	}
	if perPage, err = query.DefaultInt(c, log, "per_page", "10"); err != nil {
		return
	}
	offset := (page - 1) * perPage

	user, err := auth.User(c)
	if err != nil {
		responses.UnauthorizedError(c)
		return
	}
	isAdmin := user.HasRole(roles.Admin)

	var total int64
	tx := con.db.Model(&models.Batch{})
	if !isAdmin {
		tx = tx.Where("user_id = ?", user.ID)
	}
	tx.Where("is_honeypot IS TRUE")
	tx.Count(&total)

	tx = con.db.Limit(perPage).
		Offset(offset).
		Order("created_at DESC").
		Preload("Markups.Assessment.Fields").
		Preload("MarkupTypes.Fields")
	if !isAdmin {
		tx = tx.Where("user_id = ?", user.ID)
	}
	tx.Where("is_honeypot IS TRUE")
	tx.Find(&batches)

	c.JSON(http.StatusOK, responses.Pagination(batches, total, page, perPage))
}

func (con *Honeypot) Store(c *gin.Context) {
	const op = "HoneypotController.Store"
	markupID := c.Param("id")
	log := con.log.With(slog.String("op", op), slog.String("markup_id", markupID))

	log.Info("saving honeypot")

	var markup models.Markup
	err := con.db.
		Preload("Batch.MarkupTypes.Fields").
		Preload("Assessments.Fields").
		Where("id = ?", markupID).
		First(&markup).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("markup not found")
			responses.NotFoundError(c)
			return
		}

		log.Error("failed to find markup", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	var assessment models.Assessment
	for _, a := range markup.Assessments {
		if a.IsPrior {
			assessment = a
		}
	}
	if assessment.ID == 0 {
		log.Warn("failed to find admin assessment", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нет эталонной оценки. Невозможно создать ханипот."})
		return
	}

	var markupType models.MarkupType
	for _, mt := range markup.Batch.MarkupTypes {
		if mt.ChildID == nil {
			markupType = mt
		}
	}
	if markupType.ID == 0 {
		log.Warn("failed to find markup type", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нет типа разметки. Невозможно создать ханипот."})
		return
	}

	tx := con.db.Begin()
	if err := tx.Error; err != nil {
		log.Error("failed to begin transaction", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	newBatch := models.Batch{
		Name:       fmt.Sprintf("HONEYPOT: %s", markup.Batch.Name),
		Overlaps:   1000,
		Priority:   50,
		TypeID:     markup.Batch.TypeID,
		IsHoneypot: true,
	}

	if err := tx.Create(&newBatch).Error; err != nil {
		tx.Rollback()
		log.Error("failed to create batch", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	newMarkup := models.Markup{
		BatchID:               newBatch.ID,
		StatusID:              markupStatus.Pending,
		Data:                  markup.Data,
		CorrectAssessmentHash: nil,
	}

	if err := tx.Create(&newMarkup).Error; err != nil {
		tx.Rollback()
		log.Error("failed to create markup", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	newMarkupType := models.MarkupType{
		BatchID:   &newBatch.ID,
		Name:      markupType.Name,
		ChildID:   nil,
		CreatedAt: time.Now(),
	}

	if err := tx.Create(&newMarkupType).Error; err != nil {
		tx.Rollback()
		log.Error("failed to create markup type", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	// old markupTypeFieldID -> new markupTypeFieldID
	MTFMap := make(map[uint]uint)

	newMarkupTypeFields := make([]models.MarkupTypeField, len(markupType.Fields))
	for i, field := range markupType.Fields {
		newMarkupTypeField := models.MarkupTypeField{
			MarkupTypeID:     newMarkupType.ID,
			AssessmentTypeID: field.AssessmentTypeID,
			Name:             field.Name,
			Label:            field.Label,
			GroupID:          field.GroupID,
		}
		newMarkupTypeFields[i] = newMarkupTypeField

		if err := tx.Create(&newMarkupTypeField).Error; err != nil {
			tx.Rollback()
			log.Error("failed to create markup type field", slog.Any("error", err))
			responses.InternalServerError(c)
			return
		}

		MTFMap[field.ID] = newMarkupTypeField.ID
	}

	assessmentFields := make([]models.AssessmentField, len(assessment.Fields))
	for i, field := range assessment.Fields {
		assessmentFields[i] = models.AssessmentField{
			MarkupTypeFieldID: MTFMap[field.MarkupTypeFieldID],
			Text:              field.Text,
		}
	}
	newAssessment := models.Assessment{
		MarkupID:  newMarkup.ID,
		UserID:    assessment.UserID,
		CreatedAt: time.Now(),
		IsPrior:   true,
		Fields:    assessmentFields,
	}
	hash := newAssessment.CalculateHash()
	newAssessment.Hash = &hash

	if err := tx.Create(&newAssessment).Error; err != nil {
		tx.Rollback()
		log.Error("failed to create assessment", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	newMarkup.CorrectAssessmentHash = &hash
	if err := tx.Save(&newMarkup).Error; err != nil {
		log.Error("failed to update markup hash", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("failed to commit transaction", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, "OK")
}
