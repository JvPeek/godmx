package outputs

import (
	"godmx/dmx"

	"github.com/RickHulzinga/go-simple-artnet/node"
)

// ArtNetOutput sends DMX data to an Art-Net node.
type ArtNetOutput struct {
	node               *node.ArtNetNode
	debug              bool // Added debug field
	channelMapping     string
	numChannelsPerLamp int
}

// NewArtNetOutput creates a new ArtNetOutput.
func NewArtNetOutput(targetIP string, debug bool, channelMapping string, numChannelsPerLamp int) (*ArtNetOutput, error) {
	n, err := node.NewArtNetNode(targetIP + ":6454")
	if err != nil {
		return nil, err
	}
	n.Start()
	return &ArtNetOutput{
		node:               n,
		debug:              debug,
		channelMapping:     channelMapping,
		numChannelsPerLamp: numChannelsPerLamp,
	}, nil
}

// Send sends the lamp data as DMX to the Art-Net node.
func (a *ArtNetOutput) Send(lamps []dmx.Lamp) error {
	universe := a.node.GetUniverse(0)
	for i, lamp := range lamps {
		baseChannel := i * a.numChannelsPerLamp
		if baseChannel+a.numChannelsPerLamp-1 < 512 { // Ensure we don't go out of bounds
			switch a.channelMapping {
			case "RGB":
				universe.SetChannel(baseChannel+1, int(lamp.R))
				universe.SetChannel(baseChannel+2, int(lamp.G))
				universe.SetChannel(baseChannel+3, int(lamp.B))
			case "RGBW":
				universe.SetChannel(baseChannel+1, int(lamp.R))
				universe.SetChannel(baseChannel+2, int(lamp.G))
				universe.SetChannel(baseChannel+3, int(lamp.B))
				universe.SetChannel(baseChannel+4, int(lamp.W))
			default:
				// Default to RGBW if mapping is unknown or not provided
				universe.SetChannel(baseChannel+1, int(lamp.R))
				universe.SetChannel(baseChannel+2, int(lamp.G))
				universe.SetChannel(baseChannel+3, int(lamp.B))
				universe.SetChannel(baseChannel+4, int(lamp.W))
			}
		}
	}
	return nil
}

// Close stops the Art-Net node.
func (a *ArtNetOutput) Close() {
	a.node.Stop()
}
