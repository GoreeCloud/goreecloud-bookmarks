package collections

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"
)

var (
	ErrBookmarkRequired = errors.New("bookmark identifier is required")
	ErrBookmarkNotFound = errors.New("bookmark was not found")
)

// BookmarkLookup is the minimum owner-scoped bookmark contract required by
// collection assignment. Implementations must return false for cross-owner
// lookups so assignment never needs access to another owner's bookmark data.
// Storage failures remain errors and must not be collapsed into not-found.
type BookmarkLookup interface {
	Exists(context.Context, string, string) (bool, error)
}

type Assignment struct {
	OwnerID      string    `json:"ownerId"`
	BookmarkID   string    `json:"bookmarkId"`
	CollectionID string    `json:"collectionId"`
	AssignedAt   time.Time `json:"assignedAt"`
}

type AssignmentStore struct {
	mu          sync.RWMutex
	collections *Store
	bookmarks   BookmarkLookup
	items       map[string]Assignment
	now         func() time.Time
}

func NewAssignmentStore(collections *Store, bookmarks BookmarkLookup) *AssignmentStore {
	return &AssignmentStore{
		collections: collections,
		bookmarks:   bookmarks,
		items:       map[string]Assignment{},
		now:         func() time.Time { return time.Now().UTC() },
	}
}

func (s *AssignmentStore) Assign(ctx context.Context, ownerID, bookmarkID, collectionID string) (Assignment, error) {
	ownerID = strings.TrimSpace(ownerID)
	bookmarkID = strings.TrimSpace(bookmarkID)
	collectionID = strings.TrimSpace(collectionID)
	if ownerID == "" {
		return Assignment{}, ErrOwnerRequired
	}
	if bookmarkID == "" {
		return Assignment{}, ErrBookmarkRequired
	}
	if collectionID == "" {
		return Assignment{}, ErrCollectionNotFound
	}
	if s.bookmarks == nil {
		return Assignment{}, ErrBookmarkNotFound
	}
	found, err := s.bookmarks.Exists(ctx, ownerID, bookmarkID)
	if err != nil {
		return Assignment{}, err
	}
	if !found {
		return Assignment{}, ErrBookmarkNotFound
	}
	if s.collections == nil {
		return Assignment{}, ErrCollectionNotFound
	}
	if _, ok := s.collections.Get(ownerID, collectionID); !ok {
		return Assignment{}, ErrCollectionNotFound
	}

	assignment := Assignment{
		OwnerID:      ownerID,
		BookmarkID:   bookmarkID,
		CollectionID: collectionID,
		AssignedAt:   s.now(),
	}
	s.mu.Lock()
	s.items[assignmentKey(ownerID, bookmarkID)] = assignment
	s.mu.Unlock()
	return assignment, nil
}

func (s *AssignmentStore) Get(ownerID, bookmarkID string) (Assignment, bool) {
	ownerID = strings.TrimSpace(ownerID)
	bookmarkID = strings.TrimSpace(bookmarkID)
	if ownerID == "" || bookmarkID == "" {
		return Assignment{}, false
	}
	s.mu.RLock()
	assignment, ok := s.items[assignmentKey(ownerID, bookmarkID)]
	s.mu.RUnlock()
	return assignment, ok
}

func (s *AssignmentStore) Remove(ownerID, bookmarkID string) bool {
	ownerID = strings.TrimSpace(ownerID)
	bookmarkID = strings.TrimSpace(bookmarkID)
	if ownerID == "" || bookmarkID == "" {
		return false
	}
	key := assignmentKey(ownerID, bookmarkID)
	s.mu.Lock()
	_, ok := s.items[key]
	if ok {
		delete(s.items, key)
	}
	s.mu.Unlock()
	return ok
}

func assignmentKey(ownerID, bookmarkID string) string {
	return ownerID + "\x00" + bookmarkID
}
