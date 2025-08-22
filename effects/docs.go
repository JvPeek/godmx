package effects

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// GenerateEffectDocumentation generates a Markdown file with a table of all registered effects.
func GenerateEffectDocumentation() error {
	var builder strings.Builder

	builder.WriteString("# Registered Effects\n\n")
	builder.WriteString("This document outlines all available lighting effects, their descriptions, and configurable parameters.\n\n")

	names := GetAvailableEffects()
	sort.Strings(names)

	for _, name := range names {
		metadata, ok := GetEffectMetadata(name)
		if !ok {
			continue // Should not happen if GetAvailableEffects is accurate
		}

		// Effect Heading
		builder.WriteString(fmt.Sprintf("## %s\n\n", metadata.HumanReadableName))

		// Description
		builder.WriteString(fmt.Sprintf("%s\n\n", metadata.Description))

		// Tags
		if len(metadata.Tags) > 0 {
			sort.Strings(metadata.Tags)
			builder.WriteString(fmt.Sprintf("**Tags**: %s\n\n", strings.Join(metadata.Tags, ", ")))
		}

		// Parameters Table
		if len(metadata.Parameters) > 0 {
			builder.WriteString("### Parameters\n\n")
			builder.WriteString("| Parameter Name | Display Name | Data Type | Default Value | Min Value | Max Value | Description |\n")
			builder.WriteString("|----------------|--------------|-----------|---------------|-----------|-----------|-------------|\n")

			// Sort parameters by DisplayName for consistent output
			sort.Slice(metadata.Parameters, func(i, j int) bool {
				return metadata.Parameters[i].DisplayName < metadata.Parameters[j].DisplayName
			})

			for _, p := range metadata.Parameters {
				minVal := fmt.Sprintf("%v", p.MinValue)
				if p.MinValue == nil {
					minVal = "-"
				}
				maxVal := fmt.Sprintf("%v", p.MaxValue)
				if p.MaxValue == nil {
					maxVal = "-"
				}

				builder.WriteString(fmt.Sprintf("| %s | %s | %s | %v | %s | %s | %s |\n",
					p.InternalName, p.DisplayName, p.DataType, p.DefaultValue, minVal, maxVal, p.Description))
			}
			builder.WriteString("\n") // Add a newline after the table for spacing
		}
		builder.WriteString("---\n\n") // Separator between effects
	}

	return os.WriteFile("EFFECTS.md", []byte(builder.String()), 0644)
}

