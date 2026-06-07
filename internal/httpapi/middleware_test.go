package httpapi

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDMiddlewareAddsRequestId(t *testing.T) {
	handler := requestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	requestID := response.Header().Get(requestIDHeader)
	if requestID == "" {
		t.Fatal("expected response to include request id")
	}
}

func TestRecoverMiddlewareReturnsInternalServerError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	handler := recoverMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, response.Code)
	}
}

func TestLoggingMiddlewareRecordsStatusCode(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	handler := loggingMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	request := httptest.NewRequest(http.MethodPost, "/test", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}
}
