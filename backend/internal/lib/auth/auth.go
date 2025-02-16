package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"markup/internal/domain/models"
)

func User(c *gin.Context) (models.User, error) {
	const op = "auth.User"
	// Retrieve the user from the Gin context
	user, exists := c.Get("user")
	if !exists {
		return models.User{}, fmt.Errorf("%s: %s", op, "user not found in context")
	}

	// Type assert the user to models.User
	userObj, ok := user.(models.User)
	if !ok {
		return models.User{}, fmt.Errorf("%s: %s", op, "invalid user type in context")
	}

	return userObj, nil
}
