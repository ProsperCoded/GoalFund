// MongoDB initialization script for payments service
// This script creates the necessary collections and indexes

db = db.getSiblingDB("payments_db");

// Create collections
db.createCollection("payments");
db.createCollection("webhook_events");
db.createCollection("idempotency_keys");

// Create indexes for payments collection
db.payments.createIndex({ paymentId: 1 }, { unique: true });
db.payments.createIndex(
  { paystackReference: 1 },
  { unique: true, sparse: true }
);
db.payments.createIndex({ userId: 1 });
db.payments.createIndex({ goalId: 1 });
db.payments.createIndex({ status: 1 });
db.payments.createIndex({ createdAt: -1 });

// Create indexes for webhook_events collection
db.webhook_events.createIndex({ eventId: 1 }, { unique: true });
db.webhook_events.createIndex({ event: 1 });
db.webhook_events.createIndex({ processed: 1 });
db.webhook_events.createIndex({ receivedAt: -1 });

// Create indexes for idempotency_keys collection
db.idempotency_keys.createIndex({ key: 1 }, { unique: true });
db.idempotency_keys.createIndex({ expiresAt: 1 }, { expireAfterSeconds: 0 }); // TTL index

print("MongoDB initialization complete for payments service");
