package outputs

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
	"encoding/base64"

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
func NewGoveeOutput(goveeConfig config.GoveeOutputConfig, channelMapping string, numChannelsPerLamp int) (*GoveeOutput, error) {
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

		// Send activation command
		err = goOutput.sendRazerCommand(goveeDevice{config: devConfig, conn: conn}, "uwABsQEK")
		if err != nil {
			log.Printf("Failed to send Govee activation command to %s: %v\n", devConfig.IPAddress, err)
			// Decide if you want to continue or stop if activation fails
			// For now, we'll just log and continue
		}
	}

	if len(goOutput.devices) == 0 {
		return nil, fmt.Errorf("no Govee devices successfully initialized")
	}

	return goOutput, nil
}

// Send sends the lamp data to Govee devices.
func (g *GoveeOutput) Send(lamps []dmx.Lamp) error {
	
	if len(lamps) == 0 {
		return nil
	}

	// Prepare color data for the razer packet
	// Collect all RGB values from all lamps
	var colors []byte
	for _, lamp := range lamps {
		colors = append(colors, lamp.R, lamp.G, lamp.B)
	}

	// Create the razer packet
	razerPacket, err := createRazerPacket(colors)
	if err != nil {
		log.Printf("Failed to create Govee razer packet: %v\n", err)
		return err
	}

	// Base64 encode the razer packet
	encodedPacket := base64.StdEncoding.EncodeToString(razerPacket)

	for _, dev := range g.devices {
		// Govee devices have a rate limit of 100ms per command
		if time.Since(dev.lastSent) < 100*time.Millisecond {
			continue
		}

		// Send the razer command
		err = g.sendRazerCommand(dev, encodedPacket)
		if err != nil {
			log.Printf("Failed to send Govee razer command: %v\n", err)
		} else {
			dev.lastSent = time.Now()
		}
	}
	return nil
}


// Close closes all UDP connections for Govee devices.
func (g *GoveeOutput) Close() {
	for _, dev := range g.devices {
		if dev.conn != nil {
			err := dev.conn.Close()
			if err != nil {
				log.Printf("Error closing Govee UDP connection for %s: %v\n", dev.config.IPAddress, err)
			}
		}
	}
	log.Println("GoveeOutput connections closed.")
}

// calculateXORChecksumFast calculates the XOR checksum of a byte array.
func calculateXORChecksumFast(packet []byte) byte {
	var checksum byte
	for _, b := range packet {
		checksum ^= b
	}
	return checksum
}

// createRazerPacket creates a Govee razer packet from color data.
func createRazerPacket(colors []byte) ([]byte, error) {
	// This header is based on LedFX's pre_dreams header, modified for 0 segments and 0 stretch
	// BB 00 FA B0 00 (header) + 0x04 (color triples count)
	// The 0x04 in LedFX's pre_dreams is actually the count of color triples, not a fixed value.
	// So, the last byte of the header should be len(colors) / 3.
	header := []byte{0xBB, 0x00, 0xFA, 0xB0, 0x00, byte(len(colors) / 3)}

	// Concatenate header and colors
	fullPacket := append(header, colors...)

	// Calculate checksum and append
	checksum := calculateXORChecksumFast(fullPacket)
	fullPacket = append(fullPacket, checksum)

	return fullPacket, nil
}

// sendRazerCommand sends a razer command with a base64 encoded payload to a specific Govee device.
func (g *GoveeOutput) sendRazerCommand(dev goveeDevice, payload string) error {
	cmd := map[string]interface{}{
		"msg": map[string]interface{}{
			"cmd": "razer",
			"data": map[string]interface{}{
				"pt": payload,
			},
		},
	}

	jsonCmd, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("failed to marshal razer command: %w", err)
	}

	
	_, err = dev.conn.Write(jsonCmd)
	if err != nil {
		return fmt.Errorf("failed to send Govee razer command to %s (%s): %w", dev.config.IPAddress, dev.config.MACAddress, err)
	}
	
	return nil
}

