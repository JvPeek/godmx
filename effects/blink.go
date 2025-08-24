package effects

import (
	"fmt"
	"godmx/dmx"
	"godmx/types"
)

/*
Effect Name: Blink
Description: Alternates between two colors based on the global BPM, creating a blinking effect.
Tags: [bpm_sensitive, color_source, pattern]
Parameters:
  - InternalName: divider
    DisplayName: Divider
    Description: Divides the beat into segments for faster blinking.
    DataType: int
    DefaultValue: 1
    MinValue: 1
  - InternalName: dutyCycle
    DisplayName: Duty Cycle
    Description: Percentage of the segment that Color1 is shown.
    DataType: float64
    DefaultValue: 0.5
    MinValue: 0.0
    MaxValue: 1.0
*/
func init() {
	RegisterEffect("blink", NewBlink)
	RegisterEffectMetadata("blink", types.EffectMetadata{
		HumanReadableName: "Blink",
		Description:       "Alternates between two colors based on the global BPM, creating a blinking effect.",
		Tags:              []string{"bpm_sensitive", "color_source", "pattern"},
		Parameters: []types.ParameterMetadata{
			{
				InternalName: "divider",
				DisplayName:  "Divider",
				Description:  "Divides the beat into segments for faster blinking.",
				DataType:     "int",
				DefaultValue: 1,
				MinValue:     1,
			},
			{
				InternalName: "dutyCycle",
				DisplayName:  "Duty Cycle",
				Description:  "Percentage of the segment that Color1 is shown.",
				DataType:     "float64",
				DefaultValue: 0.5,
				MinValue:     0.0,
				MaxValue:     1.0,
			},
		},
	})
}

// Blink alternates between two colors based on the global BPM.
type Blink struct {
	Divider   int
	DutyCycle float64
}

// NewBlink creates a new Blink effect.
func NewBlink(args map[string]interface{}) (types.Effect, error) {
	divider, ok := args["divider"].(float64)
	if !ok {
		return nil, fmt.Errorf("blink effect: missing or invalid 'divider' parameter")
	}
	dutyCycle, ok := args["dutyCycle"].(float64)
	if !ok {
		return nil, fmt.Errorf("blink effect: missing or invalid 'dutyCycle' parameter")
	}
	return &Blink{Divider: int(divider), DutyCycle: dutyCycle}, nil
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
