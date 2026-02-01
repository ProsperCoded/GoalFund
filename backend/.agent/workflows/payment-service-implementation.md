---
description: Complete implementation plan for Paystack payment service integration
---

# Payment Service Implementation Plan - Paystack Integration

## Executive Summary

This document outlines the complete implementation plan for integrating Paystack payment processing into the GoalFund payment service. The implementation covers:

- **Payment initialization** (for contributions)
- **Webhook verification** (for payment confirmation)
- **In-app web checkout** support for frontend
- **Refund disbursements** (already partially implemented)
- **Bank account resolution and validation**

## Key Decisions

### 1. Payment Confirmation Strategy: **Webhooks (Recommended)**

**Chosen Approach: Webhook-based confirmation with fallback verification**

**Why Webhooks are Better:**

- ✅ **Real-time**: Instant payment confirmation without polling
- ✅ **Reliable**: Paystack guarantees webhook delivery with retries
- ✅ **Scalable**: No need to poll Paystack API repeatedly
- ✅ **Complete data**: Webhooks contain full payment details
- ✅ **Industry standard**: Used by Stripe, Paystack, Flutterwave, etc.

**Alternative (Manual Transaction Log Polling):**

- ❌ Requires periodic polling (inefficient)
- ❌ Introduces delays in payment confirmation
- ❌ Higher API usage and rate limiting concerns
- ❌ More complex state management

**Implementation Strategy:**

- Primary: Webhook-based confirmation (`charge.success` event)
- Fallback: Manual verification endpoint for edge cases (webhook failures)
- Idempotency: Prevent duplicate processing using webhook event IDs

### 2. Webhook Requirements for Web Checkout

**Yes, webhooks are REQUIRED for web checkout.**

Even though Paystack redirects users back to your callback URL after payment, you **must** use webhooks because:

1. **Users may close browser**: Redirect might not happen
2. **Network issues**: Callback URL might fail to load
3. **Security**: Frontend can be manipulated; webhooks are server-verified
4. **Reliability**: Paystack guarantees webhook delivery

**Flow:**

```
User pays → Paystack processes → Webhook sent to backend → Payment verified → User redirected to callback URL
```

The frontend callback URL is for **UX only** (showing success message), not for payment confirmation.

### 3. Frontend In-App Web Checkout

**Implementation: Paystack Inline (Pop-up) Checkout**

The frontend will use Paystack's inline checkout which:

- Opens a modal/iframe within your app (no redirect)
- Supports all payment methods (card, bank transfer, USSD, QR)
- Provides seamless UX
- Returns control to your app after payment

**Frontend Integration Steps:**

1. Backend initializes payment and returns `authorization_url` and `access_code`
2. Frontend uses Paystack Inline JS library to open checkout modal
3. User completes payment in modal
4. Paystack sends webhook to backend (payment confirmation)
5. Frontend receives callback and polls backend for confirmation
6. Backend confirms payment via webhook processing

---

## Architecture Overview

### Service Responsibilities

**Payments Service:**

- Initialize Paystack transactions
- Process and verify webhooks
- Maintain payment state machine
- Handle refund disbursements
- Emit `PaymentVerified` events
- Bank account resolution

**Goals Service:**

- Listen for `PaymentVerified` events
- Create contribution records
- Update goal progress
- Trigger ledger entries

**Ledger Service:**

- Record all financial transactions
- Maintain double-entry accounting
- Compute balances

### Payment Flow (Contributions)

```
1. User initiates contribution (Frontend)
   ↓
2. Frontend calls POST /api/v1/payments/initialize
   ↓
3. Payments Service creates payment record (INITIATED)
   ↓
4. Payments Service calls Paystack Initialize Transaction API
   ↓
5. Paystack returns authorization_url and access_code
   ↓
6. Frontend opens Paystack Inline checkout with access_code
   ↓
7. User completes payment
   ↓
8. Paystack sends webhook to POST /api/v1/payments/webhook
   ↓
9. Payments Service verifies webhook signature
   ↓
10. Payments Service updates payment status (VERIFIED)
   ↓
11. Payments Service emits PaymentVerified event
   ↓
12. Goals Service receives event and creates contribution
   ↓
13. Ledger Service records transaction
   ↓
14. Frontend polls GET /api/v1/payments/{paymentId}/status
   ↓
15. Frontend shows success message
```

---

## Database Schema (MongoDB)

### Collections

#### 1. payments

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

