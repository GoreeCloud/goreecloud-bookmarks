# GoreeCloud Bookmarks — Security

## Current state

GoreeCloud Bookmarks is currently **Development** and contains a Go/PostgreSQL service and persistence foundation.

Current implemented security-relevant behavior is deliberately narrow:

- The development listener defaults to loopback (`127.0.0.1:8080`) rather than all interfaces.
- The HTTP server uses bounded read-header, read, write, and idle timeouts.
- SIGINT/SIGTERM trigger bounded graceful shutdown.
- Responses set `X-Content-Type-Options: nosniff` and `Cache-Control: no-store`.
- Readiness fails closed for absent, unreachable, migration-required, newer-than-binary, tampered, or invalid PostgreSQL schema state.
- No active credentials or secrets are embedded in the current source.
- The direct database runtime dependency is `github.com/jackc/pgx/v5` `v5.11.0`; provenance/licensing is recorded in `docs/dependencies.md`.
- Active database URLs remain external protected configuration and are not logged or committed.
- Repository CI uses pinned GitHub Action revisions and validates formatting, module tidiness, vetting, unit/integration tests against ephemeral PostgreSQL, and builds.
- A separate fail-closed supply-chain workflow generates an exact-revision CycloneDX 1.6 SBOM and runs `govulncheck` against reachable Go vulnerabilities using pinned scanner/tool versions.

These controls apply only to the current Development foundation. They do not establish production hardening, authentication, authorization, TLS termination, database deployment security, Wardveil acceptance, Privacy Shield acceptance, backup/recovery qualification, or production readiness.

## Security objectives for future implementation

The Bookmarks system must treat captured web content, imported files, metadata fetches, synchronization payloads, sharing links, automation rules, search indexes, and user annotations as potentially sensitive and, where applicable, hostile input.

Future implementation must address, as applicable:

- Strong authentication and authorization appropriate to the supported deployment model.
- Least-privilege service identities and administrative separation.
- Input validation and output encoding.
- Safe URL retrieval and protections against server-side request forgery and unsafe redirects.
- Archive sandboxing and active-content suppression by default.
- Safe parsing and extraction of remote content and uploaded/imported files.
- Protection against path traversal, archive bombs, parser abuse, malicious HTML, hostile metadata, and unsafe downloads.
- Secure synchronization with integrity, replay, conflict, and authorization controls.
- Sharing links/tokens with explicit scope and expiry/revocation where appropriate.
- Secrets kept outside source code and ordinary checked-in configuration.
- Dependency provenance, vulnerability, and supply-chain review.
- Secure update, migration, backup/restore, and rollback procedures.
- Privacy-aware logging and operational evidence.
- Security testing for archive, import/export, authorization, sync, automation, and sharing boundaries.

## Current attack-surface boundary

The current service does not expose bookmark mutation/content ingestion over HTTP, so capture/parser/archive/import/search/sync attack surfaces are not yet active user-facing capabilities. The internal PostgreSQL layer can persist test/development bookmark rows when explicitly configured.

The current HTTP foundation must still be treated as untrusted-input-facing code if an operator deliberately binds it beyond loopback. No production exposure is authorized or documented by this repository state.

`GOREECLOUD_BOOKMARKS_LISTEN_ADDR` is a non-secret configuration value. `GOREECLOUD_BOOKMARKS_DATABASE_URL` may contain credentials and is protected configuration; active values must not be committed to ordinary source history.

## Safe archival requirement

Archived webpages must not execute as trusted applications. The future default archive viewer should prevent active scripts, uncontrolled remote resources, active forms, and automatic execution of downloads. Navigation away from preserved content should be controlled and explicit.

No archive runtime is currently implemented.

## Security authority

Wardveil Security is the GoreeCloud security authority. GoreeCloud Policy may evaluate approved security rules, and GoreeCloud Observability may carry security-relevant operational evidence, but neither replaces Wardveil's security-domain authority.

No Bookmarks-specific Wardveil integration or acceptance evidence is currently verified. The presence of unit tests or CI does not create Wardveil conformance.

## Dependency and build security

The current direct Go runtime dependency is `github.com/jackc/pgx/v5` `v5.11.0`. Its MIT license/provenance baseline and transitive dependency set are recorded through `docs/dependencies.md`, `go.mod`, and `go.sum`. Dependency maintenance and vulnerability review are continuing obligations.

The validation workflow pins `actions/checkout` and `actions/setup-go` to exact Git commit revisions. `.github/workflows/security-supply-chain.yml` additionally pins `govulncheck` `v1.8.0`, `cyclonedx-gomod` `v1.10.0`, and `actions/upload-artifact` `v7.0.1` by exact commit. The workflow generates a CycloneDX 1.6 JSON SBOM, validates its basic structure, retains it as a commit-SHA-named Actions artifact, and fails when `govulncheck` reports reachable known Go vulnerabilities or cannot complete successfully. Build and CI dependencies remain subject to GoreeCloud security-update and vulnerability-management requirements.

Future Go modules, parsers, archive libraries, frontend packages, client dependencies, production container images, and workflow dependencies must be reviewed for necessity, provenance, licensing, support, and security before they become accepted project dependencies.

## Vulnerability reporting

No Bookmarks-specific public vulnerability-reporting channel has been verified in this repository. Until one is established, repository security work must follow current GoreeCloud security and source-control governance rather than inventing an unsupported contact address.

Do not place reusable secrets, private keys, API tokens, credentials, recovery codes, protected user data, or sensitive environment data in public issue bodies or source commits.

## Security claims

Do not describe GoreeCloud Bookmarks as secure, hardened, Wardveil-conformant, production-ready, or accepted merely because the Development service/database foundation builds, tests pass, or health/readiness succeeds in CI. Those claims require implementation and attributable evidence for the exact claimed scope.
