# Datadog Monitoring Setup for GoFund

This document explains how Datadog monitoring has been integrated into the GoFund microservices project.

## Overview

GoFund now includes comprehensive Datadog monitoring with:

- **APM (Application Performance Monitoring)** - Distributed tracing across all services
- **Custom Business Metrics** - Financial correctness and KPI tracking
- **Infrastructure Monitoring** - Database, cache, and message queue metrics
- **Log Management** - Centralized logging with correlation
- **Real-time Dashboards** - Pre-configured monitoring views

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Datadog Platform                         │
│  (APM, Metrics, Logs, Dashboards, Alerts)                   │
└─────────────────────────────────────────────────────────────┘
                            ▲
                            │
                    ┌───────┴────────┐
                    │ Datadog Agent  │
                    │  (Container)   │
                    └───────┬────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
    ┌───▼───┐          ┌────▼────┐        ┌────▼────┐
    │ Users │          │ Payments│        │ Goals   │
    │Service│          │ Service │        │ Service │
    └───┬───┘          └────┬────┘        └────┬────┘
        │                   │                   │
    ┌───▼───┐          ┌────▼────┐        ┌────▼────┐
    │Postgres│         │ MongoDB │        │Postgres │
    └────────┘         └─────────┘        └─────────┘
```

## Setup Instructions

### 1. Get Datadog API Key

1. Sign up for Datadog at https://www.datadoghq.com/
2. Navigate to **Organization Settings** → **API Keys**
3. Create a new API key or copy an existing one
4. Note your Datadog site (US: `datadoghq.com`, EU: `datadoghq.eu`)

### 2. Configure Environment Variables

Create a `.env` file in the `backend` directory:

```bash
# Copy the example file
cp .env.example .env
```

Update the Datadog configuration in `.env`:

```bash
# Datadog Configuration
DD_API_KEY=your-actual-datadog-api-key-here
DD_SITE=datadoghq.com  # or datadoghq.eu for EU

# Environment and Version
DD_ENV=dev  # Options: dev, staging, production
DD_VERSION=1.0.0
```

### 3. Start the Services

```bash
# Start all services including Datadog agent
docker-compose up -d

# Verify Datadog agent is running
docker logs gofund-datadog-agent

# Check service logs for Datadog initialization
docker logs gofund-users-service | grep Datadog
```

### 4. Verify Datadog Integration

1. Go to https://app.datadoghq.com/
2. Navigate to **APM** → **Services** - You should see all 5 services
3. Navigate to **Metrics** → **Explorer** - Search for `gofund.*` metrics
4. Navigate to **Logs** → **Live Tail** - View real-time logs

## Monitored Metrics

### Business Metrics

#### Payment Metrics

- `gofund.payment.initiated.count` - Payment attempts
- `gofund.payment.success.count` - Successful payments
- `gofund.payment.failure.count` - Failed payments (tagged by reason)
- `gofund.payment.webhook.duplicate.count` - Prevented duplicate webhooks
- `gofund.payment.amount` - Payment amounts (histogram)
- `gofund.payment.processing.duration` - Payment processing time

#### Goal Metrics

- `gofund.goal.created.count` - New goals created
- `gofund.goal.funded.count` - Goals that reached target
- `gofund.goal.contribution.count` - Individual contributions
- `gofund.goal.withdrawal.requested.count` - Withdrawal requests
- `gofund.goal.proof.submitted.count` - Proof submissions
- `gofund.goal.proof.verified.count` - Verified proofs
- `gofund.goal.target_amount` - Goal target amounts
- `gofund.goal.funded.duration_days` - Time to reach funding goal

#### Ledger Metrics

- `gofund.ledger.entry.created.count` - New ledger entries (tagged by account type)
- `gofund.ledger.account.created.count` - New accounts
- `gofund.ledger.balance.computed.duration` - Balance calculation performance
- `gofund.ledger.reconciliation.mismatch` - Balance mismatches (should be 0!)

#### User Metrics

- `gofund.user.registration.count` - New user registrations
- `gofund.user.login.success.count` - Successful logins
- `gofund.user.login.failure.count` - Failed login attempts (tagged by reason)
- `gofund.user.session.created.count` - New sessions
- `gofund.user.jwt.issued.count` - JWT tokens issued

### Infrastructure Metrics

#### Event/Messaging Metrics

- `gofund.event.published.count` - Events published to RabbitMQ
- `gofund.event.consumed.count` - Events consumed
- `gofund.event.publish.duration` - Event publishing time
- `gofund.event.processing.duration` - Event processing time
- `gofund.event.processing.age` - Event lag (time between publish and process)

#### Cache Metrics

- `gofund.cache.hit.count` - Cache hits
- `gofund.cache.miss.count` - Cache misses
- `gofund.cache.operation.duration` - Cache operation time

### Automatic APM Metrics

The following are automatically collected by Datadog APM:

- HTTP request duration (by endpoint, status code)
- HTTP request count (by endpoint, method)
- HTTP error rate (4xx, 5xx)
- Database query duration
- Database query count
- Database connection pool metrics
- Service dependencies and call graphs

## Using Custom Metrics in Your Code

### Example: Track Payment Success

```go
import "github.com/gofund/shared/metrics"

