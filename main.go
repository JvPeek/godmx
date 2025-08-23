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
	"godmx/effects"
	"godmx/midi"
)

func main() {
	// Command-line flags
	debug := flag.Bool("debug", false, "Run in debug mode (exits after a short duration)")
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	webPort := flag.Int("web-port", 8080, "Port for the web UI")
	eventName := flag.String("event", "", "Name of an event to trigger on startup")
	docs := flag.Bool("docs", false, "Generate documentation for effects in EFFECTS.md")
	flag.Parse()

	// Generate documentation if -docs flag is present
	if *docs {
		fmt.Println("Generating EFFECTS.md documentation...")
		if err := effects.GenerateEffectDocumentation(); err != nil {
			fmt.Printf("Error generating documentation: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Documentation generated successfully.")
		return
	}

	fmt.Println("Starting GoDMX...")

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		// If config file doesn't exist, create a default one and try to load again
		if os.IsNotExist(err) {
			fmt.Printf("Config file not found at '%s', creating a default one.\n", *configPath)
			defaultCfg := config.CreateDefaultConfig()
			if err := config.SaveConfig(&defaultCfg, *configPath); err != nil {
				fmt.Printf("Error creating default config file: %v\n", err)
				return
			}
			// Try loading again after creating default
			cfg, err = config.LoadConfig(*configPath)
			if err != nil {
				fmt.Printf("Error loading configuration after creating default: %v\n", err)
				return
			}
		} else {
			return
		}
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
		case "ddp":
			ip, ok := chainConfig.Output.Args["ip"].(string)
			if !ok {
				fmt.Printf("DDP output 'ip' argument missing or invalid for chain %s.\n", chainConfig.ID)
				return
			}
			ddpOutput, err := outputs.NewDDPOutput(ip, *debug, chainConfig.Output.ChannelMapping, chainConfig.Output.NumChannelsPerLamp)
			if err != nil {
				fmt.Printf("Error creating DDP output for chain %s: %v\n", chainConfig.ID, err)
				return
			}
			output = ddpOutput
		case "govee":
			goveeOutput, err := outputs.NewGoveeOutput(chainConfig.Output.Govee, chainConfig.Output.ChannelMapping, chainConfig.Output.NumChannelsPerLamp)
			if err != nil {
				fmt.Printf("Error creating Govee output for chain %s: %v\n", chainConfig.ID, err)
				return
			}
			output = goveeOutput
		default:
			fmt.Printf("Unknown output type: %s for chain %s.\n", chainConfig.Output.Type, chainConfig.ID)
			return
		}

		// Create and add the chain
		chain := orchestrator.NewChain(chainConfig, orch, output)
		orch.AddChain(chain)
		chain.StartLoop()
	}

	fmt.Printf("Checking MIDI triggers. Count: %d\n", len(cfg.MidiTriggers))
	// Initialize and start MIDI controller if triggers are configured
	if len(cfg.MidiTriggers) > 0 {
		fmt.Println("MIDI triggers found. Initializing MIDI controller...")
		midiController, err := midi.NewMidiController(orch, cfg.MidiTriggers, cfg.MidiPortName)
		if err != nil {
			fmt.Printf("Error initializing MIDI controller: %v\n", err)
			// Continue without MIDI, or exit? For now, continue.
		} else {
			if err := midiController.Start(); err != nil {
				fmt.Printf("Error starting MIDI controller: %v\n", err)
				// Continue without MIDI, or exit? For now, continue.
			} else {
				defer midiController.Stop() // Ensure MIDI controller is stopped on exit
				fmt.Println("MIDI controller started successfully.")
			}
		}
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

