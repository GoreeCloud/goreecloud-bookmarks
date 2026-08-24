package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	bookmarkcore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
)

const maxCreateBodyBytes = 16 * 1024

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
	bookmarks *bookmarkcore.Service
	identity  identityResolver
	storeMode string
}

type createBookmarkRequest struct {
	URL   string   `json:"url"`
	Title string   `json:"title"`
	Note  string   `json:"note"`
	Tags  []string `json:"tags"`
}

func newServer(repository bookmarkcore.Repository, identity identityResolver, storeMode string) (server, error) {
	if identity == nil {
		return server{}, errors.New("identity resolver is required")
	}
	service, err := bookmarkcore.NewService(repository)
	if err != nil {
		return server{}, err
	}
	if storeMode == "" {
		storeMode = "unspecified"
	}
	return server{bookmarks: service, identity: identity, storeMode: storeMode}, nil
}

func main() {
	app, err := newServer(bookmarkcore.NewMemoryRepository(nil), unavailableIdentity{}, "memory-development")
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", app.health)
	mux.HandleFunc("GET /api/v1/bookmarks", app.list)
	mux.HandleFunc("POST /api/v1/bookmarks", app.create)

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
		"service":                "goreecloud-bookmarks",
		"implementation":         "native-development-foundation",
		"identity_adapter_ready": false,
		"persistent_store_ready": false,
		"store_mode":             s.storeMode,
		"production_approved":    false,
	})
}

func (s server) list(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := s.resolveOwner(w, r)
	if !ok {
		return
	}
	items, err := s.bookmarks.List(r.Context(), ownerID)
	if err != nil {
		if errors.Is(err, bookmarkcore.ErrOwnerIdentityRequired) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authenticated owner identity is required"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "bookmark storage is unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"bookmarks": items})
}

func (s server) create(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := s.resolveOwner(w, r)
	if !ok {
		return
	}

	var input createBookmarkRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxCreateBodyBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bookmark request is invalid"})
		return
	}
	if err := ensureJSONEOF(decoder); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bookmark request is invalid"})
		return
	}

	bookmark, err := s.bookmarks.Create(r.Context(), ownerID, bookmarkcore.CreateInput{
		URL: input.URL, Title: input.Title, Note: input.Note, Tags: input.Tags,
	})
	if err != nil {
		if errors.Is(err, bookmarkcore.ErrOwnerIdentityRequired) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authenticated owner identity is required"})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bookmark input is invalid"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"bookmark": bookmark})
}

func (s server) resolveOwner(w http.ResponseWriter, r *http.Request) (string, bool) {
	ownerID, err := s.identity.Resolve(r)
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "authenticated identity integration is not available"})
		return "", false
	}
	return ownerID, true
}

func ensureJSONEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); errors.Is(err, io.EOF) {
		return nil
	} else if err != nil {
		return err
	}
	return errors.New("multiple JSON values are not allowed")
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
