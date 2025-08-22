package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the overall application configuration.
type Config struct {
	Chains  []ChainConfig  `json:"chains"`
	Globals GlobalsConfig `json:"globals"`
}

// ChainConfig represents the configuration for a single chain.
type ChainConfig struct {
	ID       string        `json:"id"`
	Priority int           `json:"priority"`
	TickRate int           `json:"tickRate"`
	NumLamps int           `json:"numLamps"`
	Effects  []EffectConfig `json:"effects"`
	Output   OutputConfig  `json:"output"`
}

// EffectConfig represents the configuration for an effect.
type EffectConfig struct {
	Type string                 `json:"type"`
	Args map[string]interface{} `json:"args"`
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
	defaultConfig := Config{
		Chains: []ChainConfig{
			{
				ID:       "mainChain",
				Priority: 1,
				TickRate: 40,
				NumLamps: 50,
				Effects: []EffectConfig{
					{
						Type: "rainbow",
						Args: make(map[string]interface{}),
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