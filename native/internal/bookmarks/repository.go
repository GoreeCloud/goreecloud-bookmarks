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
	ownerID, err := normalizeOwnerID(ownerID)
	if err != nil {
		return Bookmark{}, false, err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return Bookmark{}, false, errors.New("bookmark id is required")
	}
	return s.repository.Get(ctx, ownerID, id)
}

func normalizeOwnerID(ownerID string) (string, error) {
	ownerID = strings.TrimSpace(ownerID)
	if ownerID == "" {
		return "", ErrOwnerIdentityRequired
	}
	return ownerID, nil
}
