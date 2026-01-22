package models

import (
	"database/sql"
	"time"
)

type Appointment struct {
	AppointmentID     int            `json:"appointment_id"`
	NomorRegistrasi   string         `json:"nomor_registrasi"`
	PatientID         int            `json:"patient_id"`
	DoctorID          sql.NullInt64  `json:"doctor_id"`
	TanggalKonsultasi time.Time      `json:"tanggal_konsultasi"`
	WaktuKonsultasi   sql.NullString `json:"waktu_konsultasi"`
	Status            string         `json:"status"`
	Gejala            sql.NullString `json:"gejala"`
	Diagnosa          sql.NullString `json:"diagnosa"`
	ResepObat         sql.NullString `json:"resep_obat"`
	CreatedAt         time.Time      `json:"created_at"`

	// Join fields
	NamaPasien string `json:"nama_pasien,omitempty"`
	NamaDokter string `json:"nama_dokter,omitempty"`
}

// CreateAppointment - Pasien booking konsultasi
func CreateAppointment(db *sql.DB, nomorReg string, patientID int, tanggal string) error {
	query := `INSERT INTO appointments (nomor_registrasi, patient_id, tanggal_konsultasi, status) 
	          VALUES (?, ?, ?, 'pending')`

	_, err := db.Exec(query, nomorReg, patientID, tanggal)
	return err
}

