package effects

import (
	"godmx/dmx"
	"godmx/orchestrator"
)

// Blink alternates between two colors based on the global BPM.
type Blink struct {
	counter float64
}

// NewBlink creates a new Blink effect.
func NewBlink() *Blink {
	return &Blink{}
}

// Process applies the blink effect to the lamps.
func (b *Blink) Process(lamps []dmx.Lamp, globals *orchestrator.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	// Normalize BPM to a default of 120.0 to get a speed factor.
	// A lower BPM will make the counter increment slower, and a higher BPM will make it increment faster.
	speedFactor := globals.BPM / 120.0
	b.counter += speedFactor

	// We'll use a cycle of 60 steps at 120 BPM for a full on-off blink.
	// This means the color will switch every 30 steps.
	// The modulo operator keeps the counter within the 0-59 range.
	var targetColor dmx.Lamp
	if int(b.counter)%60 < 30 {
		// Show Color1 for the first half of the cycle
		targetColor = globals.Color1
	} else {
		// Show Color2 for the second half of the cycle
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
