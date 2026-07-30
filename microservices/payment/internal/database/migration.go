package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AutoMigrate(pool *pgxpool.Pool) error {

	_, err := pool.Exec(context.Background(), `
CREATE TABLE IF NOT EXISTS payments (
	id BIGSERIAL PRIMARY KEY,
	order_id BIGINT NOT NULL,
	user_id BIGINT NOT NULL,
	amount DOUBLE PRECISION NOT NULL,
	status VARCHAR(50) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`)

	return err
}
