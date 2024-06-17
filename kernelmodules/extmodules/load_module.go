package extmodules

import (
	"fmt"

	cgowrapper "github.com/CustodiaJS/custodiajs-core/kernelmodules/extmodules/cgo_wrapper"
)

// Lädt ein Libmodule
func LoadModuleLib(pathv string) (*ExternalModule, error) {
	// Es wird mittels CGO versucht das LibModule zu laden
	result, err := cgowrapper.LoadWrappedCGOLibModule(pathv)
	if err != nil {
		return nil, fmt.Errorf("LoadModuleLib: " + err.Error())
	}

	// Das Rückgabeobjekt wird erstellt
	returnValue := &ExternalModule{result}

	// Rückgabe
	return returnValue, nil
}
