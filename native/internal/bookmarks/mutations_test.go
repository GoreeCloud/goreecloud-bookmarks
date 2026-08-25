package bookmarks

import (
	"testing"
	"time"
)

func TestStoreUpdatePreservesOwnershipAndChangesFields(t *testing.T) {
	store := NewStore()
	createdAt := time.Date(2026, 8, 25, 18, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Minute)
	store.now = func() time.Time { return createdAt }
	created, err := store.Create("owner-a", CreateInput{URL: "https://example.com/one", Title: "One", Tags: []string{"Work"}})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	store.now = func() time.Time { return updatedAt }

	updated, found, err := store.Update("owner-a", created.ID, UpdateInput{
		URL: "https://example.com/two#fragment", Title: "  Two  ", Note: "  Note  ", Tags: []string{"work", "Read"},
	})
	if err != nil || !found {
		t.Fatalf("update found=%v err=%v", found, err)
	}
	if updated.URL != "https://example.com/two" || updated.Title != "Two" || updated.Note != "Note" {
		t.Fatalf("unexpected updated bookmark: %+v", updated)
	}
	if !updated.CreatedAt.Equal(createdAt) || !updated.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("unexpected timestamps: created=%s updated=%s", updated.CreatedAt, updated.UpdatedAt)
	}

	if _, found, err := store.Update("owner-b", created.ID, UpdateInput{URL: "https://example.com/three"}); err != nil || found {
		t.Fatalf("cross-owner update found=%v err=%v", found, err)
	}
}

func TestStoreUpdateRejectsOwnerDuplicateURL(t *testing.T) {
	store := NewStore()
	first, err := store.Create("owner-a", CreateInput{URL: "https://example.com/one"})
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	second, err := store.Create("owner-a", CreateInput{URL: "https://example.com/two"})
	if err != nil {
		t.Fatalf("create second: %v", err)
	}
	if _, found, err := store.Update("owner-a", second.ID, UpdateInput{URL: first.URL}); err == nil || found {
		t.Fatalf("expected duplicate rejection, found=%v err=%v", found, err)
	}
}

func TestStoreDeleteIsOwnerScoped(t *testing.T) {
	store := NewStore()
	created, err := store.Create("owner-a", CreateInput{URL: "https://example.com/delete"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if deleted, err := store.Delete("owner-b", created.ID); err != nil || deleted {
		t.Fatalf("cross-owner delete deleted=%v err=%v", deleted, err)
	}
	if _, ok := store.Get("owner-a", created.ID); !ok {
		t.Fatal("bookmark disappeared after rejected cross-owner delete")
	}
	if deleted, err := store.Delete("owner-a", created.ID); err != nil || !deleted {
		t.Fatalf("owner delete deleted=%v err=%v", deleted, err)
	}
	if _, ok := store.Get("owner-a", created.ID); ok {
		t.Fatal("bookmark remains after delete")
	}
}

func TestServiceMutationRequiresOwnerAndID(t *testing.T) {
	service, err := NewService(NewMemoryRepository(nil))
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	if _, _, err := service.Update(t.Context(), "", "bookmark", UpdateInput{URL: "https://example.com"}); err == nil {
		t.Fatal("expected owner validation error")
	}
	if _, err := service.Delete(t.Context(), "owner", " "); err == nil {
		t.Fatal("expected id validation error")
	}
}
