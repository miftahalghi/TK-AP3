package handlers

import (
	"database/sql"
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

// AdminReschedulePage - Form reschedule appointment
func AdminReschedulePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appointmentID, _ := strconv.Atoi(vars["id"])

	// Get appointment detail
	var apt models.Appointment
	var namaPasien string
	var namaDokter sql.NullString

	err := config.DB.QueryRow(`
		SELECT 
			a.appointment_id, a.nomor_registrasi,
			a.tanggal_konsultasi, a.waktu_konsultasi,
			a.doctor_id, up.nama AS nama_pasien, 
			ud.nama AS nama_dokter
		FROM appointments a
		JOIN users up ON a.patient_id = up.user_id
		LEFT JOIN users ud ON a.doctor_id = ud.user_id
		WHERE a.appointment_id = ?
	`, appointmentID).Scan(
		&apt.AppointmentID,
		&apt.NomorRegistrasi,
		&apt.TanggalKonsultasi,
		&apt.WaktuKonsultasi,
		&apt.DoctorID,
		&namaPasien,
		&namaDokter,
	)

	if err != nil {
		http.Error(w, "Appointment tidak ditemukan: "+err.Error(), http.StatusNotFound)
		return
	}

	apt.NamaPasien = namaPasien
	if namaDokter.Valid {
		apt.NamaDokter = namaDokter.String
	} else {
		apt.NamaDokter = "Belum ditentukan"
	}

	// Get list dokter
	doctors, err := models.GetDoctors(config.DB)
	if err != nil {
		http.Error(w, "Gagal get doctors: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Appointment": apt,
		"Doctors":     doctors,
	}

	tmpl, err := template.ParseFiles("templates/admin_reschedule.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Execute error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// AdminRescheduleHandler - Proses reschedule
func AdminRescheduleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	appointmentID, _ := strconv.Atoi(vars["id"])
	doctorID, _ := strconv.Atoi(r.FormValue("doctor_id"))
	tanggal := r.FormValue("tanggal")
	waktu := r.FormValue("waktu")

	err := models.RescheduleAppointment(config.DB, appointmentID, doctorID, tanggal, waktu)
	if err != nil {
		http.Error(w, "Gagal reschedule: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

// AdminCancelAppointment - Admin cancel appointment
func AdminCancelAppointment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	appointmentID, _ := strconv.Atoi(r.FormValue("appointment_id"))

	err := models.CancelAppointment(config.DB, appointmentID)
	if err != nil {
		http.Error(w, "Gagal cancel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}
