package src

import (
	"database/sql"
    "fmt"

)

// Contact struct (matches database table columns)
type Contact struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

// GetContacts retrieves contacts with pagination from the database
func GetContacts(db *sql.DB, limit, offset int) ([]Contact, string, error) {
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
	if len(contacts) < limit {
		message = "end of table, move to the start"
	}
	return contacts, message, nil
}

// AddContact inserts a new contact into the database
func AddContact(db *sql.DB, contact Contact) (int, error) {
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
	result, err := db.Exec("DELETE FROM contacts WHERE phone_number = $1", phoneNumber)
	if err != nil {
		return 0, err
	}

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

	if len(contacts) == 0 {
		return nil, fmt.Errorf("no contacts found")
	}

	return contacts, nil
}

// EditContact updates an existing contact based on the provided phone number
func EditContact(db *sql.DB, phoneNumber string, updatedContact Contact) (int, error) {
	result, err := db.Exec(
		"UPDATE contacts SET first_name = $1, last_name = $2, phone_number = $3, address = $4 WHERE phone_number = $5",
		updatedContact.FirstName, updatedContact.LastName, updatedContact.PhoneNumber, updatedContact.Address, phoneNumber,
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rowsAffected == 0 {
		return 0, fmt.Errorf("no contacts found to update")
	}
	return int(rowsAffected), nil
}



