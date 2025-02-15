#!/bin/bash

# Define API base URL
API_URL="http://localhost:8080"

# Function to check HTTP response status
check_response() {
    url=$1
    expected_code=$2

    echo "Testing $url"
    response_code=$(curl -s -o /dev/null -w "%{http_code}" $url)

    if [ "$response_code" -ne "$expected_code" ]; then
        echo "Error: Expected status $expected_code, got $response_code"
        exit 1
    else
        echo "Success: Got expected status $response_code"
    fi
}

# Test 1: Check if the /getContacts endpoint returns a 200 status code
echo "Testing GET /getContacts"
check_response "$API_URL/getContacts" 200

# Optional: Check if the response contains contacts data (if necessary)
response=$(curl -s "$API_URL/getContacts")
if [[ "$response" == *"contacts"* ]]; then
    echo "Success: Contacts data found in response"
else
    echo "Error: No contacts found in the response"
    exit 1
fi

echo "Test 1 completed successfully."
