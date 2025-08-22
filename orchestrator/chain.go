package orchestrator

import (
	
	"godmx/dmx"
	"time"
)

// Chain represents a sequence of effects and an output.
type Chain struct {
	ID                 string
	Priority           int
	TickRate           int // FPS
	Effects            []Effect
	Output             Output
	lamps              []dmx.Lamp // Internal frame buffer for this chain
	orchestrator       *Orchestrator // Reference to the parent orchestrator
	channelMapping     string
	numChannelsPerLamp int
}

// NewChain creates a new Chain instance.
func NewChain(id string, priority, tickRate int, output Output, numLamps int, orch *Orchestrator, channelMapping string, numChannelsPerLamp int) *Chain {
	return &Chain{
		ID:                 id,
		Priority:           priority,
		TickRate:           tickRate,
		Output:             output,
		lamps:              make([]dmx.Lamp, numLamps),
		orchestrator:       orch,
		channelMapping:     channelMapping,
		numChannelsPerLamp: numChannelsPerLamp,
	}
}

// AddEffect adds an effect to the chain.
func (c *Chain) AddEffect(effect Effect) {
	c.Effects = append(c.Effects, effect)
}

// Tick processes the chain's effects and sends data to the output.
func (c *Chain) Tick() error {
	// Update global beat progress before processing effects
	c.orchestrator.UpdateBeatProgress()

	// Process effects
	globals := c.orchestrator.GetGlobals()
	globals.TickRate = c.TickRate // Get globals from the orchestrator
	for _, effect := range c.Effects {
		effect.Process(c.lamps, globals, c.channelMapping, c.numChannelsPerLamp) // Pass globals and channel info to the effect
	}

	

	// Send to output
	return c.Output.Send(c.lamps)
}

// StartLoop starts the chain's independent ticking loop.
func (c *Chain) StartLoop() {
	go func() {
		ticker := time.NewTicker(time.Duration(1000/c.TickRate) * time.Millisecond)
		defer ticker.Stop()
		defer c.Output.Close() // Close the output when the loop exits

		for range ticker.C {
			err := c.Tick()
			if err != nil {
				// Log error, but don't stop the loop
				// fmt.Printf("Chain %s error: %v\n", c.ID, err)
			}
		}
	}()
}

