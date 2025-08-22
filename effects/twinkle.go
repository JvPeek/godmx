package effects

import (
	"godmx/dmx"
	"godmx/orchestrator"
	"math/rand"
	"time"
)

func init() {
	RegisterEffect("twinkle", func(args map[string]interface{}) (orchestrator.Effect, map[string]interface{}, error) {
		effect, modifiedArgs, err := NewTwinkle(args)
		return effect, modifiedArgs, err
	})
}

// Twinkle randomly turns a percentage of lamps to white.
type Twinkle struct {
	Percentage float64
	source     rand.Source
	generator  *rand.Rand
}

// NewTwinkle creates a new Twinkle effect.
func NewTwinkle(args map[string]interface{}) (*Twinkle, map[string]interface{}, error) {
	if args == nil {
		args = make(map[string]interface{})
	}
	percentage, ok := args["percentage"].(float64)
	if !ok {
		percentage = 0.1 // Default to 10% twinkle
		args["percentage"] = percentage
	}

	src := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(src)

	return &Twinkle{
		Percentage: percentage,
		source:     src,
		generator:  gen,
	}, args, nil
}

// Process applies the twinkle effect to the lamps.
func (t *Twinkle) Process(lamps []dmx.Lamp, globals *orchestrator.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
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
}
