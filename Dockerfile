# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy only the necessary files to the container (excluding unnecessary files and directories)
COPY go.mod go.sum ./

# Download and install Go dependencies
RUN go mod download

# Copy the rest of the application code to the container
COPY . .

RUN go run main.go

# Build the Go binary
# RUN go build -o myapp

# Run the binary when the container starts
# CMD ["./myapp"]