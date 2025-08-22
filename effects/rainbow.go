package effects

import (
	"godmx/dmx"
	"godmx/types"
	"godmx/utils"
	"math"
)

func init() {
	RegisterEffect("rainbow", func(args map[string]interface{}) (types.Effect, error) {
		return NewRainbow(), nil
	})
	RegisterEffectMetadata("rainbow", types.EffectMetadata{
		HumanReadableName: "Rainbow",
		Description:       "Generates a static rainbow spectrum across the lamps.",
		Tags:              []string{"color_source", "pattern"},
		Parameters:        []types.ParameterMetadata{},
	})
}

// Rainbow creates a rainbow effect.
type Rainbow struct {
	counter float64
}

// NewRainbow creates a new Rainbow effect.
func NewRainbow() *Rainbow {
	return &Rainbow{}
}

// Process applies the rainbow effect to the lamps.
func (r *Rainbow) Process(lamps []dmx.Lamp, globals *types.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	numLamps := float64(len(lamps))
	// const fixedTickRate = 40.0 // Assuming 40 FPS as per current configs

	// Calculate how much the counter should advance per tick to complete one cycle per beat
	// counterIncrementPerTick := (numLamps * globals.BPM) / (60.0 * fixedTickRate)
	// r.counter = math.Mod(r.counter + counterIncrementPerTick, numLamps)

	phaseShift := 0.0 // Static rainbow, no phase shift

	for i := range lamps {
		// Calculate hue: current position in rainbow + offset for each lamp
		// The `r.counter` now directly represents the shift in terms of lamps.
		hue := math.Mod((float64(i) / numLamps) + phaseShift, 1.0)

		// Convert back to RGB and assign
		rgbR, rgbG, rgbB := utils.HsvToRgb(hue, 1.0, 1.0)
		lamps[i].R = rgbR
		lamps[i].G = rgbG
		lamps[i].B = rgbB
		// Set W to 0 as rainbow is typically RGB only
		lamps[i].W = 0
	}
}