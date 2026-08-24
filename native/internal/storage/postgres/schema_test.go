package postgres

import (
	"strings"
	"testing"
)

func TestMigration001PreservesOwnerScopedUniqueness(t *testing.T) {
	for _, required := range []string{
		"UNIQUE (owner_id, canonical_url)",
		"native_bookmarks_owner_created_idx",
		"native_bookmark_tags_owner_tag_idx",
		"REFERENCES native_bookmarks(id) ON DELETE CASCADE",
	} {
		if !strings.Contains(Migration001, required) {
			t.Fatalf("migration missing required ownership/storage clause %q", required)
		}
	}
}

func TestPersistentOwnerPredicateFailsClosed(t *testing.T) {
	if OwnerPredicate != "owner_id = $1" {
		t.Fatalf("owner predicate changed unexpectedly: %q", OwnerPredicate)
	}
}

func TestSchemaVersionIsExplicit(t *testing.T) {
	if SchemaVersion != 1 {
		t.Fatalf("unexpected schema version: %d", SchemaVersion)
	}
}
