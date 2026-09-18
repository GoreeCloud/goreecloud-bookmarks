# GoreeCloud Bookmarks — Dependency Record

## Status

- **Lifecycle:** Development
- **Purpose:** Record accepted implementation/build/test dependencies introduced by the current PostgreSQL data foundation.
- **Authority boundary:** This record documents repository dependency choices; it does not establish platform acceptance, production deployment, or release qualification.

## Go runtime dependency

### `github.com/jackc/pgx/v5`

- **Selected version:** `v5.11.0`
- **Role:** PostgreSQL protocol/client and connection-pool implementation for the Bookmarks server and migration command.
- **License:** MIT.
- **Selection rationale:** Maintained native Go PostgreSQL driver, current stable v5 line, compatible with the pinned Go 1.27.1 baseline, supports current PostgreSQL versions, and avoids introducing a broader application framework or ORM.
- **Security note:** v5.11.0 follows the 2026 pgx decoder/protocol hardening work and includes additional decoder hardening. Connection strings remain protected configuration.
- **Replacement boundary:** Bookmarks persistence code is isolated under `internal/database/postgres`; application-domain code must not make pgx-specific types part of the external API contract.

Direct/transitive module versions are pinned by `go.mod` and integrity-checked through `go.sum`. `go mod tidy -diff`, `go vet`, tests, and build validation run in repository CI.

## PostgreSQL server baseline

- **Development integration-test baseline:** PostgreSQL `18.6`.
- **CI image:** Docker Official Image `postgres:18.6` using an exact release tag for an ephemeral synthetic test database. It is not a production image or important-state database; production database images remain subject to stronger pinning/qualification requirements.
- **Role:** Authoritative relational/transactional server state for GoreeCloud Bookmarks.
- **Production status:** No production PostgreSQL deployment, host, database, credential, storage path, backup policy, TLS endpoint, or acceptance evidence is established by this repository state.

PostgreSQL minor releases contain bug/security fixes; supported deployments must remain on a currently supported and approved minor release rather than treating the CI baseline as a permanent freeze.

## Migration implementation

The repository does **not** add a third-party migration framework in this phase. Ordered migration SQL is embedded from `migrations/`, and the small repository-owned runner records version/name/SHA-256 checksum in `goreecloud_bookmarks_schema_migrations`, holds an advisory migration lock, and refuses incompatible history.

This choice minimizes the initial dependency surface. It may be revisited if later migration complexity provides a documented reason for a dedicated migration library.

## GitHub Actions dependencies

The validation workflow uses GitHub Actions pinned to exact commit revisions. `.github/dependabot.yml` is configured to propose weekly Go-module and GitHub Actions updates; those proposals still require ordinary review and validation before merge. Build-time dependencies remain separate from application runtime dependencies and are subject to GoreeCloud security-update governance.

## Maintenance rule

Before adding or updating a dependency:

1. Confirm necessity and the narrow responsibility it serves.
2. Verify current support, license compatibility, provenance, and relevant security information.
3. Pin the selected version/revision in the appropriate manifest/workflow.
4. Update this record when the dependency boundary materially changes.
5. Validate the exact candidate revision before merge.
6. Keep capability, security, licensing, Platform Contract, and task records truthful about what the dependency actually enables.
