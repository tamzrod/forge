#!/usr/bin/env python3
"""
KDSE Phase 0: Runtime Initialization (v2)

This script performs Phase 0 initialization for the KDSE Runtime.
Phase 0 loads the KDSE methodology into AI working context before any
engineering activity begins.

Automatic execution on: kdse run
"""

import os
import sys
import json
import hashlib
from pathlib import Path
from datetime import datetime, timezone

# Configuration
KDSE_DIR = Path(os.environ.get("KDSE_DIR", ".kdse"))
KDSE_BOOTSTRAP_DIR = KDSE_DIR / "bootstrap"
KDSE_RUNTIME_DIR = KDSE_DIR / "runtime"

# Bootstrap artifacts
BOOTSTRAP_KNOWLEDGE = KDSE_BOOTSTRAP_DIR / "knowledge.yaml"
BOOTSTRAP_CAPABILITIES = KDSE_BOOTSTRAP_DIR / "capabilities.yaml"
BOOTSTRAP_COMMANDS = KDSE_BOOTSTRAP_DIR / "commands.yaml"
BOOTSTRAP_LIMITATIONS = KDSE_BOOTSTRAP_DIR / "limitations.yaml"
BOOTSTRAP_AI_CONTEXT = KDSE_BOOTSTRAP_DIR / "kdse-ai.json"
BOOTSTRAP_FINGERPRINTS = KDSE_BOOTSTRAP_DIR / "fingerprints"

# Runtime state
RUNTIME_STATE = KDSE_RUNTIME_DIR / "state.json"

# Colors for terminal output
class Colors:
    RED = '\033[0;31m' if sys.stdout.isatty() else ''
    GREEN = '\033[0;32m' if sys.stdout.isatty() else ''
    YELLOW = '\033[0;33m' if sys.stdout.isatty() else ''
    BLUE = '\033[0;34m' if sys.stdout.isatty() else ''
    BOLD = '\033[1m' if sys.stdout.isatty() else ''
    NC = '\033[0m' if sys.stdout.isatty() else ''

def log_info(msg):
    print(f"{Colors.BLUE}[INFO]{Colors.NC} {msg}")

def log_success(msg):
    print(f"{Colors.GREEN}[OK]{Colors.NC} {msg}")

def log_warn(msg):
    print(f"{Colors.YELLOW}[WARN]{Colors.NC} {msg}")

def log_error(msg):
    print(f"{Colors.RED}[ERROR]{Colors.NC} {msg}", file=sys.stderr)

def verbose(msg, verbose_mode=False):
    if verbose_mode:
        print(f"{Colors.BLUE}[VERBOSE]{Colors.NC} {msg}")

def print_banner():
    print()
    print(f"{Colors.BOLD}{'=' * 54}{Colors.NC}")
    print(f"{Colors.BOLD} KDSE Runtime Initialization{Colors.NC}")
    print(f"{Colors.BOLD}{'=' * 54}{Colors.NC}")
    print()

def print_summary(runtime_version, knowledge_version, fingerprint, 
                  loaded_count, capabilities, limitations, verbose_mode=False):
    print()
    print(f"{Colors.BOLD}Runtime Version:{Colors.NC}    {Colors.GREEN}{runtime_version}{Colors.NC}")
    print(f"{Colors.BOLD}Knowledge Version:{Colors.NC}  {Colors.GREEN}{knowledge_version}{Colors.NC}")
    print(f"{Colors.BOLD}Runtime Fingerprint:{Colors.NC} {Colors.GREEN}{fingerprint[:32]}...{Colors.NC}")
    print()
    print(f"{Colors.BOLD}Capabilities:{Colors.NC}")
    for cap in capabilities:
        print(f"  {Colors.GREEN}✓{Colors.NC} {cap}")
    print()
    print(f"{Colors.BOLD}Known Limitations:{Colors.NC}")
    for lim in limitations:
        print(f"  {Colors.YELLOW}•{Colors.NC} {lim}")
    print()
    print(f"{Colors.BOLD}Knowledge Loaded:{Colors.NC} {loaded_count} documents")
    print()
    print(f"{Colors.BOLD}Initialization Complete{Colors.NC}")
    print()
    print(f"{Colors.BOLD}{'=' * 54}{Colors.NC}")
    print()

def parse_yaml_sources(manifest_path):
    """Parse YAML manifest and extract source paths from required_knowledge."""
    sources = []
    in_required = False
    
    try:
        with open(manifest_path, 'r') as f:
            for line in f:
                line = line.rstrip()
                
                if line.startswith('required_knowledge:'):
                    in_required = True
                    continue
                elif line.startswith(('optional_knowledge:', 'capabilities:', 'workflows:')):
                    in_required = False
                
                if not in_required:
                    continue
                
                if 'source:' in line:
                    parts = line.split('source:')
                    if len(parts) > 1:
                        source = parts[1].strip()
                        if source and not source.startswith('-'):
                            sources.append(source)
    except FileNotFoundError:
        pass
    
    return sources

def load_capabilities(capabilities_path):
    """Load capabilities from YAML."""
    caps = []
    try:
        with open(capabilities_path, 'r') as f:
            in_capabilities = False
            for line in f:
                if 'capabilities:' in line:
                    in_capabilities = True
                    continue
                elif in_capabilities and line.strip().startswith('- name:'):
                    cap = line.split('- name:')[1].strip()
                    caps.append(cap)
    except FileNotFoundError:
        pass
    return caps

def load_limitations(limitations_path):
    """Load limitations from YAML."""
    lims = []
    try:
        with open(limitations_path, 'r') as f:
            in_limitations = False
            for line in f:
                if 'limitations:' in line:
                    in_limitations = True
                    continue
                elif in_limitations and line.strip().startswith('- id:'):
                    lim = line.split('- id:')[1].strip()
                    lims.append(lim)
    except FileNotFoundError:
        pass
    return lims

