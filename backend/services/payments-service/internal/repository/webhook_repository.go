package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gofund/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// WebhookRepository handles webhook event database operations
type WebhookRepository struct {
	collection *mongo.Collection
}

// NewWebhookRepository creates a new webhook repository
func NewWebhookRepository(db *mongo.Database) *WebhookRepository {
	return &WebhookRepository{
		collection: db.Collection("webhook_events"),
	}
}

// SaveWebhookEvent saves a webhook event
func (r *WebhookRepository) SaveWebhookEvent(ctx context.Context, event *models.WebhookEvent) error {
	event.ReceivedAt = time.Now()
	event.Processed = false

	result, err := r.collection.InsertOne(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to save webhook event: %w", err)
	}

	event.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetWebhookByEventID retrieves a webhook event by its event ID
func (r *WebhookRepository) GetWebhookByEventID(ctx context.Context, eventID string) (*models.WebhookEvent, error) {
	var event models.WebhookEvent
	err := r.collection.FindOne(ctx, bson.M{"eventId": eventID}).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Not found, not an error
		}
		return nil, fmt.Errorf("failed to get webhook event: %w", err)
	}
	return &event, nil
}

// IsEventProcessed checks if a webhook event has been processed
func (r *WebhookRepository) IsEventProcessed(ctx context.Context, eventID string) bool {
	var event models.WebhookEvent
	err := r.collection.FindOne(ctx, bson.M{"eventId": eventID, "processed": true}).Decode(&event)
	return err == nil
}

// MarkWebhookProcessed marks a webhook event as processed
func (r *WebhookRepository) MarkWebhookProcessed(ctx context.Context, eventID string) error {
	now := time.Now()
	filter := bson.M{"eventId": eventID}
	update := bson.M{
		"$set": bson.M{
			"processed":   true,
			"processedAt": now,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to mark webhook as processed: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("webhook event not found")
	}

	return nil
}

// ListUnprocessedWebhooks retrieves unprocessed webhook events
func (r *WebhookRepository) ListUnprocessedWebhooks(ctx context.Context, limit int64) ([]*models.WebhookEvent, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "receivedAt", Value: 1}}).
		SetLimit(limit)

	cursor, err := r.collection.Find(ctx, bson.M{"processed": false}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list unprocessed webhooks: %w", err)
	}
	defer cursor.Close(ctx)

	var events []*models.WebhookEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("failed to decode webhook events: %w", err)
	}

	return events, nil
}

// CountWebhooksByEvent counts webhooks by event type
func (r *WebhookRepository) CountWebhooksByEvent(ctx context.Context, eventType string) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"event": eventType})
	if err != nil {
		return 0, fmt.Errorf("failed to count webhooks: %w", err)
	}
	return count, nil
}

// EnsureIndexes creates necessary indexes for the webhook_events collection
func (r *WebhookRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
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
			Keys: bson.D{{Key: "receivedAt", Value: -1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}
