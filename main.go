package main

import (
	"log"
	"menthalhealthjournal/db"
	"menthalhealthjournal/models"
	"menthalhealthjournal/router"
)

func main() {
	db.ConnectDatabase()
	defer db.DB.Close()

	// Create table if not exists
	if err := models.CreateUsersTable(db.DB); err != nil {
		log.Fatal(err)
	}
	

	r := router.SetupRouter()
	r.Run(":8080")
}
