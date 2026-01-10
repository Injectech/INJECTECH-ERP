package permission

import "time"

// Permission represents a granular capability.
type Permission struct {
	ID          string     `json:"id"`
	Code        string     `json:"code"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
