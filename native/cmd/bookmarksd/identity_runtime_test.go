package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	identitycore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/identity"
)

func TestRuntimeIdentityDefaultsUnavailable(t *testing.T) {
	resolver, mode, err := selectRuntimeIdentity(func(string) string { return "" }, defaultRepositoryMode)
	if err != nil {
		t.Fatal(err)
	}
	if mode != unavailableIdentityMode {
		t.Fatalf("unexpected identity mode: %q", mode)
	}
	if _, ok := resolver.(identitycore.UnavailableResolver); !ok {
		t.Fatalf("default identity must remain unavailable, got %T", resolver)
	}
}

func TestRuntimeIdentityRequiresExplicitDevelopmentMode(t *testing.T) {
	resolver, mode, err := selectRuntimeIdentity(func(key string) string {
		if key == identityModeEnv {
			return developmentIdentityMode
		}
		return ""
	}, defaultRepositoryMode)
	if err != nil {
		t.Fatal(err)
	}
	if mode != developmentIdentityMode {
		t.Fatalf("unexpected identity mode: %q", mode)
	}
	if _, ok := resolver.(identitycore.DevelopmentHeaderResolver); !ok {
		t.Fatalf("expected development header resolver, got %T", resolver)
	}
}

func TestRuntimeIdentityRejectsDevelopmentModeOutsideMemoryRepository(t *testing.T) {
	_, _, err := selectRuntimeIdentity(func(key string) string {
		if key == identityModeEnv {
			return developmentIdentityMode
		}
		return ""
	}, postgresRepositoryMode)
	if err == nil || !strings.Contains(err.Error(), "requires memory-development repository mode") {
		t.Fatalf("expected repository-bound development identity failure, got %v", err)
	}
}

func TestRuntimeIdentityRejectsUnknownMode(t *testing.T) {
	_, _, err := selectRuntimeIdentity(func(key string) string {
		if key == identityModeEnv {
			return "trusted-header"
		}
		return ""
	}, defaultRepositoryMode)
	if err == nil || !strings.Contains(err.Error(), "unsupported bookmark identity mode") {
		t.Fatalf("expected unsupported identity mode failure, got %v", err)
	}
}

func TestDevelopmentRuntimeIdentityExercisesOwnerScopedAPI(t *testing.T) {
	resolver, _, err := selectRuntimeIdentity(func(key string) string {
		if key == identityModeEnv {
			return developmentIdentityMode
		}
		return ""
	}, defaultRepositoryMode)
	if err != nil {
		t.Fatal(err)
	}
	repository, _, err := selectRuntimeRepository(func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	app, err := newServer(repository, resolver, defaultRepositoryMode)
	if err != nil {
		t.Fatal(err)
	}

	create := httptest.NewRequest(http.MethodPost, "/api/v1/bookmarks", strings.NewReader(`{"url":"https://example.com/native"}`))
	create.Header.Set(identitycore.DevelopmentOwnerHeader, "owner-dev")
	created := httptest.NewRecorder()
	app.create(created, create)
	if created.Code != http.StatusCreated || !strings.Contains(created.Body.String(), `"ownerId":"owner-dev"`) {
		t.Fatalf("development identity create failed: status=%d body=%s", created.Code, created.Body.String())
	}

	list := httptest.NewRequest(http.MethodGet, "/api/v1/bookmarks", nil)
	list.Header.Set(identitycore.DevelopmentOwnerHeader, "owner-dev")
	listed := httptest.NewRecorder()
	app.list(listed, list)
	if listed.Code != http.StatusOK || !strings.Contains(listed.Body.String(), "https://example.com/native") {
		t.Fatalf("development identity list failed: status=%d body=%s", listed.Code, listed.Body.String())
	}

	other := httptest.NewRequest(http.MethodGet, "/api/v1/bookmarks", nil)
	other.Header.Set(identitycore.DevelopmentOwnerHeader, "other-owner")
	otherList := httptest.NewRecorder()
	app.list(otherList, other)
	if otherList.Code != http.StatusOK || strings.Contains(otherList.Body.String(), "https://example.com/native") {
		t.Fatalf("development owner isolation failed: status=%d body=%s", otherList.Code, otherList.Body.String())
	}
}
