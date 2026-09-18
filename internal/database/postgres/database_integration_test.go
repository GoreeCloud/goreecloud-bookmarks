package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/GoreeCloud/goreecloud-bookmarks/internal/bookmarks"
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

	if _, err := database.pool.Exec(ctx, "DROP TABLE IF EXISTS bookmark_create_idempotency"); err != nil {
		t.Fatalf("drop bookmark idempotency table: %v", err)
	}
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
	if result.Applied != 2 || result.CurrentVersion != 2 {
		t.Fatalf("migration result = %+v, want two applied migrations at version 2", result)
	}

	if readiness := database.CheckReadiness(ctx); !readiness.Ready || readiness.State != string(SchemaCurrent) {
		t.Fatalf("post-migration readiness = %+v, want schema-current", readiness)
	}
}

func TestPostgreSQLMigrationsAndIdempotentBookmarkStore(t *testing.T) {
	database, ctx := openIntegrationDatabase(t)
	migrateIntegrationDatabase(t, ctx, database)

	firstID := fmt.Sprintf("integration-%d", time.Now().UnixNano())
	firstInput := bookmarks.Bookmark{
		ID:           firstID,
		OwnerID:      "identity-test-owner",
		URL:          "https://example.invalid/bookmark",
		Title:        "Integration bookmark",
		ReadState:    "unread",
		PrivacyLevel: "normal",
	}
	requestHash := strings.Repeat("a", 64)

	created, replayed, err := database.CreateBookmarkIdempotent(
		ctx,
		"identity-test-owner",
		"request-key-1",
		requestHash,
		firstInput,
	)
	if err != nil {
		t.Fatalf("create bookmark: %v", err)
	}
	if replayed {
		t.Fatal("first create reported replay")
	}
	if created.Revision != 1 || created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() {
		t.Fatalf("created bookmark has invalid revision/timestamps: %+v", created)
	}

	loaded, err := database.GetBookmark(ctx, "identity-test-owner", firstID)
	if err != nil {
		t.Fatalf("get bookmark: %v", err)
	}
	if loaded.ID != firstID || loaded.OwnerID != "identity-test-owner" || loaded.URL != created.URL {
		t.Fatalf("loaded bookmark = %+v, want persisted owner/id/url", loaded)
	}

	if _, err := database.GetBookmark(ctx, "different-owner", firstID); !errors.Is(err, bookmarks.ErrNotFound) {
		t.Fatalf("cross-owner read error = %v, want ErrNotFound", err)
	}

	retryInput := firstInput
	retryInput.ID = "retry-generated-but-unused"
	replayedBookmark, replayed, err := database.CreateBookmarkIdempotent(
		ctx,
		"identity-test-owner",
		"request-key-1",
		requestHash,
		retryInput,
	)
	if err != nil {
		t.Fatalf("replay bookmark: %v", err)
	}
	if !replayed || replayedBookmark.ID != firstID {
		t.Fatalf("replay = (%+v, %v), want original bookmark %q", replayedBookmark, replayed, firstID)
	}

	var bookmarkCount, idempotencyCount int
	if err := database.pool.QueryRow(ctx, "SELECT count(*) FROM bookmarks").Scan(&bookmarkCount); err != nil {
		t.Fatalf("count bookmarks: %v", err)
	}
	if err := database.pool.QueryRow(ctx, "SELECT count(*) FROM bookmark_create_idempotency").Scan(&idempotencyCount); err != nil {
		t.Fatalf("count bookmark idempotency rows: %v", err)
	}
	if bookmarkCount != 1 || idempotencyCount != 1 {
		t.Fatalf("replay created duplicates: bookmarks=%d idempotency=%d", bookmarkCount, idempotencyCount)
	}

	_, _, err = database.CreateBookmarkIdempotent(
		ctx,
		"identity-test-owner",
		"request-key-1",
		strings.Repeat("b", 64),
		retryInput,
	)
	if !errors.Is(err, bookmarks.ErrIdempotencyConflict) {
		t.Fatalf("conflicting replay error = %v, want ErrIdempotencyConflict", err)
	}

	otherOwnerInput := firstInput
	otherOwnerInput.ID = "other-owner-bookmark"
	otherOwnerInput.OwnerID = "other-owner"
	other, otherReplay, err := database.CreateBookmarkIdempotent(
		ctx,
		"other-owner",
		"request-key-1",
		requestHash,
		otherOwnerInput,
	)
	if err != nil {
		t.Fatalf("same key for different owner: %v", err)
	}
	if otherReplay || other.ID != "other-owner-bookmark" {
		t.Fatalf("different owner create = (%+v, %v), want independent create", other, otherReplay)
	}

	second, err := database.Migrate(ctx)
	if err != nil {
		t.Fatalf("second migrate: %v", err)
	}
	if second.Applied != 0 || second.CurrentVersion != 2 {
		t.Fatalf("second migration result = %+v, want idempotent no-op at version 2", second)
	}
}

