package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"markup/internal/domain/models"
	"markup/internal/lib/jwt"
	"net/http"
	"strings"
)

func AuthMiddleware(db *gorm.DB, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Fetch the user with roles from the database
		var user models.User
		if err := db.Preload("Roles").First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			c.Abort()
			return
		}

		// Store the user in the Gin context
		c.Set("user", user)
		c.Next()
	}
}
