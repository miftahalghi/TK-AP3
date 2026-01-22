package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error

	dsn := os.Getenv("MYSQL_URL")
	if dsn == "" {
		log.Fatal("MYSQL_URL is not set")
	}

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	log.Println("✓ Database connected successfully")
}
