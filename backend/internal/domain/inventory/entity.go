package inventory

import "time"

// Inventory captures stock information for products.
type Inventory struct {
	ID        string     `json:"id"`
	ProductID string     `json:"product_id"`
	Quantity  int64      `json:"quantity"`
	Location  string     `json:"location"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
