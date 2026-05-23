---
id: TASK-005
title: Delete all checked items
status: Done
assignee: []
created_date: '2026-05-15 08:52'
updated_date: '2026-05-23 06:47'
labels: []
dependencies: []
ordinal: 2000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add a bulk action to remove all checked/completed items from the list at once. This should be easily accessible but not prone to accidental activation.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 A "Clear" text link appears next to the "Checked" heading whenever checked items exist
- [x] #2 Tapping "Clear" immediately and permanently removes all checked items from the store and IndexedDB
<!-- AC:END -->

## Design Decisions
- Trigger placed inside "Checked" section header (contextual, only visible when actionable)
- Text link ("Clear") for unobtrusive, simple UI
- No confirmation step (reduced friction; acceptable risk for grocery items)
- No undo buffer (KISS principle; can be added later if needed)
