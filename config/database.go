package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("MYSQLUSER"),
		os.Getenv("MYSQLPASSWORD"),
		os.Getenv("MYSQLHOST"),
		os.Getenv("MYSQLPORT"),
		os.Getenv("MYSQLDATABASE"),
	)

	log.Println("Connecting to DB host:", os.Getenv("MYSQLHOST"))

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	log.Println("✓ Database connected successfully")
}
