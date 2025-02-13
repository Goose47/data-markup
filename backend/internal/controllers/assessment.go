package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"markup/internal/domain/models"
	"markup/internal/lib/responses"
	"markup/internal/lib/validation/query"
	"net/http"
	"time"
)

type Assessment struct {
	log *slog.Logger
	db  *gorm.DB
}

func NewAssessment(
	log *slog.Logger,
	db *gorm.DB,
) *Assessment {
	return &Assessment{
		log: log,
		db:  db,
	}
}

func (con *Assessment) Index(c *gin.Context) {
	const op = "AssessmentController.Index"
	log := con.log.With(slog.String("op", op))

	var userID int
	var markupID int
	var page int
	var perPage int
	var err error

	// todo: refactor query param fetching. Use query functions
	if userID, err = query.Int(c, log, "user_id", "0"); err != nil {
		return
	}
	if markupID, err = query.Int(c, log, "markup_id", "0"); err != nil {
		return
	}
	if page, err = query.Int(c, log, "page", "1"); err != nil {
		return
	}
	if perPage, err = query.Int(c, log, "per_page", "10"); err != nil {
		return
	}
	offset := (page - 1) * perPage

	var assessments []models.Assessment
	var total int64

	tx := con.db.Model(&models.Assessment{})
	if userID > 0 {
		tx = tx.Where("user_id = ?", userID)
	}
	if markupID > 0 {
		tx = tx.Where("markup_id = ?", markupID)
	}
	tx.Count(&total)

	tx = con.db.Limit(perPage).Offset(offset)
	if userID > 0 {
		tx = tx.Where("user_id = ?", userID)
	}
	if markupID > 0 {
		tx = tx.Where("markup_id = ?", markupID)
	}
	tx.Find(&assessments)

	c.JSON(http.StatusOK, responses.Pagination(assessments, total, page, perPage))
}

func (con *Assessment) Find(c *gin.Context) {
	const op = "AssessmentController.Find"
	id := c.Param("id")

	log := con.log.With(slog.String("op", op), slog.String("id", id))

	var assessment models.Assessment
	err := con.db.
		Preload("Fields").
		Where("id = ?", id).
		First(&assessment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("assessment not found")
			responses.NotFoundError(c)
			return
		}

		log.Error("failed to find assessment", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, assessment)
}

type storeAssessment struct {
	MarkupID uint                   `binding:"required" json:"markup_id"`
	Fields   []storeAssessmentField `binding:"required,dive" json:"fields"`
}

type storeAssessmentField struct {
	Text              *string `json:"text"`
	MarkupTypeFieldID uint    `binding:"required" json:"markup_type_field_id"`
}

func (con *Assessment) Store(c *gin.Context) {
	const op = "AssessmentController.Store"
	log := con.log.With(slog.String("op", op))

	var data storeAssessment

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID uint = 3 // todo retrieve from authenticated user
	var isAdmin = false // todo retrieve from authenticated user

	assessment := models.Assessment{
		CreatedAt: time.Now(),
		UserID:    userID,
		IsPrior:   isAdmin,
		MarkupID:  data.MarkupID,
	}

	assessment.Fields = make([]models.AssessmentField, len(data.Fields))
	for i, field := range data.Fields {
		assessment.Fields[i] = models.AssessmentField{
			MarkupTypeFieldID: field.MarkupTypeFieldID,
			Text:              field.Text,
		}
	}

	if err := con.db.Create(&assessment).Error; err != nil {
		log.Error("failed to create assessment", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": assessment.ID,
	})
}

type updateAssessment struct {
	Fields []struct {
		ID *uint `json:"id"`
		storeAssessmentField
	} `binding:"required,dive" json:"fields"`
}

// todo: can steal another assessment fields if their ids are passed in storeAssessmentField.
// todo: check if specified AssessmentField ids belong to respective Assessment.
// todo: same in MarkupType.Update.
func (con *Assessment) Update(c *gin.Context) {
	const op = "AssessmentController.Update"
	id := c.Param("id")

	log := con.log.With(slog.String("op", op), slog.String("id", id))

	var data updateAssessment

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var assessment models.Assessment
	err := con.db.
		Preload("Fields").
		Where("id = ?", id).
		First(&assessment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("assessment not found")
			responses.NotFoundError(c)
			return
		}

		log.Error("failed to find assessment", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	var userID uint = 3 // todo: retrieve from authenticated user
	if userID != assessment.UserID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot update assessment of another user",
		})
		return
	}

	tx := con.db.Begin()
	if err := tx.Error; err != nil {
		log.Error("failed to begin transaction", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	processedIds := make([]uint, len(data.Fields))

	assessment.Fields = make([]models.AssessmentField, len(data.Fields))
	for i, field := range data.Fields {
		nextField := models.AssessmentField{
			Text:              field.Text,
			MarkupTypeFieldID: field.MarkupTypeFieldID,
			AssessmentID:      assessment.ID,
		}
		if field.ID != nil {
			nextField.ID = *field.ID
			if err := tx.Save(&nextField).Error; err != nil {
				tx.Rollback()
				log.Error("failed to assessment type field", slog.Any("error", err))
				responses.InternalServerError(c)
				return
			}
			processedIds[i] = *field.ID
			continue
		}

		if err := tx.Create(&nextField).Error; err != nil {
			tx.Rollback()
			log.Error("failed to create assessment field", slog.Any("error", err))
			responses.InternalServerError(c)
			return
		}
		assessment.Fields[i] = nextField
	}

	if err := tx.Save(&assessment).Error; err != nil {
		tx.Rollback()
		log.Error("failed to update assessment fields", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	result := tx.Where("assessment_id = ? AND id NOT IN ?", assessment.ID, processedIds).Delete(&models.AssessmentField{})
	if err := result.Error; err != nil {
		tx.Rollback()
		log.Error("failed to delete assessment fields", slog.Any("error", err))
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

func (con *Assessment) Destroy(c *gin.Context) {
	const op = "AssessmentController.Destroy"
	id := c.Param("id")

	log := con.log.With(slog.String("op", op), slog.String("id", id))

	log.Info("AssessmentController.Destroy")

	c.JSON(http.StatusNotImplemented, nil)
}