// Indexes
paymentId: unique
paystackReference: unique
userId: 1
goalId: 1
status: 1
createdAt: -1
```

#### 2. webhook_events

```javascript
{
  _id: ObjectId,
  eventId: String,                    // Paystack event ID (for idempotency)
  event: String,                      // Event type (charge.success, transfer.success)
  data: Object,                       // Raw webhook payload
  signature: String,                  // Webhook signature
  processed: Boolean,                 // Processing status
  receivedAt: Date,
  processedAt: Date
}

// Indexes
eventId: unique
event: 1
processed: 1
receivedAt: -1
```

#### 3. idempotency_keys

```javascript
{
  _id: ObjectId,
  key: String,                        // Idempotency key
  paymentId: String,                  // Associated payment
  createdAt: Date,
  expiresAt: Date                     // TTL index (24 hours)
}

// Indexes
key: unique
expiresAt: 1 (TTL index)
```

---

## API Endpoints

### 1. Initialize Payment

**POST /api/v1/payments/initialize**

**Request:**

```json
{
  "user_id": "uuid",
  "goal_id": "uuid",
  "amount": 50000, // Amount in kobo (500 NGN)
  "currency": "NGN",
  "email": "user@example.com",
  "callback_url": "https://goalfund.com/payment/callback",
  "metadata": {
    "goal_title": "Education Fund",
    "contributor_name": "John Doe"
  }
}
```

**Response:**

```json
{
  "status": "success",
  "data": {
    "payment_id": "uuid",
    "authorization_url": "https://checkout.paystack.com/...",
    "access_code": "abc123xyz",
    "reference": "paystack_ref_123"
  }
}
```

**Paystack API Call:**

```
POST https://api.paystack.co/transaction/initialize
Headers:
  Authorization: Bearer SECRET_KEY
  Content-Type: application/json

Body:
{
  "email": "user@example.com",
  "amount": 50000,
  "currency": "NGN",
  "reference": "generated_unique_ref",
  "callback_url": "https://goalfund.com/payment/callback",
  "metadata": {
    "payment_id": "uuid",
    "goal_id": "uuid",
    "user_id": "uuid"
  }
}
```

### 2. Webhook Endpoint

**POST /api/v1/payments/webhook**

**Headers:**

```
x-paystack-signature: signature_hash
```

**Request Body (Paystack sends this):**

```json
{
  "event": "charge.success",
  "data": {
    "id": 123456,
    "reference": "paystack_ref_123",
    "amount": 50000,
    "currency": "NGN",
    "status": "success",
    "paid_at": "2026-02-01T12:00:00Z",
    "customer": {
      "email": "user@example.com"
    },
    "metadata": {
      "payment_id": "uuid",
      "goal_id": "uuid",
      "user_id": "uuid"
    }
  }
}
```

**Response:**

```json
{
  "status": "received"
}
```

**Processing Steps:**

1. Verify webhook signature using PAYSTACK_SECRET_KEY
2. Check idempotency (event.data.id)
3. Extract payment_id from metadata
4. Update payment status to VERIFIED
5. Emit PaymentVerified event
6. Return 200 OK

### 3. Verify Payment (Fallback)

**GET /api/v1/payments/{paymentId}/verify**

**Response:**

```json
{
  "status": "success",
  "data": {
    "payment_id": "uuid",
    "status": "VERIFIED",
    "amount": 50000,
    "reference": "paystack_ref_123",
    "paid_at": "2026-02-01T12:00:00Z"
  }
}
```

**Paystack API Call:**

```
GET https://api.paystack.co/transaction/verify/{reference}
Headers:
  Authorization: Bearer SECRET_KEY
```

### 4. Get Payment Status

**GET /api/v1/payments/{paymentId}/status**

**Response:**

```json
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

### 5. List Banks

**GET /api/v1/payments/banks**

**Response:**

```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "name": "Access Bank",
      "code": "044"
    },
    {
      "id": 2,
      "name": "Zenith Bank",
      "code": "057"
    }
  ]
}
```

### 6. Resolve Account Number

**GET /api/v1/payments/resolve-account?account_number=0123456789&bank_code=057**

**Response:**

