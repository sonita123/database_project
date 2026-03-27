package repository

import (
	"context"
	"database/sql"
	"unibazar/project/internal/db"
	"unibazar/project/internal/models"
)

func GetRecentUsers(limit int) ([]models.User, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT TOP (@limit) user_id, first_name, last_name, email,
		       ISNULL(password, ''), balance, user_type, created_at
		FROM users
		ORDER BY user_id DESC
	`, sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.UserID, &u.FirstName, &u.LastName, &u.Email,
			&u.Password, &u.Balance, &u.UserType, &u.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func GetUserByID(id string) (models.User, error) {
	var u models.User
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT user_id, first_name, last_name, email,
		       ISNULL(password, ''), balance, user_type, created_at
		FROM users
		WHERE user_id = @id
	`, sql.Named("id", id)).Scan(
		&u.UserID, &u.FirstName, &u.LastName, &u.Email,
		&u.Password, &u.Balance, &u.UserType, &u.CreatedAt,
	)
	return u, err
}

func GetUsersPaginated(limit, offset int) ([]models.User, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT user_id, first_name, last_name, email,
		       ISNULL(password, ''), balance, user_type, created_at
		FROM users
		ORDER BY user_id
		OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY
	`,
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.UserID, &u.FirstName, &u.LastName, &u.Email,
			&u.Password, &u.Balance, &u.UserType, &u.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func GetAllUsers() ([]models.User, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT user_id, first_name, last_name, email,
		       ISNULL(password, ''), balance, user_type, created_at
		FROM users
		ORDER BY user_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.UserID, &u.FirstName, &u.LastName, &u.Email,
			&u.Password, &u.Balance, &u.UserType, &u.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func CreateUser(firstName, lastName, email, balance, userType string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		INSERT INTO users (first_name, last_name, email, balance, user_type)
		VALUES (@first_name, @last_name, @email, @balance, @user_type)
	`,
		sql.Named("first_name", firstName),
		sql.Named("last_name", lastName),
		sql.Named("email", email),
		sql.Named("balance", balance),
		sql.Named("user_type", userType),
	)
	return err
}

func UpdateUser(id, firstName, lastName, email, balance, userType string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		UPDATE users
		SET first_name = @first_name,
		    last_name  = @last_name,
		    email      = @email,
		    balance    = @balance,
		    user_type  = @user_type
		WHERE user_id = @id
	`,
		sql.Named("first_name", firstName),
		sql.Named("last_name", lastName),
		sql.Named("email", email),
		sql.Named("balance", balance),
		sql.Named("user_type", userType),
		sql.Named("id", id),
	)
	return err
}

func DeleteUser(id string) error {
	_, err := db.Conn.ExecContext(context.Background(),
		`DELETE FROM users WHERE user_id = @id`,
		sql.Named("id", id),
	)
	return err
}

func UpdateUserType(userID int, userType string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		UPDATE users SET user_type = @user_type WHERE user_id = @user_id
	`,
		sql.Named("user_type", userType),
		sql.Named("user_id", userID),
	)
	return err
}

func GetTopSellingProducts() ([]models.TopSellingProduct, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT TOP 10 product_id, name, price, stall_id, total_sold
		FROM top_selling_products
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.TopSellingProduct
	for rows.Next() {
		var p models.TopSellingProduct
		if err := rows.Scan(&p.ProductID, &p.Name, &p.Price, &p.StallID, &p.TotalSold); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func GetLastViewedProducts(userID string) ([]models.Product, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT TOP 10 p.product_id, p.name, p.price, p.stall_id, p.stock, p.status, p.created_at
		FROM product_views v
		JOIN products p ON p.product_id = v.product_id
		WHERE v.user_id = @user_id
		ORDER BY v.viewed_at DESC
	`, sql.Named("user_id", userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ProductID, &p.Name, &p.Price, &p.StallID, &p.Stock, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func GetLastRequests(userID string) ([]models.Request, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT TOP 10 request_id, user_id, request_type, status,
		       ISNULL(description, ''), created_at
		FROM requests
		WHERE user_id = @user_id
		ORDER BY created_at DESC
	`, sql.Named("user_id", userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.Request
	for rows.Next() {
		var req models.Request
		if err := rows.Scan(
			&req.RequestID, &req.UserID, &req.RequestType,
			&req.Status, &req.Description, &req.CreatedAt,
		); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, rows.Err()
}

func RegisterUser(firstName, lastName, email, hashedPassword string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		INSERT INTO users (first_name, last_name, email, password, balance, user_type)
		VALUES (@first_name, @last_name, @email, @password, 0, 'regular')
	`,
		sql.Named("first_name", firstName),
		sql.Named("last_name", lastName),
		sql.Named("email", email),
		sql.Named("password", hashedPassword),
	)
	return err
}

// EmailExists checks whether an email is already registered
func EmailExists(email string) (bool, error) {
	var count int
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT COUNT(*) FROM users WHERE email = @email
	`, sql.Named("email", email)).Scan(&count)
	return count > 0, err
}
