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

The companion extension is maintained separately at:

`GoreeCloud/goreecloud-bookmark-browser-extension`

Its initial upstream base is `linkwarden/browser-extension` v1.5.4.

## Branch and maintenance policy

- `main` represents accepted GoreeCloud repository state and should not receive experimental identity work directly.
- `feature/*` branches are used for isolated GoreeCloud product changes.
- Upstream synchronization should be reviewed separately before integration.
- Security-relevant upstream fixes receive elevated review priority.
- GoreeCloud-specific patches should remain documented and isolated where practical.
- The exact upstream ancestry of each GoreeCloud release must remain traceable.

## Current implementation boundary

This record establishes provenance and project intent only. At this stage:

- No GoreeCloud production deployment is approved.
- No claim is made that build or automated test validation has passed on the fork.
- No authentication, authorization, storage, networking, or deployment behavior has been changed by this file.
- User-facing rebranding will be introduced incrementally on dedicated feature branches and validated before merge.
