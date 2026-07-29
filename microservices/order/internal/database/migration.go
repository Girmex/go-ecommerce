package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AutoMigrate(pool *pgxpool.Pool) error {

	_, err := pool.Exec(context.Background(), `
CREATE TABLE IF NOT EXISTS orders (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL,
	status VARCHAR(50) NOT NULL,
	total_price DOUBLE PRECISION NOT NULL DEFAULT 0,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items (
	id BIGSERIAL PRIMARY KEY,
	order_id BIGINT NOT NULL
		REFERENCES orders(id)
		ON DELETE CASCADE,
	product_id BIGINT NOT NULL,
	quantity BIGINT NOT NULL,
	unit_price DOUBLE PRECISION NOT NULL
);
`)

	return err
}