package postgres

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func openIntegrationDatabase(t *testing.T) (*Database, context.Context) {
	t.Helper()

	databaseURL := strings.TrimSpace(os.Getenv("GOREECLOUD_BOOKMARKS_TEST_DATABASE_URL"))
	if databaseURL == "" {
		t.Fatal("GOREECLOUD_BOOKMARKS_TEST_DATABASE_URL is required for PostgreSQL integration tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	database, err := Open(ctx, databaseURL)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(database.Close)

	var databaseName string
	if err := database.pool.QueryRow(ctx, "SELECT current_database()").Scan(&databaseName); err != nil {
		t.Fatalf("read current database: %v", err)
	}
	if !strings.HasSuffix(databaseName, "_test") {
		t.Fatalf("refusing destructive integration setup against database %q: expected *_test", databaseName)
	}

	resetIntegrationDatabase(t, ctx, database)
	return database, ctx
}

func resetIntegrationDatabase(t *testing.T, ctx context.Context, database *Database) {
	t.Helper()

	if _, err := database.pool.Exec(ctx, "DROP TABLE IF EXISTS bookmarks"); err != nil {
		t.Fatalf("drop bookmarks table: %v", err)
	}
	if _, err := database.pool.Exec(ctx, "DROP TABLE IF EXISTS goreecloud_bookmarks_schema_migrations"); err != nil {
		t.Fatalf("drop migration history: %v", err)
	}
}

func migrateIntegrationDatabase(t *testing.T, ctx context.Context, database *Database) {
	t.Helper()

	if readiness := database.CheckReadiness(ctx); readiness.Ready || readiness.State != string(SchemaMigrationNeeded) {
		t.Fatalf("pre-migration readiness = %+v, want migration-required", readiness)
	}
	if err := database.CheckStartupCompatibility(ctx); err != nil {
		t.Fatalf("migration-required schema should be startup-compatible: %v", err)
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
}

func TestPostgreSQLMigrationsAndBookmarkStore(t *testing.T) {
	database, ctx := openIntegrationDatabase(t)
	migrateIntegrationDatabase(t, ctx, database)

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

func TestPostgreSQLRejectsNewerSchema(t *testing.T) {
	database, ctx := openIntegrationDatabase(t)
	migrateIntegrationDatabase(t, ctx, database)

	if _, err := database.pool.Exec(ctx, `
INSERT INTO goreecloud_bookmarks_schema_migrations (version, name, checksum)
VALUES (2, 'future_migration', 'future-checksum')`); err != nil {
		t.Fatalf("insert future migration: %v", err)
	}

	if readiness := database.CheckReadiness(ctx); readiness.Ready || readiness.State != string(SchemaNewerThanBinary) {
		t.Fatalf("readiness = %+v, want schema-newer-than-binary", readiness)
	}
	if err := database.CheckStartupCompatibility(ctx); err == nil {
		t.Fatal("newer-than-binary schema was startup-compatible; want rejection")
	}
	if _, err := database.Migrate(ctx); err == nil {
		t.Fatal("Migrate() succeeded against schema newer than binary; want failure")
	}
}

func TestPostgreSQLRejectsMigrationHistoryMismatch(t *testing.T) {
	database, ctx := openIntegrationDatabase(t)
	migrateIntegrationDatabase(t, ctx, database)

	if _, err := database.pool.Exec(ctx, `
UPDATE goreecloud_bookmarks_schema_migrations
SET checksum = 'tampered'
WHERE version = 1`); err != nil {
		t.Fatalf("tamper migration history: %v", err)
	}

	if readiness := database.CheckReadiness(ctx); readiness.Ready || readiness.State != string(SchemaHistoryMismatch) {
		t.Fatalf("readiness = %+v, want migration-history-mismatch", readiness)
	}
	if err := database.CheckStartupCompatibility(ctx); err == nil {
		t.Fatal("tampered migration history was startup-compatible; want rejection")
	}
	if _, err := database.Migrate(ctx); err == nil {
		t.Fatal("Migrate() succeeded with tampered migration history; want failure")
	}
}
