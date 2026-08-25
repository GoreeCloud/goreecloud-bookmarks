package collections

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

const MaxNameRunes = 256

var (
	ErrOwnerRequired      = errors.New("collection owner identity is required")
	ErrNameRequired       = errors.New("collection name is required")
	ErrNameTooLong        = errors.New("collection name exceeds maximum length")
	ErrParentNotFound     = errors.New("collection parent was not found")
	ErrDuplicateSibling   = errors.New("collection name already exists under parent")
	ErrCollectionNotFound = errors.New("collection was not found")
	ErrHierarchyCycle     = errors.New("collection hierarchy cycle is not allowed")
)

type Collection struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"ownerId"`
	Name      string    `json:"name"`
	ParentID  string    `json:"parentId,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Store struct {
	mu    sync.RWMutex
	items map[string]Collection
	now   func() time.Time
}

func NewStore() *Store {
	return &Store{items: map[string]Collection{}, now: func() time.Time { return time.Now().UTC() }}
}

func (s *Store) Create(ownerID, name, parentID string) (Collection, error) {
	ownerID = strings.TrimSpace(ownerID)
	if ownerID == "" {
		return Collection{}, ErrOwnerRequired
	}
	name, err := normalizeName(name)
	if err != nil {
		return Collection{}, err
	}
	parentID = strings.TrimSpace(parentID)

	s.mu.Lock()
	defer s.mu.Unlock()
	if parentID != "" {
		parent, ok := s.items[parentID]
		if !ok || parent.OwnerID != ownerID {
			return Collection{}, ErrParentNotFound
		}
	}
	if s.siblingNameExistsLocked(ownerID, parentID, name, "") {
		return Collection{}, ErrDuplicateSibling
	}
	id, err := newID()
	if err != nil {
		return Collection{}, err
	}
	now := s.now()
	collection := Collection{ID: id, OwnerID: ownerID, Name: name, ParentID: parentID, CreatedAt: now, UpdatedAt: now}
	s.items[id] = collection
	return collection, nil
}

func (s *Store) List(ownerID string) []Collection {
	ownerID = strings.TrimSpace(ownerID)
	if ownerID == "" {
		return []Collection{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]Collection, 0)
	for _, item := range s.items {
		if item.OwnerID == ownerID {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].ParentID == items[j].ParentID {
			if strings.EqualFold(items[i].Name, items[j].Name) {
				return items[i].ID < items[j].ID
			}
			return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
		}
		return items[i].ParentID < items[j].ParentID
	})
	return items
}

func (s *Store) Get(ownerID, id string) (Collection, bool) {
	ownerID = strings.TrimSpace(ownerID)
	id = strings.TrimSpace(id)
	if ownerID == "" || id == "" {
		return Collection{}, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	if !ok || item.OwnerID != ownerID {
		return Collection{}, false
	}
	return item, true
}

func (s *Store) Move(ownerID, id, parentID string) (Collection, error) {
	ownerID = strings.TrimSpace(ownerID)
	id = strings.TrimSpace(id)
	parentID = strings.TrimSpace(parentID)
	if ownerID == "" {
		return Collection{}, ErrOwnerRequired
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[id]
	if !ok || item.OwnerID != ownerID {
		return Collection{}, ErrCollectionNotFound
	}
	if parentID == id {
		return Collection{}, ErrHierarchyCycle
	}
	if parentID != "" {
		parent, ok := s.items[parentID]
		if !ok || parent.OwnerID != ownerID {
			return Collection{}, ErrParentNotFound
		}
		for ancestor := parent; ancestor.ParentID != ""; {
			if ancestor.ParentID == id {
				return Collection{}, ErrHierarchyCycle
			}
			next, ok := s.items[ancestor.ParentID]
			if !ok || next.OwnerID != ownerID {
				return Collection{}, ErrParentNotFound
			}
			ancestor = next
		}
	}
	if s.siblingNameExistsLocked(ownerID, parentID, item.Name, id) {
		return Collection{}, ErrDuplicateSibling
	}
	item.ParentID = parentID
	item.UpdatedAt = s.now()
	s.items[id] = item
	return item, nil
}

func (s *Store) siblingNameExistsLocked(ownerID, parentID, name, exceptID string) bool {
	for _, item := range s.items {
		if item.ID != exceptID && item.OwnerID == ownerID && item.ParentID == parentID && strings.EqualFold(item.Name, name) {
			return true
		}
	}
	return false
}

func normalizeName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", ErrNameRequired
	}
	if len([]rune(name)) > MaxNameRunes {
		return "", ErrNameTooLong
	}
	return name, nil
}

func newID() (string, error) {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes[:]), nil
}
