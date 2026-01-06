package role

import "context"

// Repository defines persistence operations for roles.
type Repository interface {
	Create(ctx context.Context, r Role) (Role, error)
	GetByID(ctx context.Context, id string) (Role, error)
	GetByName(ctx context.Context, name string) (Role, error)
	List(ctx context.Context) ([]Role, error)
	Update(ctx context.Context, r Role) error
	SoftDelete(ctx context.Context, id string) error
}
