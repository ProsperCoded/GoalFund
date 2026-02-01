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

// PaymentRepository handles payment database operations
type PaymentRepository struct {
	collection *mongo.Collection
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *mongo.Database) *PaymentRepository {
	return &PaymentRepository{
		collection: db.Collection("payments"),
	}
}

// CreatePayment creates a new payment record
func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, payment)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	payment.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetPaymentByID retrieves a payment by its payment ID
func (r *PaymentRepository) GetPaymentByID(ctx context.Context, paymentID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.collection.FindOne(ctx, bson.M{"paymentId": paymentID}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	return &payment, nil
}

// GetPaymentByReference retrieves a payment by Paystack reference
func (r *PaymentRepository) GetPaymentByReference(ctx context.Context, reference string) (*models.Payment, error) {
	var payment models.Payment
	err := r.collection.FindOne(ctx, bson.M{"paystackReference": reference}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	return &payment, nil
}

// UpdatePayment updates an existing payment
func (r *PaymentRepository) UpdatePayment(ctx context.Context, payment *models.Payment) error {
	payment.UpdatedAt = time.Now()

	filter := bson.M{"paymentId": payment.PaymentID}
	update := bson.M{
		"$set": bson.M{
			"status":             payment.Status,
			"paystackReference":  payment.PaystackReference,
			"paystackData":       payment.PaystackData,
			"updatedAt":          payment.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// UpdatePaymentStatus updates only the payment status
func (r *PaymentRepository) UpdatePaymentStatus(ctx context.Context, paymentID string, status models.PaymentStatus) error {
	filter := bson.M{"paymentId": paymentID}
	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// ListPaymentsByUser retrieves all payments for a user
func (r *PaymentRepository) ListPaymentsByUser(ctx context.Context, userID string, limit, offset int64) ([]*models.Payment, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetLimit(limit).
		SetSkip(offset)

	cursor, err := r.collection.Find(ctx, bson.M{"userId": userID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	defer cursor.Close(ctx)

	var payments []*models.Payment
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, fmt.Errorf("failed to decode payments: %w", err)
	}

	return payments, nil
}

// ListPaymentsByGoal retrieves all payments for a goal
func (r *PaymentRepository) ListPaymentsByGoal(ctx context.Context, goalID string, limit, offset int64) ([]*models.Payment, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetLimit(limit).
		SetSkip(offset)

	cursor, err := r.collection.Find(ctx, bson.M{"goalId": goalID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	defer cursor.Close(ctx)

	var payments []*models.Payment
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, fmt.Errorf("failed to decode payments: %w", err)
	}

	return payments, nil
}

// CountPaymentsByStatus counts payments by status
func (r *PaymentRepository) CountPaymentsByStatus(ctx context.Context, status models.PaymentStatus) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"status": status})
	if err != nil {
		return 0, fmt.Errorf("failed to count payments: %w", err)
	}
	return count, nil
}

// EnsureIndexes creates necessary indexes for the payments collection
func (r *PaymentRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
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
			Keys: bson.D{{Key: "createdAt", Value: -1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}
