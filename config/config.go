package config

import (
	"encoding/json"
	"fmt"
	"godmx/effects"
	"log"
	"os"
)

// ActionConfig represents a single action to be performed by an event.
type ActionConfig struct {
	Type     string                 	`json:"type"`
	ChainID  string                 	`json:"chain_id,omitempty"`
	EffectID string                 	`json:"effect_id,omitempty"`
	Params   map[string]interface{} 	`json:"params,omitempty"`
}

// MidiTriggerConfig represents a single MIDI message that triggers an event.
type MidiTriggerConfig struct {
	MessageType string 	`json:"message_type"` // e.g., "cc", "note_on", "note_off"
	Number      int    	`json:"number"`       // CC number or note number
	Value       int    	`json:"value"`        // CC value or velocity (0-127). Use -1 for any value.
	EventName   string 	`json:"event_name"`   // The name of the event to trigger
}

// Config represents the overall application configuration.
type Config struct {
	Globals      GlobalsConfig            	`json:"globals"`
	Chains       []ChainConfig            	`json:"chains"`
	Actions      map[string][]ActionConfig 	`json:"actions"` // Renamed from Events
	Triggers     []MidiTriggerConfig      	`json:"triggers"` // Renamed from MidiTriggers
	MidiPortName string                 	`json:"midi_port_name,omitempty"`
}

// ChainConfig represents the configuration for a single chain.
type ChainConfig struct {
	ID       string         	`json:"id"`
	Priority int            	`json:"priority"`
	TickRate int            	`json:"tickRate"`
	NumLamps int            	`json:"numLamps"`
	Effects  []EffectConfig 	`json:"effects"`
	Output   OutputConfig   	`json:"output"`
}

// EffectConfig represents the configuration for an effect.
type EffectConfig struct {
	ID      string                 	`json:"id"`
	Type    string                 	`json:"type"`
	Args    map[string]interface{} 	`json:"args"`
	Enabled *bool                  	`json:"enabled,omitempty"`
	Group   string                 	`json:"group,omitempty"`
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
	Type               string                 	`json:"type"`
	Args               map[string]interface{} 	`json:"args"`
	ChannelMapping     string                 	`json:"channelMapping"`
	NumChannelsPerLamp int                    	`json:"numChannelsPerLamp"`
	Govee              GoveeOutputConfig      	`json:"govee,omitempty"`
}

// GoveeDeviceConfig represents a single Govee device.
type GoveeDeviceConfig struct {
	MACAddress string 	`json:"mac_address"`
	IPAddress  string 	`json:"ip_address"`
}

// GoveeOutputConfig represents the configuration for Govee output.
type GoveeOutputConfig struct {
	APIKey  string              	`json:"api_key,omitempty"`
	Devices []GoveeDeviceConfig 	`json:"devices"`
}

// GlobalsConfig represents the global parameters configuration.
type GlobalsConfig struct {
	BPM    float64 	`json:"bpm"`
	Color1 string  	`json:"color1"`
	Color2 string  	`json:"color2"`
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
		if os.IsNotExist(err) {
			log.Printf("Config file not found at %s. Creating default config.\n", filePath)
			cfg := CreateDefaultConfig()
			if err := SaveConfig(&cfg, filePath); err != nil {
				return nil, fmt.Errorf("failed to save default config file: %w", err)
			}
			return &cfg, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config JSON: %w", err)
	}

	// Create a default config to merge missing values from
	defaultCfg := CreateDefaultConfig()

	// Merge default values into the loaded config
	configModified := mergeConfigs(&cfg, &defaultCfg)

	// Augment effect arguments with default values from metadata
	for i := range cfg.Chains {
		chain := &cfg.Chains[i]
		for j := range chain.Effects {
			effect := &chain.Effects[j]
			metadata, ok := effects.GetEffectMetadata(effect.Type)
			if !ok {
				log.Printf("Warning: Metadata not found for effect type '%s'. Skipping default arg augmentation.\n", effect.Type)
				continue
			}

			if effect.Args == nil {
				effect.Args = make(map[string]interface{})
			}

			for _, param := range metadata.Parameters {
				if _, exists := effect.Args[param.InternalName]; !exists {
					effect.Args[param.InternalName] = param.DefaultValue
					configModified = true
				}
			}
		}
	}

	// Save config if any changes were made (either by mergeConfigs or effect arg augmentation)
	if configModified {
		log.Println("Updating config file with missing default values or augmented effect arguments.")
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

	// Ensure Chains, Actions, and Triggers are initialized if they are nil after unmarshaling
	if loaded.Chains == nil {
		loaded.Chains = []ChainConfig{}
		changed = true
	}
	if loaded.Actions == nil {
		loaded.Actions = make(map[string][]ActionConfig)
		changed = true
	}
	if loaded.Triggers == nil {
		loaded.Triggers = []MidiTriggerConfig{}
		changed = true
	}

	// Merge MidiPortName
	if loaded.MidiPortName == "" {
		loaded.MidiPortName = defaults.MidiPortName
		changed = true
	}

	return changed
}

// CreateDefaultConfig creates a default Config struct.
func CreateDefaultConfig() Config {
	return Config{
		Globals: GlobalsConfig{
			BPM:    174,
			Color1: "FFA000",
			Color2: "000000",
		},
		Chains:       []ChainConfig{},
		Actions:      make(map[string][]ActionConfig),
		Triggers:     []MidiTriggerConfig{},
		MidiPortName: "", // Default to empty, will be merged if default has one
	}
}