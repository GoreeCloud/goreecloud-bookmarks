# GoreeCloud Bookmarks — Capabilities

## Overview

This file records the current verified capability state represented by the source in `GoreeCloud/goreecloud-bookmarks`.

**Verified state date:** 2026-09-17  
**Release lifecycle:** Experimental  
**Implementation status:** Minimal Go service foundation implemented; no end-user bookmarking workflow or persistent Bookmarks data layer is implemented.  
**Platform Contract:** `0.4`, conformance `unverified`.

The planned product vision remains defined by `SPECIFICATIONS.md`, `ARCHITECTURE.md`, the repository API/data contracts, and the canonical GoreeCloud project specification. Planned functionality must not be interpreted as current capability.

## Verified experimental service foundation

The repository now contains a minimal executable Go service foundation using the Go standard library.

Implemented behavior is limited to:

- `GET /api/v1/health` returning HTTP `200` with bounded JSON process-health state.
- `GET /api/v1/ready` returning HTTP `503` with `ready: false` while the required Bookmarks data layer is not configured.
- `Cache-Control: no-store` and `X-Content-Type-Options: nosniff` response headers on the current HTTP surface.
- Method-aware routing that rejects unsupported methods for the implemented routes.
- A local-development default listen address of `127.0.0.1:8080` with an optional `GOREECLOUD_BOOKMARKS_LISTEN_ADDR` override.
- Bounded HTTP timeouts and graceful SIGINT/SIGTERM shutdown.
- Unit tests for health, fail-closed readiness, method handling, and listen-address behavior.
- Repository CI that checks the pinned Go toolchain, formatting, vetting, tests, and service build.

This foundation is an engineering prototype. Health means only that the process can answer its bounded health request. Readiness intentionally remains non-passing because no authoritative bookmark persistence layer exists.

## User capabilities

No end-user bookmark creation, Inbox, collections, tags, favorites, Read Later, metadata extraction, search, archival, reader, annotation, synchronization, offline library, sharing, collaboration, reminders, feeds, rediscovery, automation, import, export, web interface, desktop client, mobile client, or GoreeCloud Browser integration is currently verified as implemented.

## Data and persistence

No PostgreSQL schema, database connection, durable bookmark record, migration runner, SQLite client database, archive store, search index, or synchronization store is currently implemented.

The planned data model and migration rules remain documented in `docs/data-model.md` and `docs/migrations.md`; documentation is not persistence evidence.

## API

The planned v1 API contract remains in `docs/api/openapi.yaml` and `docs/api/README.md`.

Only the bounded health and readiness routes listed above are currently implemented. All bookmark, collection, tag, note, highlight, archive, reminder, sharing, automation, link-health, saved-search, and synchronization operations in the planned contract remain unimplemented.

## Administrative capabilities

No administrative interface, account administration, processing queue, retention control, storage-management interface, link-health administration, search-index administration, or multi-user administration is implemented.

The current health/readiness endpoints are engineering service signals, not an accepted GoreeCloud Manager or GoreeCloud Observability integration.

## Integral Platform Systems

Current GoreeCloud governance requires evaluation against exactly nine Integral Platform Systems:

1. GoreeCloud Manager
2. Privacy Shield
3. Wardveil Security
4. Everkeep
5. Glaze UI
6. GoreeCloud Mesh
7. GoreeCloud Identity
8. GoreeCloud Policy
9. GoreeCloud Observability

The Platform Contract declaration keeps all nine `applicable-blocked`. No Bookmarks-specific acceptance evidence is currently verified for any of them. GoreeCloud Sync remains separately governed and is not a tenth Integral Platform System.

## Security and privacy

The current service foundation does not process user bookmark content, accounts, credentials, archives, searches, or synchronization data. It contains no active secret values and has no third-party runtime module dependency.

That narrow state does not establish that GoreeCloud Bookmarks is secure, hardened, private-by-default in a complete product sense, Wardveil-conformant, Privacy-Shield-conformant, production-ready, or accepted. Those claims require implementation and attributable evidence at the relevant revision.

## Resilience, backup, and recovery

Graceful process shutdown is implemented for the experimental service foundation. No Bookmarks data backup, restore, archive recovery, migration recovery, tested rollback, or Everkeep integration is currently implemented because no authoritative Bookmarks data layer exists yet.

## Accessibility and user interface

No Bookmarks user interface is currently implemented, so no Bookmarks-specific accessibility or Glaze UI acceptance claim is established.

## Build and validation boundary

The repository pins Go `1.27.1` and includes `.github/workflows/validate-go.yml` to validate formatting, `go vet`, unit tests, and service compilation on exact pull-request/source revisions.

Source presence, passing CI, merge state, release state, deployment state, platform acceptance, production readiness, and Stable qualification are distinct states and must not be conflated.

## Capability validation rule

A future capability may be added here only when source and applicable tests/evidence support the exact claimed scope. Product behavior, platform integration, release, deployment, and production acceptance must remain separately evidenced.
