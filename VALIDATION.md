# GoreeCloud Bookmarks Validation Baseline

This document defines the minimum validation expected before the current GoreeCloud identity and Glaze UI integration candidate can be considered for merge. It does not replace the upstream test suite, project specification, or release procedure.

## Current scope

- Repository: `GoreeCloud/goreecloud-bookmarks`
- Working branch: `feature/integrate-goreecloud-ui-foundation`
- Upstream basis: Linkwarden v2.16.0
- Baseline commit: `62f1b81ff7f66001b0f5f613202f87771f3186ee`
- Current authoritative base: GoreeCloud Bookmarks `2.16.0-gc.1-rc.1` source at `2152cf2030e6fdc4ed26da7edd61206bbdf5b9c1`
- Candidate scope: GoreeCloud identity, primary responsive workspace, Glaze modal foundation, bookmark-import resilience, and confirmation safety reconciled with the current capture, portability, CI, architecture, and release foundation

## Reproducible install and build

The repository pins Yarn 4.12.0 through `packageManager` in the root `package.json`.

```bash
corepack enable
corepack prepare yarn@4.12.0 --activate
yarn --version
yarn install --immutable
yarn prisma:generate
yarn web:build
```

A successful build must complete without changing dependency lockfiles or requiring reusable credentials to be committed.

## Automated tests

Run the repository test suite after dependencies and Prisma artifacts are prepared:

```bash
yarn test --run
```

The first-party GitHub Actions workflow must run the complete test suite against PostgreSQL 16, apply all Prisma migrations, and complete the production web build on the exact candidate head. The dedicated Playwright workflow must independently validate its browser suite on that same head.

## Identity validation

Before merge, verify the built application rather than source text alone.

1. Browser title reads **GoreeCloud Bookmarks**.
2. Login and authentication screens use GoreeCloud Bookmarks identity.
3. The GoreeCloud Bookmarks favicon/application mark appears where supported.
4. PWA install metadata identifies GoreeCloud Bookmarks and does not show upstream Linkwarden screenshots.
5. The profile menu does not promote Linkwarden documentation, support, or release pages as the GoreeCloud product destination.
6. The application does not poll Linkwarden's announcement feed.
7. No ordinary user-facing screen presents Linkwarden as the primary product brand, except intentional interoperability or attribution references.
8. Required upstream license and provenance records remain intact.

## Functional regression checks

At minimum, validate:

- Sign in and sign out.
- Registration when enabled.
- Password reset/email verification when configured.
- Create, edit, delete, and search bookmarks.
- Collections and tags.
- Pinned links.
- Import/export paths relevant to the deployment.
- Preservation/archive views used by GoreeCloud.
- Shared collection behavior and user separation.
- Light, dark, and system theme behavior.
- Desktop and mobile navigation.
- Keyboard focus and understandable accessibility labels on modified surfaces.

## Release gate

Do not merge this branch until both required GitHub workflows pass on the exact candidate head and the applicable code, privacy, security, accessibility, responsive, light/dark, and interaction reviews are complete. Do not deploy it as a GoreeCloud release until the separate release and deployment gates are authorized and recorded. A green draft pull request is not release, deployment, production acceptance, or Stable approval.
