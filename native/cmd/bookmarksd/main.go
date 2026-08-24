package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	bookmarkcore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
)

type identityResolver interface {
	Resolve(*http.Request) (string, error)
}

type unavailableIdentity struct{}

func (unavailableIdentity) Resolve(*http.Request) (string, error) {
	return "", errIdentityUnavailable
}

var errIdentityUnavailable = &identityError{"production identity adapter is not integrated"}

type identityError struct{ message string }

func (e *identityError) Error() string { return e.message }

type server struct {
	store    *bookmarkcore.Store
	identity identityResolver
}

func main() {
	app := server{store: bookmarkcore.NewStore(), identity: unavailableIdentity{}}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", app.health)
	mux.HandleFunc("GET /api/v1/bookmarks", app.list)

	addr := os.Getenv("GOREECLOUD_BOOKMARKS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           securityHeaders(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("GoreeCloud Bookmarks native development service listening on %s", addr)
	log.Fatal(httpServer.ListenAndServe())
}

func (s server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"service":                  "goreecloud-bookmarks",
		"implementation":           "native-development-foundation",
		"identity_adapter_ready":   false,
		"persistent_store_ready":   false,
		"production_approved":      false,
	})
}

func (s server) list(w http.ResponseWriter, r *http.Request) {
	ownerID, err := s.identity.Resolve(r)
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "authenticated identity integration is not available"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"bookmarks": s.store.List(ownerID)})
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("encode response: %v", err)
	}
}
