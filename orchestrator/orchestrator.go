package orchestrator

import (
	"encoding/json"
	"fmt"
	"godmx/config"
	"godmx/dmx"
	"godmx/types"
	"godmx/utils"
	"time"
)

// Output defines the interface for all lighting outputs.
type Output interface {
	Send(lamps []dmx.Lamp) error
	Close() // Add Close method to the interface
}

// Orchestrator manages chains, global parameters, and overall system flow.
type Orchestrator struct {
	chains       []*Chain
	config       *config.Config
	globals      types.OrchestratorGlobals
	lastBeatTime time.Time
}

// NewOrchestrator creates a new Orchestrator instance.
func NewOrchestrator(cfg *config.Config) *Orchestrator {
	return &Orchestrator{
		config: cfg,
		globals: types.OrchestratorGlobals{
			BPM:       120.0,
			Color1:    dmx.Lamp{R: 255, G: 0, B: 0, W: 0},
			Color2:    dmx.Lamp{R: 0, G: 0, B: 255, W: 0},
		},
		lastBeatTime: time.Now(),
	}
}

// AddChain adds a new chain to the orchestrator.
func (o *Orchestrator) AddChain(chain *Chain) {
	o.chains = append(o.chains, chain)
}

func (o *Orchestrator) findChain(chainID string) (*Chain, error) {
	for _, chain := range o.chains {
		if chain.ID == chainID {
			return chain, nil
		}
	}
	return nil, fmt.Errorf("runtime chain with id '%s' not found", chainID)
}

// SetBPM sets the global BPM.
func (o *Orchestrator) SetBPM(bpm float64) {
	o.globals.BPM = bpm
}

// SetColor1 sets the global Color1.
func (o *Orchestrator) SetColor1(color dmx.Lamp) {
	o.globals.Color1 = color
}

// SetColor2 sets the global Color2.
func (o *Orchestrator) SetColor2(color dmx.Lamp) {
	o.globals.Color2 = color
}

// GetGlobals returns a pointer to the orchestrator's global parameters.
func (o *Orchestrator) GetGlobals() *types.OrchestratorGlobals {
	return &o.globals
}

// UpdateBeatProgress calculates and updates the global beat progress.
func (o *Orchestrator) UpdateBeatProgress() {
	elapsed := time.Since(o.lastBeatTime)
	beatDuration := time.Duration((60.0 / o.globals.BPM) * float64(time.Second))
	o.globals.BeatProgress = float64(elapsed) / float64(beatDuration)

	if o.globals.BeatProgress >= 1.0 {
		o.lastBeatTime = time.Now()
		o.globals.BeatProgress = 0.0
	}
}

// TriggerEvent finds an event by name in the config and executes its actions.
func (o *Orchestrator) TriggerEvent(eventName string) {
	actions, ok := o.config.Actions[eventName]
	if !ok {
		fmt.Printf("Event '%s' not found.\n", eventName)
		return
	}

	fmt.Printf("Triggering event '%s'\n", eventName)
	for _, action := range actions {
		if err := o.executeAction(action); err != nil {
			fmt.Printf("  - Error executing action '%s': %v\n", action.Type, err)
		}
	}
}

func mapToEffectConfig(params map[string]interface{}) (config.EffectConfig, error) {
	var effectConfig config.EffectConfig
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return effectConfig, fmt.Errorf("failed to marshal effect params: %w", err)
	}
	if err := json.Unmarshal(jsonBytes, &effectConfig); err != nil {
		return effectConfig, fmt.Errorf("failed to unmarshal effect params: %w", err)
	}
	return effectConfig, nil
}

// executeAction executes a single action from an event.
func (o *Orchestrator) executeAction(action config.ActionConfig) error {
	fmt.Printf("  - Executing action: %s\n", action.Type)
	var err error

	switch action.Type {
	case "add_effect":
		effectConfig, err := mapToEffectConfig(action.Params)
		if err != nil {
			return err
		}
		err = o.config.AddEffectToChain(action.ChainID, effectConfig)
		if err == nil {
			chain, findErr := o.findChain(action.ChainID)
			if findErr == nil {
				chain.SetDirty(true)
			}
		}
	case "remove_effect":
		err = o.config.RemoveEffectFromChain(action.ChainID, action.EffectID)
		if err == nil {
			chain, findErr := o.findChain(action.ChainID)
			if findErr == nil {
				chain.SetDirty(true)
			}
		}
	case "toggle_effect":
		enabled, ok := action.Params["enabled"].(bool)
		if !ok {
			return fmt.Errorf("missing or invalid 'enabled' param for toggle_effect")
		}
		err = o.config.ToggleEffectInChain(action.ChainID, action.EffectID, enabled)
		if err == nil {
			chain, findErr := o.findChain(action.ChainID)
			if findErr == nil {
				chain.SetDirty(true)
			}
		}
	case "set_global":
		for key, val := range action.Params {
			err = o.config.SetGlobal(key, val)
			// Also update the running orchestrator's globals
			if err == nil {
				o.SetBPM(o.config.Globals.BPM)
				color1, err1 := utils.ParseHexColor(o.config.Globals.Color1)
				if err1 == nil {
					o.SetColor1(color1)
				}
				color2, err2 := utils.ParseHexColor(o.config.Globals.Color2)
				if err2 == nil {
					o.SetColor2(color2)
				}
			}
		}
	default:
		err = fmt.Errorf("unknown action type: %s", action.Type)
	}

	return err
}
