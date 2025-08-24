package effects

import (
	"godmx/dmx"
	"godmx/types"
)

/*
Effect Name: Solid Color
Description: Sets all lamps to a single color defined by global Color1.
Tags: [color_source]
Parameters: []
*/
func init() {
	RegisterEffect("solidColor", NewSolidColor)
	RegisterEffectMetadata("solidColor", types.EffectMetadata{
		HumanReadableName: "Solid Color",
		Description:       "Sets all lamps to a single color defined by global Color1.",
		Tags:              []string{"color_source"},
		Parameters:        []types.ParameterMetadata{},
	})
}

// SolidColor sets all lamps to a single color from global parameters.
type SolidColor struct {
	// No fields needed, color comes from globals
}

// NewSolidColor creates a new SolidColor effect.
func NewSolidColor(args map[string]interface{}) (types.Effect, error) {
	return &SolidColor{}, nil
}

// Process applies the solid color to the lamps using globals.Color1.
func (s *SolidColor) Process(lamps []dmx.Lamp, globals *types.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	for i := range lamps {
		lamps[i].R = globals.Color1.R
		lamps[i].G = globals.Color1.G
		lamps[i].B = globals.Color1.B
		// Only set W if the channel mapping is RGBW, otherwise set to 0
		if numChannelsPerLamp == 4 && channelMapping == "RGBW" {
			lamps[i].W = globals.Color1.W
		} else {
			lamps[i].W = 0
		}
	}
}