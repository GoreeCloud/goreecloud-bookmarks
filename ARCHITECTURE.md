# GoreeCloud Bookmarks — Architecture

## Status and authority

- **Product:** GoreeCloud Bookmarks
- **Repository:** `GoreeCloud/goreecloud-bookmarks`
- **Release lifecycle:** Concept
- **Architecture state:** Selected implementation direction; no runtime implementation or platform acceptance is established by this document.
- **Authoritative product scope:** `SPECIFICATIONS.md` and `GoreeCloud/Projects/Project Specification — Bookmarks.md`
- **Current Platform Contract:** `0.4`
- **Current repository license:** `AGPL-3.0-or-later` fallback unless an authorized Bookmarks-specific decision supersedes it.

This document selects the initial technical architecture for GoreeCloud Bookmarks. It converts previously open technology questions into implementation direction while preserving the current verified reality: the repository does not yet contain a Bookmarks application, service, client, build, deployment, or accepted runtime.

## 1. Architectural principles

The initial architecture must remain:

- self-hostable and open source;
- private by default;
- service-backed with one authoritative Bookmarks data model;
- offline-capable on installed clients;
- portable across approved GoreeCloud infrastructure;
- recoverable independently of any individual client;
- versioned through documented APIs rather than direct client database access;
- modular enough to split later without requiring premature microservices;
- dependency-conscious and simple enough to operate, back up, restore, test, and migrate;
- compatible with all nine Integral Platform Systems without representing unimplemented integrations as accepted.

GoreeCloud Browser, desktop, mobile, and web experiences are clients and integration surfaces. They must not become competing authorities for the bookmark library.

## 2. Initial system shape

GoreeCloud Bookmarks will begin as a **modular service-backed application**, not as a collection of independently deployed microservices.

The first server implementation will use one GoreeCloud Bookmarks codebase with explicit internal modules for the product domains already described in the product specification:

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

These are logical service boundaries. They do **not** require separate network services, repositories, databases, containers, or release lifecycles at the initial stage.

A logical module may be extracted into a separate process or service only when verified requirements justify the added operational boundary, such as materially different scaling, security isolation, failure containment, resource use, deployment cadence, or maintenance ownership.

## 3. Server implementation

### Language

The initial Bookmarks server will use **Go**.

Go is selected for the API, synchronization coordination, capture orchestration, durable background processing, health endpoints, and other long-running network-service responsibilities.

The server should prefer the Go standard library and small, well-reviewed dependencies. A heavy application framework is not an architectural requirement.

Exact Go and dependency versions must be pinned and recorded when implementation begins; this Concept-stage decision does not invent version values.

### Server authority

The server is authoritative for account-level Bookmarks state that is meant to synchronize across devices, including bookmarks, collections, tags, notes, highlights, archive metadata, sharing state, automation rules, reminder state, server-visible preferences, and synchronization revisions.

Clients may retain offline state and pending mutations, but they must reconcile through the supported Bookmarks API rather than writing directly to the server database.

## 4. Web client

The web application will use **TypeScript** and the current approved Glaze UI contract applicable at implementation time.

The web experience is a first-class client of the same versioned Bookmarks API used by other clients. It must not depend on direct database access or undocumented server internals.

The exact TypeScript UI framework remains an implementation-level selection. It may be selected only after evaluating accessibility, Glaze UI compatibility, maintenance burden, build reproducibility, security, and long-term technology independence. The architecture does not require React, Vue, Svelte, or another specific framework merely for convention.

For the initial deployment, compiled web assets should be served by the Bookmarks server when practical so the product does not require a separate frontend runtime solely to deliver static application assets.

## 5. Native clients

### Android

The Android client will use **Kotlin with native Android architecture and platform APIs**.

This direction supports share-target capture, offline storage, notifications, background synchronization, device authentication, file/image capture, and other Android-integrated behaviors without reducing the app to a WebView wrapper.

### Apple platforms

