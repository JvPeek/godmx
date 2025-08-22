package effects

import (
	"fmt"
	"godmx/dmx"
	"godmx/types"
	"math"
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
		},
	})
}

// Shift effect shifts the DMX data left or right.
type Shift struct {
	Direction string  // "left" or "right"
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

	return &Shift{Direction: direction},
		nil
}

// Process applies the shift effect to the lamps.
func (s *Shift) Process(lamps []dmx.Lamp, globals *types.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	numLamps := float64(len(lamps))

	step := globals.BeatProgress * numLamps // Shift one full length of lamps per beat

	// Ensure step wraps around the number of lamps
	step = math.Mod(step, numLamps)

	shiftedLamps := make([]dmx.Lamp, int(numLamps))

	for i := 0; i < int(numLamps); i++ {
		var sourceIndex int
		if s.Direction == "left" {
			sourceIndex = int(math.Round(float64(i) + step)) % int(numLamps)
		} else { // right
			sourceIndex = int(math.Round(float64(i) - step))
			// Handle negative indices for right shift
			if sourceIndex < 0 {
				sourceIndex = int(numLamps) + sourceIndex % int(numLamps)
			}
			sourceIndex = sourceIndex % int(numLamps)
		}
		shiftedLamps[i] = lamps[sourceIndex]
	}

	// Copy shifted lamps back to original lamps array
	for i := range lamps {
		lamps[i] = shiftedLamps[i]
	}
}