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
	collectioncore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/collections"
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
	bookmarks   *bookmarkcore.Service
	collections *collectioncore.Store
	assignments *collectioncore.AssignmentStore
	identity    identitycore.Resolver
	storeMode   string
}

type createBookmarkRequest struct {
	URL   string   `json:"url"`
	Title string   `json:"title"`
	Note  string   `json:"note"`
	Tags  []string `json:"tags"`
}

type updateBookmarkRequest struct {
	URL   string   `json:"url"`
	Title string   `json:"title"`
	Note  string   `json:"note"`
	Tags  []string `json:"tags"`
}

type createCollectionRequest struct {
	Name     string `json:"name"`
	ParentID string `json:"parentId"`
}

type moveCollectionRequest struct {
	ParentID string `json:"parentId"`
}

type assignBookmarkCollectionRequest struct {
	CollectionID string `json:"collectionId"`
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
	collections := collectioncore.NewStore()
	var assignments *collectioncore.AssignmentStore
	if lookup, ok := repository.(collectioncore.BookmarkLookup); ok {
		assignments = collectioncore.NewAssignmentStore(collections, lookup)
	}
	return server{
		bookmarks:   service,
		collections: collections,
		assignments: assignments,
		identity:    identity,
		storeMode:   storeMode,
	}, nil
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
	mux.HandleFunc("GET /", app.library)
	mux.HandleFunc("GET /assets/library.css", libraryStyles)
	mux.HandleFunc("GET /healthz", app.health)
	mux.HandleFunc("GET /readyz", app.ready)
	mux.HandleFunc("GET /api/v1/bookmarks", app.list)
	mux.HandleFunc("GET /api/v1/bookmarks/{id}", app.get)
	mux.HandleFunc("POST /api/v1/bookmarks", app.create)
	mux.HandleFunc("PATCH /api/v1/bookmarks/{id}", app.update)
	mux.HandleFunc("DELETE /api/v1/bookmarks/{id}", app.delete)
	mux.HandleFunc("GET /api/v1/bookmarks/{id}/collection", app.getBookmarkCollection)
	mux.HandleFunc("PUT /api/v1/bookmarks/{id}/collection", app.assignBookmarkCollection)
	mux.HandleFunc("DELETE /api/v1/bookmarks/{id}/collection", app.removeBookmarkCollection)
	mux.HandleFunc("GET /api/v1/collections", app.listCollections)
	mux.HandleFunc("GET /api/v1/collections/{id}", app.getCollection)
	mux.HandleFunc("POST /api/v1/collections", app.createCollection)
	mux.HandleFunc("PATCH /api/v1/collections/{id}", app.moveCollection)

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
		"service":                     "goreecloud-bookmarks",
		"implementation":              "native-development-foundation",
		"identity_adapter_ready":      false,
		"persistent_store_ready":      false,
		"collection_assignment_ready": s.assignments != nil,
		"store_mode":                  s.storeMode,
		"production_approved":         false,
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
	if err := decodeRequestJSON(w, r, &input); err != nil {
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

func (s server) update(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := s.resolveOwner(w, r)
	if !ok {
		return
	}

	var input updateBookmarkRequest
	if err := decodeRequestJSON(w, r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bookmark request is invalid"})
		return
	}

	bookmark, found, err := s.bookmarks.Update(r.Context(), ownerID, r.PathValue("id"), bookmarkcore.UpdateInput{
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
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "bookmark not found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"bookmark": bookmark})
}

func (s server) delete(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := s.resolveOwner(w, r)
	if !ok {
		return
	}

	bookmarkID := r.PathValue("id")
	deleted, err := s.bookmarks.Delete(r.Context(), ownerID, bookmarkID)
	if err != nil {
		if errors.Is(err, bookmarkcore.ErrOwnerIdentityRequired) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authenticated owner identity is required"})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bookmark request is invalid"})
		return
	}
	if !deleted {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "bookmark not found"})
		return
	}
	if s.assignments != nil {
		s.assignments.Remove(ownerID, bookmarkID)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s server) getBookmarkCollection(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := s.resolveOwner(w, r)
	if !ok {
		return
	}
	if s.assignments == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "collection assignment storage is unavailable"})
		return
	}
	assignment, found := s.assignments.Get(ownerID, r.PathValue("id"))
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
	if s.assignments == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "collection assignment storage is unavailable"})
		return
	}
	var input assignBookmarkCollectionRequest
	if err := decodeRequestJSON(w, r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "collection assignment request is invalid"})
		return
	}
	assignment, err := s.assignments.Assign(ownerID, r.PathValue("id"), input.CollectionID)
	if err != nil {
		switch {
		case errors.Is(err, collectioncore.ErrBookmarkNotFound), errors.Is(err, collectioncore.ErrCollectionNotFound):
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "bookmark or collection not found"})
		default:
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "collection assignment is invalid"})
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
	if s.assignments == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "collection assignment storage is unavailable"})
		return
	}
	if !s.assignments.Remove(ownerID, r.PathValue("id")) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "bookmark collection assignment not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s server) listCollections(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := s.resolveOwner(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"collections": s.collections.List(ownerID)})
}

func (s server) getCollection(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := s.resolveOwner(w, r)
	if !ok {
		return
	}
	collection, found := s.collections.Get(ownerID, r.PathValue("id"))
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "collection not found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"collection": collection})
}

func (s server) createCollection(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := s.resolveOwner(w, r)
	if !ok {
		return
	}
	var input createCollectionRequest
	if err := decodeRequestJSON(w, r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "collection request is invalid"})
		return
	}
	collection, err := s.collections.Create(ownerID, input.Name, input.ParentID)
	if err != nil {
		switch {
		case errors.Is(err, collectioncore.ErrDuplicateSibling):
			writeJSON(w, http.StatusConflict, map[string]string{"error": "collection name already exists under parent"})
		case errors.Is(err, collectioncore.ErrParentNotFound):
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "collection parent not found"})
		default:
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "collection input is invalid"})
		}
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"collection": collection})
}

func (s server) moveCollection(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := s.resolveOwner(w, r)
	if !ok {
		return
	}
	var input moveCollectionRequest
	if err := decodeRequestJSON(w, r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "collection request is invalid"})
		return
	}
	collection, err := s.collections.Move(ownerID, r.PathValue("id"), input.ParentID)
	if err != nil {
		switch {
		case errors.Is(err, collectioncore.ErrCollectionNotFound), errors.Is(err, collectioncore.ErrParentNotFound):
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "collection or parent not found"})
		case errors.Is(err, collectioncore.ErrDuplicateSibling), errors.Is(err, collectioncore.ErrHierarchyCycle):
			writeJSON(w, http.StatusConflict, map[string]string{"error": "collection move conflicts with hierarchy"})
		default:
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "collection move is invalid"})
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"collection": collection})
}

func (s server) resolveOwner(w http.ResponseWriter, r *http.Request) (string, bool) {
	ownerID, err := s.identity.Resolve(r)
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "authenticated identity integration is not available"})
		return "", false
	}
	return ownerID, true
}

func decodeRequestJSON(w http.ResponseWriter, r *http.Request, target any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxCreateBodyBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	return ensureJSONEOF(decoder)
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
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self'; img-src 'self' data:; form-action 'self'; frame-ancestors 'none'; base-uri 'none'")
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("encode response: %v", err)
	}
}
