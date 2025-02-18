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
	"math/rand"
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
		Preload("Fields.MarkupTypeField").
		Preload("User").
		Preload("Markup.Batch.MarkupTypes.Fields")
	if userID > 0 {
		tx = tx.Where("user_id = ?", userID)
	}
	if markupID > 0 {
		tx = tx.Where("markup_id = ?", markupID)
	}
	tx = tx.Where("hash IS NOT NULL")
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

	// Delete admin assessment.
	if err := tx.
		Where("markup_id = ? AND is_prior IS TRUE", assessment.MarkupID).
		Delete(models.Assessment{}).Error; err != nil {
		tx.Rollback()
		log.Error("failed to delete other admins' assessments", slog.Any("error", err))
		responses.InternalServerError(c)
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

	log.Info("searching for pending assesment")
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
		log.Info("pending assesment found", slog.Any("assessment_id", pendingAssessment.ID))
		c.JSON(http.StatusOK, formatNextResponse(pendingAssessment))
		return
	}

	log.Info("fetching priorities")
	var priorities []int
	//err = con.db.
	//	Model("batches b").
	//	Select("priority").
	//	Joins("JOIN markups m ON m.batch_id = b.id").
	//	Order("priority desc").
	//	Distinct("priority").
	//	Pluck("priority", &priorities).Error
	err = con.db.
		Table("markups m").
		Select("b.priority priority,a.id, m.id, b.overlaps").
		Joins("JOIN batches b ON m.batch_id = b.id").
		Joins("LEFT JOIN assessments a ON a.markup_id = m.id").
		Where("m.status_id = ? and b.is_active IS TRUE", markupStatus.Pending).
		Group("m.id, b.overlaps, b.priority").
		//Having("COUNT(a.id) < b.overlaps").
		Having("NOT EXISTS (SELECT 1 FROM assessments a3 WHERE a3.markup_id = m.id AND a3.hash IS NULL)").
		Having("NOT EXISTS (SELECT 1 FROM assessments a2 WHERE a2.markup_id = m.id AND a2.user_id = ?)", user.ID).
		Distinct("priority").
		Pluck("priority", &priorities).Error

	if err != nil {
		log.Error("unable to fetch priorities", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}
	if len(priorities) == 0 {
		log.Error("failed to find markup", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	priority := weightedRandomChoice(priorities)
	log.Info("selected priority", slog.Any("priority", priority))

	var res struct {
		MarkupID         uint `gorm:"column:id"`
		AssessmentsCount int64
	}
	err = con.db.
		Table("markups").
		Select("markups.id, COUNT(assessments.id), batches.overlaps, markups.status_id").
		Joins("JOIN batches ON markups.batch_id = batches.id").
		Joins("LEFT JOIN assessments ON assessments.markup_id = markups.id").
		Where("status_id = ? and batches.priority = ? and batches.is_active IS TRUE", markupStatus.Pending, priority).
		Group("markups.id, markups.status_id, batches.overlaps").
		//Having("COUNT(assessments.id) < batches.overlaps").
		Having("NOT EXISTS (SELECT 1 FROM assessments a3 WHERE a3.markup_id = m.id AND a3.hash IS NULL)").
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

func weightedRandomChoice(priorities []int) int {
	// Create weight slices based on the values
	totalWeight := 0
	for _, p := range priorities {
		totalWeight += p
	}

	rand.Seed(time.Now().UnixNano())
	randomVal := rand.Intn(totalWeight) // Pick a random number in the total weight range

	// Select the number based on weighted distribution
	currentWeight := 0
	for _, p := range priorities {
		currentWeight += p
		if randomVal < currentWeight {
			return p
		}
	}

	return priorities[0] // Fallback (should never reach here)
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

	if !isAdmin && assessment.UpdatedAt.Add(30*time.Minute).Before(time.Now()) {
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

	result := tx.Where("assessment_id = ?", assessment.ID).Delete(&models.AssessmentField{})
	if err := result.Error; err != nil {
		tx.Rollback()
		log.Error("failed to delete assessment fields", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	assessment.Fields = make([]models.AssessmentField, len(data.Fields))
	for i, field := range data.Fields {
		assessment.Fields[i] = models.AssessmentField{
			Text:              field.Text,
			MarkupTypeFieldID: field.MarkupTypeFieldID,
		}
	}

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
