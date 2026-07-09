# Protocol Architecture

## Philosophy

**Protocols are NOT part of the simulation runtime.**

The simulation runtime publishes data to MMA2. MMA2 exposes protocols.

## Two Systems

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Simulation Runtime                                     │
│                                                                         │
│   ┌─────────────────────────────────────────────────────────────────┐ │
│   │                    Device Memory                                    │ │
│   │  (private, internal)                                              │ │
│   └─────────────────────────────────────────────────────────────────┘ │
│                              │                                           │
│                              ▼                                           │
│   ┌─────────────────────────────────────────────────────────────────┐ │
│   │                  Raw Ingest Publisher                              │ │
│   │  (publishes to MMA2)                                             │ │
│   └─────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                              MMA2                                         │
│                                                                         │
│   ┌─────────────────────────────────────────────────────────────────┐ │
│   │                   Operational Memory                                │ │
│   │  (shared, visible to all)                                        │ │
│   └─────────────────────────────────────────────────────────────────┘ │
│                              │                                           │
│                              ▼                                           │
│   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │
│   │  Modbus TCP   │  │     DNP3     │  │   REST API   │            │
│   └──────────────┘  └──────────────┘  └──────────────┘            │
│                                                                         │
│   Protocols are in MMA2, not in the simulation runtime                │
└─────────────────────────────────────────────────────────────────────────┘
```

## The Runtime Does Not Expose Protocols

The simulation runtime:
- Maintains device memory
- Publishes via Raw Ingest
- **Does NOT expose Modbus, DNP3, or other protocols**

MMA2:
- Owns operational memory
- Exposes protocols (Modbus, DNP3, REST, MQTT)
- Is the integration point for external systems

## Raw Ingest

Raw Ingest is the official interface between the simulation runtime and MMA2.

```go
type RawIngestPublisher interface {
    Publish(tag string, value interface{}, quality Quality) error
}
```

Devices publish operational data:

```go
behavior.Tick() {
    // Compute from device memory
    irradiance := device.Memory().ReadFloat32("input_registers", irradianceAddr)
    
    // Publish to MMA2
    publisher.Publish("weather/irradiance", irradiance, QualityGood)
}
```

## Why This Separation

```
Real Device                          Virtual Device
     │                                   │
     ▼                                   ▼
Replicator                        Raw Ingest
     │                                   │
     └───────────────┬───────────────────┘
                     │
                     ▼
              ┌─────────────┐
              │    MMA2     │
              │  (shared    │
              │   state)    │
              └─────────────┘
                     │
     ┌───────────────┼───────────────┐
     ▼               ▼               ▼
Modbus TCP        DNP3            REST
```

Atlas-PPC, SCADA, HMIs, and Historians read from MMA2. They cannot distinguish between real and virtual device origins.

## Key Principle

**The simulation runtime publishes data. MMA2 exposes protocols.**

The runtime does not implement Modbus servers. The runtime does not implement DNP3 masters. The runtime only publishes to MMA2 via Raw Ingest.

## Benefits

1. **Clean separation** - Simulation is decoupled from integration
2. **Real/virtual equivalence** - Both publish via Raw Ingest
3. **MMA2 owns protocols** - Single point for protocol configuration
4. **Scalability** - Multiple simulation runtimes can publish to one MMA2
5. **Consistency** - Real and virtual devices appear identical
