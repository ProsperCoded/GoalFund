# KYC & Email Verification - Quick Reference

## User Model Fields

```go
type User struct {
    // ... other fields ...

    // Email Verification
    EmailVerified   bool       `json:"email_verified"`

    // KYC Verification
    NIN             string     `json:"nin,omitempty"`
    KYCVerified     bool       `json:"kyc_verified"`
    KYCVerifiedAt   *time.Time `json:"kyc_verified_at,omitempty"`
}
```

## API Endpoints

### Submit NIN

```http
POST /api/v1/users/kyc/submit-nin
Authorization: Bearer {token}
Content-Type: application/json

{
  "nin": "12345678901"
}
```

**Response:**

```json
{
  "kyc_verified": true,
  "kyc_verified_at": "2026-01-31T16:54:19Z",
  "nin": "*******8901"
}
```

### Get KYC Status

```http
GET /api/v1/users/kyc/status
Authorization: Bearer {token}
```

**Response:**

```json
{
  "kyc_verified": true,
  "kyc_verified_at": "2026-01-31T16:54:19Z",
  "nin": "*******8901"
}
```

## User Response (All Auth Endpoints)

Login, Register, GetProfile, and UpdateProfile now return:

```json
{
  "id": "uuid",
  "email": "user@example.com",
  "username": "johndoe",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+2348012345678",
  "email_verified": true, // ← Email verification status
  "phone_verified": false,
  "kyc_verified": true, // ← KYC verification status
  "kyc_verified_at": "2026-01-31T16:54:19Z",
  "role": "user",
  "created_at": "2026-01-15T10:30:00Z"
}
```

## Validation Rules

### NIN

- ✅ Must be exactly 11 digits
- ✅ Numeric only (no letters or special characters)
- ✅ Cannot be duplicate (already registered to another user)
- ✅ Automatically masked in responses (shows last 4 digits)

### Email Verification

- ✅ Tracked via `email_verified` boolean
- ✅ Defaults to `false` on registration
- ✅ Can be updated through email verification flow

## Error Codes

| Status | Error Message                               | Cause                      |
| ------ | ------------------------------------------- | -------------------------- |
| 400    | `invalid NIN format - must be 11 digits`    | NIN is not 11 digits       |
| 400    | `user is already KYC verified`              | User already completed KYC |
| 400    | `NIN already registered to another account` | Duplicate NIN              |
| 401    | `unauthorized`                              | Missing or invalid token   |

## Events Emitted

### KYCVerified

```json
{
  "id": "event-uuid",
  "user_id": "user-uuid",
  "email": "user@example.com",
  "username": "johndoe",
  "created_at": 1706709259
}
```

**Published to:** `gofund.events` exchange
**Event Type:** `KYCVerified`

## Database Migration

Run this migration to add KYC fields:

```bash
# Location: backend/services/users-service/migrations/001_add_kyc_fields.sql
psql -U postgres -d users_db -f migrations/001_add_kyc_fields.sql
```

## Frontend Integration Example

```javascript
// Check if user needs KYC
if (!user.kyc_verified) {
  // Show KYC prompt
  showKYCModal();
}

// Submit NIN
async function submitKYC(nin) {
  try {
    const response = await fetch("/api/v1/users/kyc/submit-nin", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ nin }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }

    const data = await response.json();
    console.log("KYC verified:", data);
    return data;
  } catch (error) {
    console.error("KYC submission failed:", error.message);
    throw error;
  }
}

// Display verification status
function displayVerificationStatus(user) {
  return `
    <div class="verification-status">
      <div>
        Email: ${user.email_verified ? "✅ Verified" : "❌ Not Verified"}
      </div>
      <div>
        KYC: ${user.kyc_verified ? "✅ Verified" : "❌ Not Verified"}
      </div>
    </div>
  `;
}
```

## Testing with cURL

```bash
# 1. Register a user
curl -X POST http://localhost:8084/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'

# 2. Login and get token
curl -X POST http://localhost:8084/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# 3. Submit NIN for KYC
curl -X POST http://localhost:8084/users/kyc/submit-nin \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nin": "12345678901"
  }'

# 4. Check KYC status
curl -X GET http://localhost:8084/users/kyc/status \
  -H "Authorization: Bearer YOUR_TOKEN"

# 5. Get updated profile
curl -X GET http://localhost:8084/users/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Important Notes

⚠️ **Dummy Implementation**: This is a practice project. All NIN submissions are automatically approved without actual verification.

🔒 **Privacy**: NIN is always masked in API responses (e.g., `*******8901`)

📧 **Email Verification**: The `email_verified` field is tracked but the actual verification flow needs to be implemented separately.

🎯 **Production Ready**: For production use, integrate with actual NIN verification services and implement proper email verification.
