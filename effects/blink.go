package effects

import (
	"godmx/dmx"
	"godmx/orchestrator"
)

func init() {
	RegisterEffect("blink", func(args map[string]interface{}) (orchestrator.Effect, error) {
		return NewBlink(), nil
	})
	RegisterEffectParameters("blink", make(map[string]interface{}))
}

// Blink alternates between two colors based on the global BPM.
type Blink struct {
	// No fields needed, uses globals.BeatProgress
}

// NewBlink creates a new Blink effect.
func NewBlink() *Blink {
	return &Blink{}
}

// Process applies the blink effect to the lamps.
func (b *Blink) Process(lamps []dmx.Lamp, globals *orchestrator.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	var targetColor dmx.Lamp
	if globals.BeatProgress < 0.5 {
		// Show Color1 for the first half of the beat
		targetColor = globals.Color1
	} else {
		// Show Color2 for the second half of the beat
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
