# Go image
FROM golang:1.20-alpine

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the application
COPY . .

# Build the application
RUN go build -o task-service .

# Expose port
EXPOSE 8082

# Run the application
CMD ["./task-service"]
