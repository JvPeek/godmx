# GoDMX

## What is this project?

`GoDMX` is a powerful and flexible DMX lighting control system built in Go. Think of it as a modular, software-defined Eurorack for lights, allowing you to define complex lighting animations and effects through a chain-based configuration, synchronized with BPM.

## Warning: For Technical Users Only!

This software is built for technical people who want raw, command-line control over their lights. It does *not* provide a full graphical user interface (GUI) for configuration or control. If you're comfortable with terminals, configuration files, and a bit of DIY, you'll feel right at home. Otherwise, you might find `GoDMX` a bit... opinionated.

## Features

- [ ] **Chain-based DMX control:** Organize your lights into logical groups.
- [ ] **Extensible Effect System:** Apply various lighting effects (rainbow, solid color, shift, hueshift, blink, etc.).
- [ ] **Event-driven Automation:** Trigger complex sequences of actions via MIDI or Web UI.

- [ ] **Cross-platform Compilation:** Build for Linux, macOS, and Windows.
- [ ] **Web-based UI:** Monitor and control your setup from a browser.
- [ ] **ArtNet Output:** Control DMX fixtures over ArtNet.
- [ ] **Govee Output:** Control Govee smart lights.

## Core Principles

`GoDMX` is built with a few core philosophies in mind:

*   **Small Footprint, High Performance:** Designed to be lightweight and efficient, ideal for embedded systems or setups where resources are limited.
*   **Raw Control, Not Hand-holding:** This isn't a drag-and-drop GUI tool. It provides powerful building blocks and expects you to configure them via text files and command-line flags. If you know a bit of coding, you'll find it highly customizable.
*   **Building Blocks, Not Ready-made Effects:** Instead of a library of pre-defined, rigid effects, `GoDMX` offers fundamental components that you can combine and customize to create unique lighting animations. This gives you ultimate flexibility.

## Installing

`GoDMX` can be compiled for various operating systems using the provided `Makefile`. The compiled binaries will be placed in the `build/` directory.

### Prerequisites

- Go (version 1.18 or higher recommended)

### Using the Makefile

To compile `GoDMX` for your desired platform(s), navigate to the project's root directory in your terminal and use the `make` command:

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

- **Install `GoDMX` to `/usr/local/bin` (Linux only, tested):**
  ```bash
  sudo make install
  ```
  > **Note:** The `make install` target has only been tested on Linux. Windows users are welcome to try it and report back any issues. macOS users... well, you know the drill.

After compilation, you will find the executable(s) in the `build/` directory.

## Configuring

`GoDMX` is configured via a JSON file. By default, `GoDMX` looks for `config.json` in the current working directory.

### Chain and Effect Structure

At the core of `GoDMX` are **Chains**. A chain represents a sequence of DMX lamps that are processed together. Each chain has its own set of effects and an output configuration. This allows for modular and scalable lighting setups.

Key properties of a chain:
- `id`: A unique identifier for the chain.
- `priority`: Determines the order of processing if multiple chains are active.
- `tickRate`: How often the chain's effects are processed per second.
- `numLamps`: The total number of DMX lamps in this chain.
- `effects`: An array of effects applied to the lamps in this chain.
- `output`: Defines how the processed DMX data is sent out (e.g., ArtNet, Govee).

Effects are the building blocks of your lighting animations. They are configured within the `effects` array of each chain in your configuration file.

Each effect has the following common properties:
- `id`: A unique identifier for the effect within its chain.
- `type`: The type of effect (e.g., `rainbow`, `solidColor`, `shift`, `hueshift`).
- `args`: A JSON object containing parameters specific to the effect. These parameters control the behavior and appearance of the effect.
- `enabled`: A boolean indicating whether the effect is currently active.
- `group`: (Optional) A string to group related effects. If an effect has a `group` defined, only one effect within that group can be enabled at any given time. Enabling an effect in a group will automatically disable all other effects in the same group. This is enforced by the `EnforceGroupRules` logic within the orchestrator.

### Example Effect Configuration

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

## Triggers and Actions

`GoDMX` allows you to define custom **Events** that can be triggered by various sources (like MIDI messages or the Web UI). Each event consists of one or more **Actions** that `GoDMX` will perform when the event is triggered.

### Event Structure

Events are defined in the `events` section of your configuration file as a map, where the key is the event name (string) and the value is an array of actions.

Example `events` configuration:

