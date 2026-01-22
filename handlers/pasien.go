package handlers

import (
	"fmt"
	"html/template"
	"klinik-app/config"
	"klinik-app/middleware"
	"klinik-app/models"
	"net/http"
	"strconv"
	"time"
)

// PasienDashboard - Dashboard untuk pasien
func PasienDashboard(w http.ResponseWriter, r *http.Request) {
	sess := middleware.GetSession(r)

	// Get active appointments
	appointments, err := models.GetPatientActiveAppointments(config.DB, sess["UserID"].(int))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Nama":         sess["Nama"],
		"Appointments": appointments,
	}

	tmpl, err := template.ParseFiles("templates/pasien_dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// PasienCancelAppointment - Pasien cancel appointment
func PasienCancelAppointment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/pasien/dashboard", http.StatusSeeOther)
		return
	}

	appointmentID, _ := strconv.Atoi(r.FormValue("appointment_id"))

	err := models.CancelAppointment(config.DB, appointmentID)
	if err != nil {
		http.Error(w, "Gagal cancel appointment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/pasien/dashboard", http.StatusSeeOther)
}

// PasienBookingPage - Tampilkan form booking
func PasienBookingPage(w http.ResponseWriter, r *http.Request) {
	sess := middleware.GetSession(r)

	tmpl, err := template.ParseFiles("templates/pasien_booking.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, sess)
}

// PasienBookingHandler - Proses booking konsultasi
func PasienBookingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/pasien/booking", http.StatusSeeOther)
		return
	}

	sess := middleware.GetSession(r)
	tanggal := r.FormValue("tanggal")

	// Generate nomor registrasi
	nomorReg := fmt.Sprintf("REG-%d-%s", sess["UserID"], time.Now().Format("20060102150405"))

	// Simpan ke database
	err := models.CreateAppointment(config.DB, nomorReg, sess["UserID"].(int), tanggal)
	if err != nil {
		http.Error(w, "Gagal booking: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Tampilkan halaman sukses
	data := map[string]interface{}{
		"NomorReg": nomorReg,
		"Tanggal":  tanggal,
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head><title>Booking Berhasil</title></head>
<body>
	<h2>Booking Berhasil!</h2>
	<p>Nomor Registrasi: <strong>{{.NomorReg}}</strong></p>
	<p>Tanggal Konsultasi: <strong>{{.Tanggal}}</strong></p>
	<p>Status: <strong>Menunggu Persetujuan Admin</strong></p>
	<br>
	<a href="/pasien/dashboard">Kembali ke Dashboard</a>
</body>
</html>`

	t := template.Must(template.New("success").Parse(tmpl))
	t.Execute(w, data)
}

// PasienRiwayat - Tampilkan riwayat konsultasi
func PasienRiwayat(w http.ResponseWriter, r *http.Request) {
	sess := middleware.GetSession(r)

	// Get riwayat dari database
	history, err := models.GetPatientHistory(config.DB, sess["UserID"].(int))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Nama":    sess["Nama"],
		"History": history,
	}

	tmpl, err := template.ParseFiles("templates/pasien_riwayat.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
