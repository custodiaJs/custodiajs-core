package localgrpc

import (
	"vnh1/types"
)

func (o *HostCliService) SetupCore(coreObj types.CoreInterface) error {
	// Der Core wird abgespeichert
	o.core = coreObj

	// Der Vorgang ist ohne fehler durchgeführt wurden
	return nil
}
