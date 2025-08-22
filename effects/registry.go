package effects

import (
	"fmt"
	"godmx/orchestrator"
)

// EffectConstructor is a function type that constructs an Effect from a map of arguments.
type EffectConstructor func(args map[string]interface{}) (orchestrator.Effect, map[string]interface{}, error)

var (
	effectRegistry = make(map[string]EffectConstructor)
)

// RegisterEffect registers an effect constructor with a given name.
func RegisterEffect(name string, constructor EffectConstructor) {
	if _, exists := effectRegistry[name]; exists {
		panic(fmt.Sprintf("Effect '%s' already registered", name))
	}
	effectRegistry[name] = constructor
}

// GetEffectConstructor retrieves an effect constructor by name.
func GetEffectConstructor(name string) (EffectConstructor, bool) {
	constructor, ok := effectRegistry[name]
	return constructor, ok
}

// GetAvailableEffects returns a slice of all registered effect names.
func GetAvailableEffects() []string {
	names := make([]string, 0, len(effectRegistry))
	for name := range effectRegistry {
		names = append(names, name)
	}
	return names
}
