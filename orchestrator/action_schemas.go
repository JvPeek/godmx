package orchestrator

// ActionParameter describes a single parameter for an action.
type ActionParameter struct {
	InternalName string      `json:"internal_name"`
	DisplayName  string      `json:"display_name"`
	Description  string      `json:"description"`
	DataType     string      `json:"data_type"` // e.g., "string", "float64", "int", "bool"
	DefaultValue interface{} `json:"default_value,omitempty"`
	Options      []string    `json:"options,omitempty"` // For enum-like parameters
}

// ActionSchema holds comprehensive metadata about an action type.
type ActionSchema struct {
	HumanReadableName string            `json:"human_readable_name"`
	Description       string            `json:"description"`
	Parameters        []ActionParameter `json:"parameters"`
}

// ActionSchemas is a map of action type names to their schemas.
var ActionSchemas = map[string]ActionSchema{
	"add_effect": {
		HumanReadableName: "Add Effect",
		Description:       "Adds a new effect to a specified chain.",
		Parameters: []ActionParameter{
			{InternalName: "chain_id", DisplayName: "Chain ID", Description: "The ID of the chain to add the effect to.", DataType: "string"},
			{InternalName: "params", DisplayName: "Effect Parameters", Description: "The full configuration of the effect to add.", DataType: "object"}, // This will be complex, might need a nested schema or special handling
		},
	},
	"remove_effect": {
		HumanReadableName: "Remove Effect",
		Description:       "Removes an effect from a specified chain.",
		Parameters: []ActionParameter{
			{InternalName: "chain_id", DisplayName: "Chain ID", Description: "The ID of the chain to remove the effect from.", DataType: "string"},
			{InternalName: "effect_id", DisplayName: "Effect ID", Description: "The ID of the effect to remove.", DataType: "string"},
		},
	},
	"toggle_effect": {
		HumanReadableName: "Toggle Effect",
		Description:       "Toggles the enabled state of an existing effect in a chain.",
		Parameters: []ActionParameter{
			{InternalName: "chain_id", DisplayName: "Chain ID", Description: "The ID of the chain containing the effect.", DataType: "string"},
			{InternalName: "effect_id", DisplayName: "Effect ID", Description: "The ID of the effect to toggle.", DataType: "string"},
			{InternalName: "enabled", DisplayName: "Enabled", Description: "Whether the effect should be enabled or disabled.", DataType: "bool", DefaultValue: true},
		},
	},
	"set_global": {
		HumanReadableName: "Set Global Parameter",
		Description:       "Sets a global parameter (like BPM, Color1, Color2, Intensity).",
		Parameters: []ActionParameter{
			{InternalName: "bpm", DisplayName: "BPM", Description: "Global Beats Per Minute.", DataType: "float64", DefaultValue: 120.0},
			{InternalName: "intensity", DisplayName: "Intensity", Description: "Global intensity (0-255).", DataType: "int", DefaultValue: 255},
			{InternalName: "color1", DisplayName: "Color 1", Description: "Global Color 1 (hex string, e.g., #FF0000).", DataType: "string", DefaultValue: "#FF0000"},
			{InternalName: "color2", DisplayName: "Color 2", Description: "Global Color 2 (hex string, e.g., #0000FF).", DataType: "string", DefaultValue: "#0000FF"},
		},
	},
}
