package outputs

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"godmx/dmx"
	"godmx/config"
)

const (
	goveeControlPort = 4003
)

// GoveeOutput sends DMX data to Govee devices via LAN control.
type GoveeOutput struct {
	devices []goveeDevice
}

type goveeDevice struct {
	config config.GoveeDeviceConfig
	conn   *net.UDPConn
	lastSent time.Time
}

// NewGoveeOutput creates a new GoveeOutput.
func NewGoveeOutput(goveeConfig config.GoveeOutputConfig) (*GoveeOutput, error) {
	goOutput := &GoveeOutput{} 

	for _, devConfig := range goveeConfig.Devices {
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", devConfig.IPAddress, goveeControlPort))
		if err != nil {
			log.Printf("Failed to resolve Govee device address %s: %v\n", devConfig.IPAddress, err)
			continue
		}

		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			log.Printf("Failed to dial Govee device %s: %v\n", devConfig.IPAddress, err)
			continue
		}

		goOutput.devices = append(goOutput.devices, goveeDevice{
			config: devConfig,
			conn:   conn,
			lastSent: time.Now(),
		})
	}

	if len(goOutput.devices) == 0 {
		return nil, fmt.Errorf("no Govee devices successfully initialized")
	}

	return goOutput, nil
}

// Process sends the lamp data to Govee devices.
func (g *GoveeOutput) Process(lamps []dmx.Lamp) {
	// For simplicity, we'll just use the first lamp's color for all Govee devices
	// In a more complex setup, you'd map DMX channels to individual Govee devices
	if len(lamps) == 0 {
		return
	}

	// Govee brightness is 1-100, DMX is 0-255. Scale it.
	brightness := uint8(float64(lamps[0].R+lamps[0].G+lamps[0].B) / 3.0 / 2.55) // Average RGB and scale to 0-100
	if brightness == 0 { brightness = 1 } // Govee brightness can't be 0, min is 1
	if brightness > 100 { brightness = 100 }

	r := lamps[0].R
	g := lamps[0].G
	b := lamps[0].B

	// Construct Govee JSON command
	cmd := map[string]interface{}{
		"msg": map[string]interface{}{
			"cmd": "color",
			"data": map[string]interface{}{
				"r": r,
				"g": g,
				"b": b,
				"brightness": brightness,
			},
		},
	}

	jsonCmd, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Failed to marshal Govee command: %v\n", err)
		return
	}

	for _, dev := range g.devices {
		// Add MAC address to the command for Govee devices
		// This is often required by Govee's protocol, even for LAN control
		// The exact placement might vary, but it's common to include it in the top-level msg
		// Let's assume it needs to be in the data field for now, based on some examples
		// Re-marshal with MAC if needed, or modify the map before marshalling
		// For now, let's just send the color command as is, assuming MAC is handled by device discovery
		// or not strictly required for direct control.

		// Govee devices have a rate limit of 100ms per command
		if time.Since(dev.lastSent) < 100*time.Millisecond {
			continue
		}

		_, err := dev.conn.Write(jsonCmd)
		if err != nil {
			log.Printf("Failed to send Govee command to %s (%s): %v\n", dev.config.IPAddress, dev.config.MACAddress, err)
		} else {
			// log.Printf("Sent Govee command to %s: %s\n", dev.config.IPAddress, string(jsonCmd))
			dev.lastSent = time.Now()
		}
	}
}
