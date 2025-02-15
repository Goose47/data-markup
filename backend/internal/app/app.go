// Package app defines application model.
package app

import (
	"context"
	"log/slog"
	serverapp "markup/internal/app/server"
	"markup/internal/background"
	"markup/internal/config"
	"markup/internal/controllers"
	//"markup/internal/db/mysql"
	"markup/internal/db/postgres"
	"markup/internal/repos"
	"markup/internal/server"
	"markup/internal/services"
)

type clearer interface {
	Clear(ctx context.Context)
}

// App represents application.
type App struct {
	Server      *serverapp.Server
	TaskManager *background.TaskManager
	Cache       clearer
}

// New creates all dependencies for App and returns new App instance.
func New(
	log *slog.Logger,
	env string,
	port int,
	dbConfig config.DB,
) *App {
	//db, err := mysql.New(dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Pass, dbConfig.DBName)
	db, err := postgres.New(dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Pass, dbConfig.DBName)
	if err != nil {
		panic(err)
	}

	helloRepo := repos.NewHello(db)

	helloService := services.NewHelloService(log, helloRepo)

	helloCon := controllers.NewHelloController(log, helloService)
	markupTypeCon := controllers.NewMarkupType(log, db)
	batchCon := controllers.NewBatch(log, db)
	markupCon := controllers.NewMarkup(log, db)
	assessmentCon := controllers.NewAssessment(log, db)

	router := server.NewRouter(
		log,
		env,
		helloCon,
		markupTypeCon,
		batchCon,
		markupCon,
		assessmentCon,
	)
	serverApp := serverapp.New(log, port, router)

	tm := background.NewTaskManager(log, db)

	return &App{
		Server:      serverApp,
		TaskManager: tm,
	}
}
