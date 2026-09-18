# GoreeCloud Bookmarks — API Contract

## Status

- **Lifecycle:** Development
- **API major:** `v1`
- **Contract status:** Planned v1 architecture contract with health/readiness implemented; bookmark-domain operations remain unimplemented.
- **Machine-readable contract:** `docs/api/openapi.yaml`
- **OpenAPI description version:** `3.2.1`
- **Authoritative application data owner:** GoreeCloud Bookmarks

This document defines the initial GoreeCloud Bookmarks API boundary. Health/readiness are implemented engineering endpoints; other resource operations remain planned. Contract presence is not deployment, acceptance, or production-readiness evidence.

## 1. Authority model

GoreeCloud Bookmarks remains authoritative for bookmark data. Browser, web, desktop, mobile, integrations, and GoreeCloud Sync consume documented interfaces rather than owning or directly modifying the Bookmarks server database.

The API is the normal application boundary for synchronized account state. Direct cross-application database access is not part of the architecture.

## 2. Protocol and versioning

The initial API uses:

- HTTPS transport.
- REST-style resources.
- JSON request and response bodies unless an endpoint explicitly transfers another media type.
- Versioned paths beginning with `/api/v1/`.
- OpenAPI as the machine-readable interface description.

The API major version is independent from the GoreeCloud Bookmarks application version.

Backward-compatible additions may remain within `/api/v1/`. A change that requires existing conforming clients to change behavior or payload structure must use a new API major rather than silently changing the meaning of v1.

## 3. Authentication and authorization

Authenticated application operations require an access token accepted by GoreeCloud Bookmarks. GoreeCloud Identity/OIDC/OAuth integration is the intended platform direction where implemented, but no identity-provider URL or production token endpoint is invented in this Development-stage contract.

Authentication never replaces application authorization. Every resource operation must enforce owner, share, role, privacy-level, and administrative boundaries inside Bookmarks.

Health/readiness operations may be unauthenticated only when the deployed publication model explicitly permits them and their responses contain no sensitive application data.

## 4. Resource identifiers and revisions

Resource IDs are opaque stable server-generated strings. Clients must not derive meaning, ownership, creation time, shard, database key, or security properties from their syntax.

Mutable synchronized resources carry an integer `revision` that increases when synchronization-relevant state changes.

Where a mutation can overwrite existing state, the API uses an entity tag / revision precondition model. Clients send `If-Match` using the current entity tag or equivalent revision representation. A stale mutation returns a conflict/precondition response instead of silently replacing newer user data.

## 5. Idempotency

Operations that may be retried after an uncertain network result and could create duplicate state use an `Idempotency-Key` request header.

The initial contract requires idempotency support for at least:

- bookmark creation;
- archive capture requests;
- synchronization mutation batches;
- future import/export job creation;
- other bulk or create operations where duplicate execution could be harmful.

An idempotency key is scoped to the authenticated principal plus operation boundary and must not be treated as an authentication secret.

## 6. Pagination

Collection endpoints use cursor pagination rather than offset pagination for synchronized or frequently changing datasets.

List responses include:

- `items`;
- `nextCursor` when another page exists;
- `hasMore`.

Clients must treat cursors as opaque.

## 7. Error envelope

JSON errors use a stable envelope:

```json
{
  "error": {
    "code": "bookmark_not_found",
    "message": "Bookmark not found.",
    "requestId": "opaque-request-id",
    "details": {}
  }
}
```

`code` is the machine-readable contract. `message` is human-readable and may change without becoming an API breaking change. `details` is optional structured context and must not disclose secrets or another user's private state.

Expected HTTP semantics include:

- `400` malformed or invalid request;
- `401` authentication required/invalid;
- `403` authenticated but not authorized;
- `404` resource unavailable to the caller, including cases where revealing existence would violate privacy;
- `409` logical conflict that is not expressed as an HTTP precondition;
- `412` failed revision/ETag precondition;
- `422` structurally valid request that violates resource validation rules;
- `429` rate limit or bounded-work protection;
- `5xx` server-side failure.

## 8. Core v1 resource surface

The initial contract defines these resource families:

