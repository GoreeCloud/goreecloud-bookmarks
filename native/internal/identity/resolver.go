package identity

import (
	"errors"
	"net/http"
)

var ErrUnavailable = errors.New("authenticated identity integration is unavailable")

// Resolver is the native GoreeCloud Bookmarks authentication boundary.
// Implementations must return a stable authenticated owner identifier and must
// fail closed when authentication cannot be established. Request headers or
// other caller-controlled values are not production authentication unless a
// reviewed GoreeCloud Identity adapter has authenticated and asserted them.
type Resolver interface {
	Resolve(*http.Request) (string, error)
}

// UnavailableResolver is the production-safe default until the reviewed
// GoreeCloud Identity adapter is implemented and accepted.
type UnavailableResolver struct{}

func (UnavailableResolver) Resolve(*http.Request) (string, error) {
	return "", ErrUnavailable
}

var _ Resolver = DevelopmentHeaderResolver{}
var _ Resolver = UnavailableResolver{}
