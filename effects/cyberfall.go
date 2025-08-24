package effects

import (
	"fmt"
	"math/rand"
	"time"

	"godmx/dmx"
	"godmx/types"
)

/*
Effect Name: Cyberfall
Description: Simulates digital rain, acting as a brightness mask over existing colors.
Tags: [transparent, brightness_mask, random]
Parameters:
  - InternalName: speed
    DisplayName: Speed
    Description: How fast the 'rain' falls.
    DataType: float64
    DefaultValue: 1.0
    MinValue: 0.0
  - InternalName: density
    DisplayName: Density
    Description: How many 'active' columns are falling (0.0 - 1.0).
    DataType: float64
    DefaultValue: 0.5
    MinValue: 0.0
    MaxValue: 1.0
  - InternalName: trail_length
    DisplayName: Trail Length
    Description: How long the 'tail' of the falling light is (in lamps).
    DataType: int
    DefaultValue: 10
    MinValue: 0
  - InternalName: min_brightness
    DisplayName: Min Brightness
    Description: Minimum brightness for dark parts (0-255).
    DataType: int
    DefaultValue: 0
    MinValue: 0
    MaxValue: 255
  - InternalName: max_brightness
    DisplayName: Max Brightness
    Description: Maximum brightness for bright parts (0-255).
    DataType: int
    DefaultValue: 255
    MinValue: 0
    MaxValue: 255
  - InternalName: flicker_intensity
    DisplayName: Flicker Intensity
    Description: Random variation applied to brightness (0.0 - 1.0).
    DataType: float64
    DefaultValue: 0.1
    MinValue: 0.0
    MaxValue: 1.0
*/
func init() {
	rand.Seed(time.Now().UnixNano())
	RegisterEffect("cyberfall", NewCyberfall)
	RegisterEffectMetadata("cyberfall", types.EffectMetadata{
		HumanReadableName: "Cyberfall",
		Description:       "Simulates digital rain, acting as a brightness mask over existing colors.",
		Tags:              []string{"transparent", "brightness_mask", "random"},
		Parameters: []types.ParameterMetadata{
			{
				InternalName: "speed",
				DisplayName:  "Speed",
				Description:  "How fast the 'rain' falls.",
				DataType:     "float64",
				DefaultValue: 1.0,
				MinValue:     0.0,
			},
			{
				InternalName: "density",
				DisplayName:  "Density",
				Description:  "How many 'active' columns are falling (0.0 - 1.0).",
				DataType:     "float64",
				DefaultValue: 0.5,
				MinValue:     0.0,
				MaxValue:     1.0,
			},
			{
				InternalName: "trail_length",
				DisplayName:  "Trail Length",
				Description:  "How long the 'tail' of the falling light is (in lamps).",
				DataType:     "int",
				DefaultValue: 10,
				MinValue:     0,
			},
			{
				InternalName: "min_brightness",
				DisplayName:  "Min Brightness",
				Description:  "Minimum brightness for dark parts (0-255).",
				DataType:     "int",
				DefaultValue: 0,
				MinValue:     0,
				MaxValue:     255,
			},
			{
				InternalName: "max_brightness",
				DisplayName:  "Max Brightness",
				Description:  "Maximum brightness for bright parts (0-255).",
				DataType:     "int",
				DefaultValue: 255,
				MinValue:     0,
				MaxValue:     255,
			},
			{
				InternalName: "flicker_intensity",
				DisplayName:  "Flicker Intensity",
				Description:  "Random variation applied to brightness (0.0 - 1.0).",
				DataType:     "float64",
				DefaultValue: 0.1,
				MinValue:     0.0,
				MaxValue:     1.0,
			},
		},
	})
}

// Cyberfall effect simulates digital rain, acting as a brightness mask.
type Cyberfall struct {
	Speed         float64 // How fast the "rain" falls (e.g., 1.0 for normal, 2.0 for double speed)
	Density       float64 // How many "active" columns are falling (0.0 - 1.0)
	TrailLength   int     // How long the "tail" of the falling light is (in lamps)
	MinBrightness uint8   // Minimum brightness for dark parts (0-255)
	MaxBrightness uint8   // Maximum brightness for bright parts (0-255)
	FlickerIntensity float64 // Random variation applied to brightness (0.0 - 1.0)

	// Internal state for each lamp's "rain" column
	// Value represents the current position of the "head" of the rain for that column
	// -1 means no rain falling in that column
	lampStates []float64 
	lastUpdate time.Time
}

