package syncstate

import (
	"errors"
	"strings"
	"time"

	bookmarkcore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
)

var ErrInvalidBookmarkItem = errors.New("invalid bookmark sync item")

// ItemRecord is the Bookmarks-owned payload for the bookmarks.items dataset.
// Owner identity is intentionally excluded: authenticated Sync account context
// carries ownership and prevents leaking account identifiers into record data.
type ItemRecord struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Title     string    `json:"title,omitempty"`
	Note      string    `json:"note,omitempty"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ExportBookmarkItem(bookmark bookmarkcore.Bookmark) (ItemRecord, error) {
	if strings.TrimSpace(bookmark.ID) == "" || strings.TrimSpace(bookmark.URL) == "" || bookmark.CreatedAt.IsZero() || bookmark.UpdatedAt.IsZero() {
		return ItemRecord{}, ErrInvalidBookmarkItem
	}
	return ItemRecord{
		ID: bookmark.ID, URL: bookmark.URL, Title: bookmark.Title, Note: bookmark.Note,
		Tags: append([]string(nil), bookmark.Tags...), CreatedAt: bookmark.CreatedAt.UTC(), UpdatedAt: bookmark.UpdatedAt.UTC(),
	}, nil
}
