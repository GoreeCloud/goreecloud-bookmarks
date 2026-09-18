# GoreeCloud Bookmarks — Versioning, Migration, Compatibility, and Rollback

## Status

- **Lifecycle:** Development
- **Policy status:** Migration framework implemented for PostgreSQL schema version `1`; destructive migration recovery and production rollback remain unqualified.
- **API major:** `v1`
- **Server database direction:** PostgreSQL
- **Installed-client database direction:** SQLite
- **Complete archive container:** WARC 1.1 / ISO 28500:2017

This document defines how Bookmarks data and interfaces are allowed to evolve. The repository now implements an explicit PostgreSQL migration runner and schema version `1`; it does not establish production deployment, backup/restore qualification, client migrations, or release acceptance.

## 1. Separate version domains

The following versions are distinct and must not be conflated:

1. **Application version** — the GoreeCloud Bookmarks release version.
2. **API major version** — the external application contract path such as `/api/v1/`.
3. **PostgreSQL schema version** — ordered server-database migration state.
4. **SQLite client schema version** — local per-client database migration state.
5. **Archive representation version** — WARC/derived-representation format metadata.
6. **Synchronization protocol/schema version** — mutation/change-feed contract within the API major where separately needed.
7. **Platform Contract schema version** — GoreeCloud platform declaration/conformance metadata.

A change in one version domain does not automatically require a change in every other domain.

## 2. API compatibility

The initial API major is `/api/v1/`.

Within v1, compatible changes may include:

- adding new optional response fields;
- adding new optional request fields;
- adding new endpoints/resources;
- adding filters/sorts whose absence preserves existing behavior;
- adding new documented error codes for newly introduced operations;
- extending explicitly extensible metadata structures when clients are required to ignore unknown fields safely.

Breaking changes include:

- removing or renaming an existing field or endpoint;
- changing an existing field to an incompatible type;
- making an optional request field required without an alternate compatible path;
- changing an operation's meaning so an existing conforming client would perform a materially different action;
- changing privacy/authorization semantics in a way that invalidates existing client assumptions;
- silently changing pagination/cursor semantics so stored cursors become unsafe;
- changing conflict or deletion semantics in a way that could lose user data.

Breaking changes require a new API major such as `/api/v2/` unless the old behavior was never part of a released/accepted contract and governance explicitly permits correction before release.

## 3. API deprecation

When a future API major supersedes v1:

- the prior major remains documented while supported;
- a deprecation notice identifies the replacement and affected clients/integrations;
- supported first-party clients must have a verified upgrade path before retirement;
- a retirement decision must account for offline clients that may reconnect after extended periods;
- removal requires authoritative release/change evidence rather than a documentation-only declaration.

No fixed calendar deprecation interval is invented before Bookmarks has an established release cadence. Each retirement must record its actual support window and evidence.

## 4. PostgreSQL schema migrations

Server schema changes use ordered, version-controlled migration files committed with the source revision that requires them. The current implementation embeds the ordered SQL set and records applied version/name/SHA-256 checksum in `goreecloud_bookmarks_schema_migrations`.

Rules:

- applied migration history is append-only;
- an already applied production migration must not be rewritten to make history look cleaner;
- schema state must be queryable through a migration-history table or equivalent durable mechanism;
- the server must refuse unsafe startup when the database schema is newer than the running binary understands;
- database changes must be reproducible from a fresh supported database plus the ordered migration set;
- direct manual production DDL is not a normal deployment mechanism and, if emergency use is ever required, must be reconciled back into authoritative migration history.

## 5. Expand → migrate → contract pattern

Potentially breaking relational changes should use an expand/migrate/contract sequence when practical:

### Expand

Add the new table/column/index/relationship in a form that does not immediately break the currently supported application revision.

### Migrate

Backfill or transform data with bounded, resumable, observable operations where dataset size or failure risk warrants it.

### Contract

