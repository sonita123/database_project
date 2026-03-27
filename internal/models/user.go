package models

import "time"

type User struct {
	UserID    int
	FirstName string
	LastName  string
	Email     string
	Password  string
	Balance   float64
	UserType  string
	CreatedAt time.Time
}
