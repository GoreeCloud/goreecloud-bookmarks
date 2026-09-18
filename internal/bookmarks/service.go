package bookmarks

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	maxURLLength            = 8192
	maxTitleLength          = 4096
	maxDescriptionLength    = 16384
	minIdempotencyKeyLength = 8
	maxIdempotencyKeyLength = 200
)

type idGenerator func() (string, error)

// Service implements bookmark-domain behavior independently from HTTP and
// identity-provider mechanics. An authorized caller must supply ownerID.
type Service struct {
	store      Store
	generateID idGenerator
}

// NewService constructs the Bookmark domain service.
func NewService(store Store) (*Service, error) {
	return newService(store, newOpaqueID)
}

func newService(store Store, generateID idGenerator) (*Service, error) {
	if store == nil {
		return nil, errors.New("bookmark store is required")
	}
	if generateID == nil {
		return nil, errors.New("bookmark id generator is required")
	}
	return &Service{store: store, generateID: generateID}, nil
}

// Create durably persists the bookmark before any future enrichment work. The
// same owner/idempotency-key/request combination returns the original bookmark
// instead of creating a duplicate.
func (s *Service) Create(
	ctx context.Context,
	ownerID string,
	idempotencyKey string,
	input CreateInput,
) (CreateResult, error) {
	ownerID = strings.TrimSpace(ownerID)
	if ownerID == "" {
		return CreateResult{}, errors.New("owner id is required")
	}

	idempotencyKey = strings.TrimSpace(idempotencyKey)
	if len(idempotencyKey) < minIdempotencyKeyLength || len(idempotencyKey) > maxIdempotencyKeyLength {
		return CreateResult{}, fmt.Errorf(
			"idempotency key must be between %d and %d characters",
			minIdempotencyKeyLength,
			maxIdempotencyKeyLength,
		)
	}

	bookmark, err := normalizeCreateInput(ownerID, input)
	if err != nil {
		return CreateResult{}, err
	}

	requestHash, err := createRequestHash(bookmark)
	if err != nil {
		return CreateResult{}, err
	}

	bookmark.ID, err = s.generateID()
	if err != nil {
		return CreateResult{}, errors.New("generate bookmark id failed")
	}
	if strings.TrimSpace(bookmark.ID) == "" {
		return CreateResult{}, errors.New("generated bookmark id is empty")
	}

	created, replayed, err := s.store.CreateBookmarkIdempotent(
		ctx,
		ownerID,
		idempotencyKey,
		requestHash,
		bookmark,
	)
	if err != nil {
		return CreateResult{}, err
	}

	return CreateResult{Bookmark: created, Replayed: replayed}, nil
}

// Get returns one active bookmark inside the supplied owner boundary.
func (s *Service) Get(ctx context.Context, ownerID, bookmarkID string) (Bookmark, error) {
	ownerID = strings.TrimSpace(ownerID)
	bookmarkID = strings.TrimSpace(bookmarkID)
	if ownerID == "" || bookmarkID == "" {
		return Bookmark{}, ErrNotFound
	}
	return s.store.GetBookmark(ctx, ownerID, bookmarkID)
}

func normalizeCreateInput(ownerID string, input CreateInput) (Bookmark, error) {
	if len(input.CollectionIDs) > 0 || len(input.TagIDs) > 0 || input.ArchiveMode != nil {
		return Bookmark{}, errors.New("collection, tag, and archive create behavior is not implemented")
	}

	rawURL := strings.TrimSpace(input.URL)
	if rawURL == "" || len(rawURL) > maxURLLength {
		return Bookmark{}, errors.New("bookmark URL is required and must not exceed 8192 characters")
	}
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || !parsed.IsAbs() || parsed.Scheme == "" {
		return Bookmark{}, errors.New("bookmark URL must be an absolute URI")
	}

	title := rawURL
	if input.Title != nil {
		if len(*input.Title) > maxTitleLength {
			return Bookmark{}, errors.New("bookmark title exceeds 4096 characters")
		}
		if strings.TrimSpace(*input.Title) != "" {
			title = *input.Title
		}
	}

	if input.Description != nil && len(*input.Description) > maxDescriptionLength {
		return Bookmark{}, errors.New("bookmark description exceeds 16384 characters")
	}

	readState := input.ReadState
	if readState == "" {
		readState = "unread"
	}
	switch readState {
	case "unread", "reading", "completed", "archived":
	default:
		return Bookmark{}, errors.New("unsupported bookmark read state")
	}

	privacyLevel := input.PrivacyLevel
	if privacyLevel == "" {
		privacyLevel = "normal"
	}
	switch privacyLevel {
	case "normal", "sensitive", "private-vault", "shared":
	default:
		return Bookmark{}, errors.New("unsupported bookmark privacy level")
	}

	return Bookmark{
		OwnerID:          ownerID,
		URL:              rawURL,
		Title:            title,
		Description:      input.Description,
		Favorite:         input.Favorite,
		ReadState:        readState,
		PrivacyLevel:     privacyLevel,
		OfflineRequested: input.OfflineRequested,
	}, nil
}

func createRequestHash(bookmark Bookmark) (string, error) {
	canonical := struct {
		URL              string  `json:"url"`
		Title            string  `json:"title"`
		Description      *string `json:"description"`
		Favorite         bool    `json:"favorite"`
		ReadState        string  `json:"readState"`
		PrivacyLevel     string  `json:"privacyLevel"`
		OfflineRequested bool    `json:"offlineRequested"`
	}{
		URL:              bookmark.URL,
		Title:            bookmark.Title,
		Description:      bookmark.Description,
		Favorite:         bookmark.Favorite,
		ReadState:        bookmark.ReadState,
		PrivacyLevel:     bookmark.PrivacyLevel,
		OfflineRequested: bookmark.OfflineRequested,
	}

	encoded, err := json.Marshal(canonical)
	if err != nil {
		return "", errors.New("encode bookmark create request failed")
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}

func newOpaqueID() (string, error) {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes[:]), nil
}