Apple clients will use **Swift with native Apple frameworks** where iOS, iPadOS, or macOS support is approved.

The client must use the same authoritative Bookmarks service and versioned API as other clients while adapting Glaze UI to native Apple interaction conventions.

### Linux desktop

The Linux desktop client direction is **Rust with GTK 4**.

Rust and GTK 4 are selected for a native Linux experience with efficient local behavior, filesystem and desktop integration, notifications, keyboard workflows, drag-and-drop, protocol handling, offline storage, and Glaze UI adaptation without relying on a browser wrapper as the product architecture.

`libadwaita` may be evaluated as an implementation dependency only if it remains compatible with approved Glaze UI behavior, accessibility, distribution targets, and technology-independence requirements. It is not required by this architecture decision.

The planned Linux release hierarchy remains `.deb`, Flatpak, then AppImage when each format is technically appropriate and separately validated.

## 6. Server database and relational state

The initial authoritative relational database will be **PostgreSQL**.

This is a Bookmarks-specific architecture selection rather than a GoreeCloud-wide mandate.

PostgreSQL will hold relational and transactional state such as:

- users and service-facing identity references;
- devices;
- bookmarks and metadata;
- collections and membership relationships;
- tags and bookmark-tag relationships;
- notes and highlights;
- reminders;
- shares and permissions;
- automation rules;
- archive metadata and archive-version metadata;
- link-health state;
- synchronization revisions and mutation records;
- durable job state;
- saved searches and server-visible preferences.

The Bookmarks deployment should use a dedicated database owned by the Bookmarks stack unless a later approved shared-database design provides a documented operational benefit without weakening isolation, backup, migration, or recovery.

Schema evolution must use ordered, version-controlled migrations. Applied production migrations must not be casually rewritten.

## 7. Search architecture

**PostgreSQL full-text search** is the initial search/indexing engine for Phase 1 metadata and extracted text.

The initial architecture intentionally avoids requiring a second search service before measured requirements justify one.

Search code must use an internal abstraction so a future dedicated search engine can be introduced for materially larger libraries, specialized indexing, semantic retrieval, or other demonstrated requirements without changing the authoritative bookmark model.

Optional semantic or meaning-based search may later use a separate vector or semantic index. It must remain optional and must not become necessary for core bookmark capture, organization, exact-text search, archival, import/export, or synchronization.

## 8. Archive and preservation architecture

Archive metadata belongs in PostgreSQL. Preserved payloads do not.

Preserved content will use a storage abstraction whose initial self-hosted implementation writes to a dedicated persistent Bookmarks archive location outside the container writable layer. The abstraction must permit a future object-storage backend without making an external object store mandatory for initial operation.

### Preservation representations

A saved item may have multiple related representations:

- original metadata and source URL;
- normalized readable content;
- visual snapshot where enabled;
- complete web-preservation capture;
- derived previews or thumbnails;
- later archive versions.

### Complete capture format

The target standard container for complete web-preservation captures is **WARC 1.1 / ISO 28500:2017**.

WARC is used for captured protocol responses, payloads, related metadata, and preservation-oriented records. Readable text and visual snapshots may remain separately addressable derived representations linked to the same archive version rather than being forced into one presentation format.

Every persisted archive object must have a stable identifier and integrity metadata such as a cryptographic digest. Archive metadata must identify ownership, source bookmark, capture time, representation type, storage location, size, format, integrity state, and retention state where applicable.

An archive capture is not considered usable merely because bytes exist. Safe-viewer behavior, integrity validation, restoration, and recovery remain separate requirements.

## 9. Durable background work

The initial background-job mechanism will be a **PostgreSQL-backed durable job queue** processed by Go workers.

This avoids adding Redis, RabbitMQ, or another queueing platform before a verified requirement justifies it.

Jobs must support, as applicable:

