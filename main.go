package main

import (
	"flag"
	"fmt"
	"godmx/config"
	"godmx/orchestrator"
	"godmx/outputs"
	"godmx/utils"
	"godmx/webui"
	"os"
	"time"
)

func main() {
	// Command-line flags
	debug := flag.Bool("debug", false, "Run in debug mode (exits after a short duration)")
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	webPort := flag.Int("web-port", 8080, "Port for the web UI")
	eventName := flag.String("event", "", "Name of an event to trigger on startup")
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
	orch := orchestrator.NewOrchestrator(cfg)

	// Set initial global parameters from config
	orch.SetBPM(cfg.Globals.BPM)
	color1, _ := utils.ParseHexColor(cfg.Globals.Color1)
	orch.SetColor1(color1)
	color2, _ := utils.ParseHexColor(cfg.Globals.Color2)
	orch.SetColor2(color2)
	orch.SetIntensity(cfg.Globals.Intensity)

	// --- Build Chains from config ---
	for i := range cfg.Chains {
		chainConfig := &cfg.Chains[i]

		// Create Output for the chain
		var output orchestrator.Output
		switch chainConfig.Output.Type {
		case "artnet":
			ip, ok := chainConfig.Output.Args["ip"].(string)
			if !ok {
				fmt.Printf("ArtNet output 'ip' argument missing or invalid for chain %s.\n", chainConfig.ID)
				return
			}
			artNetOutput, err := outputs.NewArtNetOutput(ip, *debug, chainConfig.Output.ChannelMapping, chainConfig.Output.NumChannelsPerLamp)
			if err != nil {
				fmt.Printf("Error creating Art-Net output for chain %s: %v\n", chainConfig.ID, err)
				return
			}
			output = artNetOutput
		default:
			fmt.Printf("Unknown output type: %s for chain %s.\n", chainConfig.Output.Type, chainConfig.ID)
			return
		}

		// Create and add the chain
		chain := orchestrator.NewChain(chainConfig, orch, output)
		orch.AddChain(chain)
		chain.StartLoop()
	}

	// Start the web UI server
	webui.StartWebServer(orch, cfg, *webPort)

	fmt.Println("Orchestrator running.")

	// If an event is specified, trigger it now
	if *eventName != "" {
		orch.TriggerEvent(*eventName)
	}

	if *debug {
		fmt.Println("Debug mode: Running for 10 seconds...")
		select {
		case <-time.After(10 * time.Second):
			fmt.Println("Exiting debug mode.")
		}
	} else {
		fmt.Println("Press Ctrl+C to exit.")
		select {}
	}
}