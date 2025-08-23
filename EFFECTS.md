# Registered Effects

This document outlines all available lighting effects, their descriptions, and configurable parameters.

## Blink

Alternates between two colors based on the global BPM, creating a blinking effect.

**Tags**: bpm_sensitive, color_source, pattern

### Parameters

| Parameter Name | Display Name | Data Type | Default Value | Min Value | Max Value | Description |
|----------------|--------------|-----------|---------------|-----------|-----------|-------------|
| divider | Divider | int | 1 | 1 | - | Divides the beat into segments for faster blinking. |
| dutyCycle | Duty Cycle | float64 | 0.5 | 0 | 1 | Percentage of the segment that Color1 is shown. |

---

## Cyberfall

Simulates digital rain, acting as a brightness mask over existing colors.

**Tags**: brightness_mask, random, transparent

### Parameters

| Parameter Name | Display Name | Data Type | Default Value | Min Value | Max Value | Description |
|----------------|--------------|-----------|---------------|-----------|-----------|-------------|
| density | Density | float64 | 0.5 | 0 | 1 | How many 'active' columns are falling (0.0 - 1.0). |
| flicker_intensity | Flicker Intensity | float64 | 0.1 | 0 | 1 | Random variation applied to brightness (0.0 - 1.0). |
| max_brightness | Max Brightness | int | 255 | 0 | 255 | Maximum brightness for bright parts (0-255). |
| min_brightness | Min Brightness | int | 0 | 0 | 255 | Minimum brightness for dark parts (0-255). |
| speed | Speed | float64 | 1 | 0 | - | How fast the 'rain' falls. |
| trail_length | Trail Length | int | 10 | 0 | - | How long the 'tail' of the falling light is (in lamps). |

---

## Darkwave

Creates a dark wave that travels across the lamps, dimming them based on a sine wave.

**Tags**: bpm_sensitive, brightness_mask, pattern, transparent

### Parameters

| Parameter Name | Display Name | Data Type | Default Value | Min Value | Max Value | Description |
|----------------|--------------|-----------|---------------|-----------|-----------|-------------|
| percentage | Percentage | float64 | 0.5 | 0 | 1 | The maximum percentage of dimming applied by the wave (0.0 - 1.0). |
| speed | Speed | float64 | 1 | 0 | - | How fast the dark wave travels. |

---

## Dim

Dims all lamps by a specified percentage.

**Tags**: brightness_mask, transparent

### Parameters

| Parameter Name | Display Name | Data Type | Default Value | Min Value | Max Value | Description |
|----------------|--------------|-----------|---------------|-----------|-----------|-------------|
| percentage | Percentage | float64 | 0.5 | 0 | 1 | The percentage to dim the lamps by (0.0 - 1.0). |

---

## Gradient

Creates a smooth color gradient across the lamps, interpolating between global Color1 and Color2.

**Tags**: color_source, pattern

---

## Hue Shift

Shifts the hue of the DMX data across the lamps, synchronized with the BPM.

**Tags**: bpm_sensitive, color, pattern, transparent

### Parameters

| Parameter Name | Display Name | Data Type | Default Value | Min Value | Max Value | Description |
|----------------|--------------|-----------|---------------|-----------|-----------|-------------|
| beatspan | Beat Span | float64 | 1 | - | - | The number of beats for a full hue rotation. |
| direction | Direction | string | left | - | - | The direction to shift the hue ('left' or 'right'). |
| huerange | Hue Range | float64 | 360 | - | - | The total hue shift in degrees (0-360) over the beatspan. |

---

## Rainbow

Generates a static rainbow spectrum across the lamps.

**Tags**: color_source, pattern

---

## Shift

Shifts the DMX data (colors) across the lamps either left or right, synchronized with the BPM.

**Tags**: bpm_sensitive, pattern, transform, transparent

### Parameters

| Parameter Name | Display Name | Data Type | Default Value | Min Value | Max Value | Description |
|----------------|--------------|-----------|---------------|-----------|-----------|-------------|
| direction | Direction | string | left | - | - | The direction to shift the lamps ('left' or 'right'). |
| speed | Speed | float64 | 1 | - | - | The speed of the shift, from 0 to 1 (1 being 1 shift per beat). |

---

## Solid Color

Sets all lamps to a single color defined by global Color1.

**Tags**: color_source

---

## Twinkle

Randomly turns a percentage of lamps to white at the beginning of each beat, creating a twinkling effect.

**Tags**: bpm_sensitive, color_source, pattern, random

### Parameters

| Parameter Name | Display Name | Data Type | Default Value | Min Value | Max Value | Description |
|----------------|--------------|-----------|---------------|-----------|-----------|-------------|
| percentage | Percentage | float64 | 0.1 | 0 | 1 | The percentage of lamps to twinkle (0.0 - 1.0). |

---

## Whiteout

Sets all lamps to full white, overriding any previous colors.

**Tags**: color_source

---