- `/api/v1/bookmarks`
- `/api/v1/collections`
- `/api/v1/tags`
- `/api/v1/notes`
- `/api/v1/highlights`
- `/api/v1/archives`
- `/api/v1/archive-versions`
- `/api/v1/reminders`
- `/api/v1/shares`
- `/api/v1/permissions`
- `/api/v1/automation-rules`
- `/api/v1/link-health`
- `/api/v1/saved-searches`
- `/api/v1/sync/changes`
- `/api/v1/sync/mutations`

The OpenAPI file is the authoritative wire-schema source for fields represented there. `docs/data-model.md` describes logical relationships and persistence boundaries.

## 9. Synchronization contract

Synchronization uses the same versioned HTTPS API rather than direct database access or a separate mandatory persistent socket protocol.

### Change feed

`GET /api/v1/sync/changes` returns ordered changes after an opaque cursor. Each change identifies:

- object type;
- object ID;
- object revision;
- change type (`upsert` or `delete`);
- change timestamp;
- next synchronization cursor.

### Mutation batches

`POST /api/v1/sync/mutations` accepts a bounded batch of client mutations. Each mutation includes:

- a durable client mutation ID;
- object type;
- object ID when known;
- operation;
- base revision when modifying synchronized state;
- payload appropriate to the operation.

The server returns an outcome for every mutation rather than silently dropping part of a batch.

Conflict handling prioritizes preservation. Independent deterministic field changes may merge. Ambiguous conflicting user-authored content must not be discarded solely by last-write-wins.

## 10. Deletion and tombstones

A client-facing `DELETE` operation means the resource becomes deleted from normal product views; it does not imply immediate physical erasure from every storage layer.

Synchronized deletion produces a tombstone/change event so offline clients can observe the deletion. Physical purge follows the applicable retention, privacy, recovery, and device-resume rules.

Permanent deletion endpoints or administrative purge behavior, when implemented, must be separately documented because they have different recovery consequences.

## 11. Files, archives, and large payloads

Metadata APIs do not embed arbitrarily large archive payloads in normal JSON objects.

Archive and file transfer endpoints must, when implemented:

- authorize the caller against the owning resource;
- validate content type and size;
- use safe filenames and paths;
- provide integrity metadata;
- support resumability where operationally required;
- avoid reflecting private backend storage paths;
- prevent automatic execution of active or hostile archived content.

Complete web-preservation captures use WARC 1.1 / ISO 28500:2017 as selected in `ARCHITECTURE.md`. The API exposes archive metadata and authorized retrieval, not the storage backend itself.

## 12. Privacy-sensitive responses

Responses must minimize disclosure. In particular:

- a permission failure may intentionally appear as `404` when revealing existence would leak private state;
- logs and errors must not include access tokens, passwords, private keys, raw authorization headers, full private content, or unnecessary sensitive URLs;
- external metadata-fetch details must not expose another user's library;
- Private Vault payloads must remain opaque to server operations that are not authorized by the eventual cryptographic design.

## 13. Events and webhooks

The v1 core contract does not require public webhooks or an external event bus.

Internal event boundaries may exist for Bookmarks modules and GoreeCloud Mesh integration later. If public or cross-application events are added, they must use separately documented schemas, authorization, replay/idempotency behavior, delivery guarantees, and privacy filtering.

## 14. Compatibility rules

Within API v1:

- adding an optional response field is compatible;
- adding an optional request field is compatible;
- adding a new resource or operation is compatible;
- broadening an enum must be treated carefully because clients must ignore or safely handle unknown values where the contract says extensible;
- removing or renaming a field is breaking;
- changing a field type or requiredness incompatibly is breaking;
- changing authorization or privacy semantics in a way that invalidates conforming clients requires explicit version/migration handling;
- reinterpreting an existing field with materially different semantics is breaking even if its JSON shape stays the same.

A future `/api/v2/` may coexist with v1. Retirement of v1 requires an explicit deprecation record and a verified upgrade path for supported first-party clients; no fixed retirement date is invented before GoreeCloud Bookmarks has an established release cadence.

## 15. Contract validation

When implementation begins, CI must validate the OpenAPI description and detect unreviewed breaking changes. API tests must cover authentication, authorization, validation, pagination, idempotency, precondition/conflict behavior, deletion/tombstones, synchronization, and privacy-sensitive error behavior.

No such CI or API implementation is currently verified.