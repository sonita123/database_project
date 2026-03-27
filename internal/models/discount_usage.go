package models

import "time"

type DiscountUsage struct {
	UsageID    int
	DiscountID int
	UserID     int
	UsedAt     time.Time
}
