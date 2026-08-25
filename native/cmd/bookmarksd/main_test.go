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

func TestGetUsesResolvedOwnerAndMasksCrossOwner(t *testing.T) {
	repository := bookmarkcore.NewMemoryRepository(nil)
	service, err := bookmarkcore.NewService(repository)
	if err != nil {
		t.Fatal(err)
	}
	created, err := service.Create(t.Context(), "owner-a", bookmarkcore.CreateInput{URL: "https://example.com/private"})
	if err != nil {
		t.Fatal(err)
	}

	ownerApp, err := newServer(repository, fixedIdentity{owner: "owner-a"}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/api/v1/bookmarks/"+created.ID, nil)
	request.SetPathValue("id", created.ID)
	response := httptest.NewRecorder()
	ownerApp.get(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "https://example.com/private") {
		t.Fatalf("owner lookup failed: status=%d body=%s", response.Code, response.Body.String())
	}

	otherApp, err := newServer(repository, fixedIdentity{owner: "owner-b"}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	request = httptest.NewRequest(http.MethodGet, "/api/v1/bookmarks/"+created.ID, nil)
	request.SetPathValue("id", created.ID)
	response = httptest.NewRecorder()
	otherApp.get(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("cross-owner lookup must be masked as 404, got %d", response.Code)
	}
	if strings.Contains(response.Body.String(), "owner-a") || strings.Contains(response.Body.String(), "private") {
		t.Fatalf("cross-owner lookup leaked bookmark data: %s", response.Body.String())
	}
}

func TestGetFailsClosedWithoutUsableIdentity(t *testing.T) {
	app, err := newServer(bookmarkcore.NewMemoryRepository(nil), fixedIdentity{err: errors.New("identity backend details")}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/api/v1/bookmarks/missing", nil)
	request.SetPathValue("id", "missing")
	response := httptest.NewRecorder()
	app.get(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected unavailable identity to return 503, got %d", response.Code)
	}
	if strings.Contains(response.Body.String(), "backend details") {
		t.Fatal("identity backend details leaked to client")
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

func TestCreateUsesResolvedOwnerAndNormalizesInput(t *testing.T) {
	repository := bookmarkcore.NewMemoryRepository(nil)
	app, err := newServer(repository, fixedIdentity{owner: "owner-a"}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/v1/bookmarks", strings.NewReader(`{"url":"https://example.com/page#fragment","title":" Example ","tags":["Cloud","cloud"," search "]}`))
	response := httptest.NewRecorder()
	app.create(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", response.Code, response.Body.String())
	}
	body := response.Body.String()
	for _, expected := range []string{`"ownerId":"owner-a"`, `"url":"https://example.com/page"`, `"title":"Example"`, `"tags":["Cloud","search"]`} {
		if !strings.Contains(body, expected) {
			t.Fatalf("create response missing %s: %s", expected, body)
		}
	}
}

func TestCreateRejectsMalformedAndUnknownJSON(t *testing.T) {
	app, err := newServer(bookmarkcore.NewMemoryRepository(nil), fixedIdentity{owner: "owner-a"}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	for _, body := range []string{
		`{"url":"https://example.com","unexpected":true}`,
		`{"url":"https://example.com"} {"url":"https://example.org"}`,
		`{"url":`,
	} {
		request := httptest.NewRequest(http.MethodPost, "/api/v1/bookmarks", strings.NewReader(body))
		response := httptest.NewRecorder()
		app.create(response, request)
		if response.Code != http.StatusBadRequest {
			t.Fatalf("expected invalid JSON to return 400, got %d for %q", response.Code, body)
		}
		if strings.Contains(response.Body.String(), "unexpected") {
			t.Fatal("decoder internals leaked to client")
		}
	}
}

func TestCreateFailsClosedWithoutIdentity(t *testing.T) {
	app, err := newServer(bookmarkcore.NewMemoryRepository(nil), fixedIdentity{err: errors.New("backend details")}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/v1/bookmarks", strings.NewReader(`{"url":"https://example.com"}`))
	response := httptest.NewRecorder()
	app.create(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", response.Code)
	}
	if strings.Contains(response.Body.String(), "backend details") {
		t.Fatal("identity backend details leaked to client")
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

func TestRuntimeRepositoryDefaultsToMemoryDevelopment(t *testing.T) {
	repository, mode, err := selectRuntimeRepository(func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	if repository == nil || mode != defaultRepositoryMode {
		t.Fatalf("unexpected default repository selection: mode=%q repository=%T", mode, repository)
	}
}

func TestRuntimeRepositoryRejectsUnknownMode(t *testing.T) {
	_, _, err := selectRuntimeRepository(func(key string) string {
		if key == repositoryModeEnv {
			return "sqlite"
		}
		return ""
	})
	if err == nil || !strings.Contains(err.Error(), "unsupported bookmark repository mode") {
		t.Fatalf("expected bounded unsupported-mode error, got %v", err)
	}
}

func TestPostgresDevelopmentSelectionFailsClosed(t *testing.T) {
	_, _, err := selectRuntimeRepository(func(key string) string {
		if key == repositoryModeEnv {
			return postgresRepositoryMode
		}
		return ""
	})
	if err == nil || !strings.Contains(err.Error(), "requires an explicit database URL") {
		t.Fatalf("expected missing database URL failure, got %v", err)
	}

	secretURL := "postgres://user:secret@example.invalid/bookmarks"
	_, _, err = selectRuntimeRepository(func(key string) string {
		switch key {
		case repositoryModeEnv:
			return postgresRepositoryMode
		case postgresDatabaseURLEnv:
			return secretURL
		default:
			return ""
		}
	})
	if err == nil || !strings.Contains(err.Error(), "driver is not integrated") {
		t.Fatalf("expected unavailable driver failure, got %v", err)
	}
	if strings.Contains(err.Error(), secretURL) || strings.Contains(err.Error(), "secret") {
		t.Fatalf("database URL leaked through configuration error: %v", err)
	}
}
