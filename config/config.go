package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// ActionConfig represents a single action to be performed by an event.
type ActionConfig struct {
	Type     string                 `json:"type"`
	ChainID  string                 `json:"chain_id,omitempty"`
	EffectID string                 `json:"effect_id,omitempty"`
	Params   map[string]interface{} `json:"params,omitempty"`
}

// MidiTriggerConfig represents a single MIDI message that triggers an event.
type MidiTriggerConfig struct {
	MessageType string `json:"message_type"` // e.g., "cc", "note_on", "note_off"
	Number      int    `json:"number"`       // CC number or note number
	Value       int    `json:"value"`        // CC value or velocity (0-127). Use -1 for any value.
	EventName   string `json:"event_name"`   // The name of the event to trigger
}

// Config represents the overall application configuration.
type Config struct {
	Chains      []ChainConfig             `json:"chains"`
	Globals     GlobalsConfig            `json:"globals"`
	Events      map[string][]ActionConfig `json:"events"`
	MidiTriggers []MidiTriggerConfig      `json:"midi_triggers,omitempty"`
	MidiPortName string                 `json:"midi_port_name,omitempty"`
}

// ChainConfig represents the configuration for a single chain.
type ChainConfig struct {
	ID       string         `json:"id"`
	Priority int            `json:"priority"`
	TickRate int            `json:"tickRate"`
	NumLamps int            `json:"numLamps"`
	Effects  []EffectConfig `json:"effects"`
	Output   OutputConfig   `json:"output"`
}

// EffectConfig represents the configuration for an effect.
type EffectConfig struct {
	ID      string                 `json:"id"`
	Type    string                 `json:"type"`
	Args    map[string]interface{} `json:"args"`
	Enabled *bool                  `json:"enabled,omitempty"`
	Group   string                 `json:"group,omitempty"`
}

// UnmarshalJSON for EffectConfig to default Enabled to true if not present
func (e *EffectConfig) UnmarshalJSON(data []byte) error {
	type Alias EffectConfig
	alias := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	if e.Enabled == nil {
		trueVal := true
		e.Enabled = &trueVal
	}
	return nil
}

// OutputConfig represents the configuration for an output.
type OutputConfig struct {
	Type               string                 `json:"type"`
	Args               map[string]interface{} `json:"args"`
	ChannelMapping     string                 `json:"channelMapping"`
	NumChannelsPerLamp int                    `json:"numChannelsPerLamp"`
	Govee              GoveeOutputConfig      `json:"govee,omitempty"`
}

// GoveeDeviceConfig represents a single Govee device.
type GoveeDeviceConfig struct {
	MACAddress string `json:"mac_address"`
	IPAddress  string `json:"ip_address"`
}

// GoveeOutputConfig represents the configuration for Govee output.
type GoveeOutputConfig struct {
	APIKey  string              `json:"api_key,omitempty"`
	Devices []GoveeDeviceConfig `json:"devices"`
}

// GlobalsConfig represents the global parameters configuration.
type GlobalsConfig struct {
	BPM       float64 `json:"bpm"`
	Color1    string  `json:"color1"`
	Color2    string  `json:"color2"`
	Intensity uint8   `json:"intensity"`
}

// --- Config Manipulation Functions ---

func (c *Config) findChain(chainID string) (*ChainConfig, error) {
	for i := range c.Chains {
		if c.Chains[i].ID == chainID {
			return &c.Chains[i], nil
		}
	}
	return nil, fmt.Errorf("chain with id '%s' not found", chainID)
}

// AddEffectToChain adds a new effect to a chain configuration.
func (c *Config) AddEffectToChain(chainID string, effect EffectConfig) error {
	chain, err := c.findChain(chainID)
	if err != nil {
		return err
	}
	chain.Effects = append(chain.Effects, effect)
	return nil
}

// RemoveEffectFromChain removes an effect from a chain configuration by its ID.
func (c *Config) RemoveEffectFromChain(chainID, effectID string) error {
	chain, err := c.findChain(chainID)
	if err != nil {
		return err
	}
	newEffects := []EffectConfig{}
	found := false
	for _, effect := range chain.Effects {
		if effect.ID == effectID {
			found = true
			continue
		}
		newEffects = append(newEffects, effect)
	}
	if !found {
		return fmt.Errorf("effect with id '%s' not found in chain '%s'", effectID, chainID)
	}
	chain.Effects = newEffects
	return nil
}

