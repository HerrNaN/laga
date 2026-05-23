---
id: TASK-014
title: Add Scaleway Edge Services for custom domain and CDN
status: To Do
assignee: []
created_date: '2026-05-23 14:15'
updated_date: '2026-05-23 14:19'
labels:
  - deploy
dependencies:
  - TASK-012
priority: low
ordinal: 14000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add Scaleway Edge Services in front of the Serverless Container to provide custom domain, TLS, CDN caching, and WAF protection.

Depends on TASK-012 (container deployment). This is optional — the app works on *.scaleway.io without it. Add this when you want a custom domain or CDN caching for reduced cold-start impact.

## Architecture (extends TASK-012)

User → Edge Services (custom domain, TLS, CDN caching, WAF) → Serverless Container (Go binary with embedded SPA)

## What Edge Services Adds

### Custom Domain
- Your own domain (e.g., laga.yourdomain.com) with automatic TLS
- Replaces the default *.scaleway.io subdomain

### CDN Caching
- Static assets cached at edge locations
- Configured cache rules:
  - /assets/* → aggressive caching (immutable hashed filenames)
  - /api/* → no cache
  - / → short cache or stale-while-revalidate
- Reduces container invocations and cold-start impact for new users

### WAF Protection
- Web Application Firewall for DDoS and threat protection

## Deployment Config
- Scaleway CLI or Terraform for provisioning:
  - Edge Services (TLS, CDN, WAF)
  - DNS configuration for custom domain

## Cost
- ~€1-2/month additional
- Total with TASK-012: ~€1-3/month

## Cold Start Impact
- Before Edge Services: every uncached request hits container, cold start ~1-5s for new users
- After Edge Services: CDN caches static assets at edge, reducing container invocations significantly
- PWA service worker still handles client-side caching regardless
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Scaleway Edge Services configured with custom domain and TLS
- [ ] #2 CDN cache rules configured (/assets/* → immutable, /api/* → no-cache, / → stale-while-revalidate)
- [ ] #3 DNS records point custom domain to Edge Services
- [ ] #4 Static assets are cached at CDN edge after first request (verified via cache headers)
<!-- AC:END -->
