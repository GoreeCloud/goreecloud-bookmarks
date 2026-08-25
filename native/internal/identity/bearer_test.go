package identity

import (
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
)

type validatorFunc func(context.Context, string) (AuthenticatedIdentity, error)

func (fn validatorFunc) Validate(ctx context.Context, token string) (AuthenticatedIdentity, error) {
	return fn(ctx, token)
}

func TestBearerResolverRequiresCredential(t *testing.T) {
	resolver := BearerResolver{Validator: validatorFunc(func(context.Context, string) (AuthenticatedIdentity, error) {
		t.Fatal("validator must not run without a bearer credential")
		return AuthenticatedIdentity{}, nil
	})}

	request := httptest.NewRequest("GET", "/api/v1/bookmarks", nil)
	if _, err := resolver.Resolve(request); !errors.Is(err, ErrCredentialRequired) {
		t.Fatalf("expected credential-required error, got %v", err)
	}
}

func TestBearerResolverRequiresValidator(t *testing.T) {
	request := httptest.NewRequest("GET", "/api/v1/bookmarks", nil)
	request.Header.Set("Authorization", "Bearer opaque-token")

	if _, err := (BearerResolver{}).Resolve(request); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("expected unavailable identity error, got %v", err)
	}
}

func TestBearerResolverUsesValidatedSubject(t *testing.T) {
	var observedToken string
	resolver := BearerResolver{Validator: validatorFunc(func(_ context.Context, token string) (AuthenticatedIdentity, error) {
		observedToken = token
		return AuthenticatedIdentity{Subject: "identity-user-123"}, nil
	})}
	request := httptest.NewRequest("GET", "/api/v1/bookmarks", nil)
	request.Header.Set("Authorization", "Bearer opaque-token")

	ownerID, err := resolver.Resolve(request)
	if err != nil {
		t.Fatalf("resolve authenticated owner: %v", err)
	}
	if ownerID != "identity-user-123" {
		t.Fatalf("unexpected owner identifier %q", ownerID)
	}
	if observedToken != "opaque-token" {
		t.Fatalf("validator received unexpected token %q", observedToken)
	}
}

func TestBearerResolverNormalizesValidatorFailure(t *testing.T) {
	resolver := BearerResolver{Validator: validatorFunc(func(context.Context, string) (AuthenticatedIdentity, error) {
		return AuthenticatedIdentity{}, errors.New("sensitive upstream validator detail")
	})}
	request := httptest.NewRequest("GET", "/api/v1/bookmarks", nil)
	request.Header.Set("Authorization", "Bearer rejected-token")

	_, err := resolver.Resolve(request)
	if !errors.Is(err, ErrAuthenticationFailed) {
		t.Fatalf("expected bounded authentication error, got %v", err)
	}
	if strings.Contains(err.Error(), "upstream") {
		t.Fatalf("authentication error leaked validator detail: %v", err)
	}
}

func TestBearerResolverRejectsInvalidSubject(t *testing.T) {
	resolver := BearerResolver{Validator: validatorFunc(func(context.Context, string) (AuthenticatedIdentity, error) {
		return AuthenticatedIdentity{Subject: "invalid subject"}, nil
	})}
	request := httptest.NewRequest("GET", "/api/v1/bookmarks", nil)
	request.Header.Set("Authorization", "Bearer opaque-token")

	if _, err := resolver.Resolve(request); !errors.Is(err, ErrAuthenticationFailed) {
		t.Fatalf("expected invalid authenticated subject to fail closed, got %v", err)
	}
}

func TestBearerResolverRejectsOversizedCredential(t *testing.T) {
	resolver := BearerResolver{Validator: validatorFunc(func(context.Context, string) (AuthenticatedIdentity, error) {
		t.Fatal("validator must not run for oversized credential")
		return AuthenticatedIdentity{}, nil
	})}
	request := httptest.NewRequest("GET", "/api/v1/bookmarks", nil)
	request.Header.Set("Authorization", "Bearer "+strings.Repeat("x", maxBearerTokenBytes+1))

	if _, err := resolver.Resolve(request); !errors.Is(err, ErrAuthenticationFailed) {
		t.Fatalf("expected oversized credential failure, got %v", err)
	}
}
