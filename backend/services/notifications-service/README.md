# Notifications Service

The Notifications Service handles all notification-related functionality for the GoFund platform, including email notifications and notification persistence.

## Features

- **Email Notifications**: Send transactional emails via SMTP
- **Notification Persistence**: Store notification history in PostgreSQL
- **Event-Driven**: Consume events from RabbitMQ
- **User Preferences**: Manage user notification preferences
- **REST API**: Query and manage notifications

## Architecture

```
┌─────────────────────────────────────────────────────┐
│           Notifications Service                      │
├─────────────────────────────────────────────────────┤
│                                                       │
│  ┌──────────────┐      ┌──────────────┐            │
│  │   RabbitMQ   │──────▶│Event Handler │            │
│  │   Consumer   │      │   (Router)   │            │
│  └──────────────┘      └──────┬───────┘            │
│                               │                     │
│                        ┌──────▼────────┐           │
│                        │ Notification  │           │
│                        │   Service     │           │
│                        └──────┬────────┘           │
│                               │                     │
│                        ┌──────▼────────┐           │
│                        │ Email Service │           │
│                        └───────────────┘           │
│                                                     │
│    ┌──────────────────────────────────────────┐   │
│    │      PostgreSQL Database                 │   │
│    │  - notifications                         │   │
│    │  - notification_preferences              │   │
│    └──────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
```

## Events Consumed

The service listens to the following RabbitMQ events:

| Event                        | Description                       | Recipients           |
| ---------------------------- | --------------------------------- | -------------------- |
| `PaymentVerified`            | Payment successfully verified     | Contributor          |
| `ContributionConfirmed`      | Contribution confirmed            | Goal Owner           |
| `WithdrawalRequested`        | Withdrawal requested              | Goal Owner           |
| `WithdrawalCompleted`        | Withdrawal completed              | Goal Owner           |
| `ProofSubmitted`             | Proof of accomplishment submitted | Contributors         |
| `ProofVoted`                 | Vote cast on proof                | Goal Owner           |
| `GoalFunded`                 | Goal reached target               | Owner & Contributors |
| `UserSignedUp`               | New user registered               | New User             |
| `PasswordResetRequested`     | Password reset requested          | User                 |
| `EmailVerificationRequested` | Email verification requested      | User                 |
| `KYCVerified`                | KYC verification completed        | User                 |

## API Endpoints

### Notifications

- `GET /api/v1/notifications` - Get user notifications (paginated)
- `GET /api/v1/notifications/:id` - Get specific notification
- `PUT /api/v1/notifications/:id/read` - Mark notification as read
- `DELETE /api/v1/notifications/:id` - Delete notification
- `GET /api/v1/notifications/unread/count` - Get unread count

### Preferences

- `GET /api/v1/notifications/preferences` - Get user preferences
- `PUT /api/v1/notifications/preferences` - Update preferences

### Health

- `GET /api/v1/notifications/health` - Health check

## Environment Variables

```env
# Server
PORT=8085

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=notifications_db

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_EXCHANGE=gofund_events
RABBITMQ_QUEUE=notifications_queue

# Email (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@gofund.com
SMTP_FROM_NAME=GoFund

# Datadog
DD_SERVICE=notifications-service
DD_ENV=dev
DD_VERSION=1.0.0
```

## Database Schema

### notifications

| Column              | Type         | Description                    |
| ------------------- | ------------ | ------------------------------ |
| id                  | UUID         | Primary key                    |
| user_id             | UUID         | User who receives notification |
| type                | VARCHAR(50)  | Notification type              |
| title               | VARCHAR(255) | Notification title             |
| message             | TEXT         | Notification message           |
| data                | JSONB        | Additional event data          |
| email_sent          | BOOLEAN      | Whether email was sent         |
| email_sent_at       | TIMESTAMP    | When email was sent            |
| email_failed_reason | TEXT         | Reason for email failure       |
| retry_count         | INT          | Number of retry attempts       |
| is_read             | BOOLEAN      | Whether notification was read  |
| read_at             | TIMESTAMP    | When notification was read     |
| created_at          | TIMESTAMP    | Creation timestamp             |
| updated_at          | TIMESTAMP    | Last update timestamp          |

### notification_preferences

| Column                     | Type      | Description                        |
| -------------------------- | --------- | ---------------------------------- |
| id                         | UUID      | Primary key                        |
| user_id                    | UUID      | User ID (unique)                   |
| email_enabled              | BOOLEAN   | Email notifications enabled        |
| payment_notifications      | BOOLEAN   | Payment notifications enabled      |
| contribution_notifications | BOOLEAN   | Contribution notifications enabled |
| withdrawal_notifications   | BOOLEAN   | Withdrawal notifications enabled   |
| proof_notifications        | BOOLEAN   | Proof notifications enabled        |
| goal_notifications         | BOOLEAN   | Goal notifications enabled         |
| marketing_emails           | BOOLEAN   | Marketing emails enabled           |
| created_at                 | TIMESTAMP | Creation timestamp                 |
| updated_at                 | TIMESTAMP | Last update timestamp              |

## Running the Service

### Local Development

1. Set up environment variables:

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

2. Run database migrations:

   ```bash
   psql -U postgres -d notifications_db -f migrations/001_create_notifications_table.sql
   ```

3. Install dependencies:

   ```bash
   go mod download
   ```

4. Run the service:
   ```bash
   go run cmd/main.go
   ```

### Docker

```bash
docker build -t notifications-service .
docker run -p 8085:8085 --env-file .env notifications-service
```

## Testing

```bash
# Run unit tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Monitoring

The service integrates with Datadog for:

- **APM Tracing**: Distributed tracing across services
- **Metrics**: Custom metrics for notifications sent, failed, etc.
- **Logs**: Structured JSON logging

### Key Metrics

- `notification.created.count` - Notifications created
- `notification.email.sent.count` - Emails sent successfully
- `notification.email.failed.count` - Email failures
- `event.consumed.count` - Events consumed from RabbitMQ
- `event.processing.duration` - Event processing time

## Email Configuration

### Gmail

1. Enable 2-factor authentication
2. Generate an app-specific password
3. Use the app password in `SMTP_PASSWORD`

### SendGrid

```env
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=your-sendgrid-api-key
```

## Development Notes

- All notifications are persisted before email sending
- Email sending is asynchronous (goroutine)
- Failed emails are retried with exponential backoff
- User preferences are checked before sending emails
- Default preferences are created on user signup

## Future Enhancements

- [ ] WebSocket support for real-time notifications
- [ ] Push notifications (FCM/APNs)
- [ ] SMS notifications
- [ ] Notification templates with i18n
- [ ] Batch email processing
- [ ] Email delivery tracking
