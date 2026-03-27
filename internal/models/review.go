package models

import "time"

type Review struct {
	ReviewID          int
	UserID            int
	ProductID         int
	OrderItemID       *int
	Rating            int
	Comment           *string
	CreatedAt         time.Time
	ProductName       string
	ReviewerFirstName string
	ReviewerLastName  string
}
type ReviewableItem struct {
	OrderItemID  int
	ProductID    int
	ProductName  string
	ProductPrice float64
	OrderID      int
	OrderDate    time.Time
}
