// Author: fluffelpuff
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package ipnetwork

import (
	"net"
)

func NewHostNetworkManagmentUnit() *HostNetworkManagmentUnit {
	// Das HostNetworkManagmentUnit wird erzeugt
	hnmu := &HostNetworkManagmentUnit{}

	// Es wird eine neue Routine erzeugt, welche Automatisch die Lokalen IP Adressen einließt
	go routineWatcherForLocalIpAddresses(hnmu)

	// Es wird eine neue Routine erzeugt, welche Automatisch die DHCP Adressen überwach
	go routineWatcherForDHCPAddresses(hnmu)

	// Es wird eine neue Routine erzeugt, diese Routine wird verwendet um die IP-Adressen der Tor Exit Nodes zu ermitteln
	go routineWatcherForTorExitNodes(hnmu)

	// Das Objekt wird zurückgegeben
	return hnmu
}

func routineWatcherForLocalIpAddresses(hnmu *HostNetworkManagmentUnit) {

}

func routineWatcherForDHCPAddresses(hnmu *HostNetworkManagmentUnit) {

}

func routineWatcherForTorExitNodes(hnmu *HostNetworkManagmentUnit) {

}

func (o *HostNetworkManagmentUnit) isLoclhostIp(rIpAdr net.IP) bool {
	return false
}

// Wird verwendet um das Aktuelle Netzwerk Interface anhand der IP-Adresse zu ermitteln
func (o *HostNetworkManagmentUnit) GetNetworkInterfaceByLocalIp(address *IpAddress) *NetworkInterface {
	return nil
}

// Wird verwendet um eine IP-Adresse einzulesen
func (o *HostNetworkManagmentUnit) TryParseIp(ipaddr string) (*IpAddress, error) {
	// Es wird mittels "go:net" versucht die IP-Adresse einzulesen
	rIpAdr := net.ParseIP(ipaddr)
	if rIpAdr == nil {

	}

	// Es wird geprüft ob es sich um eine Lokale Adresse handelt
	addressIsLocalhostAddress := o.isLoclhostIp(rIpAdr)
	_ = addressIsLocalhostAddress

	// Es wird geprüft ob es sich um ein Privates Subbnet handelt,
	// wenn nicht wird geprüft ob es sich bei der Adresse um ein Tor Exit Node handelt.

	return nil, nil
}
