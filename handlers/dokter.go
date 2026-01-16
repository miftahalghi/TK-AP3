package handlers

import (
	"html/template"
	"klinik-app/config"
	"klinik-app/middleware"
	"klinik-app/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// DokterDashboard - Dashboard untuk dokter
func DokterDashboard(w http.ResponseWriter, r *http.Request) {
	sess := middleware.GetSession(r)

	// Get appointment hari ini
	appointments, err := models.GetTodayAppointmentsByDoctor(config.DB, sess["UserID"].(int))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Nama":         sess["Nama"],
		"Appointments": appointments,
	}

	tmpl, err := template.ParseFiles("templates/dokter_dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// DokterKonsultasiPage - Form input hasil konsultasi
func DokterKonsultasiPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appointmentID := vars["id"]

	data := map[string]interface{}{
		"AppointmentID": appointmentID,
	}

	tmpl, err := template.ParseFiles("templates/dokter_konsultasi.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// DokterKonsultasiHandler - Proses input hasil konsultasi
func DokterKonsultasiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/dokter/dashboard", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	appointmentID, _ := strconv.Atoi(vars["id"])
	gejala := r.FormValue("gejala")
	diagnosa := r.FormValue("diagnosa")
	resep := r.FormValue("resep")

	// Update appointment dengan hasil konsultasi
	err := models.CompleteConsultation(config.DB, appointmentID, gejala, diagnosa, resep)
	if err != nil {
		http.Error(w, "Gagal simpan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dokter/dashboard", http.StatusSeeOther)
}
