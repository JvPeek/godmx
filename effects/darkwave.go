package effects

import (
	"godmx/dmx"
	"godmx/orchestrator"
	"math"
)

func init() {
	RegisterEffect("darkwave", func(args map[string]interface{}) (orchestrator.Effect, error) {
		return NewDarkWave(args)
	})
	RegisterEffectParameters("darkwave", map[string]interface{}{"percentage": 0.5, "speed": 1.0})
}

// DarkWave is an effect that creates a dark wave along the strip.
	type DarkWave struct {
	Percentage float64
	Speed      float64
}

// NewDarkWave creates a new DarkWave effect.
func NewDarkWave(args map[string]interface{}) (*DarkWave, error) {
	percentage := args["percentage"].(float64)
	speed := args["speed"].(float64)

	return &DarkWave{Percentage: percentage, Speed: speed}, nil
}

// Process applies the DarkWave effect to the lamp strip.
func (dw *DarkWave) Process(lamps []dmx.Lamp, globals *orchestrator.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	step := globals.BeatProgress * 2 * math.Pi * dw.Speed
	for i := 0; i < len(lamps); i++ {
		sinValue := (math.Sin(float64(i)/float64(len(lamps))*2*math.Pi + step) + 1) / 2
		darkness := 1 - (sinValue * dw.Percentage)
		lamps[i] = scaleColor(lamps[i], darkness)
	}
}

func scaleColor(c dmx.Lamp, factor float64) dmx.Lamp {
	return dmx.Lamp{
		R: uint8(math.Min(255, float64(c.R)*factor)),
		G: uint8(math.Min(255, float64(c.G)*factor)),
		B: uint8(math.Min(255, float64(c.B)*factor)),
		W: uint8(math.Min(255, float64(c.W)*factor)),
	}
}
