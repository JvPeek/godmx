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
	builder.WriteString("| Effect Name | Description | Parameters | Tags |\n")
	builder.WriteString("|-------------|-------------|------------|------|\n")

	names := GetAvailableEffects()
	sort.Strings(names)

	for _, name := range names {
		metadata, ok := GetEffectMetadata(name)
		if !ok {
			continue // Should not happen if GetAvailableEffects is accurate
		}

		paramStr := ""
		if len(metadata.Parameters) > 0 {
			paramNames := make([]string, 0, len(metadata.Parameters))
			for _, p := range metadata.Parameters {
				paramNames = append(paramNames, fmt.Sprintf("%s (%s, default: %v)", p.DisplayName, p.DataType, p.DefaultValue))
			}
			sort.Strings(paramNames)
			paramStr = strings.Join(paramNames, ", ")
		}

		tagStr := ""
		if len(metadata.Tags) > 0 {
			sort.Strings(metadata.Tags)
			tagStr = strings.Join(metadata.Tags, ", ")
		}

		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", metadata.HumanReadableName, metadata.Description, paramStr, tagStr))
	}

	return os.WriteFile("EFFECTS.md", []byte(builder.String()), 0644)
}

