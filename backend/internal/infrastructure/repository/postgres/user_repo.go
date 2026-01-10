package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domainuser "backend/internal/domain/user"
	"backend/internal/infrastructure/repository/postgres/sqlc"
)

// UserRepository implements user.Repository using PostgreSQL.
type UserRepository struct {
	db *pgxpool.Pool
	q  *sqlc.Queries
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db, q: sqlc.New(db)}
}

func (r *UserRepository) Create(ctx context.Context, u domainuser.User) (domainuser.User, error) {
	id, err := uuid.Parse(u.ID)
	if err != nil {
		return domainuser.User{}, err
	}

	_, err = r.q.CreateUser(ctx, sqlc.CreateUserParams{
		ID:             id,
		Email:          u.Email,
		HashedPassword: u.HashedPassword,
		Name:           u.Name,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	})
	if err != nil {
		return domainuser.User{}, err
	}

	if len(u.Roles) > 0 {
		for _, roleName := range u.Roles {
			if _, err := r.db.Exec(ctx, `
				INSERT INTO user_roles (user_id, role_id)
				SELECT $1, id FROM roles WHERE name = $2
			`, id, roleName); err != nil {
				return domainuser.User{}, err
			}
		}
	}

	return u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (domainuser.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domainuser.User{}, err
	}
	u, err := r.q.GetUserByID(ctx, uid)
	if err != nil {
		return domainuser.User{}, err
	}
	return r.withRolesAndPermissions(ctx, u)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domainuser.User, error) {
	u, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return domainuser.User{}, err
	}
	return r.withRolesAndPermissions(ctx, u)
}

func (r *UserRepository) Update(ctx context.Context, u domainuser.User) error {
	uid, err := uuid.Parse(u.ID)
	if err != nil {
		return err
	}
	return r.q.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:             uid,
		Email:          u.Email,
		HashedPassword: u.HashedPassword,
		Name:           u.Name,
		UpdatedAt:      time.Now(),
	})
}

func (r *UserRepository) SoftDelete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.SoftDeleteUser(ctx, uid)
}

func (r *UserRepository) withRolesAndPermissions(ctx context.Context, u sqlc.User) (domainuser.User, error) {
	roles, err := r.q.ListUserRoles(ctx, u.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return domainuser.User{}, err
	}
	perms, err := r.q.ListUserPermissions(ctx, u.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return domainuser.User{}, err
	}

	return domainuser.User{
		ID:             u.ID.String(),
		Email:          u.Email,
		HashedPassword: u.HashedPassword,
		Name:           u.Name,
		Roles:          roles,
		Permissions:    perms,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		DeletedAt:      u.DeletedAt,
	}, nil
}
