# Deploy

Pulumi-based infrastructure for deploying laga to Scaleway Serverless Containers.

## Architecture

Two stacks, each with its own state backend and credentials:

| Stack | State Backend | Credentials | Purpose |
|-------|--------------|-------------|---------|
| `bootstrap` | Scaleway Object Storage (`s3://laga-pulumi-state-admin`) | Admin Scaleway credentials | Creates project, IAM, API key, infra state bucket + policy |

```
Scaleway Organization
  └── Project: laga
        ├── Container Namespace: laga
        │     └── Container: laga (ghcr.io/herrnan/laga:latest, port 8080)
        ├── Object Storage Bucket: laga-pulumi-state (IaC)
        │     └── Bucket Policy: laga-deploy app only
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
scw bucket create laga-pulumi-state-admin

# Apply a bucket policy allowing only your admin user.
# Find your user ID with: scw iam user info
aws s3api put-bucket-policy --bucket laga-pulumi-state-admin \
  --endpoint-url https://s3.fr-par.scw.cloud \
  --policy '{
    "Version": "2023-04-17",
    "Statement": [{
      "Sid": "AdminOnly",
      "Effect": "Allow",
      "Principal": { "SCW": "user_id:<YOUR_USER_ID>" },
      "Action": ["s3:*"],
      "Resource": ["laga-pulumi-state-admin", "laga-pulumi-state-admin/*"]
    }]
  }'
```

### 1. Bootstrap (one-time)

```bash
cd deploy/bootstrap

# Configure Pulumi to use the S3 state backend
export AWS_ACCESS_KEY_ID=<admin-access-key>
export AWS_SECRET_ACCESS_KEY=<admin-secret-key>
export AWS_ENDPOINT_URL_S3=https://s3.fr-par.scw.cloud
export AWS_REGION=fr-par

pulumi login 's3://laga-pulumi-state-admin'

# Configure Scaleway provider
export SCW_ACCESS_KEY=<admin-access-key>
export SCW_SECRET_KEY=<admin-secret-key>
export SCW_ORGANIZATION_ID=<your-org-id>
export SCW_DEFAULT_PROJECT_ID=<admin-project-id>

pulumi up
```

Note the outputs: `project_id`, `organization_id`, `access_key`, `secret_key`, `infra_bucket`.

## Mise tasks

```bash
mise bootstrap-preview   # Preview bootstrap stack changes
mise bootstrap-up         # Apply bootstrap stack
```