- durable persistence before execution;
- unique job or mutation identifiers;
- retry with bounded backoff;
- idempotent processing where retry is possible;
- explicit success/failure state;
- attempt counts and timestamps;
- dead-letter or terminal-failure state;
- safe operator inspection without exposing unnecessary bookmark content;
- cancellation where the underlying operation can be safely stopped.

The worker may initially run from the same Go codebase and image as the API with a distinct process role. A separate worker service or queue technology may be introduced later only when justified by measured operational requirements.

## 10. API contract

The primary application API will be **HTTPS + REST-style JSON** under a versioned path beginning with:

`/api/v1/`

The API will maintain an **OpenAPI** contract for implemented endpoints.

The API version is independent from the application release version.

Required architectural behavior includes:

- stable resource identifiers;
- application-level authentication and authorization;
- OIDC/OAuth 2.0 integration with GoreeCloud Identity where applicable and implemented;
- request validation and stable error structures;
- cursor-based pagination for large collections and histories where appropriate;
- explicit filtering and sorting contracts;
- entity revision values and/or ETags for conflict-aware mutation;
- `If-Match` or equivalent preconditions where stale writes must be rejected;
- idempotency keys for retry-sensitive capture, import, bulk, or other create operations where duplicate execution would be harmful;
- no ordinary client access to the Bookmarks database;
- bounded file/archive transfer endpoints with authorization, validation, integrity checks, and resumability where required.

GraphQL, gRPC, WebSockets, or another transport may be added only when a demonstrated requirement is not adequately served by the versioned REST contract.

## 11. Synchronization model

Installed clients will use a **local SQLite database** for offline metadata, searchable local state, selected offline content references, synchronization cursors, and pending mutations.

SQLite is client-local state, not a second account authority.

Initial synchronization will use the versioned HTTPS API rather than a separate custom persistent socket protocol.

The server will expose an incremental synchronization/change-feed contract under the versioned API. The precise endpoint schema will be defined before implementation, but the model must include:

- stable object IDs;
- per-object or per-change revision information;
- device identity where applicable;
- durable client mutation IDs;
- incremental cursors/checkpoints;
- retry-safe submissions;
- explicit deletion/tombstone semantics;
- conflict responses with enough information to preserve data.

Conflict handling must prefer preservation over silent overwrite:

- independent field changes may merge when the rules are deterministic and safe;
- stale writes to the same logical field should produce an explicit conflict rather than silently discarding user work;
- ambiguous user-authored text conflicts must preserve variants or enough revision history for recovery;
- last-write-wins may be used only for fields where the consequence is explicitly considered safe.

GoreeCloud Sync remains a separately governed GoreeCloud capability. Any later integration must preserve Bookmarks as the authority for bookmark data and must not introduce direct cross-application database ownership.

## 12. Initial deployment model

The initial self-hosted server deployment will use a **Docker Compose Bookmarks stack**.

The minimal production topology is intended to be:

1. Bookmarks application service — Go API plus compiled web assets.
2. Dedicated PostgreSQL service.
3. Background worker process/service from the same Bookmarks codebase when asynchronous work requires operational separation.

The design must avoid additional infrastructure containers unless a verified requirement justifies them.

Deployment rules include:

- Bookmarks application and PostgreSQL communicate over a purpose-specific internal Docker network;
- PostgreSQL must not be published as a general host/public port;
- the application may join the approved Caddy/proxy network when it needs web access through the approved publication path;
- application state, database state, and archive payloads must survive container replacement in approved persistent locations;
- active secrets and environment-specific protected values remain outside ordinary Git history;
- production images must be pinned according to current Docker governance;
- the service must expose appropriate health and readiness behavior before production qualification;
- test and production data, credentials, persistent storage, and deployment state remain separated;
- no deployment hostname, host, DNS entry, Caddy route, NetBird policy, port, secret, or production path is claimed by this architecture document until separately implemented and verified.

## 13. Backup, export, and recovery boundaries

Synchronization is not backup.

The server recovery model will include distinct protection for:

