#!/bin/bash
# rebuild.sh - Clean rebuild: stop, rebuild, restart
# Usage: ./deploy/scripts/rebuild.sh
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Rebuilding application..."

"$SCRIPT_DIR/down.sh"
echo "Building images..."
docker compose -f deploy/docker/docker-compose.yml build --no-cache
docker compose -f deploy/docker/docker-compose.yml up -d --remove-orphans

echo "Rebuild complete."
echo "Run './deploy/scripts/health.sh' to verify."
