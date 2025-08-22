package orchestrator

import (
	"fmt"
	"godmx/config"
	"godmx/dmx"
	"godmx/effects"
	"godmx/types"
	sync "sync"
	time "time"
)

// Chain represents a sequence of effects and an output.
type Chain struct {
	ID           string
	Priority     int
	TickRate     int // FPS
	Effects      []types.Effect
	Output       Output
	lamps        []dmx.Lamp // Internal frame buffer for this chain
	orchestrator *Orchestrator // Reference to the parent orchestrator
	config       *config.ChainConfig
	isDirty      bool
	mutex        sync.Mutex
}

// NewChain creates a new Chain instance.
func NewChain(cfg *config.ChainConfig, orch *Orchestrator, output Output) *Chain {
	c := &Chain{
		ID:           cfg.ID,
		Priority:     cfg.Priority,
		TickRate:     cfg.TickRate,
		lamps:        make([]dmx.Lamp, cfg.NumLamps),
		orchestrator: orch,
		config:       cfg,
		Output:       output,
		isDirty:      true, // Start dirty to force initial build
	}
	return c
}

// SetDirty marks the chain as needing a rebuild on the next tick.
func (c *Chain) SetDirty(dirty bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.isDirty = dirty
}

// rebuildEffectsFromConfig clears the current effects and rebuilds them from the config.
func (c *Chain) rebuildEffectsFromConfig() error {
	fmt.Printf("Rebuilding effects for chain '%s'...\n", c.ID)
	c.Effects = []types.Effect{}
	for _, effectConfig := range c.config.Effects {
		if effectConfig.Enabled != nil && !*effectConfig.Enabled {
			continue // Skip disabled effects
		}

		constructor, ok := effects.GetEffectConstructor(effectConfig.Type)
		if !ok {
			return fmt.Errorf("unknown effect type: %s", effectConfig.Type)
		}

		// Make sure to pass the ID to the constructor if it needs it
		args := effectConfig.Args
		if args == nil {
			args = make(map[string]interface{})
		}
		args["id"] = effectConfig.ID

		effect, err := constructor(args)
		if err != nil {
			return fmt.Errorf("error creating effect '%s': %w", effectConfig.Type, err)
		}
		c.Effects = append(c.Effects, effect)
	}
	c.isDirty = false
	return nil
}

// Tick processes the chain's effects and sends data to the output.
func (c *Chain) Tick() error {
	c.mutex.Lock()
	if c.isDirty {
		if err := c.rebuildEffectsFromConfig(); err != nil {
			c.mutex.Unlock()
			return err // Report error but don't stop the chain
		}
	}
	// Create a snapshot of the effects to process for this tick
	effectsSnapshot := make([]types.Effect, len(c.Effects))
	copy(effectsSnapshot, c.Effects)
	c.mutex.Unlock()

	// Update global beat progress before processing effects
	c.orchestrator.UpdateBeatProgress()

	// Process the snapshot of effects
	globals := c.orchestrator.GetGlobals()
	globals.TickRate = c.TickRate
	for _, effect := range effectsSnapshot {
		effect.Process(c.lamps, globals, c.config.Output.ChannelMapping, c.config.Output.NumChannelsPerLamp)
	}

	// Send to output
	return c.Output.Send(c.lamps)
}

// StartLoop starts the chain's independent ticking loop.
func (c *Chain) StartLoop() {
	go func() {
		ticker := time.NewTicker(time.Duration(1000/c.TickRate) * time.Millisecond)
		defer ticker.Stop()
		defer c.Output.Close()

		for range ticker.C {
			if err := c.Tick(); err != nil {
				fmt.Printf("Chain %s error: %v\n", c.ID, err)
			}
		}
	}()
}
