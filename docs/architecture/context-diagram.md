# Context Diagram

## Purpose

This diagram shows the **system boundary** of the Industrial Simulation Runtime and its relationships with external actors and systems. It answers: *"What is inside Forge and what is outside?"*

---

## System Context

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                     │
│                              EXTERNAL ENVIRONMENT                                    │
│                                                                                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   Software   │  │   Control    │  │    SCADA     │  │  Integration │          │
│  │  Developers  │  │   Engineers   │  │   Engineers  │  │    Teams     │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
│         │                 │                 │                 │                    │
│         │  Build & Test   │  Validate      │  Configure     │  Verify          │
│         │  Applications   │  Controllers   │  & Test HMI    │  Protocols       │
│         │                 │                 │                 │                    │
│         └─────────────────┼─────────────────┼─────────────────┘                    │
│                           │                 │                                     │
│                           ▼                 ▼                                     │
│  ┌───────────────────────────────────────────────────────────────────────────┐    │
│  │                                                                           │    │
│  │                         INDUSTRIAL APPLICATIONS                           │    │
│  │                                                                           │    │
│  │    ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────┐   │    │
│  │    │  Controller  │  │     HMI     │  │  Historian   │  │   REST  │   │    │
│  │    │     PLC      │  │   SCADA UI   │  │   Database   │  │   API   │   │    │
│  │    └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  └────┬────┘   │    │
│  │           │                  │                  │                 │         │    │
│  │           │  Read/Write      │  Read           │  Subscribe      │         │    │
│  │           │  Telemetry       │  Displays       │  Time-series    │         │    │
│  │           │                  │                 │                 │         │    │
│  └───────────┼──────────────────┼─────────────────┼─────────────────┼─────────┘    │
│              │                  │                 │                 │              │
│              │                  │                 │                 │              │
│              │     Modbus TCP / DNP3 / REST / MQTT                │              │
│              │                  │                 │                 │              │
└──────────────┼──────────────────┼─────────────────┼─────────────────┼──────────────┘
               │                  │                 │                 │
               ▼                  ▼                 ▼                 ▼
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                     │
│                              FORGE SYSTEM BOUNDARY                                   │
│                                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                             │   │
│  │                              MMA2                                            │   │
│  │                    (Modbus Memory Appliance)                                │   │
│  │                                                                             │   │
│  │  ┌─────────────────────────────────────────────────────────────────────┐ │   │
│  │  │                    Operational Memory                                 │ │   │
│  │  │                                                                     │ │   │
│  │  │   Holding Registers │ Input Registers │ Coils │ Discrete Inputs   │ │   │
│  │  │                                                                     │ │   │
│  │  └─────────────────────────────────────────────────────────────────────┘ │   │
│  │                              │                                             │   │
│  │                              │ Raw Ingest (TCP)                           │   │
│  │                              │                                             │   │
│  └──────────────────────────────┼─────────────────────────────────────────────┘   │
│                                 │                                                 │
│                                 ▼                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │                                                                             │   │
│  │                     SIMULATION RUNTIME                                       │   │
│  │                                                                             │   │
│  │  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐                   │   │
│  │  │  Scheduler   │  │   Plugin     │  │    Config     │                   │   │
│  │  │  + Clock     │  │   Loader     │  │               │                   │   │
│  │  └───────────────┘  └───────────────┘  └───────────────┘                   │   │
│  │                                                                             │   │
│  │  ┌─────────────────────────────────────────────────────────────────────┐   │   │
│  │  │                        MODEL REGISTRY                                 │   │   │
│  │  │                                                                     │   │   │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐   │   │   │
│  │  │  │  Grid   │  │   Sun   │  │  Wind   │  │ Weather │  │Reservoir│   │   │   │
│  │  │  │  Model  │  │  Model  │  │  Model  │  │  Model  │  │  Model  │   │   │   │
│  │  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘  └─────────┘   │   │   │
│  │  │                                                                     │   │   │
│  │  │               Physical World (Domain Models)                        │   │   │
│  │  └─────────────────────────────────────────────────────────────────────┘   │   │
│  │                                    │                                       │   │
│  │                                    │ Observe                               │   │
│  │                                    ▼                                       │   │
│  │  ┌─────────────────────────────────────────────────────────────────────┐   │   │
│  │  │                       DEVICE REGISTRY                               │   │   │
│  │  │                                                                     │   │   │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐   │   │   │
│  │  │  │ Weather │  │   PV   │  │Revenue  │  │  Grid  │  │  Pump  │   │   │   │
│  │  │  │ Station │  │Inverter│  │  Meter  │  │  Proxy  │  │         │   │   │   │
│  │  │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘   │   │   │
│  │  │       │             │             │             │             │         │   │   │
│  │  │       └─────────────┼─────────────┼─────────────┼─────────────┘         │   │   │
│  │  │                     │ Device Memory (Owned by each device)             │   │   │
│  │  │                     └────────────────────────────────────────────     │   │   │
│  │  │                                                                     │   │   │
│  │  │                 Virtual Industrial Equipment                         │   │   │
│  │  └─────────────────────────────────────────────────────────────────────┘   │   │
│  │                                                                             │   │
│  └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
                                                                                     
                                   ▲           ▲           ▲           ▲
                                   │           │           │           │
                                   │           │           │           │
              ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
              │  Commission- │  │   Training  │  │    Demo     │  │   Factory   │
              │    ing       │  │  Students   │  │   Attendees │  │ Acceptance  │
              │   Teams      │  │             │  │             │  │   Tests     │
              └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘
