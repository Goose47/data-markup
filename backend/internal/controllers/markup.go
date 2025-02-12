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

type Markup struct {
	log *slog.Logger
	db  *gorm.DB
}

func NewMarkup(
	log *slog.Logger,
	db *gorm.DB,
) *Markup {
	return &Markup{
		log: log,
		db:  db,
	}
}

func (con *Markup) Index(c *gin.Context) {
	const op = "MarkupController.Index"
	log := con.log.With(slog.String("op", op))

	var markups []models.Markup

	batchID, err := strconv.Atoi(c.DefaultQuery("batch_id", "0"))
	if err != nil || batchID == 0 {
		log.Warn("wrong batch_id parameter", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong batch_id parameter",
		})
		return
	}

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
	con.db.Model(&models.Markup{}).
		Where("batch_id = ?", batchID).
		Count(&total)

	con.db.Limit(perPage).
		Where("batch_id = ?", batchID).
		Offset(offset).
		Find(&markups)

	c.JSON(http.StatusOK, responses.Pagination(markups, total, page, perPage))
}

func (con *Markup) Find(c *gin.Context) {
	const op = "MarkupController.Find"
	id := c.Param("id")

	log := con.log.With(slog.String("op", op), slog.String("id", id))

	var markup models.Markup
	err := con.db.
		Preload("Assessments").
		Where("id = ?", id).
		First(&markup).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("markup not found")
			c.JSON(http.StatusNotFound, gin.H{
				"error": "markup not found",
			})
			return
		}

		log.Error("failed to find markup", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, markup)
}
