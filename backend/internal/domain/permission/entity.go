package permission

import "time"

// Permission represents a granular capability.
type Permission struct {
	ID          string
	Code        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
