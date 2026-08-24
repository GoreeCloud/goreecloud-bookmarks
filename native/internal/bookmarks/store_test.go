package bookmarks

import (
	"strings"
	"testing"
)

func TestStoreEnforcesOwnerIsolation(t *testing.T) {
	store := NewStore()
	created, err := store.Create("alice", CreateInput{URL: "https://example.com/page#section", Title: " Example ", Tags: []string{"Research", "research", " Useful "}})
	if err != nil {
		t.Fatal(err)
	}
	if created.URL != "https://example.com/page" {
		t.Fatalf("expected fragment-free URL, got %q", created.URL)
	}
	if created.Title != "Example" {
		t.Fatalf("expected trimmed title, got %q", created.Title)
	}
	if len(created.Tags) != 2 {
		t.Fatalf("expected case-insensitive tag de-duplication, got %#v", created.Tags)
	}
	if _, ok := store.Get("bob", created.ID); ok {
		t.Fatal("cross-owner read must fail closed")
	}
	if got := store.List("bob"); len(got) != 0 {
		t.Fatalf("cross-owner list leaked %d bookmarks", len(got))
	}
	if got := store.List("alice"); len(got) != 1 {
		t.Fatalf("owner list expected 1 bookmark, got %d", len(got))
	}
}

func TestStoreRejectsDuplicateForSameOwnerButAllowsDifferentOwner(t *testing.T) {
	store := NewStore()
	if _, err := store.Create("alice", CreateInput{URL: "https://example.com"}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Create("alice", CreateInput{URL: "https://example.com"}); err == nil {
		t.Fatal("expected same-owner duplicate rejection")
	}
	if _, err := store.Create("bob", CreateInput{URL: "https://example.com"}); err != nil {
		t.Fatalf("different owner should be independent: %v", err)
	}
}

func TestStoreRejectsUnsafeAndOversizedInputs(t *testing.T) {
	store := NewStore()
	cases := []CreateInput{
		{URL: "javascript:alert(1)"},
		{URL: "https://example.com", Title: strings.Repeat("x", MaxTitleRunes+1)},
		{URL: "https://example.com", Note: strings.Repeat("x", MaxNoteRunes+1)},
		{URL: "https://example.com", Tags: []string{strings.Repeat("x", MaxTagRunes+1)}},
	}
	for i, input := range cases {
		if _, err := store.Create("alice", input); err == nil {
			t.Fatalf("case %d: expected validation failure", i)
		}
	}
}
