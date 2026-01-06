package audit

import "time"

// Log represents an auditable action.
type Log struct {
	ID        string
	ActorID   string
	Action    string
	Resource  string
	Metadata  map[string]any
	CreatedAt time.Time
}
