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

  - Title, description, target amount (initial milestone)
  - Optional deadline

- **Continuous Funding Model:**

  - Goals can be funded continuously beyond the target amount
  - Multiple withdrawal cycles are supported
  - Goals remain OPEN for contributions unless owner explicitly closes them
  - Target amount serves as an initial milestone, not a hard cap

- **Milestone-Based Progress Tracking:**

  - Goals can be broken down into multiple milestones
  - Each milestone has its own target amount and description
  - Support for **recurring milestones** (e.g., semester tuition, monthly rent)
  - Recurring types: WEEKLY, MONTHLY, SEMESTER, YEARLY
  - Withdrawals can be tied to milestone completion
  - Proofs can be submitted per milestone for granular transparency

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

- Goal owners can request withdrawals at any time once funds are available
- **Bank account details required** - owners must provide bank information for disbursement
- Bank details can be added during goal creation or updated later
- **No verification required before withdrawal** - owners have direct access to contributed funds
- Multiple withdrawals are supported (continuous funding model)
- Withdrawals are ledger-backed and fully auditable
- Goal can continue receiving funds after withdrawal (unless closed by owner)
- Withdrawals can be tied to specific milestone completion

### 4.6 Proof of Accomplishment & Community Feedback

- Goal owners can submit proof **after withdrawal** (optional but encouraged)
- Proofs serve as **transparency and accountability** mechanism (not a withdrawal gate)
- Community voting reflects **satisfaction level** with how funds were used:
  - Contributors vote TRUE (satisfied) or FALSE (not satisfied)
  - Voting thresholds: Minimum 3 votes OR 5% of contributors
  - Votes are visible to all contributors for transparency
- **Key Point:** Voting does NOT block or reverse withdrawals - it's purely for reputation and trust-building

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
| API Gateway (Nginx)   | Entry point & control     | Security & routing       |
| Payments Service      | External payment boundary | Paystack integration     |
| Ledger Service        | System of record          | Money correctness        |
| Goals Service         | Business rules            | Product logic            |
| Users Service         | Identity & access         | Authentication           |
| Notifications Service | Communication             | Realtime & async updates |

---

## 6. Service Definitions (Detailed)

### 6.1 API Gateway (Nginx)

**Purpose:** Control plane & reverse proxy

**Technology:** Nginx (not Go-based service)

**Responsibilities:**

- HTTP reverse proxy & load balancing
- Rate limiting (10 req/s general, 5 req/s auth)
- Request routing to microservices
- SSL termination & security headers
- WebSocket support for notifications
- Static file serving (if needed)

**Routing:**

- `/api/v1/users/*` → Users Service (port 8084)
- `/api/v1/goals/*` → Goals Service (port 8083)
- `/api/v1/ledger/*` → Ledger Service (port 8082)
- `/api/v1/payments/*` → Payments Service (port 8081)
- `/api/v1/notifications/*` → Notifications Service (port 8085)

**Does NOT:**

- Execute business logic
- Access databases
- Handle authentication (delegated to services)

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
- **Milestone management** (including recurring milestones)
- Continuous funding tracking
- Contribution intent tracking
- Withdrawal request handling (with bank account validation)
- Proof submission
- Community voting & feedback logic

**Goal States:**

- **OPEN** - Accepting contributions (default state)
- **CLOSED** - Owner has stopped accepting new contributions
- **CANCELLED** - Goal was cancelled

**Key Behaviors:**

- Goals can receive unlimited contributions (continuous funding)
- Withdrawals can happen multiple times while still OPEN
- Owner can transition OPEN → CLOSED at any time
- Milestones can be one-time or recurring (WEEKLY, MONTHLY, SEMESTER, YEARLY)
- Bank account details required for withdrawals

**Database Tables:**

- goals (with bank account fields)
- **milestones** (supports recurring)
- contributions (with milestone_id reference)
- **withdrawals** (with bank details snapshot)
- proofs (with milestone_id reference)
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
- ContributionConfirmed
- WithdrawalRequested
- ProofVoted

---

## 7. Event-Driven Communication

**Transport:** RabbitMQ

**Core Events:**

- **PaymentVerified** - Emitted by Payments Service when payment succeeds
- **ContributionConfirmed** - Emitted by Goals Service after processing payment
- **WithdrawalRequested** - Emitted by Goals Service when owner requests withdrawal
- **WithdrawalCompleted** - Emitted by Ledger Service after successful withdrawal
- **ProofSubmitted** - Emitted by Goals Service when proof is submitted
- **ProofVoted** - Emitted when a contributor casts a vote on proof

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
- RabbitMQ (events & queues)

### Frontend

- React
- WebSockets

### Infrastructure

- Docker & Docker Compose
- API Gateway (Nginx reverse proxy)

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
