package main

import (
	"flag"
	"fmt"
	"godmx/config"
	"godmx/effects"
	"godmx/orchestrator"
	"godmx/outputs"
	"godmx/utils"
	"os"
	"time"
)

func main() {
	// Command-line flags
	debug := flag.Bool("debug", false, "Run in debug mode (exits after 10 seconds)")
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	flag.Parse()

	fmt.Println("Starting godmx...")

	// Check if config file exists, create if not
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		fmt.Printf("Config file not found at '%s', creating a default one.\n", *configPath)
		if err := config.CreateDefaultConfig(*configPath); err != nil {
			fmt.Printf("Error creating default config file: %v\n", err)
			return
		}
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return
	}

	// Create Orchestrator
	orch := orchestrator.NewOrchestrator()

	// Set global parameters from config
	orch.SetBPM(cfg.Globals.BPM)

	// Parse and set Color1
	color1, err := utils.ParseHexColor(cfg.Globals.Color1)
	if err != nil {
		fmt.Printf("Error parsing Color1 from config: %v\n", err)
		return
	}
	orch.SetColor1(color1)

	// Parse and set Color2
	color2, err := utils.ParseHexColor(cfg.Globals.Color2)
	if err != nil {
		fmt.Printf("Error parsing Color2 from config: %v\n", err)
		return
	}
	orch.SetColor2(color2)

	orch.SetIntensity(cfg.Globals.Intensity)

	// --- Build Chains, Effects, and Outputs from config ---
	// Assuming only one chain for now, as per config.json example
	if len(cfg.Chains) == 0 {
		fmt.Println("No chains defined in configuration.")
		return
	}
	for _, chainConfig := range cfg.Chains { // Iterate over all chain configs
		// Create Output
		var output orchestrator.Output
		channelMapping := chainConfig.Output.ChannelMapping
		if channelMapping == "" {
			channelMapping = "RGBW" // Default to RGBW
		}

		numChannelsPerLamp := chainConfig.Output.NumChannelsPerLamp
		if numChannelsPerLamp == 0 {
			numChannelsPerLamp = 4 // Default to 4 channels per lamp
		}

		switch chainConfig.Output.Type {
		case "artnet":
			ip, ok := chainConfig.Output.Args["ip"].(string)
			if !ok {
				fmt.Printf("ArtNet output 'ip' argument missing or invalid for chain %s.\n", chainConfig.ID)
				return
			}

			artNetOutput, err := outputs.NewArtNetOutput(ip, *debug, channelMapping, numChannelsPerLamp)
			if err != nil {
				fmt.Printf("Error creating Art-Net output for chain %s: %v\n", chainConfig.ID, err)
				return
			}
			output = artNetOutput
		default:
			fmt.Printf("Unknown output type: %s for chain %s.\n", chainConfig.Output.Type, chainConfig.ID)
			return
		}

		// Create Chain
		mainChain := orchestrator.NewChain(chainConfig.ID, chainConfig.Priority, chainConfig.TickRate, output, chainConfig.NumLamps, orch, channelMapping, numChannelsPerLamp)

		// Create Effects
		for _, effectConfig := range chainConfig.Effects {
			var effect orchestrator.Effect
			constructor, ok := effects.GetEffectConstructor(effectConfig.Type)
			if !ok {
				fmt.Printf("Unknown effect type: %s for chain %s. Available effects are: %v.\n", effectConfig.Type, chainConfig.ID, effects.GetAvailableEffects())
				return
			}
			var err error
			effect, err = constructor(effectConfig.Args)
			if err != nil {
				fmt.Printf("Error creating effect %s for chain %s: %v\n", effectConfig.Type, chainConfig.ID, err)
				return
			}
			mainChain.AddEffect(effect)
		}

		// Add Chain to Orchestrator
		orch.AddChain(mainChain)

		// Start the Chain's loop
		mainChain.StartLoop() // Start each chain's loop independently
	}

	fmt.Println("Orchestrator running.")

	if *debug {
		fmt.Println("Debug mode: Exiting in 10 seconds...")
		select {
		case <-time.After(10 * time.Second):
			fmt.Println("Exiting debug mode.")
		}
	} else {
		fmt.Println("Press Ctrl+C to exit.")
		// Keep main goroutine alive
		select {}
	}
}
