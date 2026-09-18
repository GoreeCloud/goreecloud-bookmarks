# GoreeCloud Bookmarks — Logical Data Model

## Status

- **Lifecycle:** Development
- **Model status:** Selected logical model with the initial PostgreSQL Bookmark subset implemented; other logical objects remain planned.
- **Server authority:** GoreeCloud Bookmarks
- **Authoritative wire schemas:** `docs/api/openapi.yaml`
- **Persistence direction:** PostgreSQL server state; SQLite installed-client offline state; archive payloads outside relational storage.

This record defines the logical objects and ownership relationships that implementation must preserve. Migration version `1` currently implements only the initial Bookmark row subset; it does not imply that every logical object below has a SQL representation.

## 1. Common synchronized fields

Unless an object is explicitly local-only or immutable, synchronized resources use these concepts:

- `id` — opaque stable server-generated identifier.
- `ownerId` — owning GoreeCloud identity/user reference where the object is user-owned.
- `revision` — positive integer increased when synchronization-relevant state changes.
- `createdAt` — server-recorded creation timestamp.
- `updatedAt` — server-recorded latest material update timestamp.
- `deletedAt` — deletion timestamp when the object is tombstoned; absent/null for active objects.

Opaque identifiers are not security boundaries and must not be interpreted by clients.

## 2. User

Bookmarks references GoreeCloud user/identity authority rather than owning credentials.

Logical responsibilities:

- Bookmarks-local user profile/preferences that belong to this product.
- Account/library ownership boundary.
- Quotas or storage accounting where implemented.
- Default privacy/archive/offline preferences where product scope requires them.

Passwords, primary identity credentials, identity-provider signing material, and unrelated GoreeCloud identity data are not Bookmarks-owned fields.

## 3. Device

Represents an installed client or synchronization participant when device identity is required.

Core fields:

- `id`
- `ownerId`
- `displayName`
- `platform`
- `clientVersion`
- `lastSeenAt`
- `revokedAt`
- synchronization cursor/checkpoint metadata as needed

Device records support session/device management and sync visibility but must not be used to bypass user authorization.

## 4. Bookmark

The Bookmark is the central saved-item record. The current Development schema implements its initial owner-scoped persisted subset in PostgreSQL.

Core fields:

- common synchronized fields;
- `url` — source URL when the item is URL-backed;
- `originalUrl` — optional provenance URL when canonicalization/tracking cleanup changes the working URL;
- `title`;
- `description`;
- `contentType`;
- `source`;
- `author`;
- `language`;
- `publishedAt`;
- `lastVisitedAt`;
- `favorite`;
- `pinned`;
- `readState` — unread / reading / completed / archived where applicable;
- `readingProgress` where applicable;
- `privacyLevel` — normal / sensitive / private-vault / shared;
- `previewImageRef` where permitted;
- `extractedTextRef` or derived-content relationship rather than unbounded text duplication when implementation requires it;
- archive/offline summary state derived from related records.

Collection membership, tags, notes, highlights, archive versions, reminders, sharing, and link-health records remain related objects rather than being flattened into one unbounded Bookmark row.

## 5. Collection

Represents structured organization.

Core fields:

- common synchronized fields;
- `name`;
- `description`;
- `parentId` for nested collections;
- `icon`;
- `coverRef`;
- `sortOrder`;
- default archive policy;
- default privacy policy;
- default offline policy;
- default view settings where synchronized.

Cycles in collection parentage are invalid.

## 6. CollectionMembership

Represents bookmark-to-collection membership when the implementation permits more than one collection relationship or needs explicit ordering/metadata.

Core fields:

- `id`;
- `bookmarkId`;
- `collectionId`;
- `sortOrder`;
- `createdAt`.

If the initial implementation limits bookmarks to one primary collection, the API may expose a single collection relationship while retaining migration room for a membership table internally.

## 7. Tag

Core fields:

- common synchronized fields;
- `name`;
- normalized name/key used for uniqueness rules;
- `color`;
- `description`;
- optional icon;
- aliases where supported.

Tag rename/merge operations must preserve bookmark relationships without silently duplicating user-visible tags.

## 8. BookmarkTag

Join relationship between Bookmark and Tag.

Core fields:

- `bookmarkId`;
- `tagId`;
- `createdAt`.

The pair is unique for active relationships.

## 9. Note

A note may belong to a bookmark or exist as a bounded standalone research note when product scope enables it.

Core fields:

- common synchronized fields;
- optional `bookmarkId`;
- optional `title`;
- content body;
- content-format/version marker;
- privacy level if it may differ from the bookmark;
- attachment/reference relationships where implemented.

User-authored note bodies require preservation-first conflict handling; ambiguous concurrent edits must not be discarded silently.

## 10. Highlight

Core fields:

- common synchronized fields;
- `bookmarkId`;
- selected text;
- optional user note;
- source locator/anchor data where available;
- optional label/color metadata;
- capture/source context needed to re-anchor after page changes.

A highlight may outlive the current live page and must retain enough source context to remain meaningful against an archive when possible.

## 11. Archive

Represents preservation policy/state for a Bookmark.

Core fields:

- common synchronized fields;
- `bookmarkId`;
- requested archive mode/policy;
- status;
- `latestVersionId` where available;
- latest successful capture time;
- failure state that does not invalidate the bookmark itself.

Archive metadata is relational. Large preserved payloads are not stored directly in ordinary PostgreSQL rows.

## 12. ArchiveVersion

Represents one immutable preservation event/version.

Core fields:

- stable ID;
- `archiveId`;
- `bookmarkId`;
- capture timestamp;
- source URL at capture;
- capture mode;
- status;
- WARC object reference for complete captures where available;
- readable-content object reference where available;
- visual-snapshot object reference where available;
- MIME/format metadata;
- byte size(s);
- cryptographic digest(s);
- storage backend/object key abstraction;
- integrity state;
- failure/warning metadata safe for user/operator display.