Remove obsolete structures only after the supported application/client path no longer requires them and rollback/recovery consequences have been explicitly evaluated.

For large libraries, data backfills must not assume they can complete inside one deployment transaction.

## 6. Destructive migration gate

Before a destructive or irreversible server migration:

- exact candidate application and migration revision must be identified;
- affected data classes and approximate scope must be known;
- a database-supported backup/recovery point must exist and be validated according to current recovery governance;
- archive/file relationships affected by the change must be included in the recovery plan;
- migration validation must run against representative test data;
- interruption/retry behavior must be understood;
- expected rollback path must be documented before execution;
- the migration must fail closed rather than continuing after integrity checks fail.

A copied live PostgreSQL data directory is not, by itself, acceptable proof of a recoverable migration backup.

## 7. Database rollback model

Rollback does **not** assume that every migration has a safe automatic `down` migration.

Rollback options, in order of safety for the actual change, may include:

1. roll back application code while leaving a backward-compatible expanded schema in place;
2. apply a forward corrective migration;
3. restore from a verified pre-migration database backup/recovery point;
4. restore both relational and archive/file state when the migration changed relationships across those stores.

A destructive downgrade that would discard newer user data is prohibited merely to return to an older binary.

## 8. Application/database compatibility gate

Every server release that changes database requirements must declare the schema versions it can start against.

The startup path must distinguish at least:

- schema current and supported;
- schema older and automatically migratable by the governed deployment process;
- schema newer than the binary and therefore unsafe;
- migration failed/incomplete and therefore not ready.

Readiness remains false when migrations are absent/incomplete, when the database is unreachable, when applied history is newer than the binary, when recorded name/checksum differs from the embedded migration, or when required schema objects are structurally absent.

### Current implemented PostgreSQL migration boundary

- Schema version `1` creates the initial `bookmarks` relation; schema version `2` adds owner-scoped durable create-idempotency state and the composite owner/bookmark reference used for atomic retry-safe capture.
- `cmd/bookmarks-migrate` explicitly applies pending migrations; normal service startup never auto-mutates schema.
- A PostgreSQL advisory lock serializes migration execution.
- Each migration is transactional and is recorded only after its SQL succeeds.
- Applied name/checksum mismatch and newer-than-binary history fail closed.
- CI integration tests validate first application through version `2`, idempotent migration re-run, current readiness, newer-schema rejection, tampered-history rejection, required-schema-object validation, sequential replay/conflict behavior, owner-scoped key independence, and concurrent retry deduplication against PostgreSQL `18.6`.
- No destructive migration exists yet, so restore-backed destructive migration qualification remains future work.

## 9. SQLite client migrations

Installed clients maintain their own SQLite schema version.

Client migrations must be transactional where the platform/database operation permits it and must preserve unsynchronized user work.

Local data is divided into two recovery classes:

### Rebuildable cache

Data that can be safely re-fetched from the authorized server may be discarded/rebuilt after integrity failure when doing so does not erase user intent or violate offline expectations.

### Non-disposable local state

Pending mutations, unsynchronized notes/highlights/edits, locally created objects awaiting server IDs, conflict copies, and other user-authored unsynchronized state must not be deleted as a cache-reset shortcut.

Before a client schema reset/rebuild, non-disposable state must be migrated, exported to a recoverable local form, or explicitly surfaced to the user if automated preservation is impossible.

## 10. Client compatibility with API versions

Each shipped first-party client must declare the API major(s) it understands.

A server must not silently serve an incompatible contract to an older client. The client/server handshake or normal API response path must make incompatibility explicit and provide a safe upgrade requirement rather than allowing undefined synchronization behavior.

Offline clients that reconnect after an extended period must either:

- synchronize through a still-supported API major; or
- perform an explicit upgrade/migration flow before submitting queued writes.

Queued writes must never be blindly replayed against an incompatible API major.

## 11. Synchronization revision rules

Synchronization-relevant resources use stable IDs and monotonically increasing logical revisions within each resource's history.

