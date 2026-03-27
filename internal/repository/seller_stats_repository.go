package repository

import (
	"context"
	"database/sql"
	"unibazar/project/internal/db"
	"unibazar/project/internal/models"
)

// GetSellerStats returns aggregated sales stats for a seller
func GetSellerStats(userID int) (models.SellerStats, error) {
	var s models.SellerStats

	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT
		    ISNULL(SUM(oi.quantity * oi.price), 0)  AS total_revenue,
		    COUNT(DISTINCT o.order_id)               AS total_orders,
		    ISNULL(SUM(oi.quantity), 0)              AS total_items_sold
		FROM order_items oi
		JOIN orders o   ON o.order_id     = oi.order_id
		JOIN products p ON p.product_id   = oi.product_id
		JOIN stalls s   ON s.stall_id     = p.stall_id
		JOIN sellers sl ON sl.seller_id   = s.seller_id
		WHERE sl.user_id  = @user_id
		  AND o.status   != 'cancelled'
	`, sql.Named("user_id", userID)).Scan(
		&s.TotalRevenue, &s.TotalOrders, &s.TotalItemsSold,
	)
	if err != nil {
		return s, err
	}

	// Product counts
	db.Conn.QueryRowContext(context.Background(), `
		SELECT
		    COUNT(*)                                          AS total_products,
		    SUM(CASE WHEN p.status = 'active' THEN 1 ELSE 0 END) AS active_products
		FROM products p
		JOIN stalls s   ON s.stall_id   = p.stall_id
		JOIN sellers sl ON sl.seller_id = s.seller_id
		WHERE sl.user_id = @user_id
	`, sql.Named("user_id", userID)).Scan(
		&s.TotalProducts, &s.ActiveProducts,
	)

	// Stall counts
	db.Conn.QueryRowContext(context.Background(), `
		SELECT
		    COUNT(*)                                            AS total_stalls,
		    SUM(CASE WHEN s.status = 'active' THEN 1 ELSE 0 END) AS active_stalls
		FROM stalls s
		JOIN sellers sl ON sl.seller_id = s.seller_id
		WHERE sl.user_id = @user_id
	`, sql.Named("user_id", userID)).Scan(
		&s.TotalStalls, &s.ActiveStalls,
	)

	return s, nil
}

// GetSellerTopProducts returns top 5 best-selling products for a seller
func GetSellerTopProducts(userID int) ([]models.SellerTopProduct, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT TOP 5
		    p.product_id,
		    p.name                             AS product_name,
		    st.name                            AS stall_name,
		    ISNULL(SUM(oi.quantity), 0)        AS total_sold,
		    ISNULL(SUM(oi.quantity * oi.price), 0) AS revenue
		FROM products p
		JOIN stalls st  ON st.stall_id    = p.stall_id
		JOIN sellers sl ON sl.seller_id   = st.seller_id
		LEFT JOIN order_items oi ON oi.product_id = p.product_id
		LEFT JOIN orders o       ON o.order_id    = oi.order_id
		                        AND o.status     != 'cancelled'
		WHERE sl.user_id = @user_id
		GROUP BY p.product_id, p.name, st.name
		ORDER BY total_sold DESC
	`, sql.Named("user_id", userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.SellerTopProduct
	for rows.Next() {
		var p models.SellerTopProduct
		if err := rows.Scan(
			&p.ProductID, &p.ProductName, &p.StallName,
			&p.TotalSold, &p.Revenue,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// GetSellerRecentOrders returns the 10 most recent orders containing seller products
func GetSellerRecentOrders(userID int) ([]models.Order, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT DISTINCT TOP 10
		    o.order_id, o.user_id, o.cart_id, o.address_id,
		    o.total_price, o.status, o.order_date,
		    u.first_name, u.last_name, u.email,
		    ISNULL(a.street, ''), ISNULL(a.city, ''), ISNULL(a.postal_code, '')
		FROM orders o
		JOIN order_items oi ON oi.order_id   = o.order_id
		JOIN products p     ON p.product_id  = oi.product_id
		JOIN stalls s       ON s.stall_id    = p.stall_id
		JOIN sellers sl     ON sl.seller_id  = s.seller_id
		JOIN users u        ON u.user_id     = o.user_id
		LEFT JOIN addresses a ON a.address_id = o.address_id
		WHERE sl.user_id = @user_id
		ORDER BY o.order_date DESC
	`, sql.Named("user_id", userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.OrderID, &o.UserID, &o.CartID, &o.AddressID,
			&o.TotalPrice, &o.Status, &o.OrderDate,
			&o.UserFirstName, &o.UserLastName, &o.UserEmail,
			&o.DeliveryStreet, &o.DeliveryCity, &o.DeliveryPostalCode,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}
