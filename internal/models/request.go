package models

import "time"

type Request struct {
	RequestID     int
	UserID        int
	RequestType   string
	Status        string
	Description   *string
	HandledBy     *int
	ResolvedAt    *time.Time
	CreatedAt     time.Time
	UserFirstName string
	UserLastName  string
	UserEmail     string
}
