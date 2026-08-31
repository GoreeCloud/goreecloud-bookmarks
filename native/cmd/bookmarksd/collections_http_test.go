package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	bookmarkcore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
	identitycore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/identity"
)

func newCollectionTestHandler(t *testing.T) http.Handler {
	t.Helper()
	app, err := newServer(bookmarkcore.NewMemoryRepository(nil), identitycore.DevelopmentHeaderResolver{}, defaultRepositoryMode)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/collections", app.listCollections)
	mux.HandleFunc("GET /api/v1/collections/{id}", app.getCollection)
	mux.HandleFunc("POST /api/v1/collections", app.createCollection)
	mux.HandleFunc("PATCH /api/v1/collections/{id}", app.moveCollection)
	return mux
}

func collectionRequest(t *testing.T, handler http.Handler, method, path, owner, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if owner != "" {
		req.Header.Set(identitycore.DevelopmentOwnerHeader, owner)
	}
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	return recorder
}

func TestCollectionHTTPRoutesAreOwnerScopedAndMovable(t *testing.T) {
	handler := newCollectionTestHandler(t)

	rootResponse := collectionRequest(t, handler, http.MethodPost, "/api/v1/collections", "owner-a", `{"name":"Research"}`)
	if rootResponse.Code != http.StatusCreated {
		t.Fatalf("create root status = %d, body=%s", rootResponse.Code, rootResponse.Body.String())
	}
	var root struct {
		Collection struct {
			ID string `json:"id"`
		} `json:"collection"`
	}
	if err := json.Unmarshal(rootResponse.Body.Bytes(), &root); err != nil || root.Collection.ID == "" {
		t.Fatalf("decode root collection: %v body=%s", err, rootResponse.Body.String())
	}

	childResponse := collectionRequest(t, handler, http.MethodPost, "/api/v1/collections", "owner-a", `{"name":"Docs","parentId":"`+root.Collection.ID+`"}`)
	if childResponse.Code != http.StatusCreated {
		t.Fatalf("create child status = %d, body=%s", childResponse.Code, childResponse.Body.String())
	}
	var child struct {
		Collection struct {
			ID string `json:"id"`
		} `json:"collection"`
	}
	if err := json.Unmarshal(childResponse.Body.Bytes(), &child); err != nil || child.Collection.ID == "" {
		t.Fatalf("decode child collection: %v body=%s", err, childResponse.Body.String())
	}

	otherOwnerGet := collectionRequest(t, handler, http.MethodGet, "/api/v1/collections/"+child.Collection.ID, "owner-b", "")
	if otherOwnerGet.Code != http.StatusNotFound {
		t.Fatalf("cross-owner get status = %d, body=%s", otherOwnerGet.Code, otherOwnerGet.Body.String())
	}

	moveResponse := collectionRequest(t, handler, http.MethodPatch, "/api/v1/collections/"+child.Collection.ID, "owner-a", `{"parentId":""}`)
	if moveResponse.Code != http.StatusOK {
		t.Fatalf("move status = %d, body=%s", moveResponse.Code, moveResponse.Body.String())
	}

	listResponse := collectionRequest(t, handler, http.MethodGet, "/api/v1/collections", "owner-a", "")
	if listResponse.Code != http.StatusOK {
		t.Fatalf("list status = %d, body=%s", listResponse.Code, listResponse.Body.String())
	}
	var listed struct {
		Collections []json.RawMessage `json:"collections"`
	}
	if err := json.Unmarshal(listResponse.Body.Bytes(), &listed); err != nil || len(listed.Collections) != 2 {
		t.Fatalf("expected two owner-a collections: err=%v body=%s", err, listResponse.Body.String())
	}
}

func TestCollectionHTTPRejectsHierarchyConflictAndMalformedBodies(t *testing.T) {
	handler := newCollectionTestHandler(t)

	rootResponse := collectionRequest(t, handler, http.MethodPost, "/api/v1/collections", "owner-a", `{"name":"Root"}`)
	var root struct {
		Collection struct {
			ID string `json:"id"`
		} `json:"collection"`
	}
	if err := json.Unmarshal(rootResponse.Body.Bytes(), &root); err != nil {
		t.Fatal(err)
	}

	selfMove := collectionRequest(t, handler, http.MethodPatch, "/api/v1/collections/"+root.Collection.ID, "owner-a", `{"parentId":"`+root.Collection.ID+`"}`)
	if selfMove.Code != http.StatusConflict {
		t.Fatalf("self move status = %d, body=%s", selfMove.Code, selfMove.Body.String())
	}

	unknownField := collectionRequest(t, handler, http.MethodPost, "/api/v1/collections", "owner-a", `{"name":"Bad","unexpected":true}`)
	if unknownField.Code != http.StatusBadRequest {
		t.Fatalf("unknown field status = %d, body=%s", unknownField.Code, unknownField.Body.String())
	}

	trailingJSON := collectionRequest(t, handler, http.MethodPost, "/api/v1/collections", "owner-a", `{"name":"Bad"}{"name":"Extra"}`)
	if trailingJSON.Code != http.StatusBadRequest {
		t.Fatalf("trailing JSON status = %d, body=%s", trailingJSON.Code, trailingJSON.Body.String())
	}

	missingIdentity := collectionRequest(t, handler, http.MethodGet, "/api/v1/collections", "", "")
	if missingIdentity.Code != http.StatusServiceUnavailable {
		t.Fatalf("missing identity status = %d, body=%s", missingIdentity.Code, missingIdentity.Body.String())
	}
}
