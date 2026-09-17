# GoreeCloud Bookmarks — Repository Specification

## Status and authority

- **Product:** GoreeCloud Bookmarks
- **Repository:** `GoreeCloud/goreecloud-bookmarks`
- **Lifecycle/state represented here:** Planned product specification; implementation state must be established separately from repository evidence.
- **Canonical project record:** `GoreeCloud/Projects/Project Specification — Bookmarks.md`
- **Related repository:** `GoreeCloud/goreecloud-bookmark-browser-extension`
- **Last reconciled:** 2026-09-17

> This file records intended product scope and architecture. It must not be used as evidence that a planned capability is implemented, deployed, validated, production-ready, or accepted.

## 1. Product vision

GoreeCloud Bookmarks is intended to be a private, cross-platform bookmarking, web-archiving, read-later, annotation, and personal-knowledge service for saving anything worth returning to.

The product principle is:

> **Save anything. Preserve what matters. Find it again instantly.**

The preferred workflow is:

> **Capture instantly → enrich automatically → organize when useful → preserve when important → rediscover when relevant.**

The product should remain simple enough for one-click bookmarking while supporting large research libraries, durable web preservation, offline access, annotations, automation, sharing, and advanced retrieval.

## 2. First-class experiences

The planned service should expose one underlying library through multiple first-class interfaces:

- Deep native integration with GoreeCloud Browser.
- Desktop application for large-library management, research, bulk operations, offline archives, multi-window workflows, and keyboard-driven use.
- Mobile application for rapid capture, share-target integration, offline reading, search, highlights, notes, reminders, and quick organization.
- Web application providing nearly the complete library-management experience from a modern browser.

All interfaces should operate on the same synchronized library while retaining useful local data and offline behavior on installed clients.

## 3. Core content model

A saved item may represent a webpage, article, document, image, video, audio item, uploaded file, note, text snippet, highlight, saved search, web-feed entry, archived page, or manually entered URL.

The shared metadata model should be able to represent, as applicable:

- Title, URL, description, content type, collection, and tags.
- Personal notes and highlights.
- Favorite, pinned, read, archive, and offline state.
- Save, publication, last-visited, and last-modified times.
- Source, author/creator, preview image, extracted text, and language.
- Reminder, privacy level, sharing state, archive versions, link health, and custom properties.

The model should be extensible without forcing unrelated concerns into one monolithic Bookmark object.

## 4. Capture and browser integration

Saving should succeed immediately before optional enrichment begins.

Planned capture paths include:

- One-click quick save of the current page.
- Advanced save with collection, tags, notes, favorite, read-later, reminder, archive, offline, and privacy options.
- Saving selected text with surrounding context and annotation.
- Saving images and their source pages.
- Saving links without opening them first.
- Saving the current tab, selected tabs, all tabs in a window, or an entire browsing session.

GoreeCloud Browser integration is intended to include:

- A built-in bookmark control with visible saved/favorite/archive/offline/shared state.
- A lightweight save panel that opens after the bookmark is already persisted.
- A Bookmarks side panel for recent items, collections, search, favorites, read-later, annotations, related items, and saved-page history.
- Optional address-bar results from the private bookmark library.
- Optional synchronization with conventional browser bookmarks using disabled, import-only, export-only, or two-way modes.
- Deterministic conflict handling and synchronization history for two-way browser-bookmark synchronization.
- A possible future replacement of the conventional browser bookmark manager with GoreeCloud Bookmarks collections.

## 5. Organization

The planned organization model includes:

- A default Inbox for capture-first, organize-later workflows.
- Collections as the primary structured organization mechanism.
- Nested collections for hierarchical organization.
- Unlimited many-to-many tags with autocomplete, suggestions, aliases, descriptions, merging, usage statistics, and bulk operations.
- Smart Collections backed by dynamic saved queries rather than duplicate items.
- Saved searches and reusable filters.

Smart-query filters may include collection, tag, domain, content type, author, save/publication date, read state, favorite state, archive state, sharing state, link health, offline availability, and reminder state, with AND/OR/NOT composition.

## 6. Search and retrieval

Search is a core differentiator and should cover:

- Titles, URLs, descriptions, tags, collection names, notes, and highlights.
- Readable article text, archived pages, documents, captured text, annotations, and extracted metadata when available.
- Visual filters for collection, tag, date, site, type, read/favorite/archive/offline/shared state, and owner where applicable.
- Compact power-user operators such as `tag:research unread:true`.

