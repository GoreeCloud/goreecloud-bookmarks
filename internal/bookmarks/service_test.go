package bookmarks

import (
	"context"
	"errors"
	"testing"
)

type fakeStore struct {
	create func(context.Context, string, string, string, Bookmark) (Bookmark, bool, error)
	get    func(context.Context, string, string) (Bookmark, error)
}

func (f *fakeStore) CreateBookmarkIdempotent(
	ctx context.Context,
	ownerID, key, requestHash string,
	bookmark Bookmark,
) (Bookmark, bool, error) {
	return f.create(ctx, ownerID, key, requestHash, bookmark)
}

func (f *fakeStore) GetBookmark(ctx context.Context, ownerID, bookmarkID string) (Bookmark, error) {
	if f.get == nil {
		return Bookmark{}, ErrNotFound
	}
	return f.get(ctx, ownerID, bookmarkID)
}

func TestCreateNormalizesAndPersistsBeforeEnrichment(t *testing.T) {
	var captured Bookmark
	var capturedOwner, capturedKey, capturedHash string
	store := &fakeStore{create: func(
		_ context.Context,
		ownerID, key, requestHash string,
		bookmark Bookmark,
	) (Bookmark, bool, error) {
		capturedOwner, capturedKey, capturedHash, captured = ownerID, key, requestHash, bookmark
		bookmark.Revision = 1
		return bookmark, false, nil
	}}
	service, err := newService(store, func() (string, error) { return "opaque-bookmark-id", nil })
	if err != nil {
		t.Fatalf("newService: %v", err)
	}

	result, err := service.Create(context.Background(), " owner-1 ", " request-key-1 ", CreateInput{
		URL: "https://example.invalid/article",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if result.Replayed {
		t.Fatal("new create reported replay")
	}
	if result.Bookmark.ID != "opaque-bookmark-id" || result.Bookmark.Revision != 1 {
		t.Fatalf("unexpected create result: %+v", result)
	}
	if capturedOwner != "owner-1" || capturedKey != "request-key-1" {
		t.Fatalf("owner/key not normalized: owner=%q key=%q", capturedOwner, capturedKey)
	}
	if len(capturedHash) != 64 {
		t.Fatalf("request hash length = %d, want 64", len(capturedHash))
	}
	if captured.Title != captured.URL || captured.ReadState != "unread" || captured.PrivacyLevel != "normal" {
		t.Fatalf("defaults not normalized before persistence: %+v", captured)
	}
}

func TestCreateRejectsUnsupportedFieldsInsteadOfDroppingThem(t *testing.T) {
	store := &fakeStore{create: func(
		context.Context, string, string, string, Bookmark,
	) (Bookmark, bool, error) {
		t.Fatal("store should not be called")
		return Bookmark{}, false, nil
	}}
	service, _ := newService(store, func() (string, error) { return "unused", nil })

	_, err := service.Create(context.Background(), "owner-1", "request-key-1", CreateInput{
		URL:           "https://example.invalid/",
		CollectionIDs: []string{"collection-1"},
	})
	if err == nil {
		t.Fatal("Create succeeded with unimplemented collection behavior")
	}
}

func TestCreateRejectsInvalidBoundaryValues(t *testing.T) {
	store := &fakeStore{create: func(
		context.Context, string, string, string, Bookmark,
	) (Bookmark, bool, error) {
		t.Fatal("store should not be called")
		return Bookmark{}, false, nil
	}}
	service, _ := newService(store, func() (string, error) { return "unused", nil })

	tests := []struct {
		name  string
		owner string
		key   string
		input CreateInput
	}{
		{name: "missing owner", key: "request-key-1", input: CreateInput{URL: "https://example.invalid/"}},
		{name: "short key", owner: "owner-1", key: "short", input: CreateInput{URL: "https://example.invalid/"}},
		{name: "relative URL", owner: "owner-1", key: "request-key-1", input: CreateInput{URL: "/relative"}},
		{name: "bad read state", owner: "owner-1", key: "request-key-1", input: CreateInput{URL: "https://example.invalid/", ReadState: "later"}},
		{name: "bad privacy", owner: "owner-1", key: "request-key-1", input: CreateInput{URL: "https://example.invalid/", PrivacyLevel: "public-by-default"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := service.Create(context.Background(), tc.owner, tc.key, tc.input); err == nil {
				t.Fatal("Create succeeded; want validation failure")
			}
		})
	}
}

func TestCreatePreservesStoreIdempotencyConflict(t *testing.T) {
	store := &fakeStore{create: func(
		context.Context, string, string, string, Bookmark,
	) (Bookmark, bool, error) {
		return Bookmark{}, false, ErrIdempotencyConflict
	}}
	service, _ := newService(store, func() (string, error) { return "opaque-id", nil })

	_, err := service.Create(context.Background(), "owner-1", "request-key-1", CreateInput{
		URL: "https://example.invalid/",
	})
	if !errors.Is(err, ErrIdempotencyConflict) {
		t.Fatalf("Create error = %v, want ErrIdempotencyConflict", err)
	}
}

func TestNewOpaqueIDProducesOpaqueValue(t *testing.T) {
	first, err := newOpaqueID()
	if err != nil {
		t.Fatalf("newOpaqueID: %v", err)
	}
	second, err := newOpaqueID()
	if err != nil {
		t.Fatalf("newOpaqueID second: %v", err)
	}
	if first == "" || second == "" || first == second {
		t.Fatalf("opaque IDs not unique/non-empty: %q %q", first, second)
	}
}
