package effects

import (
	"godmx/dmx"
	"godmx/orchestrator"
)

func init() {
	RegisterEffect("blink", func(args map[string]interface{}) (orchestrator.Effect, error) {
		return NewBlink(args), nil
	})
	RegisterEffectParameters("blink", map[string]interface{}{"divider": 1}) // Default divider is 1
}

// Blink alternates between two colors based on the global BPM.
type Blink struct {
	Divider int
}

// NewBlink creates a new Blink effect.
func NewBlink(args map[string]interface{}) *Blink {
	divider := 1 // Default value
	if d, ok := args["divider"].(float64); ok { // JSON numbers are float64
		divider = int(d)
	}
	return &Blink{Divider: divider}
}

// Process applies the blink effect to the lamps.
func (b *Blink) Process(lamps []dmx.Lamp, globals *orchestrator.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	var targetColor dmx.Lamp

	// Calculate the current segment within the beat, considering the divider
	// Each beat is divided into Divider * 2 segments (on/off cycles)
	segment := int(globals.BeatProgress * float64(b.Divider*2))

	if segment%2 == 0 {
		// Show Color1 for even segments
		targetColor = globals.Color1
	} else {
		// Show Color2 for odd segments
		targetColor = globals.Color2
	}

	for i := range lamps {
		lamps[i].R = targetColor.R
		lamps[i].G = targetColor.G
		lamps[i].B = targetColor.B
		// Only set W if the channel mapping is RGBW, otherwise set to 0
		if numChannelsPerLamp == 4 && channelMapping == "RGBW" {
			lamps[i].W = targetColor.W
		} else {
			lamps[i].W = 0
		}
	}
}
