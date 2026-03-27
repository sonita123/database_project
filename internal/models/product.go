package models

import "time"

type Product struct {
	ProductID int
	StallID   int
	Name      string
	Price     float64
	Stock     int
	Status    string
	CreatedAt time.Time
}
