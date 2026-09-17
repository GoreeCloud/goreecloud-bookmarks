# GoreeCloud Bookmarks — Privacy Policy

## Status

This repository currently represents a **Concept-stage** product. No GoreeCloud Bookmarks runtime, hosted service, application telemetry pipeline, account database, archive processor, or production data-processing operation is verified as implemented.

This document therefore records the required privacy model for future implementation and the current no-runtime state. It must be updated before any production processing begins.

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

Neither mode is currently implemented or accepted.

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

## Administration

Administrative authority and content access must remain separate concepts. A service administrator must not automatically gain permission to inspect private bookmark contents solely because they operate the deployment.

## Retention and deletion

Future implementation must define retention and deletion behavior for bookmark metadata, archives, archive versions, local caches, search indexes, derived previews, backups, collaboration records, and operational evidence. Deletion claims must reflect actual downstream retention and backup behavior.

## Current data-processing statement

At the verified repository state represented by this policy, no Bookmarks-specific application or service processing of user bookmark data is established. Repository hosting and development-platform processing are governed separately from the future Bookmarks product runtime.