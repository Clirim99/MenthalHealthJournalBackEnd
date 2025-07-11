package repositories

import (
	"database/sql"
	"fmt"
	"menthalhealthjournal/db"
	"menthalhealthjournal/models"
//	"menthalhealthjournal/services"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func GetUserByEmail(email string) (models.User, error) {
	fmt.Println("Looking up user by email:", email)
	query := `
		SELECT * FROM users WHERE email = $1
	`
	var user models.User
	err := db.DB.QueryRow(query, email).Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found")
		}
		return models.User{}, fmt.Errorf("error querying user: %v", err)
	}

	return user, nil
}

func CreateUser(c *gin.Context, user models.User) {
	query := `
        INSERT INTO users (first_name, last_name, username, email, password)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at
    `

	err := db.DB.QueryRow(query, user.FirstName, user.LastName, user.Username, user.Email, user.Password).
		Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		log.Println("DB insert error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user: " + err.Error()})
		return
	}

	// Clear password before responding
	user.Password = ""

	c.JSON(http.StatusCreated, user)
}

func AuthenticateUser(email, password string) (models.User, error) {
	user, err := GetUserByEmailForLogin(email)
	if err != nil {
		return user, errors.New("invalid email ")
	}

	// Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, errors.New("invalid  or password")
	}

	return user, nil
}

func GetUserByEmailForLogin(email string) (models.User, error) {
	var user models.User

	query := `SELECT id, first_name, last_name, username, email, password, created_at FROM users WHERE email = $1`
	row := db.DB.QueryRow(query, email)

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("user not found")
		}
		return user, err
	}

	return user, nil
}