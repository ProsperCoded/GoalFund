# =========================
# Build stage
# =========================
FROM  golang:1.24-alpine AS builder

WORKDIR /app

# Build argument
ARG APP_NAME

# ---- Dependency layer (MAX cache reuse) ----
# Copy ONLY module files
COPY shared/go.mod shared/go.sum ./shared/
COPY services/${APP_NAME}/go.mod services/${APP_NAME}/go.sum ./services/${APP_NAME}/

# Download deps
WORKDIR /app/services/${APP_NAME}
RUN go mod download

# ---- Source code layer ----
WORKDIR /app
COPY shared/ ./shared/
COPY services/${APP_NAME}/ ./services/${APP_NAME}/

# ---- Build ----
WORKDIR /app/services/${APP_NAME}
RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -buildvcs=false -ldflags="-s -w" \
    -o main ./cmd/main.go

# =========================
# Runtime stage
# =========================
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

ARG APP_NAME

COPY --from=builder /app/services/${APP_NAME}/main .

CMD ["./main"]
