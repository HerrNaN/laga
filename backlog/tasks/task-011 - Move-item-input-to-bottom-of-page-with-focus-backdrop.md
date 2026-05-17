---
id: TASK-011
title: Move item input to bottom of page with focus backdrop
status: Done
assignee: []
created_date: '2026-05-17 12:00'
updated_date: '2026-05-17 16:14'
labels: []
dependencies: []
ordinal: 11000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Relocate the 'Add item' input bar from the top of the app (below the header) to the bottom of the page. When the input is focused, the rest of the UI (header + item list) should be dimmed by a semi-transparent backdrop overlay. Tapping the backdrop unfocuses the input and dismisses the keyboard. The input bar must remain visible above the mobile virtual keyboard when focused. The header remains sticky at the top and the input bar sits sticky at the bottom, framing the scrollable item list between them. The whole page scrolls naturally (Option X) for reliable mobile keyboard handling. The Add button stays inside the input bar alongside the input field.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Input bar is rendered below the item list at the bottom of the page
- [x] #2 Header remains position: sticky at the top with z-index: 1
- [x] #3 Input bar has position: sticky; bottom: 0; z-index: 20 with a top box-shadow for floating appearance
- [x] #4 When input is focused, a semi-transparent backdrop (rgba(0,0,0,0.4)) covers the entire viewport
- [x] #5 Backdrop has z-index: 10 and unfocuses the input on click/tap
- [x] #6 Input bar stays above the mobile virtual keyboard when focused (whole-page scroll, no overflow: hidden on main)
- [x] #7 Viewport meta tag updated with interactive-widget=resizes-content
- [x] #8 App renders header -> item list -> input bar in DOM order
- [x] #9 No extra bottom padding is added to the list
- [x] #10 Add button remains inside the sticky input bar, not the backdrop
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Update viewport meta tag in index.html with interactive-widget=resizes-content\n2. Reorder App.svelte DOM to header -> ItemList -> AddItem\n3. Add sticky positioning and z-index to header and AddItem section\n4. Add box-shadow to AddItem bar\n5. Add focus state tracking in AddItem.svelte and expose to App.svelte\n6. Add backdrop overlay to App.svelte with click-to-unfocus logic\n7. Verify on mobile that keyboard opens and input stays visible
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Bugfix: Changed Add button from onclick to onpointerup for reliable mobile touch handling. Bugfix: Replaced conditional {#if} rendering of backdrop with permanent DOM element toggled via CSS opacity/pointer-events to prevent DOM-removal race condition during touch sequences that caused Add button to be unresponsive on mobile.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented TASK-011: Moved AddItem input bar from top to bottom of page. App.svelte now renders header -> ItemList -> AddItem in DOM order. Header is position: sticky at top with z-index: 1. Input bar is position: sticky at bottom with z-index: 20 and a top box-shadow. Focus state is tracked via onFocusChange callback from AddItem to App. When focused, a semi-transparent fixed backdrop (z-index: 10) covers the viewport and unfocuses the input on click. Viewport meta tag updated with interactive-widget=resizes-content for reliable mobile keyboard handling. No extra bottom padding added to list. Build passes cleanly.
<!-- SECTION:FINAL_SUMMARY:END -->
