package effects

import (
	"fmt"
	"godmx/dmx"
	"godmx/orchestrator"
	"math"
)

func init() {
	RegisterEffect("dim", func(args map[string]interface{}) (orchestrator.Effect, error) {
		return NewDim(args)
	})
	RegisterEffectParameters("dim", map[string]interface{}{"percentage": 0.5})
}

// Dim effect dims all lamps by a specified percentage.
type Dim struct {
	Percentage float64
}

// NewDim creates a new Dim effect.
func NewDim(args map[string]interface{}) (*Dim, error) {
	percentage := args["percentage"].(float64)

	if percentage < 0 || percentage > 1.0 {
		return nil, fmt.Errorf("percentage for dim effect must be between 0.0 and 1.0, got %f", percentage)
	}

	return &Dim{Percentage: percentage}, nil
}

// Process applies the dim effect to the lamps.
func (d *Dim) Process(lamps []dmx.Lamp, globals *orchestrator.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	for i := range lamps {
		lamps[i].R = uint8(math.Round(float64(lamps[i].R) * d.Percentage))
		lamps[i].G = uint8(math.Round(float64(lamps[i].G) * d.Percentage))
		lamps[i].B = uint8(math.Round(float64(lamps[i].B) * d.Percentage))
		lamps[i].W = uint8(math.Round(float64(lamps[i].W) * d.Percentage))
	}
}
