# GoreeCloud Bookmarks UI Foundation

This document records the visual-architecture foundation for GoreeCloud Bookmarks. The UI work is intentionally stacked on `feature/goreecloud-identity` so the tested identity baseline remains reviewable on its own.

## Direction

The interface should feel like a focused private library rather than an administration dashboard. The design direction combines:

- Clear information hierarchy and predictable controls.
- Refined spacing, typography, rounded geometry, and restrained depth.
- Strong light and dark theme behavior using the existing semantic theme tokens.
- A calm navigation rail inspired by modern bookmark/read-later applications without copying another product's branding or proprietary artwork.
- Accessible focus states, labels, and touch targets.
- Consistent information hierarchy across card, masonry, list, dashboard, search, and capture surfaces.

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

- Collections and tags navigation refinement.
- Empty-state and modal consistency beyond the dashboard/search/capture surfaces.
- Mobile navigation and touch ergonomics.
- Light/dark visual QA and keyboard accessibility.
- Firefox extension visual alignment after the web application shell stabilizes.

## Validation

Every code-bearing head of this branch must pass the GoreeCloud web build and inherited login regression suite before it is considered for merge into the identity branch. Broader responsive, visual, keyboard, and interaction testing is required before any production release.
