package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type readinessChecker interface {
	Ping(ctx context.Context) error
}

func NewRouter(logger *slog.Logger, checker readinessChecker, accountService accountService) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthzHandler)
	mux.HandleFunc("GET /readyz", readyzHandler(checker))
	mux.HandleFunc("POST /v1/accounts", createAccountHandler(accountService))
	mux.HandleFunc("GET /v1/accounts/{id}", getAccountHandler(accountService))

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

func readyzHandler(checker readinessChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := checker.Ping(ctx); err != nil {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready\n"))
	}
}
