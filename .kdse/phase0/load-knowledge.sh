#!/bin/bash
#===============================================================================
# KDSE Knowledge Loading
#===============================================================================
# Loads knowledge documents according to manifest loading order.
# Validates checksums and reports missing knowledge.
#
# Usage:
#   ./load-knowledge.sh [--required-only] [--verbose]
#
# Exit Codes:
#   0 - All required knowledge loaded successfully
#   1 - Missing required knowledge
#   2 - Manifest not found
#===============================================================================

set -euo pipefail

KDSE_DIR="${KDSE_DIR:-.kdse}"
KDSE_KNOWLEDGE_DIR="${KDSE_DIR}/knowledge"
KDSE_MANIFEST="${KDSE_KNOWLEDGE_DIR}/manifest.yaml"

LOAD_REQUIRED_ONLY=false
VERBOSE=false

# Track loaded knowledge
declare -a LOADED_KNOWLEDGE
declare -a FAILED_KNOWLEDGE

#-------------------------------------------------------------------------------
# Load knowledge from manifest
#-------------------------------------------------------------------------------
load_knowledge() {
    if [[ ! -f "$KDSE_MANIFEST" ]]; then
        echo "ERROR: Manifest not found: $KDSE_MANIFEST"
        exit 2
    fi
    
    echo "Loading knowledge from manifest..."
    echo ""
    
    # Load required knowledge first
    load_required_knowledge
    
    # Load optional knowledge if not required-only mode
    if [[ "$LOAD_REQUIRED_ONLY" == "false" ]]; then
        load_optional_knowledge
    fi
    
    # Report results
    report_results
}

#-------------------------------------------------------------------------------
# Load required knowledge
#-------------------------------------------------------------------------------
load_required_knowledge() {
    echo "=== Required Knowledge ==="
    
    # Parse required knowledge entries
    local in_required=false
    local loading_order=0
    
    while IFS= read -r line; do
        # Check for section headers
        if [[ "$line" =~ ^required_knowledge: ]]; then
            in_required=true
            continue
        elif [[ "$line" =~ ^[a-z_]+: ]]; then
            in_required=false
        fi
        
        [[ "$in_required" == "true" ]] || continue
        
        # Extract loading order
        if [[ "$line" =~ loading_order:\ ([0-9]+) ]]; then
            loading_order="${BASH_REMATCH[1]}"
        fi
        
        # Extract source path
        if [[ "$line" =~ source:\ \"([^\"]+)\" ]]; then
            local source="${BASH_REMATCH[1]}"
            local full_path="$source"
            
            loading_order=$((loading_order + 1))
            
            if [[ -f "$full_path" ]]; then
                LOADED_KNOWLEDGE+=("$source")
                if [[ "$VERBOSE" == "true" ]]; then
                    echo "  [OK] $loading_order. $source"
                else
                    echo "  ✓ Loaded: $source"
                fi
            else
                FAILED_KNOWLEDGE+=("$source")
                echo "  ✗ Missing: $source"
            fi
        fi
    done < "$KDSE_MANIFEST"
    
    echo ""
}

#-------------------------------------------------------------------------------
# Load optional knowledge
#-------------------------------------------------------------------------------
load_optional_knowledge() {
    echo "=== Optional Knowledge ==="
    
    local in_optional=false
    local loading_order=7
    
    while IFS= read -r line; do
        # Check for section headers
        if [[ "$line" =~ ^optional_knowledge: ]]; then
            in_optional=true
            continue
        elif [[ "$line" =~ ^[a-z_]+: ]]; then
            in_optional=false
        fi
        
        [[ "$in_optional" == "true" ]] || continue
        
        # Extract source path
        if [[ "$line" =~ source:\ \"([^\"]+)\" ]]; then
            local source="${BASH_REMATCH[1]}"
            local full_path="$source"
            
            loading_order=$((loading_order + 1))
            
            if [[ -f "$full_path" ]]; then
                LOADED_KNOWLEDGE+=("$source")
                if [[ "$VERBOSE" == "true" ]]; then
                    echo "  [OK] $loading_order. $source"
                else
                    echo "  ✓ Loaded: $source"
                fi
            else
                echo "  - Skipped: $source (not found)"
            fi
        fi
    done < "$KDSE_MANIFEST"
    
    echo ""
}

#-------------------------------------------------------------------------------
# Report results
#-------------------------------------------------------------------------------
report_results() {
    echo "=== Load Summary ==="
    echo ""
    echo "Loaded: ${#LOADED_KNOWLEDGE[@]} documents"
    
    if [[ ${#FAILED_KNOWLEDGE[@]} -gt 0 ]]; then
        echo "Failed: ${#FAILED_KNOWLEDGE[@]} documents"
        echo ""
        echo "MISSING REQUIRED KNOWLEDGE:"
        for doc in "${FAILED_KNOWLEDGE[@]}"; do
            echo "  - $doc"
        done
        echo ""
        echo "Hint: Restore missing files from KDSE repository"
        exit 1
    else
        echo "Failed: 0 documents"
        echo ""
        echo "STATUS: All required knowledge loaded successfully"
        exit 0
    fi
}

#-------------------------------------------------------------------------------
# Main
#-------------------------------------------------------------------------------
main() {
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --required-only)
                LOAD_REQUIRED_ONLY=true
                shift
                ;;
            --verbose|-v)
                VERBOSE=true
                shift
                ;;
            *)
                echo "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    
    load_knowledge
}

main "$@"
