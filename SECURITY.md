# GoreeCloud Bookmarks — Security

## Current state

GoreeCloud Bookmarks is currently a Concept-stage repository with no verified application or service implementation. No runtime security control, security certification, Wardveil integration, deployment hardening, or production acceptance is established by this document.

## Security objectives for future implementation

The Bookmarks system must treat captured web content, imported files, metadata fetches, synchronization payloads, sharing links, automation rules, search indexes, and user annotations as potentially sensitive and, where applicable, potentially hostile input.

Future implementation should include:

- Strong authentication and authorization appropriate to the supported deployment model.
- Least-privilege service identities and administrative separation.
- Input validation and output encoding.
- Safe URL retrieval and protections against server-side request forgery and unsafe redirects.
- Archive sandboxing and active-content suppression by default.
- Safe parsing and extraction of remote content and uploaded/imported files.
- Protection against path traversal, archive bombs, parser abuse, malicious HTML, hostile metadata, and unsafe downloads.
- Secure synchronization with integrity, replay, conflict, and authorization controls.
- Sharing tokens/links with explicit scope, expiry/revocation capability where appropriate, and no parent-scope leakage.
- Secrets kept outside source code and ordinary checked-in configuration.
- Dependency and supply-chain review.
- Secure update, migration, backup/restore, and rollback procedures.
- Privacy-aware logging and operational evidence.
- Security testing for archive, import/export, authorization, sync, automation, and sharing boundaries.

## Safe archival requirement

Archived webpages must not execute as trusted applications. The default archive viewer should prevent active scripts, uncontrolled remote resources, active forms, and automatic execution of downloads. Navigation away from preserved content should be controlled and explicit.

## Security authority

Wardveil Security is the GoreeCloud security authority. GoreeCloud Policy may evaluate approved security rules, and GoreeCloud Observability may carry security-relevant operational evidence, but neither replaces Wardveil's security-domain authority.

No Bookmarks-specific Wardveil integration or acceptance evidence is currently verified.

## Vulnerability reporting

No Bookmarks-specific public vulnerability-reporting channel has been verified in this repository. Until one is established, repository security work must follow current GoreeCloud security and source-control governance rather than inventing an unsupported contact address.

Do not place reusable secrets, private keys, API tokens, credentials, recovery codes, or sensitive environment data in public issue bodies or source commits.

## Security claims

Do not describe GoreeCloud Bookmarks as secure, hardened, Wardveil-conformant, production-ready, or accepted based solely on this policy. Security claims require implementation and attributable evidence at the exact relevant revision.