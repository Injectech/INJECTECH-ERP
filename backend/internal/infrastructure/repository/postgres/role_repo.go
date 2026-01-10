package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domainrole "backend/internal/domain/role"
	"backend/internal/infrastructure/repository/postgres/sqlc"
)

// RoleRepository implements role.Repository via PostgreSQL.
type RoleRepository struct {
	db *pgxpool.Pool
	q  *sqlc.Queries
}

func NewRoleRepository(db *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{db: db, q: sqlc.New(db)}
}

func (r *RoleRepository) Create(ctx context.Context, role domainrole.Role) (domainrole.Role, error) {
	rid, err := uuid.Parse(role.ID)
	if err != nil {
		return domainrole.Role{}, err
	}
	created, err := r.q.CreateRole(ctx, sqlc.CreateRoleParams{
		ID:          rid,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	})
	if err != nil {
		return domainrole.Role{}, err
	}

	for _, code := range role.Permissions {
		if _, err := r.db.Exec(ctx, `
			INSERT INTO role_permissions (role_id, permission_id)
			SELECT $1, p.id FROM permissions p WHERE p.code = $2
			ON CONFLICT (role_id, permission_id) DO NOTHING
		`, rid, code); err != nil {
			return domainrole.Role{}, err
		}
	}

	role.CreatedAt = created.CreatedAt
	role.UpdatedAt = created.UpdatedAt
	return role, nil
}

func (r *RoleRepository) GetByID(ctx context.Context, id string) (domainrole.Role, error) {
	rid, err := uuid.Parse(id)
	if err != nil {
		return domainrole.Role{}, err
	}
	rl, err := r.q.GetRoleByID(ctx, rid)
	if err != nil {
		return domainrole.Role{}, err
	}
	return r.withPermissions(ctx, rl)
}

func (r *RoleRepository) GetByName(ctx context.Context, name string) (domainrole.Role, error) {
	rl, err := r.q.GetRoleByName(ctx, name)
	if err != nil {
		return domainrole.Role{}, err
	}
	return r.withPermissions(ctx, rl)
}

func (r *RoleRepository) List(ctx context.Context) ([]domainrole.Role, error) {
	dbRoles, err := r.q.ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	var roles []domainrole.Role
	for _, rl := range dbRoles {
		mapped, err := r.withPermissions(ctx, rl)
		if err != nil {
			return nil, err
		}
		roles = append(roles, mapped)
	}
	return roles, nil
}

func (r *RoleRepository) Update(ctx context.Context, role domainrole.Role) error {
	rid, err := uuid.Parse(role.ID)
	if err != nil {
		return err
	}
	return r.q.UpdateRole(ctx, sqlc.UpdateRoleParams{
		ID:          rid,
		Name:        role.Name,
		Description: role.Description,
		UpdatedAt:   time.Now(),
	})
}

func (r *RoleRepository) SoftDelete(ctx context.Context, id string) error {
	rid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.SoftDeleteRole(ctx, rid)
}

func (r *RoleRepository) withPermissions(ctx context.Context, rl sqlc.Role) (domainrole.Role, error) {
	perms, err := r.q.ListRolePermissions(ctx, rl.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return domainrole.Role{}, err
	}
	return domainrole.Role{
		ID:          rl.ID.String(),
		Name:        rl.Name,
		Description: rl.Description,
		Permissions: perms,
		CreatedAt:   rl.CreatedAt,
		UpdatedAt:   rl.UpdatedAt,
		DeletedAt:   rl.DeletedAt,
	}, nil
}
