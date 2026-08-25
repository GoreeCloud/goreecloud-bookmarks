package identity

import (
	"errors"
	"net/http/httptest"
	"testing"
)

func TestUnavailableResolverFailsClosed(t *testing.T) {
	ownerID, err := (UnavailableResolver{}).Resolve(httptest.NewRequest("GET", "/", nil))
	if ownerID != "" {
		t.Fatalf("ownerID = %q, want empty", ownerID)
	}
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("error = %v, want ErrUnavailable", err)
	}
}

func TestDevelopmentResolverSatisfiesSharedContract(t *testing.T) {
	var resolver Resolver = DevelopmentHeaderResolver{}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set(DevelopmentOwnerHeader, "owner-123")
	ownerID, err := resolver.Resolve(r)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if ownerID != "owner-123" {
		t.Fatalf("ownerID = %q, want owner-123", ownerID)
	}
}
