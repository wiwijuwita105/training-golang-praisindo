# Stage 1: Build stage
FROM golang:1.22.5-alpine AS build

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.sum ./

# Copy the go.mod and go.sum files
COPY go.mod  ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY cashflow-svc/ ./

RUN go mod tidy

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o go-cashflow-svc .

# Stage 2: Final stage
FROM alpine:edge

# Set the working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/go-cashflow-svc .

# menambahkan ca-certificates
RUN apk add --no-cache ca-certificates

# Set the entrypoint command
CMD ["./go-cashflow-svc"]
