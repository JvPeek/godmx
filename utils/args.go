package utils

// GetFloatArg safely retrieves a float64 argument from a map, with a default value.
func GetFloatArg(args map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := args[key]; ok {
		if fVal, isFloat := val.(float64); isFloat {
			return fVal
		}
		// Handle int to float conversion if necessary
		if iVal, isInt := val.(int); isInt {
			return float64(iVal)
		}
	}
	return defaultValue
}

// GetIntArg safely retrieves an int argument from a map, with a default value.
func GetIntArg(args map[string]interface{}, key string, defaultValue int) int {
	if val, ok := args[key]; ok {
		if iVal, isInt := val.(int); isInt {
			return iVal
		}
		// Handle float to int conversion (truncation) if necessary
		if fVal, isFloat := val.(float64); isFloat {
			return int(fVal)
		}
	}
	return defaultValue
}
