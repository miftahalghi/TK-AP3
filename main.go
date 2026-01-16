package main

import (
	"klinik-app/config"
	"klinik-app/handlers"
	"klinik-app/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	config.InitDB()
	defer config.DB.Close()

	// Setup router
	r := mux.NewRouter()

	// Auth routes (public)
	r.HandleFunc("/", handlers.LoginPage).Methods("GET")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("GET")

	// Pasien routes (protected)
	r.HandleFunc("/pasien/dashboard",
		middleware.RequireAuth(
			middleware.RequireRole("pasien", handlers.PasienDashboard),
		),
	).Methods("GET")

	r.HandleFunc("/pasien/booking",
		middleware.RequireAuth(
			middleware.RequireRole("pasien", handlers.PasienBookingPage),
		),
	).Methods("GET")

	r.HandleFunc("/pasien/booking",
		middleware.RequireAuth(
			middleware.RequireRole("pasien", handlers.PasienBookingHandler),
		),
	).Methods("POST")

	r.HandleFunc("/pasien/riwayat",
		middleware.RequireAuth(
			middleware.RequireRole("pasien", handlers.PasienRiwayat),
		),
	).Methods("GET")

	// Admin routes (protected)
	r.HandleFunc("/admin/dashboard",
		middleware.RequireAuth(
			middleware.RequireRole("admin", handlers.AdminDashboard),
		),
	).Methods("GET")

	r.HandleFunc("/admin/approve/{id}",
		middleware.RequireAuth(
			middleware.RequireRole("admin", handlers.AdminApprovePage),
		),
	).Methods("GET")

	r.HandleFunc("/admin/approve/{id}",
		middleware.RequireAuth(
			middleware.RequireRole("admin", handlers.AdminApproveHandler),
		),
	).Methods("POST")

	// Dokter routes (protected)
	r.HandleFunc("/dokter/dashboard",
		middleware.RequireAuth(
			middleware.RequireRole("dokter", handlers.DokterDashboard),
		),
	).Methods("GET")

	r.HandleFunc("/dokter/konsultasi/{id}",
		middleware.RequireAuth(
			middleware.RequireRole("dokter", handlers.DokterKonsultasiPage),
		),
	).Methods("GET")

	r.HandleFunc("/dokter/konsultasi/{id}",
		middleware.RequireAuth(
			middleware.RequireRole("dokter", handlers.DokterKonsultasiHandler),
		),
	).Methods("POST")

	// Start server
	log.Println("========================================")
	log.Println("🏥 Sistem Informasi Klinik")
	log.Println("========================================")
	log.Println("Server running on: http://localhost:8080")
	log.Println("========================================")
	log.Fatal(http.ListenAndServe(":8080", r))
}
