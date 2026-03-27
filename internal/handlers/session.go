package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const SessionCookie = "unibazar_session"

func SetSession(w http.ResponseWriter, id int, role string) {

	value := strconv.Itoa(id) + ":" + role

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}

func GetSessionRole(r *http.Request) string {

	cookie, err := r.Cookie(SessionCookie)
	if err != nil {
		return ""
	}

	parts := strings.Split(cookie.Value, ":")

	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}

func ClearSession(w http.ResponseWriter) {

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
func GetSessionID(r *http.Request) int {

	cookie, err := r.Cookie(SessionCookie)
	if err != nil {
		return 0
	}

	parts := strings.Split(cookie.Value, ":")

	if len(parts) != 2 {
		return 0
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}

	return id
}
