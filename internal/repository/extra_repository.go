package repository

import (
	"context"
	"database/sql"
	"errors"
	"unibazar/project/internal/db"
	"unibazar/project/internal/models"
)

// ── Discount codes for user portal ───────────────────────────

// GetUserDiscountCodes returns all active discount codes with usage status for user
func GetUserDiscountCodes(userID string) ([]models.DiscountCodeUsage, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT dc.discount_id, dc.code, dc.discount_type,
		       dc.percentage, dc.fixed_amount, dc.expiration_date,
		       dc.is_active,
		       CASE WHEN du.usage_id IS NOT NULL THEN 1 ELSE 0 END AS used_by_user
		FROM discount_codes dc
		LEFT JOIN discount_usage du
		       ON du.discount_id = dc.discount_id
		      AND du.user_id = @user_id
		WHERE dc.is_active = 1
		  AND dc.expiration_date > GETDATE()
		ORDER BY dc.expiration_date ASC
	`, sql.Named("user_id", userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []models.DiscountCodeUsage
	for rows.Next() {
		var c models.DiscountCodeUsage
		if err := rows.Scan(
			&c.DiscountID, &c.Code, &c.DiscountType,
			&c.Percentage, &c.FixedAmount, &c.ExpirationDate,
			&c.IsActive, &c.UsedByUser,
		); err != nil {
			return nil, err
		}
		codes = append(codes, c)
	}
	return codes, rows.Err()
}

// ValidateDiscountCode checks the code exists, is active, not expired,
// user hasn't used it, and max_uses not exceeded.
// Returns the discount record if valid, error otherwise.
func ValidateDiscountCode(code string, userID int) (models.DiscountCode, error) {
	var d models.DiscountCode
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT discount_id, code, discount_type,
		       percentage, fixed_amount, expiration_date,
		       is_active, max_uses
		FROM discount_codes
		WHERE code = @code
		  AND is_active = 1
		  AND expiration_date > GETDATE()
	`, sql.Named("code", code)).Scan(
		&d.DiscountID, &d.Code, &d.DiscountType,
		&d.Percentage, &d.FixedAmount, &d.ExpirationDate,
		&d.IsActive, &d.MaxUses,
	)
	if err != nil {
		return d, errors.New("invalid or expired discount code")
	}

	// Check user hasn't already used it
	var usedCount int
	db.Conn.QueryRowContext(context.Background(), `
		SELECT COUNT(*) FROM discount_usage
		WHERE discount_id = @discount_id AND user_id = @user_id
	`,
		sql.Named("discount_id", d.DiscountID),
		sql.Named("user_id", userID),
	).Scan(&usedCount)
	if usedCount > 0 {
		return d, errors.New("you have already used this discount code")
	}

	// Check max_uses not exceeded
	if d.MaxUses != nil {
		var totalUsed int
		db.Conn.QueryRowContext(context.Background(), `
			SELECT COUNT(*) FROM discount_usage
			WHERE discount_id = @discount_id
		`, sql.Named("discount_id", d.DiscountID)).Scan(&totalUsed)
		if totalUsed >= *d.MaxUses {
			return d, errors.New("this discount code has reached its usage limit")
		}
	}

	return d, nil
}

// RecordDiscountUsage marks a code as used by a user and
// auto-deactivates it if max_uses reached
func RecordDiscountUsage(discountID, userID int) error {
	tx, err := db.Conn.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(context.Background(), `
		INSERT INTO discount_usage (discount_id, user_id)
		VALUES (@discount_id, @user_id)
	`,
		sql.Named("discount_id", discountID),
		sql.Named("user_id", userID),
	)
	if err != nil {
		return err
	}

	// Deactivate if max_uses now reached
	_, err = tx.ExecContext(context.Background(), `
		UPDATE discount_codes
		SET is_active = 0
		WHERE discount_id = @discount_id
		  AND max_uses IS NOT NULL
		  AND (
		      SELECT COUNT(*) FROM discount_usage
		      WHERE discount_id = @discount_id
		  ) >= max_uses
	`, sql.Named("discount_id", discountID))
	if err != nil {
		return err
	}

	return tx.Commit()
}

// ApplyDiscount calculates the discounted total
func ApplyDiscount(originalTotal float64, d models.DiscountCode) float64 {
	switch d.DiscountType {
	case "percentage":
		if d.Percentage != nil {
			discount := originalTotal * (*d.Percentage / 100)
			result := originalTotal - discount
			if result < 0 {
				return 0
			}
			return result
		}
	case "fixed":
		if d.FixedAmount != nil {
			result := originalTotal - *d.FixedAmount
			if result < 0 {
				return 0
			}
			return result
		}
	}
	return originalTotal
}

// ── Balance ───────────────────────────────────────────────────

// GetUserBalance returns the current balance for a user
func GetUserBalance(userID int) (float64, error) {
	var balance float64
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT balance FROM users WHERE user_id = @user_id
	`, sql.Named("user_id", userID)).Scan(&balance)
	return balance, err
}

// DeductBalance deducts amount from user balance — returns error if insufficient
func DeductBalance(userID int, amount float64) error {
	result, err := db.Conn.ExecContext(context.Background(), `
		UPDATE users
		SET balance = balance - @amount
		WHERE user_id = @user_id AND balance >= @amount
	`,
		sql.Named("amount", amount),
		sql.Named("user_id", userID),
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("insufficient balance")
	}
	return nil
}

// ── Recommended products ─────────────────────────────────────

func GetRecommendedProducts(userID string) ([]models.RecommendedProduct, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		WITH user_purchases AS (
		    SELECT DISTINCT oi.product_id
		    FROM order_items oi
		    JOIN orders o ON o.order_id = oi.order_id
		    WHERE o.user_id = @user_id
		),
		similar_users AS (
		    SELECT TOP 5
		        o.user_id,
		        COUNT(DISTINCT oi.product_id) AS shared_count
		    FROM order_items oi
		    JOIN orders o ON o.order_id = oi.order_id
		    WHERE oi.product_id IN (SELECT product_id FROM user_purchases)
		      AND o.user_id <> @user_id
		    GROUP BY o.user_id
		    ORDER BY shared_count DESC
		),
		similar_purchases AS (
		    SELECT
		        oi.product_id,
		        COUNT(DISTINCT o.user_id) AS buyers_count
		    FROM order_items oi
		    JOIN orders o ON o.order_id = oi.order_id
		    WHERE o.user_id IN (SELECT user_id FROM similar_users)
		      AND oi.product_id NOT IN (SELECT product_id FROM user_purchases)
		    GROUP BY oi.product_id
		)
		SELECT TOP 20
		    p.product_id, p.name, p.price, p.stall_id,
		    sp.buyers_count
		FROM similar_purchases sp
		JOIN products p ON p.product_id = sp.product_id
		WHERE p.status = 'active'
		ORDER BY sp.buyers_count DESC
	`, sql.Named("user_id", userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.RecommendedProduct
	for rows.Next() {
		var p models.RecommendedProduct
		if err := rows.Scan(
			&p.ProductID, &p.Name, &p.Price,
			&p.StallID, &p.BuyersCount,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// AddBalance tops up a user's balance
func AddBalance(userID int, amount float64) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		UPDATE users SET balance = balance + @amount WHERE user_id = @user_id
	`,
		sql.Named("amount", amount),
		sql.Named("user_id", userID),
	)
	return err
}
