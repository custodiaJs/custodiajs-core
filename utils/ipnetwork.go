package utils

import (
	"net"
)

// Funktion zum Überprüfen, ob eine IP-Adresse eine private IP ist
func IsPrivateIP(ip net.IP) bool {
	// IPv4 private address ranges (RFC 1918):
	// - 10.0.0.0/8
	// - 172.16.0.0/12
	// - 192.168.0.0/16
	private := []net.IPNet{
		// 10.0.0.0/8
		{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)},
		// 172.16.0.0/12
		{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)},
		// 192.168.0.0/16
		{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)},
	}

	for _, network := range private {
		if network.Contains(ip) {
			return true
		}
	}

	// IPv6 unique local addresses (RFC 4193):
	// - fc00::/7
	// Note: There are other IPv6 private ranges, but fc00::/7 is the unique local address range.
	// For simplicity, we check only this range here.
	// Check if it's IPv6 and is in the fc00::/7 range
	if ip.To4() == nil && (ip[0]&0xfe == 0xfc) {
		return true
	}

	return false
}
