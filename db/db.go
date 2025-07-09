// db/db.go
package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB // Global DB instance

func ConnectDatabase() {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "root"
		dbname   = "MenthalHealthCare"
	)

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("✅ Connected to PostgreSQL database!")
}

func CreateUsersTable() error {
    createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

    _, err := DB.Exec(createTableSQL)
    return err
}

func InsertUser(username, email string) error {
    insertSQL := `INSERT INTO users (username, email) VALUES ($1, $2)`
    _, err := DB.Exec(insertSQL, username, email)
    return err
}
