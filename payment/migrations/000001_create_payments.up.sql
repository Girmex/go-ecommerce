CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY,
    order_id VARCHAR(255) NOT NULL,
    user_id BIGINT NOT NULL,
    amount NUMERIC(12, 2) NOT NULL,
    method VARCHAR(50) NOT NULL,
    status VARCHAR(30) NOT NULL,
    gateway_txn_ref VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_payments_order_id
    ON payments(order_id);

CREATE INDEX IF NOT EXISTS idx_payments_user_id
    ON payments(user_id);

CREATE INDEX IF NOT EXISTS idx_payments_status
    ON payments(status);