package effects

import (
	"fmt"
	"godmx/dmx"
	"godmx/orchestrator"
	"math"
)

func init() {
	RegisterEffect("shift", func(args map[string]interface{}) (orchestrator.Effect, map[string]interface{}, error) {
		effect, modifiedArgs, err := NewShift(args)
		return effect, modifiedArgs, err
	})
}

// Shift effect shifts the DMX data left or right.
type Shift struct {
	Direction string  // "left" or "right"
	step      float64 // Current fractional shift step
}

// NewShift creates a new Shift effect.
func NewShift(args map[string]interface{}) (*Shift, map[string]interface{}, error) {
	if args == nil {
		args = make(map[string]interface{})
	}
	direction, ok := args["direction"].(string)
	if !ok {
		direction = "left" // Default direction
		args["direction"] = direction
	}

	if direction != "left" && direction != "right" {
		return nil, args, fmt.Errorf("invalid direction for shift effect: %s. Must be 'left' or 'right'", direction)
	}

	return &Shift{Direction: direction},
		args, nil
}

// Process applies the shift effect to the lamps.
func (s *Shift) Process(lamps []dmx.Lamp, globals *orchestrator.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	numLamps := float64(len(lamps))
	const fixedTickRate = 40.0 // Assuming 40 FPS as per current configs

	// Calculate how much the counter should advance per tick to complete one full shift per beat
	// shiftPerTick = (numLamps / TicksPerBeat)
	// TicksPerBeat = (fixedTickRate * 60.0) / globals.BPM
	shiftPerTick := (numLamps * globals.BPM) / (fixedTickRate * 60.0)

	s.step += shiftPerTick

	// Ensure step wraps around the number of lamps
	s.step = math.Mod(s.step, numLamps)

	shiftedLamps := make([]dmx.Lamp, int(numLamps))

	for i := 0; i < int(numLamps); i++ {
		var sourceIndex int
		if s.Direction == "left" {
			sourceIndex = int(math.Round(float64(i) + s.step)) % int(numLamps)
		} else { // right
			sourceIndex = int(math.Round(float64(i) - s.step))
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
