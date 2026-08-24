package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
)

const uniqueViolationSQLState = "23505"

var ErrBookmarkAlreadyExists = errors.New("bookmark already exists for owner")

type Repository struct {
	db  *sql.DB
	now func() time.Time
}

func NewRepository(db *sql.DB) (*Repository, error) {
	if db == nil {
		return nil, errors.New("postgres database is required")
	}
	return &Repository{db: db, now: func() time.Time { return time.Now().UTC() }}, nil
}

func (r *Repository) Create(ctx context.Context, ownerID string, input bookmarks.CreateInput) (bookmarks.Bookmark, error) {
	bookmark, err := bookmarks.NewBookmark(ownerID, input, r.now())
	if err != nil {
		return bookmarks.Bookmark{}, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return bookmarks.Bookmark{}, errors.New("begin bookmark transaction failed")
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
INSERT INTO native_bookmarks (id, owner_id, canonical_url, title, note, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`, bookmark.ID, bookmark.OwnerID, bookmark.URL, bookmark.Title, bookmark.Note, bookmark.CreatedAt, bookmark.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return bookmarks.Bookmark{}, ErrBookmarkAlreadyExists
		}
		return bookmarks.Bookmark{}, errors.New("create bookmark failed")
	}

	for _, tag := range bookmark.Tags {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO native_bookmark_tags (bookmark_id, owner_id, tag)
VALUES ($1, $2, $3)
`, bookmark.ID, bookmark.OwnerID, tag); err != nil {
			return bookmarks.Bookmark{}, errors.New("create bookmark tags failed")
		}
	}

	if err := tx.Commit(); err != nil {
		return bookmarks.Bookmark{}, errors.New("commit bookmark transaction failed")
	}
	return bookmark, nil
}

func (r *Repository) List(ctx context.Context, ownerID string) ([]bookmarks.Bookmark, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, owner_id, canonical_url, title, note, created_at, updated_at
FROM native_bookmarks
WHERE owner_id = $1
ORDER BY created_at DESC, id DESC
`, ownerID)
	if err != nil {
		return nil, errors.New("list bookmarks failed")
	}
	defer rows.Close()

	items := make([]bookmarks.Bookmark, 0)
	for rows.Next() {
		item, err := scanBookmark(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("iterate bookmarks failed")
	}

	for index := range items {
		tags, err := r.listTags(ctx, ownerID, items[index].ID)
		if err != nil {
			return nil, err
		}
		items[index].Tags = tags
	}
	return items, nil
}

func (r *Repository) Get(ctx context.Context, ownerID, id string) (bookmarks.Bookmark, bool, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, owner_id, canonical_url, title, note, created_at, updated_at
FROM native_bookmarks
WHERE owner_id = $1 AND id = $2
`, ownerID, id)

	item, err := scanBookmark(row)
	if errors.Is(err, sql.ErrNoRows) {
		return bookmarks.Bookmark{}, false, nil
	}
	if err != nil {
		return bookmarks.Bookmark{}, false, err
	}

	item.Tags, err = r.listTags(ctx, ownerID, id)
	if err != nil {
		return bookmarks.Bookmark{}, false, err
	}
	return item, true, nil
}

func (r *Repository) listTags(ctx context.Context, ownerID, bookmarkID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT tag
FROM native_bookmark_tags
WHERE owner_id = $1 AND bookmark_id = $2
ORDER BY lower(tag), tag
`, ownerID, bookmarkID)
	if err != nil {
		return nil, errors.New("list bookmark tags failed")
	}
	defer rows.Close()

	tags := make([]string, 0)
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, errors.New("scan bookmark tag failed")
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("iterate bookmark tags failed")
	}
	return tags, nil
}

type scanner interface {
	Scan(...any) error
}

func scanBookmark(row scanner) (bookmarks.Bookmark, error) {
	var item bookmarks.Bookmark
	if err := row.Scan(&item.ID, &item.OwnerID, &item.URL, &item.Title, &item.Note, &item.CreatedAt, &item.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return bookmarks.Bookmark{}, sql.ErrNoRows
		}
		return bookmarks.Bookmark{}, errors.New("scan bookmark failed")
	}
	item.Tags = []string{}
	return item, nil
}

type sqlStateError interface {
	SQLState() string
}

func isUniqueViolation(err error) bool {
	var state sqlStateError
	return errors.As(err, &state) && state.SQLState() == uniqueViolationSQLState
}

var _ bookmarks.Repository = (*Repository)(nil)

func (r *Repository) String() string {
	return fmt.Sprintf("postgres.Repository(schema=%d)", SchemaVersion)
}
