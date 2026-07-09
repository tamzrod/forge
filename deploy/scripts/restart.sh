#!/bin/bash
# restart.sh - Restart the application
# Usage: ./deploy/scripts/restart.sh
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Restarting application..."

"$SCRIPT_DIR/down.sh"
sleep 1
"$SCRIPT_DIR/up.sh"

echo "Application restarted."
