#!/bin/bash

# Step 1: Start the Docker containers using Docker Compose
echo "Starting Docker containers..."
docker-compose up -d --build

# Step 2: Wait for the services to be up and running (you can adjust the sleep time)
echo "Waiting for the services to be ready..."
sleep 10

# Step 3: Run the test file for the backend API endpoints (adjust the test file command if needed)
echo "Running the tests..."
go test -v ./tests/end_to_end_test.go

# Step 4: Shut down the Docker containers after the tests have completed
echo "Shutting down Docker containers..."
docker-compose down
