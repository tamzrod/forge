#!/bin/bash
# shell.sh - Open shell inside runtime container
# Usage: ./deploy/scripts/shell.sh [command]
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$DEPLOY_DIR")"

cd "$PROJECT_DIR"

CONTAINER=$(docker ps --format '{{.Names}}' | grep -E '(runtime|app|forge)' | head -1)

if [[ -z "$CONTAINER" ]]; then
    echo "Error: Runtime container not found. Is the application running?"
    echo "Run './deploy/scripts/up.sh' first."
    exit 1
fi

if [[ -n "$1" ]]; then
    docker exec -it "$CONTAINER" "$@"
else
    docker exec -it "$CONTAINER" sh
fi
