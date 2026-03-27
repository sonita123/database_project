package repository

import (
	"context"
	"database/sql"
	"unibazar/project/internal/db"
	"unibazar/project/internal/models"
)

// GetReviewableProducts returns delivered order items the user hasn't reviewed yet
func GetReviewableProducts(userID int) ([]models.ReviewableItem, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT oi.order_item_id, oi.product_id, p.name, p.price,
		       o.order_id, o.order_date
		FROM order_items oi
		JOIN orders o ON o.order_id = oi.order_id
		JOIN products p ON p.product_id = oi.product_id
		WHERE o.user_id = @user_id
		  AND o.status = 'delivered'
		  AND NOT EXISTS (
		      SELECT 1 FROM reviews r
		      WHERE r.user_id = @user_id
		        AND r.product_id = oi.product_id
		  )
		ORDER BY o.order_date DESC
	`, sql.Named("user_id", userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ReviewableItem
	for rows.Next() {
		var item models.ReviewableItem
		if err := rows.Scan(
			&item.OrderItemID, &item.ProductID, &item.ProductName, &item.ProductPrice,
			&item.OrderID, &item.OrderDate,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// GetUserReviews returns all reviews written by a user
func GetUserReviews(userID int) ([]models.Review, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT r.review_id, r.user_id, r.product_id, r.order_item_id,
		       r.rating, ISNULL(r.comment, ''), r.created_at,
		       p.name
		FROM reviews r
		JOIN products p ON p.product_id = r.product_id
		WHERE r.user_id = @user_id
		ORDER BY r.created_at DESC
	`, sql.Named("user_id", userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var rev models.Review
		if err := rows.Scan(
			&rev.ReviewID, &rev.UserID, &rev.ProductID, &rev.OrderItemID,
			&rev.Rating, &rev.Comment, &rev.CreatedAt,
			&rev.ProductName,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, rows.Err()
}

// CreateReview inserts a new review
func CreateReview(userID, productID, orderItemID, rating int, comment string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		INSERT INTO reviews (user_id, product_id, order_item_id, rating, comment)
		VALUES (@user_id, @product_id, @order_item_id, @rating, NULLIF(@comment, ''))
	`,
		sql.Named("user_id", userID),
		sql.Named("product_id", productID),
		sql.Named("order_item_id", orderItemID),
		sql.Named("rating", rating),
		sql.Named("comment", comment),
	)
	return err
}

// DeleteReview removes a review (admin use)
func DeleteReview(reviewID string) error {
	_, err := db.Conn.ExecContext(context.Background(),
		`DELETE FROM reviews WHERE review_id = @id`,
		sql.Named("id", reviewID),
	)
	return err
}

// GetProductReviews returns all reviews for a product (for shop display)
func GetProductReviews(productID int) ([]models.Review, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT r.review_id, r.user_id, r.product_id, r.order_item_id,
		       r.rating, ISNULL(r.comment, ''), r.created_at,
		       u.first_name, u.last_name
		FROM reviews r
		JOIN users u ON u.user_id = r.user_id
		WHERE r.product_id = @product_id
		ORDER BY r.created_at DESC
	`, sql.Named("product_id", productID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var rev models.Review
		if err := rows.Scan(
			&rev.ReviewID, &rev.UserID, &rev.ProductID, &rev.OrderItemID,
			&rev.Rating, &rev.Comment, &rev.CreatedAt,
			&rev.ReviewerFirstName, &rev.ReviewerLastName,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, rows.Err()
}

// CountAllReviews for admin settings page
func CountAllReviews() (int, error) {
	var count int
	return count, db.Conn.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM reviews`,
	).Scan(&count)
}

// GetAllReviewsPaginated for admin settings page
func GetAllReviewsPaginated(limit, offset int) ([]models.Review, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT r.review_id, r.user_id, r.product_id, r.order_item_id,
		       r.rating, ISNULL(r.comment, ''), r.created_at,
		       p.name, u.first_name, u.last_name
		FROM reviews r
		JOIN products p ON p.product_id = r.product_id
		JOIN users u ON u.user_id = r.user_id
		ORDER BY r.created_at DESC
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
		var rev models.Review
		if err := rows.Scan(
			&rev.ReviewID, &rev.UserID, &rev.ProductID, &rev.OrderItemID,
			&rev.Rating, &rev.Comment, &rev.CreatedAt,
			&rev.ProductName, &rev.ReviewerFirstName, &rev.ReviewerLastName,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, rows.Err()
}

// RecordProductView logs that a user viewed a product
func RecordProductView(userID, productID int) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		INSERT INTO product_views (user_id, product_id)
		VALUES (@user_id, @product_id)
	`,
		sql.Named("user_id", userID),
		sql.Named("product_id", productID),
	)
	return err
}

// GetReviewableOrderItem returns the order_item_id for a delivered product
// that the user hasn't reviewed yet — nil if none exists
func GetReviewableOrderItem(userID, productID int) (int, bool) {
	var orderItemID int
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT TOP 1 oi.order_item_id
		FROM order_items oi
		JOIN orders o ON o.order_id = oi.order_id
		WHERE o.user_id = @user_id
		  AND oi.product_id = @product_id
		  AND o.status = 'delivered'
		  AND NOT EXISTS (
		      SELECT 1 FROM reviews r
		      WHERE r.user_id = @user_id
		        AND r.product_id = @product_id
		  )
	`,
		sql.Named("user_id", userID),
		sql.Named("product_id", productID),
	).Scan(&orderItemID)
	if err != nil {
		return 0, false
	}
	return orderItemID, true
}

// HasReviewedProduct checks if user already reviewed a product
func HasReviewedProduct(userID, productID int) bool {
	var count int
	db.Conn.QueryRowContext(context.Background(), `
		SELECT COUNT(*) FROM reviews
		WHERE user_id = @user_id AND product_id = @product_id
	`,
		sql.Named("user_id", userID),
		sql.Named("product_id", productID),
	).Scan(&count)
	return count > 0
}

// GetReviewsForSellerProducts returns all reviews for products in seller's stalls
func GetReviewsForSellerProducts(userID int) ([]models.Review, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT r.review_id, r.user_id, r.product_id, r.order_item_id,
		       r.rating, ISNULL(r.comment, ''), r.created_at,
		       p.name, u.first_name, u.last_name
		FROM reviews r
		JOIN products p ON p.product_id = r.product_id
		JOIN stalls s ON s.stall_id = p.stall_id
		JOIN sellers sel ON sel.seller_id = s.seller_id
		JOIN users u ON u.user_id = r.user_id
		WHERE sel.user_id = @user_id
		ORDER BY r.created_at DESC
	`, sql.Named("user_id", userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var rev models.Review
		if err := rows.Scan(
			&rev.ReviewID, &rev.UserID, &rev.ProductID, &rev.OrderItemID,
			&rev.Rating, &rev.Comment, &rev.CreatedAt,
			&rev.ProductName, &rev.ReviewerFirstName, &rev.ReviewerLastName,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, rows.Err()
}
