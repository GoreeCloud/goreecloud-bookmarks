package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/GoreeCloud/goreecloud-bookmarks/internal/bookmarks"
	"github.com/jackc/pgx/v5"
)

// CreateBookmarkIdempotent atomically claims an owner-scoped idempotency key
// and persists a bookmark. Replays of the same normalized request return the
// original bookmark; reuse of the key for different content fails closed.
func (d *Database) CreateBookmarkIdempotent(
	ctx context.Context,
	ownerID string,
	idempotencyKey string,
	requestHash string,
	bookmark bookmarks.Bookmark,
) (bookmarks.Bookmark, bool, error) {
	if d == nil || d.pool == nil {
		return bookmarks.Bookmark{}, false, errors.New("database is not initialized")
	}

	ownerID = strings.TrimSpace(ownerID)
	idempotencyKey = strings.TrimSpace(idempotencyKey)
	if ownerID == "" ||
		bookmark.OwnerID != ownerID ||
		strings.TrimSpace(bookmark.ID) == "" ||
		strings.TrimSpace(bookmark.URL) == "" ||
		len(idempotencyKey) < 8 ||
		len(idempotencyKey) > 200 ||
		len(requestHash) != 64 {
		return bookmarks.Bookmark{}, false, errors.New("bookmark create persistence input is invalid")
	}

	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return bookmarks.Bookmark{}, false, errors.New("begin bookmark create transaction failed")
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	tag, err := tx.Exec(ctx, `
INSERT INTO bookmark_create_idempotency (
    owner_id, idempotency_key, request_hash, bookmark_id
) VALUES ($1, $2, $3, $4)
ON CONFLICT (owner_id, idempotency_key) DO NOTHING`,
		ownerID,
		idempotencyKey,
		requestHash,
		bookmark.ID,
	)
	if err != nil {
		return bookmarks.Bookmark{}, false, errors.New("claim bookmark idempotency key failed")
	}

	if tag.RowsAffected() == 1 {
		created, err := insertBookmark(ctx, tx, bookmark)
		if err != nil {
			return bookmarks.Bookmark{}, false, err
		}
		if err := tx.Commit(ctx); err != nil {
			return bookmarks.Bookmark{}, false, errors.New("commit bookmark create failed")
		}
		return created, false, nil
	}

	var existingHash, existingBookmarkID string
	err = tx.QueryRow(ctx, `
SELECT request_hash, bookmark_id
FROM bookmark_create_idempotency
WHERE owner_id = $1 AND idempotency_key = $2
FOR UPDATE`, ownerID, idempotencyKey).Scan(&existingHash, &existingBookmarkID)
	if err != nil {
		return bookmarks.Bookmark{}, false, errors.New("read bookmark idempotency key failed")
	}
	if existingHash != requestHash {
		return bookmarks.Bookmark{}, false, bookmarks.ErrIdempotencyConflict
	}

	existing, err := getBookmark(ctx, tx, ownerID, existingBookmarkID)
	if err != nil {
		return bookmarks.Bookmark{}, false, err
	}
	if err := tx.Commit(ctx); err != nil {
		return bookmarks.Bookmark{}, false, errors.New("commit bookmark replay read failed")
	}
	return existing, true, nil
}

// GetBookmark returns one active bookmark only when both ID and owner match.
func (d *Database) GetBookmark(
	ctx context.Context,
	ownerID string,
	bookmarkID string,
) (bookmarks.Bookmark, error) {
	if d == nil || d.pool == nil {
		return bookmarks.Bookmark{}, errors.New("database is not initialized")
	}
	return getBookmark(ctx, d.pool, ownerID, bookmarkID)
}

type bookmarkQueryer interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func insertBookmark(
	ctx context.Context,
	q bookmarkQueryer,
	bookmark bookmarks.Bookmark,
) (bookmarks.Bookmark, error) {
	err := q.QueryRow(ctx, `
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
	).Scan(
		&bookmark.Revision,
		&bookmark.CreatedAt,
		&bookmark.UpdatedAt,
		&bookmark.DeletedAt,
	)
	if err != nil {
		return bookmarks.Bookmark{}, errors.New("create bookmark failed")
	}
	return bookmark, nil
}

func getBookmark(
	ctx context.Context,
	q bookmarkQueryer,
	ownerID string,
	bookmarkID string,
) (bookmarks.Bookmark, error) {
	var bookmark bookmarks.Bookmark
	err := q.QueryRow(ctx, `
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
		return bookmarks.Bookmark{}, bookmarks.ErrNotFound
	}
	if err != nil {
		return bookmarks.Bookmark{}, errors.New("read bookmark failed")
	}
	return bookmark, nil
}
