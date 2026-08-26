package syncstate

import (
	"testing"
	"time"

	bookmarkcore "github.com/GoreeCloud/goreecloud-bookmarks/native/internal/bookmarks"
)

func TestExportBookmarkItemOmitsOwnerIdentity(t *testing.T) {
	when := time.Unix(300, 0).UTC()
	item, err := ExportBookmarkItem(bookmarkcore.Bookmark{
		ID: "bookmark-1", OwnerID: "owner-secret", URL: "https://example.com",
		Title: "Example", Note: "note", Tags: []string{"one", "two"}, CreatedAt: when, UpdatedAt: when,
	})
	if err != nil {
		t.Fatalf("ExportBookmarkItem: %v", err)
	}
	if item.ID != "bookmark-1" || item.URL != "https://example.com" || len(item.Tags) != 2 {
		t.Fatalf("unexpected item: %+v", item)
	}
	item.Tags[0] = "changed"
	if item.Tags[0] != "changed" {
		t.Fatal("expected isolated tag copy")
	}
}
