package effects

import (
	"godmx/dmx"
	"godmx/types"
	"math"
)

func init() {
	RegisterEffect("darkwave", func(args map[string]interface{}) (types.Effect, error) {
		return NewDarkWave(args)
	})
	RegisterEffectMetadata("darkwave", types.EffectMetadata{
		HumanReadableName: "Darkwave",
		Description:       "Creates a dark wave that travels across the lamps, dimming them based on a sine wave.",
		Tags:              []string{"bpm_sensitive", "transparent", "brightness_mask", "pattern"},
		Parameters: []types.ParameterMetadata{
			{
				InternalName: "percentage",
				DisplayName:  "Percentage",
				Description:  "The maximum percentage of dimming applied by the wave (0.0 - 1.0).",
				DataType:     "float64",
				DefaultValue: 0.5,
				MinValue:     0.0,
				MaxValue:     1.0,
			},
			{
				InternalName: "speed",
				DisplayName:  "Speed",
				Description:  "How fast the dark wave travels.",
				DataType:     "float64",
				DefaultValue: 1.0,
				MinValue:     0.0,
			},
		},
	})
}

// DarkWave is an effect that creates a dark wave along the strip.
	type DarkWave struct {
	Percentage float64
	Speed      float64
}

// NewDarkWave creates a new DarkWave effect.
func NewDarkWave(args map[string]interface{}) (*DarkWave, error) {
	percentage := 0.5 // Default value
	if p, ok := args["percentage"].(float64); ok {
		percentage = p
	}

	speed := 1.0 // Default value
	if s, ok := args["speed"].(float64); ok {
		speed = s
	}

	return &DarkWave{Percentage: percentage, Speed: speed}, nil
}

// Process applies the DarkWave effect to the lamp strip.
func (dw *DarkWave) Process(lamps []dmx.Lamp, globals *types.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
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