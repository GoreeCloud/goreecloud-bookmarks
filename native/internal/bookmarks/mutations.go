package bookmarks

import (
	"errors"
	"strings"
)

type UpdateInput struct {
	URL   string
	Title string
	Note  string
	Tags  []string
}

func (s *Store) Update(ownerID, id string, input UpdateInput) (Bookmark, bool, error) {
	ownerID = strings.TrimSpace(ownerID)
	id = strings.TrimSpace(id)
	if ownerID == "" {
		return Bookmark{}, false, ErrOwnerIdentityRequired
	}
	if id == "" {
		return Bookmark{}, false, errors.New("bookmark id is required")
	}

	normalizedURL, err := normalizeURL(input.URL)
	if err != nil {
		return Bookmark{}, false, err
	}
	title := strings.TrimSpace(input.Title)
	if len([]rune(title)) > MaxTitleRunes {
		return Bookmark{}, false, errors.New("bookmark title exceeds maximum length")
	}
	note := strings.TrimSpace(input.Note)
	if len([]rune(note)) > MaxNoteRunes {
		return Bookmark{}, false, errors.New("bookmark note exceeds maximum length")
	}
	tags, err := normalizeTags(input.Tags)
	if err != nil {
		return Bookmark{}, false, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.items[id]
	if !ok || current.OwnerID != ownerID {
		return Bookmark{}, false, nil
	}
	for existingID, existing := range s.items {
		if existingID != id && existing.OwnerID == ownerID && existing.URL == normalizedURL {
			return Bookmark{}, false, errors.New("bookmark already exists for owner")
		}
	}
	current.URL = normalizedURL
	current.Title = title
	current.Note = note
	current.Tags = tags
	current.UpdatedAt = s.now()
	s.items[id] = current
	current.Tags = append([]string(nil), current.Tags...)
	return current, true, nil
}

func (s *Store) Delete(ownerID, id string) (bool, error) {
	ownerID = strings.TrimSpace(ownerID)
	id = strings.TrimSpace(id)
	if ownerID == "" {
		return false, ErrOwnerIdentityRequired
	}
	if id == "" {
		return false, errors.New("bookmark id is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.items[id]
	if !ok || current.OwnerID != ownerID {
		return false, nil
	}
	delete(s.items, id)
	return true, nil
}
