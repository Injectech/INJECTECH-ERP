package product

import "time"

// Product represents catalog data.
type Product struct {
	ID          string
	SKU         string
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
