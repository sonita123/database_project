package models

import "time"

type Order struct {
	OrderID    int
	UserID     int
	CartID     *int
	AddressID  *int
	TotalPrice float64
	Status     string
	OrderDate  time.Time
	// Joined user fields
	UserFirstName string
	UserLastName  string
	UserEmail     string
	// Joined address fields
	DeliveryStreet     string
	DeliveryCity       string
	DeliveryPostalCode string
}
