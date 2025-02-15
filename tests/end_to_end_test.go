package tests

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "testing"
)

const baseURL = "http://localhost:8080"

// Contact structure for JSON requests
type Contact struct {
    FirstName   string `json:"first_name"`
    LastName    string `json:"last_name"`
    PhoneNumber string `json:"phone_number"`
    Address     string `json:"address"`
}

// Test adding a contact
func TestAddContact(t *testing.T) {
    log.Println("Running TestAddContact...")

    contact := Contact{
        FirstName:   "John",
        LastName:    "Doe",
        PhoneNumber: "1234567890",
        Address:     "123 Main St",
    }

    body, _ := json.Marshal(contact)
    resp, err := http.Post(baseURL+"/addContact", "application/json", bytes.NewBuffer(body))
    if err != nil {
        t.Fatalf("❌ TestAddContact failed: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        log.Println("✅ TestAddContact passed")
    } else {
        t.Errorf("❌ TestAddContact failed: expected status 200, got %d", resp.StatusCode)
    }
}

// Test retrieving contacts
func TestGetContacts(t *testing.T) {
    log.Println("Running TestGetContacts...")

    resp, err := http.Get(baseURL + "/getContacts")
    if err != nil {
        t.Fatalf("❌ TestGetContacts failed: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        log.Println("✅ TestGetContacts passed")
    } else {
        t.Errorf("❌ TestGetContacts failed: expected status 200, got %d", resp.StatusCode)
    }
}

// Test searching for a contact
func TestSearchContact(t *testing.T) {
    log.Println("Running TestSearchContact...")

    resp, err := http.Get(baseURL + "/searchContact/1234567890")
    if err != nil {
        t.Fatalf("❌ TestSearchContact failed: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        log.Println("✅ TestSearchContact passed")
    } else {
        t.Errorf("❌ TestSearchContact failed: expected status 200, got %d", resp.StatusCode)
    }
}

// Test editing a contact
func TestEditContact(t *testing.T) {
    log.Println("Running TestEditContact...")

    updatedContact := Contact{
        FirstName:   "John",
        LastName:    "Doe",
        PhoneNumber: "1234567890",
        Address:     "456 Elm St",
    }

    body, _ := json.Marshal(updatedContact)
    req, _ := http.NewRequest(http.MethodPut, baseURL+"/editContact/1234567890", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        t.Fatalf("❌ TestEditContact failed: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        log.Println("✅ TestEditContact passed")
    } else {
        t.Errorf("❌ TestEditContact failed: expected status 200, got %d", resp.StatusCode)
    }
}

// Test deleting a contact
func TestDeleteContact(t *testing.T) {
    log.Println("Running TestDeleteContact...")

    req, _ := http.NewRequest(http.MethodDelete, baseURL+"/deleteContact/1234567890", nil)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        t.Fatalf("❌ TestDeleteContact failed: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        log.Println("✅ TestDeleteContact passed")
    } else {
        t.Errorf("❌ TestDeleteContact failed: expected status 200, got %d", resp.StatusCode)
    }
}
