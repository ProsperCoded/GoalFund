# KYC System Deployment Checklist

## Pre-Deployment

### Database

- [ ] Review migration file: `backend/services/users-service/migrations/001_add_kyc_fields.sql`
- [ ] Backup users database before running migration
- [ ] Run migration on development environment first
- [ ] Verify new columns exist: `nin`, `kyc_verified`, `kyc_verified_at`
- [ ] Verify indexes created: `idx_users_nin`, `idx_users_kyc_verified_at`

### Code Review

- [ ] Review all modified files (9 files)
- [ ] Review all new files (4 files)
- [ ] Verify no breaking changes to existing endpoints
- [ ] Check error handling in KYC service
- [ ] Verify NIN masking implementation
- [ ] Review event publishing logic

### Testing

- [ ] Test NIN submission with valid 11-digit number
- [ ] Test NIN submission with invalid format
- [ ] Test duplicate NIN prevention
- [ ] Test already-verified user flow
- [ ] Test unauthorized access (no token)
- [ ] Verify NIN masking in responses
- [ ] Test KYC status endpoint
- [ ] Verify KYC fields in login response
- [ ] Verify KYC fields in register response
- [ ] Verify KYC fields in profile endpoints
- [ ] Test event emission (check RabbitMQ)
- [ ] Verify Datadog metrics tracking

### Documentation

- [ ] Review README.md updates
- [ ] Review KYC_API.md documentation
- [ ] Review KYC_IMPLEMENTATION_SUMMARY.md
- [ ] Review KYC_QUICK_REFERENCE.md
- [ ] Update API documentation (Swagger/Postman)

## Deployment Steps

### 1. Database Migration

```bash
# Connect to users database
psql -U postgres -d users_db

# Run migration
\i backend/services/users-service/migrations/001_add_kyc_fields.sql

# Verify columns
\d users

# Check indexes
\di idx_users_nin
\di idx_users_kyc_verified_at
```

### 2. Build & Deploy Service

```bash
# Navigate to users service
cd backend/services/users-service

# Build the service
go build -o users-service ./cmd

# Or build with Docker
docker build -t gofund/users-service:latest .

# Deploy to your environment
# (specific steps depend on your deployment method)
```

### 3. Environment Variables

Ensure these are set (if needed):

```bash
# Database
USERS_DB_HOST=localhost
USERS_DB_PORT=5432
USERS_DB_USER=postgres
USERS_DB_PASSWORD=your_password
USERS_DB_NAME=users_db

# JWT
JWT_SECRET=your-secret-key

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_EXCHANGE=gofund.events

# Datadog (optional)
DD_SERVICE=users-service
DD_ENV=production
DD_VERSION=1.0.0
```

### 4. Service Restart

```bash
# Restart users service
systemctl restart users-service

# Or with Docker
docker-compose restart users-service

# Verify service is running
curl http://localhost:8084/health
```

### 5. Verify Deployment

```bash
# Check health endpoint
curl http://localhost:8084/health

# Test KYC status endpoint (with valid token)
curl -H "Authorization: Bearer TOKEN" \
  http://localhost:8084/users/kyc/status

# Check logs for errors
tail -f /var/log/users-service.log
# Or with Docker
docker logs -f users-service
```

## Post-Deployment

### Monitoring

- [ ] Check Datadog dashboard for `user.kyc.verified.count` metric
- [ ] Monitor error rates in logs
- [ ] Check RabbitMQ for `KYCVerified` events
- [ ] Monitor API response times
- [ ] Set up alerts for KYC failures

### Validation

- [ ] Create test user account
- [ ] Submit test NIN
- [ ] Verify auto-approval works
- [ ] Check NIN masking in response
- [ ] Verify event published to RabbitMQ
- [ ] Check Datadog metric incremented
- [ ] Test duplicate NIN prevention
- [ ] Verify profile endpoints return KYC status

### Documentation

- [ ] Update team wiki/docs with KYC feature
- [ ] Share API documentation with frontend team
- [ ] Update Postman/Swagger collections
- [ ] Notify stakeholders of new feature

## Rollback Plan

If issues occur:

### 1. Rollback Code

```bash
# Revert to previous version
git revert HEAD
git push

# Redeploy previous version
# (specific steps depend on deployment method)
```

### 2. Rollback Database (if needed)

```sql
-- Remove KYC columns (only if necessary)
ALTER TABLE users DROP COLUMN IF EXISTS nin;
ALTER TABLE users DROP COLUMN IF EXISTS kyc_verified;
ALTER TABLE users DROP COLUMN IF EXISTS kyc_verified_at;

-- Drop indexes
DROP INDEX IF EXISTS idx_users_nin;
DROP INDEX IF EXISTS idx_users_kyc_verified_at;
```

⚠️ **Warning**: Only rollback database if no KYC data has been collected. Otherwise, you'll lose user verification data.

## Known Limitations

- ✅ **Dummy Implementation**: Auto-approves all NIN submissions
- ✅ **No External Verification**: Doesn't check against government databases
- ✅ **Email Verification**: Tracked but not yet implemented
- ✅ **No Admin Dashboard**: No UI for viewing KYC statistics

## Future Enhancements

### Phase 2 (Recommended)

- [ ] Implement actual NIN verification API integration
- [ ] Add email verification flow
- [ ] Create admin dashboard for KYC management
- [ ] Add document upload support
- [ ] Implement verification expiry

### Phase 3 (Advanced)

- [ ] Multi-level verification (basic, advanced)
- [ ] Biometric verification
- [ ] Address verification
- [ ] Business/corporate KYC
- [ ] International ID support

## Support & Troubleshooting

### Common Issues

**Issue**: Migration fails

- **Solution**: Check database permissions, verify syntax

**Issue**: NIN not masked in response

- **Solution**: Check `maskNIN()` function in kyc_service.go

**Issue**: Events not publishing

- **Solution**: Verify RabbitMQ connection, check event service initialization

**Issue**: Duplicate NIN not detected

- **Solution**: Verify `GetUserByNIN()` repository method, check database index

### Logs to Check

```bash
# Service logs
tail -f /var/log/users-service.log

# Database logs
tail -f /var/log/postgresql/postgresql.log

# RabbitMQ logs
tail -f /var/log/rabbitmq/rabbit@hostname.log

# Datadog APM
# Check Datadog dashboard for traces and errors
```

## Sign-off

- [ ] Technical Lead approval
- [ ] QA testing completed
- [ ] Security review completed
- [ ] Documentation reviewed
- [ ] Deployment plan approved
- [ ] Rollback plan tested
- [ ] Monitoring configured
- [ ] Team trained on new feature

---

**Deployment Date**: ********\_********

**Deployed By**: ********\_********

**Version**: 1.0.0

**Status**: ⬜ Pending | ⬜ In Progress | ⬜ Completed | ⬜ Rolled Back
