package main

import (
	"database/sql"
	"fmt"
	"log"

	//_ "github.com/lib/pq"
)

func main() {
	// Update these with your PostgreSQL credentials
	const (
		host     = "localhost"
		port     = 5432
		user     = "your_db_user"
		password = "your_db_password"
		dbname   = "your_db_name"
	)

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// Open connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
}
