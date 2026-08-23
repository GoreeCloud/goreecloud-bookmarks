# GoreeCloud Bookmarks 2.16.0-gc.1-rc.1

Lifecycle: Release Candidate

This is the first GoreeCloud Bookmarks application release candidate based on the Linkwarden 2.16.0 lineage. It is intended for controlled acceptance and must not be represented as Stable.

## Included

- GoreeCloud Bookmarks maintained-fork application foundation and GoreeCloud product experience.
- Hardened portable import/export contracts and bounded migration behavior.
- Authorization regression coverage for bookmark capture and collection access.
- Authenticated Firefox extension capture compatibility endpoint.
- Native GoreeCloud Browser `POST /api/v1/bookmarks/capture` contract while retaining the Firefox compatibility route.
- First-party GoreeCloud Bookmarks CI and Playwright validation gates.

## Release artifact

The release workflow publishes a multi-platform OCI image for `linux/amd64` and `linux/arm64` to GHCR using the immutable RC version tag. The workflow records the exact source revision and image digest and requests BuildKit SBOM and provenance attestations.

## Known acceptance gates

This candidate is not Stable and does not authorize production deployment by itself. Remaining acceptance includes production-representative runtime validation, multi-user authorization and sharing checks, migration and rollback validation, backup and restoration proof, Glaze UI and accessibility validation, and release-critical Firefox integration validation. Mozilla signing and persistent-installation acceptance for the canonical Firefox extension remain a separate release boundary.

`production_approved` remains false.
