package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/GoreeCloud/goreecloud-bookmarks/migrations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	historyTable     = "goreecloud_bookmarks_schema_migrations"
	migrationLockKey = int64(724436281)
)

// Database owns the PostgreSQL connection pool for the Bookmarks service.
type Database struct {
	pool *pgxpool.Pool
}

// Open parses the protected connection string and creates a pool. It does not
// apply migrations implicitly; schema mutation is performed only by Migrate.
func Open(ctx context.Context, databaseURL string) (*Database, error) {
	if strings.TrimSpace(databaseURL) == "" {
		return nil, errors.New("database URL is required")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		// Do not wrap parser output because malformed connection strings can
		// contain credentials or other protected configuration.
		return nil, errors.New("database configuration is invalid")
	}

	config.ConnConfig.ConnectTimeout = 5 * time.Second
	config.MinConns = 0
	config.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.New("database pool initialization failed")
	}

	return &Database{pool: pool}, nil
}

// Close releases all PostgreSQL pool resources.
func (d *Database) Close() {
	if d != nil && d.pool != nil {
		d.pool.Close()
	}
}

// Ping verifies that PostgreSQL is reachable through the configured pool.
func (d *Database) Ping(ctx context.Context) error {
	if d == nil || d.pool == nil {
		return errors.New("database is not initialized")
	}
	return d.pool.Ping(ctx)
}

// SchemaState describes the governed relationship between the database schema
// and the migration set embedded in this application revision.
type SchemaState string

const (
	SchemaCurrent          SchemaState = "schema-current"
	SchemaMigrationNeeded  SchemaState = "migration-required"
	SchemaNewerThanBinary  SchemaState = "schema-newer-than-binary"
	SchemaHistoryMismatch  SchemaState = "migration-history-mismatch"
	SchemaInvalid          SchemaState = "schema-invalid"
)

// ReadinessResult is deliberately bounded so operational endpoints never emit
// raw PostgreSQL errors, connection strings, credentials, or SQL details.
type ReadinessResult struct {
	Ready bool
	State string
}

// CheckReadiness verifies connectivity and exact migration compatibility.
func (d *Database) CheckReadiness(ctx context.Context) ReadinessResult {
	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := d.Ping(checkCtx); err != nil {
		return ReadinessResult{Ready: false, State: "database-unreachable"}
	}

	state, err := d.SchemaState(checkCtx)
	if err != nil {
		return ReadinessResult{Ready: false, State: string(SchemaInvalid)}
	}
	if state != SchemaCurrent {
		return ReadinessResult{Ready: false, State: string(state)}
	}

	return ReadinessResult{Ready: true, State: string(SchemaCurrent)}
}

// CheckStartupCompatibility rejects schema states that this binary must never
// operate against. A temporarily unreachable database is allowed to keep the
// process alive but not ready; newer, tampered, or structurally invalid schema
// state is an unsafe startup condition.
func (d *Database) CheckStartupCompatibility(ctx context.Context) error {
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := d.Ping(checkCtx); err != nil {
		return nil
	}

	state, err := d.SchemaState(checkCtx)
	if err != nil {
		return errors.New("database schema compatibility check failed")
	}

	switch state {
	case SchemaCurrent, SchemaMigrationNeeded:
		return nil
	case SchemaNewerThanBinary, SchemaHistoryMismatch, SchemaInvalid:
		return fmt.Errorf("unsafe database schema state: %s", state)
	default:
		return errors.New("unknown database schema state")
	}
}

// SchemaState returns the migration compatibility state without mutating the
// database.
func (d *Database) SchemaState(ctx context.Context) (SchemaState, error) {
	if d == nil || d.pool == nil {
		return SchemaInvalid, errors.New("database is not initialized")
	}
	return schemaState(ctx, d.pool)
}

// MigrationResult records the bounded outcome of one explicit migration run.
type MigrationResult struct {
	Applied        int
	CurrentVersion int64
}

