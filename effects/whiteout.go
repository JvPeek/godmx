package effects

import (
	"godmx/dmx"
	"godmx/types"
)

/*
Effect Name: Whiteout
Description: Sets all lamps to full white, overriding any previous colors.
Tags: [color_source]
Parameters: []
*/
func init() {
	RegisterEffect("whiteout", NewWhiteout)
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
func NewWhiteout(args map[string]interface{}) (types.Effect, error) {
	return &Whiteout{}, nil
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