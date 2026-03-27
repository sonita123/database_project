package repository

import (
	"context"
	"database/sql"
	"unibazar/project/internal/db"
	"unibazar/project/internal/models"
)

func GetUserByEmail(email string) (models.User, error) {
	var u models.User
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT user_id, first_name, last_name, email,
		       ISNULL(password, ''), balance, user_type, created_at
		FROM users
		WHERE email = @email
	`, sql.Named("email", email)).Scan(
		&u.UserID, &u.FirstName, &u.LastName, &u.Email,
		&u.Password, &u.Balance, &u.UserType, &u.CreatedAt,
	)
	return u, err
}

func SetUserPassword(userID int, hashedPassword string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		UPDATE users SET password = @password WHERE user_id = @user_id
	`,
		sql.Named("password", hashedPassword),
		sql.Named("user_id", userID),
	)
	return err
}

func GetSupporterByEmail(email string) (models.Supporter, error) {
	var s models.Supporter
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT supporter_id, first_name, last_name,
		       ISNULL(email, ''), ISNULL(image_url, ''), ISNULL(password, '')
		FROM supporters
		WHERE email = @email
	`, sql.Named("email", email)).Scan(
		&s.SupporterID, &s.FirstName, &s.LastName,
		&s.Email, &s.ImageURL, &s.Password,
	)
	return s, err
}

func GetSupporterByIDForAuth(supporterID string) (models.Supporter, error) {
	var s models.Supporter
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT supporter_id, first_name, last_name,
		       ISNULL(email, ''), ISNULL(image_url, ''), ISNULL(password, '')
		FROM supporters
		WHERE supporter_id = @supporter_id
	`, sql.Named("supporter_id", supporterID)).Scan(
		&s.SupporterID, &s.FirstName, &s.LastName,
		&s.Email, &s.ImageURL, &s.Password,
	)
	return s, err
}

func SetSupporterPassword(supporterID int, hashedPassword string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		UPDATE supporters SET password = @password WHERE supporter_id = @supporter_id
	`,
		sql.Named("password", hashedPassword),
		sql.Named("supporter_id", supporterID),
	)
	return err
}

func GetAdminByUsername(username string) (models.Admin, error) {
	var a models.Admin
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT admin_id, username, password, created_at
		FROM admins
		WHERE username = @username
	`, sql.Named("username", username)).Scan(
		&a.AdminID, &a.Username, &a.Password, &a.CreatedAt,
	)
	return a, err
}

func CreateAdmin(username, hashedPassword string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		INSERT INTO admins (username, password)
		VALUES (@username, @password)
	`,
		sql.Named("username", username),
		sql.Named("password", hashedPassword),
	)
	return err
}

func GetAdminByID(id int) (models.Admin, error) {
	var a models.Admin
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT admin_id, username, password, created_at
		FROM admins
		WHERE admin_id = @id
	`, sql.Named("id", id)).Scan(
		&a.AdminID, &a.Username, &a.Password, &a.CreatedAt,
	)
	return a, err
}
