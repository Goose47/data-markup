package controllers

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log/slog"
	"markup/internal/domain/enums/markupStatus"
	"markup/internal/domain/enums/roles"
	"markup/internal/domain/models"
	"markup/internal/lib/auth"
	"markup/internal/lib/responses"
	"net/http"
	"strconv"
	"time"
)

type Batch struct {
	log *slog.Logger
	db  *gorm.DB
}

func NewBatch(
	log *slog.Logger,
	db *gorm.DB,
) *Batch {
	return &Batch{
		log: log,
		db:  db,
	}
}

func (con *Batch) Index(c *gin.Context) {
	const op = "BatchController.Index"
	log := con.log.With(slog.String("op", op))

	var batches []models.Batch

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
	tx.Count(&total)

	tx = con.db.Limit(perPage).
		Offset(offset).
		Order("created_at DESC")
	if !isAdmin {
		tx = tx.Where("user_id = ?", user.ID)
	}
	tx.Find(&batches)

	c.JSON(http.StatusOK, responses.Pagination(batches, total, page, perPage))
}

func (con *Batch) Find(c *gin.Context) {
	const op = "BatchController.Find"
	id := c.Param("id")

	log := con.log.With(slog.String("op", op), slog.String("id", id))

	var batch models.Batch
	err := con.db.
		Table("batches b").
		Select("b.id,b.name,b.overlaps,b.priority,b.created_at,b.is_active,b.type_id").
		Where("id = ?", id).
		First(&batch).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("batch not found")
			responses.NotFoundError(c)
			return
		}

		log.Error("failed to find batch", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	var markupCount int64
	con.db.
		Model(models.Markup{}).
		Where("batch_id = ?", batch.ID).
		Count(&markupCount)

	var processedMarkupCount int64
	con.db.
		Model(models.Markup{}).
		Where("batch_id = ? AND status_id = ?", batch.ID, markupStatus.Processed).
		Count(&processedMarkupCount)

	var assessmentCount int64
	con.db.
		Table("assessments a").
		Select("DISTINCT a.id").
		Joins("JOIN markups m ON a.markup_id = m.id").
		Where("m.batch_id = ?", batch.ID).
		Count(&assessmentCount)

	var correctAssessmentCount int64
	con.db.
		Table("assessments a").
		Select("DISTINCT a.id").
		Joins("JOIN markups m ON a.markup_id = m.id AND a.hash = m.correct_assessment_hash").
		Where("m.batch_id = ?", batch.ID).
		Count(&correctAssessmentCount)

	var res struct {
		models.Batch
		MarkupCount            int64 `json:"markup_count"`
		ProcessedMarkupCount   int64 `json:"processed_markup_count"`
		AssessmentCount        int64 `json:"assessment_count"`
		CorrectAssessmentCount int64 `json:"correct_assessment_count"`
	}

	res.Batch = batch
	res.MarkupCount = markupCount
	res.ProcessedMarkupCount = processedMarkupCount
	res.AssessmentCount = assessmentCount
	res.CorrectAssessmentCount = correctAssessmentCount

	c.JSON(http.StatusOK, res)
}

type storeBatchType struct {
	Name     string `binding:"required" form:"name"`
	Overlaps int    `binding:"required" form:"overlaps"`
	Priority int    `binding:"required,gte=1,lte=10" form:"priority"`
	TypeID   uint   `binding:"required" form:"type_id"`
}

