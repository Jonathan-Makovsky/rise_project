# rise_project
Rise home assignment


**Installations**
Docker:
Install Docker on your system. You can download it from Docker's official website.
Docker download link: https://www.docker.com/products/docker-desktop/

**Instructions for Running the Backend**
Stop and Remove Existing Containers:
If you have any running containers, stop and remove them with:
**docker-compose down -v**

Build and start the application containers with:
**docker-compose up --build**

Access the UI:
Open the **index.html** file in your preferred web browser to interact with the frontend.

**Tests**
Local Tests:
To run unit tests for the repository functions, run:
**go test tests/repository_test.go**

End-to-End Tests:
These tests run automatically with Docker Compose

**Programming Languages and Technologies**
Backend: Go (Golang)
Web Framework: Gorilla Mux
Frontend: HTML
Database: SQL (PostgreSQL)

**Project hierarchy:**
Rise/ 
├── src/                   # Source files
│   ├── handler.go         # API handler functions for CRUD operations
│   └── repository.go      # Database interaction functions
├── setup/                 # Docker setup files
│   ├── Dockerfile         # Dockerfile for building the application container
│   └── docker-compose.yml # Docker Compose configuration for services
├── database/              # Database-related files
│   └── init.sql           # SQL schema to initialize the database
├── frontend/              # UI files
│   └── index.html         # Frontend HTML file
├── tests/                 # Test files
│   └── repository_test.go # Unit tests for repository functions
├── main.go                # Entry point for the API server and routes
├── go.mod                 # Go module dependencies
├── go.sum                 # Go module checksum
└── README.md              # Project documentation
