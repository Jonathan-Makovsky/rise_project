package src

import (
	"database/sql"
    "fmt"

)

// Contact struct represents a contact entry in the database
type Contact struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

// GetContacts retrieves contacts with pagination from the database
func GetContacts(db *sql.DB, limit, offset int) ([]Contact, string, error) {
	// Query database for contacts with limit and offset
	rows, err := db.Query(
		"SELECT id, first_name, last_name, phone_number, address FROM contacts LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var contact Contact
		if err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.PhoneNumber, &contact.Address); err != nil {
			return nil, "", err
		}
		contacts = append(contacts, contact)
	}

	message := ""
	// If fewer contacts are returned than requested, it's the end of the table
	if len(contacts) < limit {
		message = "end of table, move to the start"
	}
	return contacts, message, nil
}

// AddContact inserts a new contact into the database
func AddContact(db *sql.DB, contact Contact) (int, error) {
	// Insert the contact and get the generated ID
	err := db.QueryRow(
		"INSERT INTO contacts (first_name, last_name, phone_number, address) VALUES ($1, $2, $3, $4) RETURNING id",
		contact.FirstName, contact.LastName, contact.PhoneNumber, contact.Address,
	).Scan(&contact.ID)

	if err != nil {
		return 0, err
	}
	return contact.ID, nil
}

// DeleteContact removes a contact by phone number and returns the number of deleted rows
func DeleteContact(db *sql.DB, phoneNumber string) (int, error) {
	// Delete contact by phone number
	result, err := db.Exec("DELETE FROM contacts WHERE phone_number = $1", phoneNumber)
	if err != nil {
		return 0, err
	}
	// Get the number of rows deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rowsAffected == 0 {
		return 0, fmt.Errorf("contact not found")
	}
	return int(rowsAffected), nil
}

// SearchContact retrieves all contacts with the given phone number
func SearchContact(db *sql.DB, phoneNumber string) ([]Contact, error) {
	// Query database for contacts with the given phone number
	rows, err := db.Query(
		"SELECT id, first_name, last_name, phone_number, address FROM contacts WHERE phone_number = $1",
		phoneNumber,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var contact Contact
		if err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.PhoneNumber, &contact.Address); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	// Return an error if no contacts are found
	if len(contacts) == 0 {
		return nil, fmt.Errorf("no contacts found")
	}

	return contacts, nil
}

// EditContact updates an existing contact based on the provided phone number
func EditContact(db *sql.DB, phoneNumber string, updatedContact Contact) (int, error) {
	// Update contact's details based on phone number
	result, err := db.Exec(
		"UPDATE contacts SET first_name = $1, last_name = $2, phone_number = $3, address = $4 WHERE phone_number = $5",
		updatedContact.FirstName, updatedContact.LastName, updatedContact.PhoneNumber, updatedContact.Address, phoneNumber,
	)
	if err != nil {
		return 0, err
	}
	// Get the number of rows affected (updated)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rowsAffected == 0 {
		return 0, fmt.Errorf("no contacts found to update")
	}
	return int(rowsAffected), nil
}



