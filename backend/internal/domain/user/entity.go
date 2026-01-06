package user

import "time"

// User represents an account in the system.
type User struct {
	ID             string
	Email          string
	HashedPassword string
	Name           string
	Roles          []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}
