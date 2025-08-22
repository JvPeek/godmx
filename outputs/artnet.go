package outputs

import (
	"godmx/dmx"

	"github.com/RickHulzinga/go-simple-artnet/node"
)

// ArtNetOutput sends DMX data to an Art-Net node.
type ArtNetOutput struct {
	node  *node.ArtNetNode
	debug bool // Added debug field
}

// NewArtNetOutput creates a new ArtNetOutput.
func NewArtNetOutput(targetIP string, debug bool) (*ArtNetOutput, error) { // Added debug argument
	n, err := node.NewArtNetNode(targetIP + ":6454")
	if err != nil {
		return nil, err
	}
	n.Start()
	return &ArtNetOutput{
		node:  n,
		debug: debug, // Set debug field
	}, nil
}

// Send sends the lamp data as DMX to the Art-Net node.
func (a *ArtNetOutput) Send(lamps []dmx.Lamp) error {
	universe := a.node.GetUniverse(0)
	for i, lamp := range lamps {
		if i*4+3 < 512 {
			universe.SetChannel(i*4+1, int(lamp.R))
			universe.SetChannel(i*4+2, int(lamp.G))
			universe.SetChannel(i*4+3, int(lamp.B))
			universe.SetChannel(i*4+4, int(lamp.W))
		}
	}
	return nil
}