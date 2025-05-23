FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate

# Use a smaller image for the final container
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/server .
COPY --from=builder /app/migrate .

# Create a script to run migrations and then start the server
RUN echo '#!/bin/sh \
./migrate\n\
if [ $? -eq 0 ]; then\
    echo "Migrations completed successfully"\
    echo "Starting server..."\
    ./server\
else\
    echo "Migrations failed"\
    exit 1\
fi' > /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh

# Expose the application port
EXPOSE 8080

# Run the entrypoint script
CMD ["/app/server"]
