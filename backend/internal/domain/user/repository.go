package user

import "context"

// Repository defines persistence operations for users.
type Repository interface {
	Create(ctx context.Context, u User) (User, error)
	GetByID(ctx context.Context, id string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, u User) error
	SoftDelete(ctx context.Context, id string) error
}
