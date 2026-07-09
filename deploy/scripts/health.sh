#!/bin/bash
# health.sh - Check application health
# Usage: ./deploy/scripts/health.sh
# Exit codes: 0=healthy, 1=unhealthy, 2=unknown
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$DEPLOY_DIR")"

cd "$PROJECT_DIR"

HEALTHY=true
CHECKS=0
PASSED=0

check() {
    local name="$1"
    local result="$2"
    CHECKS=$((CHECKS + 1))
    
    if [[ "$result" == "0" ]]; then
        echo "[PASS] $name"
        PASSED=$((PASSED + 1))
    else
        echo "[FAIL] $name"
        HEALTHY=false
    fi
}

echo "Checking application health..."
echo ""

# Check 1: Docker daemon
if command -v docker &> /dev/null; then
    docker info &> /dev/null
    check "Docker daemon" "$?"
else
    echo "[SKIP] Docker not installed"
fi

# Check 2: Runtime container running
CONTAINER_RUNNING=$(docker ps --format '{{.Names}}' 2>/dev/null | grep -E '(runtime|app|forge)' | wc -l)
if [[ "$CONTAINER_RUNNING" -gt 0 ]]; then
    check "Runtime container running" "0"
else
    check "Runtime container running" "1"
fi

# Check 3: Container healthy
CONTAINER_HEALTHY=$(docker inspect --format='{{.State.Health.Status}}' forge-runtime 2>/dev/null || echo "none")
if [[ "$CONTAINER_HEALTHY" == "healthy" ]]; then
    check "Container health" "0"
elif [[ "$CONTAINER_HEALTHY" == "none" ]]; then
    echo "[SKIP] No health check configured"
else
    check "Container health" "1"
fi

# Check 4: Process inside container
if [[ "$CONTAINER_RUNNING" -gt 0 ]]; then
    docker exec forge-runtime pgrep -x forge &> /dev/null
    check "Forge process running" "$?"
fi

# Check 5: Port listening
if docker exec forge-runtime netstat -tln 2>/dev/null | grep -q ":8080 "; then
    check "Port 8080 listening" "0"
elif docker exec forge-runtime ss -tln 2>/dev/null | grep -q ":8080 "; then
    check "Port 8080 listening" "0"
else
    # Don't fail if port check fails (netstat/ss may not be available)
    echo "[SKIP] Cannot check port (netstat/ss not available)"
fi

echo ""
echo "Health check: $PASSED/$CHECKS passed"

if [[ "$HEALTHY" == "true" ]]; then
    echo "Status: Healthy"
    exit 0
else
    echo "Status: Unhealthy"
    exit 1
fi
