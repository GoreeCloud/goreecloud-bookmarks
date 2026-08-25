package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	bookmarkcore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
	identitycore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/identity"
)

func TestLibraryRendersOnlyResolvedOwnerBookmarks(t *testing.T) {
	app, err := newServer(bookmarkcore.NewMemoryRepository(nil), identitycore.DevelopmentHeaderResolver{}, defaultRepositoryMode)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := app.bookmarks.Create(t.Context(), "owner-a", bookmarkcore.CreateInput{URL: "https://example.com/a", Title: "Owner A", Tags: []string{"Read"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := app.bookmarks.Create(t.Context(), "owner-b", bookmarkcore.CreateInput{URL: "https://example.com/b", Title: "Owner B"}); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(identitycore.DevelopmentOwnerHeader, "owner-a")
	response := httptest.NewRecorder()
	app.library(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/html") {
		t.Fatalf("content type = %q", got)
	}
	body := response.Body.String()
	if !strings.Contains(body, "Owner A") || !strings.Contains(body, "Read") {
		t.Fatalf("owner bookmark missing from library: %s", body)
	}
	if strings.Contains(body, "Owner B") || strings.Contains(body, "example.com/b") {
		t.Fatalf("cross-owner bookmark leaked into library: %s", body)
	}
}

func TestLibraryFailsClosedWhenIdentityUnavailable(t *testing.T) {
	app, err := newServer(bookmarkcore.NewMemoryRepository(nil), identitycore.UnavailableResolver{}, defaultRepositoryMode)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	app.library(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusServiceUnavailable)
	}
	if !strings.Contains(response.Body.String(), "Bookmark library unavailable") {
		t.Fatalf("bounded unavailable state missing: %s", response.Body.String())
	}
	if strings.Contains(response.Body.String(), "X-GoreeCloud-Development-Owner") {
		t.Fatal("development identity mechanism leaked into user-facing unavailable state")
	}
}

func TestLibraryStylesAreLocalCSS(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/assets/library.css", nil)
	response := httptest.NewRecorder()
	libraryStyles(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/css") {
		t.Fatalf("content type = %q", got)
	}
	body := response.Body.String()
	for _, required := range []string{".app-shell", "prefers-reduced-motion", "forced-colors", "prefers-color-scheme"} {
		if !strings.Contains(body, required) {
			t.Fatalf("library stylesheet missing %q", required)
		}
	}
}

func TestSecurityHeadersDoNotForceJSONOnHTML(t *testing.T) {
	handler := securityHeaders(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
	if got := response.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Fatalf("content type = %q", got)
	}
	if response.Header().Get("Content-Security-Policy") == "" {
		t.Fatal("content security policy missing")
	}
}