func ProcessPayment(amount float64, currency string) error {
    start := time.Now()

    // Your payment processing logic
    err := paymentProvider.Process(amount, currency)

    duration := time.Since(start)

    if err != nil {
        metrics.TrackPaymentFailure("processing_error", currency)
        return err
    }

    metrics.TrackPaymentSuccess(amount, currency, duration)
    return nil
}
```

### Example: Track Goal Funding

```go
import "github.com/gofund/shared/metrics"

func CheckGoalFunding(goal *Goal) {
    if goal.CurrentAmount >= goal.TargetAmount {
        durationDays := int(time.Since(goal.CreatedAt).Hours() / 24)
        metrics.TrackGoalFunded(goal.ID, goal.CurrentAmount, goal.Currency, durationDays)
    }
}
```

### Example: Track Ledger Entry

```go
import "github.com/gofund/shared/metrics"

func CreateLedgerEntry(accountType string, amount float64, currency string) error {
    // Create ledger entry
    err := ledgerRepo.Create(entry)

    if err == nil {
        metrics.TrackLedgerEntryCreated(accountType, amount, currency)
    }

    return err
}
```

## Dashboards

### Creating Custom Dashboards

1. Go to **Dashboards** → **New Dashboard**
2. Add widgets for key metrics:
   - **Timeseries**: `gofund.payment.success.count` vs `gofund.payment.failure.count`
   - **Query Value**: `sum:gofund.payment.amount{*}`
   - **Heatmap**: `gofund.payment.processing.duration`
   - **Top List**: `gofund.user.login.failure.count` by reason

### Recommended Dashboard Widgets

#### Financial Health Dashboard

- Payment success rate (success / (success + failure))
- Total payment volume
- Payment processing duration (p50, p95, p99)
- Webhook duplicate prevention count
- Ledger reconciliation mismatches

#### User Activity Dashboard

- User registrations (timeseries)
- Login success vs failure
- Active sessions
- JWT token issuance rate

#### Goal Performance Dashboard

- Goals created vs funded
- Average time to funding
- Contribution count
- Withdrawal requests
- Proof verification rate

## Alerts

### Recommended Alerts

#### Critical Alerts

1. **Ledger Reconciliation Mismatch**

   - Metric: `gofund.ledger.reconciliation.mismatch`
   - Threshold: > 0
   - Action: Page on-call engineer immediately

2. **Payment Failure Rate**

   - Metric: `gofund.payment.failure.count / (gofund.payment.success.count + gofund.payment.failure.count)`
   - Threshold: > 10% over 5 minutes
   - Action: Alert payment team

3. **Service Error Rate**
   - Metric: APM error rate
   - Threshold: > 5% over 5 minutes
   - Action: Alert engineering team

#### Warning Alerts

1. **Event Processing Lag**

   - Metric: `gofund.event.processing.age`
   - Threshold: > 30 seconds
   - Action: Notify engineering team

2. **Database Query Performance**

   - Metric: Database query duration p95
   - Threshold: > 500ms
   - Action: Notify database team

3. **Cache Miss Rate**
   - Metric: `gofund.cache.miss.count / (gofund.cache.hit.count + gofund.cache.miss.count)`
   - Threshold: > 50%
   - Action: Review caching strategy

## Troubleshooting

### Datadog Agent Not Connecting

```bash
# Check agent status
docker exec gofund-datadog-agent agent status

# Check agent logs
docker logs gofund-datadog-agent

# Verify API key
docker exec gofund-datadog-agent agent config | grep api_key
```

### No Traces Appearing

1. Verify services are sending traces:

   ```bash
   docker logs gofund-users-service | grep "Datadog initialized"
   ```

2. Check APM configuration:

   ```bash
   docker exec gofund-datadog-agent agent status | grep -A 10 "APM Agent"
   ```

3. Verify DD_AGENT_HOST is set correctly in service containers

### Metrics Not Showing Up

1. Check if metrics are being sent:

   ```bash
   docker logs gofund-users-service | grep -i metric
   ```

2. Verify DogStatsD is listening:

   ```bash
   docker exec gofund-datadog-agent agent status | grep -A 5 "DogStatsD"
   ```

3. Check for metric name typos in code

## Performance Considerations

### Sampling

For high-traffic production environments, consider adjusting trace sampling:

```go
// In metrics.InitDatadog()
tracer.Start(
    tracer.WithServiceName(serviceName),
    tracer.WithSampleRate(0.1), // Sample 10% of traces
)
```

### Metric Aggregation

Custom metrics are aggregated on the Datadog side. No need to pre-aggregate in your application.

### Cost Optimization

- Use tags wisely (avoid high-cardinality tags like user IDs)
- Sample traces in production (keep 100% in dev/staging)
- Set retention policies for logs
- Use metric monitors instead of log-based monitors when possible

## Next Steps

1. **Create Dashboards**: Build dashboards for your team's specific needs
2. **Set Up Alerts**: Configure alerts for critical business metrics
3. **Add More Metrics**: Instrument additional business logic as needed
4. **Monitor Costs**: Track your Datadog usage and optimize as needed
5. **Train Team**: Ensure your team knows how to use Datadog effectively

## Resources

- [Datadog APM Documentation](https://docs.datadoghq.com/tracing/)
- [Datadog Go Tracer](https://docs.datadoghq.com/tracing/setup_overview/setup/go/)
- [DogStatsD](https://docs.datadoghq.com/developers/dogstatsd/)
- [Datadog Dashboards](https://docs.datadoghq.com/dashboards/)
- [Datadog Monitors](https://docs.datadoghq.com/monitors/)
