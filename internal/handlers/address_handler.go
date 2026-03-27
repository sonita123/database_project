package handlers

import (
	"net/http"
	"strconv"

	"unibazar/project/internal/models"
	"unibazar/project/internal/repository"
	"unibazar/project/internal/session"

	"github.com/gorilla/mux"
)

// AddressesHandler — GET /portal/addresses
func AddressesHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)

	addresses, _ := repository.GetUserAddresses(userID)
	if addresses == nil {
		addresses = []models.Address{}
	}

	data := mergeBase(PortalBase(r), map[string]any{
		"Title":     "My Addresses",
		"Addresses": addresses,
	})
	Templates["addresses"].ExecuteTemplate(w, "portal_layout", data)
}

// AddAddressHandler — GET/POST /portal/addresses/add
func AddAddressHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)

	if r.Method == "GET" {
		next := r.URL.Query().Get("next")
		data := mergeBase(PortalBase(r), map[string]any{
			"Title": "Add Address",
			"Next":  next,
		})
		Templates["address_form"].ExecuteTemplate(w, "portal_layout", data)
		return
	}

	city := r.FormValue("city")
	street := r.FormValue("street")
	postalCode := r.FormValue("postal_code")
	isDefault := r.FormValue("is_default") == "on"

	if city == "" || street == "" {
		data := mergeBase(PortalBase(r), map[string]any{
			"Title": "Add Address",
			"Error": "City and street are required",
			"Next":  r.FormValue("next"),
		})
		Templates["address_form"].ExecuteTemplate(w, "portal_layout", data)
		return
	}

	_, err := repository.CreateAddress(userID, city, street, postalCode, isDefault)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	next := r.FormValue("next")
	if next == "" {
		next = "/portal/addresses"
	}
	http.Redirect(w, r, next, http.StatusSeeOther)
}

// EditAddressHandler — GET/POST /portal/addresses/{id}/edit
func EditAddressHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)
	addressID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid address", 400)
		return
	}

	address, err := repository.GetAddressByID(addressID, userID)
	if err != nil {
		http.Error(w, "Address not found", 404)
		return
	}

	if r.Method == "GET" {
		data := mergeBase(PortalBase(r), map[string]any{
			"Title":   "Edit Address",
			"Address": address,
		})
		Templates["address_form"].ExecuteTemplate(w, "portal_layout", data)
		return
	}

	city := r.FormValue("city")
	street := r.FormValue("street")
	postalCode := r.FormValue("postal_code")
	isDefault := r.FormValue("is_default") == "on"

	if city == "" || street == "" {
		data := mergeBase(PortalBase(r), map[string]any{
			"Title":   "Edit Address",
			"Address": address,
			"Error":   "City and street are required",
		})
		Templates["address_form"].ExecuteTemplate(w, "portal_layout", data)
		return
	}

	if err := repository.UpdateAddress(addressID, userID, city, street, postalCode, isDefault); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/portal/addresses", http.StatusSeeOther)
}

// DeleteAddressHandler — POST /portal/addresses/{id}/delete
func DeleteAddressHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)
	addressID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid address", 400)
		return
	}
	repository.DeleteAddress(addressID, userID)
	http.Redirect(w, r, "/portal/addresses", http.StatusSeeOther)
}

// SetDefaultAddressHandler — POST /portal/addresses/{id}/default
func SetDefaultAddressHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)
	addressID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid address", 400)
		return
	}
	repository.SetDefaultAddress(addressID, userID)
	http.Redirect(w, r, "/portal/addresses", http.StatusSeeOther)
}
