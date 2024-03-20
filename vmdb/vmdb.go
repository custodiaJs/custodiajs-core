package vmdb

import (
	"fmt"
	"vnh1/static"
)

type VmDatabase struct {
	vmMap     map[string]*VmDBEntry
	vmRootDir string
}

func (o *VmDatabase) LoadAllVirtualMachines() ([]*VmDBEntry, error) {
	// Die VM's werden geladen
	vms, err := static.ScanVmDir(o.vmRootDir)
	if err != nil {
		return nil, fmt.Errorf("LoadAllVirtualMachines: " + err.Error())
	}

	// Es wird versucht die einzelenen VMS zu Laden
	retrivedVms := []*VmDBEntry{}
	for _, item := range vms {
		// Es wird versucht die VM zu laden
		vm, err := tryToLoadVM(item)
		if err != nil {
			return nil, fmt.Errorf("LoadAllVirtualMachines: " + err.Error())
		}

		// Die VM wird zwischengespeichert
		retrivedVms = append(retrivedVms, vm)
	}

	// Die Virtuellen Machinen werden zurückgegeben
	return retrivedVms, nil
}

func OpenFilebasedVmDatabase() (*VmDatabase, error) {
	return &VmDatabase{vmMap: map[string]*VmDBEntry{}, vmRootDir: "/var/lib/vnh1"}, nil
}