---
id: TASK-012
title: Deploy to Scaleway Serverless
status: To Do
assignee: []
created_date: '2026-05-23 14:06'
updated_date: '2026-05-23 14:19'
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
