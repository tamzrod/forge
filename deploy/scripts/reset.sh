#!/bin/bash
# reset.sh - Factory reset: remove runtime data, preserve configuration
# Usage: ./deploy/scripts/reset.sh
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$DEPLOY_DIR")"

echo "WARNING: This will remove runtime data but preserve configuration."
echo ""
echo "Removing:"
echo "  - databases"
echo "  - caches"
echo "  - temporary files"
echo "  - logs"
echo ""
echo "Preserving:"
echo "  - Configuration files"
echo "  - Source code"
echo ""

read -p "Continue? [y/N] " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 0
fi

cd "$PROJECT_DIR"

echo "Stopping services..."
"$SCRIPT_DIR/down.sh" 2>/dev/null || true

echo "Removing runtime data..."

# Remove data directory contents
rm -rf "$DEPLOY_DIR/data"/*
rm -rf "$DEPLOY_DIR/logs"/*

# Remove Docker volumes (runtime data only)
docker volume rm forge_runtime-data 2>/dev/null || true

# Recreate empty directories
mkdir -p "$DEPLOY_DIR/data"
mkdir -p "$DEPLOY_DIR/logs"

echo "Reset complete."
echo "Run './deploy/scripts/up.sh' to start fresh."
