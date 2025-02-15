// tests/repository_tests.go

package tests

import (
    "fmt"
    "log"
    "regexp"
    "testing"
    "database/sql"
    "github.com/DATA-DOG/go-sqlmock"
    "Rise/src" // Update with your project's path
)

// Test function to run all tests
func TestRepository(t *testing.T) {
    t.Run("Test Add, Search, Delete Contact", testAddSearchDeleteContact)
    t.Run("Test Pagination with Multiple Contacts", testPaginationWithMultipleContacts)
    t.Run("Test Add and Delete Contacts", testAddDeleteContacts)
}

// Test adding, searching, and deleting a contact
func testAddSearchDeleteContact(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    newContact := your_project_path.Contact{
        FirstName:   "Jonathan",
        LastName:    "Makovsky",
        PhoneNumber: "0543435590",
        Address:     "Tel Aviv",
    }

    mock.ExpectQuery(regexp.QuoteMeta(
        "INSERT INTO contacts (first_name, last_name, phone_number, address) VALUES ($1, $2, $3, $4) RETURNING id",
    )).
        WithArgs(newContact.FirstName, newContact.LastName, newContact.PhoneNumber, newContact.Address).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

    id, err := your_project_path.AddContact(db, newContact)
    if err != nil {
        t.Fatalf("Failed to add contact: %v", err)
    }
    newContact.ID = id

    mock.ExpectQuery(regexp.QuoteMeta(
        "SELECT COUNT(*) FROM contacts",
    )).
        WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

    row := db.QueryRow("SELECT COUNT(*) FROM contacts")
    var count int
    if err := row.Scan(&count); err != nil {
        t.Fatalf("Failed to scan count: %v", err)
    }
    if count != 1 {
        t.Fatalf("Expected 1 contact, found %d", count)
    }

    mock.ExpectQuery(regexp.QuoteMeta(
        "SELECT id, first_name, last_name, phone_number, address FROM contacts WHERE phone_number = $1",
    )).
        WithArgs(newContact.PhoneNumber).
        WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "phone_number", "address"}).
            AddRow(newContact.ID, newContact.FirstName, newContact.LastName, newContact.PhoneNumber, newContact.Address))

    contacts, err := your_project_path.SearchContact(db, newContact.PhoneNumber)
    if err != nil {
        t.Fatalf("Failed to search contact: %v", err)
    }
    if len(contacts) != 1 || contacts[0] != newContact {
        t.Fatalf("Expected contact %+v, got %+v", newContact, contacts[0])
    }

    mock.ExpectExec(regexp.QuoteMeta(
        "DELETE FROM contacts WHERE phone_number = $1",
    )).
        WithArgs(newContact.PhoneNumber).
        WillReturnResult(sqlmock.NewResult(0, 1))

    deletedCount, err := your_project_path.DeleteContact(db, newContact.PhoneNumber)
    if err != nil {
        t.Fatalf("Failed to delete contact: %v", err)
    }
    if deletedCount != 1 {
        t.Fatalf("Expected to delete 1 contact, deleted %d", deletedCount)
    }

    mock.ExpectQuery(regexp.QuoteMeta(
        "SELECT COUNT(*) FROM contacts",
    )).
        WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

    row = db.QueryRow("SELECT COUNT(*) FROM contacts")
    if err := row.Scan(&count); err != nil {
        t.Fatalf("Failed to scan count: %v", err)
    }
    if count != 0 {
        t.Fatalf("Expected 0 contacts, found %d", count)
    }

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

    var contactsToAdd []your_project_path.Contact
    for i := 1; i <= 25; i++ {
        contactsToAdd = append(contactsToAdd, your_project_path.Contact{
            FirstName:   fmt.Sprintf("FirstName%d", i),
            LastName:    fmt.Sprintf("LastName%d", i),
            PhoneNumber: fmt.Sprintf("12345678%02d", i),
            Address:     fmt.Sprintf("%d Elm St", i),
        })
    }

    for i, contact := range contactsToAdd {
        mock.ExpectQuery(regexp.QuoteMeta(
            "INSERT INTO contacts (first_name, last_name, phone_number, address) VALUES ($1, $2, $3, $4) RETURNING id",
        )).
            WithArgs(contact.FirstName, contact.LastName, contact.PhoneNumber, contact.Address).
            WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(i + 1))

        id, err := your_project_path.AddContact(db, contact)
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
        )).
            WithArgs(limit, offset).
            WillReturnRows(rows)

        retrievedContacts, message, err := your_project_path.GetContacts(db, limit, offset)
        if err != nil {
            t.Fatalf("Failed to retrieve contacts: %v", err)
        }
        if len(retrievedContacts) != expectedCount {
            t.Fatalf("Expected to retrieve %d contacts, but got %d", expectedCount, len(retrievedContacts))
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

    contactsToAdd := []your_project_path.Contact{
        {FirstName: "Alice", LastName: "Smith", PhoneNumber: "1111111111", Address: "123 Maple St"},
        {FirstName: "Bob", LastName: "Johnson", PhoneNumber: "2222222222", Address: "456 Oak St"},
    }

    for i, contact := range contactsToAdd {
        mock.ExpectQuery(regexp.QuoteMeta(
            "INSERT INTO contacts (first_name, last_name, phone_number, address) VALUES ($1, $2, $3, $4) RETURNING id",
        )).
            WithArgs(contact.FirstName, contact.LastName, contact.PhoneNumber, contact.Address).
            WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(i + 1))

        id, err := your_project_path.AddContact(db, contact)
        if err != nil {
            t.Fatalf("Failed to add contact: %v", err)
        }
        contactsToAdd[i].ID = id
    }

    mock.ExpectExec(regexp.QuoteMeta(
        "DELETE FROM contacts WHERE phone_number = $1",
    )).
        WithArgs(contactsToAdd[0].PhoneNumber).
        WillReturnResult(sqlmock.NewResult(0, 1))

    deletedCount, err := your_project_path.DeleteContact(db, contactsToAdd[0].PhoneNumber)
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
