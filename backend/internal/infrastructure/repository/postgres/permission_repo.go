package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domainpermission "backend/internal/domain/permission"
	"backend/internal/infrastructure/repository/postgres/sqlc"
)

// PermissionRepository implements permission.Repository via PostgreSQL.
type PermissionRepository struct {
	db *pgxpool.Pool
	q  *sqlc.Queries
}

func NewPermissionRepository(db *pgxpool.Pool) *PermissionRepository {
	return &PermissionRepository{db: db, q: sqlc.New(db)}
}

func (r *PermissionRepository) Create(ctx context.Context, p domainpermission.Permission) (domainpermission.Permission, error) {
	pid, err := uuid.Parse(p.ID)
	if err != nil {
		return domainpermission.Permission{}, err
	}
	created, err := r.q.CreatePermission(ctx, sqlc.CreatePermissionParams{
		ID:          pid,
		Code:        p.Code,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	})
	if err != nil {
		return domainpermission.Permission{}, err
	}
	p.CreatedAt = created.CreatedAt
	p.UpdatedAt = created.UpdatedAt
	return p, nil
}

func (r *PermissionRepository) GetByCode(ctx context.Context, code string) (domainpermission.Permission, error) {
	p, err := r.q.GetPermissionByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainpermission.Permission{}, err
		}
		return domainpermission.Permission{}, err
	}
	return domainpermission.Permission{
		ID:          p.ID.String(),
		Code:        p.Code,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		DeletedAt:   p.DeletedAt,
	}, nil
}

func (r *PermissionRepository) List(ctx context.Context) ([]domainpermission.Permission, error) {
	dbPerms, err := r.q.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}

	var result []domainpermission.Permission
	for _, p := range dbPerms {
		result = append(result, domainpermission.Permission{
			ID:          p.ID.String(),
			Code:        p.Code,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
			DeletedAt:   p.DeletedAt,
		})
	}
	return result, nil
}