An optional intelligence layer may add meaning-based retrieval and natural-language queries. It must remain optional, indicate when semantic retrieval affected a result, and never become a dependency for core bookmark functionality.

## 7. Automatic organization and automation

Optional assistance may suggest or derive:

- Tags.
- Collection placement.
- Descriptions and summaries.
- Content type and important topics.
- Related bookmarks.

Automatic organization must support user-selectable modes, remain reversible, and preserve user control.

A rules engine is planned for condition/action automation based on URL, domain, title, content type, tags, source application, and save method. Actions may add/remove tags, move collections, favorite, archive, download offline, mark read-later, set reminders, or set privacy.

## 8. Web preservation

GoreeCloud Bookmarks should optionally preserve content when saved.

Planned archive modes are:

- Metadata only.
- Readable archive.
- Visual snapshot.
- Complete archive sufficient to reproduce as much of the original page as safely practical.

Archive policy should be configurable globally and per collection.

Important bookmarks may support multiple archive versions over time, comparison of versions, and change detection for text, title, images, and metadata.

Archived content must open in an isolated viewer. By default, active scripts should not execute, remote resources should not automatically load, forms should be disabled, navigation should be controlled, and dangerous downloads must not execute automatically.

## 9. Link health and duplicates

Planned library-health states include healthy, redirected, changed, temporarily unavailable, permanently unavailable, blocked, and unknown.

A Library Health view should surface dead links, redirects, duplicates, missing archives, missing metadata, and unreachable pages while retaining archived copies when the source disappears.

Duplicate detection should consider exact URLs, tracking-parameter variants, canonical-equivalent URLs, and potential duplicate content. The user should be able to keep both, merge, replace, move annotations, combine tags, or choose a preferred URL.

Known nonessential tracking parameters should optionally be removed at save time while retaining the original URL internally for provenance when desired.

## 10. Reading, highlights, and notes

Reader mode should provide distraction-free access to saved readable content with font, size, spacing, content-width, theme, progress, table-of-contents, and estimated-reading-time controls.

Highlights should preserve selected text, optional personal notes, creation time, source location, and optional labels. A dedicated Highlights view should support independent search.

Each bookmark should support personal notes, and standalone notes should be possible without a URL. Notes may support basic formatting, lists, links, quotations, attachments, and references to other bookmarks while remaining intentionally lighter than a full note-taking suite.

## 11. Rediscovery and related content

The product should help users recover forgotten value from their library through configurable mechanisms such as:

- Related bookmarks.
- On This Day.
- Forgotten Bookmarks.
- Random Bookmark.
- Related Discovery while reading.
- Weekly Review of newly saved and unprocessed Inbox items.

Related-item processing must remain optional.

## 12. Read Later, favorites, pins, and reminders

Read Later should support unread, reading, completed, and archived states, reading progress, estimated time, length-based sorting, recency sorting, oldest-unread sorting, and offline availability.

Favorites identify important items. Pins keep items or collections prominently visible.

Bookmarks may have one-time or recurring reminders and may be snoozed until a selected date.

## 13. Offline-first client behavior

Installed clients should maintain useful local data rather than acting only as remote wrappers.

Local storage should include, at minimum, bookmark metadata, collections, tags, searchable recent-library data, pending mutations, and content selected for offline availability.

Users should be able to make individual bookmarks, Read Later, entire collections, favorites, and Smart Collections available offline.

Offline changes must synchronize when connectivity returns.

## 14. Synchronization and conflicts

Synchronization should cover bookmarks, collections, tags, notes, highlights, read state, favorites, pins, reminders, archive metadata, and preferences.

Every mutation should carry enough version information to detect conflicts.

The system should automatically merge independent edits when safe, use last-change handling only for fields where it is safe, and preserve multiple versions for ambiguous conflicting edits. User content must not silently disappear because of synchronization.

A sync-history view should expose device, action, object, timestamp, status, and conflict state, with deeper failed-operation inspection for advanced users.

## 15. Sharing and collaboration

All content must be private by default.

Intentional sharing may apply to individual bookmarks, collections, or read-only archives with access modes such as private, specific people, shared group, anyone-with-link, or public.

Sharing must not implicitly expose parent collections, unrelated tags, private notes, or internal metadata.

Shared collections may eventually support viewer, contributor, editor, manager, and owner roles plus meaningful activity history.

