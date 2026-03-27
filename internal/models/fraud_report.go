package models

import "time"

type FraudReport struct {
	ReportID    int
	StallID     int
	ReporterID  int
	Description *string
	ReportedAt  time.Time
}
