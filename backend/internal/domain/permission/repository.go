package permission

import "context"

// Repository defines persistence operations for permissions.
type Repository interface {
    Create(ctx context.Context, p Permission) (Permission, error)
    GetByCode(ctx context.Context, code string) (Permission, error)
    List(ctx context.Context) ([]Permission, error)
}
