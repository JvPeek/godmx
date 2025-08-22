package midi

import (
	"fmt"
	"log"

	"godmx/orchestrator"
	"godmx/config"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // Required for MIDI driver
)

// MidiController manages MIDI input and triggers orchestrator events.
type MidiController struct {
	orch *orchestrator.Orchestrator
	triggers []config.MidiTriggerConfig
	stopListen func()
	midiPortName string
}

// NewMidiController creates a new MidiController.
func NewMidiController(orch *orchestrator.Orchestrator, triggers []config.MidiTriggerConfig, midiPortName string) (*MidiController, error) {
	return &MidiController{
		orch: orch,
		triggers: triggers,
		midiPortName: midiPortName,
	},
	nil
}

// Start initializes MIDI and begins listening for messages.
func (mc *MidiController) Start() error {
	// Find an input port by name
	in, err := midi.FindInPort(mc.midiPortName)
	if err != nil {
		return fmt.Errorf("can't find MIDI input port %s: %w", mc.midiPortName, err)
	}
	log.Printf("Found MIDI input device: %s\n", in.String())


	// Listen for MIDI messages
	mc.stopListen, err = midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var bt []byte
		var ch, key, vel uint8
		switch {
		case msg.GetSysEx(&bt):
			log.Printf("got sysex: % X\n", bt)
		case msg.GetNoteStart(&ch, &key, &vel):
			log.Printf("starting note %s on channel %v with velocity %v\n", midi.Note(key), ch, vel)
			mc.matchAndTrigger("note_on", int64(key), int64(vel))
		case msg.GetNoteEnd(&ch, &key):
			log.Printf("ending note %s on channel %v\n", midi.Note(key), ch)
			mc.matchAndTrigger("note_off", int64(key), 0) // Velocity is 0 for note off
		case msg.GetControlChange(&ch, &key, &vel):
			log.Printf("MIDI CC: Controller=%d, Value=%d\n", key, vel)
			mc.matchAndTrigger("cc", int64(key), int64(vel))
		default:
			log.Printf("Unhandled MIDI event: % X\n", msg.Bytes())
		}
	}, midi.UseSysEx()) // UseSysEx to enable SysEx messages

	if err != nil {
		return fmt.Errorf("failed to start MIDI listener: %w", err)
	}

	log.Println("Listening for MIDI messages...")

	return nil
}

// matchAndTrigger checks if a MIDI event matches any configured trigger and triggers the event.
func (mc *MidiController) matchAndTrigger(messageType string, number int64, value int64) {
	for _, trigger := range mc.triggers {
		if trigger.MessageType == messageType &&
			trigger.Number == int(number) &&
			(trigger.Value == -1 || trigger.Value == int(value)) {
			log.Printf("MIDI Trigger matched: Type=%s, Number=%d, Value=%d -> Triggering event '%s'\n",
				messageType, number, value, trigger.EventName)
			mc.orch.TriggerEvent(trigger.EventName)
			return // Trigger only the first matching event
		}
	}
}

// Stop terminates the MIDI input stream.
func (mc *MidiController) Stop() {
	if mc.stopListen != nil {
		mc.stopListen()
	}
	midi.CloseDriver()
}
