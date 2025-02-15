package src

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Pagination control
var paginationOffset int
var paginationLock sync.Mutex

// GetContactsHandler handles the HTTP request for retrieving contacts
func GetContactsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const limit = 10 // Lock to control concurrent access to pagination
		paginationLock.Lock()
		offset := paginationOffset
		paginationOffset += limit
		paginationLock.Unlock()

		contacts, message, err := GetContacts(db, limit, offset)
		if err != nil {
			http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if len(contacts) < limit {
			paginationOffset = 0 // Reset pagination if fewer contacts are returned
		}
		
		// Send response with contacts and any message
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
}

// AddContactHandler handles the HTTP request for adding contact
func AddContactHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var contact Contact

		// Decode the incoming JSON body into a Contact struct
		if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
			// Return a good response instead of an error
			response := struct {
				Message string `json:"message"`
			}{
				Message: "Invalid request body. Please provide correct JSON format.",
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Check for empty fields
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

		// If any fields are empty, return a good response with a message
		if emptyCount > 0 {
			response := struct {
				Message string `json:"message"`
			}{
				Message: fmt.Sprintf("%d field(s) are empty. Please provide all required fields.", emptyCount),
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Insert the contact into the database
		id, err := AddContact(db, contact)
		if err != nil {
			// Return a success response with an appropriate message
			response := struct {
				Message string `json:"message"`
			}{
				Message: "Database error occurred while adding the contact.",
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}

		contact.ID = id

		// Return success message
		response := struct {
			Message string `json:"message"`
		}{
			Message: "Contact was added successfully",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// DeleteContactHandler handles the HTTP request for deleting contact by his exact phone number
func DeleteContactHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		phoneNumber := vars["phone_number"]
		if phoneNumber == "" {
			// Return a proper message if no phone number was provided
			response := struct {
				Message string `json:"message"`
			}{
				Message: "No number was given",
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Attempt to delete the contact from the db
		rowsDeleted, err := DeleteContact(db, phoneNumber)
		if err != nil {
			// If no rows were deleted, it means the number is not in the phone book
			response := struct {
				Message string `json:"message"`
			}{
				Message: "The number provided is not in the phone book",
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Return success message with the number of deleted contacts
		response := struct {
			Message string `json:"message"`
		}{
			Message: fmt.Sprintf("%d contact(s) were deleted", rowsDeleted),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// SearchContactHandler handles the HTTP request for searching contact by his phone number
func SearchContactHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		phoneNumber := vars["phone_number"]
		if phoneNumber == "" {
			http.Error(w, "Phone number not provided", http.StatusBadRequest)
			return
		}

		contacts, err := SearchContact(db, phoneNumber)
		if err != nil {
			// // Return a message if no contacts are found
			response := struct {
				Message  string    `json:"message"`
				Contacts []Contact `json:"contacts"`
			}{
				Message:  "No contacts were found with the given phone number",
				Contacts: []Contact{},
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Return the found contacts with a message
		response := struct {
			Message  string    `json:"message"`
			Contacts []Contact `json:"contacts"`
		}{
			Message:  fmt.Sprintf("%d contacts with the given phone number were found", len(contacts)),
			Contacts: contacts,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// EditContactHandler handles the HTTP request for editing contact by his phone number
func EditContactHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		phoneNumber := vars["phone_number"]
		if phoneNumber == "" {
			// Return a message if no phone number is provided
			response := struct {
				Message string `json:"message"`
			}{
				Message: "No number was given",
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}

		var updatedContact Contact
		if err := json.NewDecoder(r.Body).Decode(&updatedContact); err != nil {
			// Return a message if the request body is invalid
			response := struct {
				Message string `json:"message"`
			}{
				Message: "Invalid request body. Please provide correct JSON format.",
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Validate that all fields are provided (not empty)
		emptyCount := 0
		if updatedContact.FirstName == "" {
			emptyCount++
		}
		if updatedContact.LastName == "" {
			emptyCount++
		}
		if updatedContact.PhoneNumber == "" {
			emptyCount++
		}
		if updatedContact.Address == "" {
			emptyCount++
		}

		// If any fields are empty, return a success response with a message
		if emptyCount > 0 {
			response := struct {
				Message string `json:"message"`
			}{
				Message: fmt.Sprintf("%d field(s) are empty. Please provide all required fields.", emptyCount),
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Attempt to update the contact
		rowsUpdated, err := EditContact(db, phoneNumber, updatedContact)
		if err != nil {
			// If no rows were updated, it means the number is not in the phone book
			response := struct {
				Message string `json:"message"`
			}{
				Message: "The number provided is not in the phone book",
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Return success message
		response := struct {
			Message string `json:"message"`
		}{
			Message: fmt.Sprintf("%d contact(s) were updated successfully", rowsUpdated),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
