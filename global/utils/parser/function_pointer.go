package parser

import (
	"fmt"
	"regexp"
	"strings"
)

func ParseFunctionSignature(signature string) (funcName string, argTypes []string, err error) {
	// Regex to match the function name and the arguments inside the parentheses
	re := regexp.MustCompile(`^(\w+)\(([\w\s,]*)\)$`)
	matches := re.FindStringSubmatch(signature)
	if matches == nil || len(matches) != 3 {
		return "", nil, fmt.Errorf("invalid function signature")
	}

	// Extract the function name
	funcName = matches[1]

	// Extract argument types and handle the case of no arguments
	argString := matches[2]
	if argString != "" {
		argTypes = strings.Split(argString, ",")
		for i, arg := range argTypes {
			argTypes[i] = strings.TrimSpace(arg)
		}
	}

	// Rückgabe
	return funcName, argTypes, nil
}
