#!/bin/bash
# restore.sh - Restore from backup
# Usage: ./deploy/scripts/restore.sh [backup_file]
#   Without argument: interactive selection
#   With argument: restore specified backup

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$DEPLOY_DIR")"

cd "$PROJECT_DIR"

BACKUP_DIR="$DEPLOY_DIR/backup"

# List available backups
list_backups() {
    echo "Available backups:"
    local i=1
    for f in "$BACKUP_DIR"/*.tar.gz; do
        if [[ -f "$f" ]]; then
            local name=$(basename "$f")
            local size=$(du -h "$f" | cut -f1)
            echo "  [$i] $name ($size)"
            BACKUPS[$i]="$f"
            i=$((i + 1))
        fi
    done
    
    if [[ ${#BACKUPS[@]} -eq 0 ]]; then
        echo "  No backups found in $BACKUP_DIR"
        exit 1
    fi
}

declare -A BACKUPS

# Select backup
if [[ -n "$1" ]]; then
    # Use provided argument
    BACKUP_FILE="$1"
    if [[ ! -f "$BACKUP_FILE" ]]; then
        echo "Error: Backup file not found: $BACKUP_FILE"
        exit 1
    fi
else
    # Interactive selection
    list_backups
    echo ""
    read -p "Select backup [1-${#BACKUPS[@]}]: " selection
    
    if [[ -z "${BACKUPS[$selection]}" ]]; then
        echo "Invalid selection"
        exit 1
    fi
    
    BACKUP_FILE="${BACKUPS[$selection]}"
fi

BACKUP_NAME=$(basename "$BACKUP_FILE")

echo ""
echo "WARNING: This will replace current data with backup: $BACKUP_NAME"
echo ""
read -p "Continue? [y/N] " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 0
fi

# Stop services
echo "Stopping services..."
"$SCRIPT_DIR/down.sh" 2>/dev/null || true

# Remove current data
rm -rf "$DEPLOY_DIR/data"/*
rm -rf "$DEPLOY_DIR/config"/*

# Extract backup
echo "Restoring from backup..."
tar -xzf "$BACKUP_FILE" -C "$DEPLOY_DIR"

# Restart services
"$SCRIPT_DIR/up.sh" 2>/dev/null || true

echo ""
echo "Restore complete."
echo "Run './deploy/scripts/health.sh' to verify."
