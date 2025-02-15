// Show/hide sections dynamically
function showSection(section) {
    document.querySelectorAll(".section").forEach(s => s.style.display = "none");
    document.getElementById(section).style.display = "block";
}

// Load contacts and display them in the table
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

// Search contact by phone number
async function searchContact() {
    const phoneNumber = document.getElementById("searchPhoneNumber").value;
    const tableBody = document.querySelector("#searchResults tbody");
    tableBody.innerHTML = ""; // Clear table before loading new data

    try {
        const response = await fetch(`http://localhost:8080/searchContact/${phoneNumber}`);

        if (!response.ok) {
            throw new Error("Failed to search contacts.");
        }

        const data = await response.json();

        // If the backend returns a message (no contact found)
        if (data.message) {
            // Show the message in the table
            const noResultsRow = document.createElement("tr");
            noResultsRow.innerHTML = `<td colspan="5" style="text-align: center;">${data.message}</td>`;
            tableBody.appendChild(noResultsRow);
        } else {
            // If contacts are found, display them
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
        }
    } catch (error) {
        alert("Error: " + error.message);
    }
}

// Delete a contact (same logic for showing messages as search)
async function deleteContact() {
    const phoneNumber = document.getElementById("deletePhoneNumber").value;

    try {
        const response = await fetch(`http://localhost:8080/deleteContact/${phoneNumber}`, { method: "DELETE" });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.message || "Failed to delete contact.");
        }

        alert(data.message);
    } catch (error) {
        alert("Error: " + error.message);
    }
}

// Event listeners for buttons

// Fetch contacts when "Get Contacts" button is clicked
document.getElementById("getContactsBtn").addEventListener("click", function() {
    loadContacts();
});

// Add a new contact when "Add Contact" button is clicked
document.getElementById("addContactBtn").addEventListener("click", function() {
    addContact();
});

// Delete contact when "Delete Contact" button is clicked
document.getElementById("deleteContactBtn").addEventListener("click", function() {
    const phoneNumber = prompt("Enter phone number of the contact to delete:");
    if (phoneNumber) {
        deleteContact(phoneNumber);
    }
});

// Search contact by phone number when "Search Contact" button is clicked
document.getElementById("searchContactBtn").addEventListener("click", function() {
    const phoneNumber = prompt("Enter phone number to search:");
    if (phoneNumber) {
        searchContact(phoneNumber);
    }
});

// Default section to show when page loads
document.addEventListener("DOMContentLoaded", () => showSection("viewContacts"));
