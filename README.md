# GoreeCloud Bookmarks

GoreeCloud Bookmarks is the planned GoreeCloud bookmarking, web-archiving, read-later, annotation, and personal web-memory platform.

> **Release lifecycle:** Concept  
> **Current capability state:** Documentation only. No GoreeCloud Bookmarks application or service implementation is currently verified.

## Product direction

The intended product combines one-click capture, deep GoreeCloud Browser integration, collections and tags, full-content search, web preservation, Reader workflows, highlights and notes, offline-first clients, conflict-tolerant synchronization, import/export, private sharing, automation, rediscovery, and optional intelligence.

The guiding product concept is:

> **Save anything. Preserve what matters. Find it again instantly.**

The preferred workflow is:

> **Capture instantly → enrich automatically → organize when useful → preserve when important → rediscover when relevant.**

## Repository documentation

- [`SPECIFICATIONS.md`](SPECIFICATIONS.md) — planned product requirements and architecture.
- [`ARCHITECTURE.md`](ARCHITECTURE.md) — selected implementation architecture and technology direction; not runtime evidence.
- [`CAPABILITIES.md`](CAPABILITIES.md) — current verified capability state.
- [`FEATURES.md`](FEATURES.md) — planned feature catalog.
- [`FEATURE-ROADMAP.md`](FEATURE-ROADMAP.md) — phased implementation roadmap.
- [`BENEFITS.md`](BENEFITS.md) — intended user and platform benefits.
- [`COMPETITIVE-OBJECTIVES.md`](COMPETITIVE-OBJECTIVES.md) — product differentiation objectives and guardrails.
- [`BRANDING.md`](BRANDING.md) — product naming and visual-direction guidance.
- [`USER-MANUAL.md`](USER-MANUAL.md) — Concept-stage intended usage model; not production instructions.
- [`PRIVACY POLICY.md`](PRIVACY%20POLICY.md) — required privacy model and current no-runtime processing state.
- [`SECURITY.md`](SECURITY.md) — security requirements and current no-runtime state.
- [`NOTES.md`](NOTES.md) — verified repository notes and remaining open decisions.
- [`LICENSE`](LICENSE) — current repository rights notice applying the GoreeCloud fallback `AGPL-3.0-or-later` license.
- [`LICENSE-DECISION.md`](LICENSE-DECISION.md) — licensing decision record and supersession rule.
- [`goreecloud.platform.yaml`](goreecloud.platform.yaml) — GoreeCloud Platform Contract `0.4` declaration with unverified conformance.

The canonical GoreeCloud product record is `GoreeCloud/Projects/Project Specification — Bookmarks.md` in the authorized GoreeCloud documentation system.

## Selected architecture direction

The Concept-stage implementation direction is documented in [`ARCHITECTURE.md`](ARCHITECTURE.md). In summary, it selects:

- one authoritative modular GoreeCloud Bookmarks service rather than premature microservices;
- Go for the initial server/API/background-processing implementation;
- TypeScript for the first-class web client;
- PostgreSQL for authoritative relational server state and initial full-text search;
- SQLite for installed-client offline state and pending mutations;
- REST-style HTTPS/JSON under `/api/v1/` with an OpenAPI contract;
- Kotlin/native Android, Swift/native Apple clients, and Rust + GTK 4 for the Linux desktop direction;
- WARC 1.1 / ISO 28500:2017 for complete web-preservation capture containers;
- a Docker Compose self-hosted server stack with dedicated PostgreSQL and persistent archive storage.

These are architecture selections only. They do not establish that any implementation, deployment, synchronization path, archive pipeline, native client, or acceptance evidence currently exists.

## Platform governance

The current Platform Contract baseline is `0.4`, which evaluates exactly nine Integral Platform Systems: GoreeCloud Manager, Privacy Shield, Wardveil Security, Everkeep, Glaze UI, GoreeCloud Mesh, GoreeCloud Identity, GoreeCloud Policy, and GoreeCloud Observability.

GoreeCloud Sync is separately governed and is not a tenth Integral Platform System.

The current contract reference target for Glaze UI is `1.5.1`. No Bookmarks-specific platform-system implementation or acceptance evidence is currently verified.

## Related project

The separate `GoreeCloud/goreecloud-bookmark-browser-extension` repository represents a related browser-extension capture surface. Native GoreeCloud Browser integration remains part of the broader Bookmarks direction and must be reconciled with Browser authority and implementation state rather than creating conflicting bookmark/capture ownership.

## Implementation status

Authoritative repository `main` contains the Concept-stage documentation and governance baseline, but no application source, server source, client source, build system, tests, deployment configuration, runtime evidence, or release artifact establishing a GoreeCloud Bookmarks product implementation.

Repository documentation, architecture decisions, manifests, roadmap entries, and governance records must not be treated as runtime implementation or acceptance evidence.

## Licensing

The active GoreeCloud Software Licensing Policy establishes `AGPL-3.0-or-later` as the default fallback when no more appropriate project-specific license has been selected. No Bookmarks-specific superseding decision is currently recorded, so this repository applies that fallback through [`LICENSE`](LICENSE).

A future authorized Bookmarks-specific license decision may supersede the fallback. Any such change must update the repository and applicable canonical GoreeCloud records through the governed process and account for third-party compatibility and previously distributed code where applicable.

## Repository workflow

Material repository changes should use short-lived purpose-specific branches and pull requests into `main`. The repository includes `.github/PULL_REQUEST_TEMPLATE.md` to preserve scope, status, validation, and exact-revision traceability.

Default-branch protection remains required by GoreeCloud governance. Its final GitHub administrative configuration is tracked separately until authoritative GitHub state verifies that protection is active.

## Truthful status rule

Planned capabilities, architecture decisions, roadmap entries, specifications, mockups, manifests, metadata, or documentation presence must never be represented as implemented, accepted, production-ready, or Stable without authoritative implementation and evidence.