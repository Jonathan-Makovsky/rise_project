// tests/repository_tests.go

package tests

import (
    "fmt"
    "regexp"
    "testing"
    "github.com/DATA-DOG/go-sqlmock"
    "Rise/src" 
) 
// Test function to run all tests
func TestRepository(t *testing.T) {
    t.Run("Test Add, Search, Delete Contact", testAddSearchDeleteContact)
    t.Run("Test Pagination with Multiple Contacts", testPaginationWithMultipleContacts)
    t.Run("Test Add and Delete Contacts", testAddDeleteContacts)
    t.Run("Test add and edit Contacts", testEditContact)

    
}

// Test adding, searching, and deleting a contact
func testAddSearchDeleteContact(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    newContact := src.Contact{
        FirstName:   "Jonathan",
        LastName:    "Makovsky",
        PhoneNumber: "0543435590",
        Address:     "Tel Aviv",
    }

    // Mock the insert query
    mock.ExpectQuery(regexp.QuoteMeta(
        "INSERT INTO contacts (first_name, last_name, phone_number, address) VALUES ($1, $2, $3, $4) RETURNING id",
    )).WithArgs(newContact.FirstName, newContact.LastName, newContact.PhoneNumber, newContact.Address).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

    // Add the contact and check for errors
    id, err := src.AddContact(db, newContact)
    if err != nil {
        t.Fatalf("Failed to add contact: %v", err)
    }
    newContact.ID = id

    // Mock the query to count contacts
    mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM contacts")).
        WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

    row := db.QueryRow("SELECT COUNT(*) FROM contacts")
    var count int
    if err := row.Scan(&count); err != nil {
        t.Fatalf("Failed to scan count: %v", err)
    }
    if count != 1 {
        t.Fatalf("Expected 1 contact, found %d", count)
    }

    // Mock the search query by phone number
    mock.ExpectQuery(regexp.QuoteMeta(
        "SELECT id, first_name, last_name, phone_number, address FROM contacts WHERE phone_number = $1",
    )).WithArgs(newContact.PhoneNumber).
        WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "phone_number", "address"}).
            AddRow(newContact.ID, newContact.FirstName, newContact.LastName, newContact.PhoneNumber, newContact.Address))

    // Search for the contact and check the result
    contacts, err := src.SearchContact(db, newContact.PhoneNumber)
    if err != nil {
        t.Fatalf("Failed to search contact: %v", err)
    }
    if len(contacts) != 1 || contacts[0] != newContact {
        t.Fatalf("Expected contact %+v, got %+v", newContact, contacts[0])
    }

    // Mock the delete query
    mock.ExpectExec(regexp.QuoteMeta(
        "DELETE FROM contacts WHERE phone_number = $1",
    )).WithArgs(newContact.PhoneNumber).
        WillReturnResult(sqlmock.NewResult(0, 1))

    // Delete the contact and check for success
    deletedCount, err := src.DeleteContact(db, newContact.PhoneNumber)
    if err != nil {
        t.Fatalf("Failed to delete contact: %v", err)
    }
    if deletedCount != 1 {
        t.Fatalf("Expected to delete 1 contact, deleted %d", deletedCount)
    }

    // Mock the query to check the count of contacts after deletion
    mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM contacts")).
        WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

    row = db.QueryRow("SELECT COUNT(*) FROM contacts")
    if err := row.Scan(&count); err != nil {
        t.Fatalf("Failed to scan count: %v", err)
    }
    if count != 0 {
        t.Fatalf("Expected 0 contacts, found %d", count)
    }

    // check all mock expectations were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("There were unfulfilled expectations: %s", err)
    }
}

// Test pagination with multiple contacts
func testPaginationWithMultipleContacts(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    var contactsToAdd []src.Contact
    for i := 1; i <= 25; i++ {
        contactsToAdd = append(contactsToAdd, src.Contact{
            FirstName:   fmt.Sprintf("FirstName%d", i),
            LastName:    fmt.Sprintf("LastName%d", i),
            PhoneNumber: fmt.Sprintf("12345678%02d", i),
            Address:     fmt.Sprintf("%d Elm St", i),
        })
    }
    // Mock the insert query
    for i, contact := range contactsToAdd {
        mock.ExpectQuery(regexp.QuoteMeta(
            "INSERT INTO contacts (first_name, last_name, phone_number, address) VALUES ($1, $2, $3, $4) RETURNING id",
        )).WithArgs(contact.FirstName, contact.LastName, contact.PhoneNumber, contact.Address).
            WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(i + 1))

        id, err := src.AddContact(db, contact)
        if err != nil {
            t.Fatalf("Failed to add contact: %v", err)
        }
        contactsToAdd[i].ID = id
    }

    pageSize := 10
    totalContacts := len(contactsToAdd)
    totalPages := (totalContacts + pageSize - 1) / pageSize

    for page := 0; page < totalPages; page++ {
        offset := page * pageSize
        limit := pageSize

        expectedCount := pageSize
        if offset+limit > totalContacts {
            expectedCount = totalContacts - offset
        }

        rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "phone_number", "address"})
        for i := offset; i < offset+expectedCount; i++ {
            contact := contactsToAdd[i]
            rows.AddRow(contact.ID, contact.FirstName, contact.LastName, contact.PhoneNumber, contact.Address)
        }
        mock.ExpectQuery(regexp.QuoteMeta(
            "SELECT id, first_name, last_name, phone_number, address FROM contacts LIMIT $1 OFFSET $2",
        )).WithArgs(limit, offset).
            WillReturnRows(rows)

        // Correct variable assignment
        contacts, message, err := src.GetContacts(db, limit, offset)
        if err != nil {
            t.Fatalf("Failed to retrieve contacts: %v", err)
        }
        if len(contacts) != expectedCount {
            t.Fatalf("Expected to retrieve %d contacts, but got %d", expectedCount, len(contacts))
        }

        // Optionally, handle the message if needed
        if message != "" {
            fmt.Printf("Message: %s\n", message)
        }
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("There were unfulfilled expectations: %s", err)
    }
}

