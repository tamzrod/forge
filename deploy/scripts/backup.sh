#!/bin/bash
# backup.sh - Create backup of persistent data
# Usage: ./deploy/scripts/backup.sh [name]
#   Without name: creates timestamped backup
#   With name: creates named backup

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$DEPLOY_DIR")"

cd "$PROJECT_DIR"

# Create backup directory
BACKUP_DIR="$DEPLOY_DIR/backup"
mkdir -p "$BACKUP_DIR"

# Generate backup name
if [[ -n "$1" ]]; then
    BACKUP_NAME="$1"
else
    BACKUP_NAME="backup-$(date +%Y%m%d-%H%M%S)"
fi

BACKUP_FILE="$BACKUP_DIR/${BACKUP_NAME}.tar.gz"

echo "Creating backup: $BACKUP_NAME"

# Stop services to ensure consistent backup
"$SCRIPT_DIR/down.sh" 2>/dev/null || true

# Create backup
tar -czf "$BACKUP_FILE" \
    -C "$DEPLOY_DIR" \
    config/ \
    data/ \
    2>/dev/null || true

# Restart services
"$SCRIPT_DIR/up.sh" 2>/dev/null || true

if [[ -f "$BACKUP_FILE" ]]; then
    SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    echo "Backup created: $BACKUP_FILE ($SIZE)"
else
    echo "Error: Backup failed"
    exit 1
fi

# List available backups
echo ""
echo "Available backups:"
ls -lh "$BACKUP_DIR"/*.tar.gz 2>/dev/null | awk '{print $9, $5}'
