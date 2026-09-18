package migrations

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
)

// FS contains the ordered PostgreSQL schema migrations shipped with the exact
// application revision. Applied migration files are append-only history and
// must not be rewritten after release/deployment.
//
//go:embed *.sql
var FS embed.FS

// Migration identifies one ordered schema migration.
type Migration struct {
	Version int64
	Name    string
	File    string
}

// Ordered is the authoritative migration order for this application revision.
var Ordered = []Migration{
	{Version: 1, Name: "initial_bookmarks", File: "000001_initial_bookmarks.sql"},
	{Version: 2, Name: "bookmark_create_idempotency", File: "000002_bookmark_create_idempotency.sql"},
}

// CurrentVersion returns the highest schema version understood by this binary.
func CurrentVersion() int64 {
	if len(Ordered) == 0 {
		return 0
	}
	return Ordered[len(Ordered)-1].Version
}

// SQL returns the exact SQL bytes represented by one embedded migration.
func SQL(m Migration) (string, error) {
	content, err := FS.ReadFile(m.File)
	if err != nil {
		return "", fmt.Errorf("read migration %s: %w", m.File, err)
	}
	return string(content), nil
}

// Checksum returns the SHA-256 digest of the exact migration SQL bytes.
func Checksum(m Migration) (string, error) {
	content, err := FS.ReadFile(m.File)
	if err != nil {
		return "", fmt.Errorf("read migration %s: %w", m.File, err)
	}
	digest := sha256.Sum256(content)
	return hex.EncodeToString(digest[:]), nil
}
