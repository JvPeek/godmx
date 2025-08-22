package effects

import (
	"godmx/dmx"
	"godmx/types"
)

func init() {
	RegisterEffect("blink", func(args map[string]interface{}) (types.Effect, error) {
		return NewBlink(args), nil
	})
	RegisterEffectParameters("blink", map[string]interface{}{
		"divider":   1,
		"dutyCycle": 0.5, // Default duty cycle is 0.5
	})
}

// Blink alternates between two colors based on the global BPM.
type Blink struct {
	Divider   int
	DutyCycle float64
}

// NewBlink creates a new Blink effect.
func NewBlink(args map[string]interface{}) *Blink {
	divider := 1 // Default value
	if d, ok := args["divider"].(float64); ok {
		divider = int(d)
	}
	dutyCycle := 0.5 // Default value
	if dc, ok := args["dutyCycle"].(float64); ok {
		dutyCycle = dc
	}
	return &Blink{Divider: divider, DutyCycle: dutyCycle}
}

// Process applies the blink effect to the lamps.
func (b *Blink) Process(lamps []dmx.Lamp, globals *types.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	var targetColor dmx.Lamp

	// Calculate the current segment within the beat, considering the divider
	// Each beat is divided into Divider * 2 segments (on/off cycles)
	segment := int(globals.BeatProgress * float64(b.Divider*2))

	// Calculate progress within the current segment (0.0 to 1.0)
	progressInSegment := (globals.BeatProgress * float64(b.Divider*2)) - float64(segment)

	if progressInSegment < b.DutyCycle {
		// Show Color1 for the duration of the duty cycle within the segment
		targetColor = globals.Color1
	} else {
		// Show Color2 for the remainder of the segment
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