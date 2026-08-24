package bookmarks

import (
	"errors"
	"strings"
	"time"
)

// NewBookmark validates and normalizes a new owner-scoped bookmark before it
// reaches any concrete repository implementation.
func NewBookmark(ownerID string, input CreateInput, now time.Time) (Bookmark, error) {
	ownerID = strings.TrimSpace(ownerID)
	if ownerID == "" {
		return Bookmark{}, ErrOwnerIdentityRequired
	}

	cleanURL, err := normalizeURL(input.URL)
	if err != nil {
		return Bookmark{}, err
	}
	title := strings.TrimSpace(input.Title)
	note := strings.TrimSpace(input.Note)
	if len([]rune(title)) > MaxTitleRunes {
		return Bookmark{}, errors.New("title exceeds maximum length")
	}
	if len([]rune(note)) > MaxNoteRunes {
		return Bookmark{}, errors.New("note exceeds maximum length")
	}
	tags, err := normalizeTags(input.Tags)
	if err != nil {
		return Bookmark{}, err
	}

	id, err := newID()
	if err != nil {
		return Bookmark{}, err
	}
	if now.IsZero() {
		now = time.Now().UTC()
	} else {
		now = now.UTC()
	}

	return Bookmark{
		ID:        id,
		OwnerID:   ownerID,
		URL:       cleanURL,
		Title:     title,
		Note:      note,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
