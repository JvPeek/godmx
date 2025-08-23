package effects

import (
	"fmt"
	"godmx/dmx"
	"godmx/types"
	"godmx/utils"
	"math"
)

func init() {
	RegisterEffect("hueshift", func(args map[string]interface{}) (types.Effect, error) {
		return NewHueShift(args)
	})
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
func NewHueShift(args map[string]interface{}) (*HueShift, error) {
	direction, ok := args["direction"].(string)
	if !ok || direction == "" {
		direction = "left" // Default to "left" if not provided or empty
	}

	if direction != "left" && direction != "right" {
		return nil, fmt.Errorf("invalid direction for hueshift effect: %s. Must be 'left' or 'right'", direction)
	}

	beatSpan := 1.0 // Default beatspan: 1 beat for a full rotation
	if bs, ok := args["beatspan"].(float64); ok {
		if bs > 0 {
			beatSpan = bs
		} else {
			return nil, fmt.Errorf("invalid beatspan for hueshift effect: %f. Must be greater than 0", bs)
		}
	} else if bs, ok := args["beatspan"].(int); ok {
		if float64(bs) > 0 {
			beatSpan = float64(bs)
		} else {
			return nil, fmt.Errorf("invalid beatspan for hueshift effect: %d. Must be greater than 0", bs)
		}
	}

	hueRange := 360.0 // Default huerange: full rotation
	if hr, ok := args["huerange"].(float64); ok {
		if hr >= 0 && hr <= 360 {
			hueRange = hr
		} else {
			return nil, fmt.Errorf("invalid huerange for hueshift effect: %f. Must be between 0 and 360", hr)
		}
	} else if hr, ok := args["huerange"].(int); ok {
		if float64(hr) >= 0 && float64(hr) <= 360 {
			hueRange = float64(hr)
		} else {
			return nil, fmt.Errorf("invalid huerange for hueshift effect: %d. Must be between 0 and 360", hr)
		}
	}

	return &HueShift{Direction: direction, BeatSpan: beatSpan, HueRange: hueRange, accumulatedHueShift: 0.0, LastBeatProgress: 0.0},
		nil
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
