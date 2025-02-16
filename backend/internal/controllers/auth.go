package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log/slog"
	"markup/internal/domain/enums/roles"
	"markup/internal/domain/models"
	"markup/internal/lib/jwt"
	"markup/internal/lib/responses"
	"net/http"
	"strings"
)

type Auth struct {
	log       *slog.Logger
	db        *gorm.DB
	jwtSecret string
}

func NewAuth(
	log *slog.Logger,
	db *gorm.DB,
	jwtSecret string,
) *Auth {
	return &Auth{
		log:       log,
		db:        db,
		jwtSecret: jwtSecret,
	}
}

type register struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (con *Auth) Register(c *gin.Context) {
	const op = "AuthController.Register"
	log := con.log.With(slog.String("op", op))

	var data register
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	user.Email = data.Email

	passHash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Password = string(passHash)

	// Find the roles
	var assessorRole models.Role
	if err := con.db.Find(&assessorRole, roles.Assessor).Error; err != nil {
		responses.InternalServerError(c)
		return
	}

	user.Roles = []models.Role{
		assessorRole,
	}

	if err := con.db.Create(&user).Error; err != nil {
		log.Error("failed to create user", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, "OK")
}

type login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (con *Auth) Login(c *gin.Context) {
	const op = "AuthController.Login"
	log := con.log.With(slog.String("op", op))

	var data login
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := con.db.Where("email = ?", data.Email).First(&user).Error; err != nil {
		log.Warn("user not found", slog.Any("error", err))
		responses.UnauthorizedError(c)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		log.Error("invalid password", slog.Any("error", err))
		responses.UnauthorizedError(c)
		return
	}

	token, err := jwt.GenerateToken(user.ID, con.jwtSecret)
	if err != nil {
		log.Error("failed to generate token", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (con *Auth) Refresh(c *gin.Context) {
	const op = "AuthController.Refresh"
	log := con.log.With(slog.String("op", op))

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := jwt.ValidateToken(tokenString, con.jwtSecret)
	if err != nil {
		log.Error("invalid token", slog.Any("error", err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	newToken, err := jwt.GenerateToken(claims.UserID, con.jwtSecret)
	if err != nil {
		log.Error("failed to generate token", slog.Any("error", err))
		responses.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}

func (con *Auth) Me(c *gin.Context) {
	const op = "AuthController.Me"
	log := con.log.With(slog.String("op", op))

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := jwt.ValidateToken(tokenString, con.jwtSecret)
	if err != nil {
		log.Error("invalid token", slog.Any("error", err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	var user models.User
	if err := con.db.Preload("Roles").First(&user, claims.UserID).Error; err != nil {
		log.Error("user not found", slog.Any("error", err))
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
