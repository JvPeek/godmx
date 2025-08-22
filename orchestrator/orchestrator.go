package orchestrator

import (
	"fmt"
	"godmx/dmx"
	"time"
)

// OrchestratorGlobals holds the global parameters managed by the orchestrator.
type OrchestratorGlobals struct {
	BPM       float64
	Color1    dmx.Lamp
	Color2    dmx.Lamp
	Intensity uint8
	TotalLamps int // Total number of lamps across all chains
	TickRate int // Current chain's tick rate (FPS)
	BeatProgress float64 // New: Current position within the beat (0.0 to 1.0)
}

// Effect defines the interface for all lighting effects.
type Effect interface {
	Process(lamps []dmx.Lamp, globals *OrchestratorGlobals, channelMapping string, numChannelsPerLamp int)
}

// Output defines the interface for all lighting outputs.
type Output interface {
	Send(lamps []dmx.Lamp) error
	Close() // Add Close method to the interface
}

// Orchestrator manages chains, global parameters, and overall system flow.
type Orchestrator struct {
	chains []*Chain
	// Global Parameters
	globals OrchestratorGlobals // Embed the globals struct
	lastBeatTime time.Time // New: Time when the last beat started
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
		lastBeatTime: time.Now(), // Initialize lastBeatTime
	}
}

// AddChain adds a new chain to the orchestrator.
func (o *Orchestrator) AddChain(chain *Chain) {
	o.chains = append(o.chains, chain)
}

// SetBPM sets the global BPM.
func (o *Orchestrator) SetBPM(bpm float64) {
	o.globals.BPM = bpm
	fmt.Printf("SetBPM: Setting BPM to %.2f\n", bpm) // Add this line
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

// UpdateBeatProgress calculates and updates the global beat progress.
func (o *Orchestrator) UpdateBeatProgress() {
	
	// Calculate time since last beat
	elapsed := time.Since(o.lastBeatTime)

	

	// Calculate duration of one beat
	beatDuration := time.Duration((60.0 / o.globals.BPM) * float64(time.Second))
	

	// Calculate beat progress (0.0 to 1.0)
	o.globals.BeatProgress = float64(elapsed) / float64(beatDuration)
	

	// If we've passed a full beat, reset lastBeatTime
	if o.globals.BeatProgress >= 1.0 {
		o.lastBeatTime = time.Now()
		o.globals.BeatProgress = 0.0 // Reset for the new beat
		fmt.Println("Beat: Reset!")
	}
}

// Run starts the orchestrator's main loop.
func (o *Orchestrator) Run() {
	// This will be the main loop that ticks chains
	// We'll implement the actual ticking logic later.
}
