package repository

import (
	"context"
	"database/sql"
	"unibazar/project/internal/db"
	"unibazar/project/internal/models"
)

// GetProductsByStall returns all products in a given stall
func GetProductsByStall(stallID int) ([]models.Product, error) {
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
		if err := rows.Scan(
			&p.ProductID, &p.StallID, &p.Name, &p.Price,
			&p.Stock, &p.Status, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// GetProductByIDInt fetches a single product by integer ID
func GetProductByIDInt(productID int) (models.Product, error) {
	var p models.Product
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT product_id, stall_id, name, price, stock, status, created_at
		FROM products
		WHERE product_id = @product_id
	`, sql.Named("product_id", productID)).Scan(
		&p.ProductID, &p.StallID, &p.Name, &p.Price,
		&p.Stock, &p.Status, &p.CreatedAt,
	)
	return p, err
}

// CreateProductInStall adds a new product to a stall
func CreateProductInStall(stallID int, name string, price float64, stock int) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		INSERT INTO products (stall_id, name, price, stock, status)
		VALUES (@stall_id, @name, @price, @stock, 'active')
	`,
		sql.Named("stall_id", stallID),
		sql.Named("name", name),
		sql.Named("price", price),
		sql.Named("stock", stock),
	)
	return err
}

// UpdateProductInStall updates an existing product
func UpdateProductInStall(productID int, name string, price float64, stock int, status string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		UPDATE products
		SET name   = @name,
		    price  = @price,
		    stock  = @stock,
		    status = @status
		WHERE product_id = @product_id
	`,
		sql.Named("name", name),
		sql.Named("price", price),
		sql.Named("stock", stock),
		sql.Named("status", status),
		sql.Named("product_id", productID),
	)
	return err
}

// DeleteProductFromStall removes a product — only if it belongs to the stall
func DeleteProductFromStall(productID, stallID int) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		DELETE FROM products
		WHERE product_id = @product_id
		  AND stall_id   = @stall_id
	`,
		sql.Named("product_id", productID),
		sql.Named("stall_id", stallID),
	)
	return err
}

// GetProductsBySellerUserID returns all active products across all stalls
// belonging to the given user (used to show a seller their own shop view)
func GetProductsBySellerUserID(userID, limit, offset int) ([]models.Product, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT p.product_id, p.stall_id, p.name, p.price, p.stock, p.status, p.created_at
		FROM products p
		JOIN stalls s   ON s.stall_id   = p.stall_id
		JOIN sellers sl ON sl.seller_id = s.seller_id
		WHERE sl.user_id = @user_id
		  AND p.status   = 'active'
		  AND p.stock    > 0
		  AND s.status   = 'active'
		ORDER BY p.created_at DESC
		OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY
	`,
		sql.Named("user_id", userID),
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ProductID, &p.StallID, &p.Name, &p.Price,
			&p.Stock, &p.Status, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// CountProductsBySellerUserID counts products for pagination
func CountProductsBySellerUserID(userID int) (int, error) {
	var count int
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT COUNT(*)
		FROM products p
		JOIN stalls s   ON s.stall_id   = p.stall_id
		JOIN sellers sl ON sl.seller_id = s.seller_id
		WHERE sl.user_id = @user_id
		  AND p.status   = 'active'
		  AND p.stock    > 0
		  AND s.status   = 'active'
	`, sql.Named("user_id", userID)).Scan(&count)
	return count, err
}
