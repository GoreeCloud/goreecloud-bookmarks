# GoreeCloud Bookmarks — Privacy Policy

## Status

GoreeCloud Bookmarks is currently a **Development** project with a Go/PostgreSQL service and persistence foundation.

The implemented service exposes bounded health/readiness behavior and an internal PostgreSQL Bookmark persistence layer used through source/tests. It does **not** expose an authenticated bookmark-domain HTTP API, account system, archive processing, search indexing, synchronization, telemetry collection, sharing, or a production user-data-processing path.

This document therefore records both the required privacy model for future Bookmarks functionality and the narrow current processing boundary. It must continue to be updated before user bookmark data is processed by newly implemented subsystems.

## Current Development processing boundary

The current service foundation:

- can store the implemented Bookmark-row subset when PostgreSQL is explicitly configured, but does not expose a user-facing bookmark ingestion API;
- does not implement a user/account database;
- does not implement archive capture or content extraction;
- does not implement search queries or indexes;
- does not implement synchronization or sharing;
- does not implement application telemetry;
- uses `github.com/jackc/pgx/v5` `v5.11.0` for PostgreSQL connectivity; no bookmark content is sent to that library as an external network service beyond the configured PostgreSQL connection;
- exposes health/readiness JSON containing only bounded service state;
- logs service startup/shutdown/error information without a designed user-content logging path.

The optional `GOREECLOUD_BOOKMARKS_LISTEN_ADDR` variable configures a listener address and is not a secret. `GOREECLOUD_BOOKMARKS_DATABASE_URL` may contain database credentials and must remain protected external configuration outside source control under GoreeCloud sensitive-information requirements.

Repository hosting, GitHub Actions, and development-platform processing are separate from a future Bookmarks user-data runtime and remain governed by their applicable systems and policies.

## Core privacy requirements

GoreeCloud Bookmarks must be private by default. Bookmark libraries can reveal sensitive interests, research, health concerns, finances, relationships, work activity, browsing behavior, and other personal information; the product must treat saved-content data as sensitive user-controlled information.

Future implementation must:

- Avoid advertising dependency and behavioral profiling requirements.
- Minimize metadata collection and unnecessary external requests.
- Use explicit sharing rather than implicit exposure.
- Give users clear retention, deletion, export, and portability controls.
- Keep optional intelligence and semantic analysis optional.
- Avoid silently transmitting sensitive bookmark content to external processors.
- Respect applicable Privacy Shield requirements and GoreeCloud Policy decisions.
- Preserve provenance and purpose information for material processing where required.

## Planned privacy modes

### Standard Private Library

Transport and stored-data protections may permit approved server-side indexing, extraction, archival, and organization processing.

### Private Vault

Selected data may be encrypted before leaving the user's device so the server primarily synchronizes opaque encrypted data. Capabilities that require server-side plaintext processing may be limited, delayed, device-local, or unavailable.

Neither privacy mode is currently implemented or accepted.

## Planned privacy levels

- Normal private.
- Sensitive.
- Private Vault.
- Shared.

Sensitive content should be able to disable or restrict preview generation, external metadata fetching, automatic archival, or remote analysis unless the user explicitly allows the processing.

## Sharing

All content should default to private. Future sharing must not expose parent collections, unrelated tags, private notes, internal metadata, or other library content merely because one item or collection is shared.

## Archival privacy

Capturing a page must not grant archived content trusted-application status. Future archive viewers should disable active scripts and uncontrolled remote-resource loading by default and must apply privacy requirements to any enrichment or fetching process.

## Telemetry and observability

Operational evidence must minimize private content. GoreeCloud Observability integration must preserve useful health and diagnostic information without converting bookmark contents, URLs, annotations, search terms, or other sensitive data into unnecessary telemetry.

The current health/readiness implementation is not Bookmarks-specific Observability acceptance evidence.

## Administration

Administrative authority and content access must remain separate concepts. A service administrator must not automatically gain permission to inspect private bookmark contents solely because they operate the deployment.

## Retention and deletion

Future implementation must define retention and deletion behavior for bookmark metadata, archives, archive versions, local caches, search indexes, derived previews, backups, collaboration records, and operational evidence. Deletion claims must reflect actual downstream retention and backup behavior.

## Current data-processing statement

At the Development source state represented by this policy, GoreeCloud Bookmarks has an executable service and internal PostgreSQL bookmark persistence foundation but no authenticated end-user bookmark-data plane or production deployment. No claim of complete Privacy Shield integration, production privacy acceptance, or Private Vault behavior is established.
