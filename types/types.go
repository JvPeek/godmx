package types

import "godmx/dmx"

// OrchestratorGlobals holds the global parameters managed by the orchestrator.
type OrchestratorGlobals struct {
	BPM          float64
	Color1       dmx.Lamp
	Color2       dmx.Lamp
	Intensity    uint8
	TotalLamps   int
	TickRate     int
	BeatProgress float64
}

// Effect defines the interface for all lighting effects.
type Effect interface {
	Process(lamps []dmx.Lamp, globals *OrchestratorGlobals, channelMapping string, numChannelsPerLamp int)
}
