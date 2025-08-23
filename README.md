# godmx

`godmx` is a powerful and flexible DMX lighting control system built in Go. It allows you to define complex lighting animations and effects through a chain-based configuration, synchronized with BPM.

## Chains

At the core of `godmx` are **Chains**. A chain represents a sequence of DMX lamps that are processed together. Each chain has its own set of effects and an output configuration. This allows for modular and scalable lighting setups.

Key properties of a chain:
- `id`: A unique identifier for the chain.
- `priority`: Determines the order of processing if multiple chains are active.
- `tickRate`: How often the chain's effects are processed per second.
- `numLamps`: The total number of DMX lamps in this chain.
- `effects`: An array of effects applied to the lamps in this chain.
- `output`: Defines how the processed DMX data is sent out (e.g., ArtNet, Govee).

## Effects Configuration

Effects are the building blocks of your lighting animations. They are configured within the `effects` array of each chain in your `config.json` (or `config.home.json`).

Each effect has the following common properties:
- `id`: A unique identifier for the effect within its chain.
- `type`: The type of effect (e.g., `rainbow`, `solidColor`, `shift`, `hueshift`).
- `args`: A JSON object containing parameters specific to the effect. These parameters control the behavior and appearance of the effect.
- `enabled`: A boolean indicating whether the effect is currently active.
- `group`: (Optional) A string to group related effects. If an effect has a `group` defined, only one effect within that group can be enabled at any given time. Enabling an effect in a group will automatically disable all other effects in the same group. This is enforced by the `EnforceGroupRules` logic within the orchestrator.

### Example Effect Configuration (from `config.home.json`)

Here are examples of how `shift` and `hueshift` effects are configured:

```json
{
  "chains": [
    {
      "id": "mainChain",
      "numLamps": 4,
      "effects": [
        {
          "id": "myShiftEffect",
          "type": "shift",
          "args": {
            "direction": "left",
            "speed": 0.5
          },
          "enabled": true
        },
        {
          "id": "myHueShiftEffect",
          "type": "hueshift",
          "args": {
            "direction": "right",
            "beatspan": 4.0,
            "huerange": 90.0
          },
          "enabled": true
        }
      ]
    }
  ]
}
```

- **`shift` effect parameters:**
  - `direction` (string): "left" or "right". The direction to shift the DMX data.
  - `speed` (float64): From 0.0 to 1.0. Controls the speed of the shift (1.0 being one full lamp length shift per beat).

- **`hueshift` effect parameters:**
  - `direction` (string): "left" or "right". The direction to shift the hue.
  - `beatspan` (float64): The number of beats over which the `huerange` animation completes. For example, `4.0` means the animation takes 4 beats to complete one cycle.
  - `huerange` (float64): The total hue shift in degrees (0-360) that occurs over the `beatspan`. For example, `90.0` means the hue will shift by 90 degrees over the defined `beatspan`.

## Building the Project

`godmx` can be compiled for various operating systems using the provided `Makefile`. The compiled binaries will be placed in the `build/` directory.

### Prerequisites

- Go (version 1.18 or higher recommended)

### Using the Makefile

To compile `godmx` for your desired platform(s), navigate to the project's root directory in your terminal and use the `make` command:

- **Compile for Linux (your current OS):**
  ```bash
  make linux
  ```

- **Compile for macOS:**
  ```bash
  make macos
  ```
  > **Note:** MIDI functionality is not supported on macOS due to underlying driver limitations.


- **Compile for Windows:**
  ```bash
  make windows
  ```

- **Compile for all supported platforms (Linux, macOS, Windows):**
  ```bash
  make all
  ```

After compilation, you will find the executable(s) in the `build/` directory.

## Running `godmx`

Once compiled, you can run `godmx` by executing the binary from the `build/` directory. You will typically need to provide the path to your configuration file:

```bash
./build/godmx_linux -config config.home.json
```

(Replace `godmx_linux` with `godmx_macos` or `godmx_windows.exe` as appropriate for your platform.)
