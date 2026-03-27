package models

import "time"

type ProductView struct {
	ViewID    int
	UserID    int
	ProductID int
	ViewedAt  time.Time
}
