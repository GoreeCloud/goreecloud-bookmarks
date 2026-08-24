package postgres

import (
	"errors"
	"strings"
	"testing"
)

type fakeSQLStateError struct{ state string }

func (e fakeSQLStateError) Error() string    { return "database failure" }
func (e fakeSQLStateError) SQLState() string { return e.state }

func TestPersistentReadQueriesRequireOwnerScope(t *testing.T) {
	for name, query := range map[string]string{
		"list bookmarks": listBookmarksQuery,
		"get bookmark":   getBookmarkQuery,
		"list tags":      listTagsQuery,
	} {
		if !strings.Contains(query, "owner_id = $1") {
			t.Fatalf("%s must scope reads by authenticated owner: %q", name, query)
		}
	}
}

func TestPersistentWritesCarryOwnerIdentity(t *testing.T) {
	for name, query := range map[string]string{
		"bookmark insert": insertBookmarkQuery,
		"tag insert":      insertTagQuery,
	} {
		if !strings.Contains(query, "owner_id") {
			t.Fatalf("%s must persist owner identity: %q", name, query)
		}
	}
}

func TestUniqueViolationClassificationIsBounded(t *testing.T) {
	if !isUniqueViolation(fakeSQLStateError{state: uniqueViolationSQLState}) {
		t.Fatal("expected PostgreSQL unique violation classification")
	}
	if isUniqueViolation(fakeSQLStateError{state: "22001"}) {
		t.Fatal("non-unique SQL state must not be classified as duplicate")
	}
	if isUniqueViolation(errors.New("duplicate key value exposes database detail")) {
		t.Fatal("error text must not be parsed to classify persistence conflicts")
	}
}

func TestRepositoryRejectsNilDatabase(t *testing.T) {
	if _, err := NewRepository(nil); err == nil {
		t.Fatal("expected nil database rejection")
	}
}
