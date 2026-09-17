# GoreeCloud Bookmarks — Notes

## Current verified repository state

- Product lifecycle: **Concept**.
- Repository: `GoreeCloud/goreecloud-bookmarks`.
- Product implementation: no application or service implementation currently verified.
- Canonical product specification: `GoreeCloud/Projects/Project Specification — Bookmarks.md`.
- Repository product specification: `SPECIFICATIONS.md`.
- Repository architecture record: `ARCHITECTURE.md`.
- Planned API contract: `docs/api/README.md` and `docs/api/openapi.yaml`.
- Logical data model: `docs/data-model.md`.
- Migration/compatibility rules: `docs/migrations.md`.
- Related capture repository: `GoreeCloud/goreecloud-bookmark-browser-extension`.
- Current Platform Contract baseline: `0.4`.
- Current Integral Platform System model: exactly nine systems.
- Current Stable Glaze UI consumer target referenced by Platform Contract `0.4`: `1.5.1`.
- Current governing repository license: `AGPL-3.0-or-later` through the GoreeCloud default fallback; no Bookmarks-specific superseding license decision is currently recorded.
- Default branch `main` is not currently protected in verified GitHub state; protection remains a repository-administration task.

## Important authority boundaries

- Planned feature, architecture, API, data-model, or migration documentation is not implementation evidence.
- `CAPABILITIES.md` records current verified capability state.
- `FEATURE-ROADMAP.md` orders planned work but does not establish completion.
- `ARCHITECTURE.md` records selected implementation direction, not deployed state.
- `docs/api/openapi.yaml` records the planned application wire contract; it does not assert that an endpoint exists or is reachable.
- `docs/data-model.md` is a logical ownership/model record, not proof of a PostgreSQL or SQLite schema.
- `docs/migrations.md` defines required evolution/rollback behavior, not proof that a migration runner or tested recovery path exists.
- GoreeCloud Bookmarks remains authoritative for bookmark data.
- GoreeCloud Sync is separately governed and is not an Integral Platform System.
- GoreeCloud Browser integration must be reconciled with the Browser project and the separate bookmark-browser-extension repository rather than creating competing capture or bookmark authority.
- Optional intelligence must not become a dependency for core bookmarking.
- The fallback license may be superseded only by an authorized Bookmarks-specific decision reconciled across repository and canonical project records.

## Selected architecture direction

The initial implementation direction is selected in `ARCHITECTURE.md`:

- One authoritative modular Bookmarks service; logical product modules begin inside one codebase instead of being deployed as premature microservices.
- Go server for API, synchronization coordination, capture orchestration, jobs, and long-running service responsibilities.
- TypeScript first-class web client using the applicable Glaze UI contract.
- Kotlin/native Android client.
- Swift/native Apple clients when those platforms are implemented.
- Rust + GTK 4 Linux desktop client direction.
- PostgreSQL as the authoritative relational server database.
- PostgreSQL full-text search for the initial core search implementation, with an internal abstraction for future dedicated search/semantic indexes if justified.
- SQLite for installed-client offline metadata, synchronization cursors, and pending mutations.
- HTTPS REST-style JSON API under `/api/v1/` with OpenAPI, explicit authorization, conflict-aware revisions/ETags, and retry-safe mutation design where required.
- PostgreSQL-backed durable background jobs processed by Go workers before introducing a separate queue platform.
- Persistent archive payload storage outside the container writable layer, separate from PostgreSQL archive metadata.
- WARC 1.1 / ISO 28500:2017 as the target container for complete web-preservation captures, with readable and visual derived representations linked to archive versions.
- Docker Compose as the initial self-hosted server deployment boundary, with the Bookmarks application, dedicated PostgreSQL, and a worker role only when asynchronous processing requires it.
- Database-native backup plus independent archive/data protection and tested recovery; synchronization is not backup.

## Selected API, data, and migration direction

The initial planned contract is now recorded under `docs/`:

- OpenAPI `3.2.1` describes the planned `/api/v1/` resource surface without a server URL or deployment claim.
- Core schemas cover Bookmark, Collection, Tag, Note, Highlight, Archive, ArchiveVersion, Reminder, Share, Permission, SyncEvent/mutations, AutomationRule, LinkHealthRecord, and SavedSearch.
- Resource identifiers are opaque; mutable synchronized resources carry revision state and overwrite-sensitive mutations use ETag/`If-Match`-style preconditions.
- Retry-sensitive creates and synchronization mutation batches use idempotency keys.
- List and synchronization feeds use opaque cursor pagination/checkpoints.
- Normal deletion produces synchronization-visible tombstone state before physical purge when required for offline consistency.
- PostgreSQL is authoritative for relational/transactional server state; archive payloads remain outside ordinary relational rows; search is derived/rebuildable state.
- Installed-client SQLite state distinguishes rebuildable cache from unsynchronized user-authored data, which is non-disposable.
- Server migrations are ordered and append-only after application; potentially breaking changes use expand → migrate → contract where practical.
- Database rollback does not assume unsafe automatic down migrations; rollback may use compatible application rollback, forward correction, or verified recovery-point restoration according to the actual change.
- Historical WARC/archive bytes remain immutable; conversions or storage moves create/verify new derived representations instead of rewriting preservation history invisibly.
- API, application, PostgreSQL schema, SQLite schema, archive representation, sync schema, and Platform Contract versions are distinct version domains.

## Decisions intentionally deferred

The following remain open until implementation or a more specific governed decision establishes them:

- Exact language, framework, database, client-toolchain, and dependency versions.
- Exact TypeScript UI framework.
- Exact dependency libraries and package managers used to implement the planned API/data model.
- Exact WARC capture/replay library and compression implementation.
- Exact Private Vault cryptographic protocol, key lifecycle, and recovery model.
- Exact production deployment host, hostname, DNS, Caddy route, NetBird policy, secrets, ports, resource limits, storage paths, and operational schedules.
- Exact backup frequency and retention values.
- Exact release signing, package publication, and client-update infrastructure.
- Whether measured scale or isolation requirements justify extracting search, archive processing, jobs, or another module into a separate deployed service.
- Whether a Bookmarks-specific license should supersede the current `AGPL-3.0-or-later` fallback after dependency, distribution, interoperability, and product-architecture review.

## Documentation maintenance

When implementation begins, update `CAPABILITIES.md`, `USER-MANUAL.md`, `PRIVACY POLICY.md`, `SECURITY.md`, `goreecloud.platform.yaml`, `ARCHITECTURE.md`, the `docs/` API/data/migration contracts, and the feature roadmap only from verified implementation and evidence. Do not infer current capability from specifications, architecture decisions, API schemas, mockups, issue labels, or intended technology choices.