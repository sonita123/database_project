package models

type CartItem struct {
	CartItemID   int
	CartID       int
	ProductID    int
	Quantity     int
	ProductName  string
	ProductPrice float64
}
