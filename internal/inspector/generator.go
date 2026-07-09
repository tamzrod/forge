package inspector

import (
	"fmt"
	"time"

	"github.com/tamzrod/forge/internal/devices/weatherstation"
)

// Generator builds GenericInspectorData for different object types.
type Generator struct {
	view *View
}

// NewGenerator creates a new generator with the given view.
func NewGenerator(view *View) *Generator {
	return &Generator{view: view}
}

// Inspect returns inspection data for the specified object.
func (g *Generator) Inspect(objectID string) (*GenericInspectorData, error) {
	switch objectID {
	case "world":
		return g.InspectWorld()
	case "clock":
		return g.InspectClock()
	case "sun":
		return g.InspectSun()
	case "weather":
		return g.InspectWeather()
	case "grid":
		return g.InspectGrid()
	default:
		// Check if it's a device
		if len(objectID) > 7 && objectID[:7] == "device-" {
			deviceID := objectID[7:]
			return g.InspectDevice(deviceID)
		}
		// Check for nested device objects
		if len(objectID) > 8 {
			parts := splitObjectID(objectID)
			if len(parts) >= 2 && parts[0] == "device" {
				deviceID := parts[1]
				subPath := ""
				if len(parts) > 2 {
					subPath = parts[2]
				}
				return g.InspectDeviceChild(deviceID, subPath)
			}
		}
		return nil, fmt.Errorf("unknown object: %s", objectID)
	}
}

