---
id: TASK-004
title: Delete individual items
status: Done
assignee: []
created_date: '2026-05-15 08:52'
updated_date: '2026-05-15 09:38'
labels: []
dependencies: []
ordinal: 4000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add the ability to remove a single item from the shopping list. Include a delete button or swipe gesture on each list item.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Each item has a delete action
- [x] #2 Deleting an item removes it from the store and IndexedDB
- [x] #3 Deletion requires confirmation or is undoable
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Refined to swipe-to-delete gesture with undo. Swiping left or right reveals delete action. Releasing past threshold marks item for deletion with 3s undo window. No stacking of pending deletions — starting a new swipe commits any existing pending deletion.

Refactored to CSS scroll-snap approach inspired by https://dev.to/eduferfer/how-to-css-table-swipe-interaction-3ih3. Uses native overflow-x scroll with scroll-snap-type: x mandatory instead of manual pointer-event tracking. Three sibling panels per item (delete-left, content, delete-right), each min-width: 100%. JS reduced to single onscrollend handler that detects swipe direction from scrollLeft. Much simpler CSS (~60 lines vs ~70, but dramatically simpler JS).
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented swipe-to-delete with undo. Added deleteItem store method that auto-persists to IndexedDB via existing subscription. Rewrote ItemList.svelte with pointer-event-based swipe gesture (both directions supported). Swiping past 100px threshold reveals a pending-delete state with 3s undo window. No stacking: starting a swipe on a new item immediately commits any existing pending deletion. pnpm check passes with zero errors.
<!-- SECTION:FINAL_SUMMARY:END -->
