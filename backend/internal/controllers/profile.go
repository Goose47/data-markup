package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"markup/internal/domain/models"
	"markup/internal/lib/auth"
	"markup/internal/lib/responses"
	"net/http"
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
	AssessmentCount         int64 `json:"assessment_count"`
	CorrectAssessmentCount  int64 `json:"correct_assessment_count"`
	AssessmentCount2        int64 `json:"assessment_count2"`
	CorrectAssessmentCount2 int64 `json:"correct_assessment_count2"`
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

	// todo count by batch_type_id
	var assessmentCount int64
	con.db.
		Table("assessments a").
		Select("DISTINCT a.id").
		Joins("JOIN markups m ON a.markup_id = m.id").
		Joins("JOIN batches b ON m.batch_id = b.id").
		Where("b.type_id = ?", 1).
		Where("a.user_id = ?", user.ID).
		Count(&assessmentCount)

	var correctAssessmentCount int64
	con.db.
		Table("assessments a").
		Select("DISTINCT a.id").
		Joins("JOIN markups m ON a.markup_id = m.id AND a.hash = m.correct_assessment_hash").
		Joins("JOIN batches b ON m.batch_id = b.id").
		Where("b.type_id = ?", 1).
		Where("a.user_id = ?", user.ID).
		Count(&correctAssessmentCount)

	var assessmentCount2 int64
	con.db.
		Table("assessments a").
		Select("DISTINCT a.id").
		Joins("JOIN markups m ON a.markup_id = m.id").
		Joins("JOIN batches b ON m.batch_id = b.id").
		Where("b.type_id = ?", 2).
		Where("a.user_id = ?", user.ID).
		Count(&assessmentCount2)

	var correctAssessmentCount2 int64
	con.db.
		Table("assessments a").
		Select("DISTINCT a.id").
		Joins("JOIN markups m ON a.markup_id = m.id AND a.hash = m.correct_assessment_hash").
		Joins("JOIN batches b ON m.batch_id = b.id").
		Where("b.type_id = ?", 2).
		Where("a.user_id = ?", user.ID).
		Count(&correctAssessmentCount2)

	res := profileResponse{
		user,
		assessmentCount,
		correctAssessmentCount,
		assessmentCount2,
		correctAssessmentCount2,
	}

	c.JSON(http.StatusOK, res)
}

func (con *Profile) Index(c *gin.Context) {
	//const op = "ProfileController.Index"
	//log := con.log.With(slog.String("op", op))
	//
	//user, err := auth.User(c)
	//if err != nil {
	//	responses.UnauthorizedError(c)
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{})
}
