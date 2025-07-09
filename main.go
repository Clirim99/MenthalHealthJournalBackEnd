// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Clirim99/MenthalHealthJournalBackEnd"
)

func main() {
	// Connect to the database
	db.ConnectDatabase()
	defer db.DB.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from MenthalHealthJournalBackEnd!")
		// You can use db.DB here to query the database
	})

	fmt.Println("🌐 Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
