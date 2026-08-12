# GoreeCloud Bookmarks UI Foundation

This document records the visual-architecture foundation for GoreeCloud Bookmarks. The UI work is intentionally stacked on `feature/goreecloud-identity` so the tested identity baseline remains reviewable on its own.

## Direction

The interface should feel like a focused private library rather than an administration dashboard. The design direction combines:

- Clear information hierarchy and predictable controls.
- Refined spacing, typography, rounded geometry, and restrained depth.
- Strong light and dark theme behavior using the existing semantic theme tokens.
- A calm navigation rail inspired by modern bookmark/read-later applications without copying another product's branding or proprietary artwork.
- Accessible focus states, labels, and touch targets.
- Consistent information hierarchy across card, masonry, list, dashboard, search, capture, collection, and tag surfaces.

## Completed foundation work

### Application shell

1. `apps/web/layouts/MainLayout.tsx`
   - Introduces a neutral application background and inset content surface on desktop.
   - Keeps the existing responsive sidebar width calculations and page behavior.
   - Adds a rounded, bordered primary content surface without changing routing or page data.

2. `apps/web/components/SidebarShell.tsx`
   - Adds GoreeCloud Bookmarks identity to the persistent navigation shell.
   - Adds a compact "Private library" descriptor when expanded.
   - Refines backdrop, separators, focus treatment, mobile overlay, and control geometry.
   - Preserves existing navigation content, profile behavior, responsive collapse behavior, and route-closing behavior.

### Dashboard hierarchy

3. `apps/web/components/DashboardItem.tsx`
   - Replaces the heavy gradient statistic tile with a quieter card treatment.
   - Uses semantic theme tokens, tabular numerals, restrained hover feedback, and a more legible label/value hierarchy.

4. `apps/web/pages/dashboard.tsx`
   - Establishes a stronger page title and description hierarchy.
   - Groups dashboard view/layout controls into a compact control surface.
   - Gives dashboard sections consistent inset surfaces, section headers, descriptions, and accessible "View all" actions.
   - Replaces heavy empty states and loading placeholders with calmer surfaces that match the final information architecture.
   - Preserves dashboard section order, saved layout behavior, routing, and drag/drop behavior.

### Bookmark views

5. `apps/web/components/LinkViews/LinkComponents/LinkCard.tsx`
   - Makes bookmark title the primary visual element and source/domain secondary.
   - Refines preview presentation, selected state, preserved-format indicators, footer metadata, hover/focus actions, and optional description rendering.
   - Preserves open, edit-selection, pin, actions, and drag behavior.

6. `apps/web/components/LinkViews/LinkComponents/LinkMasonry.tsx`
   - Applies the same bookmark hierarchy to variable-height cards.
   - Refines descriptions and tag chips while preserving tag navigation and existing display preferences.

7. `apps/web/components/LinkViews/LinkComponents/LinkList.tsx`
   - Reworks the dense row view into a compact library-style result with clear icon, title, source, collection, date, preservation, and action hierarchy.
   - Preserves selection, open, pin, and drag behavior.

8. `apps/web/components/LinkViews/LinkComponents/LinkTypeBadge.tsx`
   - Standardizes domain/type display into a restrained metadata badge.
   - Preserves direct external-link behavior for URL sources.

9. `apps/web/components/LinkViews/LinkComponents/LinkPin.tsx`
   - Refines pin affordance and adds explicit accessible labels while preserving pin behavior.

10. `apps/web/components/DashboardLinks.tsx`
    - Aligns dashboard bookmark cards with the canonical card-view hierarchy and actions.

### Search and quick capture

11. `apps/web/components/SearchBar.tsx`
    - Rebuilds the inline search field as a focused library-search control with clearer focus treatment, a compact clear action, and consistent advanced-operator suggestions.
    - Preserves public-collection search routing, private `/search` routing, percent-symbol validation, and existing advanced-search syntax.
    - Removes the product-facing dependency on the upstream Linkwarden advanced-search documentation link from the suggestion surface.

12. `apps/web/components/ModalContent/SearchModal.tsx`
    - Reworks the `Ctrl/⌘+K` experience into a command-style search surface with a prominent input, recent searches, advanced operators, keyboard cues, and stronger focus behavior.
    - Preserves browser-local recent-search history and existing search-query routing.
    - Keeps advanced operators as additive helpers rather than changing search semantics.

13. `apps/web/components/ModalContent/NewLinkModal.tsx`
    - Rebuilds quick capture around URL and collection as the primary decisions, with title, tags, and description moved into an expandable secondary layer.
    - Adds a lightweight client-side source/domain preview without introducing a new metadata API or network request.
    - Converts the capture surface to a semantic form so Enter can submit from the primary URL field while preserving schema validation and the existing add-link mutation.
    - Preserves collection-context defaults when capture begins from a collection route.

14. `apps/web/pages/search.tsx`
    - Turns search results into a dedicated library workspace with a prominent editable query field, integrated view/sort controls, and a calm no-results state.
    - Normalizes query decoding before passing the value into the existing links query.
    - Preserves result sorting, card/masonry/list switching, bulk selection, pagination, and link rendering behavior.

