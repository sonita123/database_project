package handlers

import (
	"net/http"
	"unibazar/project/internal/repository"

	"github.com/gorilla/mux"
)

type UserPageData struct {
	Title          string
	User           interface{}
	Requests       interface{}
	ViewedProducts interface{}
	TopProducts    interface{}
	CurrentPath    string
}

func UserPageHandler(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	user, _ := repository.GetUserByID(id)

	requests, _ := repository.GetLastRequests(id)

	views, _ := repository.GetLastViewedProducts(id)

	topProducts, _ := repository.GetTopSellingProducts()

	data := UserPageData{
		Title:          "User Profile",
		User:           user,
		Requests:       requests,
		ViewedProducts: views,
		TopProducts:    topProducts,
		CurrentPath:    r.URL.Path,
	}

	//Templates.ExecuteTemplate(w, "layout.gohtml", data)
	_ = Templates["user"].ExecuteTemplate(w, "layout", data)

}
