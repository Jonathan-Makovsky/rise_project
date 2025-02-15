#!/bin/bash

# Step 1: Start the Docker containers using Docker Compose
echo "Starting Docker containers..."
docker-compose up -d --build

# Step 2: Wait for the services to be up and running (you can adjust the sleep time)
echo "Waiting for the services to be ready..."
sleep 10  # You can increase this value if needed to ensure the containers are ready

# Step 3: Run the test file for the backend API endpoints
echo "Running the tests..."
go test -v ./tests/end_to_end_test.go -count=1  # This ensures tests are run fresh without cache

# Step 4: Shut down the Docker containers after the tests have completed
echo "Shutting down Docker containers..."
docker-compose down

# Optionally, you can check the exit code of the tests to decide what to do next
if [ $? -eq 0 ]; then
  echo "Tests passed successfully!"
else
  echo "Tests failed. Check the output for errors."
fi
