# Stage 1: Build the application
FROM golang:1.22.5-alpine AS builder

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /app

# Copy go mod and sum files to leverage Docker cache
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source code into the container.
COPY . .

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /main ./cashflow-svc/cmd/web/main.go

# Stage 2: Set up the runtime container
FROM alpine:latest

# Install necessary dependencies (if any) for the binary to run.
RUN apk --no-cache add ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /main /main

# Set the binary as the entrypoint of the container.
ENTRYPOINT ["/main"]

# Expose the port the app runs on
EXPOSE 8080
