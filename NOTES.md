# GoreeCloud Bookmarks — Notes

## Current verified repository state

- Product lifecycle: **Concept**.
- Repository: `GoreeCloud/goreecloud-bookmarks`.
- Product implementation: no application or service implementation currently verified.
- Canonical product specification: `GoreeCloud/Projects/Project Specification — Bookmarks.md`.
- Repository product specification: `SPECIFICATIONS.md`.
- Related capture repository: `GoreeCloud/goreecloud-bookmark-browser-extension`.
- Current Platform Contract baseline reviewed during this documentation pass: `0.4`.
- Current Integral Platform System model: exactly nine systems.
- Current Stable Glaze UI consumer target referenced by Platform Contract `0.4`: `1.5.1`.

## Important authority boundaries

- Planned feature documentation is not implementation evidence.
- `CAPABILITIES.md` records current verified capability state.
- `FEATURE-ROADMAP.md` orders planned work but does not establish completion.
- GoreeCloud Sync is separately governed and is not an Integral Platform System.
- GoreeCloud Browser integration must be reconciled with the Browser project and the separate bookmark-browser-extension repository rather than creating competing capture or bookmark authority.
- Optional intelligence must not become a dependency for core bookmarking.

## Open architecture decisions

The following remain intentionally undecided until authoritative design work selects and records them:

- Backend language/framework.
- Web client technology.
- Desktop client technology.
- Mobile client technology.
- Primary database/storage model.
- Search/index engine.
- Archive format and storage strategy.
- Background job system.
- Synchronization algorithm and transport.
- Deployment/container model.
- API versioning and compatibility policy.
- Concrete licensing record for this repository if the current central licensing policy has not yet been explicitly applied to Bookmarks.

## Documentation maintenance

When implementation begins, update `CAPABILITIES.md`, `USER-MANUAL.md`, `PRIVACY POLICY.md`, `SECURITY.md`, `goreecloud.platform.yaml`, and the feature roadmap only from verified implementation and evidence. Do not infer current capability from specifications, mockups, issue labels, or intended architecture.