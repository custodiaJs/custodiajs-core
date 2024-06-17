package utils

import (
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/CustodiaJS/custodiajs-core/utils/parser"
)

func ValidateDatatypeString(dType string) bool {
	switch strings.ToLower(dType) {
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

func ValidateContainerHexID(hxId string) bool {
	return true
}

func ValidateFunctionSignature(funcPtr string) bool {
	// Es wird versucht die FunctionsSignatur einzulesen
	fname, argTypes, err := parser.ParseFunctionSignature(funcPtr)
	if err != nil {
		return false
	}

	// Es wird geprüft ob es sich um einen Zulässigen Funktionsnamen handelt
	if !ValidateFunctionName(fname) {
		return false
	}

	// Es werden alle Parameter geprüft, es dürfen nur Datentypen angegeben werden
	for _, item := range argTypes {
		if item != "string" && item != "array" && item != "object" && item != "number" && item != "bool" {
			return false
		}
	}

	// Es handelt sich um einen Zulässigen Funktionspointer
	return true
}

func CheckHostInWhitelist(host string, whitelist []string) bool {
	for _, allowed := range whitelist {
		if strings.HasPrefix(allowed, "*.") {
			wildcardDomain := strings.TrimPrefix(allowed, "*.")
			if strings.HasSuffix(host, "."+wildcardDomain) {
				return true
			}
		} else if allowed == host {
			return true
		}
	}
	return false
}
