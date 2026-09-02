package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	bookmarkcore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
)

type failingGetRepository struct {
	*bookmarkcore.MemoryRepository
	err error
}

func (r failingGetRepository) Get(context.Context, string, string) (bookmarkcore.Bookmark, bool, error) {
	return bookmarkcore.Bookmark{}, false, r.err
}

func newBookmarkCollectionFixture(t *testing.T) (server, bookmarkcore.Bookmark, string) {
	t.Helper()
	repository := bookmarkcore.NewMemoryRepository(nil)
	service, err := bookmarkcore.NewService(repository)
	if err != nil {
		t.Fatal(err)
	}
	bookmark, err := service.Create(t.Context(), "owner-a", bookmarkcore.CreateInput{URL: "https://example.com/assigned"})
	if err != nil {
		t.Fatal(err)
	}
	app, err := newServer(repository, fixedIdentity{owner: "owner-a"}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}
	collection, err := app.collections.Create("owner-a", "Research", "")
	if err != nil {
		t.Fatal(err)
	}
	return app, bookmark, collection.ID
}

func TestBookmarkCollectionHTTPAssignmentLifecycleIsOwnerScoped(t *testing.T) {
	app, bookmark, collectionID := newBookmarkCollectionFixture(t)

	request := httptest.NewRequest(http.MethodPut, "/api/v1/bookmarks/"+bookmark.ID+"/collection", strings.NewReader(`{"collectionId":"`+collectionID+`"}`))
	request.SetPathValue("id", bookmark.ID)
	response := httptest.NewRecorder()
	app.assignBookmarkCollection(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("assign status=%d body=%s", response.Code, response.Body.String())
	}
	for _, expected := range []string{bookmark.ID, collectionID, `"ownerId":"owner-a"`} {
		if !strings.Contains(response.Body.String(), expected) {
			t.Fatalf("assignment response missing %q: %s", expected, response.Body.String())
		}
	}

	request = httptest.NewRequest(http.MethodGet, "/api/v1/bookmarks/"+bookmark.ID+"/collection", nil)
	request.SetPathValue("id", bookmark.ID)
	response = httptest.NewRecorder()
	app.getBookmarkCollection(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), collectionID) {
		t.Fatalf("get assignment status=%d body=%s", response.Code, response.Body.String())
	}

	other := app
	other.identity = fixedIdentity{owner: "owner-b"}
	request = httptest.NewRequest(http.MethodGet, "/api/v1/bookmarks/"+bookmark.ID+"/collection", nil)
	request.SetPathValue("id", bookmark.ID)
	response = httptest.NewRecorder()
	other.getBookmarkCollection(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("cross-owner assignment lookup must be 404, got %d", response.Code)
	}
	if strings.Contains(response.Body.String(), "owner-a") || strings.Contains(response.Body.String(), collectionID) {
		t.Fatalf("cross-owner assignment lookup leaked data: %s", response.Body.String())
	}

	request = httptest.NewRequest(http.MethodDelete, "/api/v1/bookmarks/"+bookmark.ID+"/collection", nil)
	request.SetPathValue("id", bookmark.ID)
	response = httptest.NewRecorder()
	other.removeBookmarkCollection(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("cross-owner removal must be 404, got %d", response.Code)
	}

	request = httptest.NewRequest(http.MethodDelete, "/api/v1/bookmarks/"+bookmark.ID+"/collection", nil)
	request.SetPathValue("id", bookmark.ID)
	response = httptest.NewRecorder()
	app.removeBookmarkCollection(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("owner removal status=%d body=%s", response.Code, response.Body.String())
	}
}

func TestBookmarkDeleteRemovesCollectionAssignment(t *testing.T) {
	app, bookmark, collectionID := newBookmarkCollectionFixture(t)
	if _, err := app.assignments.Assign(t.Context(), "owner-a", bookmark.ID, collectionID); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodDelete, "/api/v1/bookmarks/"+bookmark.ID, nil)
	request.SetPathValue("id", bookmark.ID)
	response := httptest.NewRecorder()
	app.delete(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("bookmark delete status=%d body=%s", response.Code, response.Body.String())
	}
	if _, found := app.assignments.Get("owner-a", bookmark.ID); found {
		t.Fatal("bookmark deletion left a stale collection assignment")
	}
}

func TestBookmarkCollectionAssignmentFailsClosedOnRepositoryError(t *testing.T) {
	backendErr := errors.New("private backend details")
	repository := failingGetRepository{
		MemoryRepository: bookmarkcore.NewMemoryRepository(nil),
		err:              backendErr,
	}
	app, err := newServer(repository, fixedIdentity{owner: "owner-a"}, "memory-development")
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPut, "/api/v1/bookmarks/bookmark/collection", strings.NewReader(`{"collectionId":"collection"}`))
	request.SetPathValue("id", "bookmark")
	response := httptest.NewRecorder()
	app.assignBookmarkCollection(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("repository failure status=%d body=%s", response.Code, response.Body.String())
	}
	if strings.Contains(response.Body.String(), "private backend details") {
		t.Fatalf("repository details leaked to client: %s", response.Body.String())
	}
}

func TestBookmarkCollectionAssignmentRejectsMalformedRequests(t *testing.T) {
	app, bookmark, _ := newBookmarkCollectionFixture(t)
	for _, body := range []string{
		`{"collectionId":"collection","unexpected":true}`,
		`{"collectionId":`,
		`{"collectionId":"one"} {"collectionId":"two"}`,
	} {
		request := httptest.NewRequest(http.MethodPut, "/api/v1/bookmarks/"+bookmark.ID+"/collection", strings.NewReader(body))
		request.SetPathValue("id", bookmark.ID)
		response := httptest.NewRecorder()
		app.assignBookmarkCollection(response, request)
		if response.Code != http.StatusBadRequest {
			t.Fatalf("invalid body %q returned %d", body, response.Code)
		}
		if strings.Contains(response.Body.String(), "unexpected") {
			t.Fatalf("decoder details leaked to client: %s", response.Body.String())
		}
	}
}
