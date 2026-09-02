package bookmarks

import (
	"context"
	"errors"
	"strings"
)

var ErrOwnerIdentityRequired = errors.New("owner identity is required")

type Repository interface {
	Create(context.Context, string, CreateInput) (Bookmark, error)
	List(context.Context, string) ([]Bookmark, error)
	Get(context.Context, string, string) (Bookmark, bool, error)
	Update(context.Context, string, string, UpdateInput) (Bookmark, bool, error)
	Delete(context.Context, string, string) (bool, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) (*Service, error) {
	if repository == nil {
		return nil, errors.New("bookmark repository is required")
	}
	return &Service{repository: repository}, nil
}

func (s *Service) Create(ctx context.Context, ownerID string, input CreateInput) (Bookmark, error) {
	ownerID, err := normalizeOwnerID(ownerID)
	if err != nil {
		return Bookmark{}, err
	}
	return s.repository.Create(ctx, ownerID, input)
}

func (s *Service) List(ctx context.Context, ownerID string) ([]Bookmark, error) {
	ownerID, err := normalizeOwnerID(ownerID)
	if err != nil {
		return nil, err
	}
	return s.repository.List(ctx, ownerID)
}

func (s *Service) Get(ctx context.Context, ownerID, id string) (Bookmark, bool, error) {
	ownerID, id, err := normalizeMutationIdentity(ownerID, id)
	if err != nil {
		return Bookmark{}, false, err
	}
	return s.repository.Get(ctx, ownerID, id)
}

// Exists is the minimal owner-scoped lookup used by cross-domain relationships
// such as collection assignment. It preserves repository failures so callers
// can fail closed instead of silently converting storage errors into not-found.
func (s *Service) Exists(ctx context.Context, ownerID, id string) (bool, error) {
	_, found, err := s.Get(ctx, ownerID, id)
	return found, err
}

func (s *Service) Update(ctx context.Context, ownerID, id string, input UpdateInput) (Bookmark, bool, error) {
	ownerID, id, err := normalizeMutationIdentity(ownerID, id)
	if err != nil {
		return Bookmark{}, false, err
	}
	return s.repository.Update(ctx, ownerID, id, input)
}

func (s *Service) Delete(ctx context.Context, ownerID, id string) (bool, error) {
	ownerID, id, err := normalizeMutationIdentity(ownerID, id)
	if err != nil {
		return false, err
	}
	return s.repository.Delete(ctx, ownerID, id)
}

func normalizeOwnerID(ownerID string) (string, error) {
	ownerID = strings.TrimSpace(ownerID)
	if ownerID == "" {
		return "", ErrOwnerIdentityRequired
	}
	return ownerID, nil
}

func normalizeMutationIdentity(ownerID, id string) (string, string, error) {
	ownerID, err := normalizeOwnerID(ownerID)
	if err != nil {
		return "", "", err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return "", "", errors.New("bookmark id is required")
	}
	return ownerID, id, nil
}
