# Payments Service

The Payments Service handles all payment processing for GoalFund using Paystack as the payment provider. It manages payment initialization, verification, webhook processing, and refund disbursements.

## Features

### ✅ Implemented

- **Payment Initialization**: Initialize payments with Paystack for goal contributions
- **Instant Verification**: Verify payments immediately after user completes checkout
- **Webhook Processing**: Handle Paystack webhooks as backup confirmation with idempotency
- **Bank Operations**: List banks and resolve account numbers
- **Refund Disbursements**: Process refunds via Paystack Transfer API
- **Event Publishing**: Emit `PaymentVerified` events to RabbitMQ
- **Idempotency**: Prevent duplicate payment processing
- **Comprehensive Logging**: Detailed logging with Datadog integration
- **Metrics Tracking**: Payment success/failure rates, webhook processing, API latency

## Architecture

### Hybrid Verification Strategy

The service uses a **hybrid approach** for payment confirmation:

1. **Primary**: Instant verification via API when user completes payment
2. **Backup**: Webhook processing for reliability (handles edge cases)
3. **Idempotency**: Both methods update the same payment record safely

This provides:

- ✅ Instant user feedback (no waiting)
- ✅ 99.9% reliability (webhook safety net)
- ✅ No duplicate contributions

### Payment Flow

```
1. Frontend → POST /api/v1/payments/initialize
2. Backend → Creates payment (INITIATED)
3. Backend → Calls Paystack Initialize API
4. Backend → Returns authorization_url + access_code
5. Frontend → Opens Paystack Inline checkout
6. User → Completes payment
7. Frontend → Calls GET /api/v1/payments/verify/:reference
8. Backend → Verifies with Paystack API
9. Backend → Updates payment (VERIFIED)
10. Backend → Emits PaymentVerified event
11. [Later] Paystack → Sends webhook (idempotent backup)
12. Goals Service → Creates contribution
13. Ledger Service → Records transaction
```

## API Endpoints

### Payment Endpoints

#### Initialize Payment

```http
POST /api/v1/payments/initialize

Request:
{
  "user_id": "uuid",
  "goal_id": "uuid",
  "amount": 50000,              // Amount in kobo (500 NGN)
  "currency": "NGN",
  "email": "user@example.com",
  "callback_url": "https://goalfund.com/payment/callback",
  "metadata": {
    "goal_title": "Education Fund"
  }
}

Response:
{
  "status": "success",
  "data": {
    "payment_id": "uuid",
    "authorization_url": "https://checkout.paystack.com/...",
    "access_code": "abc123xyz",
    "reference": "PAY-abc123xyz"
  }
}
```

#### Verify Payment (Instant)

```http
GET /api/v1/payments/verify/:reference

Response:
{
  "status": "success",
  "data": {
    "payment_id": "uuid",
    "reference": "PAY-abc123xyz",
    "status": "VERIFIED",
    "amount": 50000,
    "currency": "NGN",
    "paid_at": "2026-02-01T12:00:00Z",
    "channel": "card"
  }
}
```

#### Get Payment Status

```http
GET /api/v1/payments/:paymentId/status

Response:
{
  "status": "success",
  "data": {
    "payment_id": "uuid",
    "status": "VERIFIED",
    "amount": 50000,
    "currency": "NGN",
    "created_at": "2026-02-01T12:00:00Z",
    "updated_at": "2026-02-01T12:01:00Z"
  }
}
```

#### List Banks

```http
GET /api/v1/payments/banks?country=nigeria

Response:
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "name": "Access Bank",
      "code": "044"
    }
  ]
}
```

#### Resolve Account Number

```http
GET /api/v1/payments/resolve-account?account_number=0123456789&bank_code=057

Response:
{
  "status": "success",
  "data": {
    "account_number": "0123456789",
    "account_name": "John Doe",
    "bank_code": "057",
    "bank_name": "Zenith Bank"
  }
}
```

### Webhook Endpoint

#### Process Webhook

```http
POST /api/v1/payments/webhook
Headers:
  x-paystack-signature: <signature>

Request (from Paystack):
{
  "event": "charge.success",
  "data": {
    "reference": "PAY-abc123xyz",
    "amount": 50000,
    "status": "success",
    ...
  }
}

Response:
{
  "status": "received"
}
```

## Database Schema (MongoDB)

### Collections

#### payments

```javascript
{
  _id: ObjectId,
  paymentId: String (UUID),           // Unique payment identifier
  paystackReference: String,          // Paystack reference
  userId: String (UUID),              // User making payment
  goalId: String (UUID),              // Goal being funded
  amount: Number (int64),             // Amount in kobo
  currency: String,                   // "NGN"
  status: String,                     // INITIATED, PENDING, VERIFIED, FAILED
  paystackData: Object,               // Raw Paystack response
  createdAt: Date,
  updatedAt: Date
}
```

#### webhook_events

```javascript
{
  _id: ObjectId,
  eventId: String,                    // Paystack event ID (for idempotency)
  event: String,                      // Event type (charge.success, etc.)
  data: Object,                       // Raw webhook payload
  signature: String,                  // Webhook signature
  processed: Boolean,                 // Processing status
  receivedAt: Date,
  processedAt: Date
}
```

