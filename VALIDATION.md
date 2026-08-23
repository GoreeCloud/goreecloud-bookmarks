# GoreeCloud Bookmarks Validation Baseline

This document defines the minimum validation expected before the current GoreeCloud identity branch can be considered for merge. It does not replace the upstream test suite or release process.

## Current scope

- Repository: `GoreeCloud/goreecloud-bookmarks`
- Working branch: `feature/goreecloud-identity`
- Upstream basis: Linkwarden v2.16.0
- Baseline commit: `62f1b81ff7f66001b0f5f613202f87771f3186ee`

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

The existing GitHub Actions Playwright workflow is also expected to validate the supported login path when repository Actions are enabled and a pull-request run is available.

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

Do not merge this branch or deploy it as a GoreeCloud release until the required build, automated tests, and interactive regression checks are recorded as successful. A draft pull request is not release approval.
