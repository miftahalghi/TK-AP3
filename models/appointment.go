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
		err := rows.Scan(
			&apt.TanggalKonsultasi,
			&apt.Status,
			&apt.Gejala,
			&apt.Diagnosa,
			&apt.ResepObat,
			&apt.NamaDokter,
		)
		if err != nil {
			return nil, err
		}
		history = append(history, apt)
	}

	return history, nil
}
