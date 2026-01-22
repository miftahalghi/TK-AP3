package handlers

import (
	"html/template"
	"klinik-app/config"
	"klinik-app/middleware"
	"klinik-app/models"
	"log"
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
		// Debug: log jika user tidak ditemukan
		log.Printf("❌ User not found for NIK: %s", nik)
		http.Error(w, "NIK atau password salah", http.StatusUnauthorized)
		return
	}

	// Debug: log hash comparison
	log.Printf("🔍 Login attempt for: %s (Role: %s)", user.Nama, user.Role)
	log.Printf("🔍 Password hash in DB: %s", user.Password)

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("❌ Password verification failed: %v", err)
		http.Error(w, "NIK atau password salah", http.StatusUnauthorized)
		return
	}

	log.Printf("✅ Login successful: %s", user.Nama)

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

// RegisterPage - Tampilkan halaman registrasi
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// RegisterHandler - Proses registrasi pasien baru
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	nik := r.FormValue("nik")
	nama := r.FormValue("nama")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	// Validasi input
	if len(nik) != 16 {
		http.Error(w, "NIK harus 16 digit", http.StatusBadRequest)
		return
	}

	if password != confirmPassword {
		http.Error(w, "Password dan konfirmasi password tidak sama", http.StatusBadRequest)
		return
	}

	if len(password) < 6 {
		http.Error(w, "Password minimal 6 karakter", http.StatusBadRequest)
		return
	}

	// Cek NIK sudah terdaftar atau belum
	existingUser, _ := models.GetUserByNIK(config.DB, nik)
	if existingUser != nil {
		http.Error(w, "NIK sudah terdaftar. Silakan login.", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Gagal memproses password", http.StatusInternalServerError)
		return
	}

	// Insert user baru dengan role pasien
	_, err = config.DB.Exec(
		"INSERT INTO users (nik, nama, password, role) VALUES (?, ?, ?, 'pasien')",
		nik, nama, string(hashedPassword),
	)

	if err != nil {
		log.Printf("❌ Registration failed: %v", err)
		http.Error(w, "Gagal registrasi. Silakan coba lagi.", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ New user registered: %s (NIK: %s)", nama, nik)

	// Auto login setelah registrasi
	user, _ := models.GetUserByNIK(config.DB, nik)
	session, _ := middleware.Store.Get(r, "session-klinik")
	session.Values["authenticated"] = true
	session.Values["user_id"] = user.UserID
	session.Values["nama"] = user.Nama
	session.Values["role"] = user.Role
	session.Save(r, w)

	http.Redirect(w, r, "/pasien/dashboard", http.StatusSeeOther)
}
