package effects

import (
	"fmt"
	"godmx/orchestrator"
)

// EffectConstructor is a function type that constructs an Effect from a map of arguments.
type EffectConstructor func(args map[string]interface{}) (orchestrator.Effect, error)

var (
	effectRegistry = make(map[string]EffectConstructor)
	effectParameterRegistry = make(map[string]map[string]interface{})
)

// RegisterEffect registers an effect constructor with a given name.
func RegisterEffect(name string, constructor EffectConstructor) {
	if _, exists := effectRegistry[name]; exists {
		panic(fmt.Sprintf("Effect '%s' already registered", name))
	}
	effectRegistry[name] = constructor
}

// RegisterEffectParameters registers the default parameters for an effect.
func RegisterEffectParameters(name string, params map[string]interface{}) {
	if _, exists := effectParameterRegistry[name]; exists {
		panic(fmt.Sprintf("Effect parameters for '%s' already registered", name))
	}
	effectParameterRegistry[name] = params
}

// GetEffectConstructor retrieves an effect constructor by name.
func GetEffectConstructor(name string) (EffectConstructor, bool) {
	constructor, ok := effectRegistry[name]
	return constructor, ok
}

// GetEffectParameters retrieves the default parameters for an effect by name.
func GetEffectParameters(name string) (map[string]interface{}, bool) {
	params, ok := effectParameterRegistry[name]
	return params, ok
}

// GetAvailableEffects returns a slice of all registered effect names.
func GetAvailableEffects() []string {
	names := make([]string, 0, len(effectRegistry))
	for name := range effectRegistry {
		names = append(names, name)
	}
	return names
}
