# Backend Dockerfile
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies and swag CLI
RUN apk add --no-cache git && \
    go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger documentation
RUN swag init -g cmd/main.go --parseDependency --parseInternal

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/

# Copy binary and Swagger docs from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

# Change ownership to non-root user
RUN chown -R appuser:appgroup . && \
    chmod -R 755 ./docs

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]

