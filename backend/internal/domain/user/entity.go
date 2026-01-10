package user

import "time"

// User represents an account in the system.
type User struct {
	ID             string     `json:"id"`
	Email          string     `json:"email"`
	HashedPassword string     `json:"-"`
	Name           string     `json:"name"`
	Roles          []string   `json:"roles"`
	Permissions    []string   `json:"permissions"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
