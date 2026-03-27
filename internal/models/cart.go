package models

import "time"

type Cart struct {
	CartID    int
	UserID    int
	Locked    bool
	CreatedAt time.Time
}
