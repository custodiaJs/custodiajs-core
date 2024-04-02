package utils

import (
	"strings"
)

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
