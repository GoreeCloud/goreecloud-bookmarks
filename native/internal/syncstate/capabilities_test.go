package syncstate

import "testing"

func TestCapabilitiesExposeStableBookmarksDatasets(t *testing.T) {
	got := Capabilities()
	if len(got) != 3 {
		t.Fatalf("capability count = %d, want 3", len(got))
	}
	want := []string{"bookmarks.items", "bookmarks.collections", "bookmarks.assignments"}
	for i, dataset := range want {
		if got[i].Dataset != dataset {
			t.Fatalf("capability[%d] = %q, want %q", i, got[i].Dataset, dataset)
		}
		if got[i].Application != "bookmarks" || got[i].SchemaVersion != 1 {
			t.Fatalf("unexpected capability metadata: %+v", got[i])
		}
		if !got[i].Read || !got[i].Write || !got[i].Delete {
			t.Fatalf("unexpected dataset permissions: %+v", got[i])
		}
	}

	got[0].Dataset = "mutated"
	if Capabilities()[0].Dataset != "bookmarks.items" {
		t.Fatal("Capabilities returned shared mutable storage")
	}
}
