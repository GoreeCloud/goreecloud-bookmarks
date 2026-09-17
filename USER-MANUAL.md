# GoreeCloud Bookmarks — User Manual

## Current availability

GoreeCloud Bookmarks is currently documented at the **Concept** lifecycle stage. No production application, service, web interface, desktop client, mobile client, or supported end-user workflow is currently verified in this repository.

This manual therefore describes the intended interaction model and clearly marks it as planned. It must be revised against real implementation before being used as installation, operation, support, or production-use instructions.

## Intended core workflow

The planned default workflow is:

1. Save a page, link, selection, image, file, or note.
2. Confirm capture immediately.
3. Allow optional enrichment such as metadata extraction, readable-content extraction, archival, search indexing, and organization suggestions to continue independently.
4. Organize later using Inbox, collections, tags, Smart Collections, favorites, Read Later, or search.
5. Retrieve the saved item through search, browsing, related-item discovery, reminders, or rediscovery features.

## Planned primary areas

- **Home** — recent, pinned, reading, reminder, and rediscovery surfaces.
- **Inbox** — unprocessed captures.
- **Library** — the complete saved-item collection.
- **Read Later** — reading queue and progress.
- **Highlights** — cross-library highlight search and review.
- **Collections** — structured organization.
- **Tags** — flexible categorization.
- **Search** — metadata and full-content retrieval.
- **Shared** — explicitly shared items and collections.
- **Archive** — preserved copies and versions.
- **Library Health** — dead links, redirects, missing archives, duplicates, and metadata issues.
- **Import / Export** — migration and portability tools.
- **Settings** — preferences, privacy, offline, archival, automation, synchronization, and processing controls.

## Planned capture behavior

A normal save should require one action and should persist the bookmark before optional enrichment begins. An advanced save panel may then expose collection, tags, notes, Read Later, favorite, reminder, archive, offline, and privacy controls.

Planned GoreeCloud Browser workflows include saving the current page, selected text, links, images, multiple tabs, and browsing sessions without requiring the Bookmarks application to replace the active page.

## Planned privacy behavior

Items are intended to be private by default. Sharing should require explicit user action. Higher-privacy modes may restrict preview generation, metadata fetching, archival, remote analysis, or server-side indexing.

No current privacy implementation should be inferred from these planned controls. See `PRIVACY POLICY.md` and `CAPABILITIES.md` for current truth.

## Planned offline behavior

Installed clients are intended to maintain local metadata, pending changes, searchable recent-library data, and user-selected offline content. Offline changes should synchronize after connectivity returns without silently discarding ambiguous conflicts.

## Planned import and export

Import should preview detected items, hierarchy, duplicates, and errors before committing a large migration. Export should support portable browser-compatible forms, structured data, selected scopes, and complete-account export where applicable.

## Support and troubleshooting

There is currently no supported Bookmarks runtime to troubleshoot. Repository documentation and planning issues must not be represented as end-user product support incidents.

When implementation begins, this manual must be expanded with verified installation, sign-in, capture, organization, offline, synchronization, sharing, backup/recovery, update, troubleshooting, accessibility, and removal/retirement procedures for each supported client and deployment mode.