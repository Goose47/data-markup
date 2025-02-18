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

type assessmentsResponse struct {
	models.Assessment
	MarkupType models.MarkupType `json:"markup_type"`
	User       models.User       `json:"user"`
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
	if userID, err = query.DefaultInt(c, log, "user_id", "0"); err != nil {
		return
	}
	if markupID, err = query.DefaultInt(c, log, "markup_id", "0"); err != nil {
		return
	}
	if page, err = query.DefaultInt(c, log, "page", "1"); err != nil {
		return
	}
	if perPage, err = query.DefaultInt(c, log, "per_page", "10"); err != nil {
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

	tx = con.db.Limit(perPage).
		Offset(offset).
		Preload("Fields").
		Preload("User").
		Preload("Markup.Batch.MarkupTypes.Fields")
	if userID > 0 {
		tx = tx.Where("user_id = ?", userID)
	}
	if markupID > 0 {
		tx = tx.Where("markup_id = ?", markupID)
	}
	tx.Find(&assessments)

	resp := make([]assessmentsResponse, len(assessments))
	for i, assessment := range assessments {
		// search for markup type that corresponds with each assessment.
		var respectiveMarkupType models.MarkupType
		for _, mt := range assessment.Markup.Batch.MarkupTypes {
			if mt.ChildID == nil {
				respectiveMarkupType = mt
				break
			}
		}

		resp[i] = assessmentsResponse{
			assessment,
			respectiveMarkupType,
			assessment.User,
		}
	}

	c.JSON(http.StatusOK, responses.Pagination(resp, total, page, perPage))
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

// Store creates models.Assessment for given models.Markup. Is used only by admins.
// todo: ensure there is only one models.Assessment for models.Markup for every models.Assessment.UserID
// todo: forbid to create models.Assessment when models.Markup is already processed if models.User is not admin
func (con *Assessment) Store(c *gin.Context) {
	const op = "AssessmentController.Store"
	log := con.log.With(slog.String("op", op))

	user, err := auth.User(c)
	if err != nil {
		responses.UnauthorizedError(c)
		return
	}
	var isAdmin = user.HasRole(roles.Admin)
	if !isAdmin {
		responses.ForbiddenError(c)
		return
	}

	var data storeAssessment

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assessment := models.Assessment{
		CreatedAt: time.Now(),
		UserID:    user.ID,
		IsPrior:   true,
		MarkupID:  data.MarkupID,
	}

	assessment.Fields = make([]models.AssessmentField, len(data.Fields))
	for i, field := range data.Fields {
		assessment.Fields[i] = models.AssessmentField{
			MarkupTypeFieldID: field.MarkupTypeFieldID,
			Text:              field.Text,
		}
	}

	hash := assessment.CalculateHash()
	assessment.Hash = &hash

	tx := con.db.Begin()
	if err := tx.Error; err != nil {
		log.Error("failed to begin transaction", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	// Save assessment.
	if err := tx.Create(&assessment).Error; err != nil {
		tx.Rollback()
		log.Error("failed to create assessment", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	if err := updateCorrectAssessment(log, tx, assessment, true); err != nil {
		tx.Rollback()
		responses.InternalServerError(c)
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("failed to commit transaction", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": assessment.ID,
	})
}

// Next selects available markup and creates empty models.Assessment for it.
func (con *Assessment) Next(c *gin.Context) {
	const op = "AssessmentController.Next"
	log := con.log.With(slog.String("op", op))

	user, err := auth.User(c)
	if err != nil {
		responses.UnauthorizedError(c)
		return
	}
	var isAssessor = user.HasRole(roles.Assessor)
	if !isAssessor {
		responses.ForbiddenError(c)
		return
	}

	var priorities []int
	err = con.db.
		Model(models.Batch{}).
		Order("priority desc").
		Distinct("priority").
		Pluck("priority", &priorities).Error

	if err != nil {
		log.Error("unable to fetch priorities", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}
	// todo do smth w priorities

	var pendingAssessment models.Assessment
	err = con.db.Model(models.Assessment{}).
		Preload("Markup.Batch.MarkupTypes.Fields.AssessmentType").
		Where("hash IS NULL and user_id = ?", user.ID).
		First(&pendingAssessment).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("failed to search for pending assessment", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}
	if err == nil {
		c.JSON(http.StatusOK, formatNextResponse(pendingAssessment))
		return
	}

	var res struct {
		MarkupID         uint `gorm:"column:id"`
		AssessmentsCount int64
	}
	err = con.db.
		Table("markups").
		Select("markups.id, COUNT(assessments.id), batches.overlaps, markups.status_id").
		Joins("JOIN batches ON markups.batch_id = batches.id").
		Joins("LEFT JOIN assessments ON assessments.markup_id = markups.id").
		Where("status_id = ? and batches.is_active IS TRUE", markupStatus.Pending).
		Group("markups.id, markups.status_id, batches.overlaps").
		Having("COUNT(assessments.id) < batches.overlaps").
		Having("NOT EXISTS (SELECT 1 FROM assessments a WHERE a.markup_id = markups.id AND a.user_id = ?)", user.ID).
		First(&res).Error

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

	assessment := models.Assessment{
		UserID:    user.ID,
		MarkupID:  res.MarkupID,
		CreatedAt: time.Now(),
		IsPrior:   false,
		Hash:      nil,
	}

	// Save assessment.
	if err := con.db.Create(&assessment).Error; err != nil {
		log.Error("failed to create assessment", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	con.db.Preload("Markup.Batch.MarkupTypes.Fields.AssessmentType").First(&assessment)

	c.JSON(http.StatusCreated, formatNextResponse(assessment))
}

func formatNextResponse(assessment models.Assessment) *gin.H {
	var markupType models.MarkupType
	for _, mt := range assessment.Markup.Batch.MarkupTypes {
		if mt.ChildID == nil {
			markupType = mt
		}
	}

	return &gin.H{
		"assessment_id": assessment.ID,
		"markup_type":   markupType,
		"data":          assessment.Markup.Data,
	}
}

// updateCorrectAssessment checks whether there is enough identical assessments or an models.Assessment is made by admin
// to mark models.Markup as markupStatus.Processed and set models.Markup.CorrectAssessmentHash.
func updateCorrectAssessment(
	log *slog.Logger,
	tx *gorm.DB,
	assessment models.Assessment,
	isAdmin bool,
) error {
	const op = "Assessment.UpdateCorrectAssessment"
	log = log.With(slog.String("op", op))

	isCorrectAssessment := false
	if isAdmin {
		isCorrectAssessment = true
	} else {
		// Check if enough assessments are made.
		tx.Preload("Markup.Batch").First(&assessment)
		var count int64
		tx.Model(&Assessment{}).
			Where("hash = ? AND markup_id = ?", assessment.Hash, assessment.MarkupID).
			Count(&count)

		if count >= int64(assessment.Markup.Batch.Overlaps) {
			isCorrectAssessment = true
		}
	}

	if isCorrectAssessment {
		// Mark Markup as processed.
		tx.Model(&Markup{}).
			Where("id = ?", assessment.MarkupID).
			Updates(map[string]interface{}{
				"status_id":               markupStatus.Processed,
				"correct_assessment_hash": assessment.Hash,
			})
		if err := tx.Error; err != nil {
			log.Error("failed to update markup", slog.Any("error", err))
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
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

	user, err := auth.User(c)
	if err != nil {
		responses.UnauthorizedError(c)
		return
	}
	var isAdmin = user.HasRole(roles.Admin)

	if user.ID != assessment.UserID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot update assessment of another user",
		})
		return
	}

	if assessment.UpdatedAt.Add(30 * time.Minute).Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot update assessment if > 30 minutes have passed",
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
		processedIds[i] = nextField.ID
	}

	result := tx.Where("assessment_id = ? AND id NOT IN ?", assessment.ID, processedIds).Delete(&models.AssessmentField{})
	if err := result.Error; err != nil {
		tx.Rollback()
		log.Error("failed to delete assessment fields", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	// load current fields
	tx.Preload("Fields").First(&assessment)
	hash := assessment.CalculateHash()
	assessment.Hash = &hash
	if err := tx.Save(&assessment).Error; err != nil {
		tx.Rollback()
		log.Error("failed to update assessment", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	if err := updateCorrectAssessment(log, tx, assessment, isAdmin); err != nil {
		tx.Rollback()
		responses.InternalServerError(c)
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
