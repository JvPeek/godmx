package effects

import (
	"fmt"
	"godmx/dmx"
	"godmx/types"
	"math/rand"
	"time"
)

/*
Effect Name: Twinkle
Description: Randomly turns a percentage of lamps to white at the beginning of each beat, creating a twinkling effect.
Tags: [bpm_sensitive, color_source, random, pattern]
Parameters:
  - InternalName: percentage
    DisplayName: Percentage
    Description: The percentage of lamps to twinkle (0.0 - 1.0).
    DataType: float64
    DefaultValue: 0.1
    MinValue: 0.0
    MaxValue: 1.0
*/
func init() {
	RegisterEffect("twinkle", NewTwinkle)
	RegisterEffectMetadata("twinkle", types.EffectMetadata{
		HumanReadableName: "Twinkle",
		Description:       "Randomly turns a percentage of lamps to white at the beginning of each beat, creating a twinkling effect.",
		Tags:              []string{"bpm_sensitive", "color_source", "random", "pattern"},
		Parameters: []types.ParameterMetadata{
			{
				InternalName: "percentage",
				DisplayName:  "Percentage",
				Description:  "The percentage of lamps to twinkle (0.0 - 1.0).",
				DataType:     "float64",
				DefaultValue: 0.1,
				MinValue:     0.0,
				MaxValue:     1.0,
			},
		},
	})
}

// Twinkle randomly turns a percentage of lamps to white.
type Twinkle struct {
	Percentage        float64
	source            rand.Source
	generator         *rand.Rand
	lastBeatTriggered bool // New: Flag to track if twinkle was triggered on the current beat
}

// NewTwinkle creates a new Twinkle effect.
func NewTwinkle(args map[string]interface{}) (types.Effect, error) {
	percentage, ok := args["percentage"].(float64)
	if !ok {
		return nil, fmt.Errorf("twinkle effect: missing or invalid 'percentage' parameter")
	}

	src := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(src)

	return &Twinkle{
		Percentage:        percentage,
		source:            src,
		generator:         gen,
		lastBeatTriggered: false,
	}, nil
}

// Process applies the twinkle effect to the lamps.
func (t *Twinkle) Process(lamps []dmx.Lamp, globals *types.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	// Trigger twinkle only once per beat, when BeatProgress crosses a threshold (e.g., 0.0)
	// and it hasn't been triggered yet for this beat.
	if globals.BeatProgress < 0.1 && !t.lastBeatTriggered { // Trigger at the beginning of the beat
		numToTwinkle := int(float64(len(lamps)) * t.Percentage)

		// Create a permutation of lamp indices and pick the first `numToTwinkle`.
		// This ensures we don't pick the same lamp twice in one frame.
		indices := t.generator.Perm(len(lamps))

		for i := 0; i < numToTwinkle; i++ {
			lampi := indices[i]
			if numChannelsPerLamp == 3 && channelMapping == "RGB" {
				lamps[lampi] = dmx.Lamp{R: 255, G: 255, B: 255, W: 0} // Set RGB to white, W to 0
			} else {
				lamps[lampi] = dmx.Lamp{R: 255, G: 255, B: 255, W: 255} // Default to RGBW white
			}
		}
		t.lastBeatTriggered = true // Mark as triggered for this beat
	} else if globals.BeatProgress >= 0.1 {
		// Reset the flag once we've passed the trigger threshold for the current beat
		t.lastBeatTriggered = false
	}
}
