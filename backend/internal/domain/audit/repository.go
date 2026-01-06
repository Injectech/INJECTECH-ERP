package audit

import "context"

// Repository defines persistence operations for audit logs.
type Repository interface {
	Create(ctx context.Context, log Log) error
	List(ctx context.Context, actorID string) ([]Log, error)
}
