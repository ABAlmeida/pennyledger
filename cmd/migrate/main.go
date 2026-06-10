package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/ABAlmeida/pennyledger/internal/config"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	settings := config.Load()

	migrations, err := migrate.New(
		"file://db/migrations",
		settings.DatabaseURL,
	)

	if err != nil {
		logger.Error("create migrator failed", "error", err)
		os.Exit(1)
	}

	if err := migrations.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("migrations failed", "error", err)
		os.Exit(1)
	}

	logger.Info("migrations complete")
}
