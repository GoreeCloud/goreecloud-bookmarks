package main

import (
	"encoding/json"
	"net/http"
	"testing"

	bookmarkcore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
	identitycore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/identity"
)

func newAssignmentTestHandler(t *testing.T) http.Handler {
	t.Helper()
	app, err := newServer(bookmarkcore.NewMemoryRepository(nil), identitycore.DevelopmentHeaderResolver{}, defaultRepositoryMode)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/bookmarks", app.create)
	mux.HandleFunc("DELETE /api/v1/bookmarks/{id}", app.delete)
	mux.HandleFunc("POST /api/v1/collections", app.createCollection)
	mux.HandleFunc("GET /api/v1/bookmarks/{id}/collection", app.getBookmarkCollection)
	mux.HandleFunc("PUT /api/v1/bookmarks/{id}/collection", app.assignBookmarkCollection)
	mux.HandleFunc("DELETE /api/v1/bookmarks/{id}/collection", app.removeBookmarkCollection)
	return mux
}

func createAssignmentTestBookmark(t *testing.T, handler http.Handler, owner string) string {
	t.Helper()
	response := collectionRequest(t, handler, http.MethodPost, "/api/v1/bookmarks", owner, `{"url":"https://example.com/","title":"Example"}`)
	if response.Code != http.StatusCreated {
		t.Fatalf("create bookmark status = %d, body=%s", response.Code, response.Body.String())
	}
	var decoded struct {
		Bookmark struct {
			ID string `json:"id"`
		} `json:"bookmark"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &decoded); err != nil || decoded.Bookmark.ID == "" {
		t.Fatalf("decode bookmark: %v body=%s", err, response.Body.String())
	}
	return decoded.Bookmark.ID
}

func createAssignmentTestCollection(t *testing.T, handler http.Handler, owner, name string) string {
	t.Helper()
	response := collectionRequest(t, handler, http.MethodPost, "/api/v1/collections", owner, `{"name":"`+name+`"}`)
	if response.Code != http.StatusCreated {
		t.Fatalf("create collection status = %d, body=%s", response.Code, response.Body.String())
	}
	var decoded struct {
		Collection struct {
			ID string `json:"id"`
		} `json:"collection"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &decoded); err != nil || decoded.Collection.ID == "" {
		t.Fatalf("decode collection: %v body=%s", err, response.Body.String())
	}
	return decoded.Collection.ID
}

func TestBookmarkCollectionAssignmentHTTPIsOwnerScopedAndReplaceable(t *testing.T) {
	handler := newAssignmentTestHandler(t)
	bookmarkID := createAssignmentTestBookmark(t, handler, "owner-a")
	firstCollection := createAssignmentTestCollection(t, handler, "owner-a", "Research")
	secondCollection := createAssignmentTestCollection(t, handler, "owner-a", "Reference")
	otherCollection := createAssignmentTestCollection(t, handler, "owner-b", "Private")

	assign := collectionRequest(t, handler, http.MethodPut, "/api/v1/bookmarks/"+bookmarkID+"/collection", "owner-a", `{"collectionId":"`+firstCollection+`"}`)
	if assign.Code != http.StatusOK {
		t.Fatalf("assign status = %d, body=%s", assign.Code, assign.Body.String())
	}

	get := collectionRequest(t, handler, http.MethodGet, "/api/v1/bookmarks/"+bookmarkID+"/collection", "owner-a", "")
	if get.Code != http.StatusOK {
		t.Fatalf("get assignment status = %d, body=%s", get.Code, get.Body.String())
	}
	var decoded struct {
		Assignment struct {
			CollectionID string `json:"collectionId"`
		} `json:"assignment"`
	}
	if err := json.Unmarshal(get.Body.Bytes(), &decoded); err != nil || decoded.Assignment.CollectionID != firstCollection {
		t.Fatalf("unexpected assignment: err=%v body=%s", err, get.Body.String())
	}

	reassign := collectionRequest(t, handler, http.MethodPut, "/api/v1/bookmarks/"+bookmarkID+"/collection", "owner-a", `{"collectionId":"`+secondCollection+`"}`)
	if reassign.Code != http.StatusOK {
		t.Fatalf("reassign status = %d, body=%s", reassign.Code, reassign.Body.String())
	}

	crossOwnerBookmark := collectionRequest(t, handler, http.MethodGet, "/api/v1/bookmarks/"+bookmarkID+"/collection", "owner-b", "")
	if crossOwnerBookmark.Code != http.StatusNotFound {
		t.Fatalf("cross-owner assignment read status = %d, body=%s", crossOwnerBookmark.Code, crossOwnerBookmark.Body.String())
	}
	crossOwnerCollection := collectionRequest(t, handler, http.MethodPut, "/api/v1/bookmarks/"+bookmarkID+"/collection", "owner-a", `{"collectionId":"`+otherCollection+`"}`)
	if crossOwnerCollection.Code != http.StatusNotFound {
		t.Fatalf("cross-owner collection assignment status = %d, body=%s", crossOwnerCollection.Code, crossOwnerCollection.Body.String())
	}

	remove := collectionRequest(t, handler, http.MethodDelete, "/api/v1/bookmarks/"+bookmarkID+"/collection", "owner-a", "")
	if remove.Code != http.StatusNoContent {
		t.Fatalf("remove assignment status = %d, body=%s", remove.Code, remove.Body.String())
	}
	missing := collectionRequest(t, handler, http.MethodGet, "/api/v1/bookmarks/"+bookmarkID+"/collection", "owner-a", "")
	if missing.Code != http.StatusNotFound {
		t.Fatalf("removed assignment status = %d, body=%s", missing.Code, missing.Body.String())
	}
}

func TestBookmarkDeleteRemovesAssignmentAndAssignmentJSONIsStrict(t *testing.T) {
	handler := newAssignmentTestHandler(t)
	bookmarkID := createAssignmentTestBookmark(t, handler, "owner-a")
	collectionID := createAssignmentTestCollection(t, handler, "owner-a", "Research")

	unknownField := collectionRequest(t, handler, http.MethodPut, "/api/v1/bookmarks/"+bookmarkID+"/collection", "owner-a", `{"collectionId":"`+collectionID+`","unexpected":true}`)
	if unknownField.Code != http.StatusBadRequest {
		t.Fatalf("unknown field status = %d, body=%s", unknownField.Code, unknownField.Body.String())
	}
	trailing := collectionRequest(t, handler, http.MethodPut, "/api/v1/bookmarks/"+bookmarkID+"/collection", "owner-a", `{"collectionId":"`+collectionID+`"}{}`)
	if trailing.Code != http.StatusBadRequest {
		t.Fatalf("trailing JSON status = %d, body=%s", trailing.Code, trailing.Body.String())
	}

	assign := collectionRequest(t, handler, http.MethodPut, "/api/v1/bookmarks/"+bookmarkID+"/collection", "owner-a", `{"collectionId":"`+collectionID+`"}`)
	if assign.Code != http.StatusOK {
		t.Fatalf("assign status = %d, body=%s", assign.Code, assign.Body.String())
	}
	deleted := collectionRequest(t, handler, http.MethodDelete, "/api/v1/bookmarks/"+bookmarkID, "owner-a", "")
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("delete bookmark status = %d, body=%s", deleted.Code, deleted.Body.String())
	}
	assignment := collectionRequest(t, handler, http.MethodGet, "/api/v1/bookmarks/"+bookmarkID+"/collection", "owner-a", "")
	if assignment.Code != http.StatusNotFound {
		t.Fatalf("assignment must be removed with bookmark: status=%d body=%s", assignment.Code, assignment.Body.String())
	}
}
