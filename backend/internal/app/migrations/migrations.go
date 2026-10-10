package migrations

import (
	"embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var files embed.FS

func Up(pool *pgxpool.Pool) (err error) {
	// Close only the database/sql adapter after migrations. Its connections are
	// borrowed from the pool, which stays open for the application's lifetime.
	db := stdlib.OpenDBFromPool(pool)
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close migration adapter: %w", closeErr))
		}
	}()

	goose.SetBaseFS(files)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set migration dialect: %w", err)
	}
	if err := goose.Up(db, "."); err != nil && !errors.Is(err, goose.ErrNoNextVersion) {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
