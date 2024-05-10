package main

import (
	"database/sql"
	
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// DSN string anatomy
	// Linux/Mac: postgres://username:password@host:port/database?sslmode=disable
	// Windows: postgresql://username:password@host:port/database?sslmode=disable

	// Uncomment the following code to check if the database connection is working
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to DB!")
}
