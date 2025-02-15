// Show/hide sections dynamically
function showSection(section) {
    document.querySelectorAll(".section").forEach(s => s.style.display = "none");
    document.getElementById(section).style.display = "block";
}

// Load contacts and display in table
async function loadContacts() {
    const tableBody = document.querySelector("#contactsTable tbody");
    tableBody.innerHTML = ""; // Clear table before loading new data

    try {
        const response = await fetch("http://localhost:8080/getContacts");
        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.message || "Failed to fetch contacts.");
        }

        data.contacts.forEach(contact => {
            const row = document.createElement("tr");
            row.innerHTML = `
                <td>${contact.id}</td>
                <td>${contact.first_name}</td>
                <td>${contact.last_name}</td>
                <td>${contact.phone_number}</td>
                <td>${contact.address}</td>
            `;
            tableBody.appendChild(row);
        });
    } catch (error) {
        alert("Error: " + error.message);
    }
}

// Add a new contact
async function addContact() {
    const contact = {
        first_name: document.getElementById("firstName").value,
        last_name: document.getElementById("lastName").value,
        phone_number: document.getElementById("phoneNumber").value,
        address: document.getElementById("address").value,
    };

    try {
        const response = await fetch("http://localhost:8080/addContact", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(contact),
        });
        const result = await response.json();
        alert(result.message);
    } catch (error) {
        alert("Error: " + error.message);
    }
}

// Edit an existing contact
async function editContact() {
    const phoneNumber = document.getElementById("editPhoneNumber").value;
    const updatedContact = {
        first_name: document.getElementById("newFirstName").value,
        last_name: document.getElementById("newLastName").value,
        phone_number: phoneNumber, // Keep the same phone number
        address: document.getElementById("newAddress").value,
    };

    try {
        const response = await fetch(`http://localhost:8080/editContact/${phoneNumber}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(updatedContact),
        });
        const result = await response.json();
        alert(result.message);
    } catch (error) {
        alert("Error: " + error.message);
    }
}

// Delete a contact
async function deleteContact() {
    const phoneNumber = document.getElementById("deletePhoneNumber").value;

    try {
        const response = await fetch(`http://localhost:8080/deleteContact/${phoneNumber}`, { method: "DELETE" });
        const result = await response.json();
        alert(result.message);
    } catch (error) {
        alert("Error: " + error.message);
    }
}

// Search for a contact
async function searchContact() {
    const phoneNumber = document.getElementById("searchPhoneNumber").value;
    const tableBody = document.querySelector("#searchResults tbody");
    tableBody.innerHTML = ""; // Clear table before loading new data

    try {
        const response = await fetch(`http://localhost:8080/searchContact/${phoneNumber}`);
        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.message || "Failed to search contacts.");
        }

        data.contacts.forEach(contact => {
            const row = document.createElement("tr");
            row.innerHTML = `
                <td>${contact.id}</td>
                <td>${contact.first_name}</td>
                <td>${contact.last_name}</td>
                <td>${contact.phone_number}</td>
                <td>${contact.address}</td>
            `;
            tableBody.appendChild(row);
        });
    } catch (error) {
        alert("Error: " + error.message);
    }
}

// Set the default section when the page loads
document.addEventListener("DOMContentLoaded", () => showSection("viewContacts"));
