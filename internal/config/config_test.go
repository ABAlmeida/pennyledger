package config

import (
	"testing"
	"time"
)

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

func TestLoadUsesDefaultShutdownTimeout(t *testing.T) {
	t.Setenv("SHUTDOWN_TIMEOUT", "")
	settings := Load()
	if settings.ShutdownTimeout != 5*time.Second {
		t.Fatalf("expected ShutdownTimeout to be 5s, got %q", settings.ShutdownTimeout)
	}
}

func TestLoadUsesShutdownTimeoutFromEnvironment(t *testing.T) {
	t.Setenv("SHUTDOWN_TIMEOUT", "10s")
	settings := Load()
	if settings.ShutdownTimeout != 10*time.Second {
		t.Fatalf("expected ShutdownTimeout to be 10s, got %q", settings.ShutdownTimeout)
	}
}
