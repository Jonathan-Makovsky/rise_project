# rise_project
Rise home assignment


**Installations:**
1. Docker application

**Instruction for running backend:**
1. Stop and remove all containers: "**docker-compose down -v**"
2. Docker build: "**docker-compose up --build**"
3. **Using UI:**
  Open the index.html file by your preferred web browser

**Tests:**
local tests
1. Run test/repository.tests file: "XXXX"

End to End tests:
Runs by the docker compose


Programming languages:
Backend - GO
Framework - gorilla/mux
Client - HTML, js

**Project hierarchy:**
Rise/ 
├── src/ # Source files
│ ├── handler.go # API handler functions for CRUD operations
│ ├── repository.go # Database interaction functions
├── setup/ # setup files (Docker)
│ ├── Dockerfile
│ ├── docker-compose.yml 
├── database/ # complete
│ ├── init.sql
├── frontend/ # UI
│ ├── index.html
├── tests/ # Testing
│ ├── repository_test.go
├── main.go # Entry point for the API server and routes
├── go.mod # Go module dependencies
├── go.sum # Go module checksum
├── README.md # Project documentation
