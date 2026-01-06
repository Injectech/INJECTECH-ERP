package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	domainuser "backend/internal/domain/user"
)

// UserRepository implements user.Repository using PostgreSQL (sqlc-generated queries expected).
type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u domainuser.User) (domainuser.User, error) {
	// TODO: implement with sqlc generated code
	return u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (domainuser.User, error) {
	return domainuser.User{}, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domainuser.User, error) {
	return domainuser.User{}, nil
}

func (r *UserRepository) Update(ctx context.Context, u domainuser.User) error {
	return nil
}

func (r *UserRepository) SoftDelete(ctx context.Context, id string) error {
	return nil
}
