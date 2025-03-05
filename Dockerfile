# Step 1: Build the Go application
FROM golang:1.23.6-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download all dependencies. Dependencies are cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Step 2: Create a smaller image to run the application
FROM alpine:latest  

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 (or any port your app is using)
EXPOSE 8080

# Run the Go binary
CMD ["./main"]