package identity

import (
	"errors"
	"net/http"
	"strings"
	"unicode"
)

const (
	// DevelopmentOwnerHeader is intentionally development-only. Production
	// identity must come from the reviewed GoreeCloud identity integration.
	DevelopmentOwnerHeader = "X-GoreeCloud-Development-Owner"
	maxDevelopmentOwnerRunes = 128
)

var (
	ErrOwnerRequired = errors.New("development owner identity is required")
	ErrOwnerInvalid  = errors.New("development owner identity is invalid")
)

// DevelopmentHeaderResolver provides an explicit local-development seam for
// exercising owner-scoped application behavior before the production identity
// adapter exists. It must not be treated as production authentication.
type DevelopmentHeaderResolver struct{}

func (DevelopmentHeaderResolver) Resolve(r *http.Request) (string, error) {
	if r == nil {
		return "", ErrOwnerRequired
	}
	ownerID := strings.TrimSpace(r.Header.Get(DevelopmentOwnerHeader))
	if ownerID == "" {
		return "", ErrOwnerRequired
	}
	if len([]rune(ownerID)) > maxDevelopmentOwnerRunes {
		return "", ErrOwnerInvalid
	}
	for _, value := range ownerID {
		if unicode.IsControl(value) || unicode.IsSpace(value) {
			return "", ErrOwnerInvalid
		}
	}
	return ownerID, nil
}