```json
"events": {
  "strobe_on": [
    {
      "type": "add_effect",
      "chain_id": "mainChain",
      "params": {
        "id": "strobeEffect",
        "type": "blink",
        "enabled": true,
        "args": { "divider": 4, "dutyCycle": 0.1 }
      }
    }
  ],
  "strobe_off": [
    {
      "type": "remove_effect",
      "chain_id": "mainChain",
      "effect_id": "strobeEffect"
    }
  ],
  "rainbow_on": [
    {
      "type": "toggle_effect",
      "chain_id": "mainChain",
      "effect_id": "defaultRainbow",
      "params": { "enabled": true }
    }
  ],
  "set_fast_bpm": [
    {
      "type": "set_global",
      "params": { "bpm": 180.0 }
    }
  ]
}
```

### Action Types

Each action within an event has a `type` property, which determines what `GoDMX` will do. Common action types include:

*   `"add_effect"`: Adds a new effect to a specified chain. Requires `chain_id` and `params` (which should contain the full effect configuration, including `id`, `type`, `args`, `enabled`, and `group`).
*   `"remove_effect"`: Removes an effect from a specified chain. Requires `chain_id` and `effect_id`.
*   `"toggle_effect"`: Toggles the `enabled` state of an existing effect in a chain. Requires `chain_id`, `effect_id`, and `params` (with an `enabled` boolean, e.g., `{"enabled": true}`).
*   `"set_global"`: Sets a global parameter (like `bpm`, `color1`, `color2`, `intensity`). Requires `params` with the global setting(s) to change.

## MIDI Configuration

`GoDMX` can listen for MIDI messages to trigger events defined in your configuration.

### `midi_port_name`

This field specifies the name of the MIDI input port `GoDMX` should listen to. You can find available MIDI port names on your system using tools like `aconnect -l` (Linux) or by checking your DAW/MIDI software.

Example:
```json
"midi_port_name": "Midi Through Port-0"
```

### `midi_triggers`

This is an array of objects that define which incoming MIDI messages should trigger specific `GoDMX` events. Each trigger has the following properties:

*   `message_type` (string): The type of MIDI message to listen for.
    *   `"cc"`: Control Change message.
    *   `"note_on"`: Note On message.
    *   `"note_off"`: Note Off message.
*   `number` (integer): The MIDI control change number (0-127) for `cc` messages, or the MIDI note number (0-127) for `note_on`/`note_off` messages.
*   `value` (integer): The value of the MIDI message (0-127). For `cc` messages, this is the CC value. For `note_on` messages, this is the velocity. Use `-1` to match any value.
*   `event_name` (string): The name of the event (as defined in the `events` section of your configuration) to trigger when this MIDI message is received.

Example `midi_triggers` configuration:

```json
"midi_triggers": [
  {
    "message_type": "cc",
    "number": 1,
    "value": -1,
    "event_name": "strobe_on"
  },
  {
    "message_type": "cc",
    "number": 2,
    "value": -1,
    "event_name": "strobe_off"
  },
  {
    "message_type": "note_on",
    "number": 60,
    "value": -1,
    "event_name": "rainbow_on"
  }
]
```

## Web UI

`GoDMX` includes a simple web-based user interface for monitoring and controlling your lighting setup.

### Accessing the Web UI

The web UI runs on a specified port, which defaults to `8080`. You can change this port using the `-web-port` command-line flag.

To access the web UI, open your web browser and navigate to:

`http://localhost:8080` (or your chosen port)

### Features

The web UI provides:

*   **Real-time Monitoring:** View the status of your configured chains and effects.
*   **BPM Control:** Adjust the global BPM.
*   **Event Triggering:** Manually trigger any defined events.

The web UI is served from the `web/` directory in the project.

## Command-Line Flags

`GoDMX` supports the following command-line flags:

*   `-debug`: Run in debug mode (exits after a short duration).
*   `-config <path>`: Path to the configuration file. By default, `GoDMX` looks for `config.json` in the current working directory.
*   `-web-port <port>`: Port for the web UI (default: `8080`).
*   `-event <name>`: Name of an event to trigger on startup.
*   `-docs`: Generate documentation for effects in `EFFECTS.md`.

### Workflow and Examples

`GoDMX` is designed to be run from your project directory. By default, it will look for a `config.json` file in the current directory.

**Typical Workflow:**

1.  `cd /path/to/your/GoDMX/project`
2.  `./build/godmx_linux` (or `godmx_macos`, `godmx_windows.exe`)

If you have multiple setups or wish to use a different configuration file, use the `-config` flag:

```bash
./build/godmx_linux -config path/to/your/custom_config.json
```

(Replace `godmx_linux` with `godmx_macos` or `godmx_windows.exe` as appropriate for your platform.)