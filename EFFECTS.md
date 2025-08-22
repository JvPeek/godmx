# Registered Effects

| Effect Name | Description | Parameters | Tags |
|-------------|-------------|------------|------|
| Blink | Alternates between two colors based on the global BPM, creating a blinking effect. | Divider (int, default: 1), Duty Cycle (float64, default: 0.5) | bpm_sensitive, color_source, pattern |
| Darkwave | Creates a dark wave that travels across the lamps, dimming them based on a sine wave. | Percentage (float64, default: 0.5), Speed (float64, default: 1) | bpm_sensitive, brightness_mask, pattern, transparent |
| Dim | Dims all lamps by a specified percentage. | Percentage (float64, default: 0.5) | brightness_mask, transparent |
| Gradient | Creates a smooth color gradient across the lamps, interpolating between global Color1 and Color2. |  | color_source, pattern |
| Rainbow | Generates a static rainbow spectrum across the lamps. |  | color_source, pattern |
| Shift | Shifts the DMX data (colors) across the lamps either left or right, synchronized with the BPM. | Direction (string, default: left) | bpm_sensitive, pattern, transform, transparent |
| Solid Color | Sets all lamps to a single color defined by global Color1. |  | color_source |
| Twinkle | Randomly turns a percentage of lamps to white at the beginning of each beat, creating a twinkling effect. | Percentage (float64, default: 0.1) | bpm_sensitive, color_source, pattern, random |
| Whiteout | Sets all lamps to full white, overriding any previous colors. |  | color_source |
