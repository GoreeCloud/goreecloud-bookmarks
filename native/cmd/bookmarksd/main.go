package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	bookmarkcore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
	identitycore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/identity"
)

const (
	maxCreateBodyBytes      = 16 * 1024
	defaultRepositoryMode   = "memory-development"
	postgresRepositoryMode  = "postgres-development"
	repositoryModeEnv       = "GOREECLOUD_BOOKMARKS_STORE"
	postgresDatabaseURLEnv  = "GOREECLOUD_BOOKMARKS_DATABASE_URL"
	identityModeEnv         = "GOREECLOUD_BOOKMARKS_IDENTITY"
	unavailableIdentityMode = "unavailable"
	developmentIdentityMode = "development-header"
)

type server struct {
	bookmarks *bookmarkcore.Service
	identity  identitycore.Resolver
	storeMode string
}

type createBookmarkRequest struct {
	URL   string   `json:"url"`
	Title string   `json:"title"`
	Note  string   `json:"note"`
	Tags  []string `json:"tags"`
}

func newServer(repository bookmarkcore.Repository, identity identitycore.Resolver, storeMode string) (server, error) {
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

func selectRuntimeRepository(getenv func(string) string) (bookmarkcore.Repository, string, error) {
	if getenv == nil {
		return nil, "", errors.New("environment reader is required")
	}
	mode := strings.ToLower(strings.TrimSpace(getenv(repositoryModeEnv)))
	if mode == "" {
		mode = defaultRepositoryMode
	}
	switch mode {
	case defaultRepositoryMode:
		return bookmarkcore.NewMemoryRepository(nil), defaultRepositoryMode, nil
	case postgresRepositoryMode:
		if strings.TrimSpace(getenv(postgresDatabaseURLEnv)) == "" {
			return nil, "", errors.New("postgres development mode requires an explicit database URL")
		}
		return nil, "", errors.New("postgres development runtime driver is not integrated")
	default:
		return nil, "", errors.New("unsupported bookmark repository mode")
	}
}

func selectRuntimeIdentity(getenv func(string) string, storeMode string) (identitycore.Resolver, string, error) {
	if getenv == nil {
		return nil, "", errors.New("environment reader is required")
	}
	mode := strings.ToLower(strings.TrimSpace(getenv(identityModeEnv)))
	if mode == "" {
		mode = unavailableIdentityMode
	}
	switch mode {
	case unavailableIdentityMode:
		return identitycore.UnavailableResolver{}, unavailableIdentityMode, nil
	case developmentIdentityMode:
		if storeMode != defaultRepositoryMode {
			return nil, "", errors.New("development header identity requires memory-development repository mode")
		}
		return identitycore.DevelopmentHeaderResolver{}, developmentIdentityMode, nil
	default:
		return nil, "", errors.New("unsupported bookmark identity mode")
	}
}

func main() {
	repository, storeMode, err := selectRuntimeRepository(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}
	identity, identityMode, err := selectRuntimeIdentity(os.Getenv, storeMode)
	if err != nil {
		log.Fatal(err)
	}
	app, err := newServer(repository, identity, storeMode)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", app.health)
	mux.HandleFunc("GET /readyz", app.ready)
	mux.HandleFunc("GET /api/v1/bookmarks", app.list)
	mux.HandleFunc("GET /api/v1/bookmarks/{id}", app.get)
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
	log.Printf("GoreeCloud Bookmarks native development service listening on %s with identity mode %s", addr, identityMode)
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

func (s server) get(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := s.resolveOwner(w, r)
	if !ok {
		return
	}
	bookmark, found, err := s.bookmarks.Get(r.Context(), ownerID, r.PathValue("id"))
	if err != nil {
		if errors.Is(err, bookmarkcore.ErrOwnerIdentityRequired) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authenticated owner identity is required"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "bookmark storage is unavailable"})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "bookmark not found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"bookmark": bookmark})
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
