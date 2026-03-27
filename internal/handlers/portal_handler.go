package handlers

import (
	"net/http"
	"strconv"

	"unibazar/project/internal/models"
	"unibazar/project/internal/repository"
	"unibazar/project/internal/session"
)

// safeUsers returns an empty slice instead of nil so {{len}} never panics
func safeRequests(s []models.Request) []models.Request {
	if s == nil {
		return []models.Request{}
	}
	return s
}
func safeStalls(s []models.Stall) []models.Stall {
	if s == nil {
		return []models.Stall{}
	}
	return s
}
func safeProducts(s []models.Product) []models.Product {
	if s == nil {
		return []models.Product{}
	}
	return s
}

// ── User portal ───────────────────────────────────────────────

func UserPortalHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)

	id := strconv.Itoa(userID)
	user, err := repository.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", 404)
		return
	}

	requests, _ := repository.GetRequestsByUser(userID)
	discountCodes, _ := repository.GetUserDiscountCodes(id)
	viewedProducts, _ := repository.GetLastViewedProducts(id)
	topProducts, _ := repository.GetTopSellingProducts()
	recommendedProducts, _ := repository.GetRecommendedProducts(id)
	isSeller, _ := repository.IsSeller(userID)
	isVIP, _ := repository.IsVIP(userID)

	// Fetch user's stalls (only if they are a seller)
	var userStalls []models.Stall
	if isSeller {
		userStalls, _ = repository.GetStallsByUserID(userID)
	}

	// Ensure all slices are non-nil so {{len .X}} never panics in templates
	if requests == nil {
		requests = []models.Request{}
	}
	if userStalls == nil {
		userStalls = []models.Stall{}
	}
	if viewedProducts == nil {
		viewedProducts = []models.Product{}
	}
	if topProducts == nil {
		topProducts = []models.TopSellingProduct{}
	}
	if recommendedProducts == nil {
		recommendedProducts = []models.RecommendedProduct{}
	}
	if discountCodes == nil {
		discountCodes = []models.DiscountCodeUsage{}
	}

	data := mergeBase(PortalBase(r), map[string]any{
		"Title":               "My Account",
		"User":                user,
		"Requests":            requests,
		"UserStalls":          userStalls,
		"DiscountCodes":       discountCodes,
		"ViewedProducts":      viewedProducts,
		"TopProducts":         topProducts,
		"RecommendedProducts": recommendedProducts,
		"IsSeller":            isSeller,
		"IsVIP":               isVIP,
		"BalanceAdded":        r.URL.Query().Get("balance_added") == "1",
		"BalanceError":        r.URL.Query().Get("balance_error"),
	})

	if err := Templates["user_portal"].ExecuteTemplate(w, "portal_layout", data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// ── Supporter portal ──────────────────────────────────────────

func SupporterPortalHandler(w http.ResponseWriter, r *http.Request) {
	supporterID := session.GetUserID(r)

	id := strconv.Itoa(supporterID)
	supporter, err := repository.GetSupporterByID(id)
	if err != nil {
		http.Error(w, "Supporter not found", 404)
		return
	}

	kpi, _ := repository.GetSupporterKPIs(id)
	pendingOrders, _ := repository.GetPendingOrders()
	openRequests, _ := repository.GetOpenRequests()
	pendingStalls, _ := repository.GetPendingStalls()

	data := mergeBase(PortalBase(r), map[string]any{
		"Title":         "My Dashboard",
		"Supporter":     supporter,
		"KPI":           kpi,
		"PendingOrders": len(pendingOrders),
		"OpenRequests":  len(openRequests),
		"PendingStalls": len(pendingStalls),
	})

	if err := Templates["supporter_portal"].ExecuteTemplate(w, "portal_layout", data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// AddBalanceHandler — POST /portal/balance/add
func AddBalanceHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)

	amountStr := r.FormValue("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		http.Redirect(w, r, "/portal?balance_error=invalid", http.StatusSeeOther)
		return
	}
	// Cap single top-up at 10,000
	if amount > 10000 {
		http.Redirect(w, r, "/portal?balance_error=too_large", http.StatusSeeOther)
		return
	}

	if err := repository.AddBalance(userID, amount); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/portal?balance_added=1", http.StatusSeeOther)
}