```

---

## External Actors

| Actor | Role | Interface | Use Case |
|-------|------|-----------|----------|
| **Software Developers** | Build applications | Build tools, REST API | Develop and test without hardware |
| **Control Engineers** | Validate logic | HMI, Controller PLC | Test controller behavior |
| **SCADA Engineers** | Configure systems | HMI, Modbus | Configure displays and alarms |
| **Integration Teams** | Verify protocols | Modbus TCP, DNP3, MQTT | Test protocol implementations |
| **Commissioning Teams** | Site preparation | All interfaces | Practice before deployment |
| **Training Teams** | Teach concepts | HMI, REST API | Safe learning environment |
| **FAT Systems** | Validate deliverables | Modbus, REST | Factory Acceptance Testing |

---

## External Systems

| System | Protocol | Direction | Purpose |
|--------|----------|-----------|---------|
| **MMA2** | Raw Ingest (TCP) | Inbound | Receives device telemetry |
| **SCADA/HMI** | Modbus TCP, DNP3, REST | Read/Write | Human operator interface |
| **PLC/Controller** | Modbus TCP, DNP3 | Read/Write | Automated control |
| **Historian** | MQTT, REST | Subscribe | Store time-series data |
| **Atlas-PPC** | Modbus TCP | Read | Power system analysis |

---

## Data Flows

### Inbound (External → Forge)

| Flow | Source | Destination | Content |
|------|--------|-------------|---------|
| **Configuration** | Operators | Runtime Config | Simulation parameters |
| **Plugin Loading** | Developers | Plugin Loader | Domain models and devices |
| **Clock Control** | Operators | Simulation Clock | Start, stop, speed |
| **Scenario Injection** | Test Framework | Scenario Engine | Time-series events |

### Outbound (Forge → External)

| Flow | Source | Destination | Content |
|------|--------|-------------|---------|
| **Telemetry** | Device Memory | MMA2 | Register values, quality |
| **Protocol Data** | MMA2 | SCADA/PLC | Modbus registers |
| **REST Response** | MMA2 | Applications | JSON-formatted data |
| **MQTT Publish** | MMA2 | Historians | Time-series values |

---

## System Boundary Rules

### Inside the Boundary (Forge)

- **Simulation Models** - Grid, Sun, Wind, Weather, Reservoir
- **Virtual Devices** - Weather Station, PV Inverter, Revenue Meter
- **Device Memory** - Owned by each device
- **Runtime** - Scheduler, Plugin Loader, Configuration
- **MMA2** - Operational memory appliance
- **Raw Ingest** - TCP transport to MMA2

### Outside the Boundary

- **Industrial Applications** - SCADA, HMI, PLC, Historian
- **External Protocols** - Modbus TCP server, DNP3 server, REST server
- **Human Operators** - Configuration, monitoring
- **Physical Systems** - Real industrial equipment

---

## Key Boundaries

```
┌─────────────────────────────────────────────────────────────────────┐
│                     PHYSICS vs EQUIPMENT                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   Simulation Models              Virtual Devices                     │
│   (Inside Forge)                (Inside Forge)                      │
│                                                                     │
│   • Grid                        • Weather Station                   │
│   • Sun                         • PV Inverter                       │
│   • Wind                        • Revenue Meter                     │
│   • Weather                     • Relay                             │
│                                                                     │
│   These are physics.            These are equipment.                │
│   They don't have addresses.    They have registers.               │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│                     DEVICE vs PROTOCOL                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   Device Memory                 Protocol Adapters                   │
│   (Inside Device)               (Inside Device)                     │
│                                                                     │
│   • Temperature                 • Modbus TCP                       │
│   • Voltage                     • DNP3                             │
│   • Status                      • REST                             │
│                                                                     │
│   Device owns memory.           Protocols expose memory.            │
│   Device owns encoding.         Protocol just serializes.          │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│                     SIMULATION vs REALITY                           │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   Forge Runtime                 External World                      │
│   (Inside Boundary)             (Outside Boundary)                  │
│                                                                     │
│   • Deterministic models       • Real hardware                     │
│   • Controlled timing          • Physical variability              │
│   • Reproducible results       • Production behavior               │
│                                                                     │
│   SCADA/PLC/HMI cannot distinguish real from virtual origins.     │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Verification Point

To verify this diagram is accurate, check:

1. **Every arrow crosses the system boundary exactly once**
2. **No arrow connects two external entities**
3. **MMA2 is inside Forge, Modbus servers are outside**
4. **Simulation Models never expose protocols**
5. **Virtual Devices never access other devices directly**
6. **Applications never access internal state**

---

## Related Documents

| Document | Purpose |
|----------|---------|
| [Overview](overview.md) | High-level architecture |
| [Runtime](runtime.md) | Runtime component details |
| [Device Model](device-model.md) | Device anatomy |
| [MMA2 Integration](mma2-integration.md) | MMA2 protocol |
| [Component Diagram](component-diagram.md) | Module relationships |

---

*Created: 2026-07-13*  
*Type: Architecture Artifact*  
*Status: Initial*
