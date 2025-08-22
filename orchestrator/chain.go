package orchestrator

import (
	"fmt"
	"godmx/dmx"
	"time"
)

// Chain represents a sequence of effects and an output.
type Chain struct {
	ID       string
	Priority int
	TickRate int // FPS
	Effects  []Effect
	Output   Output
	lamps    []dmx.Lamp // Internal frame buffer for this chain
	orchestrator *Orchestrator // Reference to the parent orchestrator
}

// NewChain creates a new Chain instance.
func NewChain(id string, priority, tickRate int, output Output, numLamps int, orch *Orchestrator) *Chain {
	return &Chain{
		ID:       id,
		Priority: priority,
		TickRate: tickRate,
		Output:   output,
		lamps:    make([]dmx.Lamp, numLamps),
		orchestrator: orch,
	}
}

// AddEffect adds an effect to the chain.
func (c *Chain) AddEffect(effect Effect) {
	c.Effects = append(c.Effects, effect)
}

// Tick processes the chain's effects and sends data to the output.
func (c *Chain) Tick() error {
	// Process effects
	globals := c.orchestrator.GetGlobals() // Get globals from the orchestrator
	for _, effect := range c.Effects {
		effect.Process(c.lamps, globals) // Pass globals to the effect
	}

	// Debug print the first few lamps' data after effect processing
	fmt.Printf("DEBUG: Chain Tick - Lamps after effect (first 4 lamps):\n")
	for i := 0; i < 4 && i < len(c.lamps); i++ {
		fmt.Printf("  Lamp %d: R=%d, G=%d, B=%d, W=%d\n", i, c.lamps[i].R, c.lamps[i].G, c.lamps[i].B, c.lamps[i].W)
	}

	// Send to output
	return c.Output.Send(c.lamps)
}

// StartLoop starts the chain's independent ticking loop.
func (c *Chain) StartLoop() {
	go func() {
		ticker := time.NewTicker(time.Duration(1000/c.TickRate) * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			err := c.Tick()
			if err != nil {
				// Log error, but don't stop the loop
				// fmt.Printf("Chain %s error: %v\n", c.ID, err)
			}
		}
	}()
}
