package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	bookmarkcore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
)

func TestUpdateUsesResolvedOwnerAndMasksCrossOwner(t *testing.T) {
	repository := bookmarkcore.NewMemoryRepository(nil)
	service, err := bookmarkcore.NewService(repository)
	if err != nil {
		t.Fatal(err)
	}
	created, err := service.Create(t.Context(), "owner-a", bookmarkcore.CreateInput{URL: "https://example.com/original"})
	if err != nil {
		t.Fatal(err)
	}

	ownerApp, err := newServer(repository, fixedIdentity{owner: "owner-a"}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPatch, "/api/v1/bookmarks/"+created.ID, strings.NewReader(`{"url":"https://example.com/updated#fragment","title":" Updated ","tags":["Cloud","cloud"]}`))
	request.SetPathValue("id", created.ID)
	response := httptest.NewRecorder()
	ownerApp.update(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", response.Code, response.Body.String())
	}
	body := response.Body.String()
	for _, expected := range []string{`"ownerId":"owner-a"`, `"url":"https://example.com/updated"`, `"title":"Updated"`, `"tags":["Cloud"]`} {
		if !strings.Contains(body, expected) {
			t.Fatalf("update response missing %s: %s", expected, body)
		}
	}

	otherApp, err := newServer(repository, fixedIdentity{owner: "owner-b"}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	request = httptest.NewRequest(http.MethodPatch, "/api/v1/bookmarks/"+created.ID, strings.NewReader(`{"url":"https://attacker.example/"}`))
	request.SetPathValue("id", created.ID)
	response = httptest.NewRecorder()
	otherApp.update(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("cross-owner update must be masked as 404, got %d", response.Code)
	}
	if strings.Contains(response.Body.String(), "owner-a") || strings.Contains(response.Body.String(), "updated") {
		t.Fatalf("cross-owner update leaked bookmark data: %s", response.Body.String())
	}
}

func TestUpdateRejectsInvalidJSONAndFailsClosedWithoutIdentity(t *testing.T) {
	app, err := newServer(bookmarkcore.NewMemoryRepository(nil), fixedIdentity{owner: "owner-a"}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPatch, "/api/v1/bookmarks/item", strings.NewReader(`{"url":"https://example.com","unexpected":true}`))
	request.SetPathValue("id", "item")
	response := httptest.NewRecorder()
	app.update(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}

	unavailable, err := newServer(bookmarkcore.NewMemoryRepository(nil), fixedIdentity{err: errors.New("identity backend details")}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	request = httptest.NewRequest(http.MethodPatch, "/api/v1/bookmarks/item", strings.NewReader(`{"url":"https://example.com"}`))
	request.SetPathValue("id", "item")
	response = httptest.NewRecorder()
	unavailable.update(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", response.Code)
	}
	if strings.Contains(response.Body.String(), "backend details") {
		t.Fatal("identity backend details leaked to client")
	}
}

func TestDeleteUsesResolvedOwnerAndMasksCrossOwner(t *testing.T) {
	repository := bookmarkcore.NewMemoryRepository(nil)
	service, err := bookmarkcore.NewService(repository)
	if err != nil {
		t.Fatal(err)
	}
	created, err := service.Create(t.Context(), "owner-a", bookmarkcore.CreateInput{URL: "https://example.com/private"})
	if err != nil {
		t.Fatal(err)
	}

	otherApp, err := newServer(repository, fixedIdentity{owner: "owner-b"}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodDelete, "/api/v1/bookmarks/"+created.ID, nil)
	request.SetPathValue("id", created.ID)
	response := httptest.NewRecorder()
	otherApp.delete(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("cross-owner delete must be masked as 404, got %d", response.Code)
	}

	ownerApp, err := newServer(repository, fixedIdentity{owner: "owner-a"}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	request = httptest.NewRequest(http.MethodDelete, "/api/v1/bookmarks/"+created.ID, nil)
	request.SetPathValue("id", created.ID)
	response = httptest.NewRecorder()
	ownerApp.delete(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", response.Code, response.Body.String())
	}
	if response.Body.Len() != 0 {
		t.Fatalf("204 response must not contain a body: %q", response.Body.String())
	}

	if _, found, err := service.Get(t.Context(), "owner-a", created.ID); err != nil || found {
		t.Fatalf("bookmark remained after delete: found=%v err=%v", found, err)
	}
}

func TestDeleteFailsClosedWithoutIdentity(t *testing.T) {
	app, err := newServer(bookmarkcore.NewMemoryRepository(nil), fixedIdentity{err: errors.New("backend details")}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodDelete, "/api/v1/bookmarks/item", nil)
	request.SetPathValue("id", "item")
	response := httptest.NewRecorder()
	app.delete(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", response.Code)
	}
	if strings.Contains(response.Body.String(), "backend details") {
		t.Fatal("identity backend details leaked to client")
	}
}
