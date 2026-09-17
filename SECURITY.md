# GoreeCloud Bookmarks — Security

## Current state

GoreeCloud Bookmarks is currently **Experimental** and contains a minimal Go service foundation.

Current implemented security-relevant behavior is deliberately narrow:

- The development listener defaults to loopback (`127.0.0.1:8080`) rather than all interfaces.
- The HTTP server uses bounded read-header, read, write, and idle timeouts.
- SIGINT/SIGTERM trigger bounded graceful shutdown.
- Responses set `X-Content-Type-Options: nosniff` and `Cache-Control: no-store`.
- Readiness fails closed with HTTP `503` while the required Bookmarks data layer does not exist.
- No active credentials or secrets are embedded in the current source.
- The executable service currently uses only the Go standard library at runtime.
- Repository CI uses pinned GitHub Action revisions and validates formatting, vetting, tests, and build.

These controls apply only to the current experimental foundation. They do not establish production hardening, authentication, authorization, TLS termination, database security, Wardveil acceptance, Privacy Shield acceptance, deployment security, or production readiness.

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

The current service does not implement bookmark mutation or content ingestion, so capture/parser/archive/import/search/sync attack surfaces are not yet active Bookmarks runtime capabilities.

The current HTTP foundation must still be treated as untrusted-input-facing code if an operator deliberately binds it beyond loopback. No production exposure is authorized or documented by this repository state.

`GOREECLOUD_BOOKMARKS_LISTEN_ADDR` is a non-secret configuration value. Future credential-bearing configuration must follow GoreeCloud sensitive-information separation requirements and must not be committed to ordinary source history.

## Safe archival requirement

Archived webpages must not execute as trusted applications. The future default archive viewer should prevent active scripts, uncontrolled remote resources, active forms, and automatic execution of downloads. Navigation away from preserved content should be controlled and explicit.

No archive runtime is currently implemented.

## Security authority

Wardveil Security is the GoreeCloud security authority. GoreeCloud Policy may evaluate approved security rules, and GoreeCloud Observability may carry security-relevant operational evidence, but neither replaces Wardveil's security-domain authority.

No Bookmarks-specific Wardveil integration or acceptance evidence is currently verified. The presence of unit tests or CI does not create Wardveil conformance.

## Dependency and build security

The first service foundation has no third-party Go runtime module dependency. This reduces, but does not eliminate, supply-chain risk.

The validation workflow pins `actions/checkout` and `actions/setup-go` to exact Git commit revisions. Build and CI dependencies remain subject to GoreeCloud security-update and vulnerability-management requirements.

Future Go modules, parsers, database drivers, archive libraries, frontend packages, client dependencies, container images, and workflow dependencies must be reviewed for necessity, provenance, licensing, support, and security before they become accepted project dependencies.

## Vulnerability reporting

No Bookmarks-specific public vulnerability-reporting channel has been verified in this repository. Until one is established, repository security work must follow current GoreeCloud security and source-control governance rather than inventing an unsupported contact address.

Do not place reusable secrets, private keys, API tokens, credentials, recovery codes, protected user data, or sensitive environment data in public issue bodies or source commits.

## Security claims

Do not describe GoreeCloud Bookmarks as secure, hardened, Wardveil-conformant, production-ready, or accepted merely because the experimental service builds, tests pass, or health returns `200`. Those claims require implementation and attributable evidence for the exact claimed scope.
