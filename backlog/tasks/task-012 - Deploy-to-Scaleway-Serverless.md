---
id: TASK-012
title: Deploy to Scaleway Serverless
status: Done
assignee: []
created_date: '2026-05-23 14:06'
updated_date: '2026-05-29 16:44'
labels:
  - deploy
dependencies: []
priority: high
ordinal: 12000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Deploy laga as a single Serverless Container on Scaleway, accessible via the default *.scaleway.io subdomain with HTTPS.

The app works offline-first with IndexedDB, so no PostgreSQL needed for initial deployment. DB persistence will be added in a separate task (TASK-013). Custom domain, CDN caching, and WAF will be added in TASK-014 (Edge Services).

## Architecture

User → Serverless Container (Go binary with embedded SPA, served over HTTPS via *.scaleway.io)

## Components

### Scaleway Serverless Containers
- Deploy Go binary as Docker image (FROM scratch, copy binary, expose 8080)
- No code refactoring needed — existing http.ListenAndServe pattern works as-is
- Scales to zero when idle, wakes on request (cold start ~1-5s)
- Pricing: €0.10/100k GB-s compute + €1.0/100k vCPU-s, generous free tier
- Auto-provisioned HTTPS endpoint at *.scaleway.io

## Code Changes Required

### Go backend
- Add Cache-Control headers for static assets:
  - /assets/* → Cache-Control: public, max-age=31536000, immutable
  - index.html → Cache-Control: no-cache (PWA must revalidate)
  - /api/* → Cache-Control: no-store
  - These headers will also benefit Edge Services CDN when added later

### Dockerfile
- Multi-stage build: compile Go binary, copy to scratch image
- Result is essentially just the statically-linked binary

### PWA
- Already using vite-plugin-pwa — service worker caches static assets after first visit
- No changes needed

### Deployment config
- Scaleway CLI or Terraform for provisioning:
  - Container namespace + deployment

## Estimated Monthly Cost (light traffic)
- Container: near-zero (free tier covers most traffic)
- Total: ~€0-1/month

## Cold Start Strategy
- First visit: cold container start ~1-5s
- Return visits: PWA service worker serves from cache → instant
- Without Edge Services, no CDN caching layer — all requests hit the container
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Scaleway Serverless Container deployed with Go binary
- [ ] #2 Dockerfile exists with multi-stage build producing minimal image
- [ ] #3 Go backend emits correct Cache-Control headers (immutable for /assets/*, no-cache for index.html, no-store for /api/*)
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
---
## Deployment Infrastructure

### Stack structure
- **deploy/bootstrap/** — Creates Scaleway Project, IAM Application, Policy, and API Key. Uses local Pulumi state (admin credentials).
- **deploy/infra/** — Creates Container Namespace and Serverless Container. Uses S3 state backend (deploy key credentials).

### Bootstrap stack outputs
- `project_id` — Scaleway project ID
- `organization_id` — Scaleway organization ID
- `access_key` / `secret_key` — Deploy pipeline API key (scoped to laga project)

### IAM Policy
- Application: laga-deploy
- Permission sets: ContainersFullAccess, ObjectStorageFullAccess
- Scoped to: laga project only

### Deployment workflow
1. Bootstrap (one-time, admin creds): `cd deploy/bootstrap && pulumi login --local && pulumi up`
2. Create state bucket (manual, one-time): `scw bucket create laga-pulumi-state --project-id <project_id>`
3. Infra (uses deploy key): `cd deploy/infra && pulumi login s3://laga-pulumi-state && pulumi up`

### Customize image tag
`pulumi config set laga:imageTag <sha>` or `pulumi up --config laga:imageTag=<sha>`

---
## Two-stack deployment architecture

### Bootstrap stack (`deploy/bootstrap/`)
- Creates: Scaleway Project, IAM Application, IAM Policy, API Key, infra state bucket + bucket policy
- State backend: `s3://laga-pulumi-state-bootstrap` (manually created)
- Bucket policy on bootstrap bucket: admin user only (set manually)

### Infra stack (`deploy/infra/`)
- Creates: Container Namespace, Serverless Container
- State backend: `s3://laga-pulumi-state-infra` (created by bootstrap stack)
- Bucket policy on infra bucket: deploy app only (set by Pulumi)

### Access isolation
- Deploy key cannot access bootstrap state (denied by bucket policy omission)
- Admin user cannot access infra state (denied by bucket policy omission)
- Must use deploy key for infra operations
<!-- SECTION:NOTES:END -->
