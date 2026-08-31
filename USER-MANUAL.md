# GoreeCloud Bookmarks User Manual

## Current status

GoreeCloud Bookmarks is under active native migration and is not Stable or production-approved. The repository contains a GoreeCloud-owned native development service plus a transitional Linkwarden-derived compatibility surface.

This manual describes currently supportable behavior without presenting development-only modes as production authentication, storage, migration, or recovery.

## Running the native development service

From the repository:

```bash
cd native
go test ./...
go run ./cmd/bookmarksd
```

The service binds to `127.0.0.1:8080` by default unless `GOREECLOUD_BOOKMARKS_ADDR` is set.

By default, the native service uses:

- in-memory development bookmark storage;
- unavailable production identity integration;
- no production approval.

## Local development identity

To exercise owner-scoped behavior locally, select the explicit development identity mode while retaining the in-memory development repository:

```bash
export GOREECLOUD_BOOKMARKS_STORE=memory-development
export GOREECLOUD_BOOKMARKS_IDENTITY=development-header
```

Requests then require:

```text
X-GoreeCloud-Development-Owner: <development-owner-id>
```

This header is a local development seam only. It is not a production login, token, or authorization mechanism.

## Bookmark API

Current native bookmark routes include:

- `GET /api/v1/bookmarks`
- `GET /api/v1/bookmarks/{id}`
- `POST /api/v1/bookmarks`
- `PATCH /api/v1/bookmarks/{id}`
- `DELETE /api/v1/bookmarks/{id}`

Example create body:

```json
{
  "url": "https://example.com/",
  "title": "Example",
  "note": "Reference",
  "tags": ["research"]
}
```

URLs must use HTTP or HTTPS. URL fragments are removed. Unknown JSON fields, trailing JSON values, malformed input, and oversized request bodies are rejected.

## Collections

The native development service supports owner-scoped collection organization.

Routes:

- `GET /api/v1/collections`
- `GET /api/v1/collections/{id}`
- `POST /api/v1/collections`
- `PATCH /api/v1/collections/{id}`

Create a top-level collection:

```json
{
  "name": "Research"
}
```

Create a nested collection:

```json
{
  "name": "Documentation",
  "parentId": "<parent-collection-id>"
}
```

Move a collection to another parent:

```json
{
  "parentId": "<new-parent-id>"
}
```

Use an empty `parentId` to move a collection to the root.

Collection rules:

- names are required and limited to 256 runes;
- sibling names are unique case-insensitively;
- parents must belong to the same owner;
- self-parenting and hierarchy cycles are rejected;
- another owner's collection is reported as not found rather than exposed.

## Bookmark-to-collection assignment

In the current `memory-development` service, one bookmark can be assigned to one collection at a time through these owner-scoped routes:

- `GET /api/v1/bookmarks/{id}/collection`
- `PUT /api/v1/bookmarks/{id}/collection`
- `DELETE /api/v1/bookmarks/{id}/collection`

Assign or move a bookmark to a collection:

```json
{
  "collectionId": "<collection-id>"
}
```

`PUT` replaces the bookmark's current assignment. A bookmark and target collection must belong to the same resolved owner; another owner's bookmark, collection, or assignment is represented as not found rather than disclosed. Deleting a bookmark removes its in-memory collection assignment.

Collection hierarchy and assignment state are currently in-memory development state. They are lost when the native development service restarts. Persistent relationship storage, collection deletion/cascade policy, multi-collection assignment, sharing, and migration semantics are not production-ready.

## Transitional Linkwarden surface

The inherited Linkwarden-derived application may still be required for workflows not yet migrated into native GoreeCloud Bookmarks. Do not interpret the presence of native bookmark, collection, or assignment routes as authorization to remove the inherited runtime or migrate production data.

## Sync

The repository contains first-party GoreeCloud Sync source foundations for bookmark records. Production synchronization still depends on accepted identity, durable storage, deployment, recovery, and protocol/runtime evidence.

## Privacy and security

Privacy Shield governs privacy and retained-data behavior. Wardveil Security governs GoreeCloud security/protection state. Development routes do not create independent production privacy or security approval.

## Backup and recovery

The in-memory native development service is not a backup system. Everkeep is the authoritative continuity system. Production migration requires accepted durable storage, export/import, backup, restore, rollback, and inherited-data migration evidence.

## What is not yet production-ready

- production GoreeCloud Identity authentication;
- durable native bookmark, collection, and assignment persistence;
- collection deletion/cascade semantics;
- shared collections and permissions;
- retained Linkwarden data migration;
- accepted import/export and preserved-content migration;
- production Sync operation;
- complete Glaze UI, Privacy Shield, Wardveil Security, and Everkeep acceptance;
- production deployment or Stable release.

## Troubleshooting

If bookmark or collection requests return an identity-integration error, confirm that you intentionally selected `development-header` for local testing and that the repository mode is `memory-development`.

If a collection move returns a conflict, verify that the requested parent does not create a cycle or duplicate the collection name under that parent.

If a bookmark-to-collection assignment returns not found, verify that both IDs exist under the same development owner identity.

If a collection or assignment is missing after restarting the service, that is expected in the current in-memory development mode.
