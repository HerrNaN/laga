# Serverless SQL deployment TODO

Run the commands below from the repository root.

## 1. Update bootstrap

- [ ] Use administrator Scaleway credentials and the existing Pulumi backend/passphrase.
- [ ] Apply the bootstrap stack to grant CI SQL permissions and create database runtime credentials:

```sh
pulumi -C deploy/bootstrap up --stack prod
```

## 2. Configure database credentials

- [ ] Copy the bootstrap outputs into the infra stack configuration:

```sh
pulumi -C deploy/infra config set --stack prod laga:databaseApplicationId \
  "$(pulumi -C deploy/bootstrap stack output --stack prod database_application_id)"

pulumi -C deploy/infra config set --stack prod --secret laga:databaseSecretKey \
  "$(pulumi -C deploy/bootstrap stack output --stack prod --show-secrets database_secret_key)"
```

Do not run these commands with shell tracing (`set -x`) enabled. Use the existing backend and passphrase for each stack.

## 3. Finish WebAuthn environment configuration

- [ ] Wire the following environment variables into the Pulumi container configuration:
  - `WEBAUTHN_RP_ID`: the public hostname, without `https://`.
  - `WEBAUTHN_ORIGIN`: `https://<your-public-hostname>`.

These are not wired into Pulumi yet and are required for successful application startup.

## 4. Deploy and verify

- [ ] Preview the infrastructure changes:

```sh
pulumi -C deploy/infra preview --stack prod
```

- [ ] Apply the infrastructure changes:

```sh
pulumi -C deploy/infra up --stack prod
```

The backend applies database migrations automatically on startup.

- [ ] Verify the application starts and connects to the database successfully.
- [ ] Commit the updated `deploy/infra/Pulumi.prod.yaml` so CI has the encrypted credentials.
