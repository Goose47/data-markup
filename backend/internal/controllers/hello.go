package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type HelloProvider interface {
	Hello(ctx context.Context) (string, error)
}

type HelloController struct {
	log   *slog.Logger
	hello HelloProvider
}

func NewHelloController(
	log *slog.Logger,
	hello HelloProvider,
) *HelloController {
	return &HelloController{
		log:   log,
		hello: hello,
	}
}

func (con *HelloController) Hello(c *gin.Context) {
	ctx := c.Request.Context()
	res, err := con.hello.Hello(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed",
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
