# Local passkeys

Start PostgreSQL:

```sh
docker compose up -d --wait
```

Run `mise run dev` and visit **http://localhost:8080**. Use the Go server's
port (8080), not Vite's port (5173): WebAuthn checks the exact origin. Opening
the app without a session redirects to `/sign-up`, the page that creates a
passkey account. The signed-in account lives on `/user`, and the header icon on
the lists page points at whichever of the two applies (sign-in icon, then user
icon).

Migrations in `backend/internal/app/migrations` run automatically at startup.
The signup flow stores an account, one passkey credential, and a session in
PostgreSQL. Shopping lists are still browser-only; returning-user sign-in is
not implemented yet. Stop the database with `docker compose down`.

After changing `backend/internal/app/auth/queries.sql` or the migration files,
regenerate the Go queries with `mise run generate`.

To run the browser signup test, leave PostgreSQL running and from `frontend`
run `pnpm exec playwright install chromium` once, then `pnpm test:e2e`.
Playwright starts fresh dev servers, so stop anything already running on ports
8080 and 5173. It creates a disposable virtual passkey and leaves its test
account in the local database.
`hk check --all` runs the same test, so it needs the database and Chromium too.
