package types

import "godmx/dmx"

// OrchestratorGlobals holds the global parameters managed by the orchestrator.
type OrchestratorGlobals struct {
	BPM          float64
	Color1       dmx.Lamp
	Color2       dmx.Lamp
	TotalLamps   int
	TickRate     int
	BeatProgress float64
}

// Effect defines the interface for all lighting effects.
type Effect interface {
	Process(lamps []dmx.Lamp, globals *OrchestratorGlobals, channelMapping string, numChannelsPerLamp int)
}

// ParameterMetadata describes a single parameter for an effect.
type ParameterMetadata struct {
	InternalName string      `json:"internal_name"`
	DisplayName  string      `json:"display_name"`
	Description  string      `json:"description"`
	DataType     string      `json:"data_type"` // e.g., "float64", "int", "string", "bool"
	DefaultValue interface{} `json:"default_value"`
	MinValue     interface{} `json:"min_value,omitempty"`
	MaxValue     interface{} `json:"max_value,omitempty"`
}

// EffectMetadata holds comprehensive metadata about an effect.
type EffectMetadata struct {
	HumanReadableName string                `json:"human_readable_name"`
	Description       string                `json:"description"`
	Tags              []string              `json:"tags"`
	Parameters        []ParameterMetadata `json:"parameters"`
}