// ToggleEffectInChain sets the enabled state of an effect in a chain configuration.
func (c *Config) ToggleEffectInChain(chainID, effectID string, enabled bool) error {
	chain, err := c.findChain(chainID)
	if err != nil {
		return err
	}

	var targetEffect *EffectConfig
	for i := range chain.Effects {
		if chain.Effects[i].ID == effectID {
			targetEffect = &chain.Effects[i]
			break
		}
	}

	if targetEffect == nil {
		return fmt.Errorf("effect with id '%s' not found in chain '%s'", effectID, chainID)
	}

	// If we are enabling this effect, disable all other effects in the same group
	if enabled && targetEffect.Group != "" {
		falseVal := false
		for i := range chain.Effects {
			effect := &chain.Effects[i]
			if effect.ID != effectID && effect.Group == targetEffect.Group && *effect.Enabled {
				effect.Enabled = &falseVal
			}
		}
	}

	// Set the enabled state of the target effect
	targetEffect.Enabled = &enabled

	return nil
}

// SetGlobal sets a global parameter.
func (c *Config) SetGlobal(key string, value interface{}) error {
	switch key {
	case "bpm":
		if bpm, ok := value.(float64); ok {
			c.Globals.BPM = bpm
		} else {
			return fmt.Errorf("invalid type for bpm: expected float64, got %T", value)
		}
	case "intensity":
		if intensity, ok := value.(float64); ok { // JSON unmarshals numbers to float64
			c.Globals.Intensity = uint8(intensity)
		} else {
			return fmt.Errorf("invalid type for intensity: expected float64, got %T", value)
		}
	case "color1":
		if color, ok := value.(string); ok {
			c.Globals.Color1 = color
		} else {
			return fmt.Errorf("invalid type for color1: expected string, got %T", value)
		}
	case "color2":
		if color, ok := value.(string); ok {
			c.Globals.Color2 = color
		} else {
			return fmt.Errorf("invalid type for color2: expected string, got %T", value)
		}
	default:
		return fmt.Errorf("unknown global parameter: %s", key)
	}
	return nil
}

// LoadConfig reads a JSON configuration file and unmarshals it into a Config struct.
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config JSON: %w", err)
	}

	// Create a default config to merge missing values from
	defaultCfg := CreateDefaultConfig()

	// Merge default values into the loaded config
	if mergeConfigs(&cfg, &defaultCfg) {
		log.Println("Updating config file with missing default values.")
		if err := SaveConfig(&cfg, filePath); err != nil {
			log.Printf("Error saving updated config file: %v\n", err)
		}
	}

	return &cfg, nil
}

