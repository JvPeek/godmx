package effects

import (
	"godmx/dmx"
	"godmx/orchestrator"
	"godmx/utils"
)

// Rainbow creates a rainbow effect.
type Rainbow struct {
	counter uint64
}

// NewRainbow creates a new Rainbow effect.
func NewRainbow() *Rainbow {
	return &Rainbow{}
}

// Process applies the rainbow effect to the lamps.
func (r *Rainbow) Process(lamps []dmx.Lamp, globals *orchestrator.OrchestratorGlobals) {
	// Use BPM to influence the speed of the rainbow
	// A higher BPM means a faster rainbow
	speedFactor := globals.BPM / 120.0 // Normalize to default BPM
	r.counter += uint64(speedFactor * 1.0) // Increment counter based on speedFactor

	for i := range lamps {
		hue := (float64(r.counter) + float64(i)*10.0) / 100.0
		r, g, b := utils.HsvToRgb(hue, 1.0, 1.0) // Changed to utils.HsvToRgb
		lamps[i].R = r
		lamps[i].G = g
		lamps[i].B = b
		lamps[i].W = 0 // No white for the rainbow
	}
}
