package handlers

import (
	"net/http"
	"strconv"

	"unibazar/project/internal/models"
	"unibazar/project/internal/repository"
	"unibazar/project/internal/session"
)

const shopPerPage = 12

// ShopHandler — public browsing for guests/buyers.
// Sellers see only their own products so they can preview their listings.
func ShopHandler(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	offset := (page - 1) * shopPerPage

	base := PortalBase(r)
	userID := session.GetUserID(r)

	var products []models.Product
	var total int
	var pageTitle string

	if !isBuyerFromBase(base) && userID != 0 {
		// Seller: show only their own stall products
		products, _ = repository.GetProductsBySellerUserID(userID, shopPerPage, offset)
		total, _ = repository.CountProductsBySellerUserID(userID)
		pageTitle = "My Listings"
	} else {
		// Guest or buyer: show all active products
		products, err = repository.GetActiveProducts(shopPerPage, offset)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		total, _ = repository.CountActiveProducts()
		pageTitle = "Shop"
	}

	if products == nil {
		products = []models.Product{}
	}

	totalPages := total / shopPerPage
	if total%shopPerPage != 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}
	pages := make([]int, totalPages)
	for i := range pages {
		pages[i] = i + 1
	}

	data := mergeBase(base, map[string]any{
		"Title":       pageTitle,
		"Products":    products,
		"Total":       total,
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
		"Pages":       pages,
	})
	Templates["shop"].ExecuteTemplate(w, "portal_layout", data)
}

// CartHandler — requires user login, sellers cannot access
func CartHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)
	base := PortalBase(r)

	if !isBuyerFromBase(base) {
		http.Redirect(w, r, "/seller/dashboard", http.StatusSeeOther)
		return
	}

	cart, err := repository.GetOrCreateCart(userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	items, _ := repository.GetCartItems(cart.CartID)

	var total float64
	for _, item := range items {
		total += item.ProductPrice * float64(item.Quantity)
	}

	data := mergeBase(base, map[string]any{
		"Title": "My Cart",
		"Cart":  cart,
		"Items": items,
		"Total": total,
	})
	Templates["cart"].ExecuteTemplate(w, "portal_layout", data)
}

// AddToCartHandler — redirects to login if not authenticated, blocks sellers
func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login?next=/portal/shop", http.StatusSeeOther)
		return
	}

	if !isBuyerFromBase(PortalBase(r)) {
		http.Redirect(w, r, "/seller/dashboard", http.StatusSeeOther)
		return
	}

	productID, err := strconv.Atoi(r.FormValue("product_id"))
	if err != nil {
		http.Error(w, "invalid product", 400)
		return
	}
	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil || quantity < 1 {
		quantity = 1
	}

	cart, err := repository.GetOrCreateCart(userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	repository.AddToCart(cart.CartID, productID, quantity)
	http.Redirect(w, r, "/portal/cart", http.StatusSeeOther)
}

// RemoveCartItemHandler
func RemoveCartItemHandler(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.Atoi(r.FormValue("cart_item_id"))
	if err != nil {
		http.Error(w, "invalid item", 400)
		return
	}
	repository.RemoveCartItem(itemID)
	http.Redirect(w, r, "/portal/cart", http.StatusSeeOther)
}

// UpdateCartItemHandler
func UpdateCartItemHandler(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.Atoi(r.FormValue("cart_item_id"))
	if err != nil {
		http.Error(w, "invalid item", 400)
		return
	}
	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil || quantity < 1 {
		http.Error(w, "invalid quantity", 400)
		return
	}
	repository.UpdateCartItemQty(itemID, quantity)
	http.Redirect(w, r, "/portal/cart", http.StatusSeeOther)
}

// CheckoutHandler — GET shows summary + address + discount; POST places the order
func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)
	base := PortalBase(r)

	if !isBuyerFromBase(base) {
		http.Redirect(w, r, "/seller/dashboard", http.StatusSeeOther)
		return
	}

	cart, err := repository.GetOrCreateCart(userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	items, _ := repository.GetCartItems(cart.CartID)
	if len(items) == 0 {
		http.Redirect(w, r, "/portal/cart", http.StatusSeeOther)
		return
	}

	var rawTotal float64
	for _, item := range items {
		rawTotal += item.ProductPrice * float64(item.Quantity)
	}

	addresses, _ := repository.GetUserAddresses(userID)
	if addresses == nil {
		addresses = []models.Address{}
	}

	balance, _ := repository.GetUserBalance(userID)

	if r.Method == "GET" {
		data := mergeBase(base, map[string]any{
			"Title":     "Checkout",
			"Items":     items,
			"Total":     rawTotal,
			"Addresses": addresses,
			"Balance":   balance,
		})
		Templates["checkout"].ExecuteTemplate(w, "portal_layout", data)
		return
	}

	// ── POST: place the order ─────────────────────────────────

	addressID, err := strconv.Atoi(r.FormValue("address_id"))
	if err != nil || addressID == 0 {
		def, defErr := repository.GetDefaultAddress(userID)
		if defErr != nil {
			http.Redirect(w, r, "/portal/addresses/add?next=/portal/checkout", http.StatusSeeOther)
			return
		}
		addressID = def.AddressID
	}

	discountCode := r.FormValue("discount_code")
	finalTotal := rawTotal

	if discountCode != "" {
		d, discErr := repository.ValidateDiscountCode(discountCode, userID)
		if discErr != nil {
			data := mergeBase(base, map[string]any{
				"Title":         "Checkout",
				"Items":         items,
				"Total":         rawTotal,
				"Addresses":     addresses,
				"Balance":       balance,
				"DiscountCode":  discountCode,
				"DiscountError": discErr.Error(),
			})
			Templates["checkout"].ExecuteTemplate(w, "portal_layout", data)
			return
		}
		finalTotal = repository.ApplyDiscount(rawTotal, d)
	}

	if balance < finalTotal {
		data := mergeBase(base, map[string]any{
			"Title":        "Checkout",
			"Items":        items,
			"Total":        rawTotal,
			"FinalTotal":   finalTotal,
			"Addresses":    addresses,
			"Balance":      balance,
			"DiscountCode": discountCode,
			"BalanceError": "Insufficient balance. Your balance is " +
				strconv.FormatFloat(balance, 'f', 0, 64) +
				" but the order total is " +
				strconv.FormatFloat(finalTotal, 'f', 0, 64) + ".",
		})
		Templates["checkout"].ExecuteTemplate(w, "portal_layout", data)
		return
	}

	orderID, _, orderErr := repository.Checkout(userID, cart.CartID, addressID, discountCode)
	if orderErr != nil {
		if orderErr == repository.ErrInsufficientBalance {
			data := mergeBase(base, map[string]any{
				"Title":        "Checkout",
				"Items":        items,
				"Total":        rawTotal,
				"FinalTotal":   finalTotal,
				"Addresses":    addresses,
				"Balance":      balance,
				"DiscountCode": discountCode,
				"BalanceError": "Insufficient balance to complete this order.",
			})
			Templates["checkout"].ExecuteTemplate(w, "portal_layout", data)
			return
		}
		http.Error(w, orderErr.Error(), 500)
		return
	}
	if orderID == 0 {
		http.Redirect(w, r, "/portal/cart", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/portal/orders", http.StatusSeeOther)
}

// UserOrdersHandler
func UserOrdersHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)
	orders, _ := repository.GetUserOrders(userID)

	data := mergeBase(PortalBase(r), map[string]any{
		"Title":  "My Orders",
		"Orders": orders,
	})
	Templates["orders"].ExecuteTemplate(w, "portal_layout", data)
}