func splitObjectID(id string) []string {
	var parts []string
	current := ""
	for i, c := range id {
		if c == '-' && i > 0 && id[i-1] != '\\' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	parts = append(parts, current)
	return parts
}

// InspectWorld returns inspection data for the simulation world.
func (g *Generator) InspectWorld() (*GenericInspectorData, error) {
	state := g.view.FullState()

	data := &GenericInspectorData{
		Object: ObjectIdentity{
			ID:   "world",
			Type: ObjectTypeSimulation,
			Name: "Simulation World",
		},
		Sections: []*Section{
			{
				ID:    SectionOverview,
				Title: "Simulation Summary",
				Icon:  "activity",
				Properties: []*Property{
					BoolProperty("Is Paused", state.Clock.IsPaused),
					DurationProperty("Elapsed", state.Clock.Elapsed),
					IntProperty("Tick Count", int64(state.Clock.TickCount)),
					IntProperty("Devices", int64(state.Devices.Count)),
				},
			},
			{
				ID:    SectionChildren,
				Title: "Simulation Models",
				Icon:  "layers",
				Children: []*ObjectRef{
					{ID: "clock", Type: ObjectTypeSimulation, Name: "Clock"},
					{ID: "sun", Type: ObjectTypeSimulation, Name: "Sun"},
					{ID: "weather", Type: ObjectTypeSimulation, Name: "Weather"},
					{ID: "grid", Type: ObjectTypeSimulation, Name: "Grid"},
				},
			},
		},
	}

	return data, nil
}

// InspectClock returns inspection data for the Clock model.
func (g *Generator) InspectClock() (*GenericInspectorData, error) {
	state := g.view.ClockState()

	var status string
	if state.IsPaused {
		status = "Paused"
	} else {
		status = "Running"
	}

	data := &GenericInspectorData{
		Object: ObjectIdentity{
			ID:   "clock",
			Type: ObjectTypeSimulation,
			Name: "Clock",
		},
		Sections: []*Section{
			{
				ID:    SectionIdentity,
				Title: "Identity",
				Icon:  "tag",
				Properties: []*Property{
					TextProperty("ID", "clock"),
					TextProperty("Type", "Simulation Clock"),
					TextProperty("Name", "Simulation Clock"),
				},
			},
			{
				ID:    SectionOverview,
				Title: "Overview",
				Icon:  "eye",
				Properties: []*Property{
					StatusProperty("Status", status),
					DurationProperty("Elapsed", state.Elapsed),
					IntProperty("Tick Count", int64(state.TickCount)),
					EnumProperty("Mode", state.Mode, []string{"Realtime", "Accelerated", "FixedStep"}),
				},
			},
			{
				ID:    SectionState,
				Title: "State",
				Icon:  "activity",
				Properties: []*Property{
					BoolProperty("Is Paused", state.IsPaused),
					DurationProperty("Elapsed (ns)", state.Elapsed),
					IntProperty("Tick Count", int64(state.TickCount)),
					TextProperty("Mode", state.Mode),
				},
			},
			{
				ID:    SectionConfiguration,
				Title: "Configuration",
				Icon:  "settings",
				Properties: []*Property{
					EnumProperty("Mode", state.Mode, []string{"Realtime", "Accelerated", "FixedStep"}),
					NumberProperty("Tick Rate", 10, "Hz"),
				},
			},
			{
				ID:    SectionDiagnostics,
				Title: "Diagnostics",
				Icon:  "alert-triangle",
				Properties: []*Property{
					StatusProperty("Health", "OK"),
					IntProperty("Total Ticks", int64(state.TickCount)),
					NumberProperty("Tick Rate", 10, "Hz"),
				},
			},
		},
	}

	return data, nil
}

// InspectSun returns inspection data for the Sun model.
func (g *Generator) InspectSun() (*GenericInspectorData, error) {
	state := g.view.SunState()

	var status string
	if state.IsDaytime {
		status = "Daytime"
	} else {
		status = "Nighttime"
	}

	data := &GenericInspectorData{
		Object: ObjectIdentity{
			ID:   "sun",
			Type: ObjectTypeSimulation,
			Name: "Sun",
		},
		Sections: []*Section{
			{
				ID:    SectionIdentity,
				Title: "Identity",
				Icon:  "tag",
				Properties: []*Property{
					TextProperty("ID", "sun"),
					TextProperty("Type", "Solar Model"),
					TextProperty("Name", "Sun Model"),
				},
			},
			{
				ID:    SectionOverview,
				Title: "Overview",
				Icon:  "eye",
				Properties: []*Property{
					StatusProperty("Status", status),
					AngleProperty("Elevation", state.Elevation),
					AngleProperty("Azimuth", state.Azimuth),
					BoolProperty("Is Daytime", state.IsDaytime),
				},
			},
			{
				ID:    SectionState,
				Title: "State",
				Icon:  "activity",
				Properties: []*Property{
					AngleProperty("Elevation", state.Elevation),
					AngleProperty("Azimuth", state.Azimuth),
					NumberProperty("GHI", state.Irradiance, "W/m²"),
					NumberProperty("DNI", state.DirectNormal, "W/m²"),
					NumberProperty("Diffuse", state.Diffuse, "W/m²"),
					BoolProperty("Is Daytime", state.IsDaytime),
				},
			},
			{
				ID:    SectionConfiguration,
				Title: "Configuration",
				Icon:  "settings",
				Properties: []*Property{
					AngleProperty("Latitude", state.Latitude),
					AngleProperty("Longitude", state.Longitude),
				},
			},
			{
				ID:    SectionDiagnostics,
				Title: "Diagnostics",
				Icon:  "alert-triangle",
				Properties: []*Property{
					StatusProperty("Model Health", "OK"),
					TimestampProperty("Last Update", time.Now().Add(-100*time.Millisecond)),
				},
			},
		},
	}

	return data, nil
}

// InspectWeather returns inspection data for the Weather model.
func (g *Generator) InspectWeather() (*GenericInspectorData, error) {
	state := g.view.WeatherState()

	data := &GenericInspectorData{
		Object: ObjectIdentity{
			ID:   "weather",
			Type: ObjectTypeSimulation,
			Name: "Weather",
		},
		Sections: []*Section{
			{
				ID:    SectionIdentity,
				Title: "Identity",
				Icon:  "tag",
				Properties: []*Property{
					TextProperty("ID", "weather"),
					TextProperty("Type", "Weather Model"),
					TextProperty("Name", "Weather Model"),
				},
			},
			{
				ID:    SectionOverview,
				Title: "Overview",
				Icon:  "eye",
				Properties: []*Property{
					NumberProperty("Temperature", state.Temperature, "°C"),
					PercentageProperty("Humidity", state.Humidity/100),
					NumberProperty("Pressure", state.Pressure, "hPa"),
					PercentageProperty("Cloud Cover", state.CloudCover),
					BoolProperty("Is Raining", state.IsRaining),
				},
			},
			{
				ID:    SectionState,
				Title: "State",
				Icon:  "activity",
				Properties: []*Property{
					NumberProperty("Temperature", state.Temperature, "°C").WithPrecision(2),
					PercentageProperty("Humidity", state.Humidity/100),
					NumberProperty("Pressure", state.Pressure, "hPa").WithPrecision(2),
					PercentageProperty("Cloud Cover", state.CloudCover),
					NumberProperty("Wind Speed", state.WindSpeed, "m/s"),
					AngleProperty("Wind Direction", state.WindDirection),
					BoolProperty("Is Raining", state.IsRaining),
				},
			},
			{
				ID:    SectionConfiguration,
				Title: "Configuration",
				Icon:  "settings",
				Properties: []*Property{
					TextProperty("Location", "40°N, 105°W"),
					NumberProperty("Elevation", 1640, "m"),
				},
			},
			{
				ID:    SectionDiagnostics,
				Title: "Diagnostics",
				Icon:  "alert-triangle",
				Properties: []*Property{
					StatusProperty("Model Health", "OK"),
					StatusProperty("Sensor Simulation", "Active"),
					TimestampProperty("Last Update", time.Now().Add(-100*time.Millisecond)),
				},
			},
		},
	}

	return data, nil
}

// InspectGrid returns inspection data for the Grid model.
func (g *Generator) InspectGrid() (*GenericInspectorData, error) {
	state := g.view.GridState()

	var status string
	if state.IsStable {
		status = "Stable"
	} else {
		status = "Unstable"
	}

	data := &GenericInspectorData{
		Object: ObjectIdentity{
			ID:   "grid",
			Type: ObjectTypeSimulation,
			Name: "Grid",
		},
		Sections: []*Section{
			{
				ID:    SectionIdentity,
				Title: "Identity",
				Icon:  "tag",
				Properties: []*Property{
					TextProperty("ID", "grid"),
					TextProperty("Type", "Grid Model"),
					TextProperty("Name", "Grid Model"),
				},
			},
			{
				ID:    SectionOverview,
				Title: "Overview",
				Icon:  "eye",
				Properties: []*Property{
					StatusProperty("Status", status),
					NumberProperty("Voltage", state.Voltage, "V"),
					NumberProperty("Frequency", state.Frequency, "Hz"),
					NumberProperty("Voltage PU", state.VoltagePU, "").WithPrecision(4),
				},
			},
			{
				ID:    SectionState,
				Title: "State",
				Icon:  "activity",
				Properties: []*Property{
					NumberProperty("Voltage", state.Voltage, "V").WithPrecision(4),
					NumberProperty("Frequency", state.Frequency, "Hz").WithPrecision(6),
					NumberProperty("Voltage PU", state.VoltagePU, "").WithPrecision(6),
					NumberProperty("Frequency PU", state.FrequencyPU, "").WithPrecision(6),
					NumberProperty("Active Balance", state.ActiveBalance, "MW"),
					NumberProperty("Reactive Balance", state.ReactiveBalance, "MVAr"),
					BoolProperty("Is Stable", state.IsStable),
				},
			},
			{
				ID:    SectionConfiguration,
				Title: "Configuration",
				Icon:  "settings",
				Properties: []*Property{
					NumberProperty("Nominal Voltage", state.NominalVoltage, "V"),
					NumberProperty("Nominal Frequency", state.NominalFrequency, "Hz"),
				},
			},
			{
				ID:    SectionDiagnostics,
				Title: "Diagnostics",
				Icon:  "alert-triangle",
				Properties: []*Property{
					StatusProperty("Model Health", status),
					StatusProperty("Power Flow", "Converged"),
					StatusProperty("Stability", status),
				},
			},
		},
	}

	return data, nil
}

// InspectDevice returns inspection data for a device.
func (g *Generator) InspectDevice(deviceID string) (*GenericInspectorData, error) {
	// Find the device
	devices := g.view.DevicesState()
	var device DeviceState
	found := false
	for _, d := range devices.Devices {
		if d.ID == deviceID {
			device = d
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("device not found: %s", deviceID)
	}

	data := &GenericInspectorData{
		Object: ObjectIdentity{
			ID:   deviceID,
			Type: ObjectTypeDevice,
			Name: device.Name,
		},
		Sections: []*Section{
			{
				ID:    SectionIdentity,
				Title: "Identity",
				Icon:  "tag",
				Properties: []*Property{
					TextProperty("ID", device.ID),
					TextProperty("Type", string(device.Type)),
					TextProperty("Name", device.Name),
				},
			},
			{
				ID:    SectionOverview,
				Title: "Overview",
				Icon:  "eye",
				Properties: []*Property{
					StatusProperty("State", device.State),
					BoolProperty("Interface Enabled", device.InterfaceEnabled),
				},
			},
			{
				ID:    SectionState,
				Title: "State",
				Icon:  "activity",
				Properties: []*Property{
					StatusProperty("Device State", device.State),
					TextProperty("Device Type", string(device.Type)),
					BoolProperty("Interface Enabled", device.InterfaceEnabled),
				},
			},
		},
	}

	// Add Communications section if interface is enabled
	if device.Interface != nil {
		iface := device.Interface
		data.Sections = append(data.Sections, &Section{
			ID:    SectionCommunications,
			Title: "Communications",
			Icon:  "radio",
			Properties: []*Property{
				BoolProperty("Interface Enabled", iface.Enabled),
				StatusProperty("Connection Status", func() string {
					if iface.Connected {
						return "Connected"
					}
					return "Disconnected"
				}()),
				IntProperty("Packets Sent", int64(iface.PacketsSent)),
				IntProperty("Errors", int64(iface.Errors)),
			},
		})

		// Add last error if present
		if iface.LastError != "" {
			lastErrSection := data.GetSection(SectionCommunications)
			if lastErrSection != nil {
				lastErrSection.Properties = append(lastErrSection.Properties,
					TextProperty("Last Error", iface.LastError))
			}
		}
	}

	// Check if this is a Weather Station and add Memory section
	if ws, ok := g.getWeatherStation(deviceID); ok {
		wsState := ws.State()
		data.Sections = append(data.Sections, &Section{
			ID:    SectionMemory,
			Title: "Device Memory",
			Icon:  "cpu",
			Properties: []*Property{
				NumberProperty("Temperature", wsState.Temperature, "°C"),
				NumberProperty("Humidity", wsState.Humidity, "%"),
				NumberProperty("Pressure", wsState.Pressure, "hPa"),
				NumberProperty("Cloud Cover", wsState.CloudCover, "%"),
				NumberProperty("Wind Speed", wsState.WindSpeed, "m/s"),
				AngleProperty("Wind Direction", wsState.WindDirection),
				BoolProperty("Rain Status", wsState.RainStatus),
				IntProperty("Tick Count", int64(wsState.TickCount)),
			},
		})

		// Add children for Virtual Firmware
		data.Sections = append(data.Sections, &Section{
			ID:    SectionChildren,
			Title: "Children",
			Icon:  "layers",
			Children: []*ObjectRef{
				{ID: fmt.Sprintf("device-%s-firmware", deviceID), Type: ObjectTypeFirmware, Name: "Virtual Firmware"},
				{ID: fmt.Sprintf("device-%s-memory", deviceID), Type: ObjectTypeMemory, Name: "Device Memory"},
				{ID: fmt.Sprintf("device-%s-interface", deviceID), Type: ObjectTypeInterface, Name: "Communication Interface"},
			},
		})
	}

	return data, nil
}

// InspectDeviceChild returns inspection data for a nested object within a device.
func (g *Generator) InspectDeviceChild(deviceID, subPath string) (*GenericInspectorData, error) {
	switch subPath {
	case "firmware":
		return g.InspectFirmware(deviceID)
	case "memory":
		return g.InspectMemory(deviceID)
	case "interface":
		return g.InspectInterface(deviceID)
	default:
		return nil, fmt.Errorf("unknown device child: %s", subPath)
	}
}

// InspectFirmware returns inspection data for a device's Virtual Firmware.
func (g *Generator) InspectFirmware(deviceID string) (*GenericInspectorData, error) {
	ws, ok := g.getWeatherStation(deviceID)
	if !ok {
		return nil, fmt.Errorf("firmware not found: %s", deviceID)
	}

	wsState := ws.State()

	data := &GenericInspectorData{
		Object: ObjectIdentity{
			ID:   fmt.Sprintf("device-%s-firmware", deviceID),
			Type: ObjectTypeFirmware,
			Name: "Virtual Firmware",
		},
		Sections: []*Section{
			{
				ID:    SectionIdentity,
				Title: "Identity",
				Icon:  "tag",
				Properties: []*Property{
					TextProperty("ID", fmt.Sprintf("device-%s-firmware", deviceID)),
					TextProperty("Type", "Virtual Firmware"),
					TextProperty("Name", "Virtual Firmware"),
					TextProperty("Device", deviceID),
				},
			},
			{
				ID:    SectionOverview,
				Title: "Overview",
				Icon:  "eye",
				Properties: []*Property{
					TextProperty("Firmware Type", "Weather Station"),
					TextProperty("Firmware Version", "1.0.0"),
					TextProperty("Manufacturer", "Forge Labs"),
					IntProperty("Sampling Interval", 1000),
					IntProperty("Memory Regions", 4),
				},
			},
			{
				ID:    SectionState,
				Title: "State",
				Icon:  "activity",
				Properties: []*Property{
					StatusProperty("Device State", wsState.DeviceState.String()),
					IntProperty("Tick Count", int64(wsState.TickCount)),
				},
			},
			{
				ID:    SectionDiagnostics,
				Title: "Diagnostics",
				Icon:  "alert-triangle",
				Properties: []*Property{
					StatusProperty("Firmware Status", "Running"),
					IntProperty("Tick Count", int64(wsState.TickCount)),
					TimestampProperty("Last Tick", time.Now().Add(-50*time.Millisecond)),
				},
			},
		},
	}

	return data, nil
}

// InspectMemory returns inspection data for a device's Device Memory.
func (g *Generator) InspectMemory(deviceID string) (*GenericInspectorData, error) {
	ws, ok := g.getWeatherStation(deviceID)
	if !ok {
		return nil, fmt.Errorf("memory not found: %s", deviceID)
	}

	wsState := ws.State()

	// Build memory region properties
	memProps := []*Property{
		NestedProperty("Measurements", []*Property{
			NumberProperty("Temperature", wsState.Temperature, "°C"),
			NumberProperty("Humidity", wsState.Humidity, "%"),
			NumberProperty("Pressure", wsState.Pressure, "hPa"),
			NumberProperty("Cloud Cover", wsState.CloudCover, "%"),
			NumberProperty("Wind Speed", wsState.WindSpeed, "m/s"),
			AngleProperty("Wind Direction", wsState.WindDirection),
		}),
		NestedProperty("Status", []*Property{
			BoolProperty("Rain Status", wsState.RainStatus),
			StatusProperty("Device Status", func() string {
				if wsState.DeviceState.String() == "Running" {
					return "OK"
				}
				return wsState.DeviceState.String()
			}()),
		}),
		NestedProperty("Diagnostics", []*Property{
			IntProperty("Tick Count", int64(wsState.TickCount)),
		}),
	}

	data := &GenericInspectorData{
		Object: ObjectIdentity{
			ID:   fmt.Sprintf("device-%s-memory", deviceID),
			Type: ObjectTypeMemory,
			Name: "Device Memory",
		},
		Sections: []*Section{
			{
				ID:    SectionIdentity,
				Title: "Identity",
				Icon:  "tag",
				Properties: []*Property{
					TextProperty("ID", fmt.Sprintf("device-%s-memory", deviceID)),
					TextProperty("Type", "Device Memory"),
					TextProperty("Name", "Device Memory"),
					TextProperty("Owner", "Virtual Firmware"),
				},
			},
			{
				ID:    SectionOverview,
				Title: "Overview",
				Icon:  "eye",
				Properties: []*Property{
					IntProperty("Total Values", int64(len(wsState.PublishingState().Enabled == false)*8+1)), // Approximate
					TextProperty("Quality", "Good"),
				},
			},
			{
				ID:         SectionMemory,
				Title:      "Memory Contents",
				Icon:       "database",
				Properties: memProps,
			},
		},
	}

	return data, nil
}

// InspectInterface returns inspection data for a device's Communication Interface.
func (g *Generator) InspectInterface(deviceID string) (*GenericInspectorData, error) {
	ws, ok := g.getWeatherStation(deviceID)
	if !ok {
		return nil, fmt.Errorf("interface not found: %s", deviceID)
	}

	pubState := ws.PublishingState()

	var status string
	if pubState.Connected {
		status = "Connected"
	} else if pubState.Enabled {
		status = "Enabled"
	} else {
		status = "Disabled"
	}

	data := &GenericInspectorData{
		Object: ObjectIdentity{
			ID:   fmt.Sprintf("device-%s-interface", deviceID),
			Type: ObjectTypeInterface,
			Name: "Communication Interface",
		},
		Sections: []*Section{
			{
				ID:    SectionIdentity,
				Title: "Identity",
				Icon:  "tag",
				Properties: []*Property{
					TextProperty("ID", fmt.Sprintf("device-%s-interface", deviceID)),
					TextProperty("Type", "Raw Ingest"),
					TextProperty("Name", "Communication Interface"),
				},
			},
			{
				ID:    SectionOverview,
				Title: "Overview",
				Icon:  "eye",
				Properties: []*Property{
					BoolProperty("Enabled", pubState.Enabled),
					StatusProperty("Status", status),
					TextProperty("Interface Type", "Raw Ingest"),
				},
			},
			{
				ID:    SectionCommunications,
				Title: "Traffic Statistics",
				Icon:  "activity",
				Properties: []*Property{
					IntProperty("Packets Sent", int64(pubState.PacketsSent)),
					IntProperty("Errors", int64(pubState.Errors)),
				},
			},
		},
	}

	// Add last error if present
	if pubState.LastError != "" {
		commSection := data.GetSection(SectionCommunications)
		if commSection != nil {
			commSection.Properties = append(commSection.Properties,
				TextProperty("Last Error", pubState.LastError))
		}
	}

	// Add diagnostics
	data.Sections = append(data.Sections, &Section{
		ID:    SectionDiagnostics,
		Title: "Diagnostics",
		Icon:  "alert-triangle",
		Properties: []*Property{
			StatusProperty("Connection Health", func() string {
				if pubState.Errors == 0 {
					return "OK"
				}
				return "Degraded"
			}()),
			TimestampProperty("Last Publish", pubState.LastPublish),
		},
	})

	return data, nil
}

// getWeatherStation retrieves a Weather Station by device ID.
func (g *Generator) getWeatherStation(deviceID string) (*weatherstation.Station, bool) {
	if g.view.registry == nil {
		return nil, false
	}

	for _, d := range g.view.registry.Devices() {
		if string(d.ID()) == deviceID {
			ws, ok := d.(*weatherstation.Station)
			return ws, ok
		}
	}

	return nil, false
}
