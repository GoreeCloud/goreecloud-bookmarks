# GoreeCloud Bookmarks — Project Record

**Repository:** `GoreeCloud/goreecloud-bookmarks`  
**Lifecycle:** Development  
**Migration baseline:** `899f2294257cf47a3846a2af9d10c1e2c84477bb`  
**Current governing license:** `AGPL-3.0-or-later` fallback  
**Canonical authority:** This file is the repository-local project record once accepted on the default branch.

## 2026-09-17 — Repository and governance foundation

The repository was created for GoreeCloud Bookmarks and established the initial documentation, architecture, Platform Contract 0.4 declaration, licensing fallback, API/data/migration planning contracts, and Development lifecycle boundaries.

The product was intentionally defined as bookmarking, web archiving, read-later, annotation, and personal web-memory software while keeping planned product scope distinct from implementation evidence.

## 2026-09-17 — Executable service foundation

Commit `b78f5d096589ff6eb3edb8932773f8b06041dc4b` added the first executable Go service foundation with bounded health/readiness behavior, tests, and exact Go 1.27.1 validation.

This was an engineering foundation, not an end-user bookmarking workflow.

## 2026-09-18 — PostgreSQL data foundation

PR #7 was squash-merged as signed commit `effd884b84802de93d1040908807c7193c7b0f72`.

It established:
- PostgreSQL connectivity through pgx/v5;
- explicit ordered/checksummed migrations;
- advisory locking;
- the first Bookmark relation;
- internal owner-scoped create/read persistence;
- database-backed fail-closed readiness;
- newer-schema and tampered-history rejection; and
- mandatory PostgreSQL integration validation.

Ordinary service startup does not auto-run migrations.

## 2026-09-18 — Firefox client ownership migration

Commit `243e9979bc97495ec1d57f4240131d31bf2b3014` moved the Firefox client into `clients/firefox/` inside the Bookmarks application repository.

Commit `f9105abd20da7693123515b668c39b814affe4bb` recorded completed ownership migration and removal of the shared transitional source copy.

The Firefox client remains a source candidate. The required application-side capture endpoint is not implemented on authoritative `main`, so this migration did not create an accepted end-to-end bookmarking workflow or a separate bookmark-data authority.

## 2026-09-19 — Supply-chain qualification gate

PR #10 merged as `899f2294257cf47a3846a2af9d10c1e2c84477bb`.

It added exact-revision CycloneDX 1.6 SBOM generation and fail-closed Go vulnerability scanning, including remediation of the then-reachable vulnerability recorded in the merge commit.

This is Development supply-chain evidence, not Production Acceptance or Stable qualification.

## Licensing

The repository currently applies `AGPL-3.0-or-later` through the GoreeCloud default fallback recorded in `LICENSE` and `LICENSE-DECISION.md`.

No Bookmarks-specific superseding license decision is currently recorded.

## Glaze UI authority reconciliation

The active Drive source and Platform Contract metadata referenced Glaze UI 1.5.1 as Stable.

Live verification of `GoreeCloud/glaze-ui` establishes V1.6 / `1.6.0` as current Official Stable and V1.7 as Development-only.

The project-governance migration updates Bookmarks references to V1.6 while retaining `applicable-blocked` / unverified Bookmarks-specific conformance.

## Drive project-specification migration

The active Drive source, **Project Specification — Bookmarks** (source ID `1SWEt2FFD8ATFJPn3CKTlnZO5E0MELR7P`), contains 69 numbered planned capability sections plus a September 2026 verified-state checkpoint.

The migration:
- consolidates those requirements with the former root `SPECIFICATIONS.md`;
- preserves newer GitHub implementation evidence over stale Drive checkpoints;
- creates root `PROJECT-SPECIFICATIONS.md`;
- creates root `PROJECT-RECORD.md`;
- removes active Drive project-specification authority from repository documentation;
- retires the old `SPECIFICATIONS.md` filename after incorporation;
- retires `ROADMAP-SYNC.md`, which required synchronized Drive roadmap updates;
- converts `FEATURE-ROADMAP.md` to non-authoritative migration input; and
- preserves archived Drive Bookmarks sources only as historical comparison material.

Drive source deletion remains blocked until this migration is accepted on authoritative `main`, read back, and verified without unresolved discrepancy.

## Separate feature/changelog migration

The mandatory repository-native `IMPLEMENTED-FEATURES.md`, `PLANNED-FEATURES.md`, and `CHANGELOGS.md` set is not yet established.

Legacy roadmap material therefore remains migration input only until that separate migration is accepted.

## Current acceptance boundary

At migration baseline, no authenticated end-to-end bookmark-management workflow, production deployment, accepted Integral Platform System integration, Release Candidate status, Production Acceptance, or Stable qualification exists.

Open dependency-update PRs remain ordinary candidates and do not alter this project record until accepted.

## Ongoing maintenance

Update this record for significant architecture, client ownership, licensing, migration, security/privacy, data-model, deployment, lifecycle, repository, major acceptance, deprecation, or retirement events.
