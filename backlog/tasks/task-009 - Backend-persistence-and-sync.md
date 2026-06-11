---
id: TASK-009
title: Backend persistence and sync
status: To Do
assignee: []
created_date: '2026-05-15 08:52'
updated_date: '2026-06-11 19:17'
labels: []
dependencies: []
ordinal: 9000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add a backend database (PostgreSQL) with API for persistent storage. Implement client-server sync for multi-device support.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Backend API supports CRUD for items
- [ ] #2 Client syncs local state to server
- [ ] #3 Server state can be restored on new devices
<!-- AC:END -->

## Comments

<!-- COMMENTS:BEGIN -->
created: 2026-06-11 19:17
---
Superseded by TASK-018 — Yjs handles client-side CRDT persistence; backend blob relay moved to Phase 2 after TASK-018 is complete
---
<!-- COMMENTS:END -->
