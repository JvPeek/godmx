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
	Type string                 `json:"type"`
	Args map[string]interface{} `json:"args"`
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