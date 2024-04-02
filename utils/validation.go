package utils

import "regexp"

func ValidateDatatypeString(dType string) bool {
	switch dType {
	case "boolean":
		return true
	case "number":
		return true
	case "string":
		return true
	case "array":
		return true
	case "object":
		return true
	default:
		return false
	}
}

func ValidateFunctionName(funcName string) bool {
	// Regulärer Ausdruck, um zu überprüfen, ob der String den Variablennamenkriterien entspricht
	validVariable := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

	// Überprüfe, ob der String den Variablennamenkriterien entspricht
	return validVariable.MatchString(funcName)
}
