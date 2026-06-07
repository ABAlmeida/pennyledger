package httpapi

import (
	"log/slog"
	"net/http"
)

func NewRouter(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthzHandler)

	return recoverMiddleware(logger)(
		requestIDMiddleware(
			loggingMiddleware(logger)(mux),
		),
	)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}
