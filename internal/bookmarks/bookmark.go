package bookmarks

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrNotFound preserves a non-enumerating resource boundary for callers.
	ErrNotFound = errors.New("bookmark not found")
	// ErrIdempotencyConflict indicates reuse of an idempotency key for a
	// materially different create request.
	ErrIdempotencyConflict = errors.New("idempotency key conflicts with an existing request")
)

// Bookmark is the currently implemented persisted subset of the planned
// Bookmark resource. Fields are added only when their behavior and storage are
// implemented and validated.
type Bookmark struct {
	ID               string
	OwnerID          string
	Revision         int64
	URL              string
	OriginalURL      *string
	Title            string
	Description      *string
	Favorite         bool
	Pinned           bool
	ReadState        string
	ReadingProgress  *float64
	PrivacyLevel     string
	OfflineRequested bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

// CreateInput represents the planned create fields at the domain boundary.
// Collection/tag/archive inputs are retained here so unimplemented behavior is
// rejected explicitly instead of being silently discarded.
type CreateInput struct {
	URL              string
	Title            *string
	Description      *string
	CollectionIDs    []string
	TagIDs            []string
	Favorite         bool
	ReadState        string
	PrivacyLevel     string
	ArchiveMode      *string
	OfflineRequested bool
}

// CreateResult distinguishes a newly committed bookmark from an idempotent
// replay of the same logical create request.
type CreateResult struct {
	Bookmark Bookmark
	Replayed bool
}

// Store is the persistence contract required by the Bookmark service. The
// owner boundary is explicit in every operation.
type Store interface {
	CreateBookmarkIdempotent(
		ctx context.Context,
		ownerID string,
		idempotencyKey string,
		requestHash string,
		bookmark Bookmark,
	) (Bookmark, bool, error)
	GetBookmark(ctx context.Context, ownerID, bookmarkID string) (Bookmark, error)
}
