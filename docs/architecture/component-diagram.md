# Component Diagram

## Purpose

This diagram shows the **module structure** of the Industrial Simulation Runtime and the relationships between components. It answers: *"What modules exist and how do they depend on each other?"*

---

## Package Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                                                                 │
│                              FORGE (Root Package)                               │
│                         github.com/tamzrod/forge                                 │
│                                                                                 │
│  Provides the public API by re-exporting types from sub-packages:                │
│                                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  Types: Device, DeviceID, Behavior, MemoryImage, Quality, Config,           │   │
│  │        Scheduler, SimulationClock                                       │   │
│  │                                                                             │   │
│  │  Functions: NewRuntime, NewRuntimeFromFile, LoadConfig, DefaultConfig,      │   │
│  │             NewDevice, NewMemory                                         │   │
│  │                                                                             │   │
│  │  Constants: QualityGood, QualityUncertain, QualityBad, QualityOffline       │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      │ imports
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                                                                 │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │     RUNTIME      │  │    SCHEDULER     │  │     MODELS      │              │
│  │   runtime.go     │  │  scheduler.go    │  │   models.go    │              │
│  │                 │  │                 │  │                 │              │
│  │  - Runtime      │  │  - Scheduler    │  │  - Model       │              │
│  │  - Config       │  │  - Simulation   │  │  - GridModel   │              │
│  │  - LoadConfig   │  │    Clock       │  │  - SunModel    │              │
│  │                 │  │                 │  │  - WindModel   │              │
│  │                 │  │                 │  │  - WeatherModel│              │
│  │                 │  │                 │  │  - Reservoir   │              │
│  │                 │  │                 │  │    Model      │              │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘              │
│           │                    │                    │                          │
│           └────────────────────┼────────────────────┘                          │
│                                │                                               │
│                                ▼                                               │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                              DEVICE                                      │   │
│  │                           device.go                                     │   │
│  │                                                                       │   │
│  │  ┌─────────────────────────────────────────────────────────────────┐ │   │
│  │  │  Device                                                           │ │   │
│  │  │                                                                   │ │   │
│  │  │  - Owns: MemoryImage, Behaviors                                   │ │   │
│  │  │  - References: Simulation Models (via ModelProvider)                │ │   │
│  │  │                                                                   │ │   │
│  │  │  ┌───────────────────┐  ┌───────────────────┐                     │ │   │
│  │  │  │     Memory      │  │    Behaviors      │                     │ │   │
│  │  │  │   (imports)     │  │   (interface)     │                     │ │   │
│  │  │  └───────────────────┘  └───────────────────┘                     │ │   │
│  │  └─────────────────────────────────────────────────────────────────┘ │   │
│  │                                                                       │   │
│  │  Sub-packages:                                                        │   │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐       │   │
│  │  │   Registry      │  │   Weather-      │  │   (Future)      │       │   │
│  │  │   registry.go   │  │   Station      │  │                 │       │   │
│  │  │                 │  │  weatherstation│  │   - PV Inverter │       │   │
│  │  │  - DeviceReg    │  │   /station.go  │  │   - RevenueMeter│       │   │
│  │  │  - Get/Set     │  │                │  │   - Relay       │       │   │
│  │  │  - List/Delete │  │  - WeatherSt   │  │   - GridProxy   │       │   │
│  │  └─────────────────┘  └─────────────────┘  └─────────────────┘       │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      │ imports
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                                                                 │
│                              MEMORY (Package)                                    │
│                           github.com/tamzrod/forge/memory                        │
│                                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  MemoryImage                                                           │   │
│  │                                                                         │   │
│  │  - regions: map[string]*Region                                        │   │
│  │  - RWMutex for thread safety                                          │   │
│  │                                                                         │   │
│  │  ┌─────────────────────────────────────────────────────────────────┐ │   │
│  │  │  Region                                                           │ │   │
│  │  │                                                                   │ │   │
│  │  │  - Name: string                                                   │ │   │
│  │  │  - Size: uint32                                                   │ │   │
│  │  │  - Values: []byte                                                │ │   │
│  │  │  - Quality: []Quality                                             │ │   │
│  │  └─────────────────────────────────────────────────────────────────┘ │   │
│  │                                                                         │   │
│  │  Quality flags: Good (0x00), Uncertain (0x40), Bad (0x80), Offline   │   │
│  │                 (0x84)                                                 │   │
│  │                                                                         │   │
│  │  Read/Write functions: Read, Write, ReadUint16, WriteUint16,           │   │
│  │                       ReadFloat32, WriteFloat32                        │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Module Dependencies

