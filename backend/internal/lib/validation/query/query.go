// Package query provides functions to retrieve and validate http query params from gin context.
package query

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

// DefaultInt retrieves key from query param, converts it to int and sends response if fails.
func DefaultInt(c *gin.Context, log *slog.Logger, key string, defaultValue string) (int, error) {
	const op = "validation.query.DefaultInt"

	value, err := strconv.Atoi(c.DefaultQuery(key, defaultValue))
	if err != nil {
		message := fmt.Sprintf("wrong %s parameter", key)
		log.Warn(message, slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return value, nil
}

// Int retrieves key from query param, converts it to int and returns nil if fails.
func Int(c *gin.Context, key string) *int {
	value, err := strconv.Atoi(c.Query(key))
	if err != nil {
		return nil
	}

	return &value
}
