package identity

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"unicode"
)

const (
	maxBearerTokenBytes = 8192
	maxSubjectRunes     = 256
)

var (
	ErrCredentialRequired   = errors.New("authenticated credential is required")
	ErrAuthenticationFailed = errors.New("authenticated credential is invalid")
)

// AuthenticatedIdentity is the minimum application-facing identity produced by
// a reviewed authentication validator. Subject must be a stable GoreeCloud
// identity identifier; application authorization remains owned by Bookmarks.
type AuthenticatedIdentity struct {
	Subject string
}

// TokenValidator validates a bearer credential against the approved identity
// authority. Production implementations are responsible for OIDC/OAuth token
// signature, issuer, audience, lifetime, and other protocol validation required
// by the approved GoreeCloud Identity integration.
type TokenValidator interface {
	Validate(context.Context, string) (AuthenticatedIdentity, error)
}

// BearerResolver converts an already validated bearer credential into the
// owner identifier consumed by native Bookmarks. It never trusts bearer token
// contents directly and normalizes validator failures to a bounded error.
type BearerResolver struct {
	Validator TokenValidator
}

func (resolver BearerResolver) Resolve(r *http.Request) (string, error) {
	if r == nil {
		return "", ErrCredentialRequired
	}
	if resolver.Validator == nil {
		return "", ErrUnavailable
	}

	parts := strings.Fields(r.Header.Get("Authorization"))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", ErrCredentialRequired
	}
	credential := parts[1]
	if credential == "" || len(credential) > maxBearerTokenBytes {
		return "", ErrAuthenticationFailed
	}

	authenticated, err := resolver.Validator.Validate(r.Context(), credential)
	if err != nil {
		return "", ErrAuthenticationFailed
	}
	if !validSubject(authenticated.Subject) {
		return "", ErrAuthenticationFailed
	}
	return authenticated.Subject, nil
}

func validSubject(subject string) bool {
	if subject == "" || strings.TrimSpace(subject) != subject || len([]rune(subject)) > maxSubjectRunes {
		return false
	}
	for _, value := range subject {
		if unicode.IsControl(value) || unicode.IsSpace(value) {
			return false
		}
	}
	return true
}

var _ Resolver = BearerResolver{}
