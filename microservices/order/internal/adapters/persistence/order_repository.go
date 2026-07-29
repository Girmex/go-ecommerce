package persistence

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/order/internal/domain"
	"github.com/Girmex/go-ecommerce/microservices/order/internal/ports"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ ports.OrderRepository = (*Repository)(nil)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateOrder(
	ctx context.Context,
	order *domain.Order,
) error {

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(
		ctx,
		`
		INSERT INTO orders (
			user_id,
			status,
			total_price
		)
		VALUES ($1, $2, $3)
		RETURNING id
		`,
		order.UserID,
		string(order.Status),
		order.TotalPrice,
	).Scan(&order.ID)

	if err != nil {
		return err
	}

	for i := range order.Items {

		order.Items[i].OrderID = order.ID

		err = tx.QueryRow(
			ctx,
			`
			INSERT INTO order_items (
				order_id,
				product_id,
				quantity,
				unit_price
			)
			VALUES ($1, $2, $3, $4)
			RETURNING id
			`,
			order.Items[i].OrderID,
			order.Items[i].ProductID,
			order.Items[i].Quantity,
			order.Items[i].UnitPrice,
		).Scan(&order.Items[i].ID)

		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *Repository) GetOrder(
	ctx context.Context,
	orderID uint,
) (*domain.Order, error) {

	order := &domain.Order{}

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			user_id,
			status,
			total_price,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
		`,
		orderID,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalPrice,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			order_id,
			product_id,
			quantity,
			unit_price
		FROM order_items
		WHERE order_id = $1
		`,
		order.ID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var item domain.OrderItem

		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
		); err != nil {
			return nil, err
		}

		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (r *Repository) GetOrdersByUser(
	ctx context.Context,
	userID uint,
) ([]domain.Order, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id
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

	var orders []domain.Order

	for rows.Next() {

		var id uint

		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		order, err := r.GetOrder(ctx, id)
		if err != nil {
			return nil, err
		}

		orders = append(orders, *order)
	}

	return orders, nil
}

func (r *Repository) UpdateOrder(
	ctx context.Context,
	order *domain.Order,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		UPDATE orders
		SET
			status = $1,
			updated_at = NOW()
		WHERE id = $2
		`,
		string(order.Status),
		order.ID,
	)

	return err
}

func (r *Repository) DeleteOrder(
	ctx context.Context,
	orderID uint,
) error {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		`
		DELETE FROM order_items
		WHERE order_id = $1
		`,
		orderID,
	)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		`
		DELETE FROM orders
		WHERE id = $1
		`,
		orderID,
	)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}