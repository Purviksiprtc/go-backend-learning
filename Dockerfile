# 1. Use official Go image
FROM golang:1.25-alpine

# 2. Set working directory inside container
WORKDIR /app

# 3. Copy go.mod and go.sum first (for dependency caching)
COPY go.mod go.sum ./

# 4. Download dependencies
RUN go mod download

# 5. Copy the rest of the application code
COPY . .

# 6. Build the Go application
RUN go build -o app

# 7. Expose port 8080
EXPOSE 8080

# 8. Run the application
CMD ["./app"]
