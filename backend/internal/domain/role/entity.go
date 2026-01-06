package role

import "time"

// Role defines a set of permissions.
type Role struct {
	ID          string
	Name        string
	Description string
	Permissions []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
