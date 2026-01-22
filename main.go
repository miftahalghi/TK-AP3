package main

import (
	"klinik-app/config"
	"klinik-app/handlers"
	"klinik-app/middleware"
	"log"
	"net/http"
	"os"

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
	r.HandleFunc("/register", handlers.RegisterPage).Methods("GET")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")

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

	r.HandleFunc("/pasien/cancel-appointment",
		middleware.RequireAuth(
			middleware.RequireRole("pasien", handlers.PasienCancelAppointment),
		),
	).Methods("POST")

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

	r.HandleFunc("/admin/reschedule/{id}",
		middleware.RequireAuth(
			middleware.RequireRole("admin", handlers.AdminReschedulePage),
		),
	).Methods("GET")

	r.HandleFunc("/admin/reschedule/{id}",
		middleware.RequireAuth(
			middleware.RequireRole("admin", handlers.AdminRescheduleHandler),
		),
	).Methods("POST")

	r.HandleFunc("/admin/cancel-appointment",
		middleware.RequireAuth(
			middleware.RequireRole("admin", handlers.AdminCancelAppointment),
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback untuk local
	}

	log.Println("========================================")
	log.Println("🏥 Sistem Informasi Klinik")
	log.Println("========================================")
	log.Println("Server running on port:", port)
	log.Println("========================================")

	log.Fatal(http.ListenAndServe(":"+port, r))

}
