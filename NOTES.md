# GoreeCloud Bookmarks — Notes

## Current verified repository state

- Product lifecycle represented by this implementation candidate: **Development**.
- Repository: `GoreeCloud/goreecloud-bookmarks`.
- Current implementation: Go/PostgreSQL service/persistence foundation plus internal retry-safe Bookmark capture service; no authenticated end-user bookmark HTTP data plane is implemented.
- Go baseline: `1.27.1`, pinned in `go.mod`.
- Current direct Go runtime dependency: `github.com/jackc/pgx/v5` `v5.11.0`, with provenance/license baseline in `docs/dependencies.md`.
- Implemented HTTP routes: `GET /api/v1/health` and `GET /api/v1/ready` only.
- Current readiness behavior: fail closed unless configured PostgreSQL is reachable and the exact embedded migration history is current and untampered.
- Development listener default: `127.0.0.1:8080`; optional non-secret override through `GOREECLOUD_BOOKMARKS_LISTEN_ADDR`.
- Repository validation: `.github/workflows/validate-go.yml` pins its GitHub Actions revisions and validates the exact Go toolchain, formatting, module tidiness, `go vet`, unit/integration tests against ephemeral PostgreSQL `18.6`, and service/migration-command builds.
- Canonical product specification: `GoreeCloud/Projects/Project Specification — Bookmarks.md`.
- Repository product specification: `SPECIFICATIONS.md`.
- Repository architecture record: `ARCHITECTURE.md`.
- Planned API contract: `docs/api/README.md` and `docs/api/openapi.yaml`.
- Logical data model: `docs/data-model.md`.
- Migration/compatibility rules: `docs/migrations.md`.
- Related capture repository: `GoreeCloud/goreecloud-bookmark-browser-extension`.
- Platform Contract baseline: `0.4`.
- Integral Platform System model: exactly nine systems.
- Stable Glaze UI consumer target referenced by Platform Contract `0.4`: `1.5.1`.
- Governing repository license: `AGPL-3.0-or-later` through the GoreeCloud default fallback; no Bookmarks-specific superseding license decision is currently recorded.
- Default branch `main` remains unprotected in verified GitHub state; protection remains a repository-administration blocker.

## Important authority boundaries

- Source implementation, CI validation, merge state, release state, deployment state, platform acceptance, and production acceptance are separate states.
- `CAPABILITIES.md` records current implemented capability scope.
- `FEATURE-ROADMAP.md` records milestone/phase work but is not release evidence.
- `ARCHITECTURE.md` records selected architecture and the boundary between the implemented foundation and still-planned components.
- `docs/api/openapi.yaml` remains the planned v1 wire contract. Contract entries beyond health/readiness must not be treated as implemented endpoints.
- `docs/data-model.md` is the logical ownership/model record. Migration version `1` implements the initial PostgreSQL Bookmark subset and migration version `2` implements owner-scoped create idempotency; the rest of the model remains planned.
- `docs/migrations.md` defines evolution/rollback behavior and now records the implemented explicit migration runner; destructive-migration recovery and production rollback remain unqualified.
- GoreeCloud Bookmarks is authoritative for the implemented PostgreSQL bookmark-domain rows; GoreeCloud Identity remains authoritative for identity/credential data.
- GoreeCloud Sync is separately governed and is not an Integral Platform System.
- GoreeCloud Browser integration must be reconciled with Browser authority and the separate bookmark-browser-extension repository.
- Optional intelligence must not become a dependency for core bookmarking.
- The fallback license may be superseded only by an authorized Bookmarks-specific decision reconciled across repository and canonical records.

## Implemented Development foundation

The current implemented foundation consists of:

