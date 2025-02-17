package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"markup/internal/domain/models"
	"markup/internal/lib/auth"
	"markup/internal/lib/responses"
	"markup/internal/lib/validation/query"
	"net/http"
	"time"
)

type MarkupType struct {
	log *slog.Logger
	db  *gorm.DB
}

func NewMarkupType(
	log *slog.Logger,
	db *gorm.DB,
) *MarkupType {
	return &MarkupType{
		log: log,
		db:  db,
	}
}

func (con *MarkupType) Index(c *gin.Context) {
	const op = "MarkupTypeController.Index"
	log := con.log.With(slog.String("op", op))

	var batchID *int
	var page int
	var perPage int
	var err error

	batchID = query.Int(c, "batch_id")
	if page, err = query.DefaultInt(c, log, "page", "1"); err != nil {
		return
	}
	if perPage, err = query.DefaultInt(c, log, "per_page", "10"); err != nil {
		return
	}
	offset := (page - 1) * perPage

	var total int64
	tx := con.db.Model(&models.MarkupType{})
	if batchID != nil {
		if *batchID > 0 {
			tx = tx.Where("batch_id = ?", batchID)
		} else {
			tx = tx.Where("batch_id IS NULL")
		}
	}
	tx.Count(&total)

	var markups []struct {
		models.MarkupType
		MarkupCount            int `json:"markup_count"`
		AssessmentCount        int `json:"assessment_count"`
		CorrectAssessmentCount int `json:"correct_assessment_count"`
	}
	tx = con.db.
		Table("markup_types mt").
		Select("mt.id,mt.batch_id,mt.name,mt.child_id,mt.user_id,mt.created_at,COUNT(DISTINCT a.markup_id) AS markup_count, COUNT(DISTINCT a.id) AS assessment_count,COUNT(DISTINCT a2.id) AS correct_assessment_count").
		Joins("LEFT JOIN markup_type_fields mtf ON mt.id = mtf.markup_type_id").
		Joins("LEFT JOIN assessment_fields af ON af.markup_type_field_id = mtf.id").
		Joins("LEFT JOIN assessments a ON af.assessment_id = a.id").
		Joins("LEFT JOIN markups m ON a.markup_id = m.id").
		Joins("LEFT JOIN assessments a2 ON af.assessment_id = a2.id and a2.hash = m.correct_assessment_hash").
		Group("mt.id").
		Limit(perPage).
		Offset(offset)
	if batchID != nil {
		if *batchID > 0 {
			tx = tx.Where("mt.batch_id = ?", batchID)
		} else {
			tx = tx.Where("mt.batch_id IS NULL")
		}
	}
	tx.Find(&markups)

	if err := tx.Error; err != nil {
		log.Error("failed to query markup_types", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, responses.Pagination(markups, total, page, perPage))
}

func (con *MarkupType) Find(c *gin.Context) {
	const op = "MarkupTypeController.Find"
	id := c.Param("id")

	log := con.log.With(slog.String("op", op), slog.String("id", id))

	var markupType models.MarkupType
	err := con.db.
		Preload("Fields.AssessmentType").
		Where("id = ?", id).
		First(&markupType).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("markupType not found")
			c.JSON(http.StatusNotFound, gin.H{
				"error": "markupType not found",
			})
			return
		}

		log.Error("failed to find markupType", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, markupType)
}

type storeMarkupType struct {
	Name   string                 `binding:"required" json:"name"`
	Fields []storeMarkupTypeField `binding:"required,dive" json:"fields"`
}

type storeMarkupTypeField struct {
	Name             *string `binding:"required" json:"name"`
	Label            *string `binding:"required" json:"label"`
	GroupID          uint    `binding:"required" json:"group_id"`
	AssessmentTypeID uint    `binding:"required" json:"assessment_type_id"`
}

