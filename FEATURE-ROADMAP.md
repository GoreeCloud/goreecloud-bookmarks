# GoreeCloud Bookmarks — Feature Roadmap

## Authority and status

**Release lifecycle:** Development  
**Roadmap status:** Legacy migration control — non-authoritative.  
**Canonical product scope:** `PROJECT-SPECIFICATIONS.md`; significant project history is in `PROJECT-RECORD.md`.

This legacy roadmap is retained temporarily as migration input for the separately governed repository-native feature/changelog migration. It is not current project-specification, release, deployment, platform-acceptance, production-readiness, or Stable evidence.

## Milestone 0 — Governed repository and executable foundation

Verified completed foundation work:

- [x] Required root documentation baseline established.
- [x] Platform Contract `0.4` declaration established without overstating conformance.
- [x] Truthful `CAPABILITIES.md` reporting established.
- [x] `AGPL-3.0-or-later` fallback recorded in `LICENSE` and `LICENSE-DECISION.md`.
- [x] Initial architecture and technology direction documented.
- [x] Planned OpenAPI v1, logical data, and migration/compatibility contracts documented.
- [x] Minimal Go service source foundation added with health/readiness behavior and unit tests.
- [x] Go validation workflow added for exact-toolchain, formatting, module-tidy, vet, test, and build checks.
- [x] PostgreSQL connectivity, explicit migration runner, migration history, initial Bookmark schema, and owner-scoped internal persistence implemented.
- [x] Mandatory PostgreSQL integration validation added, including newer-schema and migration-history-tamper rejection.
- [x] Exact-revision CycloneDX SBOM generation and pinned Go vulnerability scanning added as a fail-closed CI gate.

Foundation work that remains open:

- [ ] Activate required protection/rulesets for `main`; current GitHub state remains unprotected and the connected integration does not expose the necessary administration write.
- [ ] Remove merged short-lived branches when a branch-delete capability is available and dependencies are verified absent.
- [ ] Maintain exact-revision validation and repository-native documentation reconciliation for future changes.
- [x] Add proposal-only dependency update configuration for Go modules and GitHub Actions through `.github/dependabot.yml`; generated proposals still require review/test/merge.

The current Development foundation is intentionally not ready for production or end-user use. Readiness can pass only for a configured exact-current PostgreSQL schema; deployment, authentication, backup/recovery, user APIs, and subsequent acceptance gates remain open.

## Phase 1 — Bookmarking foundation

All Phase 1 product capabilities remain open:

- [ ] Durable one-click bookmark creation.
- [ ] Inbox.
- [ ] Collections and nested collections.
- [ ] Tags and bookmark-tag relationships.
- [ ] Favorites and Read Later.
- [ ] Basic metadata extraction.
- [ ] Search across metadata and available extracted text.
- [ ] Import and export.
- [ ] Initial web application.
- [ ] Account synchronization with conflict detection.
- [ ] Initial GoreeCloud Browser capture integration.

## Phase 2 — Native clients and offline

- [ ] Desktop application.
- [ ] Mobile application and share target.
- [ ] Local searchable cache.
- [ ] Offline-selected content.
- [ ] Conflict-resilient offline mutation queues.
- [ ] Device and synchronization-history views.

## Phase 3 — Preservation and reading

- [ ] Readable-content extraction.
- [ ] Metadata, readable, visual, and complete archive modes.
- [ ] Isolated archive viewer.
- [ ] Archive version history and comparison.
- [ ] Reader mode and reading progress.
- [ ] Highlights, annotations, and notes.
- [ ] Link-health monitoring and Library Health.

## Phase 4 — Automation and discovery

- [ ] Smart Collections and saved searches.
- [ ] Automation rules.
- [ ] Duplicate and canonical-URL detection.
- [ ] Tracking-parameter cleanup with provenance retention.
- [ ] Related bookmarks and rediscovery.
- [ ] Reminders and snooze.
- [ ] Optional feed following.

## Phase 5 — Optional intelligence

- [ ] Tag and collection suggestions.
- [ ] Summaries.
- [ ] Meaning-based search.
- [ ] Natural-language library queries.
- [ ] Related-content discovery.

Core bookmarking, organization, search, archival, import/export, and synchronization must remain usable without optional intelligence.

## Phase 6 — Collaboration

- [ ] Explicit private sharing.
- [ ] Shared collections.
- [ ] Viewer, Contributor, Editor, Manager, and Owner roles.
- [ ] Collaborative activity history.
- [ ] Public curated collections with explicit disclosure controls.

## Cross-cutting gates for every phase

Each phase must address applicable privacy, security, continuity, accessibility, testing, synchronization, migration, rollback, operational-health, policy, documentation, licensing, dependency, and exact-revision evidence requirements.

All nine Integral Platform Systems must be evaluated according to current authoritative governance. GoreeCloud Sync remains separately governed.

## Lifecycle progression

Current verified direction:

**Development** → Release Candidate → Stable

Development is appropriate to the current active implementation/validation state. Movement to Release Candidate requires an identifiable candidate plus the applicable production-readiness, platform, security, privacy, recovery, compatibility, accessibility, and acceptance evidence.


## Migration control

Do not synchronize this file with Google Drive. Remaining obligations must be reconciled into the mandatory repository-native `IMPLEMENTED-FEATURES.md`, `PLANNED-FEATURES.md`, and `CHANGELOGS.md` records before this legacy roadmap is removed.

Current project requirements belong in `PROJECT-SPECIFICATIONS.md`; significant project history belongs in `PROJECT-RECORD.md`. Completion claims require accepted repository evidence.
