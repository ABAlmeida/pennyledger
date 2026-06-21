package httpapi

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ABAlmeida/pennyledger/internal/accounts"
)

type fakeReadinessChecker struct {
	err error
}

func (f fakeReadinessChecker) Ping(ctx context.Context) error {
	return f.err
}

type fakeAccountService struct{}

func (f fakeAccountService) CreateAccount(ctx context.Context, input accounts.CreateAccountInput) (accounts.Account, error) {
	return accounts.Account{}, nil
}

func (f fakeAccountService) GetAccountByID(ctx context.Context, id string) (accounts.Account, error) {
	return accounts.Account{}, nil
}

func TestHealthzReturnsOk(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := NewRouter(logger, fakeReadinessChecker{}, fakeAccountService{})

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if response.Body.String() != "ok\n" {
		t.Fatalf("expected body %q, got %q", "ok\n", response.Body.String())
	}
}

func TestReadyzReturnsOkWhenDatabaseIsReachable(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := NewRouter(logger, fakeReadinessChecker{}, fakeAccountService{})

	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if response.Body.String() != "ready\n" {
		t.Fatalf("expected body %q, got %q", "ready\n", response.Body.String())
	}
}

func TestReadyzReturnsServiceUnavailableWhenDatabaseIsUnreachable(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := NewRouter(logger, fakeReadinessChecker{err: io.ErrUnexpectedEOF}, fakeAccountService{})

	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, response.Code)
	}

	if response.Body.String() != "not ready\n" {
		t.Fatalf("expected body %q, got %q", "not ready\n", response.Body.String())
	}
}
