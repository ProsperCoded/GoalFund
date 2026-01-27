#!/bin/bash

# Atlas schema inspection script for GoFund

set -e

echo "🔍 Inspecting database schema..."

# Check if Atlas is installed
if ! command -v atlas &> /dev/null; then
    echo "❌ Atlas CLI is not installed. Please install it first:"
    echo "   curl -sSf https://atlasgo.sh | sh"
    exit 1
fi

# Set environment
ENV=${1:-dev}
echo "📝 Using environment: $ENV"

# Inspect schema
echo "🔄 Inspecting current schema..."
atlas schema inspect --env $ENV

echo "✅ Schema inspection completed!"