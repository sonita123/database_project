package models

type Address struct {
	AddressID  int
	UserID     int
	City       string
	Street     string
	PostalCode string
	IsDefault  bool
}
