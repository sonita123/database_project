package models

import "time"

type DiscountCode struct {
	DiscountID     int
	Code           string
	DiscountType   string
	Percentage     *float64
	FixedAmount    *float64
	ExpirationDate time.Time
	MaxUses        *int
	IsActive       bool
	SupporterID    *int
}
