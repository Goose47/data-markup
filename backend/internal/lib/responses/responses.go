package responses

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

func Pagination(
	data any,
	total int64,
	page int,
	perPage int,
) gin.H {
	pagesTotal := math.Ceil(float64(total) / float64(perPage))
	return gin.H{
		"data":        data,
		"page":        page,
		"per_page":    perPage,
		"pages_total": pagesTotal,
	}
}

func InternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "internal server error",
	})
}
