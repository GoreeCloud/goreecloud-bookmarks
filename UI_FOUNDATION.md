# GoreeCloud Bookmarks UI Foundation

This document records the first visual-architecture pass for GoreeCloud Bookmarks. The UI work is intentionally stacked on `feature/goreecloud-identity` so the tested identity baseline remains reviewable on its own.

## Direction

The interface should feel like a focused private library rather than an administration dashboard. The design direction combines:

- Clear information hierarchy and predictable controls.
- Refined spacing, typography, rounded geometry, and restrained depth.
- Strong light and dark theme behavior using the existing semantic theme tokens.
- A calm navigation rail inspired by modern bookmark/read-later applications without copying another product's branding or proprietary artwork.
- Accessible focus states, labels, and touch targets.

## Initial shell changes

The first pass deliberately changes only shared layout surfaces:

1. `apps/web/layouts/MainLayout.tsx`
   - Introduces a neutral application background and inset content surface on desktop.
   - Keeps the existing responsive sidebar width calculations and page behavior.
   - Adds a rounded, bordered primary content surface without changing routing or page data.

2. `apps/web/components/SidebarShell.tsx`
   - Adds GoreeCloud Bookmarks identity to the persistent navigation shell.
   - Adds a compact "Private library" descriptor when expanded.
   - Refines backdrop, separators, focus treatment, mobile overlay, and control geometry.
   - Preserves existing navigation content, profile behavior, responsive collapse behavior, and route-closing behavior.

3. `apps/web/components/DashboardItem.tsx`
   - Replaces the heavy gradient statistic tile with a quieter card treatment.
   - Uses semantic theme tokens, tabular numerals, restrained hover feedback, and a more legible label/value hierarchy.

## Non-goals for this pass

- No backend, database, API, authentication, or authorization changes.
- No route changes.
- No bookmark/collection/tag behavior changes.
- No full dashboard redesign yet.
- No attempt to copy Raindrop.io assets, trademarks, or exact layouts.
- No change to upstream licensing or attribution.

## Next visual passes

After this shell is validated, continue with:

- Dashboard header and section hierarchy.
- Bookmark card/list density and metadata hierarchy.
- Search and create flows.
- Collections and tags navigation refinement.
- Empty states and modal consistency.
- Mobile navigation and touch ergonomics.
- Light/dark visual QA and keyboard accessibility.

## Validation

This branch must continue to pass the GoreeCloud web build and inherited login regression suite before it is considered for merge into the identity branch. Broader visual and interaction testing is required before any production release.
