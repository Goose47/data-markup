// Package server defines router settings and application endpoints.
package server

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"markup/internal/controllers"
	envpkg "markup/internal/domain/enums/env"
)

// NewRouter sets router mode based on env, registers middleware, defines handlers and options and creates new gin router.
func NewRouter(
	log *slog.Logger,
	env string,
	helloCon *controllers.HelloController,
) *gin.Engine {
	var mode string
	switch env {
	case envpkg.Local:
	case envpkg.Dev:
		mode = gin.DebugMode
	case envpkg.Prod:
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	r := gin.New()

	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	r.Use(gin.Recovery())

	r.GET("hello", helloCon.Hello)

	return r
}
