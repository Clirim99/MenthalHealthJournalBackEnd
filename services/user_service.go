package services

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"errors"


	//"menthalhealthjournal/db"
	"menthalhealthjournal/repositories"
	"menthalhealthjournal/models"
	//services "menthalhealthjournal/repositories"


	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	var user models.User

	// Bind incoming JSON to user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error decoding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Validate required fields
	if user.FirstName == "" || user.LastName == "" || user.Username == "" ||
		user.Email == "" || user.Password == "" {
		log.Println("Validation failed: missing fields")
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	user.Email = strings.ToLower(user.Email)

	// Hash the password

	tempUser, err := repositories.GetUserByEmail(user.Email)
	if err == nil {
		// userID found => user exists
		fmt.Printf("User exists with ID: %s\n", tempUser.Email)
		// Handle "email already taken" case
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Password hashing failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)

	// Create user in DB
	repositories.CreateUser(c, user)
}

func LoginUser(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	email := strings.ToLower(loginData.Email)

	user, err := AuthenticateUser(email, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// If using JWT, generate token here

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"username": user.Username,
		},
	})
}

func AuthenticateUser(email, password string) (models.User, error) {
	user, err := repositories.GetUserByEmailForLogin(email)
	if err != nil {
		return user, errors.New("invalid email or password")
	}

	// Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, errors.New("invalid email or password")
	}

	return user, nil
}