package collections

import (
	"errors"
	"testing"
	"time"
)

func TestCollectionsAreOwnerScoped(t *testing.T) {
	store := NewStore()
	collection, err := store.Create("owner-a", "Recipes", "")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if _, ok := store.Get("owner-b", collection.ID); ok {
		t.Fatal("cross-owner collection read must fail closed")
	}
	if got := store.List("owner-b"); len(got) != 0 {
		t.Fatalf("unexpected cross-owner list: %#v", got)
	}
}

func TestNestedCollectionRequiresSameOwnerParent(t *testing.T) {
	store := NewStore()
	parent, err := store.Create("owner-a", "Projects", "")
	if err != nil {
		t.Fatalf("Create parent returned error: %v", err)
	}
	child, err := store.Create("owner-a", "GoreeCloud", parent.ID)
	if err != nil {
		t.Fatalf("Create child returned error: %v", err)
	}
	if child.ParentID != parent.ID {
		t.Fatalf("child parent = %q, want %q", child.ParentID, parent.ID)
	}
	if _, err := store.Create("owner-b", "Private", parent.ID); !errors.Is(err, ErrParentNotFound) {
		t.Fatalf("cross-owner parent error = %v, want ErrParentNotFound", err)
	}
}

func TestSiblingNamesAreUniquePerOwnerAndParent(t *testing.T) {
	store := NewStore()
	if _, err := store.Create("owner-a", "Reading", ""); err != nil {
		t.Fatalf("first Create returned error: %v", err)
	}
	if _, err := store.Create("owner-a", " reading ", ""); !errors.Is(err, ErrDuplicateSibling) {
		t.Fatalf("duplicate error = %v, want ErrDuplicateSibling", err)
	}
	if _, err := store.Create("owner-b", "Reading", ""); err != nil {
		t.Fatalf("same name for another owner returned error: %v", err)
	}
}

func TestMoveRejectsHierarchyCyclesAndCrossOwnerParents(t *testing.T) {
	store := NewStore()
	root, _ := store.Create("owner-a", "Root", "")
	child, _ := store.Create("owner-a", "Child", root.ID)
	grandchild, _ := store.Create("owner-a", "Grandchild", child.ID)
	otherOwner, _ := store.Create("owner-b", "Other", "")

	if _, err := store.Move("owner-a", root.ID, grandchild.ID); !errors.Is(err, ErrHierarchyCycle) {
		t.Fatalf("cycle move error = %v, want ErrHierarchyCycle", err)
	}
	if _, err := store.Move("owner-a", child.ID, otherOwner.ID); !errors.Is(err, ErrParentNotFound) {
		t.Fatalf("cross-owner move error = %v, want ErrParentNotFound", err)
	}
}

func TestMovePreservesCreationAndAdvancesUpdate(t *testing.T) {
	store := NewStore()
	moments := []time.Time{
		time.Date(2026, 8, 25, 20, 0, 0, 0, time.UTC),
		time.Date(2026, 8, 25, 20, 1, 0, 0, time.UTC),
		time.Date(2026, 8, 25, 20, 2, 0, 0, time.UTC),
	}
	index := 0
	store.now = func() time.Time {
		value := moments[index]
		index++
		return value
	}
	collection, err := store.Create("owner-a", "Reading", "")
	if err != nil {
		t.Fatalf("Create collection returned error: %v", err)
	}
	parent, err := store.Create("owner-a", "Library", "")
	if err != nil {
		t.Fatalf("Create parent returned error: %v", err)
	}
	moved, err := store.Move("owner-a", collection.ID, parent.ID)
	if err != nil {
		t.Fatalf("Move returned error: %v", err)
	}
	if !moved.CreatedAt.Equal(collection.CreatedAt) {
		t.Fatalf("CreatedAt changed from %v to %v", collection.CreatedAt, moved.CreatedAt)
	}
	if !moved.UpdatedAt.After(collection.UpdatedAt) {
		t.Fatalf("UpdatedAt did not advance: before=%v after=%v", collection.UpdatedAt, moved.UpdatedAt)
	}
}
