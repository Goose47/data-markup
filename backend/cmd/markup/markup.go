// application entrypoint
package main

import (
	"log/slog"
	apppkg "markup/internal/app"
	"markup/internal/config"
	"markup/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.Env)
	app := apppkg.New(log, cfg.Env, cfg.Port, cfg.DB)

	app.TaskManager.Run()

	err := app.Server.Serve()
	log.Error("application has stopped: %s", slog.Any("error", err))
}
