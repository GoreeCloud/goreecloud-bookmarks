package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	bookmarkcore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
)

func TestReadinessFailsClosedWhileProductionDependenciesAreIncomplete(t *testing.T) {
	app, err := newServer(bookmarkcore.NewMemoryRepository(nil), fixedIdentity{owner: "owner"}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	response := httptest.NewRecorder()
	app.ready(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected readiness to fail closed with 503, got %d", response.Code)
	}
	body := response.Body.String()
	for _, expected := range []string{
		`"ready":false`,
		`"identity_adapter_ready":false`,
		`"persistent_store_ready":false`,
		`"store_mode":"memory-development"`,
		`"production_approved":false`,
	} {
		if !strings.Contains(body, expected) {
			t.Fatalf("readiness response missing %s: %s", expected, body)
		}
	}
}

func TestHealthRemainsLivenessSignalWhileReadinessIsBlocked(t *testing.T) {
	app, err := newServer(bookmarkcore.NewMemoryRepository(nil), fixedIdentity{owner: "owner"}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}

	healthResponse := httptest.NewRecorder()
	app.health(healthResponse, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if healthResponse.Code != http.StatusOK {
		t.Fatalf("expected liveness health to remain 200, got %d", healthResponse.Code)
	}

	readyResponse := httptest.NewRecorder()
	app.ready(readyResponse, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if readyResponse.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected readiness to remain blocked, got %d", readyResponse.Code)
	}
}
