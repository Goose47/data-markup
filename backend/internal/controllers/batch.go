package controllers

import (
	"archive/zip"
	"bytes"
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
	"slices"
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
		Where("m.batch_id = ? AND a.hash IS NOT NULL", batch.ID).
		Count(&assessmentCount)

	var correctAssessmentCount int64
	con.db.
		Table("assessments a").
		Select("DISTINCT a.id").
		Joins("JOIN markups m ON a.markup_id = m.id AND a.hash = m.correct_assessment_hash").
		Where("m.batch_id = ? AND a.hash IS NOT NULL", batch.ID).
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

func (con *Batch) Export(c *gin.Context) {
	const op = "BatchController.Export"
	id := c.Param("id")
	log := con.log.With(slog.String("op", op), slog.String("id", id))

	var batch models.Batch
	err := con.db.
		Where("id = ?", id).
		Preload("MarkupTypes.Fields").
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

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for _, mt := range batch.MarkupTypes {
		// find markups, corresponding to current markupType
		var markups []models.Markup
		err := con.db.
			Select("DISTINCT m.*").
			Table("markups m").
			Joins("LEFT JOIN assessments a ON a.markup_id = m.id AND a.hash = m.correct_assessment_hash").
			Joins("LEFT JOIN assessment_fields af ON af.assessment_id = a.id").
			Joins("LEFT JOIN markup_type_fields mtf ON af.markup_type_field_id = mtf.id").
			Where("mtf.markup_type_id = ? AND m.status_id = ?", mt.ID, markupStatus.Processed).
			Preload("Assessments.Fields").
			Find(&markups).Error
		if err != nil {
			log.Error("failed to find markups", slog.Any("error", err))
			responses.InternalServerError(c)
			return
		}

		if len(markups) == 0 {
			continue
		}

		// csv file column names
		csvHeaders := make([]string, len(mt.Fields)+1)
		csvHeaders[0] = "markup_data"

		// markupTypeField ids for corresponding column name in csv column
		markupTypeFieldIDs := make([]uint, len(mt.Fields))

		// fill csv column names
		for i, field := range mt.Fields {
			var label string
			if field.Label != nil {
				label = *field.Label
			}
			var name string
			if field.Label != nil {
				name = *field.Name
			}
			csvHeaders[i+1] = fmt.Sprintf("%d.%s.%s", field.GroupID, label, name)
			markupTypeFieldIDs[i] = field.ID
		}

		// create csv data
		csvData := make([][]string, 0, len(markups)+1)
		csvData = append(csvData, csvHeaders)

		for _, markup := range markups {
			if len(markup.Assessments) == 0 {
				log.Warn("markup is empty", slog.Int("markup_id", int(markup.ID)))
				continue
			}

			nextRow := make([]string, 0, len(csvHeaders))
			nextRow = append(nextRow, markup.Data)

			//if assessment has asssessmentField corresponding to markupTypeField, write "+"
		loop:
			for _, nextMTFID := range markupTypeFieldIDs {
				if slices.ContainsFunc(markup.Assessments[0].Fields, func(n models.AssessmentField) bool {
					return n.MarkupTypeFieldID == nextMTFID
				}) {
					nextRow = append(nextRow, "+")
					continue loop
				}
				nextRow = append(nextRow, "-")
			}

			csvData = append(csvData, nextRow)
		}

		filename := fmt.Sprintf("%d-%s-%s.csv", mt.ID, mt.Name, mt.CreatedAt.Format("2006-01-02"))

		if err := generateCSV(zipWriter, filename, csvData); err != nil {
			log.Error("Error creating ZIP", slog.Any("error", err))
			responses.InternalServerError(c)
			return
		}
	}

	// Закрываем архив
	if err := zipWriter.Close(); err != nil {
		log.Error("Error finalizing ZIP", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	// Устанавливаем заголовки для скачивания
	zipFilename := "archive_" + time.Now().Format("20060102_150405") + ".zip"
	c.Header("Content-Disposition", "attachment; filename="+zipFilename)
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Length", fmt.Sprintf("%d", buf.Len()))

	// Отправляем ZIP-файл в ответе
	c.Data(http.StatusOK, "application/zip", buf.Bytes())

	c.JSON(http.StatusOK, "OK")
}

// generateCSV создает CSV-файл и записывает его в переданный *zip.Writer
func generateCSV(zipWriter *zip.Writer, filename string, data [][]string) error {
	// Создаем файл внутри архива
	fileWriter, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	// Создаем CSV writer
	csvWriter := csv.NewWriter(fileWriter)
	defer csvWriter.Flush()

	// Записываем данные
	for _, row := range data {
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}