// Migrate applies pending ordered migrations while holding a PostgreSQL
// advisory lock. Applied history is append-only and checksummed.
func (d *Database) Migrate(ctx context.Context) (MigrationResult, error) {
	if d == nil || d.pool == nil {
		return MigrationResult{}, errors.New("database is not initialized")
	}
	if err := d.pool.Ping(ctx); err != nil {
		return MigrationResult{}, errors.New("database is unreachable")
	}

	conn, err := d.pool.Acquire(ctx)
	if err != nil {
		return MigrationResult{}, errors.New("acquire migration connection failed")
	}
	defer conn.Release()

	if _, err := conn.Exec(ctx, "SELECT pg_advisory_lock($1)", migrationLockKey); err != nil {
		return MigrationResult{}, errors.New("acquire migration lock failed")
	}
	defer func() {
		unlockCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, _ = conn.Exec(unlockCtx, "SELECT pg_advisory_unlock($1)", migrationLockKey)
	}()

	if _, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS goreecloud_bookmarks_schema_migrations (
    version bigint PRIMARY KEY,
    name text NOT NULL,
    checksum text NOT NULL,
    applied_at timestamptz NOT NULL DEFAULT now()
)`); err != nil {
		return MigrationResult{}, errors.New("create migration history failed")
	}

	state, err := schemaState(ctx, conn)
	if err != nil {
		return MigrationResult{}, err
	}
	if state == SchemaNewerThanBinary || state == SchemaHistoryMismatch || state == SchemaInvalid {
		return MigrationResult{}, fmt.Errorf("unsafe schema state: %s", state)
	}

	result := MigrationResult{CurrentVersion: migrations.CurrentVersion()}

	for _, migration := range migrations.Ordered {
		checksum, err := migrations.Checksum(migration)
		if err != nil {
			return MigrationResult{}, err
		}

		var appliedName, appliedChecksum string
		err = conn.QueryRow(ctx,
			"SELECT name, checksum FROM "+historyTable+" WHERE version = $1",
			migration.Version,
		).Scan(&appliedName, &appliedChecksum)
		if err == nil {
			if appliedName != migration.Name || appliedChecksum != checksum {
				return MigrationResult{}, fmt.Errorf("migration history mismatch at version %d", migration.Version)
			}
			continue
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return MigrationResult{}, errors.New("read migration history failed")
		}

		sqlText, err := migrations.SQL(migration)
		if err != nil {
			return MigrationResult{}, err
		}

		tx, err := conn.Begin(ctx)
		if err != nil {
			return MigrationResult{}, errors.New("begin migration transaction failed")
		}

		if _, err := tx.Exec(ctx, sqlText); err != nil {
			_ = tx.Rollback(ctx)
			return MigrationResult{}, fmt.Errorf("migration %d failed", migration.Version)
		}
		if _, err := tx.Exec(ctx,
			"INSERT INTO "+historyTable+" (version, name, checksum) VALUES ($1, $2, $3)",
			migration.Version, migration.Name, checksum,
		); err != nil {
			_ = tx.Rollback(ctx)
			return MigrationResult{}, fmt.Errorf("record migration %d failed", migration.Version)
		}
		if err := tx.Commit(ctx); err != nil {
			return MigrationResult{}, fmt.Errorf("commit migration %d failed", migration.Version)
		}
		result.Applied++
	}

	finalState, err := schemaState(ctx, conn)
	if err != nil {
		return MigrationResult{}, err
	}
	if finalState != SchemaCurrent {
		return MigrationResult{}, fmt.Errorf("post-migration schema state is %s", finalState)
	}

	return result, nil
}

type queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func schemaState(ctx context.Context, q queryer) (SchemaState, error) {
	var historyExists bool
	if err := q.QueryRow(ctx,
		"SELECT to_regclass('public."+historyTable+"') IS NOT NULL",
	).Scan(&historyExists); err != nil {
		return SchemaInvalid, errors.New("inspect migration history failed")
	}
	if !historyExists {
		return SchemaMigrationNeeded, nil
	}

	expected := make(map[int64]struct {
		name     string
		checksum string
	}, len(migrations.Ordered))
	for _, migration := range migrations.Ordered {
		checksum, err := migrations.Checksum(migration)
		if err != nil {
			return SchemaInvalid, err
		}
		expected[migration.Version] = struct {
			name     string
			checksum string
		}{name: migration.Name, checksum: checksum}
	}

	rows, err := q.Query(ctx, "SELECT version, name, checksum FROM "+historyTable+" ORDER BY version")
	if err != nil {
		return SchemaInvalid, errors.New("read migration history failed")
	}
	defer rows.Close()

	seen := make(map[int64]bool, len(expected))
	for rows.Next() {
		var version int64
		var name, checksum string
		if err := rows.Scan(&version, &name, &checksum); err != nil {
			return SchemaInvalid, errors.New("scan migration history failed")
		}
		if version > migrations.CurrentVersion() {
			return SchemaNewerThanBinary, nil
		}
		exp, ok := expected[version]
		if !ok || exp.name != name || exp.checksum != checksum {
			return SchemaHistoryMismatch, nil
		}
		seen[version] = true
	}
	if err := rows.Err(); err != nil {
		return SchemaInvalid, errors.New("iterate migration history failed")
	}

	for _, migration := range migrations.Ordered {
		if !seen[migration.Version] {
			return SchemaMigrationNeeded, nil
		}
	}

	var bookmarksExists bool
	if err := q.QueryRow(ctx, "SELECT to_regclass('public.bookmarks') IS NOT NULL").Scan(&bookmarksExists); err != nil {
		return SchemaInvalid, errors.New("inspect bookmarks table failed")
	}
	if !bookmarksExists {
		return SchemaInvalid, nil
	}

	return SchemaCurrent, nil
}
