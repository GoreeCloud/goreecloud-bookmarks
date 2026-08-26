package bookmarks

import "strings"

// Exists reports whether the bookmark exists for the supplied owner. It is
// deliberately owner-scoped so callers cannot distinguish another owner's
// bookmark from a missing bookmark.
func (s *Store) Exists(ownerID, id string) bool {
	ownerID = strings.TrimSpace(ownerID)
	id = strings.TrimSpace(id)
	if ownerID == "" || id == "" {
		return false
	}
	_, ok := s.Get(ownerID, id)
	return ok
}
