package bookmarks

import "context"

type MemoryRepository struct {
	store *Store
}

func NewMemoryRepository(store *Store) *MemoryRepository {
	if store == nil {
		store = NewStore()
	}
	return &MemoryRepository{store: store}
}

func (r *MemoryRepository) Create(_ context.Context, ownerID string, input CreateInput) (Bookmark, error) {
	return r.store.Create(ownerID, input)
}

func (r *MemoryRepository) List(_ context.Context, ownerID string) ([]Bookmark, error) {
	return r.store.List(ownerID), nil
}

func (r *MemoryRepository) Get(_ context.Context, ownerID, id string) (Bookmark, bool, error) {
	bookmark, ok := r.store.Get(ownerID, id)
	return bookmark, ok, nil
}
