package bookmarks

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	MaxTitleRunes = 512
	MaxNoteRunes  = 8192
	MaxTags       = 32
	MaxTagRunes   = 64
)

type Bookmark struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"ownerId"`
	URL       string    `json:"url"`
	Title     string    `json:"title,omitempty"`
	Note      string    `json:"note,omitempty"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateInput struct {
	URL   string
	Title string
	Note  string
	Tags  []string
}

type Store struct {
	mu    sync.RWMutex
	items map[string]Bookmark
	now   func() time.Time
}

func NewStore() *Store {
	return &Store{items: map[string]Bookmark{}, now: func() time.Time { return time.Now().UTC() }}
}

func (s *Store) Create(ownerID string, input CreateInput) (Bookmark, error) {
	bookmark, err := NewBookmark(ownerID, input, s.now())
	if err != nil {
		return Bookmark{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, existing := range s.items {
		if existing.OwnerID == bookmark.OwnerID && existing.URL == bookmark.URL {
			return Bookmark{}, errors.New("bookmark already exists for owner")
		}
	}
	s.items[bookmark.ID] = bookmark
	return bookmark, nil
}

func (s *Store) List(ownerID string) []Bookmark {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]Bookmark, 0)
	for _, item := range s.items {
		if item.OwnerID == ownerID {
			copyItem := item
			copyItem.Tags = append([]string(nil), item.Tags...)
			items = append(items, copyItem)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].ID < items[j].ID
		}
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	return items
}

func (s *Store) Get(ownerID, id string) (Bookmark, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	if !ok || item.OwnerID != ownerID {
		return Bookmark{}, false
	}
	item.Tags = append([]string(nil), item.Tags...)
	return item, true
}

func normalizeURL(raw string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", errors.New("bookmark URL must use http or https")
	}
	parsed.Fragment = ""
	return parsed.String(), nil
}

func normalizeTags(raw []string) ([]string, error) {
	seen := map[string]struct{}{}
	tags := make([]string, 0, len(raw))
	for _, value := range raw {
		tag := strings.TrimSpace(value)
		if tag == "" {
			continue
		}
		if len([]rune(tag)) > MaxTagRunes {
			return nil, errors.New("tag exceeds maximum length")
		}
		key := strings.ToLower(tag)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		tags = append(tags, tag)
	}
	if len(tags) > MaxTags {
		return nil, errors.New("too many tags")
	}
	sort.Slice(tags, func(i, j int) bool { return strings.ToLower(tags[i]) < strings.ToLower(tags[j]) })
	return tags, nil
}

func newID() (string, error) {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes[:]), nil
}
