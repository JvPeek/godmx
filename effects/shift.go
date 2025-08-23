package effects

import (
	"fmt"
	"godmx/dmx"
	"godmx/types"
)

func init() {
	RegisterEffect("shift", func(args map[string]interface{}) (types.Effect, error) {
		return NewShift(args)
	})
	RegisterEffectMetadata("shift", types.EffectMetadata{
		HumanReadableName: "Shift",
		Description:       "Shifts the DMX data (colors) across the lamps either left or right, synchronized with the BPM.",
		Tags:              []string{"bpm_sensitive", "transparent", "transform", "pattern"},
		Parameters: []types.ParameterMetadata{
			{
				InternalName: "direction",
				DisplayName:  "Direction",
				Description:  "The direction to shift the lamps ('left' or 'right').",
				DataType:     "string",
				DefaultValue: "left",
			},
			{
				InternalName: "speed",
				DisplayName:  "Speed",
				Description:  "The speed of the shift, from 0 to 1 (1 being 1 shift per beat).",
				DataType:     "float64",
				DefaultValue: 1.0,
			},
		},
	})
}

// Shift effect shifts the DMX data left or right.
type Shift struct {
	Direction string  // "left" or "right"
	Speed     float64 // 0 to 1, 1 being 1 shift per beat
}

// NewShift creates a new Shift effect.
func NewShift(args map[string]interface{}) (*Shift, error) {
	direction, ok := args["direction"].(string)
	if !ok || direction == "" {
		direction = "left" // Default to "left" if not provided or empty
	}

	if direction != "left" && direction != "right" {
		return nil, fmt.Errorf("invalid direction for shift effect: %s. Must be 'left' or 'right'", direction)
	}

	speed := 1.0 // Default speed
	if s, ok := args["speed"].(float64); ok {
		if s >= 0 && s <= 1 {
			speed = s
		} else {
			return nil, fmt.Errorf("invalid speed for shift effect: %f. Must be between 0 and 1", s)
		}
	} else if s, ok := args["speed"].(int); ok {
		if float64(s) >= 0 && float64(s) <= 1 {
			speed = float64(s)
		} else {
			return nil, fmt.Errorf("invalid speed for shift effect: %d. Must be between 0 and 1", s)
		}
	}

	return &Shift{Direction: direction, Speed: speed},
		nil
}

// Process applies the shift effect to the lamps.
func (s *Shift) Process(lamps []dmx.Lamp, globals *types.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	numLamps := float64(len(lamps))

	// Calculate the total shift amount based on beat progress and speed
	// This value accumulates over time, representing the total number of lamps shifted
	// since the beginning of the current beat cycle, scaled by speed.
	shiftAmount := globals.BeatProgress * s.Speed * numLamps

	shiftedLamps := make([]dmx.Lamp, int(numLamps))

	for i := 0; i < int(numLamps); i++ {
		var sourceIndex int
		if s.Direction == "left" {
			sourceIndex = int(float64(i) + shiftAmount)
		} else { // right
			sourceIndex = int(float64(i) - shiftAmount)
		}

		// Ensure sourceIndex wraps around the number of lamps
		sourceIndex = sourceIndex % int(numLamps)
		if sourceIndex < 0 {
			sourceIndex += int(numLamps)
		}
		shiftedLamps[i] = lamps[sourceIndex]
	}

	// Copy shifted lamps back to original lamps array
	for i := range lamps {
		lamps[i] = shiftedLamps[i]
	}
}