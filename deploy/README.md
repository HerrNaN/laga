# Scaleway deployment

`bootstrap` owns the project, IAM identities, and Pulumi state bucket. `infra`
owns the application container and its Serverless SQL PostgreSQL database.

The database scales from zero to one CPU. It is protected against accidental
deletion/replacement by Pulumi. Removing it requires explicitly disabling that
protection first. The backend runs its embedded schema migrations on startup.

## One-time database setup

1. Apply `deploy/bootstrap` with your existing administrator credentials and
   backend. This grants the CI identity `ServerlessSQLDatabaseFullAccess` and
   creates a separate database-only application/key. CI does not need IAM
   administration permissions.

   ```sh
   pulumi -C deploy/bootstrap up --stack prod
   ```

2. Copy the new outputs into the infra stack configuration, using the existing
   backend and passphrase for each stack:

   ```sh
   pulumi -C deploy/infra config set --stack prod laga:databaseApplicationId \
     "$(pulumi -C deploy/bootstrap stack output --stack prod database_application_id)"
   pulumi -C deploy/infra config set --stack prod --secret laga:databaseSecretKey \
     "$(pulumi -C deploy/bootstrap stack output --stack prod --show-secrets database_secret_key)"
   ```

   Do not run the secret-copy command with shell tracing enabled. Commit the
   resulting `deploy/infra/Pulumi.prod.yaml`; the secret is encrypted, and CI
   needs this configuration. When rotating the runtime key, repeat this copy
   and redeploy the container.

3. Preview and apply the infra stack (or use the existing deployment workflow):

   ```sh
   pulumi -C deploy/infra preview --stack prod
   pulumi -C deploy/infra up --stack prod
   ```

`DATABASE_URL` is assembled from the provisioned endpoint and runtime IAM
credentials, requires TLS, and is passed via secret environment variables. It
is not exported as a plaintext stack output. The runtime identity uses
`ServerlessSQLDatabaseReadWrite`, rather than the data-only permission set,
because startup migrations need schema administration.

The backend also requires `WEBAUTHN_RP_ID` and `WEBAUTHN_ORIGIN`; those must be
configured for the public application hostname before it can start successfully.
