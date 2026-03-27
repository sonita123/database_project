package models

import "time"

// Seller links user to seller activities
type Seller struct {
	SellerID     int       `db:"seller_id"`
	UserID       int       `db:"user_id"`
	RegisteredAt time.Time `db:"registered_at"`
	FirstName    string    `db:"first_name"` // denormalized from users
	LastName     string    `db:"last_name"`
	Email        string    `db:"email"`
}

// SellerDashboardData for portal
type SellerDashboardData struct {
	SellerID         int
	TotalSales       float64   `json:"total_sales"`
	EstimatedProfit  float64   `json:"estimated_profit"`
	TotalOrders      int       `json:"total_orders"`
	LowStockCount    int       `json:"low_stock_count"`
	Stalls           []Stall   `json:"stalls"`
	LowStockProducts []Product `json:"low_stock_products"`
}