#### idempotency_keys

```javascript
{
  _id: ObjectId,
  key: String,                        // Idempotency key
  paymentId: String,                  // Associated payment
  createdAt: Date,
  expiresAt: Date                     // TTL index (24 hours)
}
```

## Configuration

### Environment Variables

```bash
# Service Configuration
PAYMENTS_SERVICE_PORT=8081
PAYMENTS_SERVICE_ENV=development

# Paystack Configuration
PAYSTACK_SECRET_KEY=sk_test_...
PAYSTACK_PUBLIC_KEY=pk_test_...
PAYSTACK_WEBHOOK_SECRET=your_webhook_secret
PAYSTACK_BASE_URL=https://api.paystack.co

# MongoDB Configuration
MONGODB_URI=mongodb://admin:admin123@localhost:27017/payments_db?authSource=admin
PAYMENTS_MONGODB_DATABASE=payments_db

# RabbitMQ Configuration
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# Redis Configuration
REDIS_URL=redis://localhost:6379

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key

# Datadog Configuration
DD_API_KEY=your-datadog-api-key
DD_SITE=us5.datadoghq.com
DD_AGENT_HOST=datadog-agent
DD_TRACE_AGENT_PORT=8126
DD_ENV=dev
DD_VERSION=1.0.0
```

## Running the Service

### Local Development

```bash
# Navigate to service directory
cd backend/services/payments-service

# Install dependencies
go mod tidy

# Run the service
go run cmd/main.go
```

### Docker

```bash
# Build and run with docker-compose
docker-compose up --build payments-service
```

## Testing with Paystack

### Test Mode

The service is configured to use Paystack test keys. Use these test cards:

**Successful Payment:**

- Card Number: `4084 0840 8408 4081`
- CVV: `408`
- Expiry: Any future date
- PIN: `0000`
- OTP: `123456`

**Failed Payment:**

- Card Number: `5060 6666 6666 6666`

### Webhook Testing

1. **Local Testing with ngrok:**

   ```bash
   ngrok http 8081
   ```

2. **Configure Paystack Webhook:**

   - Go to Paystack Dashboard → Settings → Webhooks
   - Add webhook URL: `https://your-ngrok-url.ngrok.io/api/v1/payments/webhook`
   - Copy webhook secret to `PAYSTACK_WEBHOOK_SECRET`

3. **Test Webhook:**
   - Make a test payment
   - Check webhook events in Paystack dashboard
   - Verify webhook was processed in service logs

## Security

### Webhook Signature Verification

All webhooks are verified using HMAC SHA512:

```go
func verifySignature(payload []byte, signature, secret string) bool {
    mac := hmac.New(sha512.New, []byte(secret))
    mac.Write(payload)
    expectedSignature := hex.EncodeToString(mac.Sum(nil))
    return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
```

### Idempotency

- Webhook events are deduplicated using event IDs
- Payment verification is idempotent (safe to call multiple times)
- Idempotency keys expire after 24 hours (TTL index)

## Monitoring

### Metrics

The service tracks the following metrics:

- `payment.initialized.count` - Payments initialized
- `payment.verified.count` - Payments verified successfully
- `payment.failed.count` - Failed payments
- `webhook.received.count` - Webhooks received
- `webhook.processed.count` - Webhooks processed
- `webhook.duplicate.count` - Duplicate webhooks (idempotency)
- `webhook.signature.invalid` - Invalid webhook signatures
- `paystack.api.*.duration` - Paystack API call latency
- `refund.disbursement.initiated` - Refunds initiated

### Logging

All operations are logged with structured logging:

```go
log.Printf("[INFO] Payment verified successfully", map[string]interface{}{
    "payment_id": paymentID,
    "reference":  reference,
    "amount":     amount,
    "channel":    channel,
})
```

## Events Published

### PaymentVerified

Emitted when a payment is successfully verified:

```go
{
  "ID": "event-uuid",
  "PaymentID": "payment-uuid",
  "UserID": "user-uuid",
  "GoalID": "goal-uuid",
  "Amount": 50000,
  "CreatedAt": 1706789123
}
```

## Error Handling

The service implements comprehensive error handling:

- **Paystack API Errors**: Logged and returned with appropriate HTTP status
- **Database Errors**: Logged and returned as 500 Internal Server Error
- **Validation Errors**: Returned as 400 Bad Request
- **Webhook Signature Failures**: Returned as 401 Unauthorized
- **Idempotency Conflicts**: Handled gracefully (return existing result)

## Future Enhancements

- [ ] Support for multiple payment providers (Flutterwave)
- [ ] Recurring payment subscriptions
- [ ] Payment installment plans
- [ ] Virtual dedicated accounts per goal
- [ ] Advanced fraud detection
- [ ] Payment analytics dashboard
- [ ] Automatic retry for failed transfers

## Dependencies

- **Gin**: HTTP web framework
- **MongoDB Driver**: Database operations
- **RabbitMQ**: Event messaging
- **Datadog**: Monitoring and tracing
- **Google UUID**: Unique identifier generation

## License

This service is part of the GoalFund project.