```
                           ┌──────────────┐
                           │    FORGE     │
                           │  (root pkg)  │
                           └──────┬───────┘
                                  │
          ┌───────────────────────┼───────────────────────┐
          │                       │                       │
          ▼                       ▼                       ▼
┌──────────────────┐   ┌──────────────────┐   ┌──────────────────┐
│     RUNTIME      │   │    SCHEDULER     │   │     MODELS       │
│                  │   │                  │   │                  │
│ imports:          │   │ imports:         │   │ imports:         │
│  - device        │   │  - device        │   │  - (math, time) │
│  - models        │   │  - models        │   │                  │
│  - scheduler      │   │                  │   │ exports:         │
│  - memory        │   │ exports:         │   │  - Model interface│
│                  │   │  - Scheduler     │   │  - GridModel     │
│ exports:         │   │  - SimulationClock│   │  - SunModel     │
│  - Runtime       │   │                  │   │  - WindModel    │
│  - Config        │   │                  │   │  - WeatherModel  │
│                  │   │                  │   │  - ReservoirModel│
└────────┬─────────┘   └────────┬─────────┘   └──────────────────┘
         │                       │
         │    ┌──────────────────┘
         │    │
         ▼    ▼
┌─────────────────────────────────────────────────────────────────┐
│                          DEVICE                                   │
│                                                                  │
│  imports:                                                        │
│   - github.com/tamzrod/forge/memory                            │
│   - github.com/tamzrod/forge/models                            │
│                                                                  │
│  exports:                                                        │
│   - Device struct                                                │
│   - DeviceID type                                                │
│   - Behavior interface                                           │
│   - BehaviorFunc adapter                                         │
│   - ModelProvider interface                                      │
│                                                                  │
│  sub-packages:                                                   │
│   - registry/registry.go                                         │
│   - weatherstation/station.go                                     │
│   - (future: pv_inverter, revenue_meter, relay, etc.)            │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────┐
│                         MEMORY                                   │
│                                                                  │
│  imports:                                                        │
│   - (encoding/binary, errors, fmt, math, sync)                  │
│                                                                  │
│  exports:                                                        │
│   - MemoryImage struct                                           │
│   - Region struct                                                │
│   - Quality type + constants                                     │
│   - Error variables                                              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Component Interfaces

### Runtime (runtime.go)

```go
type Runtime struct {
    config  Config
    sched   *Scheduler
    models  map[ModelID]Model
    devices map[DeviceID]*Device
}

type Config struct {
    TickInterval time.Duration
    MaxDevices  int
}

func New(cfg Config) *Runtime
func NewFromFile(path string) (*Runtime, error)
func LoadConfig(path string) (Config, error)
func DefaultConfig() Config

func (r *Runtime) CreateDevice(id DeviceID, typeName string, memRegions map[string]uint32) *Device
func (r *Runtime) CreateModel(m Model)
func (r *Runtime) CreateGridModel(id ModelID) *GridModel
func (r *Runtime) CreateSunModel(id ModelID) *SunModel
func (r *Runtime) CreateWindModel(id ModelID) *WindModel
func (r *Runtime) CreateWeatherModel(id ModelID) *WeatherModel
func (r *Runtime) CreateReservoirModel(id ModelID, area float32) *ReservoirModel

func (r *Runtime) Device(id DeviceID) *Device
func (r *Runtime) Devices() []*Device
func (r *Runtime) Model(id ModelID) Model
func (r *Runtime) Models() []Model

func (r *Runtime) Start()
func (r *Runtime) Stop()
func (r *Runtime) Run(ctx context.Context) error
func (r *Runtime) Shutdown() error
```

### Scheduler (scheduler.go)

```go
type Scheduler struct {
    mu          sync.Mutex
    devices     []*device.Device
    models      []models.Model
    tickInterval time.Duration
    clock      SimulationClock
    running    bool
    stopCh     chan struct{}
}

type SimulationClock struct {
    elapsed   time.Duration
    tickCount uint64
}

func New(tickInterval time.Duration) *Scheduler
func (s *Scheduler) AddDevice(d *device.Device)
func (s *Scheduler) RemoveDevice(id device.DeviceID)
func (s *Scheduler) AddModel(m models.Model)
func (s *Scheduler) RemoveModel(id models.ModelID)
func (s *Scheduler) Run(ctx context.Context) error
func (s *Scheduler) Stop()
func (s *Scheduler) Pause()
func (s *Scheduler) Resume()
func (s *Scheduler) Clock() SimulationClock
```

### Device (device.go)

```go
type DeviceID string

