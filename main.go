package main
import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Importing PostgreSQL driver for SQL database interaction

	"Rise/src" // Import my source code
)


// enableCORS sets up Cross-Origin Resource Sharing (CORS) headers for all routes.
// This allows the frontend to make requests to the backend from different domains.
func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow cross-origin requests
        w.Header().Set("Access-Control-Allow-Origin", "*") 
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        // If the request method is OPTIONS, respond with a status of 200 (OK)
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
		// Otherwise, pass the request to the next handler
        next.ServeHTTP(w, r)
	})
}


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

	// Wrap router with CORS middleware
    handler := enableCORS(r)

    // Start server
    log.Printf("Server starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
