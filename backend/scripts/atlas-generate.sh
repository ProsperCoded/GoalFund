#!/bin/bash

# Atlas migration generation script for GoalFund

set -e

echo "🚀 Generating Atlas migrations for GoalFund..."

# Check if Atlas is installed
if ! command -v atlas &> /dev/null; then
    echo "❌ Atlas CLI is not installed. Please install it first:"
    echo "   curl -sSf https://atlasgo.sh | sh"
    exit 1
fi

# Set environment
ENV=${1:-dev}
echo "📝 Using environment: $ENV"

# Generate migration
echo "🔄 Generating migration from schema..."
atlas migrate diff --env $ENV

echo "✅ Migration generated successfully!"
echo ""
echo "📋 Next steps:"
echo "   1. Review the generated migration files in migrations/"
echo "   2. Apply migrations: ./scripts/atlas-apply.sh $ENV"
echo "   3. Verify schema: ./scripts/atlas-inspect.sh $ENV"