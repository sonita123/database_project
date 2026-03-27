package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"unibazar/project/internal/db"
	"unibazar/project/internal/models"

	"unibazar/project/internal/repository"

	"github.com/gorilla/mux"
)

func SupporterOrdersHandler(w http.ResponseWriter, r *http.Request) {
	tab := r.URL.Query().Get("tab")
	if tab == "" {
		tab = "pending"
	}

	var orders interface{}
	if tab == "all" {
		orders, _ = repository.GetAllOrdersForSupporter()
	} else {
		orders, _ = repository.GetPendingOrders()
	}

	data := mergeBase(PortalBase(r), map[string]any{
		"Title":  "Orders",
		"Orders": orders,
		"Tab":    tab,
		"Role":   "supporter",
	})
	Templates["supporter_orders"].ExecuteTemplate(w, "portal_layout", data)
}

func SupporterUpdateOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid order id", 400)
		return
	}

	status := r.FormValue("status")
	allowed := map[string]bool{
		"confirmed": true, "shipped": true,
		"delivered": true, "cancelled": true,
	}
	if !allowed[status] {
		http.Error(w, "invalid status", 400)
		return
	}

	if err := repository.UpdateOrderStatus(orderID, status); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/supporter/portal/orders", http.StatusSeeOther)
}

// StallsHandlerForSupporter — alias that routes to the supporter stalls page
func StallsHandlerForSupporter(w http.ResponseWriter, r *http.Request) {
	SupporterStallsHandler(w, r)
}

func GetPendingOrders() ([]models.Order, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT o.order_id, o.user_id, o.cart_id, o.address_id,
		       o.total_price, o.status, o.order_date,
		       u.first_name, u.last_name, u.email,
		       ISNULL(a.street, ''), ISNULL(a.city, ''), ISNULL(a.postal_code, '')
		FROM orders o
		JOIN users u ON u.user_id = o.user_id
		LEFT JOIN addresses a ON a.address_id = o.address_id
		WHERE o.status = 'pending'
		ORDER BY o.order_date ASC
	`)
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
			&o.UserFirstName, &o.UserLastName, &o.UserEmail,
			&o.DeliveryStreet, &o.DeliveryCity, &o.DeliveryPostalCode,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func GetAllOrdersForSupporter() ([]models.Order, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT TOP 100
		       o.order_id, o.user_id, o.cart_id, o.address_id,
		       o.total_price, o.status, o.order_date,
		       u.first_name, u.last_name, u.email,
		       ISNULL(a.street, ''), ISNULL(a.city, ''), ISNULL(a.postal_code, '')
		FROM orders o
		JOIN users u ON u.user_id = o.user_id
		LEFT JOIN addresses a ON a.address_id = o.address_id
		ORDER BY o.order_date DESC
	`)
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
			&o.UserFirstName, &o.UserLastName, &o.UserEmail,
			&o.DeliveryStreet, &o.DeliveryCity, &o.DeliveryPostalCode,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func UpdateOrderStatus(orderID int, status string) error {
	_, err := db.Conn.ExecContext(context.Background(), `
		UPDATE orders SET status = @status WHERE order_id = @order_id
	`,
		sql.Named("status", status),
		sql.Named("order_id", orderID),
	)
	return err
}

func GetOrderItems(orderID int) ([]models.OrderItem, error) {
	rows, err := db.Conn.QueryContext(context.Background(), `
		SELECT oi.order_item_id, oi.order_id, oi.product_id,
		       oi.quantity, oi.price, p.name
		FROM order_items oi
		JOIN products p ON p.product_id = oi.product_id
		WHERE oi.order_id = @order_id
	`, sql.Named("order_id", orderID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(
			&item.OrderItemID, &item.OrderID, &item.ProductID,
			&item.Quantity, &item.Price, &item.ProductName,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
