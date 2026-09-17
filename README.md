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
- [`CAPABILITIES.md`](CAPABILITIES.md) — current verified capability state.
- [`FEATURES.md`](FEATURES.md) — planned feature catalog.
- [`FEATURE-ROADMAP.md`](FEATURE-ROADMAP.md) — phased implementation roadmap.
- [`BENEFITS.md`](BENEFITS.md) — intended user and platform benefits.
- [`COMPETITIVE-OBJECTIVES.md`](COMPETITIVE-OBJECTIVES.md) — product differentiation objectives and guardrails.
- [`BRANDING.md`](BRANDING.md) — product naming and visual-direction guidance.
- [`USER-MANUAL.md`](USER-MANUAL.md) — Concept-stage intended usage model; not production instructions.
- [`PRIVACY POLICY.md`](PRIVACY%20POLICY.md) — required privacy model and current no-runtime processing state.
- [`SECURITY.md`](SECURITY.md) — security requirements and current no-runtime state.
- [`NOTES.md`](NOTES.md) — verified repository notes and open architecture decisions.
- [`goreecloud.platform.yaml`](goreecloud.platform.yaml) — GoreeCloud Platform Contract `0.4` declaration with unverified conformance.

The canonical GoreeCloud product record is `GoreeCloud/Projects/Project Specification — Bookmarks.md` in the authorized GoreeCloud documentation system.

## Platform governance

The current Platform Contract baseline is `0.4`, which evaluates exactly nine Integral Platform Systems: GoreeCloud Manager, Privacy Shield, Wardveil Security, Everkeep, Glaze UI, GoreeCloud Mesh, GoreeCloud Identity, GoreeCloud Policy, and GoreeCloud Observability.

GoreeCloud Sync is separately governed and is not a tenth Integral Platform System.

The current contract reference target for Glaze UI is `1.5.1`. No Bookmarks-specific platform-system implementation or acceptance evidence is currently verified.

## Related project

The separate `GoreeCloud/goreecloud-bookmark-browser-extension` repository represents a related browser-extension capture surface. Native GoreeCloud Browser integration remains part of the broader Bookmarks direction and must be reconciled with Browser authority and implementation state rather than creating conflicting bookmark/capture ownership.

## Implementation status

Authoritative `main` before the current governance-baseline branch is commit:

`bfe1de2287493cb959507e7417cfa957faf1323d`

That revision contains the documentation foundation merged through PR #1, but no application source, server source, client source, build system, tests, deployment configuration, runtime evidence, or release artifact establishing a Bookmarks product implementation.

The current governance branch adds the required documentation baseline and Platform Contract declaration. Those changes do not become authoritative `main` state until merged and verified.

## Licensing

GoreeCloud has an active central software-licensing policy, but this documentation pass has not found a Bookmarks-specific authoritative record explicitly assigning the repository license. No `LICENSE` file is added by inference. The license must be recorded through the governed project/repository process before any public software release that requires it.

## Truthful status rule

Planned capabilities, roadmap entries, specifications, mockups, manifests, metadata, or documentation presence must never be represented as implemented, accepted, production-ready, or Stable without authoritative implementation and evidence.