package httpapi

import (
	"encoding/json"
	"net/http"
)

// NewHandler returns the current experimental GoreeCloud Bookmarks HTTP surface.
//
// The first executable foundation intentionally exposes only bounded health and
// readiness endpoints. Bookmark resources, persistence, authentication, and
// other planned /api/v1 operations remain unimplemented.
func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", handleHealth)
	mux.HandleFunc("GET /api/v1/ready", handleReadiness)
	return responseHeaders(mux)
}

type healthResponse struct {
	Status string `json:"status"`
}

type readinessResponse struct {
	Ready  bool              `json:"ready"`
	Checks map[string]string `json:"checks"`
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{Status: "ok"})
}

func handleReadiness(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusServiceUnavailable, readinessResponse{
		Ready: false,
		Checks: map[string]string{
			"bookmarks-data": "not-configured",
		},
	})
}

func responseHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
