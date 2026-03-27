package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"unibazar/project/internal/db"
	"unibazar/project/internal/models"
)

// ── Products ─────────────────────────────────────────────────

func GetAllStalls() ([]models.Stall, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT stall_id, seller_id, name, status, approved_by, approved_at, created_at
		FROM stalls
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stalls []models.Stall
	for rows.Next() {
		var s models.Stall
		if err := rows.Scan(
			&s.StallID,
			&s.SellerID,
			&s.Name,
			&s.Status,
			&s.ApprovedBy,
			&s.ApprovedAt,
			&s.CreatedAt,
		); err != nil {
			return nil, err
		}
		stalls = append(stalls, s)
	}

	return stalls, rows.Err()
}
func GetProductsPaginated(limit, offset int) ([]models.Product, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT product_id, stall_id, name, price, stock, status, created_at
		FROM products
		ORDER BY product_id DESC 
		OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY
	`,
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
		err := rows.Scan(&p.ProductID, &p.StallID, &p.Name, &p.Price, &p.Stock, &p.Status, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func CountProducts() (int, error) {
	var count int
	err := db.Conn.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM products`,
	).Scan(&count)
	return count, err
}

func GetProductByID(id string) (models.Product, error) {
	var p models.Product
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT product_id, stall_id, name, price, stock, status, created_at
		FROM products WHERE product_id = @id
	`, sql.Named("id", id)).Scan(
		&p.ProductID, &p.StallID, &p.Name, &p.Price, &p.Stock, &p.Status, &p.CreatedAt,
	)
	return p, err
}

func CreateProduct(stallID, name, price, stock, status string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		INSERT INTO products (stall_id, name, price, stock, status) 
		VALUES (@stall_id, @name, @price, @stock, @status)
	`,
		sql.Named("stall_id", stallID),
		sql.Named("name", name),
		sql.Named("price", price),
		sql.Named("stock", stock),
		sql.Named("status", status),
	)
	return err
}

func UpdateProduct(id, name, price, stock, status string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		UPDATE products 
		SET name = @name, price = @price, stock = @stock, status = @status 
		WHERE product_id = @id
	`,
		sql.Named("name", name),
		sql.Named("price", price),
		sql.Named("stock", stock),
		sql.Named("status", status),
		sql.Named("id", id),
	)
	return err
}

func DeleteProduct(id string) error {
	_, err := db.Conn.ExecContext(context.Background(),
		`DELETE FROM products WHERE product_id = @id`,
		sql.Named("id", id),
	)
	return err
}

// ── Stalls ───────────────────────────────────────────────────

func GetStallsPaginated(limit, offset int) ([]models.Stall, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT stall_id, seller_id, name, status, approved_by, approved_at, created_at
		FROM stalls
		ORDER BY stall_id DESC 
		OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY
	`,
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stalls []models.Stall
	for rows.Next() {
		var s models.Stall
		err := rows.Scan(&s.StallID, &s.SellerID, &s.Name, &s.Status, &s.ApprovedBy, &s.ApprovedAt, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		stalls = append(stalls, s)
	}
	return stalls, rows.Err()
}

func CountStalls() (int, error) {
	var count int
	err := db.Conn.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM stalls`,
	).Scan(&count)
	return count, err
}

func GetStallByID(id string) (models.Stall, error) {
	var s models.Stall
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT stall_id, seller_id, name, status, approved_by, approved_at, created_at
		FROM stalls WHERE stall_id = @id
	`, sql.Named("id", id)).Scan(
		&s.StallID, &s.SellerID, &s.Name, &s.Status, &s.ApprovedBy, &s.ApprovedAt, &s.CreatedAt,
	)
	return s, err
}

func UpdateStallStatus(id, status string) error {
	_, err := db.Conn.ExecContext(context.Background(),
		`UPDATE stalls SET status = @status WHERE stall_id = @id`,
		sql.Named("status", status),
		sql.Named("id", id),
	)
	return err
}

func DeleteStall(id string) error {
	_, err := db.Conn.ExecContext(context.Background(),
		`DELETE FROM stalls WHERE stall_id = @id`,
		sql.Named("id", id),
	)
	return err
}

// ── Orders ───────────────────────────────────────────────────

func GetOrdersPaginated(limit, offset int) ([]models.Order, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT order_id, user_id, cart_id, address_id, total_price, status, order_date
		FROM orders
		ORDER BY order_date DESC 
		OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY
	`,
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		err := rows.Scan(&o.OrderID, &o.UserID, &o.CartID, &o.AddressID, &o.TotalPrice, &o.Status, &o.OrderDate)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func CountOrders() (int, error) {
	var count int
	err := db.Conn.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM orders`,
	).Scan(&count)
	return count, err
}

// ── Discount Codes ───────────────────────────────────────────

