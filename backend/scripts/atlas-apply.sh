#!/bin/bash

# Atlas migration application script for GoFund

set -e

echo "🚀 Applying Atlas migrations for GoFund..."

# Check if Atlas is installed
if ! command -v atlas &> /dev/null; then
    echo "❌ Atlas CLI is not installed. Please install it first:"
    echo "   curl -sSf https://atlasgo.sh | sh"
    exit 1
fi

# Set environment
ENV=${1:-dev}
echo "📝 Using environment: $ENV"

# Apply migrations
echo "🔄 Applying migrations..."
atlas migrate apply --env $ENV

echo "✅ Migrations applied successfully!"
echo ""
echo "📋 Next steps:"
echo "   1. Verify schema: ./scripts/atlas-inspect.sh $ENV"
echo "   2. Start your services"