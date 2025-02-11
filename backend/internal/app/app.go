// Package app defines application model.
package app

import (
	"context"
	"log/slog"
	serverapp "markup/internal/app/server"
	"markup/internal/config"
	"markup/internal/controllers"
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
	Server *serverapp.Server
	Cache  clearer
}

// New creates all dependencies for App and returns new App instance.
func New(
	log *slog.Logger,
	env string,
	port int,
	dbConfig config.DB,
) *App {
	db, err := postgres.New(dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Pass, dbConfig.DBName)
	if err != nil {
		panic(err)
	}

	helloRepo := repos.NewHello(db)

	helloService := services.NewHelloService(log, helloRepo)

	helloCon := controllers.NewHelloController(log, helloService)

	router := server.NewRouter(log, env, helloCon)
	serverApp := serverapp.New(log, port, router)

	return &App{
		Server: serverApp,
	}
}
