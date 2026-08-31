# GoreeCloud Bookmarks Features

This file records current functionality without treating source scaffolding, CI success, migration compatibility, or development-only runtime behavior as production or Stable acceptance.

## Implemented native foundations

- GoreeCloud-owned Go service under `native/`.
- Owner-scoped bookmark create, list, read, update, and delete behavior.
- HTTP/HTTPS bookmark URL validation and fragment removal.
- Bounded bookmark title, note, tag, and request-body handling.
- Deterministic owner-scoped bookmark listing.
- Native collection hierarchy store with owner isolation.
- Collection create, list, read, and parent-move behavior.
- Duplicate sibling-name rejection and hierarchy-cycle protection.
- Development HTTP collection routes for create/list/read/move.
- Strict JSON mutation decoding with unknown-field and trailing-value rejection.
- Explicit identity resolver boundary with fail-closed default runtime identity.
- Development-only owner-header identity seam restricted to the memory-development repository mode.
- Native GoreeCloud Sync capability, signed record submission, payload-free tombstones, authenticated transport submission, bounded retrieval/pagination, and canonical proof-vector tests.
- GoreeCloud-owned native web presentation foundation.

## Development-only / partial capabilities

- `memory-development` bookmark repository.
- In-memory collection state.
- Development-header identity testing.
- PostgreSQL repository contract/scaffolding without an accepted production runtime driver.
- Collection-to-bookmark assignment source foundation without the full user-facing/persistent HTTP product surface.
- Transitional Linkwarden runtime retained for migration and feature continuity.
- Sync source integration without completed production identity/storage/deployment acceptance.

## Planned or incomplete capabilities

- Production GoreeCloud Identity authentication and authorization.
- Accepted persistent bookmark and collection storage.
- Durable collection/bookmark assignment and collection deletion semantics.
- Sharing, invitations, collaboration, and permission management where approved.
- Preserved-page/archive content migration and native handling.
- Complete import/export and data portability.
- Linkwarden user-data and relationship migration.
- Backup, restore, rollback, and disaster-recovery acceptance through Everkeep.
- Full Privacy Shield runtime integration and privacy controls.
- Full Wardveil Security integration for fetched/preserved external content and application protection state.
- Current-Stable Glaze UI product-wide acceptance.
- Production GoreeCloud Sync operation and conflict policy.
- Representative accessibility, performance, browser/device, and target-environment acceptance.
- Production deployment and Stable lifecycle promotion.

## Feature governance

A capability is represented only at the acceptance level supported by source and evidence. Development mode, source implementation, CI validation, migration compatibility, production deployment, production acceptance, and Stable promotion remain distinct states.
