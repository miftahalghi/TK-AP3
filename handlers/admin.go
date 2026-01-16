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

// AdminDashboard - Dashboard untuk admin
func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	sess := middleware.GetSession(r)

	// Get pending appointments
	appointments, err := models.GetPendingAppointments(config.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Nama":         sess["Nama"],
		"Appointments": appointments,
	}

	tmpl, err := template.ParseFiles("templates/admin_dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// AdminApprovePage - Form approve appointment
func AdminApprovePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appointmentID := vars["id"]

	// Get list dokter
	doctors, err := models.GetDoctors(config.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"AppointmentID": appointmentID,
		"Doctors":       doctors,
	}

	tmpl, err := template.ParseFiles("templates/admin_approve.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// AdminApproveHandler - Proses approve appointment
func AdminApproveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	appointmentID, _ := strconv.Atoi(vars["id"])
	doctorID, _ := strconv.Atoi(r.FormValue("doctor_id"))
	waktu := r.FormValue("waktu")

	// Update appointment
	err := models.ApproveAppointment(config.DB, appointmentID, doctorID, waktu)
	if err != nil {
		http.Error(w, "Gagal approve: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}
