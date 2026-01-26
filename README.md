# GoFund – Product Requirements Document (PRD)

## 1. Product Overview

**Product Name:** GoFund
**Type:** Practice Fintech System (Group Funding + Accountability)
**Primary Goal:** Apply real-world fintech architecture, money correctness, and operational patterns using Go + React.

GoFund enables groups to contribute money toward shared goals with strong guarantees around **payment correctness, accounting integrity, transparency, and verifiable outcomes**.

This project is intentionally designed to mirror how real fintech systems are structured internally.

---

## 2. Problem Statement

Informal group funding systems (e.g., WhatsApp collections, screenshots, trust-based accounting) suffer from:

- No transparent contribution tracking
- No auditable money flow
- No proof that funds were used correctly
- No structured accountability

GoFund solves this by enforcing **ledger-based accounting**, **verified payments**, and **community-backed proof of completion**.

---

## 3. Non-Goals (Explicitly Out of Scope)

- No crypto or blockchain
- No lending, interest, or wallets
- No regulatory certification (KYC/AML only conceptual)
- No multi-currency settlement (single currency)

---

## 4. Core Product Features

### 4.1 Goal-Based Funding

- Users can create funding goals with:

  - Title, description, target amount
  - Deadline

- Goals progress is derived from ledger entries, not stored balances

### 4.2 Contributions

- Users initiate contributions to goals
- Contributions are pending until payment is verified

### 4.3 Payment Processing

- Integrate Paystack for payments
- Support:

  - Transaction initialization
  - Webhook verification
  - Idempotent processing

### 4.4 Ledger & Accounting

- All money movement is recorded as immutable ledger entries
- Double-entry accounting model
- Balances are computed, never stored

### 4.5 Withdrawals

- Goal owners can request withdrawals
- Withdrawals are ledger-backed and auditable

### 4.6 Proof of Accomplishment

- Goal owners submit proof after withdrawal
- Proofs require community verification

  - Minimum 3 confirmations
  - Or 5% of contributors (for large goals)

### 4.7 Real-Time Updates

- Contributors receive live updates on:

  - Contributions
  - Goal funding status
  - Proof verification

---

## 5. System Architecture (Microservices)

### 5.1 Service Overview

| Service               | Responsibility            | Core Purpose             |
| --------------------- | ------------------------- | ------------------------ |
| API Gateway           | Entry point & control     | Security & routing       |
| Payments Service      | External payment boundary | Paystack integration     |
| Ledger Service        | System of record          | Money correctness        |
| Goals Service         | Business rules            | Product logic            |
| Users Service         | Identity & access         | Authentication           |
| Notifications Service | Communication             | Realtime & async updates |

---

## 6. Service Definitions (Detailed)

### 6.1 API Gateway

**Purpose:** Control plane

**Responsibilities:**

- Authentication & authorization
- Rate limiting
- Request validation
- Trace propagation (Datadog)
- Routing to services

**Does NOT:**

- Execute business logic
- Access databases

---

### 6.2 Payments Service

**Purpose:** External money verification

**Responsibilities:**

- Paystack transaction initialization
- Webhook signature verification
- Idempotency enforcement
- Payment state machine
- Emit `PaymentVerified` events

**States:**
INITIATED → PENDING → VERIFIED → FAILED

**Database Tables:**

- payments
- webhook_events
- idempotency_keys

---

### 6.3 Ledger Service (Critical)

**Purpose:** Financial source of truth

**Responsibilities:**

- Immutable ledger entries
- Double-entry accounting
- Account management (user, goal, escrow)
- Balance computation
- Reconciliation support

**Database Tables:**

- ledger_entries
- accounts

**Rules:**

- Ledger entries are append-only
- No direct balance mutation

---

### 6.4 Goals Service

**Purpose:** Business logic & state machines

**Responsibilities:**

- Goal lifecycle management
- Contribution intent tracking
- Proof submission
- Voting & verification logic

**Goal States:**
OPEN → FUNDED → WITHDRAWN → PROOF_SUBMITTED → VERIFIED

**Database Tables:**

- goals
- contributions
- proofs
- votes

---

### 6.5 Users Service

**Purpose:** Identity infrastructure

**Responsibilities:**

- User accounts
- Authentication (JWT)
- Roles & permissions
- Optional MFA

**Database Tables:**

- users
- roles
- sessions

---

### 6.6 Notifications Service

**Purpose:** Async communication

**Responsibilities:**

- WebSocket connections
- Email / push notifications
- Event fan-out

**Consumes Events:**

- PaymentVerified
- GoalFunded
- ProofVerified

---

## 7. Event-Driven Communication

**Transport:** NATS (or internal event bus)

**Core Events:**

- PaymentVerified
- LedgerEntryCreated
- GoalFunded
- ProofSubmitted
- ProofVerified

**Rules:**

- Services never mutate other services’ databases
- Events are idempotent

---

## 8. Monitoring & Observability (Datadog)

### Tools Used:

- Datadog APM
- Distributed Tracing
- Log Management (JSON logs)
- Custom Metrics
- Dashboards & Alerts

### Key Metrics:

- payment.success.count
- payment.failure.count
- ledger.entries.created
- webhook.duplicate.count
- goal.funded.count

---

## 9. Technology Stack

### Backend

- Go
- Gin (standardized across services)
- PostgreSQL
- Redis (idempotency & caching)
- NATS (events)

### Frontend

- React
- WebSockets

### Infrastructure

- Docker & Docker Compose
- API Gateway (Go-based or Nginx initially)

### Monitoring

- Datadog (APM, logs, metrics, alerts)

---

## 10. Security Principles

- HTTPS everywhere
- Signed webhooks
- JWT authentication
- Role-based access control
- No trust in frontend data

---

## 11. Success Criteria

You succeed if:

- No duplicate payment can credit a goal twice
- All balances are derivable from ledger entries
- Payment failures do not corrupt internal state
- Every financial action is auditable
- You can trace a payment end-to-end in Datadog

---

## 12. Learning Outcomes

By completing GoFund, you will understand:

- Fintech ledger design
- Idempotent payment processing
- Event-driven systems
- Operational monitoring
- Microservice boundaries

This is the real win of the project.
