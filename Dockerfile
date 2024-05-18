# Use the official Golang image as the base image for building
FROM golang:1.22.1 as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files from the backend directory
COPY backend/ .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the backend directory into the container
COPY backend/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Use a minimal base image to reduce the image size
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder image
COPY --from=builder /app/main .

# Copy the public folder from the builder image
COPY --from=builder /app/public ./public

# Copy the .env file from the builder image
COPY --from=builder /app/.env ./.env

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
