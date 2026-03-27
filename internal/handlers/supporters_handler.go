package handlers

import (
	"net/http"
	"strconv"
	"unibazar/project/internal/models"
	"unibazar/project/internal/repository"

	"github.com/gorilla/mux"
)

type SupportersPageData struct {
	Title           string
	CurrentPath     string
	Supporters      []models.Supporter
	TotalSupporters int
	CurrentPage     int
	TotalPages      int
	PrevPage        int
	NextPage        int
	Pages           []int
}

const supportersPerPage = 5

func SupportersPageHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * supportersPerPage

	supporters, err := repository.GetSupportersPaginated(supportersPerPage, offset)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	totalSupporters, err := repository.CountSupporters()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	totalPages := totalSupporters / supportersPerPage
	if totalSupporters%supportersPerPage != 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}

	pages := make([]int, totalPages)
	for i := range pages {
		pages[i] = i + 1
	}

	data := SupportersPageData{
		Title:           "Supporters",
		CurrentPath:     r.URL.Path,
		Supporters:      supporters,
		TotalSupporters: totalSupporters,
		CurrentPage:     page,
		TotalPages:      totalPages,
		PrevPage:        page - 1,
		NextPage:        page + 1,
		Pages:           pages,
	}

	err = Templates["supporters"].ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func ViewSupporterHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	supporter, err := repository.GetSupporterByID(id)
	if err != nil {
		http.Error(w, "Supporter not found", 404)
		return
	}

	kpi, err := repository.GetSupporterKPIs(id)
	if err != nil {
		kpi = models.SupporterKPI{}
	}

	err = Templates["view_supporter"].ExecuteTemplate(w, "layout", map[string]any{
		"Title":       supporter.FirstName + " " + supporter.LastName,
		"CurrentPath": "/supporters",
		"Supporter":   supporter,
		"KPI":         kpi,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func CreateSupporterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := Templates["create_supporter"].ExecuteTemplate(w, "layout", map[string]any{
			"Title":       "Add Supporter",
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
	imageURL := r.FormValue("image_url")

	err := repository.CreateSupporter(firstName, lastName, email, imageURL)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/supporters", http.StatusSeeOther)
}

func EditSupporterHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if r.Method == "GET" {
		supporter, err := repository.GetSupporterByID(id)
		if err != nil {
			http.Error(w, "Supporter not found", 404)
			return
		}

		err = Templates["edit_supporter"].ExecuteTemplate(w, "layout", map[string]any{
			"Title":       "Edit Supporter",
			"CurrentPath": "/supporters",
			"Supporter":   supporter,
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	imageURL := r.FormValue("image_url")

	err := repository.UpdateSupporter(id, firstName, lastName, email, imageURL)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/supporters", http.StatusSeeOther)
}

func DeleteSupporterHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := repository.DeleteSupporter(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/supporters", http.StatusSeeOther)
}
