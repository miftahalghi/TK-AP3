package models

import (
	"database/sql"
	"time"
)

type User struct {
	UserID    int       `json:"user_id"`
	NIK       string    `json:"nik"`
	Nama      string    `json:"nama"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// GetUserByNIK - Mendapatkan user berdasarkan NIK
func GetUserByNIK(db *sql.DB, nik string) (*User, error) {
	var user User

	query := `SELECT user_id, nik, nama, password, role, created_at 
	          FROM users WHERE nik = ?`

	err := db.QueryRow(query, nik).Scan(
		&user.UserID,
		&user.NIK,
		&user.Nama,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetDoctors - Mendapatkan semua dokter (untuk dropdown admin)
func GetDoctors(db *sql.DB) ([]User, error) {
	query := `SELECT user_id, nama FROM users WHERE role = 'dokter'`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []User
	for rows.Next() {
		var doctor User
		err := rows.Scan(&doctor.UserID, &doctor.Nama)
		if err != nil {
			return nil, err
		}
		doctors = append(doctors, doctor)
	}

	return doctors, nil
}
