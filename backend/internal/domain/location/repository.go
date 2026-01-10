package location

import (
	"context"
	"errors"
)

// ErrNotFound is returned when a location is not found.
var ErrNotFound = errors.New("location not found")

// Repository defines persistence operations for locations.
type Repository interface {
	Create(ctx context.Context, l Location) (Location, error)
	List(ctx context.Context) ([]Location, error)
	GetDefault(ctx context.Context) (Location, error)
}
