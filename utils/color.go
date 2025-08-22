package utils

import (
	"fmt"
	"godmx/dmx"
	"strconv"
	"strings"
	"math"
)

// ParseHexColor converts a hex color string (e.g., "#RRGGBB" or "RRGGBB") to a dmx.Lamp.
// W channel is set to 0.
func ParseHexColor(hex string) (dmx.Lamp, error) {
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) != 6 {
		return dmx.Lamp{}, fmt.Errorf("invalid hex color length: %s", hex)
	}

	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return dmx.Lamp{}, fmt.Errorf("invalid red component: %w", err)
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return dmx.Lamp{}, fmt.Errorf("invalid green component: %w", err)
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return dmx.Lamp{}, fmt.Errorf("invalid blue component: %w", err)
	}

	return dmx.Lamp{R: uint8(r), G: uint8(g), B: uint8(b), W: 0}, nil
}

// HsvToRgb converts an HSV color value to RGB.
// h is from 0 to 1; s and v are from 0 to 1.
// r, g, b are from 0 to 255.
func HsvToRgb(h, s, v float64) (uint8, uint8, uint8) {
	if s == 0 {
		return uint8(v * 255), uint8(v * 255), uint8(v * 255)
	}

	h = math.Mod(h, 1.0)
	i := math.Floor(h * 6)
	f := h*6 - i
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)

	var r, g, b float64
	switch int(i) % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}

	return uint8(r * 255), uint8(g * 255), uint8(b * 255)
}

// RgbToHsv converts an RGB color value to HSV.
// r, g, b are from 0 to 255.
// h is from 0 to 1; s and v are from 0 to 1.
func RgbToHsv(r, g, b uint8) (h, s, v float64) {
	rF, gF, bF := float64(r)/255, float64(g)/255, float64(b)/255

	max := math.Max(rF, math.Max(gF, bF))
	min := math.Min(rF, math.Min(gF, bF))
	delta := max - min

	v = max

	if max == 0 {
		s = 0
	} else {
		s = delta / max
	}

	if delta == 0 {
		h = 0
	} else {
		switch max {
		case rF:
			h = (gF - bF) / delta
			if gF < bF {
				h += 6.0
			}
		case gF:
			h = (bF - rF) / delta + 2.0
		case bF:
			h = (rF - gF) / delta + 4.0
		}
		h /= 6.0
	}
	return
}