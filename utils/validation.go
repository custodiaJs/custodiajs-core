package utils

import (
	"encoding/hex"
	"regexp"
	"strings"
)

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

func ValidateVMName(vmNameString string) bool {
	return true
}

func ValidateFunctionName(funcName string) bool {
	// Regulärer Ausdruck, um zu überprüfen, ob der String den Variablennamenkriterien entspricht
	validVariable := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

	// Überprüfe, ob der String den Variablennamenkriterien entspricht
	return validVariable.MatchString(funcName)
}

func ValidateVMIdString(idString string) bool {
	if len(strings.ToLower(idString)) != 64 {
		return false
	}

	v, err := hex.DecodeString(idString)
	if err != nil {
		return false
	}

	if len(v) != 32 {
		return false
	}

	x := hex.EncodeToString(v)
	return strings.EqualFold(idString, x)
}

func ValidateExternalModuleName(funcName string) bool {
	// Regulärer Ausdruck, um zu überprüfen, ob der String den Variablennamenkriterien entspricht
	validVariable := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

	// Überprüfe, ob der String den Variablennamenkriterien entspricht
	return validVariable.MatchString(funcName)
}
