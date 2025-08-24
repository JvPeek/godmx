package webui

import (
	"embed"
	"encoding/json"
	"fmt"
	"godmx/config"
	"godmx/orchestrator"
	"io/fs"
	"log"
	"net/http"
	"path"
	"sort"
	"time"
)

//go:embed web
var content embed.FS

// staticHandler serves static files with correct MIME types.
type staticHandler struct {
	fs http.FileSystem
}

func (h *staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Prepend "web/" to the path to access files within the embedded "web" directory
	filePath := path.Join("web", r.URL.Path)

	// Open the file from the embedded file system
	f, err := h.fs.Open(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()

	// Determine content type based on file extension
	switch path.Ext(r.URL.Path) {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	default:
		// Let http.ServeContent determine the content type for other files
		// or if it's an unknown extension
	}

	// Serve the file
	http.ServeContent(w, r, r.URL.Path, time.Time{}, f)
}



// ChainConfig represents a simplified chain configuration for JSON serialization
type ChainConfig struct {
	ID        string `json:"ID"`
	Priority  int    `json:"Priority"`
	TickRate  int    `json:"TickRate"`
	NumLamps  int    `json:"NumLamps"`
	Output    OutputConfig `json:"Output"`
	Effects   []EffectConfig `json:"Effects"`
}

// EffectConfig represents a simplified effect configuration for JSON serialization
type EffectConfig struct {
	Type    string                 `json:"Type"`
	Args    map[string]interface{} `json:"Args"`
	Enabled *bool                  `json:"Enabled,omitempty"`
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
		http.Handle("/static/", http.StripPrefix("/static/", &staticHandler{http.FS(content)}))

	// Serve index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			http.NotFound(w, r)
			return
		}
		indexHTML, err := fs.ReadFile(content, "web/index.html")
		if err != nil {
			http.Error(w, "Could not read index.html", http.StatusInternalServerError)
			log.Printf("Error reading index.html: %v", err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexHTML)
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
				simplifiedEffects = append(simplifiedEffects, EffectConfig{Type: effectCfg.Type, Args: effectCfg.Args, Enabled: effectCfg.Enabled})
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

	// API endpoint to list all events
	http.HandleFunc("/api/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var eventNames []string
		for name := range cfg.Actions {
			eventNames = append(eventNames, name)
		}
		sort.Strings(eventNames) // Sort event names alphabetically
		json.NewEncoder(w).Encode(eventNames)
		})

	// API endpoint to trigger an event
	http.HandleFunc("/api/trigger", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var data struct {
			EventName string `json:"eventName"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		orch.TriggerEvent(data.EventName)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Event triggered"})
		})

	log.Printf("Web UI server starting on port %d\n", port)
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			log.Fatalf("Failed to start web server: %v", err)
		}
	}()
}
