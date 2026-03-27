package repository

import (
	"context"
	"database/sql"
	"unibazar/project/internal/db"
)

/* Total Users */

func CountUsers() (int, error) {
	var count int
	err := db.Conn.QueryRowContext(
		context.Background(),
		"SELECT COUNT(*) FROM users",
	).Scan(&count)
	return count, err
}

/* Total Revenue */

func CalculateTotalRevenue() (float64, error) {
	var revenue float64
	err := db.Conn.QueryRowContext(
		context.Background(),
		"SELECT ISNULL(SUM(total_price), 0) FROM orders",
	).Scan(&revenue)
	return revenue, err
}

/* Open Requests */

func CountOpenRequests() (int, error) {
	var count int
	err := db.Conn.QueryRowContext(
		context.Background(),
		"SELECT COUNT(*) FROM requests WHERE status = @status",
		sql.Named("status", "open"),
	).Scan(&count)
	return count, err
}
