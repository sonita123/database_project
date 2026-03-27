package repository

import (
	"context"
	"database/sql"
	"unibazar/project/internal/db"
	"unibazar/project/internal/models"
)

// GetUserAddresses returns all addresses for a user
func GetUserAddresses(userID int) ([]models.Address, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT address_id, user_id, city, street,
		       ISNULL(postal_code, ''), is_default
		FROM addresses
		WHERE user_id = @user_id
		ORDER BY is_default DESC, address_id ASC
	`, sql.Named("user_id", userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []models.Address
	for rows.Next() {
		var a models.Address
		if err := rows.Scan(
			&a.AddressID, &a.UserID, &a.City,
			&a.Street, &a.PostalCode, &a.IsDefault,
		); err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	return addresses, rows.Err()
}

// GetDefaultAddress returns the user's default address, or any address if none set as default
func GetDefaultAddress(userID int) (models.Address, error) {
	var a models.Address
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT TOP 1 address_id, user_id, city, street,
		       ISNULL(postal_code, ''), is_default
		FROM addresses
		WHERE user_id = @user_id
		ORDER BY is_default DESC, address_id ASC
	`, sql.Named("user_id", userID)).Scan(
		&a.AddressID, &a.UserID, &a.City,
		&a.Street, &a.PostalCode, &a.IsDefault,
	)
	return a, err
}

// GetAddressByID returns a single address
func GetAddressByID(addressID, userID int) (models.Address, error) {
	var a models.Address
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT address_id, user_id, city, street,
		       ISNULL(postal_code, ''), is_default
		FROM addresses
		WHERE address_id = @address_id AND user_id = @user_id
	`,
		sql.Named("address_id", addressID),
		sql.Named("user_id", userID),
	).Scan(
		&a.AddressID, &a.UserID, &a.City,
		&a.Street, &a.PostalCode, &a.IsDefault,
	)
	return a, err
}

// CreateAddress adds a new address for a user
func CreateAddress(userID int, city, street, postalCode string, isDefault bool) (int, error) {
	// If setting as default, clear existing default first
	if isDefault {
		db.Conn.ExecContext(context.Background(), `
			UPDATE addresses SET is_default = 0
			WHERE user_id = @user_id
		`, sql.Named("user_id", userID))
	}

	var addressID int
	err := db.Conn.QueryRowContext(context.Background(), `
		INSERT INTO addresses (user_id, city, street, postal_code, is_default)
		OUTPUT INSERTED.address_id
		VALUES (@user_id, @city, @street, NULLIF(@postal_code, ''), @is_default)
	`,
		sql.Named("user_id", userID),
		sql.Named("city", city),
		sql.Named("street", street),
		sql.Named("postal_code", postalCode),
		sql.Named("is_default", isDefault),
	).Scan(&addressID)
	return addressID, err
}

// UpdateAddress updates an existing address
func UpdateAddress(addressID, userID int, city, street, postalCode string, isDefault bool) error {
	if isDefault {
		db.Conn.ExecContext(context.Background(), `
			UPDATE addresses SET is_default = 0
			WHERE user_id = @user_id
		`, sql.Named("user_id", userID))
	}

	_, err := db.Conn.ExecContext(context.Background(), `
		UPDATE addresses
		SET city        = @city,
		    street      = @street,
		    postal_code = NULLIF(@postal_code, ''),
		    is_default  = @is_default
		WHERE address_id = @address_id AND user_id = @user_id
	`,
		sql.Named("city", city),
		sql.Named("street", street),
		sql.Named("postal_code", postalCode),
		sql.Named("is_default", isDefault),
		sql.Named("address_id", addressID),
		sql.Named("user_id", userID),
	)
	return err
}

// DeleteAddress removes an address
func DeleteAddress(addressID, userID int) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		DELETE FROM addresses
		WHERE address_id = @address_id AND user_id = @user_id
	`,
		sql.Named("address_id", addressID),
		sql.Named("user_id", userID),
	)
	return err
}

// SetDefaultAddress marks one address as default and clears others
func SetDefaultAddress(addressID, userID int) error {
	tx, err := db.Conn.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(context.Background(), `
		UPDATE addresses SET is_default = 0 WHERE user_id = @user_id
	`, sql.Named("user_id", userID))
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(context.Background(), `
		UPDATE addresses SET is_default = 1
		WHERE address_id = @address_id AND user_id = @user_id
	`,
		sql.Named("address_id", addressID),
		sql.Named("user_id", userID),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
