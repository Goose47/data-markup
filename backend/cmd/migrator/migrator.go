package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"markup/internal/config"
)

func main() {
	var migrationsPath string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations folder")

	cfg := config.MustLoad()

	if migrationsPath == "" {
		panic("migrations path can not be empty")
	}

	migration, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.DB.User,
			cfg.DB.Pass,
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.DBName,
		),
	)
	if err != nil {
		panic(err)
	}

	if err := migration.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("everything is up to date")

			return
		}

		panic(err)
	}

	fmt.Println("Migrations applied successfully!")
}

func rollbackMigrations(steps int) {
	migration, err := migrate.New(
		"file://db/migrations",
		"postgres://user:password@localhost:5432/yourdb?sslmode=disable",
	)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize rollback: %v", err))
	}

	if err := migration.Steps(-steps); err != nil && err != migrate.ErrNoChange {
		panic(fmt.Sprintf("failed to rollback %d migrations: %v", steps, err))
	}
	fmt.Printf("Rolled back %d migrations successfully!\n", steps)
}
