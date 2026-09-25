# GoreeCloud Bookmarks — Capabilities

## Overview

This file records the current verified capability state represented by the source in `GoreeCloud/goreecloud-bookmarks`.

**Verified state date:** 2026-09-18  
**Release lifecycle:** Development  
**Implementation status:** Go/PostgreSQL service and persistence foundation plus Firefox client source candidate implemented; no authenticated end-to-end bookmarking workflow is implemented.  
**Platform Contract:** `0.4`, conformance `unverified`.

The planned product vision remains defined by `PROJECT-SPECIFICATIONS.md`, `ARCHITECTURE.md`, and the repository API/data contracts. Planned functionality must not be interpreted as current capability.

## Verified Development service and data foundation

The repository contains an executable Go service plus a bounded PostgreSQL persistence and migration foundation. The HTTP surface remains intentionally limited to health/readiness; bookmark-domain persistence is internal and not exposed as an unauthenticated user API.

Implemented behavior is limited to:

- `GET /api/v1/health` returning HTTP `200` with bounded JSON process-health state.
- `GET /api/v1/ready` failing closed unless configured PostgreSQL is reachable and the exact embedded migration history is current and untampered.
- `Cache-Control: no-store` and `X-Content-Type-Options: nosniff` response headers on the current HTTP surface.
- Method-aware routing that rejects unsupported methods for the implemented routes.
- A local-development default listen address of `127.0.0.1:8080` with an optional `GOREECLOUD_BOOKMARKS_LISTEN_ADDR` override.
- Bounded HTTP timeouts and graceful SIGINT/SIGTERM shutdown.
- Unit tests for health, fail-closed readiness, method handling, and listen-address behavior.
- Explicit `bookmarks-migrate` command for ordered, checksummed PostgreSQL schema migrations; ordinary service startup does not mutate schema automatically.
- Initial `bookmarks` table and owner-scoped internal create/read persistence for the selected Bookmark subset.
- Repository CI that checks the pinned Go toolchain, formatting, module tidiness, vetting, unit/integration tests against PostgreSQL, and both service/migration-command builds.

This remains Development-stage engineering software. Health means only that the process can answer its bounded health request. Readiness can pass only against an explicitly configured, reachable PostgreSQL database whose migration history exactly matches the embedded schema contract.

## Firefox client source candidate

`clients/firefox/` contains GoreeCloud Bookmarks Firefox client version `0.1.1` with add-on ID `goreecloud-bookmarks@goreecloud.com`. The source includes popup/options UI, bounded Firefox permissions, source validation, JavaScript syntax validation, deterministic packaging, privacy/security records, and runtime/release acceptance plans. The client is an application-owned platform variant, not a separate repository or bookmark-data authority.

The required `/api/v1/bookmarks/extension-capture` application endpoint is not present on authoritative `main`. Consequently current Firefox source does not establish successful bookmark capture, authenticated application integration, Mozilla signing, persistent install/restart acceptance, or Stable release.

## User capabilities

No end-to-end end-user bookmark creation, Inbox, collections, tags, favorites, Read Later, metadata extraction, search, archival, reader, annotation, synchronization, offline library, sharing, collaboration, reminders, feeds, rediscovery, automation, import, export, web interface, native desktop client, native mobile client, or GoreeCloud Browser integration is currently verified as implemented. Firefox client UI/source exists, but application-backed capture remains blocked and unaccepted.

## Data and persistence

PostgreSQL connectivity, migration history, schema version checks, migration version `1`, the initial `bookmarks` relation, and internal owner-scoped bookmark create/read persistence are implemented. The migration runner rejects newer-than-binary and tampered migration histories and is invoked explicitly through `cmd/bookmarks-migrate`.

No SQLite client database, archive store, search index, synchronization store, authenticated bookmark-domain HTTP route, backup/restore path, or production database deployment is implemented.

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

The current Development source can persist bookmark records when a database is explicitly configured. Database URLs are protected external configuration and are not logged or committed. The direct Go runtime dependency is `github.com/jackc/pgx/v5` `v5.11.0`; its provenance/license baseline is recorded in `docs/dependencies.md`. Exact-revision CI also generates a CycloneDX 1.6 SBOM and runs pinned `govulncheck` reachability scanning for known Go vulnerabilities. Authentication, authorization, archives, searches, and synchronization remain unimplemented.

That narrow state does not establish that GoreeCloud Bookmarks is secure, hardened, private-by-default in a complete product sense, Wardveil-conformant, Privacy-Shield-conformant, production-ready, or accepted. Those claims require implementation and attributable evidence at the relevant revision.

## Resilience, backup, and recovery

Graceful process shutdown, ordered schema migration, exact migration-history verification, and fail-closed schema readiness are implemented. No Bookmarks database backup, tested restore, destructive-migration recovery, archive recovery, rollback qualification, or Everkeep integration is currently implemented.

## Accessibility and user interface

Firefox popup/options interface source is implemented under `clients/firefox/`, but no Firefox runtime accessibility or Glaze UI acceptance claim is established. No broader Bookmarks web, native desktop, or native mobile interface is currently implemented.

## Build and validation boundary

The repository pins Go `1.27.1` and includes `.github/workflows/validate-go.yml` to validate formatting, module tidiness, `go vet`, unit tests, mandatory PostgreSQL integration tests, and compilation of both the service and migration command on exact pull-request/source revisions. `.github/workflows/security-supply-chain.yml` separately generates and retains a CycloneDX 1.6 JSON SBOM and fails closed on reachable known Go vulnerabilities through pinned `govulncheck`. `.github/workflows/firefox-client.yml` separately validates the Firefox source contract, JavaScript syntax, deterministic packaging, and XPI integrity for Firefox-client changes.

Source presence, passing CI, merge state, release state, deployment state, platform acceptance, production readiness, and Stable qualification are distinct states and must not be conflated.

## Capability validation rule

A future capability may be added here only when source and applicable tests/evidence support the exact claimed scope. Product behavior, platform integration, release, deployment, and production acceptance must remain separately evidenced.
