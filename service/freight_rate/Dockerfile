# Use the official Golang image as the base image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project directory into the container
COPY . .

# Set the working directory to the book_service directory
WORKDIR /app/service/freight_rate

# Build the Go application
RUN go build -o freight_rate main.go

# Command to run the application
CMD ["./freight_rate"]