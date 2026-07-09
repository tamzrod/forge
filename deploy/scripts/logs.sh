#!/bin/bash
# logs.sh - Display runtime logs
# Usage: ./deploy/scripts/logs.sh [-f] [--grep pattern]
#   -f: Follow mode (tail -f)
#   --grep: Filter log lines

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$DEPLOY_DIR")"

FOLLOW=""
GREP_PATTERN=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -f|--follow)
            FOLLOW="-f"
            shift
            ;;
        --grep)
            GREP_PATTERN="$2"
            shift 2
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [-f] [--grep pattern]"
            exit 1
            ;;
    esac
done

cd "$PROJECT_DIR"

# Try docker logs first
if docker ps | grep -q forge-runtime; then
    if [[ -n "$GREP_PATTERN" ]]; then
        docker logs forge-runtime 2>&1 | grep "$GREP_PATTERN"
    else
        docker logs $FOLLOW forge-runtime 2>&1
    fi
else
    # Fall back to log file
    LOG_FILE="$DEPLOY_DIR/logs/runtime.log"
    if [[ -f "$LOG_FILE" ]]; then
        if [[ -n "$GREP_PATTERN" ]]; then
            grep "$GREP_PATTERN" "$LOG_FILE"
        elif [[ -n "$FOLLOW" ]]; then
            tail -f "$LOG_FILE"
        else
            cat "$LOG_FILE"
        fi
    else
        echo "No logs available. Is the application running?"
        exit 1
    fi
fi
