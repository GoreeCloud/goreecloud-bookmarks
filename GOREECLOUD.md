# GoreeCloud Bookmarks Fork Record

## Product identity

- **Product:** GoreeCloud Bookmarks
- **Repository:** `GoreeCloud/goreecloud-bookmarks`
- **Development model:** Maintained open-source fork of Linkwarden
- **Status:** Development and validation only; production deployment is not yet approved

## Upstream provenance

- **Upstream project:** Linkwarden
- **Upstream repository:** `linkwarden/linkwarden`
- **Initial GoreeCloud fork base:** Linkwarden v2.16.0
- **Initial upstream commit:** `62f1b81ff7f66001b0f5f613202f87771f3186ee`
- **License:** GNU Affero General Public License v3.0 (AGPL-3.0), preserved from upstream

The GoreeCloud fork must retain required upstream copyright, attribution, licensing, and source-availability obligations. Rebranding the user-facing product does not remove or replace those obligations.

## GoreeCloud project scope

GoreeCloud Bookmarks is intended to become a private, self-hosted, multi-user bookmark, read-later, web-collection, and clipping application. GoreeCloud-specific work will focus on:

- Original GoreeCloud product identity and interface design.
- A visually polished bookmark workspace inspired by the usability of Raindrop.io without copying proprietary branding or protected assets.
- Private personal bookmark libraries by default.
- Explicit shared collections and controlled multi-user permissions.
- Fast Firefox-based capture through the companion GoreeCloud browser extension.
- Search, organization, annotations, preservation, import/export, portability, backup, and recovery.
- Minimal unnecessary divergence from Linkwarden so compatible upstream security fixes and improvements can continue to be reviewed and integrated.

GoreeCloud Bookmarks is separate from GoreeCloud Research Library. Bookmarks manages ordinary saved web content and read-later material; Research Library retains research-specific verification, evidence-preservation, and source-to-decision workflows.

## Companion browser extension

The canonical Firefox extension is maintained at:

`GoreeCloud/goreecloud-firefox-extensions`, under `extensions/bookmarks/`.

The legacy `GoreeCloud/goreecloud-bookmark-browser-extension` repository remains
an upstream-provenance and behavioral reference during the controlled first-party
Firefox replacement. The initial extension ancestry is
`linkwarden/browser-extension` v1.5.4.

GoreeCloud Browser integrates separately through the versioned
`POST /api/v1/bookmarks/capture` contract. GoreeCloud Bookmarks remains the
authoritative synchronized bookmark service; signing into Browser does not
authorize silent upload of local browser bookmarks.

## Branch and maintenance policy

- `main` represents accepted GoreeCloud repository state and should not receive experimental identity work directly.
- `feature/*` branches are used for isolated GoreeCloud product changes.
- `fix/*` branches are used for compatible corrections.
- `security/*` branches are used for security-sensitive work.
- `upstream-sync/*` branches isolate reviewed Linkwarden synchronization.
- `release/*` branches prepare identifiable release candidates without implying production approval.
- Upstream synchronization should be reviewed separately before integration.
- Security-relevant upstream fixes receive elevated review priority.
- GoreeCloud-specific patches should remain documented and isolated where practical.
- The exact upstream ancestry of each GoreeCloud release must remain traceable.

## Current implementation boundary

The repository has moved beyond provenance-only setup while retaining a strict
release boundary:

- `main` contains the published `2.16.0-gc.1-rc.1` source and immutable release
  evidence. The candidate is a prerelease and `production_approved` remains
  false.
- Authenticated Firefox capture, the native GoreeCloud Browser capture contract,
  collection-authorization regression coverage, hardened portable import/export,
  and first-party application CI are integrated on `main`.
- The GoreeCloud identity, application shell, primary bookmark workspace, and
  shared Glaze modal foundation are being reconciled onto the current `main`
  foundation through `feature/integrate-goreecloud-ui-foundation`.
- Secondary localization, generated email, account, administration, and other
  inherited Linkwarden identity surfaces still require a deliberate product and
  attribution review before production acceptance.
- No production deployment, DNS, Caddy, NetBird, database cutover, or user-data
  migration is approved by this record or by the release candidate.
