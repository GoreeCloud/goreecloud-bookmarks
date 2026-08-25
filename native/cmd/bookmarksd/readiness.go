package main

import "net/http"

func (s server) ready(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusServiceUnavailable, map[string]any{
		"service":                "goreecloud-bookmarks",
		"ready":                  false,
		"identity_adapter_ready": false,
		"persistent_store_ready": false,
		"store_mode":             s.storeMode,
		"production_approved":    false,
	})
}
