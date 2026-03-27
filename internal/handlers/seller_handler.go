package handlers

import (
	"net/http"
	"strconv"

	"unibazar/project/internal/repository"
	"unibazar/project/internal/session"

	"github.com/gorilla/mux"
)

// ── User: become a seller ─────────────────────────────────────

func BecomeSellerPageHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)
	isSeller, _ := repository.IsSeller(userID)
	hasPending, _ := repository.HasOpenRequest(userID, "become_seller")

	if r.Method == "POST" {
		if !isSeller && !hasPending {
			repository.CreateRequest(userID, "become_seller", r.FormValue("description"))
		}
		http.Redirect(w, r, "/portal", http.StatusSeeOther)
		return
	}

	data := mergeBase(PortalBase(r), map[string]any{
		"Title":      "Become a Seller",
		"IsSeller":   isSeller,
		"HasPending": hasPending,
	})
	Templates["become_seller"].ExecuteTemplate(w, "portal_layout", data)
}

// ── User: request VIP ─────────────────────────────────────────

func RequestVIPHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)
	isVIP, _ := repository.IsVIP(userID)
	hasPending, _ := repository.HasOpenRequest(userID, "become_vip")

	if r.Method == "POST" {
		if !isVIP && !hasPending {
			repository.CreateRequest(userID, "become_vip", r.FormValue("description"))
		}
		http.Redirect(w, r, "/portal", http.StatusSeeOther)
		return
	}

	data := mergeBase(PortalBase(r), map[string]any{
		"Title":      "Request VIP",
		"IsVIP":      isVIP,
		"HasPending": hasPending,
	})
	Templates["request_vip"].ExecuteTemplate(w, "portal_layout", data)
}

// ── User: open a stall ────────────────────────────────────────

func CreateStallPageHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)
	isSeller, _ := repository.IsSeller(userID)
	if !isSeller {
		http.Redirect(w, r, "/portal/become-seller", http.StatusSeeOther)
		return
	}

	seller, err := repository.GetSellerByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if r.Method == "POST" {
		if name := r.FormValue("name"); name != "" {
			repository.CreateStall(seller.SellerID, name)
		}
		http.Redirect(w, r, "/portal", http.StatusSeeOther)
		return
	}

	data := mergeBase(PortalBase(r), map[string]any{
		"Title": "Open a Stall",
	})
	Templates["create_stall"].ExecuteTemplate(w, "portal_layout", data)
}

// ── Supporter: view and handle requests ──────────────────────

func SupporterRequestsHandler(w http.ResponseWriter, r *http.Request) {
	requests, _ := repository.GetOpenRequests()

	data := mergeBase(PortalBase(r), map[string]any{
		"Title":    "User Requests",
		"Requests": requests,
	})
	Templates["supporter_requests"].ExecuteTemplate(w, "portal_layout", data)
}

func SupporterHandleRequestHandler(w http.ResponseWriter, r *http.Request) {
	supporterID := session.GetUserID(r)

	requestID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid request id", 400)
		return
	}

	action := r.FormValue("action")
	status := "resolved"
	if action == "reject" {
		status = "rejected"
	}

	if err := repository.ResolveRequest(requestID, supporterID, status); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if action == "approve" {
		userID, _ := strconv.Atoi(r.FormValue("user_id"))
		switch r.FormValue("request_type") {
		case "become_seller":
			repository.CreateSeller(userID)
			repository.UpdateUserType(userID, "seller") // seller role replaces regular/vip
		case "become_vip":
			repository.GrantVIP(userID, 30)
			repository.UpdateUserType(userID, "vip")
		}
	}

	http.Redirect(w, r, "/supporter/portal/requests", http.StatusSeeOther)
}

// ── Supporter: pending stalls ─────────────────────────────────

func SupporterStallsHandler(w http.ResponseWriter, r *http.Request) {
	stalls, _ := repository.GetPendingStalls()

	data := mergeBase(PortalBase(r), map[string]any{
		"Title":  "Pending Stalls",
		"Stalls": stalls,
	})
	Templates["supporter_stalls"].ExecuteTemplate(w, "portal_layout", data)
}

func SupporterApproveStallHandler(w http.ResponseWriter, r *http.Request) {
	supporterID := session.GetUserID(r)
	id := mux.Vars(r)["id"]
	stallID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid stall id", 400)
		return
	}

	if r.FormValue("action") == "approve" {
		repository.ApproveStall(stallID, supporterID)
	} else {
		repository.UpdateStallStatus(id, "inactive")
	}

	http.Redirect(w, r, "/supporter/portal/stalls", http.StatusSeeOther)
}
