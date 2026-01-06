package inventory

import "time"

// Inventory captures stock information for products.
type Inventory struct {
	ID        string
	ProductID string
	Quantity  int64
	Location  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
