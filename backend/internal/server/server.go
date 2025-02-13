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
	markupTypeCon *controllers.MarkupType,
	batchCon *controllers.Batch,
	markupCon *controllers.Markup,
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
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			markupTypes := v1.Group("/markupTypes")
			{
				markupTypes.GET("", markupTypeCon.Index)
				markupTypes.GET("/:id", markupTypeCon.Find)
				markupTypes.POST("", markupTypeCon.Store)
				markupTypes.PUT("/:id", markupTypeCon.Update)
				markupTypes.DELETE("/:id", markupTypeCon.Destroy)
			}
			batches := v1.Group("/batches")
			{
				batches.GET("", batchCon.Index)
				batches.GET("/:id", batchCon.Find)
				batches.POST("", batchCon.Store)
				batches.PUT("/:id", batchCon.Update)
				batches.DELETE("/:id", batchCon.Destroy)

				batches.POST("/:id/markupTypes", batchCon.TieMarkupType)
			}
			markups := v1.Group("/markups")
			{
				markups.GET("", markupCon.Index)
				markups.GET("/:id", markupCon.Find)
			}
		}
	}

	return r
}