type Device struct {
    id          DeviceID
    typeName    string
    mem         *memory.MemoryImage
    behaviors   []Behavior
    modelGetter func(id models.ModelID) models.Model
    running     bool
}

func New(id DeviceID, typeName string, memRegions map[string]uint32) *Device
func (d *Device) ID() DeviceID
func (d *Device) Type() string
func (d *Device) Memory() *memory.MemoryImage
func (d *Device) AddBehavior(b Behavior)
func (d *Device) Behaviors() []Behavior
func (d *Device) Tick()
func (d *Device) Start()
func (d *Device) Stop()
func (d *Device) Running() bool
func (d *Device) Model(id models.ModelID) models.Model
```

### Behavior Interface (behavior.go)

```go
type Behavior interface {
    ID() string
    Attach(d *Device)
    Detach()
    Tick()
}

type ModelObserver interface {
    ObserveModel(id models.ModelID) models.Model
}

type BehaviorFunc struct {
    idFunc     func() string
    attachFunc func(*Device)
    detachFunc func()
    tickFunc   func()
}
```

### MemoryImage (memory.go)

```go
type MemoryImage struct {
    mu      sync.RWMutex
    regions map[string]*Region
}

type Region struct {
    Name    string
    Size    uint32
    Values  []byte
    Quality []Quality
}

type Quality uint8

const (
    QualityGood     Quality = 0x00
    QualityUncertain Quality = 0x40
    QualityBad      Quality = 0x80
    QualityOffline  Quality = 0x84
)

func New(regionDefs map[string]uint32) *MemoryImage
func (m *MemoryImage) Read(region string, address uint32, size uint32) ([]byte, error)
func (m *MemoryImage) Write(region string, address uint32, data []byte) error
func (m *MemoryImage) ReadUint16(region string, address uint32) (uint16, error)
func (m *MemoryImage) WriteUint16(region string, address uint32, value uint16) error
func (m *MemoryImage) ReadFloat32(region string, address uint32) (float32, error)
func (m *MemoryImage) WriteFloat32(region string, address uint32, value float32) error
func (m *MemoryImage) Quality(region string, address uint32) (Quality, error)
func (m *MemoryImage) SetQuality(region string, address uint32, quality Quality) error
```

---

## Execution Flow

```
Tick Loop (Scheduler)
         │
         ▼
┌─────────────────────────────────────────────────────────────────────┐
│  1. MODELS TICK (in registration order)                            │
│                                                                     │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐    │
│  │  Grid   │ │   Sun   │ │  Wind   │ │ Weather │ │Reservoir│    │
│  │  Model  │ │  Model  │ │  Model  │ │  Model  │ │  Model  │    │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘    │
│       │             │           │           │           │           │
│       └─────────────┴───────────┴───────────┴───────────┘           │
│                         │                                             │
└─────────────────────────┼───────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────────┐
│  2. DEVICES TICK (in registration order)                            │
│                                                                     │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐      │
│  │   Weather     │    │     PV       │    │   Revenue    │      │
│  │   Station     │    │   Inverter   │    │    Meter     │      │
│  │               │    │              │    │              │      │
│  │  ┌─────────┐ │    │  ┌─────────┐ │    │  ┌─────────┐ │      │
│  │  │Behavior1│ │    │  │Behavior1│ │    │  │Behavior1│ │      │
│  │  └────┬────┘ │    │  └────┬────┘ │    │  └────┬────┘ │      │
│  │       │       │    │       │       │    │       │       │      │
│  │  ┌────▼────┐ │    │  ┌────▼────┐ │    │  ┌────▼────┐ │      │
│  │  │  Read   │ │    │  │  Read   │ │    │  │  Read   │ │      │
│  │  │ Weather │ │    │  │  Sun    │ │    │  │  Grid   │ │      │
│  │  │  Model  │ │    │  │  Model  │ │    │  │  Model  │ │      │
│  │  └────┬────┘ │    │  └────┬────┘ │    │  └────┬────┘ │      │
│  │       │       │    │       │       │    │       │       │      │
│  │  ┌────▼────┐ │    │  ┌────▼────┐ │    │  ┌────▼────┐ │      │
│  │  │ Write   │ │    │  │ Calculate│ │    │  │ Integrate│ │      │
│  │  │ Memory  │ │    │  │   Power │ │    │  │  Energy │ │      │
│  │  └────┬────┘ │    │  └────┬────┘ │    │  └────┬────┘ │      │
│  │       │       │    │       │       │    │       │       │      │
│  │  ┌────▼────┐ │    │  ┌────▼────┐ │    │  ┌────▼────┐ │      │
│  │  │ Memory  │ │    │  │ Write   │ │    │  │ Write   │ │      │
│  │  │  Image  │ │    │  │ Memory  │ │    │  │ Memory  │ │      │
│  │  └─────────┘ │    │  └─────────┘ │    │  └─────────┘ │      │
│  └──────────────┘    └──────────────┘    └──────────────┘      │
└─────────────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────────┐
│  3. ADVANCE CLOCK                                                  │
│                                                                     │
│  clock.elapsed += tickInterval                                     │
│  clock.tickCount++                                                 │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Registry Pattern

