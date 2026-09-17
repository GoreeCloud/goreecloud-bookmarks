# GoreeCloud Bookmarks — User Manual

## Current availability

GoreeCloud Bookmarks is currently at the **Experimental** lifecycle stage.

A minimal Go service foundation now exists for engineering validation, but there is still no supported end-user bookmark workflow, web application, desktop client, mobile client, production service, account system, or persistent bookmark library.

The only currently implemented HTTP behavior is:

- `GET /api/v1/health` — bounded process-health response.
- `GET /api/v1/ready` — intentionally returns HTTP `503` and `ready: false` because the required Bookmarks data layer is not configured.

These endpoints are engineering signals, not an end-user Bookmarks interface.

## Experimental developer execution

With the pinned Go `1.27.1` toolchain, the service foundation can be validated with:

```bash
go vet ./...
go test ./...
go build -o ./build/goreecloud-bookmarks ./cmd/bookmarks
```

For local experimental execution:

```bash
go run ./cmd/bookmarks
```

It listens on `127.0.0.1:8080` by default. `GOREECLOUD_BOOKMARKS_LISTEN_ADDR` may explicitly override that local-development listener.

This is not a production installation procedure. No persistent data, authentication, authorization, migration, backup, reverse-proxy, TLS, database, or deployment configuration is established by the current service foundation.

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
