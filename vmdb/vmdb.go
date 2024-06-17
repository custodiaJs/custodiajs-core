package vmdb

import (
	"fmt"
	"strings"

	"github.com/CustodiaJS/custodiajs-core/filesystem"
)

// loadAllVirtualMachines lädt alle virtuellen Maschinen aus dem angegebenen Verzeichnis und speichert sie im VmDatabase-Objekt.
// Falls ein Fehler auftritt, wird eine entsprechende Fehlermeldung zurückgegeben.
func (o *VmDatabase) loadAllVirtualMachines() error {
	// Die VM's werden geladen
	vms, err := filesystem.ScanVmDir(o.vmRootDir)
	if err != nil {
		return fmt.Errorf("VmDatabase->LoadAllVirtualMachines: " + err.Error())
	}

	// Es wird versucht die einzelenen VMS zu Laden
	for _, item := range vms {
		// Es wird versucht die VM zu laden
		vm, err := tryToLoadVM(item)
		if err != nil {
			return fmt.Errorf("VmDatabase->LoadAllVirtualMachines: " + err.Error())
		}

		// Die VM wird zwischengespeichert
		o.vmMap[strings.ToLower(vm.GetVMContainerMerkleHash())] = vm
	}

	// Log
	fmt.Printf("Total VM's loaded: %d\n", len(o.vmMap))

	// Die Virtuellen Machinen werden zurückgegeben
	return nil
}

// GetAllVirtualMachines gibt alle geladenen virtuellen Maschinen als Slice von VmDBEntry-Objekten zurück.
func (o *VmDatabase) GetAllVirtualMachines() []*VmDBEntry {
	v := make([]*VmDBEntry, 0)
	for _, item := range o.vmMap {
		v = append(v, item)
	}
	return v
}

// GetAllDatabaseVMBaseData gibt die Basisdaten aller virtuellen Maschinen zurück, die in der Datenbank gespeichert sind.
func (o *VmDatabase) GetAllDatabaseVMBaseData() []*VMEntryBaseData {
	v := map[string]*VMEntryBaseData{}
	for _, item := range o.vmMap {
		for _, xtem := range item.GetAllDatabaseServices() {
			if _, found := v[string(xtem.GetDatabaseFingerprint())]; !found {
				v[string(xtem.GetDatabaseFingerprint())] = xtem
			}
		}
	}

	r := make([]*VMEntryBaseData, 0)
	for _, item := range v {
		r = append(r, item)
	}

	return r
}

// OpenFilebasedVmDatabase öffnet eine Dateibasierte VM-Datenbank und gibt ein VmDatabase-Objekt zurück.
// Dabei werden alle virtuellen Maschinen geladen und zwischengespeichert.
// Falls ein Fehler auftritt, wird eine entsprechende Fehlermeldung zurückgegeben.
func OpenFilebasedVmDatabase() (*VmDatabase, error) {
	// Das VM Datenbankobjekt wird erstellt
	resolv := &VmDatabase{
		vmMap:     map[string]*VmDBEntry{},
		vmRootDir: "/var/lib/vnh1",
	}

	// Es werden alle Virtuellen Machinen geladen und zwischengespeichert
	if err := resolv.loadAllVirtualMachines(); err != nil {
		return nil, fmt.Errorf("VmDatabase->OpenFilebasedVmDatabase: " + err.Error())
	}

	// Das Objekt wird zurückgegeben
	return resolv, nil
}
