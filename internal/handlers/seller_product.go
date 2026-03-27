package handlers

import (
	"net/http"
	"strconv"

	"unibazar/project/internal/models"
	"unibazar/project/internal/repository"
	"unibazar/project/internal/session"

	"github.com/gorilla/mux"
)

// blockIfBuyer redirects buyers/guests away from seller-only pages
func blockIfBuyer(w http.ResponseWriter, r *http.Request) bool {
	if !isBuyerFromBase(PortalBase(r)) {
		return false // is a seller — allowed
	}
	http.Redirect(w, r, "/portal", http.StatusSeeOther)
	return true
}

// blockIfNotSeller redirects non-sellers away from seller-only pages
func blockIfNotSeller(w http.ResponseWriter, r *http.Request) bool {
	userID := session.GetUserID(r)
	isSeller, _ := repository.IsSeller(userID)
	if !isSeller {
		http.Redirect(w, r, "/portal", http.StatusSeeOther)
		return true
	}
	return false
}

// ── Seller dashboard ──────────────────────────────────────────

func SellerDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if blockIfNotSeller(w, r) {
		return
	}
	userID := session.GetUserID(r)

	stalls, _ := repository.GetStallsByUserID(userID)
	if stalls == nil {
		stalls = []models.Stall{}
	}

	stats, _ := repository.GetSellerStats(userID)
	topProducts, _ := repository.GetSellerTopProducts(userID)
	recentOrders, _ := repository.GetSellerRecentOrders(userID)

	if topProducts == nil {
		topProducts = []models.SellerTopProduct{}
	}
	if recentOrders == nil {
		recentOrders = []models.Order{}
	}

	data := mergeBase(PortalBase(r), map[string]any{
		"Title":        "Seller Dashboard",
		"Stalls":       stalls,
		"Stats":        stats,
		"TopProducts":  topProducts,
		"RecentOrders": recentOrders,
	})
	Templates["seller_dashboard"].ExecuteTemplate(w, "portal_layout", data)
}

// ── Seller: manage products in their stalls ───────────────────

func SellerProductsHandler(w http.ResponseWriter, r *http.Request) {
	if blockIfNotSeller(w, r) {
		return
	}
	userID := session.GetUserID(r)
	stallID, err := strconv.Atoi(mux.Vars(r)["stall_id"])
	if err != nil {
		http.Error(w, "invalid stall", 400)
		return
	}

	stall, err := repository.GetStallByIDAndUserID(stallID, userID)
	if err != nil {
		http.Error(w, "Stall not found", 404)
		return
	}

	products, _ := repository.GetProductsByStall(stallID)
	if products == nil {
		products = []models.Product{}
	}

	data := mergeBase(PortalBase(r), map[string]any{
		"Title":    stall.Name + " — Products",
		"Stall":    stall,
		"Products": products,
	})
	Templates["seller_products"].ExecuteTemplate(w, "portal_layout", data)
}

func SellerAddProductHandler(w http.ResponseWriter, r *http.Request) {
	if blockIfNotSeller(w, r) {
		return
	}
	userID := session.GetUserID(r)
	stallID, err := strconv.Atoi(mux.Vars(r)["stall_id"])
	if err != nil {
		http.Error(w, "invalid stall", 400)
		return
	}

	stall, err := repository.GetStallByIDAndUserID(stallID, userID)
	if err != nil {
		http.Error(w, "Stall not found", 404)
		return
	}

	if r.Method == "GET" {
		data := mergeBase(PortalBase(r), map[string]any{
			"Title":   "Add Product — " + stall.Name,
			"Stall":   stall,
			"Product": models.Product{},
		})
		Templates["seller_product_form"].ExecuteTemplate(w, "portal_layout", data)
		return
	}

	name := r.FormValue("name")
	priceStr := r.FormValue("price")
	stockStr := r.FormValue("stock")

	if name == "" || priceStr == "" {
		data := mergeBase(PortalBase(r), map[string]any{
			"Title":   "Add Product — " + stall.Name,
			"Stall":   stall,
			"Product": models.Product{},
			"Error":   "Name and price are required",
		})
		Templates["seller_product_form"].ExecuteTemplate(w, "portal_layout", data)
		return
	}

	price, priceErr := strconv.ParseFloat(priceStr, 64)
	if priceErr != nil || price <= 0 {
		http.Error(w, "invalid price", 400)
		return
	}
	stock, stockErr := strconv.Atoi(stockStr)
	if stockErr != nil || stock < 0 {
		stock = 0
	}

	if err := repository.CreateProductInStall(stallID, name, price, stock); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/seller/stalls/"+mux.Vars(r)["stall_id"]+"/products", http.StatusSeeOther)
}

