// Package query provides functions to retrieve and validate http query params from gin context.
package query

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

// Int retrieves key from query param, converts it to int and sends response if fails.
func Int(c *gin.Context, log *slog.Logger, key string, defaultValue string) (int, error) {
	const op = "validation.query.Int"

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
