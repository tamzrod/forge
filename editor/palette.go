// Package editor provides the Forge Editor for creating and editing simulation models.
package editor

// Palette provides entity templates for dragging onto the canvas.
type Palette struct {
	Categories []PaletteCategory `json:"categories"`
	Items      []PaletteItem      `json:"items"`
}

// PaletteCategory organizes palette items.
type PaletteCategory struct {
	ID    EntityCategory `json:"id"`
	Name  string        `json:"name"`
	Icon  string        `json:"icon"`
	Order int           `json:"order"`
}

// NewPalette creates a new palette with default items.
func NewPalette() *Palette {
	return &Palette{
		Categories: []PaletteCategory{
			{ID: CategoryElectrical, Name: "Electrical", Icon: "⚡", Order: 1},
			{ID: CategoryEnvironment, Name: "Environment", Icon: "🌤️", Order: 2},
			{ID: CategorySimulation, Name: "Simulation", Icon: "🎬", Order: 3},
		},
		Items: []PaletteItem{
			// Electrical items
			{
				ID:          "palette-grid",
				Name:        "Utility Grid",
				Category:    CategoryElectrical,
				EntityType:  EntityTypeGrid,
				Icon:        "🔌",
				Description: "Utility grid connection point",
				DefaultProperties: Properties{
					"name":               {Value: "Utility Grid", Type: "string"},
					"nominal_voltage":    {Value: float64(69000), Type: "number", Unit: "V"},
					"nominal_frequency":  {Value: float64(60), Type: "number", Unit: "Hz"},
				},
			},
			{
				ID:          "palette-bus",
				Name:        "Bus",
				Category:    CategoryElectrical,
				EntityType:  EntityTypeBus,
				Icon:        "⚫",
				Description: "Electrical bus node",
				DefaultProperties: Properties{
					"name":            {Value: "New Bus", Type: "string"},
					"nominal_voltage": {Value: float64(480), Type: "number", Unit: "V"},
				},
			},
			{
				ID:          "palette-breaker",
				Name:        "Breaker",
				Category:    CategoryElectrical,
				EntityType:  EntityTypeBreaker,
				Icon:        "🔀",
				Description: "Circuit breaker switch",
				DefaultProperties: Properties{
					"name":    {Value: "Circuit Breaker", Type: "string"},
					"is_open": {Value: false, Type: "boolean"},
				},
			},
			{
				ID:          "palette-transformer",
				Name:        "Transformer",
				Category:    CategoryElectrical,
				EntityType:  EntityTypeTransformer,
				Icon:        "🔄",
				Description: "Power transformer",
				DefaultProperties: Properties{
					"name":       {Value: "Transformer", Type: "string"},
					"hv_voltage": {Value: float64(69000), Type: "number", Unit: "V"},
					"lv_voltage": {Value: float64(480), Type: "number", Unit: "V"},
					"rating":     {Value: float64(1000), Type: "number", Unit: "kVA"},
				},
			},
			{
				ID:          "palette-generator",
				Name:        "Virtual Generator",
				Category:    CategoryElectrical,
				EntityType:  EntityTypeGenerator,
				Icon:        "☀️",
				Description: "Solar or wind generator",
				DefaultProperties: Properties{
					"name":               {Value: "Solar Generator", Type: "string"},
					"rated_capacity":     {Value: float64(500), Type: "number", Unit: "kW", Min: floatPtr(0)},
					"available_capacity": {Value: float64(500), Type: "number", Unit: "kW", Min: floatPtr(0)},
					"is_online":          {Value: true, Type: "boolean"},
					"is_dispatchable":    {Value: true, Type: "boolean"},
				},
			},
			{
				ID:          "palette-load",
				Name:        "Virtual Load",
				Category:    CategoryElectrical,
				EntityType:  EntityTypeLoad,
				Icon:        "🏭",
				Description: "Factory or facility load",
				DefaultProperties: Properties{
					"name":                 {Value: "Factory Load", Type: "string"},
					"active_power_demand":  {Value: float64(400), Type: "number", Unit: "kW", Min: floatPtr(0)},
					"power_factor":         {Value: float64(0.9), Type: "number", Min: floatPtr(0), Max: floatPtr(1)},
					"is_connected":        {Value: true, Type: "boolean"},
				},
			},
			{
				ID:          "palette-meter",
				Name:        "Meter",
				Category:    CategoryElectrical,
				EntityType:  EntityTypeMeter,
				Icon:        "📊",
				Description: "Power measurement meter",
				DefaultProperties: Properties{
					"name": {Value: "PCC Meter", Type: "string"},
					"type": {Value: "pcc", Type: "enum", Options: []string{"pcc", "array", "feeder"}},
				},
			},
			// Environment items
			{
				ID:          "palette-sun",
				Name:        "Sun",
				Category:    CategoryEnvironment,
				EntityType:  EntityTypeSun,
				Icon:        "🌞",
				Description: "Solar position and irradiance",
				DefaultProperties: Properties{
					"latitude":     {Value: float64(35.2271), Type: "number"},
					"longitude":    {Value: float64(-80.8431), Type: "number"},
					"tilt":         {Value: float64(20), Type: "number", Unit: "°"},
					"azimuth":      {Value: float64(180), Type: "number", Unit: "°"},
				},
			},
			{
				ID:          "palette-weather",
				Name:        "Weather",
				Category:    CategoryEnvironment,
				EntityType:  EntityTypeWeather,
				Icon:        "🌤️",
				Description: "Weather conditions",
				DefaultProperties: Properties{
					"temperature":  {Value: float64(25), Type: "number", Unit: "°C"},
					"humidity":    {Value: float64(50), Type: "number", Unit: "%"},
					"cloud_cover": {Value: float64(0), Type: "number", Unit: "%", Min: floatPtr(0), Max: floatPtr(100)},
				},
			},
			{
				ID:          "palette-wind",
				Name:        "Wind",
				Category:    CategoryEnvironment,
				EntityType:  EntityTypeWind,
				Icon:        "💨",
				Description: "Wind conditions",
				DefaultProperties: Properties{
					"speed":      {Value: float64(5), Type: "number", Unit: "m/s"},
					"direction":  {Value: float64(0), Type: "number", Unit: "°"},
				},
			},
			// Simulation items
			{
				ID:          "palette-scenario",
				Name:        "Scenario",
				Category:    CategorySimulation,
				EntityType:  EntityTypeScenario,
				Icon:        "🎬",
				Description: "Test scenario",
				DefaultProperties: Properties{
					"name":        {Value: "Test Scenario", Type: "string"},
					"duration":    {Value: float64(3600), Type: "number", Unit: "s"},
					"description": {Value: "", Type: "string"},
				},
			},
			{
				ID:          "palette-clock",
				Name:        "Simulation Clock",
				Category:    CategorySimulation,
				EntityType:  EntityTypeClock,
				Icon:        "⏱️",
				Description: "Simulation time control",
				DefaultProperties: Properties{
					"start_time": {Value: "2024-01-01T08:00:00Z", Type: "string"},
					"end_time":   {Value: "2024-01-01T20:00:00Z", Type: "string"},
					"time_step":  {Value: float64(100), Type: "number", Unit: "ms"},
				},
			},
		},
	}
}

// ItemsByCategory returns palette items filtered by category.
func (p *Palette) ItemsByCategory(category EntityCategory) []PaletteItem {
	var items []PaletteItem
	for _, item := range p.Items {
		if item.Category == category {
			items = append(items, item)
		}
	}
	return items
}

// GetItem returns a palette item by ID.
func (p *Palette) GetItem(id ID) *PaletteItem {
	for _, item := range p.Items {
		if ID(item.ID) == id {
			return &item
		}
	}
	return nil
}

// DefaultPalette returns a singleton default palette.
var DefaultPalette = NewPalette()
