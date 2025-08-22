package main

import (
	"flag"
	"fmt"
	"godmx/config"
	"godmx/effects"
	"godmx/orchestrator"
	"godmx/outputs"
	"godmx/utils"
	"time"
)

func main() {
	// Command-line flags
	debug := flag.Bool("debug", false, "Run in debug mode (exits after 10 seconds)")
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	flag.Parse()

	fmt.Println("Starting godmx...")

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
	chainConfig := cfg.Chains[0] // Get the first chain config

	// Create Output
	var output orchestrator.Output
	switch chainConfig.Output.Type {
	case "artnet":
		ip, ok := chainConfig.Output.Args["ip"].(string)
		if !ok {
			fmt.Println("ArtNet output 'ip' argument missing or invalid.")
			return
		}
		artNetOutput, err := outputs.NewArtNetOutput(ip, *debug) // Pass debug flag
		if err != nil {
			fmt.Printf("Error creating Art-Net output: %v\n", err)
			return
		}
		output = artNetOutput
	default:
		fmt.Printf("Unknown output type: %s\n", chainConfig.Output.Type)
		return
	}

	// Create Chain
	mainChain := orchestrator.NewChain(chainConfig.ID, chainConfig.Priority, chainConfig.TickRate, output, chainConfig.NumLamps, orch)

	// Create Effects
	for _, effectConfig := range chainConfig.Effects {
		var effect orchestrator.Effect
		switch effectConfig.Type {
		case "rainbow":
			effect = effects.NewRainbow()
		case "solidColor":
			effect = &effects.SolidColor{}
		case "gradient":
			effect = &effects.Gradient{}
		case "blink":
			effect = effects.NewBlink()
		case "twinkle":
			var err error
			effect, err = effects.NewTwinkle(effectConfig.Args)
			if err != nil {
				fmt.Printf("Error creating twinkle effect: %v\n", err)
				return
			}
		// Add other effect types here
		default:
			fmt.Printf("Unknown effect type: %s\n", effectConfig.Type)
			return
		}
		mainChain.AddEffect(effect)
	}

	// Add Chain to Orchestrator
	orch.AddChain(mainChain)

	// Start the Chain's loop
	mainChain.StartLoop()

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
