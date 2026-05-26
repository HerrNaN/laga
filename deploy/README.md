# Deploy

Pulumi-based infrastructure for deploying laga to Scaleway Serverless Containers.

## Architecture

Two stacks, each with its own state backend and credentials:

| Stack | State Backend | Credentials | Purpose |
|-------|--------------|-------------|---------|
| `bootstrap` | Scaleway Object Storage (`s3://laga-pulumi-state-admin`) | Admin Scaleway credentials | Creates project, IAM, API key, infra state bucket + policy |
| `infra` | Scaleway Object Storage (`s3://laga-pulumi-state`) | Deploy key (scoped to laga project) | Creates container namespace and serverless container |

```
Scaleway Organization
  └── Project: laga
        ├── Container Namespace: laga
        │     └── Container: laga (ghcr.io/herrnan/laga:latest, port 8080)
        ├── Object Storage Bucket: laga-pulumi-state-admin (manual)
        ├── Object Storage Bucket: laga-pulumi-state (IaC)
        │     └── Bucket Policy: admin user + laga-deploy application
        └── IAM
              ├── Application: laga-deploy
              ├── Policy: laga-deploy-policy (ContainersFullAccess + ObjectStorageFullAccess, scoped to laga project)
              └── API Key: laga-deploy-key (default project = laga)
```

## Prerequisites

- [mise](https://mise.jdx.io/) or manual installs of: Go, Pulumi, Scaleway CLI (`scw`), AWS CLI (`aws`)
- Scaleway account with admin credentials

## Setup

### 0. Create bootstrap state bucket (one-time, manual)

This bucket is created outside Pulumi so Pulumi can use it for bootstrap state immediately.

```bash
scw object bucket create laga-pulumi-state-admin region=fr-par
```

### 1. Bootstrap (one-time)

```bash
cd deploy/bootstrap

# Configure Scaleway provider (admin credentials)
export SCW_ACCESS_KEY=<admin-access-key>
export SCW_SECRET_KEY=<admin-secret-key>
export SCW_ORGANIZATION_ID=<your-org-id>
export SCW_DEFAULT_PROJECT_ID=<admin-project-id>

# Configure Pulumi S3 backend credentials (same admin key)
export AWS_ACCESS_KEY_ID=$SCW_ACCESS_KEY
export AWS_SECRET_ACCESS_KEY=$SCW_SECRET_KEY

# Set admin user ID (find with: scw iam user list)
pulumi config set laga:adminUserId <your-user-id>

pulumi up
```

Note the outputs: `project_id`, `organization_id`, `access_key`, `secret_key`, `infra_bucket`.

### 2. Infra (deploy key credentials)

```bash
cd deploy/infra

# Configure Scaleway provider (deploy key from bootstrap outputs)
export SCW_ACCESS_KEY=<deploy-access-key>
export SCW_SECRET_KEY=<deploy-secret-key>
export SCW_ORGANIZATION_ID=<org-id from bootstrap outputs>
export SCW_DEFAULT_PROJECT_ID=<project-id from bootstrap outputs>

# Configure Pulumi S3 backend credentials (same deploy key)
export AWS_ACCESS_KEY_ID=$SCW_ACCESS_KEY
export AWS_SECRET_ACCESS_KEY=$SCW_SECRET_KEY

pulumi up
```

## Deploying a specific version

```bash
# Override image tag for one deployment
pulumi up --config laga:imageTag=abc1234

# Or persist the tag in config
pulumi config set laga:imageTag abc1234
pulumi up
```

## Mise tasks

```bash
mise bootstrap-preview   # Preview bootstrap stack changes
mise bootstrap-up         # Apply bootstrap stack
mise infra-preview        # Preview infra stack changes
mise infra-up             # Apply infra stack
```

## Rolling back

```bash
cd deploy/infra
pulumi config set laga:imageTag <previous-sha>
pulumi up
```

## Environment variables reference

| Variable | Used By | Source |
|----------|---------|--------|
| `SCW_ACCESS_KEY` | Pulumi Scaleway provider | Admin key (bootstrap) / Deploy key (infra) |
| `SCW_SECRET_KEY` | Pulumi Scaleway provider | Admin key (bootstrap) / Deploy key (infra) |
| `SCW_ORGANIZATION_ID` | Pulumi Scaleway provider | From bootstrap outputs |
| `SCW_DEFAULT_PROJECT_ID` | Pulumi Scaleway provider | Admin project (bootstrap) / laga project (infra) |
| `AWS_ACCESS_KEY_ID` | Pulumi S3 backend | Admin key (bootstrap) / Deploy key (infra) |
| `AWS_SECRET_ACCESS_KEY` | Pulumi S3 backend | Admin secret (bootstrap) / Deploy secret (infra) |