// GetPendingAppointments - Admin melihat pending appointments
func GetPendingAppointments(db *sql.DB) ([]Appointment, error) {
	query := `
		SELECT 
			a.appointment_id, a.nomor_registrasi, 
			a.tanggal_konsultasi, a.created_at,
			u.nama AS nama_pasien
		FROM appointments a
		JOIN users u ON a.patient_id = u.user_id
		WHERE a.status = 'pending'
		ORDER BY a.created_at ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var apt Appointment
		err := rows.Scan(
			&apt.AppointmentID,
			&apt.NomorRegistrasi,
			&apt.TanggalKonsultasi,
			&apt.CreatedAt,
			&apt.NamaPasien,
		)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, apt)
	}

	return appointments, nil
}

// ApproveAppointment - Admin approve dan assign dokter
func ApproveAppointment(db *sql.DB, appointmentID, doctorID int, waktu string) error {
	query := `UPDATE appointments 
	          SET doctor_id = ?, waktu_konsultasi = ?, status = 'approved' 
	          WHERE appointment_id = ?`

	_, err := db.Exec(query, doctorID, waktu, appointmentID)
	return err
}

// GetTodayAppointmentsByDoctor - Dokter melihat appointment hari ini
func GetTodayAppointmentsByDoctor(db *sql.DB, doctorID int) ([]Appointment, error) {
	query := `
		SELECT 
			a.appointment_id, a.nomor_registrasi, 
			a.waktu_konsultasi, u.nama AS nama_pasien
		FROM appointments a
		JOIN users u ON a.patient_id = u.user_id
		WHERE a.doctor_id = ? 
		  AND DATE(a.tanggal_konsultasi) = CURDATE() 
		  AND a.status = 'approved'
	`

	rows, err := db.Query(query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var apt Appointment
		err := rows.Scan(
			&apt.AppointmentID,
			&apt.NomorRegistrasi,
			&apt.WaktuKonsultasi,
			&apt.NamaPasien,
		)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, apt)
	}

	return appointments, nil
}

// CompleteConsultation - Dokter input hasil konsultasi
func CompleteConsultation(db *sql.DB, appointmentID int, gejala, diagnosa, resep string) error {
	query := `UPDATE appointments 
	          SET gejala = ?, diagnosa = ?, resep_obat = ?, status = 'completed' 
	          WHERE appointment_id = ?`

	_, err := db.Exec(query, gejala, diagnosa, resep, appointmentID)
	return err
}

// GetPatientActiveAppointments - Pasien melihat appointment aktif (pending & approved)
func GetPatientActiveAppointments(db *sql.DB, patientID int) ([]Appointment, error) {
	query := `
		SELECT 
			a.appointment_id, a.nomor_registrasi,
			a.tanggal_konsultasi, a.waktu_konsultasi,
			a.status, u.nama AS nama_dokter
		FROM appointments a
		LEFT JOIN users u ON a.doctor_id = u.user_id
		WHERE a.patient_id = ? 
		  AND a.status IN ('pending', 'approved')
		ORDER BY a.tanggal_konsultasi ASC
	`

	rows, err := db.Query(query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var apt Appointment
		var namaDokter sql.NullString

		err := rows.Scan(
			&apt.AppointmentID,
			&apt.NomorRegistrasi,
			&apt.TanggalKonsultasi,
			&apt.WaktuKonsultasi,
			&apt.Status,
			&namaDokter,
		)
		if err != nil {
			return nil, err
		}

		if namaDokter.Valid {
			apt.NamaDokter = namaDokter.String
		} else {
			apt.NamaDokter = "Belum ditentukan"
		}

		appointments = append(appointments, apt)
	}

	return appointments, nil
}

// CancelAppointment - Cancel appointment (update status jadi cancelled)
func CancelAppointment(db *sql.DB, appointmentID int) error {
	query := `UPDATE appointments SET status = 'cancelled' WHERE appointment_id = ?`
	_, err := db.Exec(query, appointmentID)
	return err
}

// RescheduleAppointment - Admin ubah jadwal appointment
func RescheduleAppointment(db *sql.DB, appointmentID, doctorID int, tanggal, waktu string) error {
	query := `UPDATE appointments 
	          SET doctor_id = ?, tanggal_konsultasi = ?, waktu_konsultasi = ? 
	          WHERE appointment_id = ?`

	_, err := db.Exec(query, doctorID, tanggal, waktu, appointmentID)
	return err
}

// GetAppointmentByID - Get detail appointment
func GetAppointmentByID(db *sql.DB, appointmentID int) (*Appointment, error) {
	var apt Appointment
	query := `
		SELECT 
			a.appointment_id, a.nomor_registrasi, a.patient_id,
			a.doctor_id, a.tanggal_konsultasi, a.waktu_konsultasi,
			a.status, u.nama AS nama_pasien
		FROM appointments a
		JOIN users u ON a.patient_id = u.user_id
		WHERE a.appointment_id = ?
	`

	var namaPasien string
	err := db.QueryRow(query, appointmentID).Scan(
		&apt.AppointmentID,
		&apt.NomorRegistrasi,
		&apt.PatientID,
		&apt.DoctorID,
		&apt.TanggalKonsultasi,
		&apt.WaktuKonsultasi,
		&apt.Status,
		&namaPasien,
	)

	if err != nil {
		return nil, err
	}

	apt.NamaPasien = namaPasien

	return &apt, nil
}

// GetPatientHistory - Pasien melihat riwayat konsultasi
func GetPatientHistory(db *sql.DB, patientID int) ([]Appointment, error) {
	query := `
		SELECT 
			a.tanggal_konsultasi, a.status, 
			a.gejala, a.diagnosa, a.resep_obat,
			u.nama AS nama_dokter
		FROM appointments a
		LEFT JOIN users u ON a.doctor_id = u.user_id
		WHERE a.patient_id = ?
		ORDER BY a.tanggal_konsultasi DESC
	`

	rows, err := db.Query(query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []Appointment
	for rows.Next() {
		var apt Appointment
		var namaDokter sql.NullString // ← UBAH: Gunakan sql.NullString untuk handle NULL

		err := rows.Scan(
			&apt.TanggalKonsultasi,
			&apt.Status,
			&apt.Gejala,
			&apt.Diagnosa,
			&apt.ResepObat,
			&namaDokter, // ← Scan ke variable temporary
		)
		if err != nil {
			return nil, err
		}

		// Convert sql.NullString ke string biasa
		if namaDokter.Valid {
			apt.NamaDokter = namaDokter.String
		} else {
			apt.NamaDokter = "Belum ditentukan" // Default value jika NULL
		}

		history = append(history, apt)
	}

	return history, nil
}

// GetAllAppointments - Admin melihat SEMUA appointments dengan berbagai status
// GetAllAppointments - Admin melihat SEMUA appointments dengan berbagai status
func GetAllAppointments(db *sql.DB) ([]Appointment, error) {
	query := `
		SELECT 
			a.appointment_id, 
			a.nomor_registrasi, 
			a.tanggal_konsultasi, 
			a.waktu_konsultasi,
			a.status, 
			a.created_at,
			up.nama AS nama_pasien,
			COALESCE(ud.nama, '') AS nama_dokter
		FROM appointments a
		JOIN users up ON a.patient_id = up.user_id
		LEFT JOIN users ud ON a.doctor_id = ud.user_id
		ORDER BY 
			CASE a.status
				WHEN 'pending' THEN 1
				WHEN 'approved' THEN 2
				WHEN 'completed' THEN 3
				WHEN 'cancelled' THEN 4
			END,
			a.tanggal_konsultasi ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var apt Appointment
		var namaDokter sql.NullString // ← Tetap perlu ini

		err := rows.Scan(
			&apt.AppointmentID,
			&apt.NomorRegistrasi,
			&apt.TanggalKonsultasi,
			&apt.WaktuKonsultasi,
			&apt.Status,
			&apt.CreatedAt,
			&apt.NamaPasien,
			&namaDokter, // ← Scan ke variable ini dulu
		)
		if err != nil {
			return nil, err
		}

		if namaDokter.Valid {
			apt.NamaDokter = namaDokter.String
		} else {
			apt.NamaDokter = "Belum ditentukan"
		}

		appointments = append(appointments, apt)
	}

	return appointments, nil
}
