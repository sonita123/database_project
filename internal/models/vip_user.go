package models

import "time"

type VipUser struct {
	VipID     int
	UserID    int
	StartDate time.Time
	EndDate   time.Time
}
