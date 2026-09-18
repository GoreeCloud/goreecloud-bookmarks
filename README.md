# GoreeCloud Bookmarks

GoreeCloud Bookmarks is the GoreeCloud bookmarking, web-archiving, read-later, annotation, and personal web-memory project.

> **Release lifecycle:** Development  
> **Current capability state:** Go/PostgreSQL service and persistence foundation. No authenticated end-user bookmarking workflow is implemented.

## Product direction

The intended product combines one-click capture, deep GoreeCloud Browser integration, collections and tags, full-content search, web preservation, Reader workflows, highlights and notes, offline-first clients, conflict-tolerant synchronization, import/export, private sharing, automation, rediscovery, and optional intelligence.

The guiding product concept is:

> **Save anything. Preserve what matters. Find it again instantly.**

The preferred workflow is:

> **Capture instantly → enrich automatically → organize when useful → preserve when important → rediscover when relevant.**

## Current Development implementation

The repository contains the current Development service/data foundation:

- Go `1.27.1` toolchain baseline.
- Standard-library HTTP service entry point under `cmd/bookmarks`.
- `GET /api/v1/health` returning bounded process-health JSON.
- `GET /api/v1/ready` failing closed unless configured PostgreSQL is reachable and its migration history exactly matches the embedded schema set.
- Loopback-only development default `127.0.0.1:8080`, with optional `GOREECLOUD_BOOKMARKS_LISTEN_ADDR` override.
- HTTP timeouts and graceful SIGINT/SIGTERM shutdown.
- PostgreSQL connectivity through `pgx/v5` `v5.11.0`.
- Explicit ordered/checksummed schema migrations through `cmd/bookmarks-migrate`; service startup never auto-migrates.
- Initial bookmark persistence schema plus internal owner-scoped create/read operations.
- Unit tests plus mandatory PostgreSQL integration tests and GitHub validation for formatting, module tidiness, vetting, tests, and builds.

No authenticated bookmark-domain HTTP endpoint, GoreeCloud Identity integration, collections, tags, search, synchronization, archival, web UI, desktop/mobile client, supported Docker deployment, backup/restore qualification, or accepted Integral Platform System integration is implemented yet.

## Development

The pinned development toolchain is Go `1.27.1`.

Validate the source with:

```bash
go version
gofmt -w ./cmd ./internal
go vet ./...
go test ./...
go build -o ./build/goreecloud-bookmarks ./cmd/bookmarks
go build -o ./build/goreecloud-bookmarks-migrate ./cmd/bookmarks-migrate
```

The repository validation workflow separately verifies formatting without modifying files, vets, tests, and builds the exact candidate revision.

For local experimental execution:

```bash
go run ./cmd/bookmarks
```

The default listener is `127.0.0.1:8080`. An explicit local override may be provided with `GOREECLOUD_BOOKMARKS_LISTEN_ADDR`. This environment variable is not a secret and the service does not automatically load a `.env` file.

Current engineering probes:

```text
GET /api/v1/health
GET /api/v1/ready
```

Without `GOREECLOUD_BOOKMARKS_DATABASE_URL`, readiness remains fail-closed. When a protected database URL is supplied, run `go run ./cmd/bookmarks-migrate` explicitly before expecting readiness to pass. The current PostgreSQL 18.6 baseline is validated only as an Experimental/Development integration-test target, not as a production deployment guarantee.

This service is not a production deployment target at the current lifecycle state.

## Repository documentation