Optional public collections may publish curated resource lists, reading lists, bibliographies, or documentation indexes with configurable title, description, cover, sort order, and layout.

## 16. Feed following

Collections may optionally follow compatible web feeds. New entries may appear temporarily or be bookmarked automatically, and users may apply rules, mark entries read, or filter by keyword.

This feature should bridge bookmarking and information discovery without turning Bookmarks into a mandatory general-purpose feed reader.

## 17. Views and bulk operations

Planned collection views include Cards, Compact, List, Headlines, and Table, with per-collection view preferences.

Bulk operations should support selecting large numbers of bookmarks to move, tag, untag, favorite, archive, mark read, download offline, change privacy, export, or delete. Large operations should be safe, show progress, and be resumable where practical.

## 18. Import, export, and portability

Import is a major onboarding capability. It should preserve as much hierarchy, title, URL, tag, creation-date, and note data as source formats provide, while previewing item counts, detected duplicates, proposed collections, and errors before commit.

Exports should support browser-compatible bookmark formats, structured machine-readable data, human-readable archives, selected items, selected collections, and complete-account export. Where permitted, complete exports may include metadata, collections, tags, notes, highlights, archive metadata, and preserved content.

Users must not be trapped in GoreeCloud Bookmarks.

## 19. Privacy and encryption

Bookmark libraries may reveal highly sensitive interests and behavior, so the product should apply a higher privacy standard than ordinary application metadata.

Core privacy principles are:

- Private by default.
- No advertising dependency.
- No behavioral-profiling requirement.
- No unnecessary external requests.
- Minimum metadata collection.
- User-controlled retention.
- Portable data.
- Explicit sharing.
- Clear deletion behavior.
- Optional intelligence rather than mandatory analysis.

Two planned privacy modes are:

### Standard Private Library

Transport encryption and protected server-side storage while permitting server-side search indexing, archive processing, automatic organization, and content extraction.

### Private Vault

Optional client-side encryption for selected data before it leaves the device, with the server primarily synchronizing opaque ciphertext. Tradeoffs must be explicit because some server-side search, processing, organization, or archival capabilities may be limited, device-local, delayed, or unavailable.

Bookmarks or collections may carry privacy levels such as Normal private, Sensitive, Private Vault, or Shared. Sensitive content may disable previews, external metadata fetches, remote analysis, or automatic archival unless the user explicitly enables them.

## 20. Proposed service architecture

The planned backend is modular rather than one monolithic bookmarking API.

Proposed services include:

- Bookmarks API for bookmarks, collections, tags, notes, highlights, sharing, and permissions.
- Capture Service for accepting saves and scheduling enrichment.
- Metadata Service for titles, descriptions, preview imagery, and structured metadata.
- Archive Service for preserved versions.
- Content Extraction Service for readable representations.
- Search Service for bookmark and archive indexing.
- Sync Service for client change coordination.
- Automation Service for user-defined rules.
- Link Health Service for URL checks and health state.
- Notification Service for reminders and relevant alerts.

The capture path must be asynchronous after durable save:

1. Create the bookmark record immediately.
2. Confirm success to the user.
3. Begin enrichment in the background.
4. Retrieve metadata.
5. Extract readable content.
6. Create an archive when enabled.
7. Update the search index.
8. Generate organization suggestions when enabled.
9. Synchronize enriched state to clients.

Failure of an enrichment stage must not delete or invalidate the original saved bookmark.

## 21. Suggested core data objects

The initial backend may center around:

- User
- Device
- Bookmark
- Collection
- CollectionMembership
- Tag
- BookmarkTag
- Note
- Highlight
- Archive
- ArchiveVersion
- Reminder
- Share
- Permission
- SyncEvent
- AutomationRule
- LinkHealthRecord
- SavedSearch

## 22. Interface and visual requirements

GoreeCloud Bookmarks should feel like a modern library rather than a traditional browser-settings screen.

The interface should prioritize saved content, preview imagery, collections, search, reading progress, tags, and quick actions while keeping heavy management controls secondary until needed.

The product should follow the current approved GoreeCloud visual system / Glaze UI contract applicable to its lifecycle. Planned treatments include layered translucent surfaces, soft depth, contextual blur, functional collection accents, rich preview cards, smooth transitions, responsive layouts, clear status indicators, consistent iconography, and light/dark modes.

