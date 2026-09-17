package postgres

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestPostgreSQLMigrationsAndBookmarkStore(t *testing.T) {
	databaseURL := strings.TrimSpace(os.Getenv("GOREECLOUD_BOOKMARKS_TEST_DATABASE_URL"))
	if databaseURL == "" {
		t.Skip("GOREECLOUD_BOOKMARKS_TEST_DATABASE_URL is not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	database, err := Open(ctx, databaseURL)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	defer database.Close()

	var databaseName string
	if err := database.pool.QueryRow(ctx, "SELECT current_database()").Scan(&databaseName); err != nil {
		t.Fatalf("read current database: %v", err)
	}
	if !strings.HasSuffix(databaseName, "_test") {
		t.Fatalf("refusing destructive integration setup against database %q: expected *_test", databaseName)
	}

	if _, err := database.pool.Exec(ctx, "DROP TABLE IF EXISTS bookmarks"); err != nil {
		t.Fatalf("drop bookmarks table: %v", err)
	}
	if _, err := database.pool.Exec(ctx, "DROP TABLE IF EXISTS goreecloud_bookmarks_schema_migrations"); err != nil {
		t.Fatalf("drop migration history: %v", err)
	}

	if readiness := database.CheckReadiness(ctx); readiness.Ready || readiness.State != string(SchemaMigrationNeeded) {
		t.Fatalf("pre-migration readiness = %+v, want migration-required", readiness)
	}

	result, err := database.Migrate(ctx)
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if result.Applied != 1 || result.CurrentVersion != 1 {
		t.Fatalf("migration result = %+v, want one applied migration at version 1", result)
	}

	if readiness := database.CheckReadiness(ctx); !readiness.Ready || readiness.State != string(SchemaCurrent) {
		t.Fatalf("post-migration readiness = %+v, want schema-current", readiness)
	}

	bookmarkID := fmt.Sprintf("integration-%d", time.Now().UnixNano())
	created, err := database.CreateBookmark(ctx, Bookmark{
		ID:           bookmarkID,
		OwnerID:      "identity-test-owner",
		URL:          "https://example.invalid/bookmark",
		Title:        "Integration bookmark",
		PrivacyLevel: "normal",
	})
	if err != nil {
		t.Fatalf("create bookmark: %v", err)
	}
	if created.Revision != 1 || created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() {
		t.Fatalf("created bookmark has invalid revision/timestamps: %+v", created)
	}

	loaded, err := database.GetBookmark(ctx, "identity-test-owner", bookmarkID)
	if err != nil {
		t.Fatalf("get bookmark: %v", err)
	}
	if loaded.ID != bookmarkID || loaded.OwnerID != "identity-test-owner" || loaded.URL != created.URL {
		t.Fatalf("loaded bookmark = %+v, want persisted owner/id/url", loaded)
	}

	if _, err := database.GetBookmark(ctx, "different-owner", bookmarkID); err != ErrBookmarkNotFound {
		t.Fatalf("cross-owner read error = %v, want ErrBookmarkNotFound", err)
	}

	second, err := database.Migrate(ctx)
	if err != nil {
		t.Fatalf("second migrate: %v", err)
	}
	if second.Applied != 0 || second.CurrentVersion != 1 {
		t.Fatalf("second migration result = %+v, want idempotent no-op", second)
	}
}
