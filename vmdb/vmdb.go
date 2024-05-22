package vmdb

import (
	"fmt"
	"strings"
	"vnh1/utils"
)

func (o *VmDatabase) loadAllVirtualMachines() error {
	// Die VM's werden geladen
	vms, err := utils.ScanVmDir(o.vmRootDir)
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

func (o *VmDatabase) GetAllVirtualMachines() []*VmDBEntry {
	v := make([]*VmDBEntry, 0)
	for _, item := range o.vmMap {
		v = append(v, item)
	}
	return v
}

func (o *VmDatabase) GetAllDatabaseConfigurations() []*VMDatabaseData {
	v := map[string]*VMDatabaseData{}
	for _, item := range o.vmMap {
		for _, xtem := range item.GetAllDatabaseServices() {
			if _, found := v[string(xtem.GetDatabaseFingerprint())]; !found {
				v[string(xtem.GetDatabaseFingerprint())] = xtem
			}
		}
	}

	r := make([]*VMDatabaseData, 0)
	for _, item := range v {
		r = append(r, item)
	}

	return r
}

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
