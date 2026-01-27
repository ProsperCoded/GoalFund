# Authentication Flow Implementation

## Overview

GoFund uses **Gateway-Level Authentication** with Nginx `auth_request` module for centralized authentication across all microservices.

## Architecture

```
Client → Nginx Gateway → Users Service (auth) → Target Service
```

## Implementation Details

### 1. Nginx Configuration

#### Rate Limiting Zones:
- `api`: 10 req/s (general API endpoints)
- `auth`: 5 req/s (authentication endpoints)
- `contribute`: 2 req/s (anonymous contributions - tight limit)
- `webhook`: 1 req/s (payment webhooks - very tight limit)

#### Connection Pooling:
- `keepalive 32` for users-service (high auth traffic)
- `keepalive 16` for other services

#### Auth Response Caching:
- Valid tokens cached for 30 seconds
- Invalid tokens cached for 10 seconds
- Cache key: Authorization header value

### 2. Route Classification

#### Public Routes (No Authentication):
- `POST /api/v1/users/login`
- `POST /api/v1/users/register`
- `POST /api/v1/users/forgot-password`
- `POST /api/v1/users/reset-password`
- `POST /api/v1/payments/webhook` (Paystack callbacks)
- `POST /api/v1/payments/contribute` (anonymous contributions)
- `GET /api/v1/goals/list` (browse goals)
- `GET /api/v1/goals/view` (view specific goal)

#### Protected Routes (Authentication Required):
- All other `/api/v1/users/*` endpoints
- All other `/api/v1/goals/*` endpoints
- All `/api/v1/ledger/*` endpoints
- All other `/api/v1/payments/*` endpoints
- All `/api/v1/notifications/*` endpoints

### 3. Authentication Flow

#### Step 1: Client Request
```http
GET /api/v1/users/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Step 2: Nginx Auth Subrequest
```http
GET /internal/verify
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
X-Original-URI: /api/v1/users/profile
X-Real-IP: 192.168.1.100
```

#### Step 3: Users Service Validation
- Parse JWT token
- Validate signature and expiration
- Extract user context
- Return headers:
  ```http
  X-User-ID: user-123
  X-User-Email: john@example.com
  X-User-Roles: user,premium
  ```

#### Step 4: Nginx Forwards to Target Service
```http
GET /users/profile
X-User-ID: user-123
X-User-Email: john@example.com
X-User-Roles: user,premium
```

#### Step 5: Target Service Response
Target service can trust the user context headers since they come from authenticated source.

### 4. Users Service Implementation

#### Internal Auth Endpoint (`/internal/verify`):
- **Purpose**: Validate JWT tokens for Nginx auth_request
- **Performance**: Optimized for speed (no DB calls)
- **Caching**: Responses cached by Nginx for 30 seconds
- **Headers**: Returns user context in response headers

#### Public Auth Endpoints:
- `POST /auth/login` - User authentication
- `POST /auth/register` - User registration
- `POST /auth/refresh` - Token refresh
- `POST /auth/logout` - Token invalidation
- `POST /auth/forgot-password` - Password reset request
- `POST /auth/reset-password` - Password reset with token

#### Protected User Endpoints:
- `GET /users/profile` - Get current user profile
- `PUT /users/profile` - Update current user profile

### 5. Performance Optimizations

#### Connection Pooling:
```nginx
upstream users-service {
    server users-service:8084;
    keepalive 32;  # Keep 32 connections alive
}
```

#### Auth Response Caching:
```nginx
proxy_cache_path /tmp/auth_cache levels=1:2 keys_zone=auth:10m max_size=100m inactive=1m;

location /auth/verify {
    proxy_cache auth;
    proxy_cache_key $http_authorization;
    proxy_cache_valid 200 30s;  # Cache valid tokens
    proxy_cache_valid 401 403 10s;  # Cache invalid tokens briefly
}
```

#### Lightweight Auth Handler:
- No database queries (JWT is self-contained)
- Minimal response (only headers)
- Fast JWT validation (~0.1-0.5ms)

### 6. Security Features

#### Rate Limiting:
- Different limits for different endpoint types
- Burst handling for legitimate traffic spikes
- IP-based limiting to prevent abuse

#### Security Headers:
- `X-Frame-Options: DENY`
- `X-Content-Type-Options: nosniff`
- `X-XSS-Protection: 1; mode=block`
- `Referrer-Policy: strict-origin-when-cross-origin`

#### Token Security:
- JWT tokens with expiration
- Secure token transmission (Authorization header)
- Token caching with appropriate TTL

### 7. Error Handling

#### Authentication Failures:
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - Valid token but insufficient permissions
- Rate limiting returns `429 Too Many Requests`

#### Service Communication:
- Internal HTTP (no HTTPS overhead)
- Connection pooling reduces latency
- Graceful fallback if auth service is down

### 8. Monitoring & Observability

#### Nginx Logs:
- Request timing including auth subrequest time
- Cache hit/miss ratios
- Rate limiting events

#### Metrics to Track:
- Auth request latency
- Cache hit ratio
- Authentication success/failure rates
- Rate limiting triggers

### 9. TODO: Implementation Tasks

#### Users Service:
1. Implement JWT token generation/validation
2. Add user repository and database operations
3. Implement password hashing and validation
4. Add email verification and password reset
5. Implement role-based access control
6. Add token blacklisting for logout
7. Add comprehensive logging and metrics

#### Other Services:
1. Update services to read user context from headers
2. Implement role-based authorization per endpoint
3. Add user context to business logic
4. Update error handling for authentication failures

This authentication system provides centralized security, excellent performance, and clear separation of concerns while maintaining the flexibility to customize authorization per service.