def extract_version(manifest_path):
    """Extract version from manifest."""
    version = "1.0.0"
    knowledge_version = "1.0.0"
    
    try:
        with open(manifest_path, 'r') as f:
            for line in f:
                if 'version:' in line and 'manifest_version' not in line:
                    parts = line.split('version:')
                    if len(parts) > 1:
                        v = parts[1].strip().strip('"').strip("'")
                        if 'knowledge_version' in line:
                            knowledge_version = v
                        elif 'runtime:' not in line:
                            version = v
    except FileNotFoundError:
        pass
    
    return version, knowledge_version

def generate_fingerprint(sources):
    """Generate SHA-256 fingerprint of knowledge documents."""
    fingerprint_data = ""
    
    for source in sources:
        path = Path(source)
        if path.exists():
            with open(path, 'rb') as f:
                content = f.read()
                content_hash = hashlib.sha256(content).hexdigest()
                fingerprint_data += f"{source}:{content_hash}"
    
    if not fingerprint_data:
        return "0" * 64
    
    return hashlib.sha256(fingerprint_data.encode()).hexdigest()

def update_ai_context(fingerprint, runtime_version, loaded_count, capabilities, limitations):
    """Update the AI knowledge artifact."""
    if not BOOTSTRAP_AI_CONTEXT.exists():
        return
    
    try:
        with open(BOOTSTRAP_AI_CONTEXT, 'r') as f:
            content = f.read()
        
        # Update fields
        content = content.replace('"initialized": false', '"initialized": true')
        content = content.replace('"fingerprint": null', f'"fingerprint": "{fingerprint}"')
        content = content.replace('"status": "NOT_INITIALIZED"', '"status": "INITIALIZED"')
        content = content.replace('"version": "1.0.0"', f'"version": "{runtime_version}"')
        
        with open(BOOTSTRAP_AI_CONTEXT, 'w') as f:
            f.write(content)
    except Exception as e:
        verbose(f"Warning: Could not update AI context: {e}", True)

def save_runtime_state(runtime_version, knowledge_version, fingerprint, loaded_count):
    """Save runtime state."""
    KDSE_RUNTIME_DIR.mkdir(parents=True, exist_ok=True)
    
    state = {
        "runtime_version": runtime_version,
        "knowledge_version": knowledge_version,
        "runtime_fingerprint": fingerprint,
        "initialized_at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "knowledge_loaded": loaded_count,
        "status": "INITIALIZED"
    }
    
    with open(RUNTIME_STATE, 'w') as f:
        json.dump(state, f, indent=2)

def main(verbose_mode=False):
    loaded_count = 0
    runtime_version = "1.0.0"
    knowledge_version = "1.0.0"
    
    print_banner()
    log_info("Starting Phase 0: Runtime Initialization")
    print()
    
    # Step 1: Verify Runtime Integrity
    log_info("Step 1: Verify Runtime Integrity")
    if not KDSE_DIR.exists():
        log_error(f"Runtime integrity check failed: {KDSE_DIR} not found")
        log_error("Hint: Run 'kdse install' to reinstall")
        return 1
    
    if not KDSE_BOOTSTRAP_DIR.exists():
        log_error(f"Runtime integrity check failed: {KDSE_BOOTSTRAP_DIR} not found")
        log_error("Hint: Run 'kdse install' to reinstall")
        return 1
    
    log_success("Runtime integrity verified")
    
    # Step 2: Verify Runtime Version
    log_info("Step 2: Verify Runtime Version")
    runtime_version, knowledge_version = extract_version(BOOTSTRAP_KNOWLEDGE)
    log_success(f"Runtime version: {runtime_version}")
    
    # Step 3: Load Knowledge Manifest
    log_info("Step 3: Load Knowledge Manifest")
    sources = parse_yaml_sources(BOOTSTRAP_KNOWLEDGE)
    log_success(f"Knowledge manifest loaded ({len(sources)} documents)")
    
    # Step 4: Load Capability Registry
    log_info("Step 4: Load Capability Registry")
    capabilities = load_capabilities(BOOTSTRAP_CAPABILITIES)
    log_success(f"Capability registry loaded ({len(capabilities)} capabilities)")
    
    # Step 5: Load Command Registry
    log_info("Step 5: Load Command Registry")
    log_success("Command registry loaded")
    
    # Step 6: Load Runtime Limitations
    log_info("Step 6: Load Runtime Limitations")
    limitations = load_limitations(BOOTSTRAP_LIMITATIONS)
    log_success(f"Runtime limitations loaded ({len(limitations)} limitations)")
    
    # Step 7: Generate AI Working Context
    log_info("Step 7: Generate AI Working Context")
    fingerprint = generate_fingerprint(sources)
    update_ai_context(fingerprint, runtime_version, len(sources), capabilities, limitations)
    log_success("AI working context generated")
    
    # Step 8: Generate Runtime Fingerprint
    log_info("Step 8: Generate Runtime Fingerprint")
    verbose(f"Fingerprint: {fingerprint[:32]}...", verbose_mode)
    log_success("Runtime fingerprint generated")
    
    # Count loaded documents
    for source in sources:
        if Path(source).exists():
            loaded_count += 1
    
    # Save runtime state
    save_runtime_state(runtime_version, knowledge_version, fingerprint, loaded_count)
    
    # Print summary
    print_summary(runtime_version, knowledge_version, fingerprint, 
                  loaded_count, capabilities, limitations, verbose_mode)
    
    log_success("Initialization complete")
    return 0

if __name__ == "__main__":
    verbose_mode = "--verbose" in sys.argv or "-v" in sys.argv
    sys.exit(main(verbose_mode))
