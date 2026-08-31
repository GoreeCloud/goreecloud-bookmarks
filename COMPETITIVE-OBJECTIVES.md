# GoreeCloud Bookmarks Competitive Objectives

## Purpose

GoreeCloud Bookmarks should meet mature bookmark-manager expectations while remaining a first-party GoreeCloud application with stronger ownership, privacy boundaries, continuity, and platform integration.

## Benchmark products

Relevant benchmarks include Linkwarden, Raindrop.io, Pocket-style read-later workflows, browser-native bookmark managers, and self-hosted bookmark/preservation tools. A benchmark is a source of product expectations, not a permanent dependency or branding authority.

## Capabilities to match

- Fast bookmark capture and editing.
- Useful titles, notes, and tags.
- Nested collection/folder organization.
- Search and filtering across a personal library.
- Import and export.
- Reliable browser integration.
- Cross-device synchronization.
- Responsive, accessible library navigation.
- Preserved or archived content where supported.
- Clear handling of duplicate or invalid records.

## Capabilities to exceed

### GoreeCloud-owned architecture

The canonical bookmark, collection, identity, synchronization, privacy, security, and recovery boundaries should remain GoreeCloud-controlled rather than inheriting the full architecture of a third-party application.

### Privacy by default

Owner isolation, minimal retained data, explicit Sync boundaries, payload-free deletion tombstones, and Privacy Shield governance should be normal behavior rather than optional add-ons.

### Continuity and migration evidence

Backup, restore, rollback, export, and inherited-runtime migration should be treated as product-quality requirements through Everkeep rather than emergency-only operational tasks.

### Failure honesty

Development-only identity, memory storage, incomplete persistence, migration gaps, or unavailable platform integrations must be represented explicitly instead of being hidden behind a superficially functional UI.

### Platform integration

Bookmarks should integrate cleanly with GoreeCloud Browser, Search, Sync, Identity, Everkeep, Privacy Shield, Wardveil Security, Glaze UI, and Mesh without duplicating those systems' authority.

## Capabilities intentionally rejected

- Advertising or sponsored bookmark placement.
- Behavioral profiling for monetization.
- Hidden cross-user data exposure.
- Treating a development identity header as production authentication.
- Silent loss of inherited user data during migration.
- Permanent coupling of GoreeCloud's canonical data model to Linkwarden internals.
- Production claims based only on CI or an in-memory development service.

## Long-term objective

The goal is a native GoreeCloud library that is private, portable, recoverable, accessible, well synchronized, and independently maintainable while preserving user data and useful inherited workflows throughout migration.
