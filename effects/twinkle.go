package effects

import (
	"fmt"
	"godmx/dmx"
	"godmx/orchestrator"
	"math/rand"
	"time"
)

func init() {
	RegisterEffect("twinkle", func(args map[string]interface{}) (orchestrator.Effect, error) {
		return NewTwinkle(args)
	})
	RegisterEffectParameters("twinkle", map[string]interface{}{"percentage": 0.1})
}

// Twinkle randomly turns a percentage of lamps to white.
type Twinkle struct {
	Percentage float64
	source     rand.Source
	generator  *rand.Rand
	lastBeatTriggered bool // New: Flag to track if twinkle was triggered on the current beat
}

// NewTwinkle creates a new Twinkle effect.
func NewTwinkle(args map[string]interface{}) (*Twinkle, error) {
	percentage := args["percentage"].(float64)

	src := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(src)

	return &Twinkle{
		Percentage: percentage,
		source:     src,
		generator:  gen,
		lastBeatTriggered: false,
	}, nil
}

// Process applies the twinkle effect to the lamps.
func (t *Twinkle) Process(lamps []dmx.Lamp, globals *orchestrator.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	fmt.Printf("Twinkle: BeatProgress=%.2f, lastBeatTriggered=%t\n", globals.BeatProgress, t.lastBeatTriggered)

	// Trigger twinkle only once per beat, when BeatProgress crosses a threshold (e.g., 0.0)
	// and it hasn't been triggered yet for this beat.
	if globals.BeatProgress < 0.1 && !t.lastBeatTriggered { // Trigger at the beginning of the beat
		fmt.Println("Twinkle: Triggered!")
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

