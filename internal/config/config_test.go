package config

import "testing"

func TestLoadUsesDefaultHTTPAddr(t *testing.T) {
	t.Setenv("HTTP_ADDR", "")

	settings := Load()

	if settings.HTTPAddr != ":8080" {
		t.Fatalf("expected HTTPAddr to be :8080, got %q", settings.HTTPAddr)
	}
}

func TestLoadUsesHTTPAddrFromEnvironment(t *testing.T) {
	t.Setenv("HTTP_ADDR", ":9090")

	settings := Load()

	if settings.HTTPAddr != ":9090" {
		t.Fatalf("expected HTTPAddr to be :9090, got %q", settings.HTTPAddr)
	}
}
