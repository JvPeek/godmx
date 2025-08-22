package webui

import (
	"encoding/json"
	"fmt"
	"godmx/config"
	"godmx/orchestrator"
	"log"
	"net/http"
	"path/filepath"
)

// ChainConfig represents a simplified chain configuration for JSON serialization
type ChainConfig struct {
	ID        string `json:"ID"`
	Priority  int    `json:"Priority"`
	TickRate  int    `json:"TickRate"`
	NumLamps  int    `json:"NumLamps"`
	Output    OutputConfig `json:"Output"`
	Effects   []EffectConfig `json:"Effects"`
}

// OutputConfig represents a simplified output configuration for JSON serialization
type EffectConfig struct {
	Type string                 `json:"Type"`
	Args map[string]interface{} `json:"Args"`
}

// OutputConfig represents a simplified output configuration for JSON serialization
type OutputConfig struct {
	Type               string                 `json:"Type"`
	Args               map[string]interface{} `json:"Args"`
	ChannelMapping     string                 `json:"ChannelMapping"`
	NumChannelsPerLamp int                    `json:"NumChannelsPerLamp"`
}



// StartWebServer starts the HTTP server for the web UI.
func StartWebServer(orch *orchestrator.Orchestrator, cfg *config.Config, port int) {
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web"))))

	// Serve index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join("web", "index.html"))
	})

	// API endpoint for chains
	http.HandleFunc("/api/chains", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Convert orchestrator chains to simplified ChainConfig for JSON serialization
		var simplifiedChains []ChainConfig
		for _, chainCfg := range cfg.Chains {
			simplifiedOutput := OutputConfig{
				Type:               chainCfg.Output.Type,
				Args:               chainCfg.Output.Args,
				ChannelMapping:     chainCfg.Output.ChannelMapping,
				NumChannelsPerLamp: chainCfg.Output.NumChannelsPerLamp,
			}
			var simplifiedEffects []EffectConfig
			for _, effectCfg := range chainCfg.Effects {
				simplifiedEffects = append(simplifiedEffects, EffectConfig{Type: effectCfg.Type, Args: effectCfg.Args})
			}
			simplifiedChains = append(simplifiedChains, ChainConfig{
				ID:        chainCfg.ID,
				Priority:  chainCfg.Priority,
				TickRate:  chainCfg.TickRate,
				NumLamps:  chainCfg.NumLamps,
				Output:    simplifiedOutput,
				Effects:   simplifiedEffects,
			})
		}

		json.NewEncoder(w).Encode(simplifiedChains)
	})

	// API endpoint for BPM
	http.HandleFunc("/api/bpm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		if r.Method == http.MethodPost {
			var data struct {
				BPM float64 `json:"bpm"`
			}
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			orch.SetBPM(data.BPM)
			log.Printf("BPM updated to: %.2f", data.BPM)
		}
		
		json.NewEncoder(w).Encode(map[string]float64{"bpm": orch.GetGlobals().BPM})
	})

	log.Printf("Web UI server starting on port %d\n", port)
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			log.Fatalf("Failed to start web server: %v", err)
		}
	}()
}
