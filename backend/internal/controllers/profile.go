package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"markup/internal/domain/models"
	"markup/internal/lib/auth"
	"markup/internal/lib/responses"
	"net/http"
	"time"
)

type Profile struct {
	log *slog.Logger
	db  *gorm.DB
}

func NewProfile(
	log *slog.Logger,
	db *gorm.DB,
) *Profile {
	return &Profile{
		log: log,
		db:  db,
	}
}

type profileResponse struct {
	models.User
	Assessments             []profileResponseAssessment `json:"assessments"`
	AssessmentCount         int64                       `json:"assessment_count"`
	CorrectAssessmentCount  int64                       `json:"correct_assessment_count"`
	AssessmentCount2        int64                       `json:"assessment_count2"`
	CorrectAssessmentCount2 int64                       `json:"correct_assessment_count2"`
}

type profileResponseAssessment struct {
	models.Assessment
	IsCorrect  bool `json:"is_correct"`
	IsEditable bool `json:"is_editable"`
}

func (con *Profile) Me(c *gin.Context) {
	const op = "ProfileController.Me"
	log := con.log.With(slog.String("op", op))

	log.Info("ignore this message")

	user, err := auth.User(c)
	if err != nil {
		responses.UnauthorizedError(c)
		return
	}
	con.db.
		Preload("Assessments.Fields").
		Preload("Assessments.Markup").
		First(&user)

	res := getProfileData(con.db, user)

	c.JSON(http.StatusOK, res)
}

func (con *Profile) Index(c *gin.Context) {
	const op = "ProfileController.Index"
	log := con.log.With(slog.String("op", op))

	log.Info("ignore this message")

	var users []models.User
	con.db.
		Find(&users)

	res := make([]profileResponse, len(users))
	for i, user := range users {
		res[i] = getProfileData(con.db, user)
	}

	c.JSON(http.StatusOK, res)
}

func (con *Profile) Find(c *gin.Context) {
	const op = "ProfileController.Find"
	id := c.Param("id")
	log := con.log.With(slog.String("op", op), slog.String("id", id))

	log.Info("ignore this message")

	var user models.User
	err := con.db.
		Preload("Roles").
		Preload("Assessments.Fields").
		Preload("Assessments.Markup").
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("user not found")
			responses.NotFoundError(c)
			return
		}

		log.Error("failed to find user", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	res := getProfileData(con.db, user)

	c.JSON(http.StatusOK, res)
}

func getProfileData(db *gorm.DB, user models.User) profileResponse {
	// todo count by batch_type_id
	var assessmentCount int64
	db.
		Table("assessments a").
		Select("DISTINCT a.id").
		Joins("JOIN markups m ON a.markup_id = m.id").
		Joins("JOIN batches b ON m.batch_id = b.id").
		Where("b.type_id = ?", 1).
		Where("a.user_id = ?", user.ID).
		Where("a.hash IS NOT NULL").
		Count(&assessmentCount)

	var correctAssessmentCount int64
	db.
		Table("assessments a").
		Select("DISTINCT a.id").
		Joins("JOIN markups m ON a.markup_id = m.id AND a.hash = m.correct_assessment_hash").
		Joins("JOIN batches b ON m.batch_id = b.id").
		Where("b.type_id = ?", 1).
		Where("a.user_id = ?", user.ID).
		Where("a.hash IS NOT NULL").
		Count(&correctAssessmentCount)

	var assessmentCount2 int64
	db.
		Table("assessments a").
		Select("DISTINCT a.id").
		Joins("JOIN markups m ON a.markup_id = m.id").
		Joins("JOIN batches b ON m.batch_id = b.id").
		Where("b.type_id = ?", 2).
		Where("a.user_id = ?", user.ID).
		Where("a.hash IS NOT NULL").
		Count(&assessmentCount2)

	var correctAssessmentCount2 int64
	db.
		Table("assessments a").
		Select("DISTINCT a.id").
		Joins("JOIN markups m ON a.markup_id = m.id AND a.hash = m.correct_assessment_hash").
		Joins("JOIN batches b ON m.batch_id = b.id").
		Where("b.type_id = ?", 2).
		Where("a.user_id = ?", user.ID).
		Where("a.hash IS NOT NULL").
		Count(&correctAssessmentCount2)

	transformedAssessments := make([]profileResponseAssessment, len(user.Assessments))
	for i, a := range user.Assessments {
		isCorrect := a.Hash != nil && a.Markup.CorrectAssessmentHash != nil && *a.Markup.CorrectAssessmentHash == *a.Hash
		transformedAssessments[i] = profileResponseAssessment{
			a,
			isCorrect,
			a.CreatedAt.After(time.Now().Add(-35 * time.Minute)),
		}
	}

	return profileResponse{
		user,
		transformedAssessments,
		assessmentCount,
		correctAssessmentCount,
		assessmentCount2,
		correctAssessmentCount2,
	}
}
