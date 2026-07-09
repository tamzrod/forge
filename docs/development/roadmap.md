# Roadmap

## Overview

This roadmap outlines the development trajectory for the Industrial Simulation Runtime.

## Version 1.0 - Foundation

**Goal:** Core runtime with basic simulation capabilities.

### Milestone 1.1 - Runtime Core

- [ ] Runtime initialization and shutdown
- [ ] Basic device lifecycle (create, destroy)
- [ ] Memory management (allocate, read, write)
- [ ] Scheduler with fixed tick interval
- [ ] Configuration from YAML

### Milestone 1.2 - Behaviors

- [ ] Behavior interface
- [ ] Behavior registration
- [ ] Behavior execution on tick
- [ ] Memory access from behaviors

### Milestone 1.3 - Protocols

- [ ] Protocol adapter interface
- [ ] Modbus TCP adapter
- [ ] Protocol binding to devices
- [ ] Memory exposure via protocol

### Milestone 1.4 - Basic Plugin System

- [ ] Plugin interface
- [ ] Plugin loading
- [ ] Device type registration
- [ ] Behavior factory registration

## Version 1.1 - Fault System

**Goal:** Fault injection capabilities.

### Milestone 1.1.1 - Fault Framework

- [ ] Fault interface
- [ ] Fault attachment to devices
- [ ] Memory access with fault wrapping
- [ ] Quality flag support

### Milestone 1.1.2 - Basic Faults

- [ ] Communication loss fault
- [ ] Frozen values fault
- [ ] Bad quality fault
- [ ] Offline fault

### Milestone 1.1.3 - Advanced Faults

- [ ] Noise fault
- [ ] Scaling error fault
- [ ] Delay fault
- [ ] Drift fault

## Version 1.2 - Scenario Engine

**Goal:** Event injection and scenario playback.

### Milestone 1.2.1 - Scenario Framework

- [ ] Scenario definition
- [ ] Event types
- [ ] Scenario loading
- [ ] Event execution on tick

### Milestone 1.2.2 - Event Actions

- [ ] Setpoint change action
- [ ] Fault injection action
- [ ] Device state change action
- [ ] Environmental change action

### Milestone 1.2.3 - Scenario Management

- [ ] Scenario composition
- [ ] Conditional events
- [ ] Scenario validation
- [ ] Scenario templates

## Version 1.3 - Energy Plugin

**Goal:** First complete domain plugin.

### Milestone 1.3.1 - Basic Energy Devices

- [ ] Weather station device
- [ ] PV model device
- [ ] Inverter device
- [ ] Revenue meter device

### Milestone 1.3.2 - Energy Behaviors

- [ ] Weather behavior
- [ ] PV power calculation behavior
- [ ] Inverter conversion behavior
- [ ] Power measurement behavior

### Milestone 1.3.3 - Grid Support

- [ ] Grid device
- [ ] Transformer device
- [ ] Relay device
- [ ] Grid connection behavior

## Version 1.4 - Protocol Expansion

**Goal:** Additional protocol support.

### Milestone 1.4.1 - DNP3 Adapter

- [ ] DNP3 protocol implementation
- [ ] Binary and analog points mapping
- [ ] Event class support
- [ ] Time synchronization

### Milestone 1.4.2 - REST API

- [ ] REST adapter
- [ ] Device management endpoints
- [ ] Memory read/write endpoints
- [ ] Scenario control endpoints

### Milestone 1.4.3 - MQTT Adapter

- [ ] MQTT protocol implementation
- [ ] Topic structure
- [ ] Quality of service
- [ ] Last will and testament

## Version 2.0 - Multi-Domain

**Goal:** Support for multiple industrial domains.

### Milestone 2.0.1 - Water Plugin

- [ ] Pump device
- [ ] Valve device
- [ ] Tank device
- [ ] Flow meter device

### Milestone 2.0.2 - Manufacturing Plugin

- [ ] PLC device
- [ ] Robot device
- [ ] Conveyor device
- [ ] Sensor device

### Milestone 2.0.3 - Building Automation Plugin

- [ ] HVAC controller
- [ ] Lighting controller
- [ ] Access control
- [ ] Fire alarm

## Version 2.1 - Advanced Features

**Goal:** Enhanced simulation capabilities.

### Milestone 2.1.1 - Parallel Execution

- [ ] Device dependency graph
- [ ] Parallel behavior execution
- [ ] Thread-safe memory access
- [ ] Performance metrics

### Milestone 2.1.2 - Time Control

- [ ] Time multiplier
- [ ] Time jumping
- [ ] Time reversal (for debugging)
- [ ] Playback from log

### Milestone 2.1.3 - Data Recording

- [ ] Memory value logging
- [ ] Event logging
- [ ] Playback from recording
- [ ] Export to CSV/Parquet

## Version 2.2 - Operational

**Goal:** Production-readiness features.

### Milestone 2.2.1 - Monitoring

- [ ] Metrics collection
- [ ] Health checks
- [ ] Alerting hooks
- [ ] Dashboard integration

### Milestone 2.2.2 - Persistence

- [ ] Simulation state save/load
- [ ] Device state export
- [ ] Configuration versioning
- [ ] Rollback support

### Milestone 2.2.3 - Security

- [ ] Protocol authentication
- [ ] TLS support
- [ ] Role-based access
- [ ] Audit logging

## Future Considerations

These items are under consideration but not yet scheduled:

### Advanced Simulation

- Co-simulation with external simulators
- Hardware-in-the-loop support
- Real-time synchronization
- Distributed simulation

### Advanced Protocols

- IEC 61850
- OPC UA
- EtherNet/IP
- PROFINET

### Advanced Faults

- Partial communication loss
- Intermittent faults
- Cascading failures
- Cyber attack simulation

### Machine Learning

- Anomaly detection
- Predictive maintenance
- Behavioral learning
- Digital twin alignment

### Visualization

- 3D visualization
- Real-time dashboards
- Timeline explorer
- State machine visualization

## Deprecation Policy

When deprecating features:

1. Announce deprecation in release notes
2. Mark with `// Deprecated:` comment
3. Keep for at least one major version
4. Remove in subsequent major version

## Versioning

We follow semantic versioning:

- **Major**: Breaking changes
- **Minor**: New features, backward compatible
- **Patch**: Bug fixes, backward compatible

## Release Cadence

- **Major releases**: Annually
- **Minor releases**: Quarterly
- **Patches**: As needed

## Contributing to Roadmap

To propose roadmap items:

1. Open a GitHub Discussion
2. Describe the use case
3. Explain the value
4. Estimate complexity
5. Suggest priority