// Test adding and deleting contacts
func testAddDeleteContacts(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    contactsToAdd := []src.Contact{
        {FirstName: "Alice", LastName: "Smith", PhoneNumber: "1111111111", Address: "123 Maple St"},
        {FirstName: "Bob", LastName: "Johnson", PhoneNumber: "2222222222", Address: "456 Oak St"},
    }

    for i, contact := range contactsToAdd {
        // Mock the insert query
        mock.ExpectQuery(regexp.QuoteMeta(
            "INSERT INTO contacts (first_name, last_name, phone_number, address) VALUES ($1, $2, $3, $4) RETURNING id",
        )).
            WithArgs(contact.FirstName, contact.LastName, contact.PhoneNumber, contact.Address).
            WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(i + 1))

        id, err := src.AddContact(db, contact)
        if err != nil {
            t.Fatalf("Failed to add contact: %v", err)
        }
        contactsToAdd[i].ID = id
    }
    // Mock the delete query
    mock.ExpectExec(regexp.QuoteMeta(
        "DELETE FROM contacts WHERE phone_number = $1",
    )).
        WithArgs(contactsToAdd[0].PhoneNumber).
        WillReturnResult(sqlmock.NewResult(0, 1))

    deletedCount, err := src.DeleteContact(db, contactsToAdd[0].PhoneNumber)
    if err != nil {
        t.Fatalf("Failed to delete contact: %v", err)
    }

    if deletedCount != 1 {
        t.Fatalf("Expected to delete 1 contact, deleted %d", deletedCount)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("There were unfulfilled expectations: %s", err)
    }
}

// Test editing a contact
func testEditContact(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    // Step 1: Add a contact
    newContact := src.Contact{
        FirstName:   "Jonathan",
        LastName:    "Makovsky",
        PhoneNumber: "0543435590",
        Address:     "Tel Aviv",
    }
    // Mock the insert query
    mock.ExpectQuery(regexp.QuoteMeta(
        "INSERT INTO contacts (first_name, last_name, phone_number, address) VALUES ($1, $2, $3, $4) RETURNING id",
    )).WithArgs(newContact.FirstName, newContact.LastName, newContact.PhoneNumber, newContact.Address).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

    id, err := src.AddContact(db, newContact)
    if err != nil {
        t.Fatalf("Failed to add contact: %v", err)
    }
    newContact.ID = id

    // Step 2: Attempt to edit someone's phone number (not possible)
    otherPhoneNumber := "0543435591"
    updatedContact := src.Contact{
        FirstName:   "Jonathan",
        LastName:    "Makovsky",
        PhoneNumber: otherPhoneNumber, // Someone else's phone number
        Address:     "New Address",
    }

    mock.ExpectExec(regexp.QuoteMeta(
        "UPDATE contacts SET first_name = $1, last_name = $2, phone_number = $3, address = $4 WHERE phone_number = $5",
    )).WithArgs(updatedContact.FirstName, updatedContact.LastName, updatedContact.PhoneNumber, updatedContact.Address, newContact.PhoneNumber).
        WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected (not possible)

    // Try to edit with a different phone number, should not succeed
    rowsUpdated, err := src.EditContact(db, newContact.PhoneNumber, updatedContact)
    if err == nil || rowsUpdated != 0 {
        t.Fatalf("Expected no rows to be updated when editing a non-existent contact, but got %d rows", rowsUpdated)
    }

    // Step 3: Try to edit the same contact's phone number (should succeed)
    updatedContact.PhoneNumber = newContact.PhoneNumber // Correct the phone number to match
    mock.ExpectExec(regexp.QuoteMeta(
        "UPDATE contacts SET first_name = $1, last_name = $2, phone_number = $3, address = $4 WHERE phone_number = $5",
    )).WithArgs(updatedContact.FirstName, updatedContact.LastName, updatedContact.PhoneNumber, updatedContact.Address, newContact.PhoneNumber).
        WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row updated

    // Now edit with the correct phone number
    rowsUpdated, err = src.EditContact(db, newContact.PhoneNumber, updatedContact)
    if err != nil || rowsUpdated != 1 {
        t.Fatalf("Expected 1 row to be updated when editing the contact, but got %d rows", rowsUpdated)
    }

    // Check all mock expectations were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("There were unfulfilled expectations: %s", err)
    }
}