// SaveConfig marshals the Config struct to JSON and writes it to the specified file path.
func SaveConfig(cfg *Config, filePath string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// mergeConfigs merges default values into the loaded config.
// It returns true if any changes were made.
func mergeConfigs(loaded *Config, defaults *Config) bool {
	changed := false

	// Merge Globals
	if loaded.Globals.BPM == 0 {
		loaded.Globals.BPM = defaults.Globals.BPM
		changed = true
	}
	if loaded.Globals.Color1 == "" {
		loaded.Globals.Color1 = defaults.Globals.Color1
		changed = true
	}
	if loaded.Globals.Color2 == "" {
		loaded.Globals.Color2 = defaults.Globals.Color2
		changed = true
	}
	if loaded.Globals.Intensity == 0 {
		loaded.Globals.Intensity = defaults.Globals.Intensity
		changed = true
	}

	// Merge MidiTriggers
	if len(loaded.MidiTriggers) == 0 && len(defaults.MidiTriggers) > 0 {
		loaded.MidiTriggers = defaults.MidiTriggers
		changed = true
	}

	// Merge MidiPortName
	if loaded.MidiPortName == "" {
		loaded.MidiPortName = defaults.MidiPortName
		changed = true
	}

	// TODO: More sophisticated merging for Chains and Events if needed

	return changed
}

// CreateDefaultConfig creates a default Config struct.
func CreateDefaultConfig() Config {
	trueVal := true
	falseVal := false
	return Config{
		Chains: []ChainConfig{
			{
				ID:       "mainChain",
				Priority: 1,
				TickRate: 40,
				NumLamps: 93,
				Effects: []EffectConfig{
					{
						ID:      "defaultRainbow",
						Type:    "rainbow",
						Args:    make(map[string]interface{}),
						Enabled: &trueVal,
						Group:   "basic_color",
					},
					{
						ID:      "defaultSolidColor",
						Type:    "solidColor",
						Args:    make(map[string]interface{}),
						Enabled: &falseVal,
						Group:   "basic_color",
					},
				},
				Output: OutputConfig{
					Type: "artnet",
					Args: map[string]interface{}{
						"ip": "192.168.125.153",
					},
					ChannelMapping:     "RGB",
					NumChannelsPerLamp: 3,
				},
			},
			{
				ID:       "goveeChain",
				Priority: 3,
				TickRate: 20,
				NumLamps: 16, // Updated to 16
				Effects: []EffectConfig{
					{
						ID:      "goveeRainbow",
						Type:    "rainbow",
						Args:    make(map[string]interface{}),
						Enabled: &trueVal,
						Group:   "govee_color_effects", // Added group
					},
					{ // NEW: goveeSolidColor effect
						ID:      "goveeSolidColor",
						Type:    "solidColor",
						Args:    make(map[string]interface{}),
						Enabled: &falseVal, // Initially disabled
						Group:   "govee_color_effects",
					},
				},
				Output: OutputConfig{
					Type: "govee",
					Govee: GoveeOutputConfig{
						// APIKey: "YOUR_GOVEE_API_KEY", // Removed APIKey
						Devices: []GoveeDeviceConfig{
							{
								MACAddress: "XX:XX:XX:XX:XX:XX", // Placeholder
								IPAddress:  "192.168.1.100",     // Placeholder
							},
						},
					},
					ChannelMapping:     "RGB", // Govee devices typically use RGB
					NumChannelsPerLamp: 3,
				},
			},
			{
				ID:       "ddpChain",
				Priority: 2,
				TickRate: 30,
				NumLamps: 170, // Max lamps in a DDP packet is 170
				Effects: []EffectConfig{
					{
						ID:      "ddpRainbow",
						Type:    "rainbow",
						Args:    make(map[string]interface{}),
						Enabled: &trueVal,
						Group:   "ddp_color_effects",
					},
				},
				Output: OutputConfig{
					Type: "ddp",
					Args: map[string]interface{}{
						"ip": "192.168.1.101", // Placeholder IP for WLED device
					},
					ChannelMapping:     "RGB",
					NumChannelsPerLamp: 3,
				},
			},
		},
		Globals: GlobalsConfig{
			BPM:       120,
			Color1:    "#FF0000",
			Color2:    "#0000FF",
			Intensity: 255,
		},
		Events: map[string][]ActionConfig{
			"strobe_on": {
				{
					Type:    "add_effect",
					ChainID: "mainChain",
					Params: map[string]interface{}{
						"id":      "strobeEffect",
						"type":    "blink",
						"enabled": true,
						"args":    map[string]interface{}{"divider": 4, "dutyCycle": 0.1},
					},
				},
			},
			"strobe_off": {
				{
					Type:     "remove_effect",
					ChainID:  "mainChain",
					EffectID: "strobeEffect",
				},
			},
			"rainbow_on": { // Updated
				{
					Type:     "toggle_effect",
					ChainID:  "mainChain",
					EffectID: "defaultRainbow",
					Params:   map[string]interface{}{"enabled": true},
				},
				{
					Type:     "toggle_effect",
					ChainID:  "goveeChain",
					EffectID: "goveeRainbow",
					Params:   map[string]interface{}{"enabled": true},
				},
			},
			"solid_color_on": { // Updated
				{
					Type:     "toggle_effect",
					ChainID:  "mainChain",
					EffectID: "defaultSolidColor",
					Params:   map[string]interface{}{"enabled": true},
				},
				{
					Type:     "toggle_effect",
					ChainID:  "goveeChain",
					EffectID: "goveeSolidColor", // Changed to toggle_effect
					Params:   map[string]interface{}{"enabled": true},
				},
			},
			"faster_bpm": {
				{
					Type:   "set_global",
					Params: map[string]interface{}{"bpm": 140},
				},
			},
		},
		MidiTriggers: []MidiTriggerConfig{
			{
				MessageType: "cc",
				Number:      1,
				Value:       -1, // Any value
				EventName:   "strobe_on",
			},
			{
				MessageType: "cc",
				Number:      2,
				Value:       -1, // Any value
				EventName:   "strobe_off",
			},
			{
				MessageType: "note_on",
				Number:      60, // Middle C
				Value:       -1, // Any velocity
				EventName:   "rainbow_on",
			},
			{
				MessageType: "note_off",
				Number:      60, // Middle C
				Value:       -1, // Any velocity
				EventName:   "rainbow_off",
			},
		},
		MidiPortName: "Midi Through Port-0",
	}
}
