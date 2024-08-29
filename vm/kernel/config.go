package kernel

import "github.com/CustodiaJS/custodiajs-core/types"

func NewFromExist(existKernelConfig *KernelConfig, extModues ...types.KernelModuleInterface) *KernelConfig {
	// Es wird eine neue Module Liste Erstellt
	modList := make([]types.KernelModuleInterface, 0)
	modList = append(modList, existKernelConfig.Modules...)
	modList = append(modList, extModues...)

	// Die Liste mit neuen Items wird zurückgegebn
	return &KernelConfig{Modules: modList}
}
