package handlers

import (
	"net/http"
	"time"
	"unibazar/project/internal/models"
	"unibazar/project/internal/repository"
	"unibazar/project/internal/session"
)

type DashboardData struct {
	Title           string
	CurrentPath     string
	AdminUsername   string
	Today           string
	TotalUsers      int
	TotalSupporters int
	TotalRevenue    float64
	OpenRequests    int
	RecentUsers     []models.User
}

func Dashboard(w http.ResponseWriter, r *http.Request) {

	adminID := session.GetUserID(r)

	adminUsername := "Admin"

	if adminID != 0 {
		admin, err := repository.GetAdminByID(adminID)
		if err == nil {
			adminUsername = admin.Username
		}
	}

	totalUsers, err := repository.CountUsers()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	totalSupporters, err := repository.CountSupporters()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	totalRevenue, err := repository.CalculateTotalRevenue()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	openRequests, err := repository.CountOpenRequests()
	if err != nil {
		openRequests = 0
	}

	recentUsers, err := repository.GetRecentUsers(5)
	if err != nil {
		recentUsers = []models.User{}
	}

	data := DashboardData{
		Title:           "Dashboard",
		CurrentPath:     r.URL.Path,
		AdminUsername:   adminUsername,
		Today:           time.Now().Format("02 Jan 2006"),
		TotalUsers:      totalUsers,
		TotalSupporters: totalSupporters,
		TotalRevenue:    totalRevenue,
		OpenRequests:    openRequests,
		RecentUsers:     recentUsers,
	}

	err = Templates["dashboard"].ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
