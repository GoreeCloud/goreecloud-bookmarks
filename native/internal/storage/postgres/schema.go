package postgres

const SchemaVersion = 1

const Migration001 = `
CREATE TABLE IF NOT EXISTS native_bookmarks (
    id TEXT PRIMARY KEY,
    owner_id TEXT NOT NULL,
    canonical_url TEXT NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT native_bookmarks_owner_url_unique UNIQUE (owner_id, canonical_url),
    CONSTRAINT native_bookmarks_owner_nonempty CHECK (length(owner_id) > 0),
    CONSTRAINT native_bookmarks_url_nonempty CHECK (length(canonical_url) > 0)
);

CREATE INDEX IF NOT EXISTS native_bookmarks_owner_created_idx
    ON native_bookmarks (owner_id, created_at DESC, id DESC);

CREATE TABLE IF NOT EXISTS native_bookmark_tags (
    bookmark_id TEXT NOT NULL REFERENCES native_bookmarks(id) ON DELETE CASCADE,
    owner_id TEXT NOT NULL,
    tag TEXT NOT NULL,
    PRIMARY KEY (bookmark_id, tag),
    CONSTRAINT native_bookmark_tags_owner_nonempty CHECK (length(owner_id) > 0),
    CONSTRAINT native_bookmark_tags_tag_nonempty CHECK (length(tag) > 0)
);

CREATE INDEX IF NOT EXISTS native_bookmark_tags_owner_tag_idx
    ON native_bookmark_tags (owner_id, tag, bookmark_id);
`

// OwnerPredicate is deliberately explicit so every persistent read and mutation
// is scoped by authenticated owner identity. Storage adapters must not replace
// this with client-provided ownership state.
const OwnerPredicate = "owner_id = $1"
