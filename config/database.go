package config

import (
	"database/sql"
	"log"
	"net/url"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	rawURL := os.Getenv("MYSQL_URL")
	if rawURL == "" {
		log.Fatal("MYSQL_URL is not set")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal("Invalid MYSQL_URL:", err)
	}

	user := u.User.Username()
	pass, _ := u.User.Password()
	host := u.Host
	dbName := strings.TrimPrefix(u.Path, "/")

	dsn := user + ":" + pass + "@tcp(" + host + ")/" + dbName + "?parseTime=true"

	log.Println("Connecting to DB host:", host)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	log.Println("✓ Database connected successfully")
}
