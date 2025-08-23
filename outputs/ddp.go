package outputs

import (
	"encoding/binary"
	"godmx/dmx"
	"net"
)

const (
	// DDP protocol constants from the spec you so graciously provided.
	ddpPort           = 4048
	ddpHeaderLen      = 10
	ddpFlags1Ver1     = 0x40
	ddpFlags1Push     = 0x01
	ddpDataTypeRGB    = 0x01
	ddpIdDisplay      = 0x01
	ddpMaxDataLen     = 1440 * 3 // A reasonable max to fit in a standard MTU
)

// DDPOutput sends DMX data to a DDP-compliant controller like WLED.
type DDPOutput struct {
	conn               net.Conn
	debug              bool
	channelMapping     string
	numChannelsPerLamp int
	sequence           byte
}

// NewDDPOutput creates a new DDPOutput.
func NewDDPOutput(targetIP string, debug bool, channelMapping string, numChannelsPerLamp int) (*DDPOutput, error) {
	conn, err := net.Dial("udp", targetIP+":4048")
	if err != nil {
		return nil, err
	}

	return &DDPOutput{
		conn:               conn,
		debug:              debug,
		channelMapping:     channelMapping,
		numChannelsPerLamp: numChannelsPerLamp,
		sequence:           0,
	}, nil
}

// Send sends the lamp data as DDP to the controller.
func (d *DDPOutput) Send(lamps []dmx.Lamp) error {
	pixelData := make([]byte, len(lamps)*3) // 3 bytes per pixel for RGB
	for i, lamp := range lamps {
		// FUCK RGBW for now, DDP is an RGB world.
		pixelData[i*3] = byte(lamp.R)
		pixelData[i*3+1] = byte(lamp.G)
		pixelData[i*3+2] = byte(lamp.B)
	}

	// Increment sequence number, wrapping around after 15.
	d.sequence = (d.sequence % 15) + 1

	// Now we build the packet, piece by miserable piece.
	header := make([]byte, ddpHeaderLen)

	// Byte 0: Flags
	header[0] = ddpFlags1Ver1 | ddpFlags1Push // Version 1, Push mode

	// Byte 1: Sequence Number
	header[1] = d.sequence

	// Byte 2: Data Type
	header[2] = ddpDataTypeRGB

	// Byte 3: Destination ID
	header[3] = ddpIdDisplay

	// Bytes 4-7: Data Offset (32-bit, MSB first)
	// We are sending the whole frame, so offset is 0.
	binary.BigEndian.PutUint32(header[4:8], 0)

	// Bytes 8-9: Data Length (16-bit, MSB first)
	binary.BigEndian.PutUint16(header[8:10], uint16(len(pixelData)))

	// Combine header and data into a single packet of glorious bytes.
	packet := append(header, pixelData...)

	// Send it to the ether.
	_, err := d.conn.Write(packet)
	return err
}

// Close closes the DDP connection.
func (d *DDPOutput) Close() {
	if d.conn != nil {
		d.conn.Close()
	}
}