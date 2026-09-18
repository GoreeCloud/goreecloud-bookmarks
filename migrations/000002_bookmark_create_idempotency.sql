CREATE UNIQUE INDEX bookmarks_owner_id_id_unique
    ON bookmarks (owner_id, id);

CREATE TABLE bookmark_create_idempotency (
    owner_id text NOT NULL,
    idempotency_key text NOT NULL,
    request_hash text NOT NULL,
    bookmark_id text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (owner_id, idempotency_key),
    FOREIGN KEY (owner_id, bookmark_id)
        REFERENCES bookmarks (owner_id, id)
        ON DELETE CASCADE
        DEFERRABLE INITIALLY DEFERRED,
    CHECK (length(owner_id) > 0),
    CHECK (length(idempotency_key) BETWEEN 8 AND 200),
    CHECK (length(request_hash) = 64)
);

CREATE INDEX bookmark_create_idempotency_bookmark_id_idx
    ON bookmark_create_idempotency (bookmark_id);
