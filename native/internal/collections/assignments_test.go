package collections

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
)

type failingBookmarkLookup struct {
	err error
}

func (lookup failingBookmarkLookup) Exists(context.Context, string, string) (bool, error) {
	return false, lookup.err
}

func newAssignmentBookmarkService(t *testing.T) *bookmarks.Service {
	t.Helper()
	service, err := bookmarks.NewService(bookmarks.NewMemoryRepository(nil))
	if err != nil {
		t.Fatal(err)
	}
	return service
}

func TestAssignmentStoreEnforcesOwnerScopedBookmarkAndCollection(t *testing.T) {
	ctx := context.Background()
	bookmarkService := newAssignmentBookmarkService(t)
	collectionStore := NewStore()
	assignmentStore := NewAssignmentStore(collectionStore, bookmarkService)

	bookmark, err := bookmarkService.Create(ctx, "owner-a", bookmarks.CreateInput{URL: "https://example.com/a"})
	if err != nil {
		t.Fatal(err)
	}
	collection, err := collectionStore.Create("owner-a", "Research", "")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := assignmentStore.Assign(ctx, "owner-b", bookmark.ID, collection.ID); err != ErrBookmarkNotFound {
		t.Fatalf("cross-owner bookmark must be masked as not found, got %v", err)
	}
	otherCollection, err := collectionStore.Create("owner-b", "Research", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := assignmentStore.Assign(ctx, "owner-a", bookmark.ID, otherCollection.ID); err != ErrCollectionNotFound {
		t.Fatalf("cross-owner collection must be masked as not found, got %v", err)
	}

	assigned, err := assignmentStore.Assign(ctx, "owner-a", bookmark.ID, collection.ID)
	if err != nil {
		t.Fatal(err)
	}
	if assigned.OwnerID != "owner-a" || assigned.BookmarkID != bookmark.ID || assigned.CollectionID != collection.ID {
		t.Fatalf("unexpected assignment: %#v", assigned)
	}
	if assigned.AssignedAt.IsZero() {
		t.Fatal("assignment timestamp must be populated")
	}
}

func TestAssignmentStoreReassignsAndRemovesWithoutCrossOwnerLeakage(t *testing.T) {
	ctx := context.Background()
	bookmarkService := newAssignmentBookmarkService(t)
	collectionStore := NewStore()
	assignmentStore := NewAssignmentStore(collectionStore, bookmarkService)
	assignmentStore.now = func() time.Time { return time.Unix(123, 0).UTC() }

	bookmark, err := bookmarkService.Create(ctx, "owner-a", bookmarks.CreateInput{URL: "https://example.com/a"})
	if err != nil {
		t.Fatal(err)
	}
	first, _ := collectionStore.Create("owner-a", "First", "")
	second, _ := collectionStore.Create("owner-a", "Second", "")

	if _, err := assignmentStore.Assign(ctx, "owner-a", bookmark.ID, first.ID); err != nil {
		t.Fatal(err)
	}
	reassigned, err := assignmentStore.Assign(ctx, "owner-a", bookmark.ID, second.ID)
	if err != nil {
		t.Fatal(err)
	}
	if reassigned.CollectionID != second.ID {
		t.Fatalf("expected reassignment to %s, got %s", second.ID, reassigned.CollectionID)
	}
	if _, ok := assignmentStore.Get("owner-b", bookmark.ID); ok {
		t.Fatal("cross-owner assignment lookup must fail closed")
	}
	if assignmentStore.Remove("owner-b", bookmark.ID) {
		t.Fatal("cross-owner removal must not report success")
	}
	if !assignmentStore.Remove("owner-a", bookmark.ID) {
		t.Fatal("owner removal should succeed")
	}
	if _, ok := assignmentStore.Get("owner-a", bookmark.ID); ok {
		t.Fatal("assignment should be removed")
	}
}

func TestAssignmentStoreRejectsMissingIdentifiers(t *testing.T) {
	ctx := context.Background()
	store := NewAssignmentStore(NewStore(), newAssignmentBookmarkService(t))
	if _, err := store.Assign(ctx, "", "bookmark", "collection"); err != ErrOwnerRequired {
		t.Fatalf("expected owner required, got %v", err)
	}
	if _, err := store.Assign(ctx, "owner", "", "collection"); err != ErrBookmarkRequired {
		t.Fatalf("expected bookmark required, got %v", err)
	}
	if _, err := store.Assign(ctx, "owner", "bookmark", ""); err != ErrCollectionNotFound {
		t.Fatalf("expected collection not found, got %v", err)
	}
}

func TestAssignmentStorePropagatesBookmarkLookupFailure(t *testing.T) {
	lookupErr := errors.New("bookmark storage unavailable")
	store := NewAssignmentStore(NewStore(), failingBookmarkLookup{err: lookupErr})
	_, err := store.Assign(context.Background(), "owner", "bookmark", "collection")
	if !errors.Is(err, lookupErr) {
		t.Fatalf("bookmark lookup failure = %v, want %v", err, lookupErr)
	}
}
