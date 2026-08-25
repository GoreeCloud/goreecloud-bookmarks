package identity

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDevelopmentHeaderResolverRequiresExplicitOwner(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	ownerID, err := (DevelopmentHeaderResolver{}).Resolve(request)
	if !errors.Is(err, ErrOwnerRequired) || ownerID != "" {
		t.Fatalf("expected missing owner failure, got owner=%q err=%v", ownerID, err)
	}
}

func TestDevelopmentHeaderResolverReturnsBoundedOwner(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set(DevelopmentOwnerHeader, "  owner-a  ")
	ownerID, err := (DevelopmentHeaderResolver{}).Resolve(request)
	if err != nil {
		t.Fatal(err)
	}
	if ownerID != "owner-a" {
		t.Fatalf("unexpected owner identity: %q", ownerID)
	}
}

func TestDevelopmentHeaderResolverRejectsWhitespaceAndControls(t *testing.T) {
	for _, ownerID := range []string{"owner a", "owner\ta", "owner\na"} {
		request := httptest.NewRequest("GET", "/", nil)
		request.Header.Set(DevelopmentOwnerHeader, ownerID)
		if _, err := (DevelopmentHeaderResolver{}).Resolve(request); !errors.Is(err, ErrOwnerInvalid) {
			t.Fatalf("expected invalid owner failure for %q, got %v", ownerID, err)
		}
	}
}

func TestDevelopmentHeaderResolverRejectsOversizedOwner(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set(DevelopmentOwnerHeader, strings.Repeat("a", maxDevelopmentOwnerRunes+1))
	if _, err := (DevelopmentHeaderResolver{}).Resolve(request); !errors.Is(err, ErrOwnerInvalid) {
		t.Fatalf("expected oversized owner failure, got %v", err)
	}
}

func TestDevelopmentHeaderResolverRejectsNilRequest(t *testing.T) {
	if _, err := (DevelopmentHeaderResolver{}).Resolve(nil); !errors.Is(err, ErrOwnerRequired) {
		t.Fatalf("expected nil request failure, got %v", err)
	}
}
