package handlers

import (
	"net/http"
	"strconv"

	"unibazar/project/internal/models"
	"unibazar/project/internal/repository"
	"unibazar/project/internal/session"

	"github.com/gorilla/mux"
)

// RequestStallPage renders the stall request form
func RequestStallPage(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sellerID, _ := repository.CreateSellerIfNotExists(userID)
	userStalls, err := repository.GetSellerStalls(sellerID)

	data := map[string]interface{}{
		"Title":  "Request Stall",
		"Role":   "user",
		"Stalls": userStalls,
	}

	err = Templates["stall_request"].ExecuteTemplate(w, "portal_layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// StallViewHandler renders single stall view
func StallViewHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	stallIDStr := vars["id"]
	stallID, err := strconv.Atoi(stallIDStr)
	if err != nil {
		http.Error(w, "Invalid stall ID", http.StatusBadRequest)
		return
	}

	stall, err := repository.GetStallByID(strconv.Itoa(stallID)) // assume exists in repo
	if err != nil {
		http.Error(w, "Stall not found", http.StatusNotFound)
		return
	}

	products, _ := repository.ListSellerProductsByStall(stall.StallID)

	data := map[string]interface{}{
		"Title":    stall.Name,
		"Role":     "user",
		"Stall":    stall,
		"Products": products,
	}

	err = Templates["stall"].ExecuteTemplate(w, "portal_layout", data) // or portal_layout
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// StallManageHandler seller manage own stall
func StallManageHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/seller/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	stallIDStr := vars["id"]
	stallID, err := strconv.Atoi(stallIDStr)
	if err != nil {
		http.Error(w, "Invalid stall ID", http.StatusBadRequest)
		return
	}

	seller, err := repository.GetSellerByUserID(userID)
	if err != nil {
		http.Error(w, "Seller not found", http.StatusForbidden)
		return
	}

	stall, err := repository.GetStallByID(strconv.Itoa(stallID))
	if err != nil || stall.SellerID != seller.SellerID {
		http.Error(w, "Stall not accessible", http.StatusForbidden)
		return
	}

	stalls := []models.Stall{stall} // single for manage
	products, _ := repository.ListSellerProductsByStall(stall.StallID)

	data := map[string]interface{}{
		"Title":    "Manage " + stall.Name,
		"Role":     "seller",
		"Stalls":   stalls,
		"Products": products,
	}

	err = Templates["seller_stalls"].ExecuteTemplate(w, "portal_layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UserRequestsHandler lists user stall/seller requests
func UserRequestsHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	requests, err := repository.GetLastRequests(strconv.Itoa(userID)) // reuse from portal
	if err != nil {
		requests = []models.Request{}
	}

	data := map[string]interface{}{
		"Title":    "My Requests",
		"Role":     "user",
		"Requests": requests,
	}

	// Reuse requests template logic or create new
	err = Templates["user_portal"].ExecuteTemplate(w, "portal_content", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
