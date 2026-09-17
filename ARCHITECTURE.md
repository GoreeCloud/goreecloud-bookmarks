# GoreeCloud Bookmarks — Architecture

## Status and authority

- **Product:** GoreeCloud Bookmarks
- **Repository:** `GoreeCloud/goreecloud-bookmarks`
- **Release lifecycle:** Experimental
- **Architecture state:** Selected architecture with a minimal executable Go service foundation; broader product runtime remains unimplemented.
- **Authoritative product scope:** `SPECIFICATIONS.md` and `GoreeCloud/Projects/Project Specification — Bookmarks.md`
- **Current Platform Contract:** `0.4`
- **Current repository license:** `AGPL-3.0-or-later` fallback unless an authorized Bookmarks-specific decision supersedes it.
- **API/data contracts:** `docs/api/README.md`, `docs/api/openapi.yaml`, `docs/data-model.md`, and `docs/migrations.md`.

This document separates the architecture that has been selected from the subset that is currently implemented. `CAPABILITIES.md` controls current capability claims.

## 1. Architectural principles

GoreeCloud Bookmarks is intended to remain:

- self-hostable and open source;
- private by default;
- service-backed with one authoritative Bookmarks data model;
- offline-capable on installed clients;
- portable across approved GoreeCloud infrastructure;
- recoverable independently of any individual client;
- versioned through documented APIs rather than direct client database access;
- modular enough to split later without premature microservices;
- dependency-conscious and understandable to future maintainers;
- compatible with all nine Integral Platform Systems without representing unimplemented integrations as accepted.

GoreeCloud Browser, desktop, mobile, and web experiences are clients/integration surfaces rather than competing authorities for bookmark state.

## 2. Initial system shape

The selected architecture is a **modular service-backed application**, not an initial fleet of independently deployed microservices.

Logical modules include:

- bookmarks, collections, tags, notes, highlights, sharing, and permissions;
- capture and enrichment coordination;
- metadata extraction;
- archival and preservation;
- readable-content extraction;
- search;
- synchronization;
- automation;
- link health;
- reminders and notifications.

These boundaries do not initially require separate repositories, network services, databases, containers, or release lifecycles. Extraction into a separate service requires evidence of a meaningful scaling, isolation, failure-containment, resource, deployment-cadence, or ownership need.

## 3. Server implementation

### Current implementation

The first executable server foundation is implemented in Go `1.27.1` using only the Go standard library at runtime.

Current source provides:

- `cmd/bookmarks/main.go` as the process entry point;
- `internal/httpapi` as the initial HTTP boundary;
- bounded server timeouts;
- graceful process shutdown;
- loopback-only default development binding;
- health and fail-closed readiness endpoints;
- unit tests and repository CI validation.

It does **not** yet provide bookmark persistence, authentication, authorization, PostgreSQL access, background jobs, archive processing, search, synchronization, or user-facing product functionality.

### Server authority

When the data layer is implemented, the Bookmarks server is selected as authoritative for account-level state intended to synchronize across devices, including bookmarks, collections, tags, notes, highlights, archive metadata, sharing state, automation rules, reminder state, server-visible preferences, and synchronization revisions.

Clients may retain offline state and pending mutations but must reconcile through supported interfaces rather than write directly to the server database.

## 4. Web client

The planned web application uses **TypeScript** and the approved Glaze UI contract applicable at implementation time.

It is intended to consume the same versioned Bookmarks API as other clients and must not depend on direct database access or undocumented server internals.

The exact TypeScript UI framework remains deferred pending accessibility, Glaze UI, security, maintenance, reproducibility, and technology-independence review. No web client is currently implemented.

## 5. Native clients

### Android

Selected direction: **Kotlin with native Android architecture and platform APIs** for share-target capture, offline state, notifications, background synchronization, device authentication, file/image capture, and platform integration.

### Apple platforms

Selected direction: **Swift with native Apple frameworks** for approved iOS/iPadOS/macOS clients.

### Linux desktop

Selected direction: **Rust with GTK 4**. `libadwaita` may be evaluated only if compatible with Glaze UI, accessibility, distribution, and independence requirements.

The planned Linux distribution preference remains `.deb`, then Flatpak, then AppImage when each is technically appropriate and validated.

No native Bookmarks client is currently implemented.

## 6. Server database and relational state

Selected authoritative relational database: **PostgreSQL**.

