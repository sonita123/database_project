package repository

import (
	"context"
	"database/sql"
	"time"
	"unibazar/project/internal/db"
	"unibazar/project/internal/models"
)

func CreateSession(token string, userID int, role string, expires time.Time) error {
	_, err := db.Conn.ExecContext(context.Background(),
		`INSERT INTO sessions (token, user_id, role, expires_at) 
		 VALUES (@token, @user_id, @role, @expires_at)`,
		sql.Named("token", token),
		sql.Named("user_id", userID),
		sql.Named("role", role),
		sql.Named("expires_at", expires),
	)
	return err
}

func GetSession(token string) (*models.Session, error) {
	var s models.Session
	err := db.Conn.QueryRowContext(context.Background(),
		`SELECT token, user_id, role, expires_at 
		 FROM sessions 
		 WHERE token = @token`,
		sql.Named("token", token),
	).Scan(&s.Token, &s.UserID, &s.Role, &s.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func DeleteSession(token string) error {
	_, err := db.Conn.ExecContext(context.Background(),
		`DELETE FROM sessions WHERE token = @token`,
		sql.Named("token", token),
	)
	return err
}