package utils

import (
	"fmt"
	"regexp"
	"strings"
	"vnh1/types"
)

// parseFunctionSignature parses a function signature string and returns a FunctionSignature struct
func ParseFunctionSignature(input string) (*types.FunctionSignature, error) {
	signature := &types.FunctionSignature{}

	// Regex patterns for different parts of the input
	vmIDRegex := regexp.MustCompile(`vmid{([^}]*)}`)
	vmNameRegex := regexp.MustCompile(`vmname{([^}]*)}`)
	funcRegex := regexp.MustCompile(`function (\w+)\((.*?)\)(?: (\w+))?`)

	// Extract VM ID or VM name if present
	vmIDMatch := vmIDRegex.FindStringSubmatch(input)
	if len(vmIDMatch) > 1 {
		signature.VMID = vmIDMatch[1]
	}

	vmNameMatch := vmNameRegex.FindStringSubmatch(input)
	if len(vmNameMatch) > 1 {
		signature.VMName = vmNameMatch[1]
	}

	// Extract function name, parameters, and return type
	funcMatch := funcRegex.FindStringSubmatch(input)
	if len(funcMatch) > 1 {
		signature.FunctionName = funcMatch[1]
		if funcMatch[2] != "" {
			signature.Params = strings.Split(funcMatch[2], ", ")
		}
		if len(funcMatch) > 3 {
			signature.ReturnType = funcMatch[3]
		}
	} else {
		// Handle cases without vmid or vmname prefix
		funcRegexNoVM := regexp.MustCompile(`^function (\w+)\((.*?)\)(?: (\w+))?`)
		funcMatchNoVM := funcRegexNoVM.FindStringSubmatch(input)
		if len(funcMatchNoVM) > 1 {
			signature.FunctionName = funcMatchNoVM[1]
			if funcMatchNoVM[2] != "" {
				signature.Params = strings.Split(funcMatchNoVM[2], ", ")
			}
			if len(funcMatchNoVM) > 3 {
				signature.ReturnType = funcMatchNoVM[3]
			}
		}
	}

	return signature, nil
}

func ParseFunctionSignatureOptionalFunction(input string) (*types.FunctionSignature, error) {
	signature := &types.FunctionSignature{}

	// Regex patterns for different parts of the input
	vmIDRegex := regexp.MustCompile(`vmid{([^}]*)}`)
	vmNameRegex := regexp.MustCompile(`vmname{([^}]*)}`)
	funcRegex := regexp.MustCompile(`(?:function )?(\w+)\((.*?)\)(?: (\w+))?`)

	// Extract VM ID or VM name if present
	vmIDMatch := vmIDRegex.FindStringSubmatch(input)
	if len(vmIDMatch) > 1 {
		signature.VMID = vmIDMatch[1]
	}

	vmNameMatch := vmNameRegex.FindStringSubmatch(input)
	if len(vmNameMatch) > 1 {
		signature.VMName = vmNameMatch[1]
	}

	// Extract function name, parameters, and return type
	funcMatch := funcRegex.FindStringSubmatch(input)
	if len(funcMatch) > 1 {
		signature.FunctionName = funcMatch[1]
		if funcMatch[2] != "" {
			signature.Params = strings.Split(funcMatch[2], ", ")
		}
		if len(funcMatch) > 3 {
			signature.ReturnType = funcMatch[3]
		}
	} else {
		// Handle cases without vmid or vmname prefix
		funcRegexNoVM := regexp.MustCompile(`^function (\w+)\((.*?)\)(?: (\w+))?`)
		funcMatchNoVM := funcRegexNoVM.FindStringSubmatch(input)
		if len(funcMatchNoVM) > 1 {
			signature.FunctionName = funcMatchNoVM[1]
			if funcMatchNoVM[2] != "" {
				signature.Params = strings.Split(funcMatchNoVM[2], ", ")
			}
			if len(funcMatchNoVM) > 3 {
				signature.ReturnType = funcMatchNoVM[3]
			}
		}
	}

	return signature, nil
}

// String generates a formatted string representation of FunctionDetails.
func String(fd *types.FunctionSignature) string {
	var sb strings.Builder
	if fd.VMID != "" && fd.VMName == "" {
		sb.WriteString(fmt.Sprintf("vmid{%s} -> ", fd.VMID))
	} else if fd.VMName != "" && fd.VMID == "" {
		sb.WriteString(fmt.Sprintf("vmname{%s} -> ", fd.VMName))
	}
	sb.WriteString(fmt.Sprintf("function %s(", fd.FunctionName))
	if len(fd.Params) > 0 {
		sb.WriteString(strings.Join(fd.Params, ", "))
	}
	sb.WriteString(")")
	if fd.ReturnType != "" {
		sb.WriteString(fmt.Sprintf(" %s", fd.ReturnType))
	}
	return sb.String()
}

func FunctionOnlySignatureString(fd *types.FunctionSignature) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("function %s(", fd.FunctionName))
	if len(fd.Params) > 0 {
		sb.WriteString(strings.Join(fd.Params, ", "))
	}
	sb.WriteString(")")
	if fd.ReturnType != "" {
		sb.WriteString(fmt.Sprintf(" %s", fd.ReturnType))
	}
	return sb.String()
}
