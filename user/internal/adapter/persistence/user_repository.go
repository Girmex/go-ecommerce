package persistence

import (
	"context"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *domain.User,
) error {
	query := `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	return r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID)
}

func (r *UserRepository) GetByID(
	ctx context.Context,
	id uint,
) (*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}