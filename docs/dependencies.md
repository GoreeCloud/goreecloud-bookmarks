# GoreeCloud Bookmarks — Dependency Record

## Status

- **Lifecycle:** Development
- **Purpose:** Record accepted runtime, build, test, and security-validation dependencies used by the current Development foundation.
- **Authority boundary:** This record documents repository dependency choices; it does not establish platform acceptance, production deployment, or release qualification.

## Go runtime dependency

### `github.com/jackc/pgx/v5`

- **Selected version:** `v5.11.0`
- **Role:** PostgreSQL protocol/client and connection-pool implementation for the Bookmarks server and migration command.
- **License:** MIT.
- **Selection rationale:** Maintained native Go PostgreSQL driver, current stable v5 line, compatible with the pinned Go 1.27.1 baseline, supports current PostgreSQL versions, and avoids introducing a broader application framework or ORM.
- **Security note:** v5.11.0 follows the 2026 pgx decoder/protocol hardening work and includes additional decoder hardening. Connection strings remain protected configuration.
- **Replacement boundary:** Bookmarks persistence code is isolated under `internal/database/postgres`; application-domain code must not make pgx-specific types part of the external API contract.

Direct/transitive module versions are pinned by `go.mod` and integrity-checked through `go.sum`. `golang.org/x/text` is explicitly held at `v0.39.0` as an indirect security floor because exact-revision `govulncheck` demonstrated a reachable `GO-2026-5970` path from Bookmarks through `pgxpool`; versions before `v0.39.0` are affected. `go mod tidy -diff`, `go vet`, tests, vulnerability scanning, and build validation run in repository CI.

## PostgreSQL server baseline

- **Development integration-test baseline:** PostgreSQL `18.6`.
- **CI image:** Docker Official Image `postgres:18.6@sha256:4ef4dbc939d61acea57712655ddb4b4ab27419c913f94cca0cd57cb3ea3c2280`, pinned by release tag and observed registry digest for the ephemeral synthetic test database. This does not qualify a production database image or deployment.
- **Role:** Authoritative relational/transactional server state for GoreeCloud Bookmarks.
- **Production status:** No production PostgreSQL deployment, host, database, credential, storage path, backup policy, TLS endpoint, or acceptance evidence is established by this repository state.

PostgreSQL minor releases contain bug/security fixes; supported deployments must remain on a currently supported and approved minor release rather than treating the CI baseline as a permanent freeze.

## Migration implementation

The repository does **not** add a third-party migration framework in this phase. Ordered migration SQL is embedded from `migrations/`, and the small repository-owned runner records version/name/SHA-256 checksum in `goreecloud_bookmarks_schema_migrations`, holds an advisory migration lock, and refuses incompatible history.

This choice minimizes the initial dependency surface. It may be revisited if later migration complexity provides a documented reason for a dedicated migration library.

## GitHub Actions dependencies

The validation workflows use GitHub Actions pinned to exact commit revisions. `.github/dependabot.yml` is configured to propose weekly Go-module and GitHub Actions updates; those proposals still require ordinary review and validation before merge. Build-time dependencies remain separate from application runtime dependencies and are subject to GoreeCloud security-update governance.

`.github/workflows/security-supply-chain.yml` additionally uses `actions/upload-artifact` at exact commit `043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (release `v7.0.1`) to retain the exact-revision CycloneDX SBOM as CI evidence. This is a CI-only GitHub-maintained dependency and does not become an application runtime dependency.

## Security and SBOM tooling

### `golang.org/x/vuln/cmd/govulncheck`

- **Selected version:** `v1.8.0`.
- **Role:** Exact-revision Go vulnerability reachability scan in CI.
- **License:** BSD-3-Clause.
- **Scope:** CI/security validation only; it is not linked into the Bookmarks runtime.
- **Security behavior:** The tool version is pinned while its vulnerability database remains current by design. A scanner/tool/database failure or a reachable known vulnerability causes the workflow to fail rather than being converted into a pass.

### `github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod`

- **Selected version:** `v1.10.0`.
- **Role:** Generate a machine-readable CycloneDX JSON SBOM for the Go module at the exact candidate revision.
- **License:** Apache-2.0.
- **Scope:** CI/security validation only; it is not linked into the Bookmarks runtime.
- **Output contract:** CycloneDX `1.6` JSON including Go standard-library and module dependency inventory, retained as a GitHub Actions artifact named with the exact Git commit SHA.

The security workflow invokes both Go security tools by exact module version through the pinned Go toolchain. It does not use `@latest`, floating container tags, or unpinned third-party actions.

The current `pgx/v5 v5.11.0` baseline is newer than the patched version floors for the reviewed 2026 pgx advisories affecting versions before `v5.9.0` and `v5.9.2`. Exact-revision automated vulnerability scanning remains required because a one-time advisory review does not establish ongoing safety.

## Maintenance rule

Before adding or updating a dependency:

1. Confirm necessity and the narrow responsibility it serves.
2. Verify current support, license compatibility, provenance, and relevant security information.
3. Pin the selected version/revision in the appropriate manifest/workflow.
4. Update this record when the dependency boundary materially changes.
5. Validate the exact candidate revision before merge.
6. Keep capability, security, licensing, Platform Contract, and task records truthful about what the dependency actually enables.
