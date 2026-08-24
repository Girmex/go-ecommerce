package persistence

import (
	"context"
	"errors"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Create(
	ctx context.Context,
	product *domain.Product,
) error {
	query := `
		INSERT INTO products (
			id,
			user_id,
			name,
			description,
			price,
			stock,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		product.ID,
		product.UserID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.CreatedAt,
		product.UpdatedAt,
	)

	return err
}

func (r *ProductRepository) GetByID(
	ctx context.Context,
	id string,
) (*domain.Product, error) {
	query := `
		SELECT
			id,
			user_id,
			name,
			description,
			price,
			stock,
			created_at,
			updated_at
		FROM products
		WHERE id = $1
	`

	product := &domain.Product{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.UserID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProductNotFound
		}

		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) List(
	ctx context.Context,
) ([]*domain.Product, error) {
	query := `
		SELECT
			id,
			user_id,
			name,
			description,
			price,
			stock,
			created_at,
			updated_at
		FROM products
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*domain.Product, 0)

	for rows.Next() {
		product := &domain.Product{}

		if err := rows.Scan(
			&product.ID,
			&product.UserID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) Update(
	ctx context.Context,
	product *domain.Product,
	userID uint,
) error {
	query := `
		UPDATE products
		SET
			name = $1,
			description = $2,
			price = $3,
			stock = $4,
			updated_at = $5
		WHERE id = $6
		AND user_id = $7
	`

	result, err := r.db.Exec(
		ctx,
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.UpdatedAt,
		product.ID,
		userID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepository) Delete(
	ctx context.Context,
	id string,
	userID uint,
) error {
	query := `
		DELETE FROM products
		WHERE id = $1
		AND user_id = $2
	`

	result, err := r.db.Exec(
		ctx,
		query,
		id,
		userID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}