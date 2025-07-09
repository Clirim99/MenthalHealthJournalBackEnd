package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"menthalhealthjournal/db"
	"menthalhealthjournal/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate
	if user.FirstName == "" || user.LastName == "" || user.Username == "" ||
		user.Email == "" || user.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	user.Email = strings.ToLower(user.Email)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert
	query := `
		INSERT INTO users (first_name, last_name, username, email, password)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err = db.DB.QueryRow(query, user.FirstName, user.LastName, user.Username, user.Email, string(hashedPassword)).
		Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create user: %v", err), http.StatusInternalServerError)
		return
	}

	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
