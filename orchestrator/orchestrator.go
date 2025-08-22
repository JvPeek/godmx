package orchestrator

import (
	"godmx/dmx"
)

// OrchestratorGlobals holds the global parameters managed by the orchestrator.
type OrchestratorGlobals struct {
	BPM       float64
	Color1    dmx.Lamp
	Color2    dmx.Lamp
	Intensity uint8
}

// Effect defines the interface for all lighting effects.
type Effect interface {
	Process(lamps []dmx.Lamp, globals *OrchestratorGlobals)
}

// Output defines the interface for all lighting outputs.
type Output interface {
	Send(lamps []dmx.Lamp) error
}

// Orchestrator manages chains, global parameters, and overall system flow.
type Orchestrator struct {
	chains []*Chain
	// Global Parameters
	globals OrchestratorGlobals // Embed the globals struct
}

// NewOrchestrator creates a new Orchestrator instance.
func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		globals: OrchestratorGlobals{
			BPM:       120.0, // Default BPM
			Color1:    dmx.Lamp{R: 255, G: 0, B: 0, W: 0}, // Default Red
			Color2:    dmx.Lamp{R: 0, G: 0, B: 255, W: 0}, // Default Blue
			Intensity: 255, // Default full intensity
		},
	}
}

// AddChain adds a new chain to the orchestrator.
func (o *Orchestrator) AddChain(chain *Chain) {
	o.chains = append(o.chains, chain)
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

// SetIntensity sets the global Intensity.
func (o *Orchestrator) SetIntensity(intensity uint8) {
	o.globals.Intensity = intensity
}

// GetGlobals returns a pointer to the orchestrator's global parameters.
func (o *Orchestrator) GetGlobals() *OrchestratorGlobals {
	return &o.globals
}

// Run starts the orchestrator's main loop.
func (o *Orchestrator) Run() {
	// This will be the main loop that ticks chains
	// We'll implement the actual ticking logic later.
}