func SellerEditProductHandler(w http.ResponseWriter, r *http.Request) {
	if blockIfNotSeller(w, r) {
		return
	}
	userID := session.GetUserID(r)

	stallID, err := strconv.Atoi(mux.Vars(r)["stall_id"])
	if err != nil {
		http.Error(w, "invalid stall", 400)
		return
	}
	productID, err := strconv.Atoi(mux.Vars(r)["product_id"])
	if err != nil {
		http.Error(w, "invalid product", 400)
		return
	}

	stall, err := repository.GetStallByIDAndUserID(stallID, userID)
	if err != nil {
		http.Error(w, "Stall not found", 404)
		return
	}

	product, err := repository.GetProductByIDInt(productID)
	if err != nil || product.StallID != stallID {
		http.Error(w, "Product not found", 404)
		return
	}

	if r.Method == "GET" {
		data := mergeBase(PortalBase(r), map[string]any{
			"Title":   "Edit Product",
			"Stall":   stall,
			"Product": product,
		})
		Templates["seller_product_form"].ExecuteTemplate(w, "portal_layout", data)
		return
	}

	name := r.FormValue("name")
	priceStr := r.FormValue("price")
	stockStr := r.FormValue("stock")
	status := r.FormValue("status")

	price, priceErr := strconv.ParseFloat(priceStr, 64)
	if priceErr != nil || price <= 0 {
		http.Error(w, "invalid price", 400)
		return
	}
	stock, stockErr := strconv.Atoi(stockStr)
	if stockErr != nil || stock < 0 {
		stock = 0
	}
	if status == "" {
		status = "active"
	}

	if err := repository.UpdateProductInStall(productID, name, price, stock, status); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/seller/stalls/"+mux.Vars(r)["stall_id"]+"/products", http.StatusSeeOther)
}

func SellerDeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if blockIfNotSeller(w, r) {
		return
	}
	userID := session.GetUserID(r)

	stallID, err := strconv.Atoi(mux.Vars(r)["stall_id"])
	if err != nil {
		http.Error(w, "invalid stall", 400)
		return
	}
	productID, err := strconv.Atoi(mux.Vars(r)["product_id"])
	if err != nil {
		http.Error(w, "invalid product", 400)
		return
	}

	// Verify ownership before deleting
	if _, err = repository.GetStallByIDAndUserID(stallID, userID); err != nil {
		http.Error(w, "Stall not found", 404)
		return
	}

	repository.DeleteProductFromStall(productID, stallID)
	http.Redirect(w, r, "/seller/stalls/"+mux.Vars(r)["stall_id"]+"/products", http.StatusSeeOther)
}

// ── Seller: view their own product reviews ────────────────────

func SellerReviewsHandler(w http.ResponseWriter, r *http.Request) {
	if blockIfNotSeller(w, r) {
		return
	}
	userID := session.GetUserID(r)

	reviews, _ := repository.GetReviewsForSellerProducts(userID)
	if reviews == nil {
		reviews = []models.Review{}
	}

	data := mergeBase(PortalBase(r), map[string]any{
		"Title":   "Product Reviews",
		"Reviews": reviews,
	})
	Templates["seller_reviews"].ExecuteTemplate(w, "portal_layout", data)
}
