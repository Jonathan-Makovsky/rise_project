package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"Rise/src"
)

func main() {
	// Database connection setup
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@db:5432/phonebook?sslmode=disable"
	}

	// Open the database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create router
	r := mux.NewRouter()

	// Use the new handler for retrieving contacts
	r.HandleFunc("/getContacts", src.GetContactsHandler(db)).Methods("GET")
	r.HandleFunc("/addContact", src.AddContactHandler(db)).Methods("POST")
	r.HandleFunc("/deleteContact/{phone_number}", src.DeleteContactHandler(db)).Methods("DELETE")
	r.HandleFunc("/searchContact/{phone_number}", src.SearchContactHandler(db)).Methods("GET")
	r.HandleFunc("/editContact/{phone_number}", src.EditContactHandler(db)).Methods("PUT")

	// Start server
	log.Printf("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
