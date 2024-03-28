package utils

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
