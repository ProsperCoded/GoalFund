# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Build argument to specify which service to build
ARG APP_NAME
ENV APP_NAME=${APP_NAME}

# Copy module files first (for dependency caching)
COPY shared/go.mod shared/go.sum ./shared/
COPY services/${APP_NAME}/go.mod services/${APP_NAME}/go.sum ./services/${APP_NAME}/

# Download dependencies (this layer will be cached unless go.mod changes)
WORKDIR /app/services/${APP_NAME}
RUN go mod download

# Copy the rest of the source code
WORKDIR /app
COPY shared/ ./shared/
COPY services/${APP_NAME}/ ./services/${APP_NAME}/

# Build the specific service with optimized flags
WORKDIR /app/services/${APP_NAME}
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Build argument needs to be redefined in final stage
ARG APP_NAME
ENV APP_NAME=${APP_NAME}

# Copy the binary from builder
COPY --from=builder /app/services/${APP_NAME}/main .

# Run the application
CMD ["./main"]
