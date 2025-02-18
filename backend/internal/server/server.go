// Package server defines router settings and application endpoints.
package server

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"markup/internal/controllers"
	envpkg "markup/internal/domain/enums/env"
	"markup/internal/server/middleware"
)

// NewRouter sets router mode based on env, registers middleware, defines handlers and options and creates new gin router.
func NewRouter(
	log *slog.Logger,
	env string,
	jwtSecret string,
	db *gorm.DB,
	helloCon *controllers.HelloController,
	markupTypeCon *controllers.MarkupType,
	batchCon *controllers.Batch,
	markupCon *controllers.Markup,
	assessmentCon *controllers.Assessment,
	authCon *controllers.Auth,
	profileCon *controllers.Profile,
	honeypotCon *controllers.Honeypot,
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
		v1protected := api.Group("/v1")
		v1protected.Use(middleware.AuthMiddleware(db, jwtSecret))
		{
			markupTypes := v1protected.Group("/markupTypes")
			{
				markupTypes.GET("", markupTypeCon.Index)
				markupTypes.GET("/:id", markupTypeCon.Find)
				markupTypes.POST("", markupTypeCon.Store)
				markupTypes.PUT("/:id", markupTypeCon.Update)
				markupTypes.DELETE("/:id", markupTypeCon.Destroy)
			}
			batches := v1protected.Group("/batches")
			{
				batches.GET("", batchCon.Index)
				batches.GET("/:id", batchCon.Find)
				batches.POST("", batchCon.Store)
				batches.PUT("/:id", batchCon.Update)
				batches.DELETE("/:id", batchCon.Destroy)

				batches.POST("/:id/markupTypes", batchCon.TieMarkupType)
				batches.PUT("/:id/toggleActive", batchCon.ToggleIsActive)

				batches.GET("/:id/export", batchCon.Export)
			}
			markups := v1protected.Group("/markups")
			{
				markups.GET("", markupCon.Index)
				markups.GET("/:id", markupCon.Find)
			}
			assessments := v1protected.Group("/assessments")
			{
				assessments.GET("", assessmentCon.Index)
				assessments.GET("/:id", assessmentCon.Find)
				assessments.POST("", assessmentCon.Store)
				assessments.PUT("/:id", assessmentCon.Update)
				assessments.DELETE("/:id", assessmentCon.Destroy)

				assessments.POST("/next", assessmentCon.Next)
			}
			auth := v1protected.Group("/auth")
			{
				auth.POST("refresh", authCon.Refresh)
				auth.GET("me", authCon.Me)
			}
			profile := v1protected.Group("/profiles")
			{
				profile.GET("/me", profileCon.Me)
				profile.GET("/:id", profileCon.Find)
				profile.GET("", profileCon.Index)
			}
			honeypots := v1protected.Group("/honeypots")
			{
				honeypots.GET("", honeypotCon.Index)
				honeypots.POST("/:id", honeypotCon.Store)
			}
		}

		v1public := api.Group("/v1")
		{
			auth := v1public.Group("/auth")
			{
				auth.POST("register", authCon.Register)
				auth.POST("login", authCon.Login)
			}
		}
	}

	return r
}
