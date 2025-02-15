package src

/*
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Contact struct (matches database table columns)
type Contact struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

// Global database instance (set in main.go)
var DB *sql.DB

// Server-side pagination state
var paginationOffset int
var paginationLock sync.Mutex

func GetContacts(w http.ResponseWriter, r *http.Request) {
	const limit = 3

	// Lock to safely access and update the offset
	paginationLock.Lock()
	offset := paginationOffset
	paginationOffset += limit
	paginationLock.Unlock()

	// Query the database with LIMIT and OFFSET
	rows, err := DB.Query(
		"SELECT id, first_name, last_name, phone_number, address FROM contacts LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Collect the results
	var contacts []Contact
	for rows.Next() {
		var contact Contact
		if err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.PhoneNumber, &contact.Address); err != nil {
			http.Error(w, "Error reading result row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		contacts = append(contacts, contact)
	}

	// Determine the message based on the number of contacts returned
	message := ""
	if len(contacts) < limit {
		paginationOffset = 0 // Reset offset
		message = "end of table, move to the start"
	} else {
		message = fmt.Sprintf("We pulled rows %d-%d", offset+1, offset+len(contacts))
	}

	// Build and return the response
	response := struct {
		Message  string    `json:"message"`
		Contacts []Contact `json:"contacts"`
	}{
		Message:  message,
		Contacts: contacts,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// AddContact adds a new contact to the database
func AddContact(w http.ResponseWriter, r *http.Request) {
	var contact Contact

	// Decode the incoming JSON body into a Contact struct
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check for empty fields and count how many are missing
	emptyCount := 0
	if contact.FirstName == "" {
		emptyCount++
	}
	if contact.LastName == "" {
		emptyCount++
	}
	if contact.PhoneNumber == "" {
		emptyCount++
	}
	if contact.Address == "" {
		emptyCount++
	}

	// If any fields are empty, return the count in the response
	if emptyCount > 0 {
		message := ""
		if emptyCount == 1 {
			message = "1 field is empty"
		} else {
			message = fmt.Sprintf("%d fields are empty", emptyCount)
		}
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	// Insert the contact into the database
	err := DB.QueryRow(
		"INSERT INTO contacts (first_name, last_name, phone_number, address) VALUES ($1, $2, $3, $4) RETURNING id",
		contact.FirstName, contact.LastName, contact.PhoneNumber, contact.Address,
	).Scan(&contact.ID)

	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the newly created contact
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(contact)
}

// deleteContact removes a contact by ID from the database
func DeleteContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Get phone number from the URL path
	phoneNumber := vars["phone_number"]
	if phoneNumber == "" {
		http.Error(w, "Phone number not provided", http.StatusBadRequest)
		return
	}

	// Attempt to delete the contact by phone number
	result, err := DB.Exec("DELETE FROM contacts WHERE phone_number = $1", phoneNumber)
	if err != nil {
		http.Error(w, "Error executing delete query: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	// Return a success status with no body
	w.WriteHeader(http.StatusNoContent)
}

// SearchContact retrieves all contacts with the given phone number
func SearchContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Get the phone number from the URL path
	phoneNumber := vars["phone_number"]
	if phoneNumber == "" {
		http.Error(w, "Phone number not provided", http.StatusBadRequest)
		return
	}

	// Query the database for all contacts matching the given phone number
	rows, err := DB.Query(
		"SELECT id, first_name, last_name, phone_number, address FROM contacts WHERE phone_number = $1",
		phoneNumber,
	)
	if err != nil {
		http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Collect all matching rows
	var contacts []Contact
	for rows.Next() {
		var contact Contact
		if err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.PhoneNumber, &contact.Address); err != nil {
			http.Error(w, "Error reading row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		contacts = append(contacts, contact)
	}

	// If no rows were found, return a 404
	if len(contacts) == 0 {
		http.Error(w, "No contacts found with that phone number", http.StatusNotFound)
		return
	}

	// Return all matching contacts as JSON
	json.NewEncoder(w).Encode(contacts)
}

// EditContact updates all fields of a contact based on their current phone number
// EditContact updates all fields of contacts with the given phone number
func EditContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Get the phone number from the URL
	phoneNumber := vars["phone_number"]
	if phoneNumber == "" {
		http.Error(w, "Phone number not provided", http.StatusBadRequest)
		return
	}

	// Decode the updated contact details from the request body
	var updatedContact Contact
	if err := json.NewDecoder(r.Body).Decode(&updatedContact); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Update all matching rows in the database
	result, err := DB.Exec(
		"UPDATE contacts SET first_name = $1, last_name = $2, phone_number = $3, address = $4 WHERE phone_number = $5",
		updatedContact.FirstName, updatedContact.LastName, updatedContact.PhoneNumber, updatedContact.Address, phoneNumber,
	)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "No contacts found to update", http.StatusNotFound)
		return
	}

	// Return success status and the number of rows updated
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%d contacts updated", rowsAffected)))
}
*/
