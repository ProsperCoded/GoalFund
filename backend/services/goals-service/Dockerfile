# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go workspace files
COPY go.work go.work
COPY shared/ ./shared/
COPY services/goals-service/ ./services/goals-service/

# Build the application
WORKDIR /app/services/goals-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/services/goals-service/main .

# Expose port
EXPOSE 8083

# Run the application
CMD ["./main"]
