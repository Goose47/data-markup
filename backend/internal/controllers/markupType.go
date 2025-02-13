package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"markup/internal/domain/models"
	"markup/internal/lib/responses"
	"net/http"
	"strconv"
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

	var markups []models.MarkupType

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		log.Warn("wrong page parameter", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong page parameter",
		})
		return
	}
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		log.Warn("wrong per_page parameter", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong per_page parameter",
		})
		return
	}
	offset := (page - 1) * perPage

	var total int64
	con.db.Model(&models.MarkupType{}).Count(&total)

	con.db.Limit(perPage).Offset(offset).Find(&markups)

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
	Fields []storeMarkupTypeField `binding:"required" json:"fields"`
}

type storeMarkupTypeField struct {
	Name             *string `binding:"required" json:"name"`
	Label            string  `binding:"required" json:"label"`
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
	if tx.Error != nil {
		log.Error("failed to begin transaction", slog.Any("error", tx.Error))
		responses.InternalServerError(c)
		return
	}
	var userID *uint // todo retrieve from authenticated user
	markupType := models.MarkupType{
		Name:   data.Name,
		UserID: userID,
	}

	if err := tx.Create(&markupType).Error; err != nil {
		tx.Rollback()
		log.Error("failed to create markup type", slog.Any("error", tx.Error))
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
		log.Error("failed to create markup type fields", slog.Any("error", tx.Error))
		responses.InternalServerError(c)
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("failed to commit transaction", slog.Any("error", tx.Error))
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
		ID           *uint `json:"id"`
		MarkupTypeID uint  `binding:"required" json:"markup_type_id"`
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
	if tx.Error != nil {
		log.Error("failed to begin transaction", slog.Any("error", tx.Error))
		responses.InternalServerError(c)
		return
	}

	markupType.Name = data.Name

	if err := tx.Save(&markupType).Error; err != nil {
		tx.Rollback()
		log.Error("failed to update markup type", slog.Any("error", tx.Error))
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
			MarkupTypeID:     field.MarkupTypeID,
		}

		if field.ID != nil {
			nextField.ID = *field.ID
			if err := tx.Save(&nextField).Error; err != nil {
				tx.Rollback()
				log.Error("failed to update markup type field", slog.Any("error", tx.Error))
				responses.InternalServerError(c)
				return
			}
			processedIds[i] = *field.ID
			continue
		}

		if err := tx.Create(&nextField).Error; err != nil {
			tx.Rollback()
			log.Error("failed to create markup type field", slog.Any("error", tx.Error))
			responses.InternalServerError(c)
			return
		}
		processedIds[i] = nextField.ID
	}

	result := tx.Where("markup_type_id = ? AND id NOT IN ?", markupType.ID, processedIds).Delete(&models.MarkupTypeField{})
	if result.Error != nil {
		tx.Rollback()
		log.Error("failed to delete markup type fields", slog.Any("error", tx.Error))
		responses.InternalServerError(c)
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("failed to commit transaction", slog.Any("error", tx.Error))
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
