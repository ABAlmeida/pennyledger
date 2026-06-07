package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/ABAlmeida/pennyledger/internal/config"
	"github.com/ABAlmeida/pennyledger/internal/httpapi"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	settings := config.Load()

	router := httpapi.NewRouter(logger)

	addr := settings.HTTPAddr
	logger.Info("starting server", "addr", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
