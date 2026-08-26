package main

import (
	"net/http"

	"github.com/GoreeCloud/goreecloud-bookmarks/native/internal/syncstate"
)

// syncCapabilities exposes protocol metadata only. It deliberately does not
// expose bookmark records, owner identifiers, storage locations or credentials.
func (s server) syncCapabilities(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"schema_version":      1,
		"application":         "bookmarks",
		"capabilities":        syncstate.Capabilities(),
		"records_exposed":     false,
		"owner_ids_exposed":   false,
		"credentials_exposed": false,
		"production_approved": false,
	})
}