Planned relational/transactional state includes users/service identity references, devices, bookmarks, collections, tags, notes, highlights, reminders, shares, permissions, automation rules, archive metadata, link-health state, synchronization revisions, durable jobs, saved searches, and server-visible preferences.

A dedicated Bookmarks database is preferred unless a later approved design provides a documented benefit without weakening isolation, backup, migration, or recovery.

Schema evolution must use ordered, version-controlled migrations. `docs/data-model.md` defines logical ownership; `docs/migrations.md` defines version, compatibility, destructive-change, rollback, and recovery requirements.

**Current state:** no PostgreSQL connection, schema, migration, or persistent bookmark record is implemented. This is why readiness intentionally fails.

## 7. Search architecture

Selected initial search engine: **PostgreSQL full-text search** for Phase 1 metadata and extracted text.

Search code should use an internal boundary so a future dedicated or semantic index can be introduced only when justified. Optional semantic retrieval must not become a dependency for core capture, organization, exact-text search, archival, import/export, or synchronization.

Search indexes are derived state and must remain rebuildable from authorized authoritative data.

**Current state:** no Bookmarks search runtime or index exists.

## 8. Archive and preservation architecture

Archive metadata belongs in PostgreSQL; preserved payload bytes do not.

Selected storage direction uses a persistent archive-storage abstraction outside the container writable layer, allowing future object-storage support without requiring an external object store for initial self-hosting.

A saved item may have original metadata/source URL, normalized readable content, visual snapshots, complete preservation captures, previews/thumbnails, and later versions.

The target complete-capture container is **WARC 1.1 / ISO 28500:2017**. Historical preservation bytes are immutable for migration/provenance purposes; conversions or relocations must preserve integrity and provenance rather than invisibly rewriting history.

Persisted archive objects require stable identifiers and integrity metadata such as cryptographic digests.

**Current state:** no archive storage or capture runtime exists.

## 9. Durable background work

Selected initial mechanism: **PostgreSQL-backed durable jobs processed by Go workers** rather than adding Redis, RabbitMQ, or another queue before a verified requirement exists.

Jobs should support durable persistence, unique IDs, bounded retry/backoff, idempotency where needed, explicit states, attempt history, terminal failure handling, safe operator inspection, and cancellation where safe.

Long-lived jobs must remain compatible across deployments or be explicitly migrated/drained/cancelled/recreated according to `docs/migrations.md`.

**Current state:** no background-job store or worker runtime exists.

## 10. API contract

The selected application API is **HTTPS + REST-style JSON** under `/api/v1/` with OpenAPI `3.2.1` planning records in `docs/api/`.

The planned resource surface includes Bookmark, Collection, Tag, Note, Highlight, Archive, ArchiveVersion, Reminder, Share, Permission, AutomationRule, LinkHealthRecord, SavedSearch, synchronization changes/mutations, and bounded health/readiness operations.

Required design behavior includes stable opaque IDs, authorization, request validation, stable errors, cursor pagination, explicit filtering/sorting, revisions/ETags, `If-Match` for overwrite-sensitive changes, idempotency for retry-sensitive operations, and privacy-sensitive error behavior.

### Currently implemented routes

Only these routes are currently implemented:

- `GET /api/v1/health` — returns HTTP `200` with bounded process-health JSON.
- `GET /api/v1/ready` — returns HTTP `503` with `ready: false` and `bookmarks-data: not-configured` until the required data layer exists.

The current health/readiness implementation does not require authentication because it returns only bounded service state and no user/content/dependency detail.

Every other operation in the OpenAPI planning contract remains unimplemented.

## 11. Synchronization model

Selected installed-client local store: **SQLite** for offline metadata, searchable local state, selected offline references, synchronization cursors, and pending mutations.

SQLite is client-local state, not a second account authority. Unsynchronized user-authored data is non-disposable even when other local cache state can be rebuilt.

The planned v1 sync contract uses `/api/v1/sync/changes` and `/api/v1/sync/mutations`, stable IDs/revisions, durable client mutation IDs, opaque cursors, retry-safe batches, deletion/tombstones, per-mutation outcomes, and preservation-first conflicts.

GoreeCloud Sync remains separately governed. Any later integration must preserve Bookmarks domain authority.

**Current state:** no SQLite client, synchronization endpoint, mutation queue, or GoreeCloud Sync integration exists.

## 12. Initial deployment model

Selected future self-hosted deployment boundary: **Docker Compose** with:

