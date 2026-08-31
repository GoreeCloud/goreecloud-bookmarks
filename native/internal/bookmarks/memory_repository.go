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

func (r *MemoryRepository) Update(_ context.Context, ownerID, id string, input UpdateInput) (Bookmark, bool, error) {
	return r.store.Update(ownerID, id, input)
}

func (r *MemoryRepository) Delete(_ context.Context, ownerID, id string) (bool, error) {
	return r.store.Delete(ownerID, id)
}

// Exists implements the owner-scoped lookup required by the in-memory
// collection-assignment development surface. Cross-owner records remain
// indistinguishable from missing records through Store.Exists.
func (r *MemoryRepository) Exists(ownerID, id string) bool {
	return r.store.Exists(ownerID, id)
}
