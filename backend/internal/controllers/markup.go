package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"markup/internal/domain/enums/markupStatus"
	"markup/internal/domain/models"
	"markup/internal/lib/responses"
	"markup/internal/lib/validation/query"
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

	var page int
	var perPage int

	if page, err = query.DefaultInt(c, log, "page", "1"); err != nil {
		return
	}
	if perPage, err = query.DefaultInt(c, log, "per_page", "10"); err != nil {
		return
	}
	offset := (page - 1) * perPage

	var total int64
	con.db.Model(&models.Markup{}).
		Where("batch_id = ?", batchID).
		Count(&total)

	con.db.Limit(perPage).
		Preload("Assessments").
		Where("batch_id = ?", batchID).
		Order("correct_assessment_hash IS NULL, id asc").
		Offset(offset).
		Find(&markups)

	c.JSON(http.StatusOK, responses.Pagination(markups, total, page, perPage))
}

type markupFindResponse struct {
	models.Markup
	CorrectAssessment *models.Assessment `json:"correct_assessment"`
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
			responses.NotFoundError(c)
			return
		}

		log.Error("failed to find markup", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	var correctAssessment *models.Assessment
	if markup.StatusID == markupStatus.Processed {
		err := con.db.
			Where("hash = ? and markup_id = ?", markup.CorrectAssessmentHash, markup.ID).
			First(&correctAssessment).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Error("correct assessment not found", slog.String("hash", *markup.CorrectAssessmentHash))
			} else {
				log.Error("failed to find markup", slog.Any("error", err))
			}
			responses.InternalServerError(c)
			return
		}
	}

	c.JSON(http.StatusOK, markupFindResponse{
		markup,
		correctAssessment,
	})
}
