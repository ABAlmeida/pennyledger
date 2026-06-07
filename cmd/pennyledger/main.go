package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/ABAlmeida/pennyledger/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	settings := config.Load()
	addr := settings.HTTPAddr
	logger.Info("starting server", "addr", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
