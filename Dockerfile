# Step 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code and Makefile into the container
COPY . .

# Build the Go binary for linux/amd64
ENV GOARCH=amd64
ENV GOOS=linux
RUN go build -o bin/sample-server sample_app/main.go

# Step 2: Create a lightweight container for running the Go app
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go binary from the builder image
COPY --from=builder /app/bin/sample-server .

# Expose the port that the Go server listens on
EXPOSE 3000

# Command to run the Go server
CMD ["./sample-server"]
