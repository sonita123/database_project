package session

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"unibazar/project/internal/models"
	"unibazar/project/internal/repository"
)

const SessionCookie = "unibazar_session"

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Create stores a new session in DB and sets the cookie
func Create(w http.ResponseWriter, userID int, role string) error {
	token := generateToken()
	exp := time.Now().Add(24 * time.Hour)

	if err := repository.CreateSession(token, userID, role, exp); err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Expires:  exp,
	})
	return nil
}

// Get retrieves and validates the session from the request cookie
func Get(r *http.Request) (*models.Session, error) {
	cookie, err := r.Cookie(SessionCookie)
	if err != nil {
		return nil, err
	}

	s, err := repository.GetSession(cookie.Value)
	if err != nil {
		return nil, err
	}

	if s.ExpiresAt.Before(time.Now()) {
		return nil, http.ErrNoCookie
	}

	return s, nil
}

// GetUserID returns the session user ID, or 0 if not logged in
func GetUserID(r *http.Request) int {
	s, err := Get(r)
	if err != nil {
		return 0
	}
	return s.UserID
}

// GetRole returns the session role, or "" if not logged in
func GetRole(r *http.Request) string {
	s, err := Get(r)
	if err != nil {
		return ""
	}
	return s.Role
}

// Destroy deletes the session from DB and clears the cookie
func Destroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(SessionCookie)
	if err == nil {
		repository.DeleteSession(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   SessionCookie,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
