package effects

import (
	"godmx/dmx"
	"godmx/types"
)

func init() {
	RegisterEffect("whiteout", func(args map[string]interface{}) (types.Effect, error) {
		return NewWhiteout(), nil
	})
	RegisterEffectMetadata("whiteout", types.EffectMetadata{
		HumanReadableName: "Whiteout",
		Description:       "Sets all lamps to full white, overriding any previous colors.",
		Tags:              []string{"color_source"},
		Parameters:        []types.ParameterMetadata{},
	})
}

// Whiteout sets all lamps to full white.
type Whiteout struct {
	// No fields needed
}

// NewWhiteout creates a new Whiteout effect.
func NewWhiteout() *Whiteout {
	return &Whiteout{}
}

// Process applies the whiteout effect to the lamps.
func (w *Whiteout) Process(lamps []dmx.Lamp, globals *types.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	for i := range lamps {
		if numChannelsPerLamp == 3 && channelMapping == "RGB" {
			lamps[i] = dmx.Lamp{R: 255, G: 255, B: 255, W: 0} // Set RGB to white, W to 0
		} else {
			lamps[i] = dmx.Lamp{R: 255, G: 255, B: 255, W: 255} // Default to RGBW white
		}
	}
}
