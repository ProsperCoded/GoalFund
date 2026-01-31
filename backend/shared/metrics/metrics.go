package metrics

import (
	"fmt"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// Client is the global DogStatsD client for custom metrics
	Client *statsd.Client
)

// InitDatadog initializes Datadog tracing and metrics
func InitDatadog(serviceName, env, version string) error {
	// Start the tracer with configuration
	tracer.Start(
		tracer.WithServiceName(serviceName),
		tracer.WithEnv(env),
		tracer.WithServiceVersion(version),
		tracer.WithRuntimeMetrics(),
	)

	// Initialize DogStatsD client for custom metrics
	var err error
	Client, err = statsd.New("datadog-agent:8125",
		statsd.WithNamespace("gofund."),
		statsd.WithTags([]string{
			fmt.Sprintf("service:%s", serviceName),
			fmt.Sprintf("env:%s", env),
			fmt.Sprintf("version:%s", version),
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to create statsd client: %w", err)
	}

	return nil
}

// StopDatadog gracefully stops Datadog tracing and metrics
func StopDatadog() {
	tracer.Stop()
	if Client != nil {
		Client.Close()
	}
}

// Business Metrics Helper Functions

// IncrementCounter increments a counter metric with optional tags
func IncrementCounter(metric string, tags ...string) {
	if Client != nil {
		Client.Incr(metric, tags, 1)
	}
}

// RecordHistogram records a histogram value with optional tags
func RecordHistogram(metric string, value float64, tags ...string) {
	if Client != nil {
		Client.Histogram(metric, value, tags, 1)
	}
}

// RecordGauge records a gauge value with optional tags
func RecordGauge(metric string, value float64, tags ...string) {
	if Client != nil {
		Client.Gauge(metric, value, tags, 1)
	}
}

// RecordDuration records the duration of an operation
func RecordDuration(metric string, start time.Time, tags ...string) {
	duration := time.Since(start).Seconds()
	RecordHistogram(metric, duration, tags...)
}

// Payment Metrics

// TrackPaymentInitiated tracks when a payment is initiated
func TrackPaymentInitiated(amount float64, currency string) {
	IncrementCounter("payment.initiated.count", fmt.Sprintf("currency:%s", currency))
	RecordHistogram("payment.amount", amount, fmt.Sprintf("currency:%s", currency), "status:initiated")
}

// TrackPaymentSuccess tracks successful payment
func TrackPaymentSuccess(amount float64, currency string, duration time.Duration) {
	IncrementCounter("payment.success.count", fmt.Sprintf("currency:%s", currency))
	RecordHistogram("payment.amount", amount, fmt.Sprintf("currency:%s", currency), "status:success")
	RecordHistogram("payment.processing.duration", duration.Seconds(), fmt.Sprintf("currency:%s", currency), "status:success")
}

// TrackPaymentFailure tracks failed payment
func TrackPaymentFailure(reason, currency string) {
	IncrementCounter("payment.failure.count", fmt.Sprintf("currency:%s", currency), fmt.Sprintf("reason:%s", reason))
}

// TrackWebhookDuplicate tracks duplicate webhook attempts
func TrackWebhookDuplicate(eventType string) {
	IncrementCounter("payment.webhook.duplicate.count", fmt.Sprintf("event_type:%s", eventType))
}

// TrackWebhookProcessed tracks processed webhooks
func TrackWebhookProcessed(eventType string, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	IncrementCounter("payment.webhook.processed.count", fmt.Sprintf("event_type:%s", eventType), fmt.Sprintf("status:%s", status))
}

// Goal Metrics

// TrackGoalCreated tracks new goal creation
func TrackGoalCreated(targetAmount float64, currency string) {
	IncrementCounter("goal.created.count", fmt.Sprintf("currency:%s", currency))
	RecordHistogram("goal.target_amount", targetAmount, fmt.Sprintf("currency:%s", currency))
}

// TrackGoalFunded tracks when a goal reaches its target
func TrackGoalFunded(goalID string, finalAmount float64, currency string, durationDays int) {
	IncrementCounter("goal.funded.count", fmt.Sprintf("currency:%s", currency))
	RecordHistogram("goal.funded.amount", finalAmount, fmt.Sprintf("currency:%s", currency))
	RecordHistogram("goal.funded.duration_days", float64(durationDays), fmt.Sprintf("currency:%s", currency))
}

// TrackContribution tracks individual contributions
func TrackContribution(amount float64, currency string, goalID string) {
	IncrementCounter("goal.contribution.count", fmt.Sprintf("currency:%s", currency))
	RecordHistogram("goal.contribution.amount", amount, fmt.Sprintf("currency:%s", currency))
}

// TrackWithdrawalRequested tracks withdrawal requests
func TrackWithdrawalRequested(amount float64, currency string) {
	IncrementCounter("goal.withdrawal.requested.count", fmt.Sprintf("currency:%s", currency))
	RecordHistogram("goal.withdrawal.amount", amount, fmt.Sprintf("currency:%s", currency))
}

// TrackProofSubmitted tracks proof submissions
func TrackProofSubmitted(goalID string) {
	IncrementCounter("goal.proof.submitted.count", fmt.Sprintf("goal_id:%s", goalID))
}

// TrackProofVerified tracks verified proofs
func TrackProofVerified(goalID string, verificationCount int) {
	IncrementCounter("goal.proof.verified.count", fmt.Sprintf("goal_id:%s", goalID))
	RecordGauge("goal.proof.verification_count", float64(verificationCount), fmt.Sprintf("goal_id:%s", goalID))
}

// Ledger Metrics

// TrackLedgerEntryCreated tracks new ledger entries
func TrackLedgerEntryCreated(accountType string, amount float64, currency string) {
	IncrementCounter("ledger.entry.created.count", fmt.Sprintf("account_type:%s", accountType), fmt.Sprintf("currency:%s", currency))
	RecordHistogram("ledger.entry.amount", amount, fmt.Sprintf("account_type:%s", accountType), fmt.Sprintf("currency:%s", currency))
}

// TrackBalanceComputation tracks balance calculation performance
func TrackBalanceComputation(accountID string, duration time.Duration, entryCount int) {
	RecordHistogram("ledger.balance.computed.duration", duration.Seconds(), fmt.Sprintf("account_id:%s", accountID))
	RecordGauge("ledger.balance.entry_count", float64(entryCount), fmt.Sprintf("account_id:%s", accountID))
}

// TrackAccountCreated tracks new account creation
func TrackAccountCreated(accountType string) {
	IncrementCounter("ledger.account.created.count", fmt.Sprintf("account_type:%s", accountType))
}

// TrackReconciliationMismatch tracks balance reconciliation mismatches
func TrackReconciliationMismatch(accountID string, expectedBalance, actualBalance float64) {
	IncrementCounter("ledger.reconciliation.mismatch", fmt.Sprintf("account_id:%s", accountID))
	RecordGauge("ledger.reconciliation.expected_balance", expectedBalance, fmt.Sprintf("account_id:%s", accountID))
	RecordGauge("ledger.reconciliation.actual_balance", actualBalance, fmt.Sprintf("account_id:%s", accountID))
}

// User Metrics

// TrackUserRegistration tracks new user registrations
func TrackUserRegistration() {
	IncrementCounter("user.registration.count")
}

// TrackLoginSuccess tracks successful login attempts
func TrackLoginSuccess(userID string) {
	IncrementCounter("user.login.success.count", fmt.Sprintf("user_id:%s", userID))
}

// TrackLoginFailure tracks failed login attempts
func TrackLoginFailure(reason string) {
	IncrementCounter("user.login.failure.count", fmt.Sprintf("reason:%s", reason))
}

// TrackSessionCreated tracks new session creation
func TrackSessionCreated(userID string) {
	IncrementCounter("user.session.created.count", fmt.Sprintf("user_id:%s", userID))
}

// TrackJWTIssued tracks JWT token issuance
func TrackJWTIssued(tokenType string) {
	IncrementCounter("user.jwt.issued.count", fmt.Sprintf("token_type:%s", tokenType))
}

// Event/Messaging Metrics

// TrackEventPublished tracks event publishing
func TrackEventPublished(eventType string, success bool, duration time.Duration) {
	status := "success"
	if !success {
		status = "failure"
	}
	IncrementCounter("event.published.count", fmt.Sprintf("event_type:%s", eventType), fmt.Sprintf("status:%s", status))
	RecordHistogram("event.publish.duration", duration.Seconds(), fmt.Sprintf("event_type:%s", eventType))
}

// TrackEventConsumed tracks event consumption
func TrackEventConsumed(eventType string, success bool, processingDuration time.Duration, eventAge time.Duration) {
	status := "success"
	if !success {
		status = "failure"
	}
	IncrementCounter("event.consumed.count", fmt.Sprintf("event_type:%s", eventType), fmt.Sprintf("status:%s", status))
	RecordHistogram("event.processing.duration", processingDuration.Seconds(), fmt.Sprintf("event_type:%s", eventType))
	RecordHistogram("event.processing.age", eventAge.Seconds(), fmt.Sprintf("event_type:%s", eventType))
}

// Cache Metrics

// TrackCacheHit tracks cache hits
func TrackCacheHit(cacheKey string) {
	IncrementCounter("cache.hit.count", fmt.Sprintf("key:%s", cacheKey))
}

// TrackCacheMiss tracks cache misses
func TrackCacheMiss(cacheKey string) {
	IncrementCounter("cache.miss.count", fmt.Sprintf("key:%s", cacheKey))
}

// TrackCacheOperation tracks cache operation duration
func TrackCacheOperation(operation string, duration time.Duration) {
	RecordHistogram("cache.operation.duration", duration.Seconds(), fmt.Sprintf("operation:%s", operation))
}
