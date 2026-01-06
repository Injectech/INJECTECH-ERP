package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	domainrole "backend/internal/domain/role"
)

// RoleRepository implements role.Repository via PostgreSQL.
type RoleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(ctx context.Context, role domainrole.Role) (domainrole.Role, error) {
	return role, nil
}

func (r *RoleRepository) GetByID(ctx context.Context, id string) (domainrole.Role, error) {
	return domainrole.Role{}, nil
}

func (r *RoleRepository) GetByName(ctx context.Context, name string) (domainrole.Role, error) {
	return domainrole.Role{}, nil
}

func (r *RoleRepository) List(ctx context.Context) ([]domainrole.Role, error) {
	return []domainrole.Role{}, nil
}

func (r *RoleRepository) Update(ctx context.Context, role domainrole.Role) error {
	return nil
}

func (r *RoleRepository) SoftDelete(ctx context.Context, id string) error {
	return nil
}
