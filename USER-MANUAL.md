# GoreeCloud Bookmarks — User Manual

## Current availability

GoreeCloud Bookmarks is currently at the **Development** lifecycle stage.

A Go/PostgreSQL service and persistence foundation now exists for engineering validation, but there is still no supported end-user bookmark workflow, web application, desktop client, mobile client, production service, account/authentication system, or supported user bookmark library.

The only currently implemented HTTP behavior is:

- `GET /api/v1/health` — bounded process-health response.
- `GET /api/v1/ready` — reports fail-closed database/schema readiness and returns HTTP `200` only when configured PostgreSQL is reachable with exact-current migration history.

These endpoints are engineering signals, not an end-user Bookmarks interface.

## Development engineer execution

With the pinned Go `1.27.1` toolchain, the service foundation can be validated with:

```bash
go vet ./...
go test ./...
go build -o ./build/goreecloud-bookmarks ./cmd/bookmarks
go build -o ./build/goreecloud-bookmarks-migrate ./cmd/bookmarks-migrate
```

For local Development execution without a database:

```bash
go run ./cmd/bookmarks
```

It listens on `127.0.0.1:8080` by default. `GOREECLOUD_BOOKMARKS_LISTEN_ADDR` may explicitly override that local-development listener.

For database-backed engineering validation, supply a protected `GOREECLOUD_BOOKMARKS_DATABASE_URL` and run `go run ./cmd/bookmarks-migrate` explicitly before expecting readiness to pass. The repository does not auto-load `.env` files or auto-run migrations on normal service startup.

This is not a production installation procedure. Authentication, authorization, backup/restore qualification, reverse-proxy/TLS configuration, supported production database deployment, and platform acceptance are not established by the current Development source.

## Intended core workflow

The planned end-user workflow remains:

1. Save a page, link, selection, image, file, or note.
2. Confirm capture immediately.
3. Allow optional enrichment such as metadata extraction, readable-content extraction, archival, search indexing, and organization suggestions to continue independently.
4. Organize later using Inbox, collections, tags, Smart Collections, favorites, Read Later, or search.
5. Retrieve the saved item through search, browsing, related-item discovery, reminders, or rediscovery features.

None of those end-user steps is currently implemented.

## Planned primary areas

- **Home** — recent, pinned, reading, reminder, and rediscovery surfaces.
- **Inbox** — unprocessed captures.
- **Library** — complete saved-item collection.
- **Read Later** — reading queue and progress.
- **Highlights** — cross-library highlight search and review.
- **Collections** — structured organization.
- **Tags** — flexible categorization.
- **Search** — metadata and full-content retrieval.
- **Shared** — explicitly shared items and collections.
- **Archive** — preserved copies and versions.
- **Library Health** — dead links, redirects, missing archives, duplicates, and metadata issues.
- **Import / Export** — migration and portability tools.
- **Settings** — preferences, privacy, offline, archival, automation, synchronization, and processing controls.

## Planned capture, privacy, offline, and portability behavior

A future normal save should persist the bookmark before optional enrichment begins. Planned GoreeCloud Browser workflows include saving the current page, selected text, links, images, multiple tabs, and browsing sessions.

Items are intended to be private by default, with explicit sharing and higher-privacy processing controls. No complete Bookmarks privacy implementation should be inferred from those requirements; see `PRIVACY POLICY.md` and `CAPABILITIES.md`.

Installed clients are intended to maintain useful local state and safely synchronize offline mutations. No client or synchronization implementation currently exists.

Import and export are planned to support migration preview, duplicate handling, portable browser-compatible forms, structured data, selected scopes, and complete-account export where applicable. No import/export runtime is currently implemented.

## Support and troubleshooting

The experimental service foundation may be validated as source code, but there is no supported user library or production Bookmarks service to troubleshoot yet.

A health response of `200` does not mean Bookmarks is ready for use. Current readiness is intentionally `503` until the required data layer and its validation exist.

This manual must continue to be expanded from verified implementation evidence as installation, sign-in, capture, organization, offline, synchronization, sharing, backup/recovery, update, troubleshooting, accessibility, and retirement behavior becomes real.