// NewCyberfall creates a new Cyberfall effect.
func NewCyberfall(args map[string]interface{}) (types.Effect, error) {
	speed, ok := args["speed"].(float64)
	if !ok {
		return nil, fmt.Errorf("cyberfall effect: missing or invalid 'speed' parameter")
	}
	density, ok := args["density"].(float64)
	if !ok {
		return nil, fmt.Errorf("cyberfall effect: missing or invalid 'density' parameter")
	}
	trailLength, ok := args["trail_length"].(float64)
	if !ok {
		return nil, fmt.Errorf("cyberfall effect: missing or invalid 'trail_length' parameter")
	}
	minBrightness, ok := args["min_brightness"].(float64)
	if !ok {
		return nil, fmt.Errorf("cyberfall effect: missing or invalid 'min_brightness' parameter")
	}
	maxBrightness, ok := args["max_brightness"].(float64)
	if !ok {
		return nil, fmt.Errorf("cyberfall effect: missing or invalid 'max_brightness' parameter")
	}
	flickerIntensity, ok := args["flicker_intensity"].(float64)
	if !ok {
		return nil, fmt.Errorf("cyberfall effect: missing or invalid 'flicker_intensity' parameter")
	}

	return &Cyberfall{
		Speed:         speed,
		Density:       density,
		TrailLength:   int(trailLength),
		MinBrightness: uint8(minBrightness),
		MaxBrightness: uint8(maxBrightness),
		FlickerIntensity: flickerIntensity,
		lastUpdate:    time.Now(),
	}, nil
}

// Process applies the Cyberfall effect as a brightness mask to the lamps.
func (c *Cyberfall) Process(lamps []dmx.Lamp, globals *types.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	numLamps := len(lamps)
	if numLamps == 0 {
		return
	}

	if c.lampStates == nil || len(c.lampStates) != numLamps {
		c.lampStates = make([]float64, numLamps)
		for i := range c.lampStates {
			c.lampStates[i] = -1.0 // Initialize to no rain
		}
	}

	// Calculate time elapsed since last update
	now := time.Now()
	deltaTime := now.Sub(c.lastUpdate).Seconds()
	c.lastUpdate = now

	// Update rain positions
	for i := 0; i < numLamps; i++ {
		// Advance existing rain
		if c.lampStates[i] >= 0 {
			c.lampStates[i] += c.Speed * deltaTime * float64(numLamps) / 5.0 // Scale speed by numLamps
			if c.lampStates[i] >= float64(numLamps + c.TrailLength) { // Rain has fallen off
				c.lampStates[i] = -1.0
			}
		}

		// Start new rain
		if c.lampStates[i] < 0 && rand.Float64() < c.Density * deltaTime * 2.0 { // Probability based on density and time
			c.lampStates[i] = 0.0 // Start at the top
		}
	}

	// Apply mask
	for i := range lamps {
		lamp := &lamps[i] // Get a pointer to modify in place
		maskBrightness := float64(c.MinBrightness) // Default to min brightness

		if c.lampStates[i] >= 0 {
			// Calculate position within the trail (0.0 at head, 1.0 at tail end)
			trailPos := (c.lampStates[i] - float64(i)) / float64(c.TrailLength)

			if trailPos >= 0.0 && trailPos <= 1.0 {
				// Brightness falls off towards the tail
				// Linear falloff for simplicity, can be changed to exponential/sinusoidal
				brightnessFactor := 1.0 - trailPos

				// Apply flicker
				if c.FlickerIntensity > 0 {
					brightnessFactor += (rand.Float64()*2 - 1) * c.FlickerIntensity
					if brightnessFactor < 0 { brightnessFactor = 0 }
					if brightnessFactor > 1 { brightnessFactor = 1 }
				}

				maskBrightness = float64(c.MinBrightness) + brightnessFactor * float64(c.MaxBrightness - c.MinBrightness)
			}
		}

		// Apply brightness mask to R, G, B, W components of dmx.Lamp
		// Scale maskBrightness to a factor between 0.0 and 1.0
		brightnessMultiplier := maskBrightness / 255.0

		lamp.R = uint8(float64(lamp.R) * brightnessMultiplier)
		lamp.G = uint8(float64(lamp.G) * brightnessMultiplier)
		lamp.B = uint8(float64(lamp.B) * brightnessMultiplier)
		lamp.W = uint8(float64(lamp.W) * brightnessMultiplier)
	}
}
