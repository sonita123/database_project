package models

import "time"

type Session struct {
	Token     string
	UserID    int
	Role      string
	ExpiresAt time.Time
}
