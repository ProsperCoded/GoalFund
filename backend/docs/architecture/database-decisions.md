# Database Architecture Decisions

## Overview

GoFund uses a polyglot persistence approach, selecting the optimal database technology for each service's specific requirements.

## Database Selection by Service

### 1. Payments Service → MongoDB

**Decision:** Use MongoDB for payment processing and webhook management.

**Rationale:**
- **JSON-native**: Paystack webhooks are JSON payloads - natural document fit
- **Flexible schema**: Payment providers frequently change webhook formats
- **Simple operations**: Primarily inserts, lookups by payment ID, and basic queries
- **Idempotency**: Easy upserts using `_id` field for duplicate prevention
- **Rapid development**: No schema migrations needed for webhook format changes

**Collections:**
- `payments` - Payment records with status, amounts, references
- `webhook_events` - Raw webhook payloads from Paystack
- `idempotency_keys` - Prevent duplicate webhook processing

### 2. Ledger Service → PostgreSQL

**Decision:** Use PostgreSQL for financial ledger and accounting.

**Rationale:**
- **ACID compliance**: Critical for financial data integrity
- **Aggregation performance**: Highly optimized `SUM()` operations for balance calculations
- **Mature indexing**: B-tree indexes excel at numerical range queries
- **Window functions**: Support complex financial calculations (running balances, reconciliation)
- **Proven in fintech**: Industry standard for financial systems
- **Append-only optimization**: PostgreSQL handles write-heavy workloads efficiently

**Tables:**
- `ledger_entries` - Immutable double-entry accounting records
- `accounts` - Account metadata (user, goal, escrow accounts)

### 3. Goals Service → PostgreSQL

**Decision:** Use PostgreSQL for goal management and business logic.

**Rationale:**
- **Relational data**: Goals have complex relationships (contributors, proofs, votes)
- **ACID transactions**: Goal state changes must be atomic
- **Complex queries**: Reporting and analytics require SQL capabilities
- **Data integrity**: Foreign key constraints ensure referential integrity

**Tables:**
- `goals` - Goal definitions and metadata
- `contributions` - Contribution tracking
- `proofs` - Proof of accomplishment submissions
- `votes` - Community verification votes

### 4. Users Service → PostgreSQL

**Decision:** Use PostgreSQL for user management and authentication.

**Rationale:**
- **Structured data**: User profiles have consistent schema
- **Security**: Mature authentication patterns and constraints
- **Indexing**: Efficient lookups by email, username, etc.
- **ACID compliance**: User operations must be consistent

**Tables:**
- `users` - User profiles and credentials
- `roles` - Role-based access control
- `sessions` - Authentication sessions

## Performance Considerations

### Balance Calculation Example

**PostgreSQL (Ledger Service):**
```sql
SELECT SUM(amount) as balance 
FROM ledger_entries 
WHERE account_id = 'user_123' 
  AND created_at <= '2024-01-01';
```
- **Performance**: ~1ms for millions of records with proper indexing
- **Consistency**: ACID guarantees accurate results

**MongoDB Alternative:**
```javascript
db.ledger_entries.aggregate([
  { $match: { account_id: "user_123", created_at: { $lte: new Date("2024-01-01") } } },
  { $group: { _id: null, balance: { $sum: "$amount" } } }
])
```
- **Performance**: ~10-50ms for same dataset
- **Consistency**: Eventually consistent by default

## Trade-offs Accepted

1. **Operational complexity**: Multiple database technologies require different expertise
2. **Data consistency**: Cross-service transactions require eventual consistency patterns
3. **Backup strategies**: Different backup/restore procedures for each database type

## Benefits Gained

1. **Optimized performance**: Each service uses the best tool for its workload
2. **Development velocity**: Teams can choose appropriate abstractions
3. **Scalability**: Independent scaling strategies per service
4. **Fault isolation**: Database issues don't cascade across all services

## Migration Strategy

- **Phase 1**: Start with current architecture (proven technologies)
- **Phase 2**: Monitor performance bottlenecks and query patterns
- **Phase 3**: Consider NoSQL alternatives only if PostgreSQL becomes a constraint

This polyglot approach balances performance optimization with operational simplicity, ensuring each service can excel at its core responsibilities.