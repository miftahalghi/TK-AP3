package handlers

import (
	"html/template"
	"klinik-app/config"
	"klinik-app/middleware"
	"klinik-app/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// LoginPage - Tampilkan halaman login
func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// LoginHandler - Proses login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	nik := r.FormValue("nik")
	password := r.FormValue("password")

	// Get user dari database
	user, err := models.GetUserByNIK(config.DB, nik)
	if err != nil {
		http.Error(w, "NIK atau password salah", http.StatusUnauthorized)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		http.Error(w, "NIK atau password salah", http.StatusUnauthorized)
		return
	}

	// Buat session
	session, _ := middleware.Store.Get(r, "session-klinik")
	session.Values["authenticated"] = true
	session.Values["user_id"] = user.UserID
	session.Values["nama"] = user.Nama
	session.Values["role"] = user.Role
	session.Save(r, w)

	// Redirect sesuai role
	switch user.Role {
	case "pasien":
		http.Redirect(w, r, "/pasien/dashboard", http.StatusSeeOther)
	case "admin":
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	case "dokter":
		http.Redirect(w, r, "/dokter/dashboard", http.StatusSeeOther)
	default:
		http.Error(w, "Role tidak dikenali", http.StatusForbidden)
	}
}

// LogoutHandler - Proses logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "session-klinik")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
