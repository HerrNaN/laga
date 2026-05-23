---
id: TASK-013
title: Add PostgreSQL persistence via Scaleway SDB
status: To Do
assignee: []
created_date: '2026-05-23 14:12'
updated_date: '2026-05-23 14:19'
labels:
  - deploy
dependencies:
  - TASK-012
priority: medium
ordinal: 13000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add server-side persistence using Scaleway Serverless SQL Database (PostgreSQL) and connect it to the Go backend.

Depends on TASK-012 (deploy infrastructure). The app currently works offline-first with IndexedDB; this task adds server-side persistence for data sync and sharing.

## Architecture (extends TASK-012)

User → Edge Services → Serverless Container (Go binary) → Serverless SQL Database (PostgreSQL)

## Components

### Scaleway Serverless SQL Database
- PostgreSQL-compatible, scales compute to zero when idle
- Pricing: €0.136/vCPU/hour (compute, 5-min min window) + €0.000272/GB/hour storage
- Set min/max vCPU thresholds (min 0, max 2-4) for cost control
- Region: Paris (fr-par)
- Connect via DATABASE_URL env var

## Code Changes Required

### Go backend
- Add PostgreSQL driver (pgx) with connection pooling
- Add DATABASE_URL env var handling
- Add database migration tooling (goose, golang-migrate, or embedded SQL)
- Add API routes for data sync/persistence

### Deployment config
- Terraform or Scaleway CLI for provisioning:
  - Serverless SQL Database
  - IAM application + API key for DB access
  - Add DATABASE_URL secret to container environment

## Estimated Added Monthly Cost (light traffic)
- Database: ~€1-3/month (1-5 GB storage, intermittent compute)
- Total with TASK-012: ~€2-5/month
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Go backend has pgx driver with connection pooling and DATABASE_URL env var
- [ ] #2 Scaleway Serverless SQL Database provisioned with min/max vCPU thresholds
- [ ] #3 Go binary connects to SDB and serves requests end-to-end
- [ ] #4 Database migration tooling is set up
<!-- AC:END -->
