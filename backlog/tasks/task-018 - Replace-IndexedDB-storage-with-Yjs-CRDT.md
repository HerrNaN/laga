---
id: TASK-018
title: Replace IndexedDB storage with Yjs CRDT
status: To Do
assignee: []
created_date: '2026-06-11 19:15'
updated_date: '2026-06-11 19:16'
labels: []
dependencies: []
references:
  - TASK-009
  - TASK-010
priority: high
ordinal: 18000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Replace the current manual IndexedDB persistence (via the idb package and full-replace strategy) with Yjs CRDT for local-first storage. This lays the groundwork for future multi-device sync by using Yjs as the data layer, enabling conflict-free merging of offline edits.

What changes:
- Drop idb package, use yjs + y-indexeddb instead
- Rewrite stores.ts to use a Y.Map backing instead of a plain array writable
- Delete db.ts (replaced by y-indexeddb persistence)
- Update App.svelte to remove hydration from old IndexedDB
- No component changes — store API surface remains identical

Reasoning (grill-me outcome):
- Yjs is the right CRDT choice — battle-tested client lib, simple blob-relay server pattern
- Data model: Y.Map of Y.Maps keyed by item ID (order/grouping is per-user in the UI)
- Server will eventually implement the Yjs sync protocol over Charge SSE (custom Go blob relay)
- Old IndexedDB data is orphaned — accepted start fresh for migration
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 idb package is removed from dependencies
- [ ] #2 yjs and y-indexeddb packages are added
- [ ] #3 db.ts is deleted — no remaining manual IndexedDB calls
- [ ] #4 stores.ts uses a Y.Doc with Y.Map(items) as the source of truth
- [ ] #5 Store mutations go through Yjs (add, toggle, delete, deleteChecked) — writable store is driven by Yjs observe
- [ ] #6 Persistence is handled automatically by y-indexeddb — no manual subscribe-and-write-back
- [ ] #7 App.svelte no longer imports from db.ts or calls items.hydrate()
- [ ] #8 All existing components work unchanged (AddItem, ItemList, ItemRow, SwipeAction)
- [ ] #9 Items created with crypto.randomUUID() instead of timestamp-based ID
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Install dependencies: pnpm add yjs y-indexeddb, pnpm remove idb (after code is migrated)
2. Rewrite stores.ts:
   - Create Y.Doc + IndexeddbPersistence at module level
   - Create Y.Map(items) as the data model
   - observe() handler rebuilds Item array and calls writable.set()
   - All mutation methods (addItem, toggleItem, deleteItem, deleteCheckedItems) operate on Y.Map
   - Switch createItemID to crypto.randomUUID()
   - Remove hydrate(), keep subscribe, remove set/update from exposed API
3. Delete db.ts (no longer needed)
4. Update App.svelte:
   - Remove import of getAllItems, hydrate call, db.ts import
   - Keep the setTimeout for body height
5. Verify: pnpm check, pnpm dev, manual test of add/toggle/swipe-delete/clear-checked
6. Remove idb from package.json
<!-- SECTION:PLAN:END -->
