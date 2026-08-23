package database

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(
	ctx context.Context,
	databaseURL string,
) (*pgxpool.Pool, error) {

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func RunMigrations(
	databaseURL string,
	migrationsPath string,
) error {

	m, err := migrate.New(
		"file://"+migrationsPath,
		databaseURL,
	)
	if err != nil {
		return err
	}

	defer func() {
		_, _ = m.Close()
	}()

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}