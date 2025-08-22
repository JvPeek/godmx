package effects

import (
	"fmt"
	"godmx/types"
)

// EffectConstructor is a function type that constructs an Effect from a map of arguments.
type EffectConstructor func(args map[string]interface{}) (types.Effect, error)



var (
	effectRegistry = make(map[string]EffectConstructor)
	effectMetadataRegistry = make(map[string]types.EffectMetadata)
)

// RegisterEffect registers an effect constructor with a given name.
func RegisterEffect(name string, constructor EffectConstructor) {
	if _, exists := effectRegistry[name]; exists {
		panic(fmt.Sprintf("Effect '%s' already registered", name))
	}
	effectRegistry[name] = constructor
}

// RegisterEffectMetadata registers comprehensive metadata for an effect.
func RegisterEffectMetadata(name string, metadata types.EffectMetadata) {
	if _, exists := effectMetadataRegistry[name]; exists {
		panic(fmt.Sprintf("Effect metadata for '%s' already registered", name))
	}
	effectMetadataRegistry[name] = metadata
}

// GetEffectConstructor retrieves an effect constructor by name.
func GetEffectConstructor(name string) (EffectConstructor, bool) {
	constructor, ok := effectRegistry[name]
	return constructor, ok
}

// GetEffectMetadata retrieves comprehensive metadata for an effect by name.
func GetEffectMetadata(name string) (types.EffectMetadata, bool) {
	metadata, ok := effectMetadataRegistry[name]
	return metadata, ok
}

// GetAvailableEffects returns a slice of all registered effect names.
func GetAvailableEffects() []string {
	names := make([]string, 0, len(effectRegistry))
	for name := range effectRegistry {
		names = append(names, name)
	}
	return names
}

// Helper functions for backward compatibility or specific needs

// GetEffectParameters retrieves the default parameters for an effect by name from its metadata.
func GetEffectParameters(name string) (map[string]interface{}, bool) {
	metadata, ok := GetEffectMetadata(name)
	if !ok {
		return nil, false
	}
	params := make(map[string]interface{})
	for _, p := range metadata.Parameters {
		params[p.InternalName] = p.DefaultValue
	}
	return params, true
}

// GetEffectTags retrieves the tags for an effect by name from its metadata.
func GetEffectTags(name string) ([]string, bool) {
	metadata, ok := GetEffectMetadata(name)
	if !ok {
		return nil, false
	}
	return metadata.Tags, true
}

// GetEffectDescription retrieves the description for an effect by name from its metadata.
func GetEffectDescription(name string) (string, bool) {
	metadata, ok := GetEffectMetadata(name)
	if !ok {
		return "", false
	}
	return metadata.Description, true
}