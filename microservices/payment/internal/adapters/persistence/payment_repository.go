package persistence

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/payment/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreatePayment(
	ctx context.Context,
	payment *domain.Payment,
) error {

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO payments (
			order_id,
			user_id,
			amount,
			status
		)
		VALUES ($1,$2,$3,$4)
		RETURNING
			id,
			created_at,
			updated_at
		`,
		payment.OrderID,
		payment.UserID,
		payment.Amount,
		string(payment.Status),
	).Scan(
		&payment.ID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	return err
}

func (r *Repository) GetPayment(
	ctx context.Context,
	id uint,
) (*domain.Payment, error) {

	payment := &domain.Payment{}

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			order_id,
			user_id,
			amount,
			status,
			created_at,
			updated_at
		FROM payments
		WHERE id=$1
		`,
		id,
	).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.Amount,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *Repository) UpdatePayment(
	ctx context.Context,
	payment *domain.Payment,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		UPDATE payments
		SET
			status=$1,
			updated_at=NOW()
		WHERE id=$2
		`,
		string(payment.Status),
		payment.ID,
	)

	return err
}
