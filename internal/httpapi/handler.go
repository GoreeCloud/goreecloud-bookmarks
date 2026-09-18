package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
)

// ReadinessResult is the bounded service state exposed by the readiness route.
type ReadinessResult struct {
	Ready  bool
	Checks map[string]string
}

// ReadinessChecker resolves current dependency/schema readiness without
// exposing protected configuration or raw database errors.
type ReadinessChecker interface {
	CheckReadiness(context.Context) ReadinessResult
}

// ReadinessFunc adapts a function to ReadinessChecker.
type ReadinessFunc func(context.Context) ReadinessResult

// CheckReadiness implements ReadinessChecker.
func (f ReadinessFunc) CheckReadiness(ctx context.Context) ReadinessResult {
	return f(ctx)
}

// NewHandler returns the current Experimental HTTP surface with fail-closed
// default readiness for callers that have not configured a data layer.
func NewHandler() http.Handler {
	return NewHandlerWithReadiness(ReadinessFunc(func(context.Context) ReadinessResult {
		return ReadinessResult{
			Ready: false,
			Checks: map[string]string{
				"bookmarks-data": "not-configured",
			},
		}
	}))
}

// NewHandlerWithReadiness binds the implemented health/readiness routes to the
// supplied dependency checker. Bookmark-domain HTTP operations remain
// unimplemented and are not inferred from the data layer.
func NewHandlerWithReadiness(readiness ReadinessChecker) http.Handler {
	if readiness == nil {
		return NewHandler()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", handleHealth)
	mux.HandleFunc("GET /api/v1/ready", func(w http.ResponseWriter, r *http.Request) {
		handleReadiness(w, r, readiness)
	})
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

func handleReadiness(w http.ResponseWriter, r *http.Request, readiness ReadinessChecker) {
	result := readiness.CheckReadiness(r.Context())
	status := http.StatusServiceUnavailable
	if result.Ready {
		status = http.StatusOK
	}
	writeJSON(w, status, readinessResponse{
		Ready:  result.Ready,
		Checks: result.Checks,
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
