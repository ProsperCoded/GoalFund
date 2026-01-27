# API Gateway (Nginx Configuration)

This directory contains Nginx configuration files for the GoFund API Gateway.

## Files

- `nginx.conf` - Main Nginx configuration with routing rules
- `proxy_params` - Common proxy parameters for upstream services

## Routing

The API Gateway routes requests based on URL patterns:

| Route Pattern | Target Service | Port |
|---------------|----------------|------|
| `/api/v1/users/*` | Users Service | 8084 |
| `/api/v1/goals/*` | Goals Service | 8083 |
| `/api/v1/ledger/*` | Ledger Service | 8082 |
| `/api/v1/payments/*` | Payments Service | 8081 |
| `/api/v1/notifications/*` | Notifications Service | 8085 |

## Features

- **Rate Limiting**: 10 req/s for general API, 5 req/s for auth endpoints
- **Security Headers**: X-Frame-Options, X-Content-Type-Options, etc.
- **WebSocket Support**: For real-time notifications
- **Health Check**: `/health` endpoint
- **Compression**: Gzip enabled for common content types
- **Logging**: Detailed access logs with upstream timing

## Usage

The configuration is automatically loaded by Docker Compose. No manual setup required.

Access the API at: `http://localhost:8080/api/v1/`