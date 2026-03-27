package repository

import (
	"context"
	"database/sql"
	"unibazar/project/internal/db"
	"unibazar/project/internal/models"
)

func GetSupportersPaginated(limit, offset int) ([]models.Supporter, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT supporter_id, first_name, last_name,
		       ISNULL(email, ''), ISNULL(image_url, ''), ISNULL(password, '')
		FROM supporters
		ORDER BY supporter_id
		OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY
	`,
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var supporters []models.Supporter
	for rows.Next() {
		var s models.Supporter
		if err := rows.Scan(
			&s.SupporterID, &s.FirstName, &s.LastName,
			&s.Email, &s.ImageURL, &s.Password,
		); err != nil {
			return nil, err
		}
		supporters = append(supporters, s)
	}
	return supporters, rows.Err()
}

func CountSupporters() (int, error) {
	var count int
	return count, db.Conn.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM supporters`,
	).Scan(&count)
}

func GetSupporterByID(id string) (models.Supporter, error) {
	var s models.Supporter
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT supporter_id, first_name, last_name,
		       ISNULL(email, ''), ISNULL(image_url, ''), ISNULL(password, '')
		FROM supporters
		WHERE supporter_id = @id
	`, sql.Named("id", id)).Scan(
		&s.SupporterID, &s.FirstName, &s.LastName,
		&s.Email, &s.ImageURL, &s.Password,
	)
	return s, err
}

func CreateSupporter(firstName, lastName, email, imageURL string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		INSERT INTO supporters (first_name, last_name, email, image_url)
		VALUES (@first_name, @last_name, NULLIF(@email, ''), NULLIF(@image_url, ''))
	`,
		sql.Named("first_name", firstName),
		sql.Named("last_name", lastName),
		sql.Named("email", email),
		sql.Named("image_url", imageURL),
	)
	return err
}

func UpdateSupporter(id, firstName, lastName, email, imageURL string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		UPDATE supporters
		SET first_name = @first_name,
		    last_name  = @last_name,
		    email      = NULLIF(@email, ''),
		    image_url  = NULLIF(@image_url, '')
		WHERE supporter_id = @id
	`,
		sql.Named("first_name", firstName),
		sql.Named("last_name", lastName),
		sql.Named("email", email),
		sql.Named("image_url", imageURL),
		sql.Named("id", id),
	)
	return err
}

func DeleteSupporter(id string) error {
	_, err := db.Conn.ExecContext(context.Background(),
		`DELETE FROM supporters WHERE supporter_id = @id`,
		sql.Named("id", id),
	)
	return err
}

func GetAllSupportersSimple() ([]models.Supporter, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT supporter_id, first_name, last_name,
		       ISNULL(email, ''), ISNULL(image_url, ''), ISNULL(password, '')
		FROM supporters
		ORDER BY supporter_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var supporters []models.Supporter
	for rows.Next() {
		var s models.Supporter
		if err := rows.Scan(
			&s.SupporterID, &s.FirstName, &s.LastName,
			&s.Email, &s.ImageURL, &s.Password,
		); err != nil {
			return nil, err
		}
		supporters = append(supporters, s)
	}
	return supporters, rows.Err()
}