Color must have functional meaning rather than acting as decoration alone.

## 23. Keyboard and accessibility requirements

Desktop and web clients should support extensive keyboard navigation and a command palette for major actions.

Accessibility must be part of the initial architecture, including:

- Full keyboard navigation.
- Screen-reader labels.
- High-contrast compatibility.
- Reduced motion.
- Scalable text.
- Logical focus order.
- Large touch targets.
- Non-color status indicators.
- Accessible reader typography.

## 24. Administration, storage, backup, and observability

Multi-user administration may expose users, storage usage, archive storage, search-index health, processing queues, failed archive jobs, link-health tasks, sync status, sharing policy, registration policy, and retention settings.

Administrative authority must not automatically grant content access.

Users should be able to inspect bookmark count, archive count, archive storage, offline storage, largest collections, and largest saved files, and remove old versions, duplicate archives, previews, offline copies, or failed captures without deleting the bookmark itself.

Backup and recovery should cover bookmarks, collections, tags, notes, highlights, sharing configuration, archive metadata, saved searches, and rules. Preserved webpage data may use a separate lifecycle due to size. Recovery must be testable.

Operational health should cover API availability, capture/archive queues, failed captures, search-index delay, sync errors, storage health, link-check backlog, and notification failures.

## 25. Planned navigation

Primary navigation may include:

- Home
- Inbox
- Library
- Read Later
- Highlights
- Collections
- Tags
- Search
- Shared
- Archive

Secondary utilities may include Library Health, Import / Export, and Settings.

## 26. Delivery phases

### Phase 1 — Bookmarking Foundation

Planned foundation: bookmark creation, collections, nested collections, tags, favorites, search, metadata extraction, import/export, GoreeCloud Browser integration, web app, and account synchronization.

### Phase 2 — Native Clients and Offline

Planned: desktop app, mobile app, share integration, local cache, offline collections/reading, conflict-resilient synchronization, and device management.

### Phase 3 — Preservation and Reading

Planned: readable extraction, page archiving, archive viewer/versioning, reader mode, highlights, notes, reading progress, and link-health monitoring.

### Phase 4 — Automation and Discovery

Planned: Smart Collections, saved searches, automation rules, duplicate detection, tracking cleanup, related bookmarks, rediscovery, reminders, and web feeds.

### Phase 5 — Optional Intelligence

Planned optional capabilities: tag suggestions, collection suggestions, summaries, meaning-based search, natural-language library queries, and related-content discovery.

### Phase 6 — Collaboration

Planned: shared collections, collaboration, fine-grained permissions, public collections, activity history, and collection-level administration.

## 27. Recommended MVP boundary

The planned MVP is:

1. Save URLs.
2. GoreeCloud Browser one-click integration.
3. Automatic title, description, and preview extraction.
4. Collections.
5. Nested collections.
6. Tags.
7. Favorites.
8. Inbox.
9. Read Later.
10. Search.
11. Basic full-content extraction.
12. Web application.
13. Desktop application.
14. Mobile application.
15. Offline-capable local cache.
16. Cross-device synchronization.
17. Import.
18. Export.
19. Basic backup support.
20. Private-by-default access control.

This MVP is a planning boundary, not a statement of current implementation.

## 28. Long-term differentiation objectives

The intended differentiation is the combination of:

- Deep GoreeCloud Browser integration.
- Private personal web archiving.
- Offline-first access.
- Strong rediscovery.
- Unified search across URLs, article text, notes, highlights, and preserved content.
- Automation that remains optional and understandable.
- Explicit privacy controls for sensitive material.
- Strong import/export/backup and data ownership.
- One synchronized library across browser, desktop, mobile, and web.

## 29. Long-term vision

At maturity, GoreeCloud Bookmarks may become a personal web-memory layer for the GoreeCloud ecosystem: a user should be able to search what they remember and retrieve the original bookmark, archived copy, related pages, notes, highlights, previous versions, associated collections, and other relevant saved material.

The product should sit between traditional bookmarks, read-later software, web archives, research libraries, and personal knowledge management while retaining one-click capture as a fundamental interaction.

## 30. Current implementation boundary

As of this reconciliation, the repository was verified to contain only its initial minimal `README.md` on `main`. This specification therefore establishes planned repository-coupled requirements only.

No product capability described above may be promoted to current/implemented status until authoritative repository, test, runtime, release, and acceptance evidence supports that claim.
