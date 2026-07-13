#!/bin/bash
#===============================================================================
# KDSE Runtime State Management
#===============================================================================
# Manages the runtime state for KDSE Phase 0 initialization.
# Provides functions to read, write, and validate runtime state.
#
# Usage:
#   source runtime-state.sh
#   runtime_state_get <key>
#   runtime_state_set <key> <value>
#   runtime_state_save
#   runtime_state_load
#===============================================================================

KDSE_RUNTIME_DIR="${KDSE_DIR:-.kdse}/runtime"
KDSE_RUNTIME_STATE="${KDSE_RUNTIME_DIR}/state.json"

# In-memory state
declare -A RUNTIME_STATE

#-------------------------------------------------------------------------------
# Initialize runtime state
#-------------------------------------------------------------------------------
runtime_state_init() {
    mkdir -p "$KDSE_RUNTIME_DIR"
    
    if [[ -f "$KDSE_RUNTIME_STATE" ]]; then
        runtime_state_load
    else
        # Initialize with defaults
        RUNTIME_STATE["runtime_version"]="1.0.0"
        RUNTIME_STATE["knowledge_version"]="1.0.0"
        RUNTIME_STATE["knowledge_fingerprint"]=""
        RUNTIME_STATE["compatible_standard"]=">= 1.0.0"
        RUNTIME_STATE["initialized_at"]=""
        RUNTIME_STATE["repository_path"]="$(pwd)"
        RUNTIME_STATE["status"]="NOT_INITIALIZED"
    fi
}

#-------------------------------------------------------------------------------
# Load state from file
#-------------------------------------------------------------------------------
runtime_state_load() {
    if [[ ! -f "$KDSE_RUNTIME_STATE" ]]; then
        return 1
    fi
    
    # Parse JSON (simple approach - for production, use jq)
    while IFS='=' read -r key value; do
        # Skip JSON syntax
        [[ "$key" =~ ^[a-z_]+[[:space:]]*: ]] || continue
        
        key=$(echo "$key" | sed 's/[[:space:]]*://' | tr -d '"')
        value=$(echo "$value" | tr -d '",' | tr -d ' ')
        
        RUNTIME_STATE["$key"]="$value"
    done < "$KDSE_RUNTIME_STATE"
    
    return 0
}

#-------------------------------------------------------------------------------
# Save state to file
#-------------------------------------------------------------------------------
runtime_state_save() {
    mkdir -p "$KDSE_RUNTIME_DIR"
    
    cat > "$KDSE_RUNTIME_STATE" << EOF
{
  "runtime_version": "${RUNTIME_STATE[runtime_version]}",
  "knowledge_version": "${RUNTIME_STATE[knowledge_version]}",
  "knowledge_fingerprint": "${RUNTIME_STATE[knowledge_fingerprint]}",
  "compatible_standard": "${RUNTIME_STATE[compatible_standard]}",
  "initialized_at": "${RUNTIME_STATE[initialized_at]}",
  "repository_path": "${RUNTIME_STATE[repository_path]}",
  "status": "${RUNTIME_STATE[status]}"
}
EOF
    
    return 0
}

#-------------------------------------------------------------------------------
# Get a state value
#-------------------------------------------------------------------------------
runtime_state_get() {
    local key="$1"
    echo "${RUNTIME_STATE[$key]:-}"
}

#-------------------------------------------------------------------------------
# Set a state value
#-------------------------------------------------------------------------------
runtime_state_set() {
    local key="$1"
    local value="$2"
    RUNTIME_STATE["$key"]="$value"
}

#-------------------------------------------------------------------------------
# Check if initialized
#-------------------------------------------------------------------------------
runtime_state_is_initialized() {
    [[ "${RUNTIME_STATE[status]}" == "READY" ]]
}

#-------------------------------------------------------------------------------
# Get status
#-------------------------------------------------------------------------------
runtime_state_get_status() {
    echo "${RUNTIME_STATE[status]:-UNKNOWN}"
}

#-------------------------------------------------------------------------------
# Initialize on source
#-------------------------------------------------------------------------------
runtime_state_init
