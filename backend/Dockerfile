# Build stage - builds all services
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go workspace files
COPY go.work go.work
COPY shared/ ./shared/
COPY services/ ./services/

# Download dependencies for all services
WORKDIR /app/services/api-gateway
RUN go mod download

WORKDIR /app/services/users-service
RUN go mod download

WORKDIR /app/services/goals-service
RUN go mod download

WORKDIR /app/services/ledger-service
RUN go mod download

WORKDIR /app/services/payments-service
RUN go mod download

WORKDIR /app/services/notifications-service
RUN go mod download

# Build all services
WORKDIR /app/services/api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/api-gateway ./cmd/main.go

WORKDIR /app/services/users-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/users-service ./cmd/main.go

WORKDIR /app/services/goals-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/goals-service ./cmd/main.go

WORKDIR /app/services/ledger-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/ledger-service ./cmd/main.go

WORKDIR /app/services/payments-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/payments-service ./cmd/main.go

WORKDIR /app/services/notifications-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/notifications-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy all binaries
COPY --from=builder /app/bin/ ./

# Default command (can be overridden in docker-compose)
CMD ["./api-gateway"]