func GetDiscountCodesPaginated(limit, offset int) ([]models.DiscountCode, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT discount_id, code, discount_type, percentage, fixed_amount,
		       expiration_date, max_uses, is_active, supporter_id
		FROM discount_codes
		ORDER BY discount_id DESC 
		OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY
	`,
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []models.DiscountCode
	for rows.Next() {
		var c models.DiscountCode
		err := rows.Scan(&c.DiscountID, &c.Code, &c.DiscountType, &c.Percentage, &c.FixedAmount, &c.ExpirationDate, &c.MaxUses, &c.IsActive, &c.SupporterID)
		if err != nil {
			return nil, err
		}
		codes = append(codes, c)
	}
	return codes, rows.Err()
}

func CountDiscountCodes() (int, error) {
	var count int
	err := db.Conn.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM discount_codes`,
	).Scan(&count)
	return count, err
}

func GetDiscountCodeByID(id string) (models.DiscountCode, error) {
	var c models.DiscountCode
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT discount_id, code, discount_type, percentage, fixed_amount,
		       expiration_date, max_uses, is_active, supporter_id
		FROM discount_codes WHERE discount_id = @id
	`, sql.Named("id", id)).Scan(
		&c.DiscountID, &c.Code, &c.DiscountType, &c.Percentage, &c.FixedAmount, &c.ExpirationDate, &c.MaxUses, &c.IsActive, &c.SupporterID,
	)
	return c, err
}

func CreateDiscountCode(code, discountType, percentage, fixedAmount, expirationDate, maxUses, supporterID string) error {
	// SQL Server requires a proper time.Time value — parse the date string first
	expDate, err := time.Parse("2006-01-02", expirationDate)
	if err != nil {
		// Try common alternative format
		expDate, err = time.Parse("2006-01-02T15:04", expirationDate)
		if err != nil {
			return fmt.Errorf("invalid expiration date format: %w", err)
		}
	}

	_, err = db.Conn.ExecContext(context.Background(), `
		INSERT INTO discount_codes
		       (code, discount_type, percentage, fixed_amount, expiration_date, max_uses, supporter_id)
		VALUES (@code, @discount_type,
		        NULLIF(@percentage, ''),
		        NULLIF(@fixed_amount, ''),
		        @expiration_date,
		        NULLIF(@max_uses, ''),
		        NULLIF(@supporter_id, ''))
	`,
		sql.Named("code", code),
		sql.Named("discount_type", discountType),
		sql.Named("percentage", percentage),
		sql.Named("fixed_amount", fixedAmount),
		sql.Named("expiration_date", expDate),
		sql.Named("max_uses", maxUses),
		sql.Named("supporter_id", supporterID),
	)
	return err
}
func ToggleDiscountCode(id string, isActive bool) error {
	_, err := db.Conn.ExecContext(context.Background(),
		`UPDATE discount_codes SET is_active = @active WHERE discount_id = @id`,
		sql.Named("active", isActive),
		sql.Named("id", id),
	)
	return err
}

func DeleteDiscountCode(id string) error {
	_, err := db.Conn.ExecContext(context.Background(),
		`DELETE FROM discount_codes WHERE discount_id = @id`,
		sql.Named("id", id),
	)
	return err
}

// ── Fraud Reports ────────────────────────────────────────────

func GetFraudReportsPaginated(limit, offset int) ([]models.FraudReport, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT report_id, stall_id, reporter_id, description, reported_at
		FROM fraud_reports
		ORDER BY reported_at DESC 
		OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY
	`,
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.FraudReport
	for rows.Next() {
		var f models.FraudReport
		err := rows.Scan(&f.ReportID, &f.StallID, &f.ReporterID, &f.Description, &f.ReportedAt)
		if err != nil {
			return nil, err
		}
		reports = append(reports, f)
	}
	return reports, rows.Err()
}

func CountFraudReports() (int, error) {
	var count int
	err := db.Conn.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM fraud_reports`,
	).Scan(&count)
	return count, err
}

func DeleteFraudReport(id string) error {
	_, err := db.Conn.ExecContext(context.Background(),
		`DELETE FROM fraud_reports WHERE report_id = @id`,
		sql.Named("id", id),
	)
	return err
}

// ── Reviews ──────────────────────────────────────────────────

func GetReviewsPaginated(limit, offset int) ([]models.Review, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT review_id, user_id, product_id, order_item_id, rating, comment, created_at
		FROM reviews
		ORDER BY created_at DESC 
		OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY
	`,
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var rv models.Review
		err := rows.Scan(&rv.ReviewID, &rv.UserID, &rv.ProductID, &rv.OrderItemID, &rv.Rating, &rv.Comment, &rv.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, rv)
	}
	return reviews, rows.Err()
}

func CountReviews() (int, error) {
	var count int
	err := db.Conn.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM reviews`,
	).Scan(&count)
	return count, err
}
