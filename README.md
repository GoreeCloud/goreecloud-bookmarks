# GoreeCloud Bookmarks

GoreeCloud Bookmarks is the first-party bookmark collection, organization, preservation, and synchronization application for GoreeCloud. The long-term application is original GoreeCloud-owned native software. The inherited Linkwarden application tree is transitional and remains only where it is still required for continuity, compatibility, migration, feature preservation, and upstream security maintenance.

## Lifecycle status

**Native migration in progress — not Stable.**

This repository currently contains both the inherited Linkwarden-based compatibility surface and a growing native GoreeCloud Bookmarks implementation. Passing source tests or retaining the inherited runtime does not make the application Stable and does not authorize retirement of migration or recovery paths.

Stable remains blocked until the native application has completed the applicable functional, data-migration, accessibility, runtime, backup/restore, and real-environment acceptance work and has current validated integration with Glaze UI, Wardveil Security, Privacy Shield, and Everkeep.

## Native implementation

The GoreeCloud-owned native application lives under `native/` and is implemented as a Go service.

Current source areas include:

- `native/cmd/bookmarksd` — native service entry point.
- `native/internal/bookmarks` — first-party bookmark records and application behavior.
- `native/internal/collections` — native collection organization boundaries.
- `native/internal/identity` — application identity/session integration boundaries.
- `native/internal/storage` — native persistence interfaces and state handling.
- `native/internal/webui` — GoreeCloud-owned web presentation.
- `native/internal/syncstate` — GoreeCloud Sync capability, signing, authenticated submission, retrieval, deletion, and validation logic.

The native Sync client advertises its application capability explicitly, requires exact negotiated schema conformance, bounds record and continuation identifiers, requires an authenticated session before transport, and preserves Privacy Shield data minimization by keeping deletion tombstones free of application payload.

## Transitional Linkwarden boundary

The inherited Linkwarden application is a migration dependency, not the target GoreeCloud product architecture. It may continue to provide retained functionality while equivalent or improved native behavior is implemented, migrated, and accepted.

Inherited capabilities must not disappear silently during migration. Data ownership, collection relationships, preserved content, import/export behavior, and other retained workflows require explicit migration or retirement decisions before the inherited runtime is removed.

Required upstream copyright, license, source-availability, and attribution obligations remain in force for retained Linkwarden-derived code. Upstream branding and hosted-service marketing do not define GoreeCloud Bookmarks identity or lifecycle status.

## Mandatory platform gates

GoreeCloud Bookmarks must continuously conform to the current approved contracts for:

- **Glaze UI** — first-party responsive presentation, accessibility, interaction, appearance, and adaptive layouts.
- **Wardveil Security** — safe content handling, application protection state, diagnostics, and security integration where applicable.
- **Privacy Shield** — data minimization, privacy-preserving synchronization, retained-data controls, and application privacy boundaries.
- **Everkeep** — preserved-content continuity, portability, backup/restore, migration, and recovery requirements.

Missing, materially outdated, or unvalidated mandatory integration keeps the application non-Stable.

## Validation

The repository currently uses both native and transitional validation while migration is in progress. Relevant workflows include:

- GoreeCloud Bookmarks Native Foundation
- GoreeCloud Bookmarks CI
- Linkwarden Playwright Tests

The native foundation gate validates GoreeCloud-owned Go source and contracts. The broader inherited application tests remain useful migration and regression evidence while the compatibility surface exists. Neither class of source validation substitutes for production migration, backup/restore, recovery, platform-system acceptance, or Stable approval.

For the native Go module:

```bash
cd native
go test ./...
go build ./cmd/bookmarksd
```

These commands validate the native source locally; they do not deploy or migrate production data.

## Branding

Canonical GoreeCloud Bookmarks artwork is owned by `GoreeCloud/goreecloud-branding-assets`. See `BRANDING.md` for the repository consumer boundary. Synchronized artwork retained here is not the branding authority.

## Repository structure

- `native/` — original GoreeCloud Bookmarks implementation.
- `apps/` and the remaining inherited web-application paths — transitional Linkwarden-derived runtime and compatibility surface where still present.
- `BRANDING.md` — canonical branding ownership and synchronization rule.
- `LICENSE.md` — licensing terms for inherited and retained source.

## Development rule

New GoreeCloud Bookmarks product work belongs in the native application architecture unless a narrowly scoped inherited or third-party foundation is technically necessary. The inherited complete-product architecture must not be expanded as the long-term GoreeCloud implementation.

Changes must preserve migration, user data, security, privacy, recoverability, and required upstream licensing while moving application-defining behavior into GoreeCloud-owned native source.
