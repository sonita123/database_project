package models

import "time"

type Stall struct {
	StallID        int
	SellerID       int
	Name           string
	Status         string
	ApprovedBy     *int
	ApprovedAt     *time.Time
	CreatedAt      time.Time
	OwnerFirstName string
	OwnerLastName  string
}
