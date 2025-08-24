package effects

import (
	"godmx/dmx"
	"godmx/types"
	"math"
	"godmx/utils"
)

/*
Effect Name: Gradient
Description: Creates a smooth color gradient across the lamps, interpolating between global Color1 and Color2.
Tags: [color_source, pattern]
Parameters: []
*/
func init() {
	RegisterEffect("gradient", NewGradient)
	RegisterEffectMetadata("gradient", types.EffectMetadata{
		HumanReadableName: "Gradient",
		Description:       "Creates a smooth color gradient across the lamps, interpolating between global Color1 and Color2.",
		Tags:              []string{"color_source", "pattern"},
		Parameters:        []types.ParameterMetadata{},
	})
}

// Gradient creates a color gradient across the lamps.
type Gradient struct {
	// No fields needed, colors come from globals
}

// NewGradient creates a new Gradient effect.
func NewGradient(args map[string]interface{}) (types.Effect, error) {
	return &Gradient{}, nil
}

// Process applies the gradient effect to the lamps using globals.Color1 and globals.Color2.
func (g *Gradient) Process(lamps []dmx.Lamp, globals *types.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	numLamps := float64(len(lamps))

	// Convert Color1 and Color2 to HSV for interpolation
	h1, s1, v1 := utils.RgbToHsv(globals.Color1.R, globals.Color1.G, globals.Color1.B)
	h2, s2, v2 := utils.RgbToHsv(globals.Color2.R, globals.Color2.G, globals.Color2.B)

	// Handle hue interpolation across the color wheel (shortest path)
	if math.Abs(h1-h2) > 0.5 {
		if h1 > h2 {
			h2 += 1.0
		} else {
			h1 += 1.0
		}
	}

	for i := range lamps {
		// Calculate interpolation factor (0.0 for Color1, 1.0 for Color2)
		factor := float64(i) / (numLamps - 1)

		// Interpolate HSV components
		h := h1*(1-factor) + h2*factor
		s := s1*(1-factor) + s2*factor
		v := v1*(1-factor) + v2*factor

		// Convert back to RGB and assign by individual components
		r, g, b := utils.HsvToRgb(h, s, v)
		lamps[i].R = r
		lamps[i].G = g
		lamps[i].B = b
		// Set W to 0 as gradient is typically RGB only
		lamps[i].W = 0
	}
}
