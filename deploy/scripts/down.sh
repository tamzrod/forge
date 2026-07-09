#!/bin/bash
# down.sh - Stop the application
# Usage: ./deploy/scripts/down.sh
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$DEPLOY_DIR")"

cd "$PROJECT_DIR"

echo "Stopping application..."

if ! command -v docker compose &> /dev/null; then
    echo "Error: docker compose not found"
    exit 1
fi

docker compose -f deploy/docker/docker-compose.yml down

echo "Application stopped."
