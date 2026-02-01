package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gofund/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IdempotencyRepository handles idempotency key operations
type IdempotencyRepository struct {
	collection *mongo.Collection
}

// NewIdempotencyRepository creates a new idempotency repository
func NewIdempotencyRepository(db *mongo.Database) *IdempotencyRepository {
	return &IdempotencyRepository{
		collection: db.Collection("idempotency_keys"),
	}
}

// CheckIdempotencyKey checks if an idempotency key exists
func (r *IdempotencyRepository) CheckIdempotencyKey(ctx context.Context, key string) (bool, string, error) {
	var idempotency models.IdempotencyKey
	err := r.collection.FindOne(ctx, bson.M{"key": key}).Decode(&idempotency)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, "", nil // Key doesn't exist
		}
		return false, "", fmt.Errorf("failed to check idempotency key: %w", err)
	}

	// Key exists, return the associated payment ID
	return true, idempotency.PaymentID, nil
}

// SaveIdempotencyKey saves an idempotency key
func (r *IdempotencyRepository) SaveIdempotencyKey(ctx context.Context, key, paymentID string) error {
	idempotency := &models.IdempotencyKey{
		Key:       key,
		PaymentID: paymentID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // Expire after 24 hours
	}

	_, err := r.collection.InsertOne(ctx, idempotency)
	if err != nil {
		return fmt.Errorf("failed to save idempotency key: %w", err)
	}

	return nil
}

// DeleteIdempotencyKey deletes an idempotency key
func (r *IdempotencyRepository) DeleteIdempotencyKey(ctx context.Context, key string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"key": key})
	if err != nil {
		return fmt.Errorf("failed to delete idempotency key: %w", err)
	}
	return nil
}

// EnsureIndexes creates necessary indexes for the idempotency_keys collection
func (r *IdempotencyRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "key", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			// TTL index to automatically delete expired keys
			Keys:    bson.D{{Key: "expiresAt", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}
