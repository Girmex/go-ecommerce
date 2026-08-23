package persistence

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/domain"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/port"
)

var _ port.PaymentRepository = (*PaymentRepository)(nil)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) Create(
	ctx context.Context,
	payment *domain.Payment,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO payments (
			id,
			order_id,
			user_id,
			amount,
			method,
			status,
			gateway_txn_ref,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`,
		payment.ID,
		payment.OrderID,
		payment.UserID,
		payment.Amount,
		payment.Method,
		payment.Status,
		payment.GatewayTxnRef,
		payment.CreatedAt,
		payment.UpdatedAt,
	)

	return err
}

func (r *PaymentRepository) GetByID(
	ctx context.Context,
	id string,
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
			method,
			status,
			gateway_txn_ref,
			created_at,
			updated_at
		FROM payments
		WHERE id = $1
		`,
		id,
	).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.Amount,
		&payment.Method,
		&payment.Status,
		&payment.GatewayTxnRef,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPaymentNotFound
		}

		return nil, err
	}

	return payment, nil
}

func (r *PaymentRepository) List(
	ctx context.Context,
) ([]*domain.Payment, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			order_id,
			user_id,
			amount,
			method,
			status,
			gateway_txn_ref,
			created_at,
			updated_at
		FROM payments
		ORDER BY created_at DESC
		`,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	payments := make([]*domain.Payment, 0)

	for rows.Next() {
		payment := &domain.Payment{}

		err := rows.Scan(
			&payment.ID,
			&payment.OrderID,
			&payment.UserID,
			&payment.Amount,
			&payment.Method,
			&payment.Status,
			&payment.GatewayTxnRef,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *PaymentRepository) Update(
	ctx context.Context,
	payment *domain.Payment,
) error {

	result, err := r.db.Exec(
		ctx,
		`
		UPDATE payments
		SET
			order_id = $2,
			user_id = $3,
			amount = $4,
			method = $5,
			status = $6,
			gateway_txn_ref = $7,
			updated_at = $8
		WHERE id = $1
		`,
		payment.ID,
		payment.OrderID,
		payment.UserID,
		payment.Amount,
		payment.Method,
		payment.Status,
		payment.GatewayTxnRef,
		payment.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPaymentNotFound
	}

	return nil
}