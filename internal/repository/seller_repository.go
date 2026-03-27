package repository

import (
	"context"
	"database/sql"
	"unibazar/project/internal/db"
	"unibazar/project/internal/models"
)

// GetSellerByEmail for login
func GetSellerByEmail(email string) (models.Seller, error) {
	var s struct {
		UserID    int
		FirstName string
		LastName  string
		Email     string
		SellerID  int
	}
	_ = db.Conn.QueryRowContext(context.Background(), `
		SELECT u.user_id, u.first_name, u.last_name, u.email, s.seller_id
		FROM users u
		JOIN sellers s ON s.user_id = u.user_id
		WHERE u.email = @email
	`, sql.Named("email", email)).Scan(
		&s.UserID, &s.FirstName, &s.LastName, &s.Email, &s.SellerID,
	)
	if s.SellerID == 0 {
		return models.Seller{}, sql.ErrNoRows
	}
	return models.Seller{
		SellerID:  s.SellerID,
		UserID:    s.UserID,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Email:     s.Email,
	}, nil
}

// CreateSellerIfNotExists creates seller if not exists for user
func CreateSellerIfNotExists(userID int) (int, error) {
	var sellerID int
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT seller_id FROM sellers WHERE user_id = @user_id
	`, sql.Named("user_id", userID)).Scan(&sellerID)
	if err == nil {
		return sellerID, nil // exists
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	// create
	err = db.Conn.QueryRowContext(context.Background(), `
		INSERT INTO sellers (user_id)
		OUTPUT INSERTED.seller_id
		VALUES (@user_id)
	`, sql.Named("user_id", userID)).Scan(&sellerID)
	return sellerID, err
}

// GetSellerStalls
func GetSellerStalls(sellerID int) ([]models.Stall, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT stall_id, seller_id, name, status, approved_by, approved_at, created_at
		FROM stalls
		WHERE seller_id = @seller_id
		ORDER BY created_at DESC
	`, sql.Named("seller_id", sellerID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stalls []models.Stall
	for rows.Next() {
		var stall models.Stall
		err := rows.Scan(&stall.StallID, &stall.SellerID, &stall.Name, &stall.Status, &stall.ApprovedBy, &stall.ApprovedAt, &stall.CreatedAt)
		if err != nil {
			return nil, err
		}
		stalls = append(stalls, stall)
	}
	return stalls, rows.Err()
}

// RequestNewStall for seller
func RequestNewStall(sellerID int, name string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		INSERT INTO stalls (seller_id, name, status) 
		VALUES (@seller_id, @name, 'pending')
	`, sql.Named("seller_id", sellerID), sql.Named("name", name))
	return err
}

// GetSellerDashboardData - main dashboard query
func GetSellerDashboardData(sellerID int) (models.SellerDashboardData, error) {
	var data models.SellerDashboardData

	// Basic seller info + stats
	err := db.Conn.QueryRowContext(context.Background(), `
		WITH seller_stalls AS (
			SELECT stall_id FROM stalls WHERE seller_id = @seller_id AND status = 'active'
		),
		sales_stats AS (
			SELECT 
				ISNULL(SUM(oi.quantity * oi.price), 0) as total_sales,
				ISNULL(SUM(oi.quantity * oi.price * 0.1), 0) as estimated_profit, -- assuming 10% commission
				COUNT(DISTINCT oi.order_id) as total_orders
			FROM order_items oi
			JOIN orders o ON o.order_id = oi.order_id
			JOIN products p ON p.product_id = oi.product_id
			WHERE p.stall_id IN (SELECT stall_id FROM seller_stalls)
		),
		low_stock AS (
			SELECT COUNT(*) as low_stock_count
			FROM products p
			JOIN seller_stalls ss ON p.stall_id = ss.stall_id
			WHERE p.stock <= 5 AND p.status = 'active'
		)
		SELECT 
			(SELECT seller_id FROM sellers WHERE seller_id = @seller_id) as seller_id,
			ss.total_sales, ss.estimated_profit, ss.total_orders,
			ls.low_stock_count
		FROM sales_stats ss, low_stock ls
	`, sql.Named("seller_id", sellerID)).Scan(
		&data.SellerID,
		&data.TotalSales, &data.EstimatedProfit, &data.TotalOrders,
		&data.LowStockCount,
	)
	if err != nil {
		return data, err
	}

	// Low stock products
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT p.product_id, p.name, p.price, p.stock, p.stall_id
		FROM products p
		JOIN stalls s ON s.stall_id = p.stall_id
		WHERE s.seller_id = @seller_id AND p.stock <= 5 AND p.status = 'active'
		ORDER BY p.stock ASC
	`, sql.Named("seller_id", sellerID))
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ProductID, &p.Name, &p.Price, &p.Stock, &p.StallID)
		if err != nil {
			return data, err
		}
		data.LowStockProducts = append(data.LowStockProducts, p)
	}

	data.Stalls, _ = GetSellerStalls(sellerID) // ignore error for now

	return data, rows.Err()
}

// ListSellerProductsByStall
func ListSellerProductsByStall(stallID int) ([]models.Product, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT product_id, stall_id, name, price, stock, status, created_at
		FROM products
		WHERE stall_id = @stall_id
		ORDER BY created_at DESC
	`, sql.Named("stall_id", stallID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ProductID, &p.StallID, &p.Name, &p.Price, &p.Stock, &p.Status, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}