- [`SPECIFICATIONS.md`](SPECIFICATIONS.md) — planned product requirements.
- [`ARCHITECTURE.md`](ARCHITECTURE.md) — selected architecture plus current implementation boundary.
- [`CAPABILITIES.md`](CAPABILITIES.md) — current verified capability state.
- [`FEATURES.md`](FEATURES.md) — current experimental foundation and planned feature catalog.
- [`FEATURE-ROADMAP.md`](FEATURE-ROADMAP.md) — phased implementation roadmap and current milestone state.
- [`BENEFITS.md`](BENEFITS.md) — intended user and platform benefits.
- [`COMPETITIVE-OBJECTIVES.md`](COMPETITIVE-OBJECTIVES.md) — differentiation objectives and guardrails.
- [`BRANDING.md`](BRANDING.md) — product naming and visual direction.
- [`USER-MANUAL.md`](USER-MANUAL.md) — current availability plus planned user workflows.
- [`PRIVACY POLICY.md`](PRIVACY%20POLICY.md) — required privacy model and current experimental-processing boundary.
- [`SECURITY.md`](SECURITY.md) — security requirements and current experimental security boundary.
- [`NOTES.md`](NOTES.md) — verified repository notes and remaining decisions.
- [`docs/api/README.md`](docs/api/README.md) — planned v1 API behavior, compatibility, synchronization, and privacy contract.
- [`docs/api/openapi.yaml`](docs/api/openapi.yaml) — OpenAPI 3.2.1 planning contract for `/api/v1/`; contract presence is not endpoint implementation evidence.
- [`docs/data-model.md`](docs/data-model.md) — logical resource model and storage/authority boundaries.
- [`docs/migrations.md`](docs/migrations.md) — versioning, migration, compatibility, and rollback rules.
- [`LICENSE`](LICENSE) — current `AGPL-3.0-or-later` fallback rights notice.
- [`LICENSE-DECISION.md`](LICENSE-DECISION.md) — licensing decision record and supersession rule.
- [`goreecloud.platform.yaml`](goreecloud.platform.yaml) — Platform Contract `0.4` declaration with unverified conformance.

The canonical GoreeCloud product record is `GoreeCloud/Projects/Project Specification — Bookmarks.md` in the authorized GoreeCloud documentation system.

## Selected architecture direction

The broader implementation direction in [`ARCHITECTURE.md`](ARCHITECTURE.md) selects:

- one authoritative modular GoreeCloud Bookmarks service rather than premature microservices;
- Go for the server/API/background-processing implementation;
- TypeScript for the first-class web client;
- PostgreSQL for authoritative relational server state and initial full-text search;
- SQLite for installed-client offline state and pending mutations;
- REST-style HTTPS/JSON under `/api/v1/` with an OpenAPI contract;
- Kotlin/native Android, Swift/native Apple clients, and Rust + GTK 4 for Linux desktop;
- WARC 1.1 / ISO 28500:2017 for complete web-preservation capture containers;
- a Docker Compose self-hosted server stack with dedicated PostgreSQL and persistent archive storage.

The Go service, PostgreSQL connectivity/migration foundation, and initial bookmark relation are implemented today. The remaining selections remain architecture direction until their source and evidence exist.

## Platform governance

Platform Contract `0.4` evaluates exactly nine Integral Platform Systems: GoreeCloud Manager, Privacy Shield, Wardveil Security, Everkeep, Glaze UI, GoreeCloud Mesh, GoreeCloud Identity, GoreeCloud Policy, and GoreeCloud Observability.

GoreeCloud Sync is separately governed and is not a tenth Integral Platform System.

All nine Bookmarks integrations remain blocked/unaccepted. The current health and readiness routes are not, by themselves, a GoreeCloud Observability or Manager integration.

## Related project

`GoreeCloud/goreecloud-bookmark-browser-extension` is a related capture surface. Native GoreeCloud Browser integration remains broader planned work and must preserve Bookmarks as the authoritative bookmark domain.

## Licensing

The active GoreeCloud Software Licensing Policy applies `AGPL-3.0-or-later` as the default fallback because no Bookmarks-specific superseding decision is currently recorded. This repository records that state through [`LICENSE`](LICENSE) and [`LICENSE-DECISION.md`](LICENSE-DECISION.md).

## Repository workflow

Material changes use short-lived purpose-specific branches and pull requests into `main`. `.github/PULL_REQUEST_TEMPLATE.md` preserves scope, validation, and exact-revision traceability.

Default-branch protection remains required by GoreeCloud governance but is not currently verified active in GitHub. That repository-administration blocker is tracked separately.

## Truthful status rule

Implemented source, passing CI, merge state, release state, deployment state, platform acceptance, production readiness, and Stable qualification are separate claims. Planned capabilities, documentation, manifests, or a successful process start must never be used to imply a more mature state than evidence supports.
