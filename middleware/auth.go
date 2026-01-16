package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("secret-key-klinik-ganti-ini"))

// RequireAuth - Middleware untuk memastikan user sudah login
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "session-klinik")

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}

// RequireRole - Middleware untuk memastikan user punya role tertentu
func RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "session-klinik")

		userRole, ok := session.Values["role"].(string)
		if !ok || userRole != role {
			http.Error(w, "Forbidden - Anda tidak punya akses", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

// GetSession - Helper untuk mendapatkan data session
func GetSession(r *http.Request) map[string]interface{} {
	session, _ := Store.Get(r, "session-klinik")
	return map[string]interface{}{
		"UserID": session.Values["user_id"],
		"Nama":   session.Values["nama"],
		"Role":   session.Values["role"],
	}
}
