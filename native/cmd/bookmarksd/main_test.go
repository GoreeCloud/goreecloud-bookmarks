package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	bookmarkcore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
)

type fixedIdentity struct {
	owner string
	err   error
}

func (i fixedIdentity) Resolve(*http.Request) (string, error) { return i.owner, i.err }

func TestNewServerRequiresRepositoryAndIdentity(t *testing.T) {
	if _, err := newServer(nil, fixedIdentity{owner: "owner"}, "memory"); err == nil {
		t.Fatal("expected repository requirement")
	}
	if _, err := newServer(bookmarkcore.NewMemoryRepository(nil), nil, "memory"); err == nil {
		t.Fatal("expected identity resolver requirement")
	}
}

func TestListUsesResolvedOwnerThroughService(t *testing.T) {
	repository := bookmarkcore.NewMemoryRepository(nil)
	service, err := bookmarkcore.NewService(repository)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Create(t.Context(), "owner-a", bookmarkcore.CreateInput{URL: "https://example.com/a"}); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Create(t.Context(), "owner-b", bookmarkcore.CreateInput{URL: "https://example.com/b"}); err != nil {
		t.Fatal(err)
	}

	app, err := newServer(repository, fixedIdentity{owner: "owner-a"}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/api/v1/bookmarks", nil)
	response := httptest.NewRecorder()
	app.list(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	body := response.Body.String()
	if !strings.Contains(body, "https://example.com/a") || strings.Contains(body, "https://example.com/b") {
		t.Fatalf("owner isolation failed in runtime response: %s", body)
	}
}

func TestListFailsClosedWithoutUsableIdentity(t *testing.T) {
	repository := bookmarkcore.NewMemoryRepository(nil)

	unavailable, err := newServer(repository, fixedIdentity{err: errors.New("identity backend unavailable")}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/api/v1/bookmarks", nil)
	response := httptest.NewRecorder()
	unavailable.list(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected unavailable identity to return 503, got %d", response.Code)
	}

	empty, err := newServer(repository, fixedIdentity{owner: "   "}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	response = httptest.NewRecorder()
	empty.list(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected empty owner to return 401, got %d", response.Code)
	}
}

func TestHealthDoesNotClaimPersistentReadiness(t *testing.T) {
	app, err := newServer(bookmarkcore.NewMemoryRepository(nil), fixedIdentity{owner: "owner"}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()
	app.health(response, request)
	body := response.Body.String()
	if !strings.Contains(body, `"persistent_store_ready":false`) || !strings.Contains(body, `"production_approved":false`) {
		t.Fatalf("health response overstated readiness: %s", body)
	}
}
