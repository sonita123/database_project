package repository

import (
	"context"
	"database/sql"
	"errors"
	"unibazar/project/internal/db"
	"unibazar/project/internal/models"
)

// ── Cart ──────────────────────────────────────────────────────

func GetOrCreateCart(userID int) (models.Cart, error) {
	var c models.Cart
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT TOP 1 cart_id, user_id, locked, created_at
		FROM carts
		WHERE user_id = @user_id AND locked = 0
		ORDER BY created_at DESC
	`, sql.Named("user_id", userID)).Scan(
		&c.CartID, &c.UserID, &c.Locked, &c.CreatedAt,
	)
	if err != nil {
		err = db.Conn.QueryRowContext(context.Background(), `
			INSERT INTO carts (user_id)
			OUTPUT INSERTED.cart_id, INSERTED.user_id, INSERTED.locked, INSERTED.created_at
			VALUES (@user_id)
		`, sql.Named("user_id", userID)).Scan(
			&c.CartID, &c.UserID, &c.Locked, &c.CreatedAt,
		)
	}
	return c, err
}

func GetCartItems(cartID int) ([]models.CartItem, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT ci.cart_item_id, ci.cart_id, ci.product_id, ci.quantity,
		       p.name, p.price
		FROM cart_items ci
		JOIN products p ON p.product_id = ci.product_id
		WHERE ci.cart_id = @cart_id
		ORDER BY ci.cart_item_id
	`, sql.Named("cart_id", cartID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		if err := rows.Scan(
			&item.CartItemID, &item.CartID, &item.ProductID, &item.Quantity,
			&item.ProductName, &item.ProductPrice,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func AddToCart(cartID, productID, quantity int) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		MERGE cart_items AS target
		USING (SELECT @cart_id AS cart_id, @product_id AS product_id) AS source
		ON target.cart_id = source.cart_id AND target.product_id = source.product_id
		WHEN MATCHED THEN
			UPDATE SET quantity = target.quantity + @quantity
		WHEN NOT MATCHED THEN
			INSERT (cart_id, product_id, quantity) VALUES (@cart_id, @product_id, @quantity);
	`,
		sql.Named("cart_id", cartID),
		sql.Named("product_id", productID),
		sql.Named("quantity", quantity),
	)
	return err
}

func UpdateCartItemQty(cartItemID, quantity int) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		UPDATE cart_items SET quantity = @quantity WHERE cart_item_id = @cart_item_id
	`,
		sql.Named("quantity", quantity),
		sql.Named("cart_item_id", cartItemID),
	)
	return err
}

func RemoveCartItem(cartItemID int) error {
	_, err := db.Conn.ExecContext(context.Background(),
		`DELETE FROM cart_items WHERE cart_item_id = @cart_item_id`,
		sql.Named("cart_item_id", cartItemID),
	)
	return err
}

// ── Products ──────────────────────────────────────────────────

func GetActiveProducts(limit, offset int) ([]models.Product, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT product_id, stall_id, name, price, stock, status, created_at
		FROM products
		WHERE status = 'active' AND stock > 0
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

func CountActiveProducts() (int, error) {
	var count int
	return count, db.Conn.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM products WHERE status = 'active' AND stock > 0`,
	).Scan(&count)
}

// ── Checkout ──────────────────────────────────────────────────
// Checkout places an order:
//  1. Validates cart is not empty
//  2. Applies discount if provided (validates, records usage)
//  3. Checks user balance >= final total
//  4. Deducts balance atomically
//  5. Creates order + order_items + decrements stock
//  6. Locks the cart

func Checkout(userID, cartID, addressID int, discountCode string) (int, float64, error) {
	tx, err := db.Conn.BeginTx(context.Background(), nil)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	// Collect cart items
	rows, err := tx.QueryContext(context.Background(), `
		SELECT ci.product_id, ci.quantity, p.price
		FROM cart_items ci
		JOIN products p ON p.product_id = ci.product_id
		WHERE ci.cart_id = @cart_id
	`, sql.Named("cart_id", cartID))
	if err != nil {
		return 0, 0, err
	}
	type lineItem struct {
		productID int
		quantity  int
		price     float64
	}
	var items []lineItem
	var rawTotal float64
	for rows.Next() {
		var li lineItem
		if err := rows.Scan(&li.productID, &li.quantity, &li.price); err != nil {
			rows.Close()
			return 0, 0, err
		}
		rawTotal += li.price * float64(li.quantity)
		items = append(items, li)
	}
	rows.Close()
	if len(items) == 0 {
		return 0, 0, nil
	}

	// Apply discount if provided
	finalTotal := rawTotal
	var appliedDiscountID int
	if discountCode != "" {
		discount, err := ValidateDiscountCode(discountCode, userID)
		if err != nil {
			return 0, 0, err
		}
		finalTotal = ApplyDiscount(rawTotal, discount)
		appliedDiscountID = discount.DiscountID
	}

	// Check balance
	var balance float64
	if err := tx.QueryRowContext(context.Background(), `
		SELECT balance FROM users WHERE user_id = @user_id
	`, sql.Named("user_id", userID)).Scan(&balance); err != nil {
		return 0, 0, err
	}
	if balance < finalTotal {
		return 0, finalTotal, ErrInsufficientBalance
	}

	// Deduct balance
	_, err = tx.ExecContext(context.Background(), `
		UPDATE users SET balance = balance - @amount WHERE user_id = @user_id
	`,
		sql.Named("amount", finalTotal),
		sql.Named("user_id", userID),
	)
	if err != nil {
		return 0, 0, err
	}

	// Create order
	var orderID int
	err = tx.QueryRowContext(context.Background(), `
		INSERT INTO orders (user_id, cart_id, address_id, total_price)
		OUTPUT INSERTED.order_id
		VALUES (@user_id, @cart_id, @address_id, @total_price)
	`,
		sql.Named("user_id", userID),
		sql.Named("cart_id", cartID),
		sql.Named("address_id", addressID),
		sql.Named("total_price", finalTotal),
	).Scan(&orderID)
	if err != nil {
		return 0, 0, err
	}

	// Insert order_items and decrement stock
	for _, li := range items {
		_, err = tx.ExecContext(context.Background(), `
			INSERT INTO order_items (order_id, product_id, quantity, price)
			VALUES (@order_id, @product_id, @quantity, @price)
		`,
			sql.Named("order_id", orderID),
			sql.Named("product_id", li.productID),
			sql.Named("quantity", li.quantity),
			sql.Named("price", li.price),
		)
		if err != nil {
			return 0, 0, err
		}
		_, err = tx.ExecContext(context.Background(), `
			UPDATE products SET stock = stock - @quantity WHERE product_id = @product_id
		`,
			sql.Named("quantity", li.quantity),
			sql.Named("product_id", li.productID),
		)
		if err != nil {
			return 0, 0, err
		}
	}

	// Lock the cart
	_, err = tx.ExecContext(context.Background(), `
		UPDATE carts SET locked = 1 WHERE cart_id = @cart_id
	`, sql.Named("cart_id", cartID))
	if err != nil {
		return 0, 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, err
	}

	// Record discount usage outside transaction (best-effort)
	if appliedDiscountID != 0 {
		RecordDiscountUsage(appliedDiscountID, userID)
	}

	return orderID, finalTotal, nil
}

// ErrInsufficientBalance is returned when user balance is too low
var ErrInsufficientBalance = errors.New("insufficient balance")

// ── Orders ────────────────────────────────────────────────────

func GetUserOrders(userID int) ([]models.Order, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT o.order_id, o.user_id, o.cart_id, o.address_id,
		       o.total_price, o.status, o.order_date,
		       ISNULL(a.street, ''), ISNULL(a.city, ''), ISNULL(a.postal_code, '')
		FROM orders o
		LEFT JOIN addresses a ON a.address_id = o.address_id
		WHERE o.user_id = @user_id
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
			&o.DeliveryStreet, &o.DeliveryCity, &o.DeliveryPostalCode,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}
