package httpapi

import (
	"fmt"
	"vnh1/types"
)

func (o *HttpApiService) SetupCore(coreObj types.CoreInterface) error {
	// Es wird geprüft ob der Core festgelegt wurde
	if o.core != nil {
		return fmt.Errorf("SetupCore: always linked with core")
	}

	// Das Objekt wird zwischengespeichert
	o.core = coreObj

	// Der Vorgang ist ohne fehler durchgeführt wurden
	return nil
}