### Collections and tags

15. `apps/web/components/CollectionCard.tsx`
    - Rebuilds collection cards as library objects with a restrained color accent, collection icon, title/description hierarchy, link count, creation date, public-state cue, and compact collaborator stack.
    - Preserves collection routing, edit/share/delete-or-leave actions, owner/member lookup, permission handling, configured icons, and configured colors.

16. `apps/web/pages/collections/index.tsx`
    - Reworks the collections index into separate owned and shared collection sections with compact sorting/creation controls and a focused first-collection empty state.
    - Preserves existing collection sorting and root-collection filtering.

17. `apps/web/pages/collections/[id].tsx`
    - Replaces the legacy page-wide gradient with a compact collection identity surface and a restrained collection-color accent.
    - Makes subcollections a distinct nested-library section and keeps bookmarks as the primary content workspace.
    - Preserves open-all, edit, sharing, subcollection creation, delete/leave, membership display, collection permissions, bookmark sorting, bookmark views, and bulk-edit capability.

18. `apps/web/components/CollectionListing.tsx`
    - Refines the nested sidebar tree with clearer active rows, disclosure chevrons, collection icons/colors, public-state cues, recursive link-count badges, and keyboard focus treatment.
    - Preserves Atlaskit tree behavior, nesting, drag/reorder behavior, collection-order persistence, recursive count calculation, ownership restrictions, and collection update mutations.

19. `apps/web/components/ModalContent/NewCollectionModal.tsx`
    - Reorganizes collection creation around collection identity, name, optional description, and existing icon/color selection.
    - Adds a clearer parent-collection context for subcollection creation and prevents blank-name submission.

20. `apps/web/components/ModalContent/EditCollectionModal.tsx`
    - Aligns collection editing with the creation surface while preserving the existing update mutation, icon/color controls, name, and description.

21. `apps/web/components/TagCard.tsx`
    - Rebuilds tag cards around a consistent tag identity, creation date, and prominent linked-bookmark count.
    - Preserves tag routing, selection mode, and deletion behavior.

22. `apps/web/pages/tags/index.tsx`
    - Reworks the tags index with compact create/sort/select controls, a calmer responsive grid, a clearer multi-select action bar, and a focused first-tag empty state.
    - Preserves pagination, sorting modes, bulk deletion, and tag merging.

23. `apps/web/pages/tags/[id].tsx`
    - Gives tag detail pages the same compact identity/workspace hierarchy used by collections and search.
    - Refines inline rename and empty states while preserving rename/delete mutations, bookmark sorting, view switching, bulk selection, and tag-filtered link queries.

24. `apps/web/components/ModalContent/NewTagModal.tsx`
    - Converts tag creation into a focused keyboard-submit form with explicit cancel/create actions while preserving the existing tag upsert mutation.

### Mobile, touch, theme, and keyboard QA

25. Shared controls and focus treatment
    - Adds a consistent visible `focus-visible` ring to the shared button primitive and text-input primitive instead of relying on individual screens to restore focus visibility.
    - Keeps the existing semantic theme tokens so the same focus, border, surface, and text hierarchy works in both light and dark themes.

26. Sidebar and nested navigation
    - Widens the compact mobile rail from 56px to 64px so its primary actions can use 44px touch targets without clipping.
    - Separates collection/tag disclosure controls from their navigation links so an interactive button is no longer nested inside an interactive link.
    - Adds Escape-to-close behavior for the expanded mobile sidebar and larger mobile create/search/sidebar controls.
    - Enlarges nested collection rows and disclosure controls while preserving Atlaskit drag, reorder, nesting, ownership, and persisted ordering behavior.

27. Bookmark, collection, tag, and search touch controls
    - Makes bookmark action and pin controls visible by default on small viewports while retaining hover/focus reveal on larger screens.
    - Increases small-screen action targets for collection and tag cards without changing their mutations or routing behavior.
    - Makes tag cards keyboard reachable in both navigation and multi-select modes.
    - Enlarges mobile search clear, recent-search delete, and operator controls while preserving search semantics and recent-search storage.

## Non-goals for this foundation

- No backend, database, API, authentication, or authorization changes.
- No route changes.
- No bookmark, collection, tag, preservation, search-query, or sharing data-model changes.
- No new metadata-fetching service for the capture modal.
- No attempt to copy Raindrop.io assets, trademarks, proprietary artwork, or exact layouts.
- No change to upstream licensing or attribution.
- No production deployment approval.

## Next visual passes

Continue with:

- Empty-state and modal consistency beyond the completed dashboard/search/capture/collection/tag surfaces.
- Broader manual responsive and light/dark visual acceptance on representative phone, tablet, laptop, and desktop widths.
- Firefox extension visual alignment after the web application shell stabilizes.

## Validation

Every code-bearing head of this branch must pass the GoreeCloud web build and inherited login regression suite before it is considered for merge into the identity branch. The pre-QA collections/tags head `6d939e5815a27ec23d00f94ae9a47b4152f8873d` passed both required GitHub Actions workflows. Broader responsive, visual, keyboard, touch, and interaction testing remains required before any production release.
