# Build stage
FROM golang:1.24.2 as builder

WORKDIR /app

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy certs and set permissions for them
COPY cert.pem key.pem ./
RUN chmod 644 cert.pem key.pem

# Copy source code
COPY . .

# Build the app
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o server main.go

# Final stage: Use debian instead of distroless
FROM debian:bullseye-slim

WORKDIR /app

# Copy the built binary and certs from the builder
COPY --from=builder /app/server .
COPY --from=builder /app/cert.pem .
COPY --from=builder /app/key.pem .

# Fix permissions for certs (just in case)
RUN chmod 644 /app/cert.pem /app/key.pem

EXPOSE $PORTRAW

ENTRYPOINT ["./server"]