func (con *MarkupType) Store(c *gin.Context) {
	const op = "MarkupTypeController.Store"
	log := con.log.With(slog.String("op", op))

	var data storeMarkupType

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := con.db.Begin()
	if err := tx.Error; err != nil {
		log.Error("failed to begin transaction", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	user, err := auth.User(c)
	if err != nil {
		responses.UnauthorizedError(c)
		return
	}

	markupType := models.MarkupType{
		Name:      data.Name,
		UserID:    &user.ID,
		CreatedAt: time.Now(),
	}

	// todo: make only one query to save all models. See Assesment.Store
	if err := tx.Create(&markupType).Error; err != nil {
		tx.Rollback()
		log.Error("failed to create markup type", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	fields := make([]models.MarkupTypeField, len(data.Fields))
	for i, field := range data.Fields {
		fields[i] = models.MarkupTypeField{
			Name:             field.Name,
			Label:            field.Label,
			GroupID:          field.GroupID,
			AssessmentTypeID: field.AssessmentTypeID,
			MarkupTypeID:     markupType.ID,
		}
	}

	if err := tx.Create(&fields).Error; err != nil {
		tx.Rollback()
		log.Error("failed to create markup type fields", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("failed to commit transaction", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": markupType.ID,
	})
}

type updateMarkupType struct {
	Name   string `binding:"required" json:"name"`
	Fields []struct {
		ID *uint `json:"id"`
		storeMarkupTypeField
	} `binding:"required" json:"fields"`
}

func (con *MarkupType) Update(c *gin.Context) {
	const op = "MarkupTypeController.Update"
	id := c.Param("id")

	log := con.log.With(slog.String("op", op), slog.String("id", id))

	var data updateMarkupType

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var markupType models.MarkupType
	err := con.db.
		Preload("Fields.AssessmentType").
		Where("id = ?", id).
		First(&markupType).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("markupType not found")
			c.JSON(http.StatusNotFound, gin.H{
				"error": "markupType not found",
			})
			return
		}

		log.Error("failed to find markupType", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	if markupType.BatchID != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot update markupType that is tied to a batch",
		})
		return
	}

	tx := con.db.Begin()
	if err := tx.Error; err != nil {
		log.Error("failed to begin transaction", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	markupType.Name = data.Name

	if err := tx.Save(&markupType).Error; err != nil {
		tx.Rollback()
		log.Error("failed to update markup type", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	processedIds := make([]uint, len(data.Fields))

	for i, field := range data.Fields {
		nextField := models.MarkupTypeField{
			Name:             field.Name,
			Label:            field.Label,
			GroupID:          field.GroupID,
			AssessmentTypeID: field.AssessmentTypeID,
			MarkupTypeID:     markupType.ID,
		}

		if field.ID != nil {
			nextField.ID = *field.ID
			if err := tx.Save(&nextField).Error; err != nil {
				tx.Rollback()
				log.Error("failed to update markup type field", slog.Any("error", err))
				responses.InternalServerError(c)
				return
			}
			processedIds[i] = *field.ID
			continue
		}

		if err := tx.Create(&nextField).Error; err != nil {
			tx.Rollback()
			log.Error("failed to create markup type field", slog.Any("error", err))
			responses.InternalServerError(c)
			return
		}
		processedIds[i] = nextField.ID
	}

	result := tx.Where("markup_type_id = ? AND id NOT IN ?", markupType.ID, processedIds).Delete(&models.MarkupTypeField{})
	if err := result.Error; err != nil {
		tx.Rollback()
		log.Error("failed to delete markup type fields", slog.Any("error", err))
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

func (con *MarkupType) Destroy(c *gin.Context) {
	const op = "MarkupTypeController.Destroy"
	id := c.Param("id")

	log := con.log.With(slog.String("op", op), slog.String("id", id))

	var markupType models.MarkupType
	err := con.db.
		Preload("Fields.AssessmentType").
		Where("id = ?", id).
		First(&markupType).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("markupType not found")
			c.JSON(http.StatusNotFound, gin.H{
				"error": "markupType not found",
			})
			return
		}

		log.Error("failed to find markupType", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	if markupType.BatchID != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot delete markupType that is tied to a batch",
		})
		return
	}

	if err := con.db.Delete(&models.MarkupType{}, id).Error; err != nil {
		log.Error("failed to delete markupType", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, "OK")
}
