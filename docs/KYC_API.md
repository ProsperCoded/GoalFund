# KYC Verification API Documentation

## Overview

The KYC (Know Your Customer) verification system allows users to submit their National Identification Number (NIN) for basic identity verification. This is a **dummy implementation** designed for practice purposes - all submissions are automatically approved without actual verification.

## Features

- ✅ NIN submission and validation
- ✅ Automatic verification (dummy implementation)
- ✅ Privacy-focused NIN masking
- ✅ Duplicate NIN prevention
- ✅ Email verification tracking
- ✅ Event-driven notifications

## Endpoints

### 1. Submit NIN for Verification

Submit a National Identification Number for KYC verification.

**Endpoint:** `POST /api/v1/users/kyc/submit-nin`

**Authentication:** Required (Bearer Token)

**Request Body:**

```json
{
  "nin": "12345678901"
}
```

**Validation Rules:**

- NIN must be exactly 11 digits
- NIN must be numeric only
- NIN cannot be already registered to another account

**Success Response (200 OK):**

```json
{
  "kyc_verified": true,
  "kyc_verified_at": "2026-01-31T16:54:19Z",
  "nin": "*******8901"
}
```

**Error Responses:**

- **400 Bad Request** - Invalid NIN format

```json
{
  "error": "invalid NIN format - must be 11 digits"
}
```

- **400 Bad Request** - Already verified

```json
{
  "error": "user is already KYC verified"
}
```

- **400 Bad Request** - Duplicate NIN

```json
{
  "error": "NIN already registered to another account"
}
```

- **401 Unauthorized** - Missing or invalid token

```json
{
  "error": "unauthorized"
}
```

---

### 2. Get KYC Status

Retrieve the current KYC verification status for the authenticated user.

**Endpoint:** `GET /api/v1/users/kyc/status`

**Authentication:** Required (Bearer Token)

**Success Response (200 OK) - Verified:**

```json
{
  "kyc_verified": true,
  "kyc_verified_at": "2026-01-31T16:54:19Z",
  "nin": "*******8901"
}
```

**Success Response (200 OK) - Not Verified:**

```json
{
  "kyc_verified": false,
  "kyc_verified_at": null,
  "nin": ""
}
```

**Error Response:**

- **401 Unauthorized** - Missing or invalid token

```json
{
  "error": "unauthorized"
}
```

---

## User Profile Integration

KYC status is automatically included in all user profile responses:

**Example User Response:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "username": "johndoe",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+2348012345678",
  "email_verified": true,
  "phone_verified": false,
  "kyc_verified": true,
  "kyc_verified_at": "2026-01-31T16:54:19Z",
  "role": "user",
  "created_at": "2026-01-15T10:30:00Z"
}
```

## Privacy & Security

### NIN Masking

- Full NIN is stored in the database
- Only last 4 digits are shown in API responses
- Format: `*******XXXX` (e.g., `*******8901`)

### Validation

- NIN must be exactly 11 numeric digits
- Regex pattern: `^\d{11}$`
- Duplicate NIN check across all users

### Events

When KYC verification is completed, a `KYCVerified` event is emitted:

```json
{
  "id": "event-uuid",
  "user_id": "user-uuid",
  "email": "user@example.com",
  "username": "johndoe",
  "created_at": 1706709259
}
```

This event can be consumed by other services (e.g., Notifications Service) to send confirmation emails or trigger other workflows.

## Email Verification

The system also tracks email verification status:

- `email_verified` boolean field in user model
- Defaults to `false` upon registration
- Can be updated through email verification flow
- Included in all user responses

## Example Usage

### cURL Example - Submit NIN

```bash
curl -X POST https://api.gofund.com/api/v1/users/kyc/submit-nin \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nin": "12345678901"
  }'
```

### cURL Example - Get Status

```bash
curl -X GET https://api.gofund.com/api/v1/users/kyc/status \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### JavaScript Example

```javascript
// Submit NIN
async function submitNIN(nin) {
  const response = await fetch("/api/v1/users/kyc/submit-nin", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${accessToken}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ nin }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error);
  }

  return await response.json();
}

// Get KYC Status
async function getKYCStatus() {
  const response = await fetch("/api/v1/users/kyc/status", {
    headers: {
      Authorization: `Bearer ${accessToken}`,
    },
  });

  return await response.json();
}
```

## Database Schema

### Users Table - KYC Fields

```sql
-- KYC verification columns
nin VARCHAR(11),                    -- National Identification Number
kyc_verified BOOLEAN DEFAULT FALSE, -- Verification status
kyc_verified_at TIMESTAMP,          -- Verification timestamp

-- Indexes
CREATE INDEX idx_users_nin ON users(nin);
CREATE INDEX idx_users_kyc_verified_at ON users(kyc_verified_at);
```

## Testing Notes

Since this is a **dummy implementation**:

- All valid NIN submissions are automatically approved
- No actual verification against government databases
- Suitable for development and testing only
- **NOT for production use without proper verification integration**

## Future Enhancements

Potential improvements for production:

- Integration with actual NIN verification API
- Document upload support (ID card, passport)
- Multi-level verification (basic, intermediate, advanced)
- Verification expiry and renewal
- Admin review and approval workflow
- Audit trail for verification attempts