The internal packages provide device type registries:

```
┌─────────────────────────────────────────────────────────────────────┐
│                     INTERNAL PACKAGES                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  internal/devices/                                                  │
│  ├── device.go         (core Device struct)                       │
│  ├── registry.go       (DeviceRegistry)                          │
│  ├── memory.go         (device memory helpers)                     │
│  └── weatherstation/                                           │
│      └── station.go    (WeatherStation type)                      │
│                                                                     │
│  internal/models/                                                  │
│  ├── clock.go          (ClockModel)                               │
│  ├── grid.go           (GridModel)                                │
│  ├── sun.go            (SunModel)                                  │
│  ├── wind.go           (WindModel)                                 │
│  └── weather.go        (WeatherModel)                              │
│                                                                     │
│  internal/publishers/                                              │
│  └── rawingest/         (MMA2 Raw Ingest)                         │
│      ├── publisher.go    (RawIngestPublisher)                      │
│      ├── protocol.go     (protocol definitions)                   │
│      └── config.go       (configuration)                          │
│                                                                     │
│  internal/inspector/                                               │
│  ├── server.go        (Inspector API)                             │
│  ├── dashboard.go      (dashboard views)                           │
│  └── view.go           (view definitions)                          │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Future Extensibility

```
Plugin System (planned)
         │
         ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Plugin Interface                                                    │
│                                                                     │
│  type Plugin interface {                                            │
│      Models() []Model                                               │
│      DeviceTypes() []DeviceType                                     │
│      Initialize(runtime *Runtime)                                   │
│  }                                                                   │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  Energy Plugin (implemented)     Water Plugin (planned)             │
│  ├── GridModel                   ├── ReservoirModel               │
│  ├── SunModel                    ├── RiverModel                    │
│  ├── WindModel                   ├── HydraulicNetworkModel         │
│  ├── WeatherModel                │                                │
│  │                              ├── PumpDevice                    │
│  ├── WeatherStation              ├── ValveDevice                   │
│  ├── PVInverter                  ├── TankDevice                    │
│  ├── RevenueMeter                └── FlowMeterDevice               │
│  └── Relay                                                        │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Key Design Decisions

| Decision | Rationale |
|----------|-----------|
| **forge package as facade** | Provides clean public API without exposing internal structure |
| **Device owns Memory** | Memory is the source of truth; devices encapsulate their state |
| **Behavior as interface** | Enables plugin architecture for new device types |
| **Scheduler controls tick order** | Deterministic execution requires centralized scheduling |
| **Models are separate from devices** | Physics (models) should not know about equipment (devices) |
| **Memory uses RWMutex** | Read-heavy workload; writers should not block readers |

---

## Verification Points

To verify this diagram is accurate, check:

1. **Every import has a valid dependency** - No circular imports
2. **Interfaces match implementations** - Device implements ModelProvider
3. **Scheduler owns tick order** - Models tick before devices
4. **Device owns memory** - No external code modifies device memory directly
5. **Behaviors observe models** - Behaviors read from models, write to memory

---

## Related Documents

| Document | Purpose |
|----------|---------|
| [Context Diagram](context-diagram.md) | System boundaries and external actors |
| [Runtime](runtime.md) | Runtime component details |
| [Device Model](device-model.md) | Device anatomy |
| [Scheduler](scheduler.md) | Time advancement details |

---

*Created: 2026-07-13*  
*Type: Architecture Artifact*  
*Status: Initial*