```json
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

---

## Implementation Structure

### Directory Structure

```
payments-service/
├── cmd/
│   └── main.go                          # Entry point
├── internal/
│   ├── config/
│   │   └── config.go                    # Configuration
│   ├── controller/
│   │   ├── payment_controller.go        # Payment endpoints
│   │   └── webhook_controller.go        # Webhook endpoint
│   ├── service/
│   │   ├── payment_service.go           # Payment business logic
│   │   ├── webhook_service.go           # Webhook processing
│   │   ├── paystack_client.go           # Paystack API client
│   │   └── refund_disbursement_service.go # Refund handling (existing)
│   ├── repository/
│   │   ├── payment_repository.go        # Payment DB operations
│   │   ├── webhook_repository.go        # Webhook DB operations
│   │   └── idempotency_repository.go    # Idempotency checks
│   ├── dto/
│   │   ├── payment.go                   # Payment DTOs
│   │   ├── webhook.go                   # Webhook DTOs
│   │   ├── bank.go                      # Bank DTOs (existing)
│   │   └── disbursement.go              # Disbursement DTOs (existing)
│   └── middleware/
│       └── webhook_auth.go              # Webhook signature verification
├── migrations/
│   └── init.js                          # MongoDB initialization
├── go.mod
└── Dockerfile
```

---

## Implementation Tasks

### Phase 1: Core Payment Infrastructure (Priority: HIGH)

#### Task 1.1: Configuration Setup

- [ ] Create `internal/config/config.go`
- [ ] Load Paystack keys from environment
- [ ] Configure MongoDB connection
- [ ] Configure RabbitMQ connection
- [ ] Add webhook secret configuration

#### Task 1.2: Paystack Client

- [ ] Create `internal/service/paystack_client.go`
- [ ] Implement `InitializeTransaction()`
- [ ] Implement `VerifyTransaction()`
- [ ] Implement `ListBanks()`
- [ ] Implement `ResolveAccountNumber()`
- [ ] Add error handling and retries
- [ ] Add request/response logging

#### Task 1.3: Repository Layer

- [ ] Create `internal/repository/payment_repository.go`
  - CreatePayment()
  - GetPaymentByID()
  - GetPaymentByReference()
  - UpdatePaymentStatus()
  - ListPaymentsByUser()
  - ListPaymentsByGoal()
- [ ] Create `internal/repository/webhook_repository.go`
  - SaveWebhookEvent()
  - MarkWebhookProcessed()
  - GetWebhookByEventID()
- [ ] Create `internal/repository/idempotency_repository.go`
  - CheckIdempotencyKey()
  - SaveIdempotencyKey()

### Phase 2: Payment Initialization (Priority: HIGH)

#### Task 2.1: Payment Service

- [ ] Create `internal/service/payment_service.go`
- [ ] Implement `InitializePayment()`
  - Validate request
  - Generate unique payment ID
  - Create payment record (INITIATED)
  - Call Paystack initialize API
  - Update payment with Paystack reference
  - Return authorization URL and access code
- [ ] Implement `GetPaymentStatus()`
- [ ] Implement `VerifyPayment()` (fallback)

#### Task 2.2: Payment Controller

- [ ] Create `internal/controller/payment_controller.go`
- [ ] Implement `POST /initialize`
- [ ] Implement `GET /{paymentId}/status`
- [ ] Implement `GET /{paymentId}/verify`
- [ ] Implement `GET /banks`
- [ ] Implement `GET /resolve-account`
- [ ] Add request validation
- [ ] Add error handling

### Phase 3: Webhook Processing (Priority: HIGH)

#### Task 3.1: Webhook Signature Verification

- [ ] Create `internal/middleware/webhook_auth.go`
- [ ] Implement signature verification using HMAC SHA512
- [ ] Validate webhook payload structure
- [ ] Add logging for failed verifications

#### Task 3.2: Webhook Service

- [ ] Create `internal/service/webhook_service.go`
- [ ] Implement `ProcessWebhook()`
  - Check idempotency (event ID)
  - Save webhook event
  - Extract payment reference
  - Update payment status
  - Emit PaymentVerified event
  - Mark webhook as processed
- [ ] Handle different event types:
  - `charge.success` (payment successful)
  - `charge.failed` (payment failed)
  - `transfer.success` (refund successful)
  - `transfer.failed` (refund failed)

#### Task 3.3: Webhook Controller

- [ ] Create `internal/controller/webhook_controller.go`
- [ ] Implement `POST /webhook`
- [ ] Apply signature verification middleware
- [ ] Return 200 OK immediately (async processing)

### Phase 4: Event Publishing (Priority: HIGH)

#### Task 4.1: Event Emitter

- [ ] Integrate RabbitMQ messaging
- [ ] Implement `EmitPaymentVerified()`
- [ ] Add retry logic for failed publishes
- [ ] Add dead letter queue for failed events

### Phase 5: Refund Disbursement (Priority: MEDIUM)

#### Task 5.1: Complete Refund Service

- [ ] Update `internal/service/refund_disbursement_service.go`
- [ ] Implement actual Paystack Transfer API calls
  - `CreateTransferRecipient()`
  - `InitiateTransfer()`
  - `VerifyTransfer()`
- [ ] Add transfer status tracking
- [ ] Emit `RefundCompleted` events

#### Task 5.2: Refund Webhook Handling

- [ ] Handle `transfer.success` webhook
- [ ] Handle `transfer.failed` webhook
- [ ] Update disbursement status
- [ ] Emit appropriate events

### Phase 6: Testing & Validation (Priority: HIGH)

#### Task 6.1: Unit Tests

- [ ] Test payment initialization
- [ ] Test webhook signature verification
- [ ] Test idempotency enforcement
- [ ] Test event publishing

#### Task 6.2: Integration Tests

- [ ] Test end-to-end payment flow
- [ ] Test webhook processing
- [ ] Test refund disbursement
- [ ] Test error scenarios

#### Task 6.3: Paystack Test Mode

- [ ] Use test API keys for development
- [ ] Test with Paystack test cards
- [ ] Verify webhook delivery in test mode

### Phase 7: Monitoring & Observability (Priority: MEDIUM)

#### Task 7.1: Datadog Metrics

- [ ] Track `payment.initialized.count`
- [ ] Track `payment.verified.count`
- [ ] Track `payment.failed.count`
- [ ] Track `webhook.received.count`
- [ ] Track `webhook.duplicate.count`
- [ ] Track `refund.initiated.count`
- [ ] Track `refund.completed.count`

#### Task 7.2: Logging

- [ ] Log all Paystack API calls
- [ ] Log webhook events
- [ ] Log payment state transitions
- [ ] Add correlation IDs for tracing

---

## Frontend Integration Guide

### 1. Install Paystack Inline

```html
<script src="https://js.paystack.co/v1/inline.js"></script>
```

### 2. Initialize Payment (React Example)

```javascript
const initiatePayment = async (goalId, amount) => {
  try {
    // Call backend to initialize payment
    const response = await fetch("/api/v1/payments/initialize", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        goal_id: goalId,
        amount: amount * 100, // Convert to kobo
        currency: "NGN",
        email: user.email,
        callback_url: `${window.location.origin}/payment/callback`,
        metadata: {
          goal_title: goalTitle,
          contributor_name: user.name,
        },
      }),
    });

    const data = await response.json();

    // Open Paystack inline checkout
    const handler = PaystackPop.setup({
      key: "pk_test_...", // Your Paystack public key
      email: user.email,
      amount: amount * 100,
      currency: "NGN",
      ref: data.data.reference,
      callback: function (response) {
        // Payment completed, verify on backend
        verifyPayment(data.data.payment_id);
      },
      onClose: function () {
        // User closed modal
        console.log("Payment cancelled");
      },
    });

    handler.openIframe();
  } catch (error) {
    console.error("Payment initialization failed:", error);
  }
};
```

### 3. Verify Payment

```javascript
const verifyPayment = async (paymentId) => {
  // Poll backend for payment status
  const maxAttempts = 10;
  let attempts = 0;

  const checkStatus = async () => {
    try {
      const response = await fetch(`/api/v1/payments/${paymentId}/status`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      const data = await response.json();

      if (data.data.status === "VERIFIED") {
        // Payment successful
        showSuccessMessage();
        redirectToGoalPage();
      } else if (data.data.status === "FAILED") {
        // Payment failed
        showErrorMessage();
      } else if (attempts < maxAttempts) {
        // Still processing, try again
        attempts++;
        setTimeout(checkStatus, 2000); // Poll every 2 seconds
      } else {
        // Timeout
        showTimeoutMessage();
      }
    } catch (error) {
      console.error("Status check failed:", error);
    }
  };

  checkStatus();
};
```

---

## Security Considerations

### 1. Webhook Signature Verification

```go
func VerifyWebhookSignature(payload []byte, signature string, secret string) bool {
    mac := hmac.New(sha512.New, []byte(secret))
    mac.Write(payload)
    expectedSignature := hex.EncodeToString(mac.Sum(nil))
    return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
```

### 2. Idempotency Enforcement

- Use Paystack event ID as idempotency key
- Store processed event IDs in database
- Reject duplicate webhook events
- Set TTL on idempotency keys (24 hours)

### 3. Amount Validation

- Always validate amounts on backend
- Never trust frontend amounts
- Verify amount matches goal contribution
- Check currency matches expected currency

### 4. API Key Security

- Store keys in environment variables
- Never commit keys to version control
- Use test keys in development
- Rotate keys periodically in production

---

## Environment Variables

Add to `.env`:

```bash
# Paystack Configuration
PAYSTACK_SECRET_KEY=sk_test_...
PAYSTACK_PUBLIC_KEY=pk_test_...
PAYSTACK_WEBHOOK_SECRET=your_webhook_secret

# MongoDB Configuration (Payments Service)
PAYMENTS_MONGODB_URI=mongodb://admin:admin123@localhost:27017/payments_db?authSource=admin

# Service Configuration
PAYMENTS_SERVICE_PORT=8081
PAYMENTS_SERVICE_ENV=development

# Webhook Configuration
WEBHOOK_TIMEOUT_SECONDS=30
WEBHOOK_MAX_RETRIES=3
```

---

## Testing Strategy

### 1. Paystack Test Cards

```
Successful Payment:
Card Number: 4084 0840 8408 4081
CVV: 408
Expiry: Any future date
PIN: 0000
OTP: 123456

Failed Payment:
Card Number: 5060 6666 6666 6666
```

### 2. Webhook Testing

- Use Paystack webhook testing tool in dashboard
- Use ngrok for local webhook testing
- Create mock webhook events for unit tests

### 3. Test Scenarios

- [ ] Successful payment flow
- [ ] Failed payment (insufficient funds)
- [ ] Duplicate webhook delivery
- [ ] Webhook signature mismatch
- [ ] Network timeout during payment
- [ ] User closes checkout modal
- [ ] Concurrent payment attempts

---

## Deployment Checklist

- [ ] Switch to production Paystack keys
- [ ] Configure production webhook URL
- [ ] Set up webhook endpoint in Paystack dashboard
- [ ] Enable webhook IP whitelisting
- [ ] Configure MongoDB indexes
- [ ] Set up Datadog monitoring
- [ ] Configure rate limiting
- [ ] Set up error alerting
- [ ] Document API endpoints
- [ ] Create runbook for payment issues

---

## Monitoring & Alerts

### Key Metrics to Track

1. Payment success rate (target: >95%)
2. Webhook processing time (target: <2s)
3. Failed payment rate
4. Duplicate webhook count
5. API error rate

### Alerts to Configure

1. Payment success rate drops below 90%
2. Webhook processing failures
3. High rate of failed payments
4. Paystack API errors
5. MongoDB connection failures

---

## Rollback Plan

If issues occur in production:

1. Switch to maintenance mode
2. Disable payment initialization endpoint
3. Continue processing pending webhooks
4. Investigate and fix issues
5. Re-enable payment initialization
6. Monitor closely for 24 hours

---

## Future Enhancements

1. **Multiple Payment Providers**: Add Flutterwave as backup
2. **Recurring Payments**: Support subscription-based contributions
3. **Payment Plans**: Allow installment payments
4. **Virtual Accounts**: Dedicated bank accounts per goal
5. **Payment Analytics**: Advanced reporting and insights
6. **Fraud Detection**: Implement fraud scoring
7. **Payment Reminders**: Notify users of pending payments

---

## References

- [Paystack API Documentation](https://paystack.com/docs/api/)
- [Paystack Inline Checkout](https://paystack.com/docs/payments/accept-payments/#popup)
- [Paystack Webhooks](https://paystack.com/docs/payments/webhooks/)
- [Paystack Transfers](https://paystack.com/docs/transfers/single-transfers/)
- [Paystack Test Cards](https://paystack.com/docs/payments/test-payments/)

---

## Summary

This implementation plan provides a complete roadmap for integrating Paystack into the GoalFund payment service. The key highlights are:

✅ **Webhook-based payment confirmation** (reliable and scalable)
✅ **In-app web checkout** support for seamless UX
✅ **Idempotent webhook processing** to prevent duplicates
✅ **Comprehensive error handling** and fallback mechanisms
✅ **Full refund disbursement** support
✅ **Production-ready security** measures
✅ **Monitoring and observability** with Datadog

The implementation is structured in phases, allowing for incremental development and testing. Start with Phase 1-4 for core payment functionality, then add refunds and monitoring in later phases.
