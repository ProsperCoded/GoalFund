package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

// MongoConfig holds MongoDB configuration
type MongoConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(cfg MongoConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI(cfg.URI)
	
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(cfg.Database)
	return db, nil
}

// SetupMongoIndexes creates necessary indexes for MongoDB collections
func SetupMongoIndexes(db *mongo.Database) error {
	ctx := context.Background()

	// Payments collection indexes
	paymentsCollection := db.Collection("payments")
	paymentsIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "paymentId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "paystackReference", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
		{
			Keys: bson.D{{Key: "userId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "goalId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "createdAt", Value: 1}},
		},
	}

	if _, err := paymentsCollection.Indexes().CreateMany(ctx, paymentsIndexes); err != nil {
		return fmt.Errorf("failed to create payments indexes: %w", err)
	}

	// Webhook events collection indexes
	webhookCollection := db.Collection("webhook_events")
	webhookIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "eventId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "event", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "processed", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "receivedAt", Value: 1}},
		},
	}

	if _, err := webhookCollection.Indexes().CreateMany(ctx, webhookIndexes); err != nil {
		return fmt.Errorf("failed to create webhook_events indexes: %w", err)
	}

	// Idempotency keys collection indexes
	idempotencyCollection := db.Collection("idempotency_keys")
	idempotencyIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "key", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "expiresAt", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0), // TTL index
		},
	}

	if _, err := idempotencyCollection.Indexes().CreateMany(ctx, idempotencyIndexes); err != nil {
		return fmt.Errorf("failed to create idempotency_keys indexes: %w", err)
	}

	// Payment attempts collection indexes
	attemptsCollection := db.Collection("payment_attempts")
	attemptsIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "paymentId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "createdAt", Value: 1}},
		},
	}

	if _, err := attemptsCollection.Indexes().CreateMany(ctx, attemptsIndexes); err != nil {
		return fmt.Errorf("failed to create payment_attempts indexes: %w", err)
	}

	// Payment requests collection indexes
	requestsCollection := db.Collection("payment_requests")
	requestsIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "paymentId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "userId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "goalId", Value: 1}},
		},
		{
			Keys:    bson.D{{Key: "expiresAt", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0), // TTL index
		},
	}

	if _, err := requestsCollection.Indexes().CreateMany(ctx, requestsIndexes); err != nil {
		return fmt.Errorf("failed to create payment_requests indexes: %w", err)
	}

	return nil
}

// SetupPaymentsDatabase initializes the payments database with indexes
func SetupPaymentsDatabase(cfg MongoConfig) (*mongo.Database, error) {
	db, err := NewMongoDB(cfg)
	if err != nil {
		return nil, err
	}

	// Setup indexes
	if err := SetupMongoIndexes(db); err != nil {
		return nil, err
	}

	return db, nil
}