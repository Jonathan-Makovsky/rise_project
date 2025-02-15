# Dockerfile
FROM golang:1.21-alpine

WORKDIR /app

# Install required dependencies
RUN apk add --no-cache git

# Copy go module files first
COPY ./go.mod ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the application
RUN go build -o main .

# Make sure the binary is executable
RUN chmod +x main

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]
