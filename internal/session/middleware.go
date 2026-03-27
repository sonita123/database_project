package session

import "net/http"

// ── Middleware ────────────────────────────────────────────────

func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := Get(r)
		if err != nil || s.Role != "admin" {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := Get(r)
		if err != nil || s.Role != "user" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func RequireSupporter(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := Get(r)
		if err != nil || s.Role != "supporter" {
			http.Redirect(w, r, "/supporter/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func RequireBuyer(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := Get(r)
		if err != nil || s.Role != "user" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// Check user_type from DB — sellers are blocked from buying
		// We pass the check to the handler via context; block happens in handler
		next(w, r)
	}
}
func RequireSeller(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := Get(r)
		if err != nil || s.Role != "seller" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