ArchiveVersion payloads are immutable. A recapture or conversion creates a new version/derived representation rather than rewriting historical capture bytes in place.

## 13. Reminder

Core fields:

- common synchronized fields;
- `bookmarkId`;
- `dueAt`;
- optional recurrence rule;
- status;
- `lastTriggeredAt`;
- notification preference/channel references where implemented.

Reminder delivery failure does not delete the reminder.

## 14. Share

Represents an intentional disclosure boundary.

Core fields:

- common synchronized fields;
- `resourceType` — bookmark / collection / archive;
- `resourceId`;
- `accessMode` — specific-people / shared-group / anyone-with-link / public;
- optional expiration;
- revocation state;
- publication/display configuration where applicable.

A share never implicitly publishes parent collections, unrelated tags, private notes, or hidden internal metadata.

## 15. Permission

Defines access granted through a Share or collaborative resource.

Core fields:

- stable ID;
- `shareId` or collaboration boundary;
- principal type;
- principal identity/group reference where applicable;
- role — viewer / contributor / editor / manager / owner as supported by the resource;
- created/updated/revoked timestamps.

Permissions are enforced server-side. Client UI visibility is not authorization.

## 16. SyncEvent

Represents a server-side synchronization change record or equivalent change-feed entry.

Core fields:

- monotonic/opaque cursor position;
- object type;
- object ID;
- resulting object revision;
- change type — upsert / delete;
- changed timestamp;
- minimal routing metadata needed by clients.

Sync events should avoid embedding unnecessary private content when clients can retrieve the authorized resource by ID.

## 17. ClientMutation

Represents an offline/retry-safe client write during synchronization processing.

Core fields:

- durable `clientMutationId` unique within its defined scope;
- device/client identity where applicable;
- object type and object ID when known;
- operation;
- `baseRevision` for conflict-aware updates;
- payload;
- submitted/processed timestamps;
- outcome/result revision;
- conflict/error state.

The server must make retry of the same accepted mutation safe.

## 18. AutomationRule

Core fields:

- common synchronized fields;
- `name`;
- enabled state;
- ordered conditions;
- ordered actions;
- execution-policy metadata;
- last execution/error summary where user-visible.

Rule condition/action types must be explicitly versioned/validated. Arbitrary executable code is not part of the initial rule model.

## 19. LinkHealthRecord

Core fields:

- stable ID;
- `bookmarkId`;
- state — healthy / redirected / changed / temporarily-unavailable / permanently-unavailable / blocked / unknown;
- checked timestamp;
- safe HTTP/status summary where applicable;
- redirect target where safe;
- previous/last healthy timestamp;
- warning/error details appropriate for user display.

Link-health processing must not become a privacy bypass for Sensitive or Private Vault items.

## 20. SavedSearch

Core fields:

- common synchronized fields;
- `name`;
- query text;
- structured filters;
- sort definition;
- optional icon/color/pin/order settings.

SavedSearch is a query definition, not a duplicated result set.

## 21. Job

Internal durable work such as metadata extraction, archive capture, indexing, imports, exports, link checks, notifications, or automation execution uses a job record.

Core fields:

- stable ID;
- job type;
- owner/resource references;
- status — pending / running / succeeded / failed / cancelled;
- attempt count;
- available/next-at timestamp;
- started/completed timestamps;
- idempotency/deduplication key where required;
- bounded structured input reference;
- result reference;
- safe error summary.

Job storage is operational state, not an API promise that every internal job is exposed to normal clients.

## 22. Privacy classification

Data handling must distinguish at least:

- ordinary private library data;
- Sensitive data with reduced external processing/preview behavior;
- Private Vault ciphertext/metadata whose server-visible scope is constrained by the future cryptographic design;
- intentionally Shared/public representations.

Changing privacy level must trigger validation of derived previews, indexes, archives, sharing state, and external-processing permissions rather than changing only a label.

## 23. Server storage boundaries

PostgreSQL stores relational metadata and transactional state. The current implementation includes only the initial `bookmarks` relation and migration-history metadata; other relational objects below remain planned until implemented through governed migrations.

Archive bodies, WARC files, visual snapshots, large extracted artifacts, and other large payloads use the archive/file storage abstraction outside ordinary relational rows. PostgreSQL stores authorized references and integrity metadata.

Search starts with PostgreSQL full-text search. A future separate search index is derived state and must be rebuildable from authoritative Bookmarks data/approved preserved content.

## 24. Client-local data

Installed clients use SQLite for local/offline state. Client databases may contain:

- cached authorized resource metadata;
- locally searchable data permitted by privacy settings;
- sync cursor/checkpoint state;
- pending client mutations;
- device-local preferences;
- references to locally stored offline content.

A client cache may be rebuilt from the server when safe, but **unsynchronized user mutations are not disposable cache data** and must be protected during client-schema migration, logout/revocation handling, failure recovery, and cache reset.

## 25. Deletion and recovery

Normal deletion creates synchronization-visible deletion/tombstone state before physical purge where required for offline consistency.

Physical erasure must respect applicable user intent, privacy obligations, backup/retention rules, active shares, archive relationships, and recovery boundaries.

Backup copies are governed recovery artifacts, not active synchronized duplicates. Their expiration and restore behavior must be documented separately from live-resource deletion.

## 26. Schema evolution rule

The logical model may evolve, but implementation must preserve stable IDs, ownership, privacy semantics, and recoverability. Breaking wire changes require API-version handling. Database changes use governed migrations described in `docs/migrations.md` rather than ad hoc production edits.