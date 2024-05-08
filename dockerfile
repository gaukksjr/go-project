# Use the official Golang image as a base image for the build stage
FROM golang:1.21-alpine as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod .
COPY go.sum .

# Download and install Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main /app

# Use a minimal image to run the application
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/main .

# Copy the static directory with HTML templates
COPY static ./static

# Expose the port the application runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
