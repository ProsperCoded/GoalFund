#!/bin/bash

# Development database setup script for GoFund

set -e

echo "🚀 Setting up GoFund development databases..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running. Please start Docker first.${NC}"
    exit 1
fi

echo -e "${YELLOW}📦 Starting infrastructure services...${NC}"

# Start infrastructure services
docker-compose up -d postgres-ledger postgres-goals postgres-users mongodb-payments redis rabbitmq

echo -e "${YELLOW}⏳ Waiting for databases to be ready...${NC}"

# Wait for PostgreSQL databases
for i in {1..30}; do
    if docker-compose exec -T postgres-ledger pg_isready -U postgres > /dev/null 2>&1 && \
       docker-compose exec -T postgres-goals pg_isready -U postgres > /dev/null 2>&1 && \
       docker-compose exec -T postgres-users pg_isready -U postgres > /dev/null 2>&1; then
        echo -e "${GREEN}✅ PostgreSQL databases are ready!${NC}"
        break
    fi
    echo "Waiting for PostgreSQL... ($i/30)"
    sleep 2
done

# Wait for MongoDB
for i in {1..30}; do
    if docker-compose exec -T mongodb-payments mongosh --eval "db.adminCommand('ping')" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ MongoDB is ready!${NC}"
        break
    fi
    echo "Waiting for MongoDB... ($i/30)"
    sleep 2
done

echo -e "${YELLOW}🔄 Running database migrations...${NC}"

# Apply Atlas migrations (if Atlas is installed)
if command -v atlas &> /dev/null; then
    echo "Applying Atlas migrations..."
    atlas migrate apply --env dev || echo "Atlas migrations failed or no migrations to apply"
else
    echo -e "${YELLOW}⚠️  Atlas CLI not found. Skipping schema migrations.${NC}"
    echo -e "${YELLOW}   Install Atlas: curl -sSf https://atlasgo.sh | sh${NC}"
fi

echo -e "${GREEN}✅ Development databases are ready!${NC}"
echo ""
echo -e "${YELLOW}📋 Database URLs:${NC}"
echo "   PostgreSQL (Users):    localhost:5435/users_db"
echo "   PostgreSQL (Goals):    localhost:5434/goals_db"
echo "   PostgreSQL (Ledger):   localhost:5433/ledger_db"
echo "   MongoDB (Payments):    localhost:27017/payments_db"
echo "   Redis:                 localhost:6379"
echo "   RabbitMQ Management:   http://localhost:15672 (guest/guest)"
echo ""
echo -e "${GREEN}🚀 You can now start your services!${NC}"