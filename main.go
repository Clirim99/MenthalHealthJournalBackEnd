package main

import (
	"log"
	"menthalhealthjournal/db"
	"menthalhealthjournal/models"
	"menthalhealthjournal/router"
	"menthalhealthjournal/services"
//	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize OpenAI client
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}
	services.InitOpenAIClient()
	log.Println("✅ OpenAI client initialized")

	// Connect to database
	db.ConnectDatabase()
	defer db.DB.Close()

	// Create all tables
	if err := models.CreateUsersTable(db.DB); err != nil {
		log.Fatal("Failed to create users table:", err)
	}
	log.Println("✅ Users table created/verified")

	if err := models.CreateEntriesTable(db.DB); err != nil {
		log.Fatal("Failed to create entries table:", err)
	}
	log.Println("✅ Entries table created/verified")

	if err := models.CreateSummariesTable(db.DB); err != nil {
		log.Fatal("Failed to create summaries table:", err)
	}
	log.Println("✅ Summaries table created/verified")

	if err := models.CreateChatSessionsTable(db.DB); err != nil {
		log.Fatal("Failed to create chat_sessions table:", err)
	}
	log.Println("✅ Chat sessions table created/verified")

	if err := models.CreateChatMessagesTable(db.DB); err != nil {
		log.Fatal("Failed to create chat_messages table:", err)
	}
	log.Println("✅ Chat messages table created/verified")

	// Setup and run router
	r := router.SetupRouter()
	log.Println("🚀 Server starting on :8080")
	r.Run(":8080")
}