- PostgreSQL through database-supported backups and independently protected recovery points;
- archive payload storage through versioned independent backup appropriate to its size and importance;
- deployment/configuration artifacts required to recreate the stack;
- protected secret recovery references without placing active secrets in the repository;
- application-native export for user portability and account-level data transfer.

Database files must not be assumed recoverable merely because their live storage directory was copied.

Recovery qualification must eventually prove that the database, archive relationships, permissions, service configuration, authentication path, search reconstruction, background jobs, and user-visible library can be restored to an approved state.

## 14. Security and privacy boundaries

The architecture requires:

- least-privilege service and database credentials;
- application-level authorization even on private networks;
- private-by-default bookmark libraries;
- no unnecessary behavioral analytics or advertising telemetry;
- no active secrets in source control or client packages;
- safe handling of hostile URLs, fetched metadata, documents, and archived content;
- isolated archive viewing with active content disabled by default;
- controlled outbound requests for capture and metadata processing;
- clear privacy behavior for Sensitive and Private Vault content;
- protection against cross-user access to bookmarks, archives, highlights, notes, and shares;
- security review before claims of Wardveil, Privacy Shield, Identity, Policy, or other platform-system conformance.

Private Vault remains an advanced privacy mode. Its cryptographic protocol and key-management design are not selected by this architecture document and must be designed separately before implementation.

## 15. Observability boundary

Operational evidence should describe service health without unnecessarily exposing bookmark content.

The implementation should eventually expose evidence for:

- API health and readiness;
- database connectivity and migration state;
- capture/enrichment queue depth and failure rate;
- archive job state;
- search indexing lag;
- synchronization failures/conflicts;
- link-health backlog;
- storage capacity and archive integrity state;
- notification failures;
- backup and restore-test status where integration permits.

Logging must avoid tokens, credentials, full private content, and unnecessary sensitive URLs or request bodies.

## 16. Planned repository implementation structure

No empty implementation directories are created by this architecture decision. When source implementation begins, the repository should use ecosystem-appropriate locations similar to:

```text
cmd/                    Go service/worker entry points
internal/               Go application domains and internal packages
web/                    TypeScript web client
clients/android/        Kotlin Android client
clients/linux/          Rust + GTK 4 Linux client
clients/apple/          Swift Apple client(s) when approved
migrations/             PostgreSQL schema migrations
docs/                   supplementary API, deployment, recovery, and architecture details
scripts/                governed development/build/test utilities when needed
```

The exact structure may adapt to selected build systems and platform conventions. Empty scaffolding must not be created solely to make the repository look complete.

## 17. Decisions intentionally deferred

The following remain separate implementation or governance decisions rather than being invented here:

- exact Go, PostgreSQL, SQLite, Rust, GTK, Kotlin, Swift, and TypeScript versions;
- exact TypeScript UI framework;
- exact dependency libraries and package managers;
- exact WARC creation/replay library and compression implementation;
- exact Private Vault encryption protocol and key-management design;
- exact production host, hostname, DNS, Caddy, NetBird, storage paths, secrets, ports, and resource limits;
- exact backup frequency and retention values;
- exact release signing and distribution infrastructure;
- whether a future Bookmarks-specific license should supersede the current `AGPL-3.0-or-later` fallback;
- whether measured scale later justifies extracting search, archive processing, workers, or another module into an independently deployed service.

Deferred decisions must be recorded and verified before they become operational requirements.

## 18. Capability and lifecycle boundary

This architecture document is **not implementation evidence**.

It does not establish that:

- the Go service exists;
- PostgreSQL or SQLite has been deployed;
- WARC capture works;
- any native client exists;
- synchronization works;
- Docker Compose deployment exists;
- any Platform Contract integration has passed acceptance;
- Bookmarks has advanced beyond Concept.

`CAPABILITIES.md` remains the repository authority for currently verified capabilities. Until implementation and exact-revision evidence exist, the selected technologies in this document remain architecture decisions only.