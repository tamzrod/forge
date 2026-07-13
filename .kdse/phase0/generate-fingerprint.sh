#!/bin/bash
#===============================================================================
# KDSE Knowledge Fingerprint Generation
#===============================================================================
# Generates a SHA-256 fingerprint of loaded knowledge documents.
# The fingerprint is used to verify knowledge integrity.
#
# Usage:
#   ./generate-fingerprint.sh [--manifest <path>]
#
# Output:
#   SHA-256 fingerprint hash
#===============================================================================

set -euo pipefail

KDSE_DIR="${KDSE_DIR:-.kdse}"
KDSE_KNOWLEDGE_DIR="${KDSE_DIR}/knowledge"
KDSE_MANIFEST="${KDSE_KNOWLEDGE_DIR}/manifest.yaml"

#-------------------------------------------------------------------------------
# Generate fingerprint
#-------------------------------------------------------------------------------
generate_fingerprint() {
    local manifest_path="${1:-$KDSE_MANIFEST}"
    
    if [[ ! -f "$manifest_path" ]]; then
        echo "ERROR: Manifest not found: $manifest_path" >&2
        exit 1
    fi
    
    local fingerprint_data=""
    
    # Parse required knowledge sources from manifest
    local in_required=false
    
    while IFS= read -r line; do
        # Check for section headers
        if [[ "$line" =~ ^required_knowledge: ]]; then
            in_required=true
            continue
        elif [[ "$line" =~ ^[a-z_]+: ]]; then
            in_required=false
        fi
        
        [[ "$in_required" == "true" ]] || continue
        
        # Extract source path
        if [[ "$line" =~ source:\ \"([^\"]+)\" ]]; then
            local source="${BASH_REMATCH[1]}"
            local full_path="$source"
            
            if [[ -f "$full_path" ]]; then
                # Hash the content and append source:hash
                local content_hash=$(sha256sum "$full_path" 2>/dev/null | awk '{print $1}')
                fingerprint_data="${fingerprint_data}${source}:${content_hash}"
            fi
        fi
    done < "$manifest_path"
    
    # Generate final fingerprint from accumulated data
    if [[ -z "$fingerprint_data" ]]; then
        echo "WARNING: No content to fingerprint" >&2
        echo "0000000000000000000000000000000000000000000000000000000000000000"
        exit 1
    fi
    
    echo "$fingerprint_data" | sha256sum | awk '{print $1}'
}

#-------------------------------------------------------------------------------
# Main
#-------------------------------------------------------------------------------
main() {
    local manifest_path=""
    
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --manifest|-m)
                manifest_path="$2"
                shift 2
                ;;
            *)
                echo "Usage: $0 [--manifest <path>]"
                exit 1
                ;;
        esac
    done
    
    generate_fingerprint "$manifest_path"
}

main "$@"
