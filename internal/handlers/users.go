package handlers

import (
	"net/http"
	"strconv"
	"unibazar/project/internal/models"
	"unibazar/project/internal/repository"

	"github.com/gorilla/mux"
)

const usersPerPage = 10

type UsersPageData struct {
	Title       string
	CurrentPath string
	Users       []models.User
	TotalUsers  int
	CurrentPage int
	TotalPages  int
	PrevPage    int
	NextPage    int
	Pages       []int
}

func UsersPageHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * usersPerPage

	users, err := repository.GetUsersPaginated(usersPerPage, offset)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	totalUsers, err := repository.CountUsers()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	totalPages := totalUsers / usersPerPage
	if totalUsers%usersPerPage != 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}

	pages := make([]int, totalPages)
	for i := range pages {
		pages[i] = i + 1
	}

	data := UsersPageData{
		Title:       "Users",
		CurrentPath: r.URL.Path,
		Users:       users,
		TotalUsers:  totalUsers,
		CurrentPage: page,
		TotalPages:  totalPages,
		PrevPage:    page - 1,
		NextPage:    page + 1,
		Pages:       pages,
	}

	err = Templates["users"].ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func ViewUserHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	user, err := repository.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", 404)
		return
	}

	requests, err := repository.GetLastRequests(id)
	if err != nil {
		requests = nil
	}

	discountCodes, err := repository.GetUserDiscountCodes(id)
	if err != nil {
		discountCodes = nil
	}

	viewedProducts, err := repository.GetLastViewedProducts(id)
	if err != nil {
		viewedProducts = nil
	}

	topProducts, err := repository.GetTopSellingProducts()
	if err != nil {
		topProducts = nil
	}

	recommendedProducts, err := repository.GetRecommendedProducts(id)
	if err != nil {
		recommendedProducts = nil
	}

	err = Templates["user"].ExecuteTemplate(w, "layout", map[string]any{
		"Title":               user.FirstName + " " + user.LastName,
		"CurrentPath":         "/users",
		"User":                user,
		"Requests":            requests,
		"DiscountCodes":       discountCodes,
		"ViewedProducts":      viewedProducts,
		"TopProducts":         topProducts,
		"RecommendedProducts": recommendedProducts,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := Templates["create_user"].ExecuteTemplate(w, "layout", map[string]any{
			"Title":       "Create User",
			"CurrentPath": r.URL.Path,
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	balance := r.FormValue("balance")
	userType := r.FormValue("user_type")

	err := repository.CreateUser(firstName, lastName, email, balance, userType)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func EditUserHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if r.Method == "GET" {
		user, err := repository.GetUserByID(id)
		if err != nil {
			http.Error(w, "User not found", 404)
			return
		}

		err = Templates["edit_user"].ExecuteTemplate(w, "layout", map[string]any{
			"Title":       "Edit User",
			"CurrentPath": r.URL.Path,
			"User":        user,
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	balance := r.FormValue("balance")
	userType := r.FormValue("user_type")

	err := repository.UpdateUser(id, firstName, lastName, email, balance, userType)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := repository.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
