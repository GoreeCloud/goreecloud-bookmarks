# GoreeCloud Bookmarks Specifications

## Product

GoreeCloud Bookmarks is GoreeCloud's first-party bookmark collection, organization, preservation, and synchronization application. The target architecture is original GoreeCloud-owned native software. Retained Linkwarden-derived code is transitional migration and compatibility material, not the target product architecture.

## Lifecycle

- Native migration: in progress
- Production approved: No
- Stable: No
- Canonical repository: `GoreeCloud/goreecloud-bookmarks`
- Native module: `native/`

A source implementation, passing CI run, or retained inherited runtime does not by itself authorize production migration or Stable promotion.

## Native service boundary

`native/cmd/bookmarksd` is the current GoreeCloud-owned development service. It uses explicit repository and identity boundaries.

Default runtime behavior is deliberately fail-closed for production claims:

- bookmark repository mode defaults to `memory-development`;
- production persistent storage is not accepted;
- GoreeCloud Identity integration is not production-ready;
- the development header identity seam is permitted only with the in-memory development repository;
- health/readiness metadata must not claim production approval.

## Bookmark contract

The native bookmark domain supports owner-scoped create, list, get, update, and delete operations.

Current input limits include:

- title: 512 runes;
- note: 8192 runes;
- tags: at most 32;
- tag length: 64 runes;
- HTTP mutation body: at most 16 KiB.

Bookmark URLs must use HTTP or HTTPS and have a host. URL fragments are removed before storage. Owner identity is required before bookmark state is accessed.

## Collection contract

The native collection domain supports owner-scoped hierarchy organization.

Current development HTTP routes are:

- `GET /api/v1/collections`
- `GET /api/v1/collections/{id}`
- `POST /api/v1/collections`
- `PATCH /api/v1/collections/{id}` for parent moves

Collection names are required and limited to 256 runes. Sibling names are unique case-insensitively within the same parent. Parent references are owner-scoped. Self-parenting and hierarchy cycles are rejected. Cross-owner collection lookups are represented as not found rather than exposing another owner's state.

Collection state currently uses the in-memory development store. This route surface is not evidence of persistent collection storage, completed bookmark-assignment HTTP behavior, shared collections, permission management, or Linkwarden migration completion.

## Request parsing

Native bookmark and collection mutation requests use bounded strict JSON decoding:

- unknown fields are rejected;
- multiple/trailing JSON values are rejected;
- request bodies are size-limited;
- malformed input receives bounded error responses rather than internal error detail.

## Identity boundary

GoreeCloud Identity is authoritative for production authentication and authorization.

The `X-GoreeCloud-Development-Owner` resolver exists only for local development. It is not a production authentication mechanism and is restricted to the `memory-development` repository mode by runtime selection policy.

## Synchronization boundary

The native GoreeCloud Sync client owns the `bookmarks.items` capability and currently includes application-side signing, authenticated submission, deletion tombstones, retrieval, pagination, schema validation, identifier limits, and canonical record-proof interoperability testing.

Sync source foundations do not establish production identity, deployment, persistent bookmark storage, complete conflict resolution, or Stable protocol acceptance.

## Privacy and security

Privacy Shield is authoritative for privacy governance. Wardveil Security is authoritative for GoreeCloud security state and protection evidence.

Current native principles include owner-scoped data access, bounded input, payload-free Sync tombstones, no production use of the development identity header, and fail-closed runtime modes when reviewed integration is absent.

## Continuity

Everkeep is authoritative for backup, recovery, portability, migration, and preserved-content continuity. Native in-memory collection state is not a backup or recovery implementation. Stable migration remains blocked on durable storage, import/export, backup/restore, rollback, and retained Linkwarden data migration acceptance.

## Glaze UI

Glaze UI is authoritative for GoreeCloud-owned presentation. Native web surfaces must track the current approved Stable Glaze UI contract and complete product-specific accessibility and real-runtime acceptance before Stable promotion.

## Migration boundary

The inherited Linkwarden runtime may remain available while required functionality and user data are migrated. Removal requires explicit evidence for retained bookmark data, collection relationships, preserved content, import/export behavior, account ownership, backup/restore, rollback, and approved feature-retirement decisions.

## Stable blockers

Stable remains blocked by at least:

- production GoreeCloud Identity integration;
- accepted persistent bookmark and collection storage;
- complete collection/bookmark relationship behavior;
- retained-data migration from the transitional runtime;
- current Glaze UI acceptance;
- current Wardveil Security acceptance;
- current Privacy Shield acceptance;
- current Everkeep backup/restore and migration acceptance;
- representative runtime, accessibility, performance, recovery, and operational evidence;
- final release governance and production approval.
