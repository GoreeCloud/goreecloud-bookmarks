# Migration Record

Legacy repository: `GoreeCloud/goreecloud-bookmark-browser-extension`

Inspected legacy baseline: main tree `3175f85191bab62d84873c46e54b32ab958f9be8`.

The legacy project is a Linkwarden-derived cross-browser extension. Its Firefox manifest identifies Linkwarden 1.5.4, uses upstream add-on ID `jordanlinkwarden@gmail.com`, requests `storage`, `scripting`, `activeTab`, `tabs`, `bookmarks`, and `contextMenus`, and grants `<all_urls>` host access. The repository also contains a Safari manifest and Xcode project material.

The owning GoreeCloud Bookmarks application repository does not import Safari/Xcode release artifacts. Instead, GoreeCloud begins a Firefox-specific first-party replacement with a narrower permission model and GoreeCloud-controlled add-on identity. The legacy project remains provenance and behavioral reference until equivalent required workflows are validated here.

No Stable release is claimed by this migration. Mozilla signing, persistent installation, runtime capture validation, server endpoint acceptance, and restart verification remain separate gates.

## Application ownership migration completion

The Firefox client source-authority migration completed on 2026-09-18.

- Owning application migration: Bookmarks PR #8 merged to authoritative `main` at `243e9979bc97495ec1d57f4240131d31bf2b3014`.
- Owning application post-merge validation: `Bookmarks Firefox Client` run #2 and `Validate Go` run #14 both passed.
- Shared Firefox cleanup: `GoreeCloud/goreecloud-firefox-extensions` PR #103 removed the transitional `extensions/bookmarks/` copy and reconciled shared policy/inventory.
- Shared repository post-merge validation: Firefox Repository run #396 passed.
- Current canonical source: `GoreeCloud/goreecloud-bookmarks/clients/firefox/`.
- Firefox add-on ID remains `goreecloud-bookmarks@goreecloud.com`.
- Current source version remains `0.1.1`; no Stable release is claimed by the repository move.

Historical source provenance remains available through Git history in the shared Firefox repository and the migration metadata in this directory. The shared repository is no longer an active source authority for GoreeCloud Bookmarks Firefox development.
