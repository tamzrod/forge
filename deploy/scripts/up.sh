#!/bin/bash
# up.sh - Start the application
# Usage: ./deploy/scripts/up.sh
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$DEPLOY_DIR")"

cd "$PROJECT_DIR"

echo "Starting application..."

# Check if docker compose is available
if ! command -v docker compose &> /dev/null; then
    echo "Error: docker compose not found"
    exit 1
fi

# Start services
docker compose -f deploy/docker/docker-compose.yml up -d

echo "Application started."
echo "Run './deploy/scripts/health.sh' to verify."
