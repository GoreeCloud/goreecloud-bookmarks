package main

import (
	"errors"
	"net/http"

	bookmarkcore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
	collectioncore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/collections"
)

type assignBookmarkCollectionRequest struct {
	CollectionID string `json:"collectionId"`
}

func (s server) getBookmarkCollection(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := s.resolveOwner(w, r)
	if !ok {
		return
	}
	bookmarkID := r.PathValue("id")
	found, err := s.bookmarks.Exists(r.Context(), ownerID, bookmarkID)
	if err != nil {
		if errors.Is(err, bookmarkcore.ErrOwnerIdentityRequired) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authenticated owner identity is required"})
			return
		}
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "bookmark storage is unavailable"})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "bookmark collection assignment not found"})
		return
	}
	if s.assignments == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "bookmark collection storage is unavailable"})
		return
	}
	assignment, found := s.assignments.Get(ownerID, bookmarkID)
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "bookmark collection assignment not found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"assignment": assignment})
}

func (s server) assignBookmarkCollection(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := s.resolveOwner(w, r)
	if !ok {
		return
	}
	var input assignBookmarkCollectionRequest
	if err := decodeRequestJSON(w, r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bookmark collection request is invalid"})
		return
	}
	if s.assignments == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "bookmark collection storage is unavailable"})
		return
	}
	assignment, err := s.assignments.Assign(r.Context(), ownerID, r.PathValue("id"), input.CollectionID)
	if err != nil {
		switch {
		case errors.Is(err, bookmarkcore.ErrOwnerIdentityRequired), errors.Is(err, collectioncore.ErrOwnerRequired):
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authenticated owner identity is required"})
		case errors.Is(err, collectioncore.ErrBookmarkRequired):
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bookmark collection request is invalid"})
		case errors.Is(err, collectioncore.ErrBookmarkNotFound), errors.Is(err, collectioncore.ErrCollectionNotFound):
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "bookmark or collection not found"})
		default:
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "bookmark storage is unavailable"})
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"assignment": assignment})
}

func (s server) removeBookmarkCollection(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := s.resolveOwner(w, r)
	if !ok {
		return
	}
	bookmarkID := r.PathValue("id")
	found, err := s.bookmarks.Exists(r.Context(), ownerID, bookmarkID)
	if err != nil {
		if errors.Is(err, bookmarkcore.ErrOwnerIdentityRequired) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authenticated owner identity is required"})
			return
		}
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "bookmark storage is unavailable"})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "bookmark collection assignment not found"})
		return
	}
	if s.assignments == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "bookmark collection storage is unavailable"})
		return
	}
	if !s.assignments.Remove(ownerID, bookmarkID) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "bookmark collection assignment not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
