package migrations

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpAndDownInIsolatedSchema(t *testing.T) {
	databaseURL := os.Getenv("LAGA_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set LAGA_TEST_DATABASE_URL to run PostgreSQL migration tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	schema := "laga_migration_test_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	config, err := pgxpool.ParseConfig(databaseURL)
	require.NoError(t, err)
	config.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	_, err = pool.Exec(ctx, "CREATE SCHEMA "+pgx.Identifier{schema}.Sanitize())
	require.NoError(t, err)
	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := pool.Exec(cleanupCtx, "DROP SCHEMA "+pgx.Identifier{schema}.Sanitize()+" CASCADE")
		assert.NoError(t, err)
	})

	require.NoError(t, Up(pool))
	// The migration adapter has closed; the application's pool must still work.
	require.NoError(t, pool.Ping(ctx))
	for _, table := range []string{"users", "credentials", "registration_attempts", "login_attempts", "sessions"} {
		var exists bool
		err := pool.QueryRow(ctx, "SELECT to_regclass($1) IS NOT NULL", table).Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists, "missing migrated table %s", table)
	}

	db := stdlib.OpenDBFromPool(pool)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })
	require.NoError(t, goose.Down(db, "."))
	var tables int
	err = pool.QueryRow(
		ctx,
		"SELECT count(*) FROM information_schema.tables WHERE table_schema = $1 AND table_name != 'goose_db_version'",
		schema,
	).Scan(&tables)
	require.NoError(t, err)
	assert.Zero(t, tables)
	require.NoError(t, Up(pool))
}
