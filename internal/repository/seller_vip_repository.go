package repository

import (
	"context"
	"database/sql"
	"unibazar/project/internal/db"
	"unibazar/project/internal/models"
)

// ── Requests ──────────────────────────────────────────────────

func CreateRequest(userID int, requestType, description string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		INSERT INTO requests (user_id, request_type, description)
		VALUES (@user_id, @request_type, @description)
	`,
		sql.Named("user_id", userID),
		sql.Named("request_type", requestType),
		sql.Named("description", description),
	)
	return err
}

func HasOpenRequest(userID int, requestType string) (bool, error) {
	var count int
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT COUNT(*) FROM requests
		WHERE user_id = @user_id
		  AND request_type = @request_type
		  AND status IN ('open', 'in_progress')
	`,
		sql.Named("user_id", userID),
		sql.Named("request_type", requestType),
	).Scan(&count)
	return count > 0, err
}

func GetRequestsByUser(userID int) ([]models.Request, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT request_id, user_id, request_type, status,
		       ISNULL(description, ''), created_at
		FROM requests
		WHERE user_id = @user_id
		ORDER BY created_at DESC
	`, sql.Named("user_id", userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reqs []models.Request
	for rows.Next() {
		var req models.Request
		if err := rows.Scan(
			&req.RequestID, &req.UserID, &req.RequestType,
			&req.Status, &req.Description, &req.CreatedAt,
		); err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}
	return reqs, rows.Err()
}

func GetOpenRequests() ([]models.Request, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT r.request_id, r.user_id, r.request_type, r.status,
		       ISNULL(r.description, ''), r.created_at,
		       u.first_name, u.last_name, u.email
		FROM requests r
		JOIN users u ON u.user_id = r.user_id
		WHERE r.status IN ('open', 'in_progress')
		ORDER BY r.created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reqs []models.Request
	for rows.Next() {
		var req models.Request
		if err := rows.Scan(
			&req.RequestID, &req.UserID, &req.RequestType,
			&req.Status, &req.Description, &req.CreatedAt,
			&req.UserFirstName, &req.UserLastName, &req.UserEmail,
		); err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}
	return reqs, rows.Err()
}

func ResolveRequest(requestID, supporterID int, status string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		UPDATE requests
		SET status = @status, handled_by = @handled_by, resolved_at = GETDATE()
		WHERE request_id = @request_id
	`,
		sql.Named("status", status),
		sql.Named("handled_by", supporterID),
		sql.Named("request_id", requestID),
	)
	return err
}

// ── Seller ────────────────────────────────────────────────────

func IsSeller(userID int) (bool, error) {
	var count int
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT COUNT(*) FROM sellers WHERE user_id = @user_id
	`, sql.Named("user_id", userID)).Scan(&count)
	return count > 0, err
}

func CreateSeller(userID int) (int, error) {
	// First check if already a seller
	var sellerID int
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT seller_id FROM sellers WHERE user_id = @user_id
	`, sql.Named("user_id", userID)).Scan(&sellerID)
	if err == nil {
		// Already exists — return existing ID
		return sellerID, nil
	}
	// Insert new seller record
	err = db.Conn.QueryRowContext(context.Background(), `
		INSERT INTO sellers (user_id)
		OUTPUT INSERTED.seller_id
		VALUES (@user_id)
	`, sql.Named("user_id", userID)).Scan(&sellerID)
	return sellerID, err
}

func GetSellerByUserID(userID int) (models.Seller, error) {
	var s models.Seller
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT seller_id, user_id, registered_at
		FROM sellers WHERE user_id = @user_id
	`, sql.Named("user_id", userID)).Scan(
		&s.SellerID, &s.UserID, &s.RegisteredAt,
	)
	return s, err
}

// ── Stalls ────────────────────────────────────────────────────

func CreateStall(sellerID int, name string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		INSERT INTO stalls (seller_id, name, status)
		VALUES (@seller_id, @name, 'pending')
	`,
		sql.Named("seller_id", sellerID),
		sql.Named("name", name),
	)
	return err
}

func ApproveStall(stallID, supporterID int) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		UPDATE stalls
		SET status = 'active', approved_by = @approved_by, approved_at = GETDATE()
		WHERE stall_id = @stall_id
	`,
		sql.Named("approved_by", supporterID),
		sql.Named("stall_id", stallID),
	)
	return err
}

func GetPendingStalls() ([]models.Stall, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT s.stall_id, s.seller_id, s.name, s.status,
		       ISNULL(s.approved_by, 0), s.created_at,
		       u.first_name, u.last_name
		FROM stalls s
		JOIN sellers sel ON sel.seller_id = s.seller_id
		JOIN users u ON u.user_id = sel.user_id
		WHERE s.status = 'pending'
		ORDER BY s.created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stalls []models.Stall
	for rows.Next() {
		var st models.Stall
		if err := rows.Scan(
			&st.StallID, &st.SellerID, &st.Name, &st.Status,
			&st.ApprovedBy, &st.CreatedAt,
			&st.OwnerFirstName, &st.OwnerLastName,
		); err != nil {
			return nil, err
		}
		stalls = append(stalls, st)
	}
	return stalls, rows.Err()
}

// ── VIP ───────────────────────────────────────────────────────

func IsVIP(userID int) (bool, error) {
	var count int
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT COUNT(*) FROM vip_users
		WHERE user_id = @user_id AND end_date > GETDATE()
	`, sql.Named("user_id", userID)).Scan(&count)
	return count > 0, err
}

func GrantVIP(userID int, days int) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		INSERT INTO vip_users (user_id, end_date)
		VALUES (@user_id, DATEADD(day, @days, GETDATE()))
	`,
		sql.Named("user_id", userID),
		sql.Named("days", days),
	)
	return err
}

// GetStallsByUserID returns all stalls belonging to a user (via their seller record)
func GetStallsByUserID(userID int) ([]models.Stall, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT s.stall_id, s.seller_id, s.name, s.status,
		       ISNULL(s.approved_by, 0), s.approved_at, s.created_at
		FROM stalls s
		JOIN sellers sel ON sel.seller_id = s.seller_id
		WHERE sel.user_id = @user_id
		ORDER BY s.created_at DESC
	`, sql.Named("user_id", userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stalls []models.Stall
	for rows.Next() {
		var st models.Stall
		if err := rows.Scan(
			&st.StallID, &st.SellerID, &st.Name, &st.Status,
			&st.ApprovedBy, &st.ApprovedAt, &st.CreatedAt,
		); err != nil {
			return nil, err
		}
		stalls = append(stalls, st)
	}
	return stalls, rows.Err()
}

// GetStallByIDAndUserID fetches a stall only if it belongs to the user
func GetStallByIDAndUserID(stallID, userID int) (models.Stall, error) {
	var st models.Stall
	err := db.Conn.QueryRowContext(context.Background(), `
		SELECT s.stall_id, s.seller_id, s.name, s.status,
		       ISNULL(s.approved_by, 0), s.approved_at, s.created_at
		FROM stalls s
		JOIN sellers sel ON sel.seller_id = s.seller_id
		WHERE s.stall_id = @stall_id AND sel.user_id = @user_id
	`,
		sql.Named("stall_id", stallID),
		sql.Named("user_id", userID),
	).Scan(
		&st.StallID, &st.SellerID, &st.Name, &st.Status,
		&st.ApprovedBy, &st.ApprovedAt, &st.CreatedAt,
	)
	return st, err
}
