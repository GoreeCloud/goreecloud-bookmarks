package main

import (
	"errors"
	"net/http"

	bookmarkcore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
	"github.com/GoreeCloud/goreecloud-bookmarks/native/internal/webui"
)

func (s server) library(w http.ResponseWriter, r *http.Request) {
	ownerID, err := s.identity.Resolve(r)
	if err != nil {
		webui.RenderLibrary(w, http.StatusServiceUnavailable, webui.LibraryPageData{
			StoreMode: s.storeMode,
			Error:     "Authenticated identity integration is not available for this native development runtime.",
		})
		return
	}
	items, err := s.bookmarks.List(r.Context(), ownerID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Bookmark storage is unavailable."
		if errors.Is(err, bookmarkcore.ErrOwnerIdentityRequired) {
			status = http.StatusUnauthorized
			message = "Authenticated owner identity is required."
		}
		webui.RenderLibrary(w, status, webui.LibraryPageData{StoreMode: s.storeMode, Error: message})
		return
	}
	webui.RenderLibrary(w, http.StatusOK, webui.LibraryPageData{Bookmarks: items, StoreMode: s.storeMode})
}

func libraryStyles(w http.ResponseWriter, r *http.Request) {
	webui.Styles(w, r)
}
