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
	panic("implement me")
}

func (r *Repository) GetOrdersByUser(
	ctx context.Context,
	userID uint,
) ([]domain.Order, error) {
	panic("implement me")
}

func (r *Repository) UpdateOrder(
	ctx context.Context,
	order *domain.Order,
) error {
	panic("implement me")
}

func (r *Repository) DeleteOrder(
	ctx context.Context,
	orderID uint,
) error {
	panic("implement me")
}