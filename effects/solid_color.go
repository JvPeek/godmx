package effects

import (
	"godmx/dmx"
	"godmx/orchestrator"
)

func init() {
	RegisterEffect("solidColor", func(args map[string]interface{}) (orchestrator.Effect, map[string]interface{}, error) {
		return &SolidColor{}, args, nil
	})
}

// SolidColor sets all lamps to a single color from global parameters.
type SolidColor struct {
	// No fields needed, color comes from globals
}

// Process applies the solid color to the lamps using globals.Color1.
func (s *SolidColor) Process(lamps []dmx.Lamp, globals *orchestrator.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
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