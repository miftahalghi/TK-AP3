package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "password123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	fmt.Println("=== HASH UNTUK PASSWORD: password123 ===")
	fmt.Println(string(hash))
	fmt.Println("\n=== COPY QUERY INI DAN JALANKAN DI DATABASE ===")
	fmt.Printf("UPDATE users SET password = '%s' WHERE nik = '1111111111111111';\n", string(hash))
}
