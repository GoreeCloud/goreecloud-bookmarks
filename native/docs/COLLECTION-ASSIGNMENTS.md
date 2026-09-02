# GoreeCloud Bookmarks Collection Assignments

## Scope

The native Bookmarks service supports one owner-scoped collection assignment per bookmark through a dedicated relationship store and bounded HTTP surface. This is a development-stage native contract and not a production-readiness claim.

## HTTP surface

The current native service exposes:

- `GET /api/v1/bookmarks/{id}/collection` — read the authenticated owner's current assignment;
- `PUT /api/v1/bookmarks/{id}/collection` — assign or reassign the authenticated owner's bookmark to a collection using `{ "collectionId": "..." }`;
- `DELETE /api/v1/bookmarks/{id}/collection` — remove the authenticated owner's current assignment.

Bookmark deletion also removes the corresponding in-memory assignment so the development relationship store does not retain stale membership for a deleted bookmark.

## Owner isolation

Every relationship operation is bound to the resolved owner identity. Bookmark existence is checked through the repository-backed bookmark service, and collection lookup is owner-scoped.

A bookmark or collection owned by another account is intentionally indistinguishable from a missing object at this API boundary. Cross-owner requests return a generic not-found response and do not expose the other owner's bookmark, collection, assignment, title, URL, or owner identifier.

GET and DELETE re-check current bookmark ownership before consulting or removing assignment state. This prevents a stale relationship record from becoming an ownership side channel.

## Fail-closed storage behavior

The bookmark lookup contract carries repository errors separately from the owner-scoped existence result. A repository failure therefore remains an availability error instead of being silently converted to “not found.”

The HTTP layer returns a sanitized service-unavailable response for such failures and does not expose backend error details.

## Request handling

Assignment writes use the same bounded strict JSON decoding used by the native bookmark service. Unknown fields, malformed JSON, oversized bodies, and multiple JSON values are rejected rather than partially interpreted.

## Persistence limitation

Collections and collection assignments are currently in-memory development foundations. The PostgreSQL development runtime driver and durable collection/assignment persistence are not yet integrated, and the current HTTP surface does not imply production persistence or Stable acceptance.