- `cmd/bookmarks/main.go` — service entry point, local listener configuration, HTTP timeouts, and graceful shutdown.
- `internal/httpapi/handler.go` — bounded health/readiness HTTP surface.
- `internal/database/postgres/` — PostgreSQL connectivity, exact schema-state checks, migration execution, and owner-scoped internal bookmark persistence.
- `migrations/000001_initial_bookmarks.sql` — initial authoritative Bookmark table.
- `migrations/000002_bookmark_create_idempotency.sql` — owner-scoped durable create-idempotency state.
- `internal/bookmarks/` — Bookmark domain/capture service, validation/defaulting, opaque ID generation, request hashing, and retry-safe store contract.
- `cmd/bookmarks-migrate/main.go` — explicit migration entry point; service startup does not auto-migrate.
- `docs/dependencies.md` — current third-party dependency provenance/licensing baseline.
- `internal/httpapi/handler_test.go` — health, readiness, method, and response-behavior tests.
- `cmd/bookmarks/main_test.go` — listen-address tests.
- `go.mod` / `go.sum` — pinned Go `1.27.1` baseline plus reviewed `pgx/v5` dependency graph.
- `.github/workflows/validate-go.yml` — exact-candidate validation.

The foundation deliberately does not expose a bookmark-domain HTTP API, invent local credentials, add placeholder clients, create a supported Docker deployment, or claim platform acceptance merely to increase apparent implementation breadth.

## Current PostgreSQL checkpoint

- PostgreSQL migration version `1` creates the initial `bookmarks` relation; version `2` adds the `bookmark_create_idempotency` relation and composite owner/bookmark reference needed for atomic retry-safe capture.
- Migration history is append-only/checksummed and guarded by a PostgreSQL advisory lock.
- The runtime rejects missing, newer-than-binary, tampered, or structurally invalid schema state through fail-closed readiness.
- Internal persistence supports owner-scoped bookmark reads and atomic idempotent creates. The same owner/key/request replays the original bookmark; different content under the same owner/key fails with an idempotency conflict; concurrent identical retries commit only one bookmark.
- Active database configuration is external through `GOREECLOUD_BOOKMARKS_DATABASE_URL`; no active database credential is stored in Git.
- CI uses an ephemeral database whose name must end in `_test` before destructive reset code can run. The integration suite also verifies sequential and concurrent idempotent retries, conflict handling, owner isolation, migration v2 idempotence, and required-schema-object checks.

## Selected broader architecture direction

The unimplemented broader direction remains:

- One authoritative modular Bookmarks service rather than premature microservices.
- Go server for API, synchronization coordination, capture orchestration, jobs, and long-running service responsibilities.
- TypeScript web client.
- Kotlin/native Android client.
- Swift/native Apple clients when implemented.
- Rust + GTK 4 Linux desktop client.
- PostgreSQL authoritative relational state and initial full-text search.
- SQLite installed-client offline state and pending mutations.
- REST/JSON under `/api/v1/` with OpenAPI, explicit authorization, revisions/ETags, and retry-safe mutation rules.
- PostgreSQL-backed durable background jobs before a separate queue platform is justified.
- Persistent archive payload storage outside the container writable layer.
- WARC 1.1 / ISO 28500:2017 for complete preservation captures.
- Docker Compose as the initial self-hosted server deployment boundary.
- Database-native backup plus independent archive/data protection and tested recovery.

## Decisions intentionally deferred

- SQLite client implementation/tooling and any future additional database/migration dependencies.
- Exact TypeScript UI framework and frontend dependency set.
- Exact WARC capture/replay library and compression implementation.
- Exact Private Vault cryptographic protocol, key lifecycle, and recovery model.
- Exact production host, hostname, DNS, Caddy route, NetBird policy, secrets, ports, resource limits, storage paths, and operational schedules.
- Exact backup frequency and retention values.
- Exact release signing, package publication, and client-update infrastructure.
- Whether measured scale or isolation requirements justify extracting search, archive processing, jobs, or another module into a separate deployed service.
- Whether a Bookmarks-specific license should supersede the current `AGPL-3.0-or-later` fallback after applicable review.

## Documentation maintenance

Every new implementation change must reconcile `CAPABILITIES.md`, `FEATURES.md`, `USER-MANUAL.md`, `PRIVACY POLICY.md`, `SECURITY.md`, `goreecloud.platform.yaml`, `ARCHITECTURE.md`, the relevant `docs/` contracts, and the feature roadmap when the verified meaning of those records changes.
