# syntax=docker/dockerfile:1

# Build Stage
FROM golang:1.22.5-alpine AS builder

# Set destination for COPY
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /shortlink

# Final Stage
FROM alpine:3.18

# Set destination for COPY
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /shortlink /app/shortlink

# Expose port 8080
EXPOSE 8080

# Run the Go application
CMD ["./shortlink"]
