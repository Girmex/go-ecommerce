package persistence

import (
	"context"
	"errors"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(
	ctx context.Context,
	order *domain.Order,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Create order
	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO orders (
			id,
			user_id,
			total,
			status,
			payment_id,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		`,
		order.ID,
		order.UserID,
		order.Total,
		order.Status,
		nullString(order.PaymentID),
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Create order items
	for _, item := range order.Items {
		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO order_items (
				order_id,
				product_id,
				quantity,
				unit_price
			)
			VALUES ($1, $2, $3, $4)
			`,
			order.ID,
			item.ProductID,
			item.Quantity,
			item.UnitPrice,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *OrderRepository) GetByID(
	ctx context.Context,
	id string,
) (*domain.Order, error) {
	order := &domain.Order{}

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			user_id,
			total,
			status,
			payment_id,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
		`,
		id,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.Total,
		&order.Status,
		&order.PaymentID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}

		return nil, err
	}

	items, err := r.getItems(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	order.Items = items

	return order, nil
}

func (r *OrderRepository) ListByUser(
	ctx context.Context,
	userID string,
) ([]*domain.Order, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			user_id,
			total,
			status,
			payment_id,
			created_at,
			updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order

	for rows.Next() {
		order := &domain.Order{}

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Total,
			&order.Status,
			&order.PaymentID,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		items, err := r.getItems(ctx, order.ID)
		if err != nil {
			return nil, err
		}

		order.Items = items

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) Update(
	ctx context.Context,
	order *domain.Order,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(
		ctx,
		`
		UPDATE orders
		SET
			total = $1,
			status = $2,
			payment_id = $3,
			updated_at = $4
		WHERE id = $5
		`,
		order.Total,
		order.Status,
		nullString(order.PaymentID),
		order.UpdatedAt,
		order.ID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrOrderNotFound
	}

	// Replace existing items.
	_, err = tx.Exec(
		ctx,
		`
		DELETE FROM order_items
		WHERE order_id = $1
		`,
		order.ID,
	)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO order_items (
				order_id,
				product_id,
				quantity,
				unit_price
			)
			VALUES ($1, $2, $3, $4)
			`,
			order.ID,
			item.ProductID,
			item.Quantity,
			item.UnitPrice,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *OrderRepository) getItems(
	ctx context.Context,
	orderID string,
) ([]domain.OrderItem, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			product_id,
			quantity,
			unit_price
		FROM order_items
		WHERE order_id = $1
		ORDER BY id
		`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.OrderItem, 0)

	for rows.Next() {
		var item domain.OrderItem

		err := rows.Scan(
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func nullString(value string) any {
	if value == "" {
		return nil
	}

	return value
}