Mutations include a `baseRevision` when updating existing synchronized state.

Server outcomes distinguish:

- accepted mutation with resulting revision;
- idempotent replay of an already accepted mutation;
- stale/conflicting mutation;
- validation failure;
- authorization failure;
- permanently unsupported mutation due to protocol/API incompatibility.

A client mutation ID remains stable across retries of the same logical mutation.

## 12. Conflict and rollback preservation

User-authored content receives preservation-first handling.

If a concurrent edit cannot be merged safely:

- the server must not silently discard either meaningful variant;
- the conflict response must identify the authoritative current revision and preserve enough submitted data for client recovery;
- a later merge/resolve action creates a new revision rather than rewriting historical evidence invisibly.

Last-write-wins is permitted only for explicitly low-risk fields where loss of the previous value has been accepted by design.

## 13. Deletion/tombstone compatibility

Deleted synchronized resources emit tombstones/change records.

Tombstones may be physically purged only when the applicable retention policy, device retirement/resume policy, backup/recovery model, and privacy requirements allow it.

A server must not purge tombstones in a way that causes a legitimately resuming client to resurrect deleted content without an explicit conflict/reconciliation path.

## 14. Archive representation migrations

Historical complete captures are immutable preservation artifacts.

Rules:

- existing WARC capture bytes are not rewritten in place merely because a newer writer/library becomes available;
- integrity digests remain associated with the exact stored bytes they validate;
- a new conversion, normalization, compression form, or replay artifact is stored as a new derived representation/version with provenance back to the source capture;
- archive metadata schema may migrate in PostgreSQL while preserving stable Archive/ArchiveVersion IDs and object references;
- storage-key moves require integrity verification before the old object is removed;
- bulk archive relocation must be resumable and auditable.

## 15. Search/index migrations

Search indexes are derived state.

The initial PostgreSQL full-text index can be rebuilt from authorized authoritative metadata and extracted content.

If a future external search/vector engine is introduced:

- its index is not the sole copy of user content;
- index schema/version is recorded;
- dual-write/backfill/cutover is used when needed to avoid search downtime or inconsistent results;
- rollback can return reads to the previous validated index while authoritative Bookmarks data remains unchanged;
- privacy-level changes must invalidate/rebuild affected derived index entries.

## 16. Job-schema migrations

Durable jobs may span deployments.

A release must not deploy a worker that can no longer understand still-pending jobs without one of these governed transitions:

- backward-compatible job payload handling;
- explicit queued-job migration;
- safe cancellation/recreation of jobs whose semantics permit it;
- draining the old queue with the compatible worker before cutover.

Job payloads that may outlive one release should include a job type/schema version.

## 17. Import/export compatibility

Structured exports must contain a format/version marker and enough metadata to identify the export schema.

Importers must reject unsupported future versions explicitly rather than silently dropping unknown critical data.

When an older export can be upgraded safely, conversion occurs through versioned import logic with preview/warnings before commit.

Round-trip tests must eventually prove preservation of the fields promised by the export format.

## 18. Recovery after failed migration

A migration is not complete merely because a command exited successfully.

Post-migration validation must check, as applicable:

- schema version and constraints;
- representative resource counts/relationships;
- authorization ownership boundaries;
- archive references and integrity metadata;
- search/index rebuild state;
- synchronization cursor/change-feed behavior;
- pending job compatibility;
- API contract smoke tests;
- backup/recovery evidence freshness.

If validation fails, the system remains not-ready until the documented repair/rollback path completes.

## 19. Exact-revision evidence

Migration, rollback, and compatibility claims must cite the exact repository/application revision and, where material, the exact database migration boundary tested.

A migration tested on one revision is not automatically evidence for a later changed migration set.

## 20. Current Concept-stage boundary

No database schema, migration runner, client database, API compatibility test, archive conversion, or rollback procedure has yet been implemented or executed. This document defines required behavior for future implementation and qualification; it does not create runtime evidence.