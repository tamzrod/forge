#!/bin/bash
# nuke.sh - Destroy everything: containers, images, volumes, generated files
# Usage: ./deploy/scripts/nuke.sh
# WARNING: Requires typing "nuke" to confirm. This is destructive.
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$DEPLOY_DIR")"

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║                    ! DANGER !                           ║"
echo "║                                                            ║"
echo "║  This will permanently destroy:                            ║"
echo "║    - All containers                                        ║"
echo "║    - All images (except base)                             ║"
echo "║    - All volumes                                          ║"
echo "║    - All runtime data                                    ║"
echo "║    - All logs                                            ║"
echo "║    - All backups                                         ║"
echo "║                                                            ║"
echo "║  This cannot be undone.                                 ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
read -p "Type 'nuke' to confirm: " confirmation

if [[ "$confirmation" != "nuke" ]]; then
    echo "Aborted. You must type 'nuke' exactly."
    exit 1
fi

cd "$PROJECT_DIR"

echo ""
echo "Destroying everything..."

# Stop and remove containers
docker compose -f deploy/docker/docker-compose.yml down --volumes --remove-orphans 2>/dev/null || true

# Remove images built by this project
docker images --format '{{.Repository}}:{{.Tag}}' | grep -E '^forge/' | xargs -r docker rmi -f 2>/dev/null || true

# Remove all project volumes
docker volume ls --format '{{.Name}}' | grep -E '^forge_' | xargs -r docker volume rm 2>/dev/null || true

# Clean generated directories
rm -rf "$DEPLOY_DIR/data"/*
rm -rf "$DEPLOY_DIR/logs"/*
rm -rf "$DEPLOY_DIR/backup"/*

echo ""
echo "Everything destroyed."
echo ""
echo "To start fresh:"
echo "  ./deploy/scripts/build.sh"
echo "  ./deploy/scripts/up.sh"
