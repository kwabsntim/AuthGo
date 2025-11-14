# Step 1: Use Go base image for building
FROM golang:1.24-alpine AS builder

# Step 2: Set working directory inside container
WORKDIR /app

# Step 3: Copy dependency files first (for better caching)
COPY go.mod go.sum ./

# Step 4: Download dependencies
RUN go mod download

# Step 5: Copy all source code
COPY . .

# Step 6: Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service .

# Step 7: Use minimal image for final container
FROM alpine:latest

# Step 8: Install certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

# Step 9: Set working directory
WORKDIR /root/

# Step 10: Copy binary from builder stage
#change the auth -service to the service name you are building
COPY --from=builder /app/auth-service .


# Step 12: Expose the port your app runs on
EXPOSE 8080

# Step 13: Command to run when container starts
#change the auth-service to the name of the service 
CMD ["./auth-service"]