1. Bookmarks application service — Go API plus compiled web assets.
2. Dedicated PostgreSQL service.
3. Worker process/service from the same Bookmarks codebase when asynchronous processing requires operational separation.

Future deployment must use purpose-specific internal networking, avoid general PostgreSQL host publication, preserve state outside replaceable containers, keep secrets outside Git, pin production images, expose truthful health/readiness, and separate test/production credentials/data/storage.

No hostname, DNS record, Caddy route, NetBird policy, public/private production port, secret, production storage path, container stack, or deployment is established by the current experimental source.

## 13. Backup, export, and recovery boundaries

Synchronization is not backup.

Future recovery must separately protect PostgreSQL, archive payload storage, deployment/configuration artifacts, protected secret recovery references, and application-native portable export.

Recovery qualification must eventually prove database, archive relationships, permissions, configuration, authentication, search reconstruction, jobs, migration state, API behavior, and user-visible library restoration.

Destructive migrations require a validated recovery point and repair/rollback plan under `docs/migrations.md`.

**Current state:** graceful process shutdown exists, but no Bookmarks data backup, restore, migration recovery, tested rollback, or Everkeep integration exists.

## 14. Security and privacy boundaries

The architecture requires least privilege, application-level authorization, private-by-default libraries, no unnecessary behavioral analytics, secrets outside source control, safe hostile-content handling, controlled outbound fetching, explicit higher-privacy behavior, cross-user isolation, privacy-sensitive error handling, and platform security/privacy review.

The current experimental process defaults to loopback, uses bounded timeouts, sets basic defensive response headers, includes no active secrets, and has no user-content data plane. These narrow properties are not production-security evidence.

Private Vault cryptographic protocol and key management remain separately deferred.

## 15. Observability boundary

Operational evidence should describe service health without unnecessarily exposing bookmark content.

Current source exposes bounded health/readiness signals only. These do not constitute GoreeCloud Observability acceptance.

Future operational evidence should cover database connectivity/migration state, capture/enrichment jobs, archive jobs, search lag, synchronization failures/conflicts, link-health backlog, storage/integrity state, notification failures, and backup/restore-test status.

Logging must avoid tokens, credentials, full private content, and unnecessary sensitive URLs/request bodies.

## 16. Repository implementation structure

Current implemented source structure begins as:

```text
cmd/bookmarks/           Go service entry point and entry-point tests
internal/httpapi/        HTTP routing, health/readiness behavior, and tests
docs/api/                planned API contract and guidance
docs/data-model.md       logical data ownership/model
docs/migrations.md       version/migration/rollback requirements
.github/workflows/       repository validation
go.mod                   pinned Go toolchain baseline
```

Future source may add conventional locations such as `web/`, `clients/android/`, `clients/linux/`, `clients/apple/`, `migrations/`, and supplementary `docs/` areas when implementation actually requires them. Empty scaffolding must not be created solely to make the repository look complete.

## 17. Decisions intentionally deferred

Still deferred:

- PostgreSQL, SQLite, Rust, GTK, Kotlin, Swift, TypeScript, database-driver, and other dependency/tool versions not yet implemented;
- exact TypeScript UI framework;
- exact WARC creation/replay library and compression implementation;
- exact Private Vault encryption protocol and key-management design;
- production host, hostname, DNS, Caddy, NetBird, storage paths, secrets, ports, and resource limits;
- backup frequency and retention values;
- release signing/distribution infrastructure;
- database migration tooling and client SQLite migration implementation;
- any future Bookmarks-specific license decision superseding the fallback;
- whether measured scale later justifies extracting search, archive processing, workers, or another module.

Go is no longer an entirely deferred tool choice: the experimental service foundation pins Go `1.27.1`.

## 18. Capability and lifecycle boundary

The executable Go service foundation is implementation evidence for its narrow source/tested behavior only.

It establishes an **Experimental** prototype foundation, not a complete Bookmarks product.

It does **not** establish that:

- bookmark creation or durable bookmark persistence works;
- PostgreSQL or SQLite is configured;
- authentication or authorization exists;
- planned bookmark/collection/tag/note/archive/sync API operations exist;
- WARC capture works;
- any web/native client exists;
- synchronization works;
- Docker Compose deployment exists;
- any of the nine Integral Platform Systems has passed Bookmarks acceptance;
- production readiness, Release Candidate, or Stable requirements are met.

`CAPABILITIES.md` remains the repository authority for current verified capability scope. Every future claim must remain attributable to source and applicable validation at an exact revision.
