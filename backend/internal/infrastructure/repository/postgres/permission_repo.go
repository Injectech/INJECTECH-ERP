package postgres

import (
    "context"

    "github.com/jackc/pgx/v5/pgxpool"

    domainpermission "backend/internal/domain/permission"
)

// PermissionRepository implements permission.Repository via PostgreSQL.
type PermissionRepository struct {
    db *pgxpool.Pool
}

func NewPermissionRepository(db *pgxpool.Pool) *PermissionRepository {
    return &PermissionRepository{db: db}
}

func (r *PermissionRepository) Create(ctx context.Context, p domainpermission.Permission) (domainpermission.Permission, error) {
    return p, nil
}

func (r *PermissionRepository) GetByCode(ctx context.Context, code string) (domainpermission.Permission, error) {
    return domainpermission.Permission{}, nil
}

func (r *PermissionRepository) List(ctx context.Context) ([]domainpermission.Permission, error) {
    return []domainpermission.Permission{}, nil
}