func (con *Batch) Store(c *gin.Context) {
	const op = "BatchController.Store"
	log := con.log.With(slog.String("op", op))

	var data storeBatchType

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle file upload separately
	file, err := c.FormFile("markups")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	tx := con.db.Begin()
	if err := tx.Error; err != nil {
		log.Error("failed to begin transaction", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	batch := models.Batch{
		Name:      data.Name,
		Overlaps:  data.Overlaps,
		Priority:  data.Priority,
		TypeID:    data.TypeID,
		CreatedAt: time.Now(),
		IsActive:  false,
	}

	user, err := auth.User(c)
	if err != nil {
		responses.UnauthorizedError(c)
		return
	}

	isAdmin := user.HasRole(roles.Admin)

	if !isAdmin {
		batch.Users = append(batch.Users, user)
	}

	//if userID != nil {
	//	var user models.User
	//	err := con.db.
	//		Where("id = ?", userID).
	//		First(&user).Error
	//	if err != nil {
	//		if errors.Is(err, gorm.ErrRecordNotFound) {
	//			log.Warn("user not found")
	//			c.JSON(http.StatusNotFound, gin.H{
	//				"error": "user not found",
	//			})
	//			return
	//		}
	//
	//		log.Error("failed to find user", slog.Any("error", err))
	//		responses.InternalServerError(c)
	//		return
	//	}
	//	batch.Users = append(batch.Users, user)
	//}

	if err := tx.Create(&batch).Error; err != nil {
		tx.Rollback()
		log.Error("failed to create batch", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	//parse csv and crate markups
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	reader := csv.NewReader(src)
	headers, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read CSV headers"})
		return
	}

	// Batch size
	batchSize := 100
	markups := make([]models.Markup, 0, batchSize)
	counter := 1

	// Read and process CSV records line by line
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // End of file reached
		}
		if err != nil {
			tx.Rollback()
			log.Error("failed to read next markup", slog.Int("counter", counter), slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error reading CSV at line %d", counter)})
			return
		}

		// Map CSV data to the model fields
		markupData := map[string]string{}
		for i, header := range headers {
			markupData[header] = record[i]
		}
		markupDataMarshalled, err := json.Marshal(markupData)
		if err != nil {
			tx.Rollback()
			log.Error("failed to marshal next markup", slog.Int("counter", counter), slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error reading CSV at line %d", counter)})
			return
		}

		markup := models.Markup{
			BatchID:               batch.ID,
			StatusID:              markupStatus.Pending,
			Data:                  string(markupDataMarshalled),
			CorrectAssessmentHash: nil,
		}

		markups = append(markups, markup)

		if len(markups) >= batchSize {
			if err := tx.Create(&markups).Error; err != nil {
				tx.Rollback()
				log.Error("failed to save next markup batch", slog.Any("error", err))
				responses.InternalServerError(c)
				return
			}

			markups = make([]models.Markup, 0, batchSize)
		}
	}

	// Insert the remaining records if any
	if len(markups) > 0 {
		if err := tx.Create(&markups).Error; err != nil {
			tx.Rollback()
			log.Error("failed to save next markup batch", slog.Any("error", err))
			responses.InternalServerError(c)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("failed to commit transaction", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": batch.ID,
	})
}

type updateBatchType struct {
	Name     string `binding:"required" json:"name"`
	Overlaps int    `binding:"required" json:"overlaps"`
	Priority int    `binding:"required,gte=1,lte=10" json:"priority"`
	TypeID   uint   `binding:"required" json:"type_id"`
	IsActive *bool  `binding:"required" json:"is_active"`
}

func (con *Batch) Update(c *gin.Context) {
	const op = "BatchController.Update"
	id := c.Param("id")

	log := con.log.With(slog.String("op", op), slog.String("id", id))

	var data updateBatchType

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var batch models.Batch
	err := con.db.
		Where("id = ?", id).
		First(&batch).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("batch not found")
			c.JSON(http.StatusNotFound, gin.H{
				"error": "batch not found",
			})
			return
		}

		log.Error("failed to find batch", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	batch.Name = data.Name
	batch.Overlaps = data.Overlaps
	batch.Priority = data.Priority
	batch.TypeID = data.TypeID
	batch.IsActive = *data.IsActive

	if err := con.db.Save(&batch).Error; err != nil {
		log.Error("failed to update batch type", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, "OK")
}

func (con *Batch) ToggleIsActive(c *gin.Context) {
	const op = "BatchController.ToggleIsActive"
	id := c.Param("id")

	log := con.log.With(slog.String("op", op), slog.String("id", id))

	var batch models.Batch
	err := con.db.
		Where("id = ?", id).
		First(&batch).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("batch not found")
			responses.NotFoundError(c)
			return
		}

		log.Error("failed to find batch", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	batch.IsActive = !batch.IsActive

	if err := con.db.Save(&batch).Error; err != nil {
		con.db.Rollback()
		log.Error("failed to toggle is_active field", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_active": batch.IsActive,
	})
}

func (con *Batch) Destroy(c *gin.Context) {
	const op = "BatchController.Destroy"
	id := c.Param("id")

	log := con.log.With(slog.String("op", op), slog.String("id", id))

	log.Info("Batchcon.destroy")

	c.JSON(http.StatusNotImplemented, nil)
}

type tieMarkupType struct {
	BatchID      uint  `binding:"required" json:"batch_id"`
	MarkupTypeID *uint `json:"markup_type_id"`
	storeMarkupType
}

func (con *Batch) TieMarkupType(c *gin.Context) {
	const op = "BatchController.TieMarkupType"
	id := c.Param("id")

	log := con.log.With(slog.String("op", op), slog.String("id", id))

	var data tieMarkupType
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

	var markupType models.MarkupType
	if data.MarkupTypeID != nil {
		var existingMarkupType models.MarkupType
		err := con.db.
			Preload("Fields.AssessmentType").
			Where("id = ?", data.MarkupTypeID).
			First(&existingMarkupType).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Warn("markupType not found")
				responses.NotFoundError(c)
				return
			}

			log.Error("failed to find markupType", slog.Any("error", err))
			responses.InternalServerError(c)
			return
		}

		markupType = models.MarkupType{
			Name: existingMarkupType.Name,
		}
		markupType.Fields = make([]models.MarkupTypeField, len(existingMarkupType.Fields))
		for i, field := range existingMarkupType.Fields {
			markupType.Fields[i] = models.MarkupTypeField{
				AssessmentTypeID: field.AssessmentTypeID,
				Name:             field.Name,
				Label:            field.Label,
				GroupID:          field.GroupID,
			}
		}
	} else {
		markupType = models.MarkupType{
			Name: data.Name,
		}
		markupType.Fields = make([]models.MarkupTypeField, len(data.Fields))
		for i, field := range data.Fields {
			markupType.Fields[i] = models.MarkupTypeField{
				AssessmentTypeID: field.AssessmentTypeID,
				Name:             field.Name,
				Label:            field.Label,
				GroupID:          field.GroupID,
			}
		}
	}

	markupType.BatchID = &data.BatchID
	markupType.CreatedAt = time.Now()
	if err := tx.Save(&markupType).Error; err != nil {
		tx.Rollback()
		log.Error("failed to save markup type", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	// Set child id of last MarkupType (parent).
	var lastMarkupType models.MarkupType

	err := tx.
		Where(
			"batch_id = ? AND id != ? AND child_id IS NULL",
			markupType.BatchID,
			markupType.ID,
		).
		First(&lastMarkupType).
		Error

	if err != nil {
		// If parent is not found than it is first MarkupType of batch. No parent exist
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			log.Error("failed to find last markup type field", slog.Any("error", err))
			responses.InternalServerError(c)
			return
		}
	} else {
		lastMarkupType.ChildID = &markupType.ID
		if err := tx.Save(&lastMarkupType).Error; err != nil {
			tx.Rollback()
			log.Error("failed to update parent markup type field", slog.Any("error", err))
			responses.InternalServerError(c)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("failed to commit transaction", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, "OK")
}
