package handlers

import (
	"net/http"
	"strconv"

	"unibazar/project/internal/repository"
	"unibazar/project/internal/session"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// ── User Login ────────────────────────────────────────────────

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	if session.GetRole(r) == "user" {
		http.Redirect(w, r, "/portal", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		Templates["login_user"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Sign In",
			"Error": "",
			"Next":  r.URL.Query().Get("next"),
		})
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := repository.GetUserByEmail(email)
	if err != nil {
		Templates["login_user"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Sign In", "Error": "Email not found", "Next": r.FormValue("next"),
		})
		return
	}
	if user.Password == "" {
		Templates["login_user"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Sign In", "Error": "No password set — contact admin", "Next": r.FormValue("next"),
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		Templates["login_user"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Sign In", "Error": "Incorrect password", "Next": r.FormValue("next"),
		})
		return
	}

	session.Create(w, user.UserID, "user")

	next := r.FormValue("next")
	if next == "" {
		next = "/portal"
	}
	http.Redirect(w, r, next, http.StatusSeeOther)
}

// ── Supporter Login ───────────────────────────────────────────

func SupporterLoginHandler(w http.ResponseWriter, r *http.Request) {
	if session.GetRole(r) == "supporter" {
		http.Redirect(w, r, "/supporter/portal", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		Templates["login_supporter"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Supporter Login", "Error": "",
		})
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	supporter, err := repository.GetSupporterByEmail(email)
	if err != nil {
		Templates["login_supporter"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Supporter Login", "Error": "Email not found",
		})
		return
	}
	if supporter.Password == "" {
		Templates["login_supporter"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Supporter Login", "Error": "No password set — contact admin",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(supporter.Password), []byte(password)); err != nil {
		Templates["login_supporter"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Supporter Login", "Error": "Incorrect password",
		})
		return
	}

	session.Create(w, supporter.SupporterID, "supporter")
	http.Redirect(w, r, "/supporter/portal", http.StatusSeeOther)
}

// ── Admin Login ───────────────────────────────────────────────

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if session.GetRole(r) == "admin" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		Templates["login_admin"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Admin Login", "Error": "",
		})
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	admin, err := repository.GetAdminByUsername(username)
	if err != nil {
		Templates["login_admin"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Admin Login", "Error": "Username not found",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		Templates["login_admin"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Admin Login", "Error": "Incorrect password",
		})
		return
	}

	session.Create(w, admin.AdminID, "admin")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ── Admin Register ────────────────────────────────────────────

func AdminRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if session.GetRole(r) == "admin" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		Templates["register_admin"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Create Admin", "Error": "",
		})
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	confirm := r.FormValue("password_confirm")

	if username == "" || password == "" {
		Templates["register_admin"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Create Admin", "Error": "Username and password are required",
		})
		return
	}
	if password != confirm {
		Templates["register_admin"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Create Admin", "Error": "Passwords do not match",
		})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := repository.CreateAdmin(username, string(hashed)); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

// ── Logout ────────────────────────────────────────────────────

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session.Destroy(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ── Set User Password ─────────────────────────────────────────

func SetUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if r.Method == "GET" {
		user, err := repository.GetUserByID(id)
		if err != nil {
			http.Error(w, "User not found", 404)
			return
		}
		Templates["set_user_password"].ExecuteTemplate(w, "layout", map[string]any{
			"Title": "Set Password", "CurrentPath": r.URL.Path, "User": user,
			"AdminUsername": AdminUsername(r),
		})
		return
	}

	password := r.FormValue("password")
	confirm := r.FormValue("password_confirm")
	if password == "" {
		http.Error(w, "Password cannot be empty", 400)
		return
	}
	if password != confirm {
		http.Error(w, "Passwords do not match", 400)
		return
	}

	userID, _ := strconv.Atoi(id)
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := repository.SetUserPassword(userID, string(hashed)); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// ── Set Supporter Password ────────────────────────────────────

func SetSupporterPasswordHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if r.Method == "GET" {
		supporter, err := repository.GetSupporterByID(id)
		if err != nil {
			http.Error(w, "Supporter not found", 404)
			return
		}
		Templates["set_supporter_password"].ExecuteTemplate(w, "layout", map[string]any{
			"Title": "Set Password", "CurrentPath": r.URL.Path, "Supporter": supporter,
			"AdminUsername": AdminUsername(r),
		})
		return
	}

	password := r.FormValue("password")
	confirm := r.FormValue("password_confirm")
	if password == "" {
		http.Error(w, "Password cannot be empty", 400)
		return
	}
	if password != confirm {
		http.Error(w, "Passwords do not match", 400)
		return
	}

	supporterID, _ := strconv.Atoi(id)
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := repository.SetSupporterPassword(supporterID, string(hashed)); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/supporters", http.StatusSeeOther)
}

// ── AdminUsername helper for layout data ──────────────────────

func AdminUsername(r *http.Request) string {
	id := session.GetUserID(r)
	if id == 0 {
		return "Admin"
	}
	admin, err := repository.GetAdminByID(id)
	if err != nil {
		return "Admin"
	}
	return admin.Username
}

func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Already logged in — send to portal
	if session.GetRole(r) == "user" {
		http.Redirect(w, r, "/portal", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		Templates["register_user"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Create Account",
			"Error": "",
		})
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm := r.FormValue("password_confirm")

	// Validate
	if firstName == "" || lastName == "" || email == "" || password == "" {
		Templates["register_user"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Create Account", "Error": "All fields are required",
			"FirstName": firstName, "LastName": lastName, "Email": email,
		})
		return
	}
	if len(password) < 6 {
		Templates["register_user"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Create Account", "Error": "Password must be at least 6 characters",
			"FirstName": firstName, "LastName": lastName, "Email": email,
		})
		return
	}
	if password != confirm {
		Templates["register_user"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Create Account", "Error": "Passwords do not match",
			"FirstName": firstName, "LastName": lastName, "Email": email,
		})
		return
	}

	// Check email not taken
	exists, err := repository.EmailExists(email)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if exists {
		Templates["register_user"].ExecuteTemplate(w, "auth", map[string]any{
			"Title": "Create Account", "Error": "An account with this email already exists",
			"FirstName": firstName, "LastName": lastName, "Email": email,
		})
		return
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Create user
	if err := repository.RegisterUser(firstName, lastName, email, string(hashed)); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Auto-login after registration
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		// Registration succeeded but login failed — redirect to login
		http.Redirect(w, r, "/login?registered=1", http.StatusSeeOther)
		return
	}

	session.Create(w, user.UserID, "user")
	http.Redirect(w, r, "/portal", http.StatusSeeOther)
}
