package postgres

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// ErrBookmarkNotFound avoids exposing whether a bookmark exists outside the
// caller's authorized owner boundary.
var ErrBookmarkNotFound = errors.New("bookmark not found")

// Bookmark is the initial persisted subset of the selected logical Bookmark
// model. Additional fields can be added through governed migrations as their
// product behavior is implemented.
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

// CreateBookmark persists one new bookmark record. It is an internal data-layer
// capability; no unauthenticated or end-user HTTP creation route is exposed by
// this package.
func (d *Database) CreateBookmark(ctx context.Context, bookmark Bookmark) (Bookmark, error) {
	if d == nil || d.pool == nil {
		return Bookmark{}, errors.New("database is not initialized")
	}
	if strings.TrimSpace(bookmark.ID) == "" || strings.TrimSpace(bookmark.OwnerID) == "" || strings.TrimSpace(bookmark.URL) == "" {
		return Bookmark{}, errors.New("bookmark id, owner id, and URL are required")
	}
	if bookmark.ReadState == "" {
		bookmark.ReadState = "unread"
	}
	if bookmark.PrivacyLevel == "" {
		bookmark.PrivacyLevel = "normal"
	}
	if bookmark.Title == "" {
		bookmark.Title = bookmark.URL
	}

	err := d.pool.QueryRow(ctx, `
INSERT INTO bookmarks (
    id, owner_id, revision, url, original_url, title, description,
    favorite, pinned, read_state, reading_progress, privacy_level,
    offline_requested
) VALUES ($1, $2, 1, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING revision, created_at, updated_at, deleted_at`,
		bookmark.ID,
		bookmark.OwnerID,
		bookmark.URL,
		bookmark.OriginalURL,
		bookmark.Title,
		bookmark.Description,
		bookmark.Favorite,
		bookmark.Pinned,
		bookmark.ReadState,
		bookmark.ReadingProgress,
		bookmark.PrivacyLevel,
		bookmark.OfflineRequested,
	).Scan(&bookmark.Revision, &bookmark.CreatedAt, &bookmark.UpdatedAt, &bookmark.DeletedAt)
	if err != nil {
		return Bookmark{}, errors.New("create bookmark failed")
	}

	return bookmark, nil
}

// GetBookmark returns one active bookmark only when both ID and owner match.
func (d *Database) GetBookmark(ctx context.Context, ownerID, bookmarkID string) (Bookmark, error) {
	if d == nil || d.pool == nil {
		return Bookmark{}, errors.New("database is not initialized")
	}

	var bookmark Bookmark
	err := d.pool.QueryRow(ctx, `
SELECT id, owner_id, revision, url, original_url, title, description,
       favorite, pinned, read_state, reading_progress, privacy_level,
       offline_requested, created_at, updated_at, deleted_at
FROM bookmarks
WHERE id = $1 AND owner_id = $2 AND deleted_at IS NULL`, bookmarkID, ownerID).Scan(
		&bookmark.ID,
		&bookmark.OwnerID,
		&bookmark.Revision,
		&bookmark.URL,
		&bookmark.OriginalURL,
		&bookmark.Title,
		&bookmark.Description,
		&bookmark.Favorite,
		&bookmark.Pinned,
		&bookmark.ReadState,
		&bookmark.ReadingProgress,
		&bookmark.PrivacyLevel,
		&bookmark.OfflineRequested,
		&bookmark.CreatedAt,
		&bookmark.UpdatedAt,
		&bookmark.DeletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Bookmark{}, ErrBookmarkNotFound
	}
	if err != nil {
		return Bookmark{}, errors.New("read bookmark failed")
	}

	return bookmark, nil
}
