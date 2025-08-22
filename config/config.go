package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// ActionConfig represents a single action to be performed by an event.
type ActionConfig struct {
	Type     string                 `json:"type"`
	ChainID  string                 `json:"chain_id,omitempty"`
	EffectID string                 `json:"effect_id,omitempty"`
	Params   map[string]interface{} `json:"params,omitempty"`
}

// Config represents the overall application configuration.
type Config struct {
	Chains  []ChainConfig             `json:"chains"`
	Globals GlobalsConfig            `json:"globals"`
	Events  map[string][]ActionConfig `json:"events"`
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
	for i := range chain.Effects {
		if chain.Effects[i].ID == effectID {
			chain.Effects[i].Enabled = &enabled
			return nil
		}
	}
	return fmt.Errorf("effect with id '%s' not found in chain '%s'", effectID, chainID)
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
	// Add cases for color1 and color2 if needed, requires parsing hex string
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

// CreateDefaultConfig creates a default config.json file.
func CreateDefaultConfig(filePath string) error {
	trueVal := true
	defaultConfig := Config{
		Chains: []ChainConfig{
			{
				ID:       "mainChain",
				Priority: 1,
				TickRate: 40,
				NumLamps: 50,
				Effects: []EffectConfig{
					{
						ID:      "defaultRainbow",
						Type:    "rainbow",
						Args:    make(map[string]interface{}),
						Enabled: &trueVal,
					},
				},
				Output: OutputConfig{
					Type: "artnet",
					Args: map[string]interface{}{
						"ip": "127.0.0.1",
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
			"rainbow_off": {
				{
					Type:     "toggle_effect",
					ChainID:  "mainChain",
					EffectID: "defaultRainbow",
					Params:   map[string]interface{}{"enabled": false},
				},
			},
			"rainbow_on": {
				{
					Type:     "toggle_effect",
					ChainID:  "mainChain",
					EffectID: "defaultRainbow",
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
	}

	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write default config file: %w", err)
	}

	return nil
}