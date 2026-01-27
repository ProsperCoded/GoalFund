# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go workspace files
COPY go.work go.work
COPY shared/ ./shared/
COPY services/ ./services/

# Build argument to specify which service to build
ARG APP_NAME
ENV APP_NAME=${APP_NAME}

# Download dependencies and build the specific service
WORKDIR /app/services/${APP_NAME}
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

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
