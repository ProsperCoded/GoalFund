# Database Schema Design

## Overview

This document outlines the complete database schema for all GoFund services, including PostgreSQL tables and MongoDB collections.

## Entity Relationship Diagram

```mermaid
erDiagram
    %% Users Service (PostgreSQL)
    users {
        uuid id PK
        string email UK
        string username UK
        string password_hash
        string first_name
        string last_name
        string phone
        boolean email_verified
        boolean phone_verified
        timestamp created_at
        timestamp updated_at
    }

    roles {
        uuid id PK
        string name UK
        string description
        json permissions
        timestamp created_at
    }

    user_roles {
        uuid user_id FK
        uuid role_id FK
        timestamp assigned_at
    }

    sessions {
        uuid id PK
        uuid user_id FK
        string token_hash UK
        timestamp expires_at
        json metadata
        timestamp created_at
    }

    %% Goals Service (PostgreSQL)
    goals {
        uuid id PK
        uuid owner_id FK
        string title
        text description
        bigint target_amount
        string currency
        timestamp deadline
        string status
        timestamp created_at
        timestamp updated_at
    }

    contributions {
        uuid id PK
        uuid goal_id FK
        uuid user_id FK
        uuid payment_id FK
        bigint amount
        string currency
        string status
        timestamp created_at
        timestamp updated_at
    }

    proofs {
        uuid id PK
        uuid goal_id FK
        uuid submitted_by FK
        string title
        text description
        json media_urls
        string status
        timestamp submitted_at
        timestamp verified_at
    }

    votes {
        uuid id PK
        uuid proof_id FK
        uuid voter_id FK
        boolean is_approved
        text comment
        timestamp voted_at
    }

    %% Ledger Service (PostgreSQL)
    accounts {
        uuid id PK
        string account_type
        uuid entity_id
        string currency
        timestamp created_at
    }

    ledger_entries {
        uuid id PK
        uuid account_id FK
        uuid transaction_id
        string entry_type
        bigint amount
        string currency
        string description
        json metadata
        timestamp created_at
    }

    %% Relationships
    users ||--o{ user_roles : has
    roles ||--o{ user_roles : assigned_to
    users ||--o{ sessions : owns
    users ||--o{ goals : creates
    users ||--o{ contributions : makes
    users ||--o{ proofs : submits
    users ||--o{ votes : casts
    goals ||--o{ contributions : receives
    goals ||--o{ proofs : has
    proofs ||--o{ votes : receives
    accounts ||--o{ ledger_entries : contains
    contributions ||--|| ledger_entries : triggers
```

## MongoDB Collections (Payments Service)

```mermaid
erDiagram
    %% Payments Service (MongoDB)
    payments {
        ObjectId _id PK
        string paymentId UK
        string paystackReference UK
        string userId FK
        string goalId FK
        long amount
        string currency
        string status
        object paystackData
        date createdAt
        date updatedAt
    }

    webhook_events {
        ObjectId _id PK
        string eventId UK
        string event
        object data
        string signature
        boolean processed
        date receivedAt
        date processedAt
    }

    idempotency_keys {
        ObjectId _id PK
        string key UK
        string paymentId FK
        date createdAt
        date expiresAt
    }

    %% Virtual relationships (cross-service)
    payments ||--o{ webhook_events : triggers
    payments ||--|| idempotency_keys : uses
```

## Cross-Service Relationships

```mermaid
graph TD
    %% Users Service
    U[Users Service - PostgreSQL]
    U_users[users table]
    
    %% Goals Service  
    G[Goals Service - PostgreSQL]
    G_goals[goals table]
    G_contributions[contributions table]
    
    %% Payments Service
    P[Payments Service - MongoDB]
    P_payments[payments collection]
    
    %% Ledger Service
    L[Ledger Service - PostgreSQL]
    L_accounts[accounts table]
    L_entries[ledger_entries table]
    
    %% Cross-service relationships
    U_users -.->|owner_id| G_goals
    U_users -.->|user_id| G_contributions
    U_users -.->|userId| P_payments
    U_users -.->|entity_id| L_accounts
    
    G_goals -.->|goalId| P_payments
    G_goals -.->|entity_id| L_accounts
    G_contributions -.->|payment_id| P_payments
    
    P_payments -.->|triggers| L_entries
    
    style U fill:#e1f5fe
    style G fill:#f3e5f5
    style P fill:#e8f5e8
    style L fill:#fff3e0
```

## Detailed Table Specifications

### Users Service Tables

#### users
- **Purpose**: Store user account information
- **Key Constraints**: 
  - Unique email and username
  - Password must be hashed
  - Email/phone verification flags

#### roles & user_roles
- **Purpose**: Role-based access control (RBAC)
- **Key Features**:
  - Many-to-many relationship between users and roles
  - JSON permissions for flexible authorization

#### sessions
- **Purpose**: Authentication session management
- **Key Features**:
  - Token-based authentication
  - Automatic expiration
  - Metadata for device tracking

### Goals Service Tables

#### goals
- **Purpose**: Core goal definitions
- **Key Features**:
  - Monetary targets with currency support
  - Status tracking (OPEN, FUNDED, WITHDRAWN, etc.)
  - Deadline enforcement

#### contributions
- **Purpose**: Track user contributions to goals
- **Key Features**:
  - Links to payment records via payment_id
  - Status tracking for contribution lifecycle

#### proofs & votes
- **Purpose**: Community verification system
- **Key Features**:
  - Proof submission with media support
  - Democratic voting mechanism
  - Approval thresholds

### Ledger Service Tables

#### accounts
- **Purpose**: Account management for double-entry bookkeeping
- **Account Types**: USER, GOAL, ESCROW, REVENUE
- **Key Features**:
  - Multi-currency support
  - Entity polymorphism (users, goals, etc.)

#### ledger_entries
- **Purpose**: Immutable financial transaction log
- **Entry Types**: DEBIT, CREDIT
- **Key Features**:
  - Append-only design
  - Rich metadata for audit trails
  - Balance calculation via aggregation

### Payments Service Collections

#### payments
- **Purpose**: Payment processing and status tracking
- **Key Features**:
  - Paystack integration data
  - Payment state machine
  - Cross-service references (userId, goalId)

#### webhook_events
- **Purpose**: Webhook processing and deduplication
- **Key Features**:
  - Raw Paystack webhook storage
  - Processing status tracking
  - Signature verification data

#### idempotency_keys
- **Purpose**: Prevent duplicate payment processing
- **Key Features**:
  - TTL expiration (24 hours)
  - Unique constraint enforcement
  - Payment correlation

## Data Flow Example

1. **User creates goal**: `users.id` → `goals.owner_id`
2. **User contributes**: Creates `contributions` record with `payment_id`
3. **Payment processed**: `payments` collection stores Paystack data
4. **Payment verified**: Creates `ledger_entries` for double-entry bookkeeping
5. **Goal funded**: Status updated when target reached
6. **Proof submitted**: `proofs` table with community voting
7. **Funds released**: Additional `ledger_entries` for withdrawal

This schema ensures data consistency, audit trails, and proper separation of concerns across microservices.