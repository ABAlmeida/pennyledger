package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ABAlmeida/pennyledger/internal/config"
	"github.com/ABAlmeida/pennyledger/internal/httpapi"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	settings := config.Load()

	router := httpapi.NewRouter(logger)

	server := &http.Server{
		Addr:    settings.HTTPAddr,
		Handler: router,
	}
	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("starting server", "addr", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdownSignals)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped", "error", err)
			os.Exit(1)
		}

	case sig := <-shutdownSignals:
		logger.Info("shutdown started", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), settings.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("shutdown failed", "error", err)
			os.Exit(1)
		}

		if err := <-serverErrors; err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped", "error", err)
			os.Exit(1)
		}

		logger.Info("shutdown complete")
	}
}