func TestPostgreSQLConcurrentIdempotentBookmarkCreate(t *testing.T) {
	database, ctx := openIntegrationDatabase(t)
	migrateIntegrationDatabase(t, ctx, database)

	type createResult struct {
		bookmark bookmarks.Bookmark
		replayed bool
		err      error
	}

	start := make(chan struct{})
	results := make(chan createResult, 2)
	var wg sync.WaitGroup

	for _, bookmarkID := range []string{"concurrent-bookmark-a", "concurrent-bookmark-b"} {
		bookmarkID := bookmarkID
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start

			bookmark, replayed, err := database.CreateBookmarkIdempotent(
				ctx,
				"concurrent-owner",
				"concurrent-request-key",
				strings.Repeat("c", 64),
				bookmarks.Bookmark{
					ID:           bookmarkID,
					OwnerID:      "concurrent-owner",
					URL:          "https://example.invalid/concurrent",
					Title:        "Concurrent bookmark",
					ReadState:    "unread",
					PrivacyLevel: "normal",
				},
			)
			results <- createResult{bookmark: bookmark, replayed: replayed, err: err}
		}()
	}

	close(start)
	wg.Wait()
	close(results)

	var firstID string
	var createdCount, replayCount int
	for result := range results {
		if result.err != nil {
			t.Fatalf("concurrent create: %v", result.err)
		}
		if firstID == "" {
			firstID = result.bookmark.ID
		} else if result.bookmark.ID != firstID {
			t.Fatalf("concurrent retries returned different bookmark IDs: %q vs %q", firstID, result.bookmark.ID)
		}
		if result.replayed {
			replayCount++
		} else {
			createdCount++
		}
	}
	if createdCount != 1 || replayCount != 1 {
		t.Fatalf("concurrent outcomes: created=%d replayed=%d, want 1/1", createdCount, replayCount)
	}

	var bookmarkCount, idempotencyCount int
	if err := database.pool.QueryRow(ctx, "SELECT count(*) FROM bookmarks").Scan(&bookmarkCount); err != nil {
		t.Fatalf("count bookmarks: %v", err)
	}
	if err := database.pool.QueryRow(ctx, "SELECT count(*) FROM bookmark_create_idempotency").Scan(&idempotencyCount); err != nil {
		t.Fatalf("count idempotency rows: %v", err)
	}
	if bookmarkCount != 1 || idempotencyCount != 1 {
		t.Fatalf("concurrent retries created duplicates: bookmarks=%d idempotency=%d", bookmarkCount, idempotencyCount)
	}
}

func TestPostgreSQLRejectsNewerSchema(t *testing.T) {
	database, ctx := openIntegrationDatabase(t)
	migrateIntegrationDatabase(t, ctx, database)

	if _, err := database.pool.Exec(ctx, `
INSERT INTO goreecloud_bookmarks_schema_migrations (version, name, checksum)
VALUES (3, 'future_migration', 'future-checksum')`); err != nil {
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

func TestPostgreSQLRejectsMissingRequiredSchemaObject(t *testing.T) {
	database, ctx := openIntegrationDatabase(t)
	migrateIntegrationDatabase(t, ctx, database)

	if _, err := database.pool.Exec(ctx, "DROP TABLE bookmark_create_idempotency"); err != nil {
		t.Fatalf("drop required idempotency table: %v", err)
	}

	if readiness := database.CheckReadiness(ctx); readiness.Ready || readiness.State != string(SchemaInvalid) {
		t.Fatalf("readiness = %+v, want schema-invalid", readiness)
	}
	if err := database.CheckStartupCompatibility(ctx); err == nil {
		t.Fatal("missing required schema object was startup-compatible; want rejection")
	}
}
