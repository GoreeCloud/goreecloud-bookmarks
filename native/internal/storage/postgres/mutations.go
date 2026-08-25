package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
)

const lockBookmarkForUpdateQuery = `
SELECT id, owner_id, canonical_url, title, note, created_at, updated_at
FROM native_bookmarks
WHERE owner_id = $1 AND id = $2
FOR UPDATE
`

const updateBookmarkQuery = `
UPDATE native_bookmarks
SET canonical_url = $3, title = $4, note = $5, updated_at = $6
WHERE owner_id = $1 AND id = $2
`

const deleteBookmarkTagsQuery = `
DELETE FROM native_bookmark_tags
WHERE owner_id = $1 AND bookmark_id = $2
`

const deleteBookmarkQuery = `
DELETE FROM native_bookmarks
WHERE owner_id = $1 AND id = $2
`

func (r *Repository) Update(ctx context.Context, ownerID, id string, input bookmarks.UpdateInput) (bookmarks.Bookmark, bool, error) {
	normalized, err := bookmarks.NewBookmark(ownerID, bookmarks.CreateInput{
		URL: input.URL, Title: input.Title, Note: input.Note, Tags: input.Tags,
	}, r.now())
	if err != nil {
		return bookmarks.Bookmark{}, false, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return bookmarks.Bookmark{}, false, errors.New("begin bookmark update transaction failed")
	}
	defer func() { _ = tx.Rollback() }()

	current, err := scanBookmark(tx.QueryRowContext(ctx, lockBookmarkForUpdateQuery, ownerID, id))
	if errors.Is(err, sql.ErrNoRows) {
		return bookmarks.Bookmark{}, false, nil
	}
	if err != nil {
		return bookmarks.Bookmark{}, false, err
	}

	updatedAt := r.now()
	if _, err := tx.ExecContext(ctx, updateBookmarkQuery,
		ownerID, id, normalized.URL, normalized.Title, normalized.Note, updatedAt,
	); err != nil {
		if isUniqueViolation(err) {
			return bookmarks.Bookmark{}, false, ErrBookmarkAlreadyExists
		}
		return bookmarks.Bookmark{}, false, errors.New("update bookmark failed")
	}
	if _, err := tx.ExecContext(ctx, deleteBookmarkTagsQuery, ownerID, id); err != nil {
		return bookmarks.Bookmark{}, false, errors.New("replace bookmark tags failed")
	}
	for _, tag := range normalized.Tags {
		if _, err := tx.ExecContext(ctx, insertTagQuery, id, ownerID, tag); err != nil {
			return bookmarks.Bookmark{}, false, errors.New("replace bookmark tags failed")
		}
	}
	if err := tx.Commit(); err != nil {
		return bookmarks.Bookmark{}, false, errors.New("commit bookmark update transaction failed")
	}

	current.URL = normalized.URL
	current.Title = normalized.Title
	current.Note = normalized.Note
	current.Tags = append([]string(nil), normalized.Tags...)
	current.UpdatedAt = updatedAt
	return current, true, nil
}

func (r *Repository) Delete(ctx context.Context, ownerID, id string) (bool, error) {
	result, err := r.db.ExecContext(ctx, deleteBookmarkQuery, ownerID, id)
	if err != nil {
		return false, errors.New("delete bookmark failed")
	}
	count, err := result.RowsAffected()
	if err != nil {
		return false, errors.New("read delete bookmark result failed")
	}
	return count > 0, nil
}
