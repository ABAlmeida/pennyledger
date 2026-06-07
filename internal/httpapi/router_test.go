package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthzReturnsOk(t *testing.T) {
	router := NewRouter()

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
