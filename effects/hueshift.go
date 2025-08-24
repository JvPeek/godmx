package effects

import (
	"fmt"
	"godmx/dmx"
	"godmx/types"
	"godmx/utils"
	"math"
)

/*
Effect Name: Hue Shift
Description: Shifts the hue of the DMX data across the lamps, synchronized with the BPM.
Tags: [bpm_sensitive, transparent, color, pattern]
Parameters:
  - InternalName: direction
    DisplayName: Direction
    Description: The direction to shift the hue ('left' or 'right').
    DataType: string
    DefaultValue: "left"
  - InternalName: beatspan
    DisplayName: Beat Span
    Description: The number of beats for a full hue rotation.
    DataType: float64
    DefaultValue: 1.0
  - InternalName: huerange
    DisplayName: Hue Range
    Description: The total hue shift in degrees (0-360) over the beatspan.
    DataType: float64
    DefaultValue: 360.0
*/
func init() {
	RegisterEffect("hueshift", NewHueShift)
	RegisterEffectMetadata("hueshift", types.EffectMetadata{
		HumanReadableName: "Hue Shift",
		Description:       "Shifts the hue of the DMX data across the lamps, synchronized with the BPM.",
		Tags:              []string{"bpm_sensitive", "transparent", "color", "pattern"},
		Parameters: []types.ParameterMetadata{
			{
				InternalName: "direction",
				DisplayName:  "Direction",
				Description:  "The direction to shift the hue ('left' or 'right').",
				DataType:     "string",
				DefaultValue: "left",
			},
			{
				InternalName: "beatspan",
				DisplayName:  "Beat Span",
				Description:  "The number of beats for a full hue rotation.",
				DataType:     "float64",
				DefaultValue: 1.0,
			},
			{
				InternalName: "huerange",
				DisplayName:  "Hue Range",
				Description:  "The total hue shift in degrees (0-360) over the beatspan.",
				DataType:     "float64",
				DefaultValue: 360.0,
			},
		},
	})
}

// HueShift effect shifts the hue of the DMX data.
type HueShift struct {
	Direction string  // "left" or "right"
	BeatSpan  float64 // Number of beats for the huerange to complete
	HueRange  float64 // Total hue shift in degrees (0-360) over the BeatSpan
	accumulatedHueShift float64 // Internal state to accumulate hue shift over beats
	LastBeatProgress float64 // Stores BeatProgress from the previous frame to detect beat transitions
}

// NewHueShift creates a new HueShift effect.
func NewHueShift(args map[string]interface{}) (types.Effect, error) {
	direction, ok := args["direction"].(string)
	if !ok {
		return nil, fmt.Errorf("hueshift effect: missing or invalid 'direction' parameter")
	}
	if direction != "left" && direction != "right" {
		return nil, fmt.Errorf("hueshift effect: invalid direction '%s'. Must be 'left' or 'right'", direction)
	}

	beatSpan, ok := args["beatspan"].(float64)
	if !ok {
		return nil, fmt.Errorf("hueshift effect: missing or invalid 'beatspan' parameter")
	}
	hueRange, ok := args["huerange"].(float64)
	if !ok {
		return nil, fmt.Errorf("hueshift effect: missing or invalid 'huerange' parameter")
	}

	return &HueShift{Direction: direction, BeatSpan: beatSpan, HueRange: hueRange, accumulatedHueShift: 0.0, LastBeatProgress: 0.0}, nil
}

// Process applies the hueshift effect to the lamps.
func (s *HueShift) Process(lamps []dmx.Lamp, globals *types.OrchestratorGlobals, channelMapping string, numChannelsPerLamp int) {
	// Update accumulatedHueShift based on beat progress
	if globals.BeatProgress < s.LastBeatProgress {
		s.accumulatedHueShift += (1.0 - s.LastBeatProgress) + globals.BeatProgress
	} else {
		s.accumulatedHueShift += (globals.BeatProgress - s.LastBeatProgress)
	}

	s.accumulatedHueShift = math.Mod(s.accumulatedHueShift, s.BeatSpan)
	beatspanProgress := s.accumulatedHueShift / s.BeatSpan
	hueShiftAmount := beatspanProgress * (s.HueRange / 360.0)

	for i := range lamps {
		// Get current color
		r, g, b := lamps[i].R, lamps[i].G, lamps[i].B

		// Convert to HSV
		h, sat, val := utils.RgbToHsv(r, g, b)

		// Apply the hue shift
		if s.Direction == "left" {
			h += hueShiftAmount
		} else { // right
			h -= hueShiftAmount
		}

		// Wrap hue
		h = math.Mod(h, 1.0)
		if h < 0 {
			h += 1.0
		}

		// Convert back to RGB
		newR, newG, newB := utils.HsvToRgb(h, sat, val)

		// Update lamp with new RGB values
		lamps[i].R = newR
		lamps[i].G = newG
		lamps[i].B = newB
	}

	// Store current BeatProgress for the next frame's calculation
	s.LastBeatProgress = globals.BeatProgress
}