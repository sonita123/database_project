package models

type OrderItem struct {
	OrderItemID int
	OrderID     int
	ProductID   int
	Quantity    int
	Price       float64
	// Joined field for display
	ProductName string
}
