# Simple Dockerfile for Transaction Reconciler Service
FROM golang:1.24.5-alpine

# Set working directory
WORKDIR /app

# Copy all files
COPY . .

# Build the application
RUN go build -o main .

# Run the application
CMD ["./main"]
